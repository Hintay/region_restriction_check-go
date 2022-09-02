package multination

import (
	"github.com/Hintay/region_restriction_check-go/medias/model"
	"github.com/valyala/fasthttp"
	"net/url"
	"strings"
)

func request(m *model.Media) (*fasthttp.Response, error) {
	m.Method = "POST"
	data := url.Values{"grant_type": {"refresh_token"}, "refresh_token": {"eyJ6aXAiOiJERUYiLCJraWQiOiJLcTYtNW1Ia3BxOXdzLUtsSUUyaGJHYkRIZFduRjU3UjZHY1h6aFlvZi04IiwiY3R5IjoiSldUIiwiZW5jIjoiQzIwUCIsImFsZyI6ImRpciJ9..OdwL8TEIFZouLDJe.wLz6zEC3PlPAGxx4X4qyP837lUbFrI_DQGnrJDMtEaQd5gsjHwaYshscoDXCYjMioU8JvsH_HKZga3fzSDEoWuMA5lgv4dyJpoB4Cqi91JjPSkqsRHKZ1I-nRoTmnSkcW3RHE-0coAqDWgK7IZ5cPiHQ-9KVRqqZkmTbEHynBdgH2y-FJP8zK0-dAynzR2krlUahhcykp7J7VqhZj_l5HVZZkPylZ6eKoK4J8fQvuGJoqMaRZTzrIH4Yk9J3GMbKnYqEG3SKRp5qAuWTtqLDOoGN0wWsUE5VRuCZxRKpxayJWABq2u4ABkAtIqUx8CPx77ZXxZVlcjRN1Xa8F2-e2mTxZq_1FgzmWECFg6onkDj_TpfBdeFoxDzhnRNceoQ-iyyNf3sgxJ_nz_bwztVZf0Vt3OR8yBnXfbkuEY7GQ4pvCuy-peW0mwJJCd2eJ9ADwDEGmoY4F47W-8rxdBhgna-0hu0FuLxt9MlmH_tGCmM_T-61xsxymLO9tlkwBnxNw4u6T9X2hcvC7-4uzr5cJiaJ3sGPMNo_ixTrP8SG9zCIse-X6_Lq0v3Uo-QOKhcD4N3gIfwZFYEvf-HVGWzFpU683q9CJfTTEXhsufj1URhSis7GdAa3nLZVt7CScsMPcYrMI317PmU-Brdvl_Ic4QeHTeF8-57kzD3mm5mrlQ7kQIXQzzQPqHYt70MzxL_scfT90cpYaSOBQnB1l--226h7X51XxSbrOcO-25zS7OSyedya8eMG6zAmgkk1zvZUzdCHZyzYD8-t0KYcfA5AwiLIFHxgqL4ni9fVy-SpYTKRwCmkp_pZOPaFwJh8zkhw8QaSLHq7ubko7H1kjJZxzsG1l4Bla1QRlj_-FVoY8GZ6okFk3Ts6A2qOK6v8UT7sL_w2zaHDQH1q2o05vsLwqIOxg3Xyey0tahzPbl-In_i1JGGvqGXOiPcKL5uOcTOo1luk32AbCS9i5mkopTS401YYYMH-Sx_krW_VJd2czpFefc0dlagtzBytqlcyscscFwq6IE6VHwG2Ij-WfO44G5hGDJFkZMZLeDUnTIyNrLe9hcfJp73koOSFnURsFWFjM2lgUIayiREAl02oh2alUyqnG09gdXufT_2W0DjA4i7qYuv6ol5NIVc389dF3x4a_7dPBvsMU3ppA1rlV04FlK6_fRv-Dk_jclXRZiQ5ul2ZO2CQ96LmrzmkdeNxFxcwaNXCJGBiRWXfMunoddIRg_LrVGuqWRgxj4DEnngZ2-qI_dliGiYraIehsHvtWeXIUWNF_FQSnQgZLg4WPekcluCecE4Iv7Sz36k9GUDqqs8hRWddirhufYem6RC84PyNqafCnwczrx5pOacmVzDl9Oi8OIhdDasdJa7gvsDoFzf6bv5st7EvbORkgPs6MK46mDMlwkL7TqjrJnSJzozCX4zLbYeyiWK6EXCehOpImMN262KLYQxnf5ugvk11gIA4NXpTbzyo4hp2LS7u8UMs5_w3t02vizxSQGokp-3qkEWmViy3pup1IXMPrcpS6KWHX0AYi1oRDZB5B8vM04pRHwYjsgMp2L-w4PMaDC4QDRU81IdvQ5VRkyLT7CL5hDlq5smXw_7wSFTWxs9vc5PmnrykSAkwFPocORC2j4T96uiu3z4gNoBu_dwKNcPi-dV7myC4iRRTmpm0V5A9IW510RGTyso_b-1hUeGvToYl9VwNgN7Impt3PjEQO2HXMU3p96tdulDEA_8bbyPdEGfxxVK3k2n_dxj_GzPKA8V4ESoNMRrV1vCuxPnrzfAOhqmNOEewTHqlxENSsZFGvfzVj1KemR7zLky14JMVslILnvxl6vuX7SbfIQ5JDktq9qKtTKo1mFrA-mBS3n00FacjPi364nnugiWQN7EwhNdEDH_KtWXGZVh-u2NM5cdoS1kAsOKSLxFTnTDG738LhoB3i_ZOjHFASKiZcsX6yD5csIP21jG5nFF9Qw2qsnqmxRuDLilIoGczEMt2Pfo180CG8Dyr7XtOYNeVU7__h9zBm9CvaAHDoQQhU4KlXM4LsljFeajw5f2wn08OmsdfkSYYl45O718QgzR_RRqwDpQH2pyKDJZ9yZt5OCyxcbnCgepjUyp6S-Pigfw73ASoCknhLLheb2mqkWIC-s3NmClpMoK-IyE57AiHHCatZfPGPnNofVioN5SbVR08mV7pdyQEhQGxGFM_LTAFFpwC48gOFTq-FWdV58muDULTqO3ImbGG6X3vV-PVbher1oJx0CFnelGGIx9lwM-yHbpVZGq9IXnKqoblCHiwuaJgbCKBnTjia2gYPNlN0Ql1ia3vQc7bybDVHyLePAVbOk10MdwHprwMGE__wsXqagElQCGJpU3ytPDktncRPCSQBQ3mw94CCIOQYEyhnA1Vik127AznwbR10Xm59diGBtix0Ao-VIrjKzQNw2hXqC_H-IgY46OT5ZndZ02SAe6AVyipq6kTui_ZyuQhy-zAOiat4t6qh-LyL1xImBuOZ7e79737LYiLHEIgHOIQ68DKcSmsIuA.gwrRhM5AiYUQ6iAbRZhxlw"}, "subject_token_type": {"urn:bamtech:params:oauth:token-type:device"}}
	m.Body = data.Encode()
	m.Headers["authorization"] = "Bearer ZGlzbmV5JmJyb3dzZXImMS4wLjA.Cu56AgSfBTDag5NiRA81oLHkDZfu5L3CKadnefEAY84"
	m.Headers["content-type"] = "application/x-www-form-urlencoded"
	return m.Do()
}

