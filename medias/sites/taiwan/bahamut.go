package taiwan

import (
	"encoding/json"
	"github.com/Hintay/region_restriction_check-go/medias/model"
	"github.com/valyala/fasthttp"
)

func CheckBahamutAnime(m *model.Media) (result *model.CheckResult) {
	m.URL = "https://ani.gamer.com.tw/ajax/token.php?adID=89422&sn=14667"
	m.Logger.Infoln("running")

	if _, ok := m.Headers["User-Agent"]; !ok {
		m.Headers["User-Agent"] = model.UaBrowser
	}
	result = &model.CheckResult{Media: m.Name, Region: m.Region}

	resp, err := m.Do()
	if err != nil {
		m.Logger.Errorln(err)
		result.Failed(err)
		return
	}
	defer fasthttp.ReleaseResponse(resp)

	switch resp.StatusCode() {
	case fasthttp.StatusOK:
		break
	case fasthttp.StatusForbidden:
		result.No()
		return
	default:
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
	if _, ok := r["animeSn"]; ok {
		result.Yes()
	} else {
		result.No()
	}

	return
}
