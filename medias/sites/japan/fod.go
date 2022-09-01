package japan

import (
	"bytes"
	"github.com/Hintay/region_restriction_check-go/medias/model"
	"github.com/valyala/fasthttp"
)

func CheckFOD(m *model.Media) (result *model.CheckResult) {
	m.URL = "https://geocontrol1.stream.ne.jp/fod-geo/check.xml?time=1624504256"
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

	if resp.StatusCode() != fasthttp.StatusOK {
		result.UnexpectedStatusCode(resp.StatusCode())
		return
	}

	if bytes.Contains(resp.Body(), []byte("true")) {
		result.Yes()
	} else {
		result.No()
	}
	return
}
