package exporter

import (
	"github.com/Hintay/region_restriction_check-go/medias/model"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	StatusMapping = map[string]int{
		model.CheckResultYes:           10,
		model.CheckResultInfo:          1,
		model.CheckResultNo:            0,
		model.CheckResultOriginalsOnly: -1,
		model.CheckResultOverseaOnly:   -2,
		model.CheckResultUnexpected:    -3,
		model.CheckResultFailed:        -4,
	}

	rrcStatus = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "rrc_unblock_status",
		Help: "Region Restriction Check Status",
	}, []string{"media", "media_readable", "task"})
)
