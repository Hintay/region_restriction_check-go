package hongkong

import (
	"encoding/json"
	"fmt"
	"github.com/Hintay/region_restriction_check-go/medias/model"
	"github.com/valyala/fasthttp"
)

func CheckViuTV(m *model.Media) (result *model.CheckResult) {
	m.URL = "https://api.viu.now.com/p8/3/getLiveURL"
	m.Method = "POST"
	m.Headers[fasthttp.HeaderContentType] = model.ContentTypeJSON
	m.Logger.Infoln("running")

	if _, ok := m.Headers["User-Agent"]; !ok {
		m.Headers["User-Agent"] = model.UaBrowser
	}
	if m.Body == "" {
		m.Body = `{"callerReferenceNo":"20210726112323","contentId":"099","contentType":"Channel","channelno":"099","mode":"prod","deviceId":"29b3cb117a635d5b56","deviceType":"ANDROID_WEB"}`
	}
	result = &model.CheckResult{Media: m.Name, Region: m.Region}

	resp, err := m.Do()
	if err != nil {
		m.Logger.Errorln(err)
		result.Failed(err)
		return
	}
	defer fasthttp.ReleaseResponse(resp)

	if resp.StatusCode() != fasthttp.StatusOK {
		result.UnexpectedStatusCode(resp.StatusCode())
	}

	r := make(map[string]interface{})
	err = json.Unmarshal(resp.Body(), &r)
	if err != nil {
		m.Logger.Errorln(err)
		result.Failed(err)
		return
	}

	if rr, ok := r["responseCode"]; ok {
		switch rr {
		case "SUCCESS", "PRODUCT_INFORMATION_INCOMPLETE":
			result.Yes()
		case "GEO_CHECK_FAIL":
			result.No()
		default:
			result.Unexpected(fmt.Sprintf("result: %s", rr))
		}
	} else {
		result.Unexpected(`key "responseCode" not found`)
	}

	return
}
