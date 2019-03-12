package monitor

import (
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"time"
)

// 初始化 web_reqeust_total， counter类型指标， 表示接收http请求总次数
var WebRequestTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "web_reqeust_total",
		Help: "Number of hello requests in total",
	},
	// 设置两个标签 请求方法和 路径 对请求总次数在两个
	[]string{"method", "endpoint"},
)

// web_request_duration_seconds，Histogram类型指标，bucket代表duration的分布区间
var WebRequestDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "web_request_duration_seconds",
		Help:    "web request duration distribution",
		Buckets: []float64{0.1, 0.3, 0.5, 0.7, 0.9, 1},
	},
	[]string{"method", "endpoint"},
)

func init() {
	// 注册监控指标
	prometheus.MustRegister(WebRequestTotal)
	prometheus.MustRegister(WebRequestDuration)
}


// 包装 handler function,不侵入业务逻辑
func Monitor(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h(w, r)
		duration := time.Since(start)
		WebRequestTotal.With(prometheus.Labels{"method": r.Method, "endpoint": r.URL.Path}).Inc()
		WebRequestDuration.With(prometheus.Labels{"method": r.Method, "endpoint": r.URL.Path}).Observe(duration.Seconds())
	}
}