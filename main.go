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
	addr            = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")
	target          = os.Getenv("TARGET")
	iterations, err = strconv.Atoi(os.Getenv("NUM"))

	reqLatency = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "backpressure_req_latency",
		Help: "Latency of a request.",
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

	go func() {
		method := strings.ToLower(os.Getenv("METHOD"))
		switch method {
		case "get":
			doGet()
		case "post":
			doPost()
		}
	}()
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))

}

func doGet() {
	if err != nil {
		panic(err)
	}

	for index := 0; index < iterations; index++ {

		req, err := http.NewRequest("GET", target, nil)
		ctx := httpstat.WithHTTPStat(req.Context(), &result)
		req = req.WithContext(ctx)
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
		result.End(time.Now())

		duration := result.Total(time.Now())
		reqLatency.Set(float64(duration))
		fmt.Printf("%+v\n", result)
		fmt.Println("...")
	}
}

func doPost() {

}
