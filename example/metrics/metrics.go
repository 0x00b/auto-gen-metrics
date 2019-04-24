package metrics

import (
	m "github.com/0x00b/auto-gen-metrics"
	"github.com/prometheus/client_golang/prometheus"
)

// 上报属性只需要两步操作如下：
// 1、在Metrics结构体中添加需要上报的属性名以及对应的类型, 使用InitMetrics注册 Metrics
// 2、使用相应的属相上报方法上报属性

//注：本例中的一切命名都可按自己的命名规范来


var PM PromMetrics

//另一种实现方式，不过不能进行功能扩展，例如公共的labels需要额外处理
//不过可以直接用prometheus的metric类型带来的全部能力
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


var M Metrics

type Metrics struct {
	//上报的属性名/prometheus类型
	RECEIVE_REQUEST_RATE    m.Gauge   //目前只支持Counter ，Gauge功能还未实现
	RECEIVE_REQUEST_TOTAL   m.Counter `pml:"label1,label2;"` //有两个标签，label1,label2
	DEAL_REQUEST_SUCC_TOTAL m.Counter
	DEAL_REQUEST_FAIL_TOTAL m.Counter

	//TODO: add your metrics

}