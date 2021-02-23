package auto_gen_metrics

import (
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/prometheus/client_golang/prometheus"
)

func InitMetrics2(prefix string, Metrics interface{}, labels prometheus.Labels, publicTags []string) {
	s := reflect.ValueOf(Metrics).Elem()
	typeOfAttr := s.Type()
	publicLabels = labels
	var publicKeys []string
	publicLvs = publicLvs[:0]
	//for key, value := range publicLabels {
	//	publicKeys = append(publicKeys, key)
	//	publicLvs = append(publicLvs, value)
	//}
	for _, value := range publicTags {
		publicKeys = append(publicKeys, value)
	}
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		f.SetInt(int64(i))
		name := typeOfAttr.Field(i).Name
		rfTags := strings.Split(typeOfAttr.Field(i).Tag.Get("pml"), ";")
		fType := f.Type().String()
		fLen := len(fType)
		switch true {
		case fType[fLen-len("Counter"):] == "Counter":
			counterMap[i] = GetCounterVec(prefix, name, labels, publicKeys, rfTags)
			prometheus.MustRegister(counterMap[i])
		case fType[fLen-len("Gauge"):] == "Gauge":
			gaugeMap[i] = GetGaugeVec(prefix, name, labels, publicKeys, rfTags)
			prometheus.MustRegister(gaugeMap[i])
		//case fType[fLen-len("Summary"):] == "Summary":
		//	summaryMap[i] = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		//		Name: name,
		//	}, tags)
		//	prometheus.MustRegister(summaryMap[i])
		//case fType[fLen-len("Histogram"):] == "Histogram":
		//	histogramMap[i] = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		//		Name: name,
		//	}, tags)
		//	prometheus.MustRegister(histogramMap[i])
		default:
			panic("invalid metric type")
		}
	}
}

//另一种实现方式，不过不能进行功能扩展，例如公共的labels需要额外处理
//不过可以直接用prometheus的metric类型
func InitMetrics(prefix string, Metrics interface{}, labels prometheus.Labels, publicTags []string) {
	s := reflect.ValueOf(Metrics).Elem()
	typeOfAttr := s.Type()
	publicLabels = labels
	var publicKeys []string
	publicLvs = publicLvs[:0]
	//for key, value := range publicLabels {
	//	publicKeys = append(publicKeys, key)
	//	publicLvs = append(publicLvs, value)
	//}
	for _, value := range publicTags {
		publicKeys = append(publicKeys, value)
	}
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)

		fName := typeOfAttr.Field(i).Name

		rfTags := strings.Split(typeOfAttr.Field(i).Tag.Get("pml"), ";")

		var counterVec prometheus.CounterVec
		var gaugeVec prometheus.GaugeVec
		var sumVec prometheus.SummaryVec
		var histVec prometheus.HistogramVec

		fType := f.Type().String()
		switch true {
		//Counter
		case reflect.TypeOf(counterVec).String() == fType:
			metric := GetCounterVec(prefix, fName, publicLabels, publicKeys, rfTags)
			*(**prometheus.CounterVec)(unsafe.Pointer(s.FieldByName(fName).Addr().Pointer())) = metric
			prometheus.MustRegister(metric)
		case "prometheus.Counter" == fType:
			metric := GetCounter(prefix, fName, publicLabels, publicKeys, rfTags)
			*(*prometheus.Counter)(unsafe.Pointer(s.FieldByName(fName).Addr().Pointer())) = metric
			prometheus.MustRegister(metric)
			//Gauge
		case reflect.TypeOf(gaugeVec).String() == fType:
			metric := GetGaugeVec(prefix, fName, publicLabels, publicKeys, rfTags)
			*(**prometheus.GaugeVec)(unsafe.Pointer(s.FieldByName(fName).Addr().Pointer())) = metric
			prometheus.MustRegister(metric)
		case "prometheus.Gauge" == fType:
			metric := GetGauge(prefix, fName, publicLabels, publicKeys, rfTags)
			*(*prometheus.Gauge)(unsafe.Pointer(s.FieldByName(fName).Addr().Pointer())) = metric
			prometheus.MustRegister(metric)
			//Summary
		case reflect.TypeOf(sumVec).String() == fType:
			metric := GetSummaryVec(prefix, fName, publicLabels, publicKeys, rfTags)
			*(**prometheus.SummaryVec)(unsafe.Pointer(s.FieldByName(fName).Addr().Pointer())) = metric
			prometheus.MustRegister(metric)
		case "prometheus.Summary" == fType:
			metric := GetSummary(prefix, fName, publicLabels, publicKeys, rfTags)
			*(*prometheus.Summary)(unsafe.Pointer(s.FieldByName(fName).Addr().Pointer())) = metric
			prometheus.MustRegister(metric)
			//Histogram
		case reflect.TypeOf(histVec).String() == fType:
			metric := GetHistogramVec(prefix, fName, publicLabels, publicKeys, rfTags)
			*(**prometheus.HistogramVec)(unsafe.Pointer(s.FieldByName(fName).Addr().Pointer())) = metric
			prometheus.MustRegister(metric)
		case "prometheus.Histogram" == fType:
			metric := GetHistogram(prefix, fName, publicLabels, publicKeys, rfTags)
			*(*prometheus.Histogram)(unsafe.Pointer(s.FieldByName(fName).Addr().Pointer())) = metric
			prometheus.MustRegister(metric)
		default:
			panic("invalid metric type")
		}
	}
}

