用来自动生成需要上报到prometheus的metrics，使用方法很简单，参考：
https://github.com/0x00b/auto-gen-metrics/blob/master/example/main.go

```

import (
	am "github.com/0x00b/auto-gen-metrics" 
	"github.com/prometheus/client_golang/prometheus"
)

// 上报属性只需要两步操作如下：
// 1、在PromMetrics结构体中添加需要上报的属性名以及对应的类型
// 2、使用相应的属相上报方法上报属性

func example() {

	am.InitMetrics("", &M, nil, nil)
	//am.InitMetrics(App+"_"+Module, &M, prometheus.Labels{"Public1": "test1", "Public2": "test2"},[]string{"AppId"})
	//am.InitMetrics(App+"_"+Module, &M, nil, []string{"AppId"})
	
	M.DEAL_REQUEST_SUCC_TOTAL.With(am.GetLabels(
			//prometheus.Labels{"AppId": "appid1"}),
		nil)).Inc()
	M.EXAMPLE_TOTAL.With(am.GetLabels(prometheus.Labels{
		//"AppId":  "appid1",
		"label1": "test1",
		"label2": "test2",
	})).Inc()
	M.RECEIVE_REQUEST_TOTAL.WithLabelValues(am.GetLvs(
		//"appid1",
		"xxx",
		"xxx")...).Inc()
	t, err := M.RECEIVE_REQUEST_TOTAL.GetMetricWith(am.GetLabels(prometheus.Labels{
		//"AppId": "appid1",
		"label1": "test1",
		"label2": "test2"}))
	if err == nil {
		t.Inc()
	}
	M.RECV_TEST_TOTAL.Inc()
}

var M PromMetrics

type PromMetrics struct {
	//反射说明：`pml:"name;label1,label2...;namespace;subsystem;example(help msg)"`
	// 1、name:
	// 		1) 反射中设置则使用反射设置，不设置则默认将变量名转换为小写。
	// 		2) InitMetrics第一个参数prefix是name前缀，若提供，则总会生效，
	// 		   可以使用prefix设置namespace和subsystem, 两者共存则最后metric为namespace_subsystem_prefix_name，
	// 		   只提供其中一个则为prefix_name,或者namespace_subsystem_name,
	// 		   可以用不提供namespace，subsystem，只提供prefix，让所有metric有相同的前缀
	// 		3) 如果要满足Go的命名规范，变量名用驼峰，不用下划线，则可以使用反射来设置name，变量名使用驼峰命名，
	// 		   但是这样带来一个不太便利的问题查找，比如我在prometheus视图中看到某个metric，
	// 		   要找到代码中对应的位置，需要先到这里找到对应变量名
	// 		4) name不能重复
	// 2、labels：
	//		1) labels可以设置多个，用","隔开，想要有label效果，则要使用TypeVec（eg:CounterVec）,如果只是Type（eg:Counter），则所有label不生效
	// 3、namespace、subsystem参加name
	// 4、help:
	// 		1)metric的help信息

	//上报的属性名/prometheus类型/反射
	EXAMPLE_TOTAL prometheus.CounterVec `pml:";label1,label2;namespace;subsystem;example(help msg)"`                   //name is example_total
	ExampleTotal  prometheus.CounterVec `pml:"example_test_total;label1,label2;namespace;subsystem;example(help msg)"` //name is example_test_total
	RECEIVE_REQUEST_RATE    prometheus.GaugeVec   //目前只支持Counter ，Gauge功能还未实现
	RECEIVE_REQUEST_TOTAL   prometheus.CounterVec `pml:";label1,label2;"` //有两个标签，label1,label2
	DEAL_REQUEST_SUCC_TOTAL prometheus.CounterVec
	DEAL_REQUEST_FAIL_TOTAL prometheus.CounterVec

	//TODO: add your metrics
	RECV_TEST_TOTAL prometheus.Counter `pml:";;;;just for test"`
}

```