package medias

import (
	"github.com/Hintay/region_restriction_check-go/medias/model"
	"github.com/Hintay/region_restriction_check-go/medias/sites"
	"github.com/Hintay/region_restriction_check-go/medias/sites/america"
	"github.com/Hintay/region_restriction_check-go/medias/sites/hongkong"
	"github.com/Hintay/region_restriction_check-go/medias/sites/japan"
	"github.com/Hintay/region_restriction_check-go/medias/sites/multination"
	"github.com/Hintay/region_restriction_check-go/medias/sites/taiwan"
)

var (
	mediasGlobal = map[string]func(*model.Media) *model.CheckResult{
		"Dazn":           multination.CheckDazn,
		"HotStar":        multination.CheckHotStar,
		"Netflix":        multination.CheckNetflix,
		"TVBAnywhere":    multination.CheckTVBAnywhere,
		"Viu.com":        multination.CheckViuCom,
		"DisneyPlus":     multination.CheckDisneyPlus,
		"YouTubePremium": multination.CheckYouTubePremium,
		"iQYI":           multination.CheckiQYI,
		// "NetflixCDN": CheckNetflixCDN,
		/*
			"YouTubeCDN": CheckYouTubeCDN,
			"PrimeVideo": PrimeVideo,
		*/
	}

	mediasJP = map[string]func(*model.Media) *model.CheckResult{
		"PCRJP":        japan.CheckPCRJP,
		"UMAJP":        japan.CheckUMAJP,
		"Kancolle":     japan.CheckKancolle,
		"KonosubaFD":   japan.CheckKonosubaFD,
		"ProjectSekai": japan.CheckProjectSekai,
		"AbemaTV":      japan.CheckAbemaTV,
		"HBOGoAsia":    sites.CheckHBOGoAsia,
		"DMM":          japan.CheckDMM,
		"Niconico":     japan.CheckNiconico,
		"Paravi":       japan.CheckParavi,
		"HuluJP":       japan.CheckHuluJP,
		"KaraokeDAM":   japan.CheckKaraokeDAM,
		"FOD":          japan.CheckFOD,
		"Radiko":       japan.CheckRadiko,
		/*"Unext":        CheckUnext,
		"TVer":   CheckTVer,
		"WOWOW":  CheckWOWOW,*/
	}

	mediasHK = map[string]func(*model.Media) *model.CheckResult{
		"BilibiliHKMCTW": sites.CheckBilibiliHKMCTW,
		"MyTVSuper":      hongkong.CheckMyTVSuper,
		"ViuTV":          hongkong.CheckViuTV,
		"NowE":           hongkong.CheckNowE,
		"HBOGoAsia":      sites.CheckHBOGoAsia,
	}

	mediasTW = map[string]func(*model.Media) *model.CheckResult{
		"BahamutAnime":   taiwan.CheckBahamutAnime,
		"BilibiliHKMCTW": sites.CheckBilibiliHKMCTW,
		"BilibiliTW":     sites.CheckBilibiliTW,
		"HBOGoAsia":      sites.CheckHBOGoAsia,
		"KKTV":           taiwan.CheckKKTV,
		"LiTV":           taiwan.CheckLiTV,
		"4GTV":           taiwan.Check4GTV,
		"LineTV":         taiwan.CheckLineTV,
		"HamiVideo":      taiwan.CheckHamiVideo,
		"Catchplay":      taiwan.CheckCatchplay,
		"ElevenSports":   taiwan.CheckElevenSports,
	}

	mediasNA = map[string]func(*model.Media) *model.CheckResult{
		"Fox": america.CheckFox,
		// "HuluUS": CheckHuluUS,
		"EPIX":   america.CheckEPIX,
		"Starz":  america.CheckStarz,
		"HBONow": america.CheckHBONow,
		"HBOMax": america.CheckHBOMax,
	}

	MediaFuncList = map[string]map[string]func(*model.Media) *model.CheckResult{
		"Global": mediasGlobal,
		"JP":     mediasJP,
		"TW":     mediasTW,
		"HK":     mediasHK,
		"NA":     mediasNA,
	}
)
