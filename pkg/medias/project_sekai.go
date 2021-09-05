package medias

import (
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

func CheckProjectSekai(m *Media) *CheckResult {
	m.Logger.Infoln("running")
	if m.URL == "" {
		m.URL = "https://game-version.sekai.colorfulpalette.org/1.8.1/3ed70b6a-8352-4532-b819-108837926ff5"
	}
	if _, ok := m.Headers["User-Agent"]; !ok {
		m.Headers["User-Agent"] = "pjsekai/48 CFNetwork/1240.0.4 Darwin/20.6.0"
	}
	result := &CheckResult{Media: m.Name, Region: m.Region, Type: "Game"}

	resp, err := m.Do()
	if err != nil {
		m.Logger.Errorln(err)
		result.Failed(err)
		return result
	}
	defer fasthttp.ReleaseResponse(resp)

	switch resp.StatusCode() {
	case fasthttp.StatusOK:
		result.Yes()
	case fasthttp.StatusForbidden:
		result.No()
	default:
		result.UnexpectedStatusCode(resp.StatusCode())
	}

	m.Logger.WithFields(log.Fields{
		"status_code": resp.StatusCode(),
		"result":      result.Result,
		"message":     result.Message,
	}).Infoln("done")
	return result
}
