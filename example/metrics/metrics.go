package metrics

import (
	m "github.com/0x00b/auto-gen-metrics"
	"github.com/prometheus/client_golang/prometheus"
)

// 上报属性只需要两步操作如下：
// 1、在Metrics结构体中添加需要上报的属性名以及对应的类型, 使用InitMetrics注册 Metrics
// 2、使用相应的属相上报方法上报属性

//注：本例中的一切命名都可按自己的命名规范来

var M Metrics

// tag：pml
// format："name;label1,label2...;NameSpace...;..."

//默认会将这里的名字转换成小写的上报到prometheus，以符合prometheus和grafana的命名规范
//如果自己有自己的名字转换规则可以设置 attr.ParseAttrName 函数替换成自己的规则
//NOTE: 所有metric都要定义在同一个Metrics（struct）中
type Metrics struct {
	//上报的属性名/prometheus类型
	RECEIVE_REQUEST_RATE    m.Gauge   //目前只支持Counter ，Gauge功能还未实现
	RECEIVE_REQUEST_TOTAL   m.Counter `pml:"label1,label2;"` //有两个标签，label1,label2
	DEAL_REQUEST_SUCC_TOTAL m.Counter
	DEAL_REQUEST_FAIL_TOTAL m.Counter

	//TODO: add your metrics

}

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
