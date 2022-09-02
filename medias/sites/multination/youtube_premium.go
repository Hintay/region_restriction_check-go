package multination

import (
	"github.com/Hintay/region_restriction_check-go/medias/model"
	"github.com/valyala/fasthttp"
	"strings"
)

func CheckYouTubePremium(m *model.Media) (result *model.CheckResult) {
	if _, ok := m.Headers["User-Agent"]; !ok {
		m.Headers["User-Agent"] = model.UaBrowser
	}
	result = &model.CheckResult{Media: m.Name, Region: m.Region}

	m.URL = "https://www.youtube.com/premium"
	m.Headers["Accept-Language"] = "en"
	resp, err := m.Do()
	if err != nil {
		m.Logger.Errorln(err)
		result.Failed(err)
		return
	}
	defer fasthttp.ReleaseResponse(resp)

	body := string(resp.Body())
	if strings.Contains(body, "www.google.cn") {
		result.No("Region: CN")
		return
	}

	if strings.Contains(body, "manageSubscriptionButton") {
		regionStart := strings.Index(body, "\"countryCode\"")
		if regionStart >= 0 {
			result.Yes("Region: ", body[regionStart+15:regionStart+17])
		} else {
			result.Yes()
		}
	} else {
		result.No()
	}

	return
}
