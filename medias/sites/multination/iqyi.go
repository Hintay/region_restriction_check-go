package multination

import (
	"github.com/Hintay/region_restriction_check-go/medias/model"
	"github.com/valyala/fasthttp"
	"strings"
)

func CheckiQYI(m *model.Media) (result *model.CheckResult) {
	if _, ok := m.Headers["User-Agent"]; !ok {
		m.Headers["User-Agent"] = model.UaBrowser
	}
	result = &model.CheckResult{Media: m.Name, Region: m.Region}

	m.URL = "https://www.iq.com/"
	resp, err := m.Do()
	if err != nil {
		m.Logger.Errorln(err)
		result.Failed(err)
		return
	}
	defer fasthttp.ReleaseResponse(resp)

	var cookie fasthttp.Cookie
	cookie.SetKey("mod")
	if resp.Header.Cookie(&cookie) {
		result.Yes("Region: ", strings.ToUpper(string(cookie.Value())))
	} else {
		result.Failed(`cookie "mod" not found`)
	}
	return
}
