package main

import (
	am "github.com/0x00b/auto-gen-metrics"
	m "github.com/0x00b/auto-gen-metrics/example/metrics"
	pm "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

const (
	//最好是小写
	Module = "mod"
	App    = "app"
)

func main() {
	//两种方式选一
	//SelfMetrics()
	PromMetrics()
}

func PromMetrics() {
	am.InitMetrics(App+"_"+Module, &m.PM, pm.Labels{"Public1": "test1", "Public2": "test2"},[]string{"AppId"})
	//am.InitMetrics2(App+"_"+Module, &m.PM, nil, []string{"AppId"})
	e := http.NewServeMux()
	e.Handle("/metrics", promhttp.Handler())
	e.HandleFunc("/test", func(writer http.ResponseWriter, request *http.Request) {
		m.PM.DEAL_REQUEST_SUCC_TOTAL.With(am.GetLabels(
			pm.Labels{"AppId": "appid1"})).Inc()

		m.PM.RECEIVE_REQUEST_TOTAL.With(am.GetLabels(pm.Labels{
			"AppId":  "appid1",
			"label1": "test1",
			"label2": "test2",
		})).Inc()

		m.PM.RECEIVE_REQUEST_TOTAL.WithLabelValues(am.GetLvs("appid1", "xxx", "xxx")...).Inc()

		t, err := m.PM.RECEIVE_REQUEST_TOTAL.GetMetricWith(am.GetLabels(pm.Labels{
			"AppId": "appid1", "label1": "test1", "label2": "test2"}))
		if err == nil {
			t.Inc()
		}
		m.PM.RecvTestTotal.Inc()
	})
	http.ListenAndServe(":8888", e)
}

func SelfMetrics() {
	am.InitMetrics2(App+"_"+Module, &m.M, pm.Labels{"Public1": "test1", "Public2": "test2"}, []string{"AppId"})
	e := http.NewServeMux()
	e.Handle("/metrics", promhttp.Handler())
	e.HandleFunc("/test", func(writer http.ResponseWriter, request *http.Request) {
		//直接使用counter的上报方法上报属性即可
		m.M.DEAL_REQUEST_SUCC_TOTAL.WithLabelValues("appid1").Inc()

		//有错panic
		m.M.RECEIVE_REQUEST_TOTAL.With(pm.Labels{
			"AppId":  "appid1",
			"label1": "test1",
			"label2": "test2",
		}).Inc()

		m.M.RECEIVE_REQUEST_TOTAL.WithLabelValues("appid1", "xxx", "xxx").Inc()

		//可以判断是否有错
		t, err := m.M.RECEIVE_REQUEST_TOTAL.GetMetricWith(pm.Labels{
			"AppId": "appid1", "label1": "test1", "label2": "test2"})
		if err == nil {
			t.Inc()
		}
	})
	http.ListenAndServe(":8888", e)
}
