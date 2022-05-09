package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var onlineUsers = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "online_users",
	Help: "Online users",
	ConstLabels: map[string]string{
		"type": "test",
	},
})

var httpRequestsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "http_requests_total",
	Help: "Count of all HTTP requests for goapp",
}, []string{})

var httpDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Name: "http_request_duration",
	Help: "Duration in seconds of all HTTP requests",
}, []string{"handler"})

func main() {
	registry := prometheus.NewRegistry()
	registry.MustRegister(onlineUsers)
	registry.MustRegister(httpRequestsTotal)
	registry.MustRegister(httpDuration)

	go func() {
		for {
			onlineUsers.Set(float64(rand.Intn(2000)))
		}
	}()

	home := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(time.Duration(rand.Intn(8))*time.Second)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello World!"))
	})

	contact := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(time.Duration(rand.Intn(5))*time.Second)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Contact"))
	})

	homeHandler := promhttp.InstrumentHandlerDuration(
		httpDuration.MustCurryWith(prometheus.Labels{"handler": "home"}),
		promhttp.InstrumentHandlerCounter(httpRequestsTotal, home),
	)

	contactHandler := promhttp.InstrumentHandlerDuration(
		httpDuration.MustCurryWith(prometheus.Labels{"handler": "contact"}),
		promhttp.InstrumentHandlerCounter(httpRequestsTotal, contact),
	)

	http.Handle("/", homeHandler)
	http.Handle("/contact", contactHandler)
	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	
	log.Fatal(http.ListenAndServe(":8181", nil))
}
