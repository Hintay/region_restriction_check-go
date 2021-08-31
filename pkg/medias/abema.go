package medias

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

// {"isoCountryCode":"JP","timeZone":"Asia/Tokyo","utcOffset":"+09:00","cdnRegionUrl":"https://ds-linear-abematv.akamaized.net/region"}
func CheckAbemaTV(m *Media) *CheckResult {
	m.Logger.Infoln("running")
	if m.URL == "" {
		m.URL = "https://api.abema.io/v1/ip/check?device=android"
	}
	if _, ok := m.Headers["User-Agent"]; !ok {
		m.Headers["User-Agent"] = "Dalvik/2.1.0 (Linux; U; Android 9; ALP-AL00 Build/HUAWEIALP-AL00)"
	}
	result := &CheckResult{Media: m.Name, Region: m.Region}

	resp, err := m.Do()
	if err != nil {
		m.Logger.Errorln(err)
		result.Error = err
		result.Result = CheckResultFailed
		return result
	}
	defer fasthttp.ReleaseResponse(resp)

	result.Result = CheckResultNo
	if resp.StatusCode() != fasthttp.StatusForbidden {
		r := make(map[string]string)
		err = json.Unmarshal(resp.Body(), &r)
		if err != nil {
			result.Error = err
			result.Result = CheckResultUnexpected
		} else {
			if reg, ok := r["isoCountryCode"]; ok {
				if reg == "JP" {
					result.Result = CheckResultYes
				} else {
					result.Result = CheckResultOverseaOnly
				}
			}
		}
	}

	m.Logger.WithFields(log.Fields{
		"status_code": resp.StatusCode(),
		"result":      result.Result,
		"error":       result.Error,
	}).Infoln("done")
	return result
}