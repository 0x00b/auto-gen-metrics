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

var M Metrics

//默认会将这里的名字转换成小写的上报到prometheus，以符合prometheus和grafana的命名规范
//如果自己有自己的名字转换规则可以设置 attr.ParseAttrName 函数替换成自己的规则
type Metrics struct {
	//上报的属性名/prometheus类型
	RECEIVE_REQUEST_RATE    m.Gauge //目前只支持Counter ，Gauge功能还未实现
	RECEIVE_REQUEST_TOTAL   m.Counter
	DEAL_REQUEST_SUCC_TOTAL m.Counter
	DEAL_REQUEST_FAIL_TOTAL m.Counter

	//TODO: add your metrics

}

func main() {
	//初始化，最终的metrics = App_Module_MetricName
	//eg:  app_mod_receive_request_total
	m.InitMetrics(App+"_"+Module, &M)

	e := http.NewServeMux()
	e.Handle("/metrics", promhttp.Handler())
	e.HandleFunc("/test", func(writer http.ResponseWriter, request *http.Request) {
		//直接使用counter的上报方法上报属性即可
		m.AttrCounterInc(M.RECEIVE_REQUEST_TOTAL)
	})
	http.ListenAndServe(":8888", e)
}
