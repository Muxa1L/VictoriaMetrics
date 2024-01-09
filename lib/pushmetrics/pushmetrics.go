package pushmetrics

import (
	"flag"
	"strings"
	"time"

	"github.com/VictoriaMetrics/VictoriaMetrics/lib/appmetrics"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/flagutil"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/logger"
	"github.com/VictoriaMetrics/metrics"
)

var (
	pushURL = flagutil.NewArrayString("pushmetrics.url", "Optional URL to push metrics exposed at /metrics page. See https://docs.victoriametrics.com/#push-metrics . "+
		"By default metrics exposed at /metrics page aren't pushed to any remote storage")
	pushInterval   = flag.Duration("pushmetrics.interval", 10*time.Second, "Interval for pushing metrics to -pushmetrics.url")
	pushExtraLabel = flagutil.NewArrayString("pushmetrics.extraLabel", "Optional labels to add to metrics pushed to -pushmetrics.url . "+
		`For example, -pushmetrics.extraLabel='instance="foo"' adds instance="foo" label to all the metrics pushed to -pushmetrics.url`)
)

func init() {
	// The -pushmetrics.url flag can contain basic auth creds, so it mustn't be visible when exposing the flags.
	flagutil.RegisterSecretFlag("pushmetrics.url")
}

var (
	// create a custom context for the pushmetrics module to close the metric reporting goroutine when the vmstorage process is shutdown.
	pushMetricsCtx, cancelPushMetric = context.WithCancel(context.Background())
)

// Init must be called after logger.Init
func Init() {
	extraLabels := strings.Join(*pushExtraLabel, ",")
	for _, pu := range *pushURL {
		opts := &metrics.PushOptions{
			ExtraLabels: extraLabels,
		}
		if err := metrics.InitPushExtWithOptions(pushMetricsCtx, pu, *pushInterval, appmetrics.WritePrometheusMetrics, opts); err != nil {
			logger.Fatalf("cannot initialize pushmetrics: %s", err)
		}
	}
}

func StopPushMetrics() {
	cancelPushMetric()
}
