package metrics

import (
	"net/http"
	"os"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	opsProcessedMutex  sync.Mutex
	infoProcessedMutex sync.Mutex
	wsMutex            sync.Mutex
	opsProcessed       = promauto.NewCounter(prometheus.CounterOpts{
		Name: "all_uses",
		Help: "The total number of active uses",
	})
	infoProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "all_info_uses",
		Help: "The total number of info uses",
	})
	conns = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "all",
		Subsystem: "mq",
		Name:      "conns",
		Help:      "All current websocket connections",
	})
)

func Incr() {
	opsProcessedMutex.Lock()
	opsProcessed.Inc()
	opsProcessedMutex.Unlock()
}

func IncrInf() {
	infoProcessedMutex.Lock()
	infoProcessed.Inc()
	infoProcessedMutex.Unlock()
}

func IncrC() {
	wsMutex.Lock()
	conns.Inc()
	wsMutex.Unlock()
}
func DecrC() {
	wsMutex.Lock()
	conns.Dec()
	wsMutex.Unlock()
}
func StartMonitor() {
	prometheus.MustRegister(conns)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":"+os.Getenv("PROM_PORT"), nil)
}
