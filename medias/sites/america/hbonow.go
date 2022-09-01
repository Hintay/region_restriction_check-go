package america

import (
	"github.com/Hintay/region_restriction_check-go/medias/model"
	"github.com/valyala/fasthttp"
)

func CheckHBONow(m *model.Media) (result *model.CheckResult) {
	m.URL = "https://play.hbonow.com/"
	m.Logger.Infoln("running")

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
		result.Yes()
	case fasthttp.StatusFound:
		redirUrl := string(resp.Header.Peek("location"))
		if redirUrl == "http://hbogeo.cust.footprint.net/hbonow/geo.html" || redirUrl == "http://geocust.hbonow.com/hbonow/geo.html" {
			result.No()
		} else {
			result.Unexpected("URL:", redirUrl)
		}
	default:
		result.UnexpectedStatusCode(resp.StatusCode())
	}

	return
}
