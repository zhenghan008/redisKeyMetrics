package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"redisKeyMetrics/parse"
	"redisKeyMetrics/statistics"
	"reflect"
	"strings"
)

var (
	redisKeyGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "redis_key_metrics",
			Help: "This is a redis Key metrics",
		},
		[]string{"key_type", "key_name", "key_unit"}, // 定义非固定的label
	)
)

func handleErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func updateGaugeHandler(redisAddr string, passwd string, sampleType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			bigKeyResult []*parse.BigKeyResult
			st           statistics.SampleType
			target       string
		)
		if redisAddr != "" {
			target = redisAddr
		} else {
			target = r.URL.Query().Get("target")
		}
		result := strings.Split(target, ":")
		if len(result) != 2 {
			return

		}
		switch sampleType {
		case "big":
			st = statistics.Bigkeys
		case "hot":
			st = statistics.Hotkeys
		case "mem":
			st = statistics.Memkeys
		default:
			st = statistics.Bigkeys
		}
		newSampleRes, err := statistics.GetSampleResult(result[0], result[1], passwd, st)
		handleErr(err)
		oldSampleRes, err := parse.ReadData()
		handleErr(err)
		b := &parse.BigKeyResult{}
		newBigKeyResult, _ := b.ParseBigKeyResult(1, newSampleRes)
		if oldSampleRes != "" {
			oldBigKeyResult, _ := b.ParseBigKeyResult(1, oldSampleRes)
			isDeepEqual := reflect.DeepEqual(newBigKeyResult, oldBigKeyResult)
			if isDeepEqual == false {
				err = statistics.ToResultFile(newSampleRes)
				handleErr(err)
				bigKeyResult = newBigKeyResult
			} else {
				bigKeyResult = oldBigKeyResult
			}

		} else {
			err = statistics.ToResultFile(newSampleRes)
			handleErr(err)
			bigKeyResult = newBigKeyResult
		}
		for _, each := range bigKeyResult {
			redisKeyGauge.With(prometheus.Labels{"key_type": each.StructureType, "key_name": strings.Replace(each.KeyName, "\"", "", -1), "key_unit": each.KeyUnit}).Set(each.KeySize)
			log.Println(each)
		}
		promhttp.Handler().ServeHTTP(w, r)
	}
}

func main() {
	redisPasswd := os.Getenv("REDIS_PASSWD")
	sampleType := os.Getenv("SAMPLE_TYPE")
	redisAddr := os.Getenv("REDIS_ADDR")
	// 创建HTTP处理程序来暴露指标
	http.HandleFunc("/metrics", updateGaugeHandler(redisAddr, redisPasswd, sampleType))
	// 启动HTTP服务器来暴露指标
	err := http.ListenAndServe(":9022", nil)
	handleErr(err)
}
