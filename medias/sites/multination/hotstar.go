package multination

import (
	"github.com/Hintay/region_restriction_check-go/medias/model"
	"github.com/valyala/fasthttp"
	"strings"
)

func CheckHotStar(m *model.Media) (result *model.CheckResult) {
	if _, ok := m.Headers["User-Agent"]; !ok {
		m.Headers["User-Agent"] = model.UaBrowser
	}
	result = &model.CheckResult{Media: m.Name, Region: m.Region}

	m.URL = "https://www.hotstar.com/"
	m.Method = "HEAD"
	resp, err := m.Do()
	if err != nil {
		m.Logger.Errorln(err)
		result.Failed(err)
		return
	}
	defer fasthttp.ReleaseResponse(resp)

	var cookie fasthttp.Cookie
	cookie.SetKey("geo")
	if !resp.Header.Cookie(&cookie) {
		result.Unexpected(`cookie "geo" not found`)
		return
	}

	if resp.StatusCode() == 475 {
		result.No()
		return
	}

	region := string(cookie.Value()[:2])
	location := string(resp.Header.Peek("Location"))
	if strings.ToUpper(location[len(location)-2:]) == region {
		result.Yes("Region: ", region)
	} else {
		result.No()
	}

	return
}
