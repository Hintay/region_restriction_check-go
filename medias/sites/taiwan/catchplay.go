package taiwan

import (
	"encoding/json"
	"fmt"
	"github.com/Hintay/region_restriction_check-go/medias/model"
	"github.com/valyala/fasthttp"
)

func CheckCatchplay(m *model.Media) (result *model.CheckResult) {
	m.URL = "https://sunapi.catchplay.com/geo"
	m.Logger.Infoln("running")

	if _, ok := m.Headers["User-Agent"]; !ok {
		m.Headers["User-Agent"] = model.UaBrowser
	}
	if _, ok := m.Headers["Authorization"]; !ok {
		m.Headers["Authorization"] = "Basic NTQ3MzM0NDgtYTU3Yi00MjU2LWE4MTEtMzdlYzNkNjJmM2E0Ok90QzR3elJRR2hLQ01sSDc2VEoy"
	}
	result = &model.CheckResult{Media: m.Name, Region: m.Region}

	resp, err := m.Do()
	if err != nil {
		m.Logger.Errorln(err)
		result.Failed(err)
		return
	}
	defer fasthttp.ReleaseResponse(resp)

	r := make(map[string]interface{})
	err = json.Unmarshal(resp.Body(), &r)
	if err != nil {
		m.Logger.Errorln(err)
		result.Failed(err)
		return
	}

	if rr, ok := r["code"]; ok {
		if rr.(string) == "0" {
			result.Yes()
		} else if rr.(string) == "100016" {
			result.No()
		} else {
			result.Unexpected(fmt.Sprintf("code: %+v", rr))
		}
	} else {
		result.Unexpected(fmt.Sprintf(`key "code" not found`))
	}

	return
}
