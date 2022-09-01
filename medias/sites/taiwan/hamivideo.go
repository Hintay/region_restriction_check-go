package taiwan

import (
	"encoding/json"
	"fmt"
	"github.com/Hintay/region_restriction_check-go/medias/model"
	"github.com/valyala/fasthttp"
)

func CheckHamiVideo(m *model.Media) (result *model.CheckResult) {
	m.URL = "https://hamivideo.hinet.net/api/play.do?id=OTT_VOD_0000249064&freeProduct=1"
	m.Logger.Infoln("running")

	if _, ok := m.Headers["User-Agent"]; !ok {
		m.Headers["User-Agent"] = model.UaBrowser
	}
	result = &model.CheckResult{Media: m.Name, Region: m.Region}

	resp, err := m.Do()
	if err != nil {
		m.Logger.Errorln(err)
		result.Failed(err)
		return result
	}
	defer fasthttp.ReleaseResponse(resp)

	if resp.StatusCode() != fasthttp.StatusOK {
		result.UnexpectedStatusCode(resp.StatusCode())
	}
	r := make(map[string]string)
	err = json.Unmarshal(resp.Body(), &r)
	if err != nil {
		m.Logger.Errorln(err)
		result.Failed(err)
		return
	}

	if c, ok := r["code"]; ok {
		if c == "06001-107" {
			result.Yes()
		} else if c == "06001-106" {
			result.No()
		} else {
			result.Unexpected(fmt.Sprintf("code: %s", c))
		}
	} else {
		result.Unexpected(`key "code" not found`)
	}

	return
}
