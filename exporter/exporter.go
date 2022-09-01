package exporter

import (
	"github.com/Hintay/region_restriction_check-go/medias/model"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func HandleStatusUpdate(result *model.CheckResult) {
	rrcStatus.WithLabelValues(
		result.Media,
		model.HumanReadableNames[result.Media],
		result.Task,
	).Set(float64(StatusMapping[result.Result]))
}

func ServeExporter(listen string) error {
	http.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(listen, nil)
}
