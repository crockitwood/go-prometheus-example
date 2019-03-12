package main

import (
	"github.com/crockitwood/go-prometheus-example/monitor"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	// expose prometheus metrics接口
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/hello", Hello)
	http.HandleFunc("/query", monitor.Monitor(Query))
	log.Fatal(http.ListenAndServe(":8080", nil))
}


// hello
func Hello(w http.ResponseWriter, r *http.Request)  {
	// 请求计数
	monitor.WebRequestTotal.With(prometheus.Labels{"method": r.Method, "endpoint": r.URL.Path}).Inc()
	_,_ = io.WriteString(w, "hello world!")
}


// query
func Query(w http.ResponseWriter, r *http.Request)  {
	//模拟业务查询耗时0~1s
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	_,_ = io.WriteString(w, "some results")
}