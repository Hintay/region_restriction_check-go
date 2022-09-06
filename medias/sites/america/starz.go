package america

import (
	"encoding/hex"
	"encoding/json"
	"github.com/Hintay/region_restriction_check-go/medias/model"
	"github.com/valyala/fasthttp"
	"math/rand"
)

func randString() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func CheckStarz(m *model.Media) (result *model.CheckResult) {
	if _, ok := m.Headers["User-Agent"]; !ok {
		m.Headers["User-Agent"] = model.UaBrowser
	}
	result = &model.CheckResult{Media: m.Name, Region: m.Region}

	m.URL = "https://www.starz.com/sapi/header/v1/starz/us/" + randString()
	m.Headers["Referer"] = "https://www.starz.com/us/en/"
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

	m.URL = "https://auth.starz.com/api/v4/User/geolocation"
	m.Headers["AuthTokenAuthorization"] = string(resp.Body())
	resp, err = m.Do()
	if err != nil {
		m.Logger.Errorln(err)
		result.Failed(err)
		return
	}
	defer fasthttp.ReleaseResponse(resp)

	var r struct {
		IsAllowedAccess   bool `json:"isAllowedAccess"`
		IsAllowedCountry  bool `json:"isAllowedCountry"`
		IsKnownProxy      bool `json:"isKnownProxy"`
		WhitelistOverride bool `json:"whitelistOverride"`
	}
	err = json.Unmarshal(resp.Body(), &r)
	if err != nil {
		m.Logger.Errorln(err)
		result.Failed(err)
		return
	}

	if r.WhitelistOverride || r.IsAllowedAccess && r.IsAllowedCountry && !r.IsKnownProxy {
		result.Yes()
	} else {
		result.No()
	}

	return
}
