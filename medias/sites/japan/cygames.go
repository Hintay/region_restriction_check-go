package japan

import (
	"github.com/Hintay/region_restriction_check-go/medias/model"
	"github.com/valyala/fasthttp"
)

func CheckPCRJP(m *model.Media) (result *model.CheckResult) {
	m.Logger.Infoln("running")
	if m.URL == "" {
		m.URL = "https://api-priconne-redive.cygames.jp/"
	}
	if _, ok := m.Headers["User-Agent"]; !ok {
		m.Headers["User-Agent"] = model.UaDalvik
	}
	m.Method = "HEAD"
	m.Http2 = true
	result = &model.CheckResult{Media: m.Name, Region: m.Region, Type: "Game"}

	resp, err := m.Do()
	if err != nil {
		m.Logger.Errorln(err)
		result.Failed(err)
		return
	}
	defer fasthttp.ReleaseResponse(resp)

	switch resp.StatusCode() {
	case fasthttp.StatusNotFound:
		result.Yes()
	case fasthttp.StatusForbidden:
		result.No()
	default:
		result.UnexpectedStatusCode(resp.StatusCode())
	}

	return
}

func CheckUMAJP(m *model.Media) *model.CheckResult {
	m.URL = "https://api-umamusume.cygames.jp/"
	return CheckPCRJP(m)
}