func GetCounter(prefix, name string, constLabel prometheus.Labels, publicKeys, rfTags []string) prometheus.Counter {
	name = getName(name, rfTags)
	namespace := ""
	if len(rfTags) > 2 {
		namespace = rfTags[2]
	}
	namespace = getNameSpace(prefix, namespace)
	subsystem := ""
	if len(rfTags) > 3 {
		subsystem = rfTags[3]
	}
	help := ""
	if len(rfTags) > 4 {
		help = rfTags[4]
	}
	return prometheus.NewCounter(
		prometheus.CounterOpts{
			Name:        name,
			Namespace:   namespace,
			Subsystem:   subsystem,
			Help:        help,
			ConstLabels: constLabel,
		})
}

func GetCounterVec(prefix, name string, constLabel prometheus.Labels, publicKeys, rfTags []string) *prometheus.CounterVec {
	tags := getTags(publicKeys, rfTags)
	name = getName(name, rfTags)
	namespace := ""
	if len(rfTags) > 2 {
		namespace = rfTags[2]
	}
	namespace = getNameSpace(prefix, namespace)
	subsystem := ""
	if len(rfTags) > 3 {
		subsystem = rfTags[3]
	}
	help := ""
	if len(rfTags) > 4 {
		help = rfTags[4]
	}
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:        name,
			Namespace:   namespace,
			Subsystem:   subsystem,
			Help:        help,
			ConstLabels: constLabel,
		}, tags)
}
func GetGaugeVec(prefix, name string, constLabel prometheus.Labels, publicKeys, rfTags []string) *prometheus.GaugeVec {
	tags := getTags(publicKeys, rfTags)
	name = getName(name, rfTags)
	namespace := ""
	if len(rfTags) > 2 {
		namespace = rfTags[2]
	}
	namespace = getNameSpace(prefix, namespace)
	subsystem := ""
	if len(rfTags) > 3 {
		subsystem = rfTags[3]
	}
	help := ""
	if len(rfTags) > 4 {
		help = rfTags[4]
	}
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        name,
			Namespace:   namespace,
			Subsystem:   subsystem,
			Help:        help,
			ConstLabels: constLabel,
		}, tags)
}

func GetGauge(prefix, name string, constLabel prometheus.Labels, publicKeys, rfTags []string) prometheus.Gauge {
	name = getName(name, rfTags)
	namespace := ""
	if len(rfTags) > 2 {
		namespace = rfTags[2]
	}
	namespace = getNameSpace(prefix, namespace)
	subsystem := ""
	if len(rfTags) > 3 {
		subsystem = rfTags[3]
	}
	help := ""
	if len(rfTags) > 4 {
		help = rfTags[4]
	}
	return prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name:        name,
			Namespace:   namespace,
			Subsystem:   subsystem,
			Help:        help,
			ConstLabels: constLabel,
		})
}

