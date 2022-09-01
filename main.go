package main

import (
	"encoding/json"
	"flag"
	"github.com/Hintay/region_restriction_check-go/exporter"
	"github.com/Hintay/region_restriction_check-go/medias"
	"github.com/Hintay/region_restriction_check-go/medias/model"
	"io/ioutil"
	"os"
	"runtime"
	"sort"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	version     = "0.0.11"
	modeChecker = "checker"
	modeMonitor = "monitor"
)

var (
	flags struct {
		ConfigFile string `json:"-"`
		Version    bool   `json:"-"`
		DNS        string `json:"dns"`
		SOCKS      string `json:"socks_proxy"`
		// Checker
		Mode    string      `json:"-"`
		Regions regionArray `json:"-"`

		// Monitor
		Interval       int             `json:"interval"`
		ExporterListen string          `json:"exporter_listen"`
		Log            Log             `json:"log"`
		Tasks          map[string]Task `json:"tasks"`
	}
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.StringVar(&flags.ConfigFile, "config.file", "config.json", "only working with monitor mode")
	flag.StringVar(&flags.DNS, "dns", "1.1.1.1:53", "default dns server")
	flag.StringVar(&flags.SOCKS, "socks_proxy", "", "socks proxy server")
	flag.StringVar(&flags.Mode, "mode", modeChecker, "[checker, monitor]")
	flag.StringVar(&flags.Log.Level, "log.level", "info", "log level")
	flag.Var(&flags.Regions, "region", "available regions: [all, Global, JP, TW, HK, NA]")
	flag.BoolVar(&flags.Version, "version", false, "display version and exit")
	flag.Parse()
}

func main() {
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{
		DisableColors:   true,
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339Nano,
	})
	log.SetOutput(os.Stdout)

	log.WithField("version", version).Infoln("Region Restriction Check")
	if flags.Version {
		return
	}

	switch flags.Mode {
	case modeChecker:
		runChecker()
	case modeMonitor:
		runMonitor()
	default:
		log.WithField("mode", flags.Mode).Fatalln("mode not found")
	}
}

func runChecker() {
	log.SetLevel(flags.Log.ParseLevel())
	flags.Regions.CheckAll()
	result := make(chan *model.CheckResult)

	cnt := 0
	checked := make(map[string]int)
	for _, region := range flags.Regions {
		if len(region) == 0 {
			continue
		}

		if mediaFuncsRegion, ok := medias.MediaFuncList[region]; ok {
			for mediaName, mediaFunc := range mediaFuncsRegion {
				if _, ok := checked[mediaName]; ok {
					continue
				}
				checked[mediaName] = 0
				cnt++
				go func(mn, rg string, f func(*model.Media) *model.CheckResult) {
					mc := model.NewMediaConf()
					mc.Name = mn
					mc.Region = rg
					mc.Logger = log.WithFields(log.Fields{
						"region": rg,
						"media":  mn,
					})
					mc.Timeout = 10
					mc.DNS = flags.DNS
					mc.SOCKS = flags.SOCKS
					r := f(mc)
					mc.Logger.WithFields(log.Fields{
						"result":  r.Result,
						"message": r.Message,
					}).Infoln("done")
					result <- r
				}(mediaName, region, mediaFunc)
			}
		} else {
			log.WithField("region", region).Errorln("region not found")
		}
	}

	var r model.CheckResultSlice
	for ; cnt > 0; cnt-- {
		r = append(r, <-result)
	}
	sort.Sort(&r)
	r.PrintTo(os.Stdout)
}

func runMonitor() {
	buf, err := ioutil.ReadFile(flags.ConfigFile)
	if err != nil {
		log.Fatalln(err)
	}

	err = json.Unmarshal(buf, &flags)
	if err != nil {
		log.Fatalln(err)
	}

	log.WithField("filename", flags.Log.Filename).Infoln("redirect log to file")
	log.SetLevel(flags.Log.ParseLevel())
	log.SetOutput(&lumberjack.Logger{
		Filename:   flags.Log.Filename,
		MaxSize:    flags.Log.MaxSize,
		MaxBackups: flags.Log.MaxBackups,
		MaxAge:     flags.Log.MaxAge,
		Compress:   false,
	})
	log.WithField("version", version).Infoln("Region Restriction Check")

	go func() {
		if err := exporter.ServeExporter(flags.ExporterListen); err != nil {
			log.Fatalln(err)
		}
		log.WithField("exporter_listen", flags.ExporterListen).Infoln("exporter listening")
	}()

	result := make(chan *model.CheckResult)

	for taskName, task := range flags.Tasks {
		if !task.Enabled {
			continue
		}

		for mediaName, mediaConf := range task.Medias {
			if !mediaConf.Enabled {
				continue
			}
			if mediaConf.Interval == 0 {
				mediaConf.Interval = task.Interval
			}
			mediaConf.Name = mediaName
			mlog := log.WithFields(log.Fields{
				"task":  taskName,
				"media": mediaName,
			})

			found := false
			for region, mediaFuncsRegion := range medias.MediaFuncList {
				if mediaFunc, ok := mediaFuncsRegion[mediaName]; ok {
					mlog = mlog.WithFields(log.Fields{
						"region": region,
					})
					mediaConf.Region = region

					go func(tn string, mc *model.Media, logger *log.Entry) {
						mc.Logger = logger
						for {
							res := mediaFunc(mc)
							res.Task = tn

							mc.Logger.WithFields(log.Fields{
								"result":  res.Result,
								"message": res.Message,
							}).Infoln("done")
							result <- res
							logger.WithField("interval", mc.Interval).Infoln("waiting")
							time.Sleep(time.Duration(mc.Interval) * time.Second)
						}
					}(taskName, mediaConf, mlog)
					found = true
					break
				}
			}
			if !found {
				mlog.Errorln("media not found")
			}
		}
	}

	for {
		res := <-result
		exporter.HandleStatusUpdate(res)
	}
}
