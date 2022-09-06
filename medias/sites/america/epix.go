package america

import (
	"encoding/json"
	"github.com/Hintay/region_restriction_check-go/medias/model"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
)

func CheckEPIX(m *model.Media) (result *model.CheckResult) {
	if _, ok := m.Headers["User-Agent"]; !ok {
		m.Headers["User-Agent"] = model.UaBrowser
	}
	result = &model.CheckResult{Media: m.Name, Region: m.Region}

	m.URL = "https://api.epix.com/v2/sessions"
	m.Method = "POST"
	m.Body = `{"device":{"guid":"` + uuid.New().String() + `","format":"console","os":"web","app_version":"1.0.2","model":"browser","manufacturer":"google"},"apikey":"f07debfcdf0f442bab197b517a5126ec","oauth":{"token":null}}`
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

	var session struct {
		DeviceSession struct {
			SessionToken string `json:"session_token"`
		} `json:"device_session"`
	}
	err = json.Unmarshal(resp.Body(), &session)
	if err != nil {
		m.Logger.Errorln(err)
		result.Failed(err)
		return
	}

	m.URL = "https://api.epix.com/v2/movies/16921/play"
	m.Body = `{}`
	m.Headers["X-Session-Token"] = session.DeviceSession.SessionToken
	resp, err = m.Do()
	if err != nil {
		m.Logger.Errorln(err)
		result.Failed(err)
		return
	}
	defer fasthttp.ReleaseResponse(resp)

	var movie struct {
		Movie struct {
			Entitlements struct {
				Status string `json:"status"`
			} `json:"entitlements"`
		} `json:"movie"`
	}
	err = json.Unmarshal(resp.Body(), &movie)
	if err != nil {
		m.Logger.Errorln(err)
		result.Failed(err)
		return
	}

	switch movie.Movie.Entitlements.Status {
	case "PROXY_DETECTED":
		result.No()
	case "GEO_BLOCKED":
		result.No()
	default:
		result.Yes()
	}
	return
}
