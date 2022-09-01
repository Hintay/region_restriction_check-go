package taiwan

import (
	"encoding/json"
	"fmt"
	"github.com/Hintay/region_restriction_check-go/medias/model"
	"github.com/valyala/fasthttp"
)

func CheckLiTV(m *model.Media) (result *model.CheckResult) {
	m.Method = "POST"
	m.URL = "https://www.litv.tv/vod/ajax/getUrl"
	m.Headers[fasthttp.HeaderContentType] = model.ContentTypeJSON
	m.Headers[fasthttp.HeaderContentType] = model.ContentTypeJSON
	m.Logger.Infoln("running")

	if _, ok := m.Headers["User-Agent"]; !ok {
		m.Headers["User-Agent"] = model.UaBrowser
	}
	if m.Body == "" {
		m.Body = `{"type":"noauth","assetId":"vod44868-010001M001_800K","puid":"6bc49a81-aad2-425c-8124-5b16e9e01337"}`
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
		return
	}

	r := make(map[string]interface{})
	err = json.Unmarshal(resp.Body(), &r)
	if err != nil {
		m.Logger.Errorln(err)
		result.Failed(err)
		return
	}

	if rr, ok := r["errorMessage"]; ok {
		if rr == nil {
			result.Yes()
		} else if rr.(string) == "vod.error.outsideregionerror" {
			result.No()
		} else {
			result.Unexpected(fmt.Sprintf("%+v", rr))
		}
	} else {
		result.Unexpected(`key "errorMessage" not found`)
	}

	return result
}
