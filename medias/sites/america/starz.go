package america

import (
	"encoding/json"
	"github.com/Hintay/region_restriction_check-go/medias/model"
	"github.com/valyala/fasthttp"
	"math/rand"
	"strings"
	"time"
)

const letterBytes = "0123456789abcdefghijklmnopqrstuvwxyz"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

// Code from https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func randString(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return sb.String()
}

func CheckStarz(m *model.Media) (result *model.CheckResult) {
	if _, ok := m.Headers["User-Agent"]; !ok {
		m.Headers["User-Agent"] = model.UaBrowser
	}
	result = &model.CheckResult{Media: m.Name, Region: m.Region}

	m.URL = "https://www.starz.com/sapi/header/v1/starz/us/" + randString(32)
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
