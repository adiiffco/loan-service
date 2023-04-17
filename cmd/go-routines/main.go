package routines

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var grTicker *time.Ticker

func InitializeMetrics() {
	fmt.Println("Initializing routines")
	grCounter := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "goroutines_count",
		Help: "Number of go routines",
	})
	grTicker = time.NewTicker(time.Duration(1) * time.Second)
	go func() {
		for range grTicker.C {
			logrus.Debug("Current count of go routines: ", runtime.NumGoroutine())
			if viper.Get("log.level") == "DEBUG" {
				pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
			}
			grCounter.Set(float64(runtime.NumGoroutine()))
		}
	}()
}

func ShutDown() {
	grTicker.Stop()
}
