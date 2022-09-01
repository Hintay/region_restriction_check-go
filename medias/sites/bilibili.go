package sites

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/Hintay/region_restriction_check-go/medias/model"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
)

func CheckBilibiliTW(m *model.Media) *model.CheckResult {
	m.Logger.Infoln("running")
	result := &model.CheckResult{Media: m.Name, Region: m.Region}

	s := randomSession()
	m.URL = fmt.Sprintf("https://api.bilibili.com/pgc/player/web/playurl?avid=50762638&cid=100279344&qn=0&type=&otype=json&ep_id=268176&fourk=1&fnver=0&fnval=16&session=%s&module=bangumi", s)

	checkBilibili(m, result)
	return result
}

func CheckBilibiliHKMCTW(m *model.Media) *model.CheckResult {
	m.Logger.Infoln("running")
	result := &model.CheckResult{Media: m.Name, Region: "HK_MC_TW"}

	s := randomSession()
	m.URL = fmt.Sprintf("https://api.bilibili.com/pgc/player/web/playurl?avid=18281381&cid=29892777&qn=0&type=&otype=json&ep_id=183799&fourk=1&fnver=0&fnval=16&session=%s&module=bangumi", s)

	checkBilibili(m, result)
	return result
}

func randomSession() string {
	u := uuid.New().String()
	d := md5.New()
	d.Write([]byte(u))
	return fmt.Sprintf("%x", d.Sum(nil))
}

func checkBilibili(m *model.Media, result *model.CheckResult) {
	if _, ok := m.Headers["User-Agent"]; !ok {
		m.Headers["User-Agent"] = model.UaBrowser
	}
	result = &model.CheckResult{Media: m.Name, Region: m.Region}

	resp, err := m.Do()
	if err != nil {
		m.Logger.Errorln(err)
		result.Failed(err)
		return
	}
	defer fasthttp.ReleaseResponse(resp)

	switch resp.StatusCode() {
	case fasthttp.StatusOK:
		var r struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		}
		err = json.Unmarshal(resp.Body(), &r)
		if err != nil {
			m.Logger.Errorln(err)
			result.Failed(err)
			return
		}

		if r.Code == 0 {
			result.Yes()
		} else if r.Code == -10403 {
			result.No()
		} else {
			result.UnexpectedStatusCode(r.Code)
		}
	case fasthttp.StatusForbidden:
		result.No()
	default:
		result.UnexpectedStatusCode(resp.StatusCode())
	}
}
