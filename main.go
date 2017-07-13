package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tcnksm/go-httpstat"
)

var (
	addr          = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")
	target        = os.Getenv("TARGET")
	iterations, _ = strconv.Atoi(os.Getenv("NUM"))
	pause, _      = strconv.Atoi(os.Getenv("PAUSE"))
	agents, _     = strconv.Atoi(os.Getenv("AGENTS"))

	reqLatency = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "backpressure_req_latency",
		Help:    "The latency of the requests.",
		Buckets: prometheus.LinearBuckets(0, 100, 20),
	})

	reqOk = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "backpressure_req_ok_total",
		Help: "Number of successful requests.",
	})
	reqFail = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "backpressure_req_errors_total",
		Help: "Number of requests with errors.",
	})
	result httpstat.Result
)

func init() {
	prometheus.MustRegister(reqLatency)
	prometheus.MustRegister(reqOk)
	prometheus.MustRegister(reqFail)
}

func main() {

	for index := 0; index < agents; index++ {
		go func() {
			method := strings.ToLower(os.Getenv("METHOD"))
			switch method {
			case "get":
				doGet()
			case "post":
				doPost()
			}
		}()
	}

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))

}

func doGet() {

	for index := 0; index < iterations; index++ {
		begin := time.Now()

		req, err := http.NewRequest("GET", target, nil)

		if err != nil {
			panic(err)
		}
		client := http.DefaultClient
		res, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		if _, err := io.Copy(ioutil.Discard, res.Body); err != nil {
			log.Fatal(err)
		}
		if res.StatusCode == 200 {
			reqOk.Inc()
		} else {
			reqFail.Inc()
		}
		res.Body.Close()
		reqLatency.Observe(float64(time.Since(begin).Nanoseconds()))
		fmt.Printf("%d - %d\n", index, time.Since(begin).Nanoseconds())
		time.Sleep(time.Millisecond * time.Duration(pause))
	}
}

func doPost() {

}