func CheckDisneyPlus(m *model.Media) (result *model.CheckResult) {
	if _, ok := m.Headers["User-Agent"]; !ok {
		m.Headers["User-Agent"] = model.UaBrowser
	}
	m.Logger.Infoln("running")
	result = &model.CheckResult{Media: m.Name, Region: m.Region}

	m.URL = "https://www.disneyplus.com"
	resp, err := request(m)
	if err != nil {
		m.Logger.Errorln(err)
		result.Failed(err)
		return
	}
	defer fasthttp.ReleaseResponse(resp)

	location := resp.Header.Peek("Location")
	if location == nil {
		setCookie := string(resp.Header.Peek("Set-Cookie"))
		startIndex := strings.Index(setCookie, "y=")
		endIndex := strings.Index(setCookie, ";")
		result.Yes("Region: ", setCookie[startIndex+2:endIndex])
	} else if string(location) == "https://disneyplus.disney.co.jp/" {
		result.Yes("Region: ", "JP")
	} else if string(location) == "https://preview.disneyplus.com/unavailable/" {
		result.No()
	} else {
		result.Unexpected("Unknown Location: ", string(location))
	}

	// Verify Authorized
	m.URL = "https://global.edge.bamgrid.com/token"
	resp, err = request(m)
	if err != nil {
		m.Logger.Errorln(err)
		result.Failed(err)
		return
	}
	defer fasthttp.ReleaseResponse(resp)
	if strings.Index(resp.String(), "unauthorized") != -1 {
		result.Message += " [Unauthorized]"
	}
	return
}
