package taiwan

import (
	"encoding/json"
	"fmt"
	"github.com/Hintay/region_restriction_check-go/medias/model"
	"github.com/valyala/fasthttp"
)

func CheckLineTV(m *model.Media) (result *model.CheckResult) {
	m.URL = "https://www.linetv.tw/api/part/11829/eps/1/part?chocomemberId="
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

	if resp.StatusCode() == fasthttp.StatusForbidden {
		result.No()
		return
	} else if resp.StatusCode() != fasthttp.StatusOK {
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

	if c, ok := r["countryCode"]; ok {
		cc := fmt.Sprintf("%v", c)
		if cc == "228" {
			result.Yes()
		} else if cc == "114" {
			result.No()
		} else {
			result.Unexpected(fmt.Sprintf("country code: %s", cc))
		}
	} else {
		result.Unexpected(`key "countryCode" not found`)
	}

	return result
}
