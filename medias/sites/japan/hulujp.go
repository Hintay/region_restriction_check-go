package japan

import (
	"bytes"
	"github.com/Hintay/region_restriction_check-go/medias/model"
	"github.com/valyala/fasthttp"
)

func CheckHuluJP(m *model.Media) (result *model.CheckResult) {
	m.URL = "https://id.hulu.jp"
	m.Method = "HEAD"
	m.Logger.Infoln("running")
	if _, ok := m.Headers["User-Agent"]; !ok {
		m.Headers["User-Agent"] = model.UaDalvik
	}
	m.Headers["Accept"] = "*/*"
	m.Headers["Accept-Encoding"] = "gzip, deflate, br"
	result = &model.CheckResult{Media: m.Name, Region: m.Region}

	resp, err := m.Do()
	if err != nil {
		m.Logger.Errorln(err)
		result.Failed(err)
		return
	}
	defer fasthttp.ReleaseResponse(resp)

	if resp.StatusCode() != fasthttp.StatusFound {
		result.UnexpectedStatusCode(resp.StatusCode())
		return
	}

	redirUrl := resp.Header.Peek("location")
	if bytes.Contains(redirUrl, []byte("login")) {
		result.Yes()
	} else if bytes.Contains(redirUrl, []byte("restrict")) {
		result.No()
	} else {
		result.Unexpected("location: " + string(redirUrl))
	}

	return result
}