func GetSummaryVec(prefix, name string, constLabel prometheus.Labels, publicKeys, rfTags []string) *prometheus.SummaryVec {
	tags := getTags(publicKeys, rfTags)
	name = getName(name, rfTags)
	namespace := ""
	if len(rfTags) > 2 {
		namespace = rfTags[2]
	}
	namespace = getNameSpace(prefix, namespace)
	subsystem := ""
	if len(rfTags) > 3 {
		subsystem = rfTags[3]
	}
	objectives := make(map[float64]float64)
	if len(rfTags) > 4 && rfTags[4] != "" {
		values := strings.Split(rfTags[4], ",")
		for _, value := range values {
			keys := strings.Split(value, ":")
			if len(keys) == 2 {
				k, _ := strconv.ParseFloat(keys[0], 64)
				v, _ := strconv.ParseFloat(keys[1], 64)
				objectives[k] = v
			}
		}
	}
	maxAge := prometheus.DefMaxAge
	if len(rfTags) > 5 && rfTags[5] != "" {
		if v, err := strconv.ParseUint(rfTags[5], 10, 64); err == nil {
			maxAge = time.Second * time.Duration(v)
		}
	}
	var ageBuckets uint32 = prometheus.DefAgeBuckets
	if len(rfTags) > 6 && rfTags[6] != "" {
		if v, err := strconv.ParseUint(rfTags[6], 10, 32); err == nil {
			ageBuckets = uint32(v)
		}
	}
	var bufCap uint32 = prometheus.DefBufCap
	if len(rfTags) > 7 && rfTags[7] != "" {
		if v, err := strconv.ParseUint(rfTags[7], 10, 32); err == nil {
			bufCap = uint32(v)
		}
	}
	help := ""
	if len(rfTags) > 8 {
		help = rfTags[8]
	}
	return prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:        name,
			Namespace:   namespace,
			Subsystem:   subsystem,
			Help:        help,
			ConstLabels: constLabel,
			Objectives:  objectives,
			MaxAge:      maxAge,
			AgeBuckets:  ageBuckets,
			BufCap:      bufCap,
		}, tags)
}

func GetSummary(prefix, name string, constLabel prometheus.Labels, publicKeys, rfTags []string) prometheus.Summary {
	name = getName(name, rfTags)
	namespace := ""
	if len(rfTags) > 2 {
		namespace = rfTags[2]
	}
	namespace = getNameSpace(prefix, namespace)
	subsystem := ""
	if len(rfTags) > 3 {
		subsystem = rfTags[3]
	}
	objectives := make(map[float64]float64)
	if len(rfTags) > 4 && rfTags[4] != "" {
		values := strings.Split(rfTags[4], ",")
		for _, value := range values {
			keys := strings.Split(value, ":")
			if len(keys) == 2 {
				k, _ := strconv.ParseFloat(keys[0], 64)
				v, _ := strconv.ParseFloat(keys[1], 64)
				objectives[k] = v
			}
		}
	}
	maxAge := prometheus.DefMaxAge
	if len(rfTags) > 5 && rfTags[5] != "" {
		if v, err := strconv.ParseUint(rfTags[5], 10, 64); err == nil {
			maxAge = time.Second * time.Duration(v)
		}
	}
	var ageBuckets uint32 = prometheus.DefAgeBuckets
	if len(rfTags) > 6 && rfTags[6] != "" {
		if v, err := strconv.ParseUint(rfTags[6], 10, 32); err == nil {
			ageBuckets = uint32(v)
		}
	}
	var bufCap uint32 = prometheus.DefBufCap
	if len(rfTags) > 7 && rfTags[7] != "" {
		if v, err := strconv.ParseUint(rfTags[7], 10, 32); err == nil {
			bufCap = uint32(v)
		}
	}
	help := ""
	if len(rfTags) > 8 {
		help = rfTags[8]
	}
	return prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name:        name,
			Namespace:   namespace,
			Subsystem:   subsystem,
			Help:        help,
			ConstLabels: constLabel,
			Objectives:  objectives,
			MaxAge:      maxAge,
			AgeBuckets:  ageBuckets,
			BufCap:      bufCap,
		})
}

func getBuckets(rfTags []string) (buckets []float64) {
	if len(rfTags) > 4 {
		values := strings.Split(rfTags[4], ",")
		if len(values) > 0 && values[0] != "" {
			for _, v := range values {
				if value, err := strconv.ParseFloat(v, 64); err == nil {
					buckets = append(buckets, value)
				}
			}
		}
	}
	return
}
func GetHistogramVec(prefix, name string, constLabel prometheus.Labels, publicKeys, rfTags []string) *prometheus.HistogramVec {
	tags := getTags(publicKeys, rfTags)
	name = getName(name, rfTags)
	namespace := ""
	if len(rfTags) > 2 {
		namespace = rfTags[2]
	}
	namespace = getNameSpace(prefix, namespace)
	subsystem := ""
	if len(rfTags) > 3 {
		subsystem = rfTags[3]
	}

	help := ""
	if len(rfTags) > 5 {
		help = rfTags[5]
	}
	return prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:        name,
			Namespace:   namespace,
			Subsystem:   subsystem,
			Help:        help,
			ConstLabels: constLabel,
			Buckets:     getBuckets(rfTags),
		}, tags)
}

