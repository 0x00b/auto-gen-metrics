package auto_gen_metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"reflect"
	"strings"
)

func InitMetrics(prefix string, Metrics interface{}) {
	s := reflect.ValueOf(Metrics).Elem()
	typeOfAttr := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		f.SetInt(int64(i))
		name := typeOfAttr.Field(i).Name
		if ParseAttrName == nil {
			name =  DefaultParseAttrName(name)
		} else {
			name =  ParseAttrName(name)
		}
		if prefix !=""{
			name = prefix + "_" +name
		}
		ftype := f.Type().String()
		flen := len(ftype)
		switch true {
		case ftype[flen-len("Counter"):] == "Counter":
			counterMap[i] = prometheus.NewCounter(
				prometheus.CounterOpts{
					Name: name,
				})
			prometheus.MustRegister(counterMap[i])
		case ftype[flen-len("Gauge"):] == "Gauge":
			gaugeMap[i] = prometheus.NewGauge(
				prometheus.GaugeOpts{
					Name: name,
				})
			prometheus.MustRegister(gaugeMap[i])
		case ftype[flen-len("Summary"):] == "Summary":
			summaryMap[i] = prometheus.NewSummary(prometheus.SummaryOpts{
				Name: name,
				Help: "no help.",
				//Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
			})
			prometheus.MustRegister(summaryMap[i])
		case ftype[flen-len("Histogram"):] == "Histogram":
			histogramMap[i] = prometheus.NewHistogram(prometheus.HistogramOpts{
				Name: name,
				Help: "no help",
				//Buckets: prometheus.LinearBuckets(*normMean-5**normDomain, .5**normDomain, 20),
			})
			prometheus.MustRegister(histogramMap[i])
		default:
		}
	}
}

var ParseAttrName ParseAttrNameFunc = nil

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
	counterMap   = make(map[int]prometheus.Counter)
	summaryMap   = make(map[int]prometheus.Summary)
	gaugeMap     = make(map[int]prometheus.Gauge)
	histogramMap = make(map[int]prometheus.Histogram)
)

func AttrCounterInc(metric AttrType) {
	if v, ok := counterMap[metric.AttrValue()]; ok {
		v.Inc()
		return
	}
	panic("no this Counter")
}

func AttrCounterAdd(metric AttrType, value int) {
	if v, ok := counterMap[metric.AttrValue()]; ok {
		v.Add(float64(value))
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
