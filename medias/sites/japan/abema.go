package japan

import (
	"encoding/json"
	"github.com/Hintay/region_restriction_check-go/medias/model"
	"github.com/valyala/fasthttp"
)

// {"isoCountryCode":"JP","timeZone":"Asia/Tokyo","utcOffset":"+09:00","cdnRegionUrl":"https://ds-linear-abematv.akamaized.net/region"}
func CheckAbemaTV(m *model.Media) (result *model.CheckResult) {
	m.URL = "https://api.abema.io/v1/ip/check?device=android"
	m.Logger.Infoln("running")

	if _, ok := m.Headers["User-Agent"]; !ok {
		m.Headers["User-Agent"] = model.UaDalvik
	}
	result = &model.CheckResult{Media: m.Name, Region: m.Region}

	resp, err := m.Do()
	if err != nil {
		m.Logger.Errorln(err)
		result.Failed(err)
		return
	}
	defer fasthttp.ReleaseResponse(resp)

	if resp.StatusCode() == fasthttp.StatusForbidden {
		result.No()
		return
	}

	var r struct {
		IsoCountryCode string `json:"isoCountryCode"`
	}
	err = json.Unmarshal(resp.Body(), &r)
	if err != nil {
		m.Logger.Errorln(err)
		result.Failed(err)
		return
	}

	switch r.IsoCountryCode {
	case "":
		result.No()
	case "Jp":
		result.Yes()
	default:
		result.Oversea()
	}
	return
}
