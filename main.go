package main

import (
	"flag"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"redisKeyMetrics/parse"
	"redisKeyMetrics/statistics"
	"reflect"
	"strconv"
	"strings"
)

var (
	redisKeyGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "redis_key_metrics",
			Help: "This is a redis Key metrics",
		},
		[]string{"key_type", "key_name", "key_unit", "sample_type"}, // 定义非固定的label
	)
)

func handleErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

// Compare old and new content and generate metrics
func compareAndGenMetrics(newSampleRes string, stFlags string) error {
	var bigKeyResult []*parse.BigKeyResult
	b := &parse.BigKeyResult{}
	oldSampleRes, err := parse.ReadData(stFlags)
	if err != nil {
		return err
	}
	newBigKeyResult, _ := b.ParseBigKeyResult(newSampleRes, stFlags)
	if oldSampleRes != "" {
		oldBigKeyResult, err := b.ParseBigKeyResult(oldSampleRes, stFlags)
		if err != nil {
			return err
		}
		isDeepEqual := reflect.DeepEqual(newBigKeyResult, oldBigKeyResult)
		if isDeepEqual == false {
			err = statistics.ToResultFile(newSampleRes, stFlags)
			if err != nil {
				return err
			}
			bigKeyResult = newBigKeyResult
		} else {
			bigKeyResult = oldBigKeyResult
		}

	} else {
		err = statistics.ToResultFile(newSampleRes, stFlags)
		if err != nil {
			return err
		}
		bigKeyResult = newBigKeyResult
	}
	for _, each := range bigKeyResult {
		redisKeyGauge.With(prometheus.Labels{"key_type": each.StructureType, "key_name": strings.Replace(each.KeyName, "\"", "", -1), "key_unit": each.KeyUnit, "sample_type": each.SampleType}).Set(each.KeySize)
		log.Println(each)
	}
	return nil

}

func updateGaugeHandler(redisAddr string, passwd string, sampleType string, isConcurrent bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			target           string
			newSampleResDict map[string]string
		)
		if redisAddr != "" {
			target = redisAddr
		} else {
			target = r.URL.Query().Get("target")
		}
		result := strings.Split(target, ":")
		if len(result) != 2 {
			log.Fatalln("not find the redis!redis addr is wrong!")
			return
		}
		stList := strings.Split(sampleType, "|")
		if len(stList) == 0 || len(stList) > 3 {
			log.Fatalln("the sampleType num is wrong")
			return
		}
		for _, each := range stList {
			if each != "big" && each != "mem" && each != "hot" {
				log.Fatalln("the sampleType only support hot mem big")
				return
			}
		}
		if isConcurrent {
			newSampleResDict = statistics.GetSampleResultConcurrent(result[0], result[1], passwd, stList)
		} else {
			newSampleResDict = statistics.GetSampleResultSerial(result[0], result[1], passwd, stList)
		}
		log.Printf("newSampleResDict: %s", newSampleResDict)
		if newSampleResDict != nil {
			for _, each := range stList {
				newSampleRes, ok := newSampleResDict[each]
				if ok {
					err := compareAndGenMetrics(newSampleRes, each)
					if err != nil {
						handleErr(err)
						return
					}
				}

			}
		}
		promhttp.Handler().ServeHTTP(w, r)
	}
}

func main() {
	port := flag.Int("P", 9022, "ports exposed by the service")
	isConcurrent := flag.Bool("c", false, "whether to enable concurrent execute command")
	redisPasswd := flag.String("p", "", "Redis password")
	sampleType := flag.String("s", "big", "Sample type example: big or big|hot or big|hot|mem")
	redisAddr := flag.String("h", "", "Redis address example: 127.0.0.1:6379")
	flag.Parse()
	// 创建HTTP处理程序来暴露指标
	http.HandleFunc("/metrics", updateGaugeHandler(*redisAddr, *redisPasswd, *sampleType, *isConcurrent))
	// 启动HTTP服务器来暴露指标
	err := http.ListenAndServe(":"+strconv.Itoa(*port), nil)
	handleErr(err)
}
