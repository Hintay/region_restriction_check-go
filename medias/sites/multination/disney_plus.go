package multination

import (
	"encoding/json"
	"errors"
	"github.com/Hintay/region_restriction_check-go/medias/model"
	"github.com/valyala/fasthttp"
	"net/url"
	"strings"
)

type disneyPlus struct {
	media  *model.Media
	result *model.CheckResult

	token string
}

var (
	returnError = errors.New("THIS ERROR IS FOR RETURN BREAK")
)

func (c *disneyPlus) Check(m *model.Media) *model.CheckResult {
	c.media = m
	c.setHeaders()
	c.result = &model.CheckResult{Media: m.Name, Region: m.Region}

	if err := c.process(); err != nil && err != returnError {
		c.media.Logger.Errorln(err)
		c.result.Failed(err)
	}
	return c.result
}

func (c *disneyPlus) process() error {
	if err := c.grantDevice(); err != nil {
		return err
	}
	if err := c.grantToken(); err != nil {
		return err
	}
	if err := c.fetchResult(); err != nil {
		return err
	}
	return nil
}

func (c *disneyPlus) setHeaders() {
	if _, ok := c.media.Headers["User-Agent"]; !ok {
		c.media.Headers["User-Agent"] = model.UaBrowser
	}
	c.media.Headers["authorization"] = "Bearer ZGlzbmV5JmJyb3dzZXImMS4wLjA.Cu56AgSfBTDag5NiRA81oLHkDZfu5L3CKadnefEAY84"
}

func (c *disneyPlus) grantDevice() error {
	c.media.URL = "https://disney.api.edge.bamgrid.com/devices"
	c.media.Method = "POST"
	c.media.Headers[fasthttp.HeaderContentType] = model.ContentTypeJSON
	c.media.Body = `{"deviceFamily":"browser","applicationRuntime":"chrome","deviceProfile":"windows","attributes":{}}`
	resp, err := c.media.Do()
	if err != nil {
		return err
	}
	defer fasthttp.ReleaseResponse(resp)

	if resp.StatusCode() == 403 {
		c.result.No()
		return returnError
	}

	var grant struct {
		Assertion string `json:"assertion"`
		Error     []struct {
			Code        string `json:"code"`
			Description string `json:"description"`
		} `json:"errors"`
	}
	err = json.Unmarshal(resp.Body(), &grant)
	if err != nil {
		return err
	}
	if grant.Assertion == "" {
		return errors.New(grant.Error[0].Description)
	}
	c.token = grant.Assertion
	return nil
}

func (c *disneyPlus) grantToken() error {
	c.media.URL = "https://disney.api.edge.bamgrid.com/token"
	c.media.Method = "POST"
	c.media.Headers[fasthttp.HeaderContentType] = model.ContentTypeURLEncode
	data := url.Values{
		"grant_type":         {"urn:ietf:params:oauth:grant-type:token-exchange"},
		"subject_token":      {c.token},
		"subject_token_type": {"urn:bamtech:params:oauth:token-type:device"},
	}
	c.media.Body = data.Encode()
	resp, err := c.media.Do()
	if err != nil {
		return err
	}
	defer fasthttp.ReleaseResponse(resp)

	var grant struct {
		RefreshToken     string `json:"refresh_token"`
		Error            string `json:"error"`
		ErrorDescription string `json:"error_description"`
	}
	err = json.Unmarshal(resp.Body(), &grant)
	if err != nil {
		return err
	}
	if grant.ErrorDescription == "forbidden-location" {
		c.result.No()
		return returnError
	}
	c.token = grant.RefreshToken
	return nil
}

func (c *disneyPlus) fetchResult() error {
	c.media.URL = "https://disney.api.edge.bamgrid.com/graph/v1/device/graphql"
	c.media.Method = "POST"
	c.media.Headers[fasthttp.HeaderContentType] = model.ContentTypeJSON
	c.media.Body = `{"query": "mutation refreshToken($input: RefreshTokenInput!) { refreshToken(refreshToken: $input){ activeSession { sessionId } } }", "variables": {"input": {"refreshToken": "` + c.token + `"} } }`
	resp, err := c.media.Do()
	if err != nil {
		return err
	}
	defer fasthttp.ReleaseResponse(resp)

	var query struct {
		Extensions struct {
			SDK struct {
				Session struct {
					InSupportedLocation bool `json:"inSupportedLocation"`
					Location            struct {
						CountryCode string `json:"countryCode"`
					} `json:"location"`
				} `json:"session"`
			} `json:"sdk"`
		} `json:"extensions"`
	}
	err = json.Unmarshal(resp.Body(), &query)
	if err != nil {
		return err
	}

	if query.Extensions.SDK.Session.InSupportedLocation {
		c.result.Yes("Region: ", query.Extensions.SDK.Session.Location.CountryCode)
		return nil
	}

	c.media.URL = "https://www.disneyplus.com"
	c.media.Method = "HEAD"
	c.media.Body = ""
	delete(c.media.Headers, fasthttp.HeaderContentType)
	resp, err = c.media.Do()
	if err != nil {
		return err
	}
	defer fasthttp.ReleaseResponse(resp)

	location := string(resp.Header.Peek("Location"))
	if strings.Contains(location, "unavailable") || strings.Contains(location, "hotstar") {
		c.result.No()
	} else {
		c.result.Info("Available For [Disney+ ",
			strings.ToUpper(query.Extensions.SDK.Session.Location.CountryCode), "] Soon")
	}

	return nil
}

func CheckDisneyPlus(m *model.Media) (result *model.CheckResult) {
	m.Logger.Infoln("running")
	var disneyPlus disneyPlus
	return disneyPlus.Check(m)
}
