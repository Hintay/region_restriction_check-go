package model

import (
	"encoding/json"
	"errors"
	"github.com/Hintay/region_restriction_check-go/medias/dialer"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"time"
)

const (
	ContentTypeJSON      = "application/json"
	ContentTypeURLEncode = "application/x-www-form-urlencoded"

	UaDalvik  = "Dalvik/2.1.0 (Linux; U; Android 9; ALP-AL00 Build/HUAWEIALP-AL00)"
	UaBrowser = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.87 Safari/537.36"
)

type Media struct {
	Enabled  bool              `json:"enabled"`
	URL      string            `json:"-"`
	Method   string            `json:"method"`
	Addr     string            `json:"addr"`
	Headers  map[string]string `json:"headers"`
	Body     string            `json:"body"`
	DNS      string            `json:"dns"`
	SOCKS    string            `json:"socks_proxy"`
	Timeout  int               `json:"timeout"`
	Interval int               `json:"interval"`
	Name     string            `json:"-"`
	Region   string            `json:"-"`
	Logger   *log.Entry        `json:"-"`
}

func NewMediaConf() *Media {
	m := Media{}
	m.Headers = make(map[string]string)
	return &m
}

func (m *Media) UnmarshalJSON(data []byte) error {
	m.Enabled = true

	var result map[string]json.RawMessage
	err := json.Unmarshal(data, &result)
	if err != nil {
		return err
	}
	m.Headers = make(map[string]string)

	for k, v := range result {
		switch k {
		case "enabled":
			err = json.Unmarshal(v, &m.Enabled)
		case "method":
			err = json.Unmarshal(v, &m.Method)
		case "body":
			err = json.Unmarshal(v, &m.Body)
		case "timeout":
			err = json.Unmarshal(v, &m.Timeout)
		case "headers":
			err = json.Unmarshal(v, &m.Headers)
		case "interval":
			err = json.Unmarshal(v, &m.Interval)
		case "addr":
			err = json.Unmarshal(v, &m.Addr)
		}
		if err != nil {
			return err
		}
	}

	if m.Timeout == 0 {
		m.Timeout = 10
	}
	if m.Method == "" {
		m.Method = "GET"
	}
	return nil
}

func (m *Media) GetDial() fasthttp.DialFunc {
	switch m.SOCKS {
	case "":
		return dialer.TCPDialer(m.DNS, m.Addr, m.Logger)
	default:
		return dialer.SocksDialer(m.SOCKS, m.Addr)
	}
}

func (m *Media) Do() (*fasthttp.Response, error) {
	client := fasthttp.Client{Dial: m.GetDial()}
	client.NoDefaultUserAgentHeader = true

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(m.URL)
	req.Header.SetMethod(m.Method)
	for k, v := range m.Headers {
		req.Header.Set(k, v)
	}
	req.SetBodyString(m.Body)

	m.Logger = m.Logger.WithFields(log.Fields{
		"url":         string(req.URI().FullURI()),
		"method":      string(req.Header.Method()),
		"req_body":    string(req.Body()),
		"user_agent":  string(req.Header.UserAgent()),
		"timeout":     m.Timeout,
		"status_code": 0,
		"dns":         m.DNS,
		"socks":       m.SOCKS,
	})

	resp := fasthttp.AcquireResponse()
	if err := client.DoDeadline(req, resp, time.Now().Add(time.Duration(m.Timeout)*time.Second)); err != nil {
		fasthttp.ReleaseResponse(resp)
		return nil, err
	}
	m.Logger = m.Logger.WithField("status_code", resp.StatusCode())
	return resp, nil
}

func (m *Media) DoRedirects() (*fasthttp.Response, error) {
	cnt := 0
	for {
		resp, err := m.Do()
		if err != nil {
			return nil, err
		}
		status := resp.StatusCode()

		if status == fasthttp.StatusFound || status == fasthttp.StatusMovedPermanently {
			cnt += 1
			if cnt > 50 {
				return nil, errors.New("too many redirects")
			}
			m.URL = string(resp.Header.Peek("location"))
			fasthttp.ReleaseResponse(resp)
			continue
		} else {
			return resp, nil
		}
	}
}
