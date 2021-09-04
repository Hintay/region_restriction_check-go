package medias

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

func CheckHamiVideo(m *Media) *CheckResult {
	m.Logger.Infoln("running")
	if m.URL == "" {
		m.URL = "https://hamivideo.hinet.net/api/play.do?id=OTT_VOD_0000249064&freeProduct=1"
	}
	if _, ok := m.Headers["User-Agent"]; !ok {
		m.Headers["User-Agent"] = UA_Dalvik
	}
	result := &CheckResult{Media: m.Name, Region: m.Region}

	resp, err := m.Do()
	if err != nil {
		m.Logger.Errorln(err)
		result.Message = err.Error()
		result.Result = CheckResultFailed
		return result
	}
	defer fasthttp.ReleaseResponse(resp)

	result.Result = CheckResultUnexpected
	if resp.StatusCode() == fasthttp.StatusOK {
		r := make(map[string]string)
		err = json.Unmarshal(resp.Body(), &r)
		if err != nil {
			result.Message = err.Error()
		} else {
			if c, ok := r["code"]; ok {
				if c == "06001-107" {
					result.Result = CheckResultYes
				} else if c == "06001-106" {
					result.Result = CheckResultNo
				} else {
					result.Message = fmt.Sprintf("code: %s", c)
				}
			} else {
				result.Message = "code not found"
			}
		}
	} else {
		result.Message = fmt.Sprintf("status code: %d", resp.StatusCode())
	}
	m.Logger.WithFields(log.Fields{
		"status_code": resp.StatusCode(),
		"result":      result.Result,
		"message":     result.Message,
	}).Infoln("done")
	return result
}