func GetHistogram(prefix, name string, constLabel prometheus.Labels, publicKeys, rfTags []string) prometheus.Histogram {
	name = getName(name, rfTags)
	namespace := ""
	if len(rfTags) > 2 {
		namespace = rfTags[2]
	}
	namespace = getNameSpace(prefix, namespace)
	subsystem := ""
	if len(rfTags) > 3 {
		subsystem = rfTags[3]
	}
	help := ""
	if len(rfTags) > 4 {
		help = rfTags[4]
	}
	return prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:        name,
			Namespace:   namespace,
			Subsystem:   subsystem,
			Help:        help,
			ConstLabels: constLabel,
			Buckets:     getBuckets(rfTags),
		})
}

func getTags(publicKeys, rfTags []string) []string {
	var tags []string
	for _, v := range publicKeys {
		tags = append(tags, v)
	}

	if len(rfTags) > 1 {
		rfLabels := strings.Split(rfTags[1], ",")
		if len(rfLabels) == 1 && rfLabels[0] == "" {
			rfLabels = nil
		}
		for _, v := range rfLabels {
			tags = append(tags, v)
		}
	}
	return tags
}
func getName(name string, rfTags []string) string {
	if len(rfTags) > 0 && len(rfTags[0]) > 0 {
		name = rfTags[0]
	} else {
		if ParseAttrName == nil {
			name = DefaultParseAttrName(name)
		} else {
			name = ParseAttrName(name)
		}
	}
	return name
}
func getNameSpace(prefix, name string) (namespace string) {
	if prefix != "" {
		namespace = prefix
	}
	if name != "" {
		if namespace != "" {
			namespace += "_" + name
		} else {
			namespace = name
		}
	}
	return
}

var (
	ParseAttrName ParseAttrNameFunc = nil
	publicLabels  prometheus.Labels //公共的label，value
	publicLvs     []string          //公共的value,有顺序的
)

//define Attr Types
type Counter int
type Gauge int
type Histogram int
type Summary int

type ParseAttrNameFunc func(string) string

func DefaultParseAttrName(name string) string {
	return strings.ToLower(name)
}

type AttrType interface {
	AttrValue() int
}

func (c Counter) AttrValue() int {
	return int(c)
}
func (c Gauge) AttrValue() int {
	return int(c)
}

func (c Histogram) AttrValue() int {
	return int(c)
}

func (c Summary) AttrValue() int {
	return int(c)
}

var (
	counterMap   = make(map[int]*prometheus.CounterVec)
	summaryMap   = make(map[int]*prometheus.SummaryVec)
	gaugeMap     = make(map[int]*prometheus.GaugeVec)
	histogramMap = make(map[int]*prometheus.HistogramVec)
)

//func GetLabels(labels prometheus.Labels) prometheus.Labels {
//	ls := make(prometheus.Labels, len(labels)+len(publicLabels))
//	for key, value := range publicLabels {
//		ls[key] = value
//	}
//	for key, value := range labels {
//		ls[key] = value
//	}
//	return ls
//}
//
//func GetLvs(lvs ...string) []string {
//	var vs []string
//	vs = append(vs, publicLvs...)
//	vs = append(vs, lvs...)
//	return vs
//}

func (c Counter) With(labels prometheus.Labels) prometheus.Counter {
	if v, ok := counterMap[c.AttrValue()]; ok {
		return v.With(labels)
	}
	panic("no this Counter")
}
func (c Counter) WithLabelValues(lvs ...string) prometheus.Counter {
	if v, ok := counterMap[c.AttrValue()]; ok {
		return v.WithLabelValues(lvs...)
	}
	panic("no this Counter")
}

func (c Counter) GetMetricWith(labels prometheus.Labels) (prometheus.Counter, error) {
	if v, ok := counterMap[c.AttrValue()]; ok {
		return v.GetMetricWith(labels)
	}
	panic("no this Counter")
}

func (c Counter) GetMetricWithLabelValues(lvs ...string) (prometheus.Counter, error) {
	if v, ok := counterMap[c.AttrValue()]; ok {
		return v.GetMetricWithLabelValues(lvs...)
	}
	panic("no this Counter")
}

//只给没有label的metric调用
func (c Counter) Inc() {
	if v, ok := counterMap[c.AttrValue()]; ok {
		v.With(publicLabels).Inc()
		return
	}
	panic("no this Counter")
}

//只给没有label的metric调用
func (c Counter) Add(metric AttrType, value int) {
	if v, ok := counterMap[metric.AttrValue()]; ok {
		v.With(publicLabels).Add(float64(value))
		return
	}
	panic("no this Counter")
}

func AttrGaugeSet(metric AttrType) {
	panic("no this Gauge")
}
func AttrGaugeInc(metric AttrType) {
	panic("no this Gauge")
}
func AttrGaugeAdd(metric AttrType) {
	panic("no this Gauge")
}
func AttrGaugeSub(metric AttrType) {
	panic("no this Gauge")
}
