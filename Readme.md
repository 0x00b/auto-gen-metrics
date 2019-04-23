用来自动生成需要上报到prometheus的metrics，使用方法很简单，参考：

```
package main

import (
	m "auto-gen-metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

const (
	//最好是小写
	Module = "mod"
	App    = "app"
)

// 上报属性只需要两步操作如下：
// 1、在Metrics结构体中添加需要上报的属性名以及对应的类型
// 2、使用相应的属相上报方法上报属性

//注：本例中的一切命名都可按自己的命名规范来


var PM PromMetrics

//另一种实现方式，不过不能进行功能扩展，例如公共的labels需要额外处理
//不过可以直接用prometheus的metric类型带来的全部能力
type PromMetrics struct {
	//上报的属性名/prometheus类型
	RECEIVE_REQUEST_RATE    prometheus.GaugeVec   //目前只支持Counter ，Gauge功能还未实现
	RECEIVE_REQUEST_TOTAL   prometheus.CounterVec `pml:";label1,label2;"` //有两个标签，label1,label2
	DEAL_REQUEST_SUCC_TOTAL prometheus.CounterVec
	DEAL_REQUEST_FAIL_TOTAL prometheus.CounterVec

	//TODO: add your metrics
	RecvTestTotal prometheus.Counter `pml:"recv_test_total;;namespace;subsys;help msg"`
}


func main() {
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


```