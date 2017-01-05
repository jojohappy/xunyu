package compute

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/xunyu/utils"
)

type AccesslogKafka struct {
	RealIp       string  `json:"real_ip"`
	Remote       string  `json:"remote"`
	Port         string  `json:"port"`
	Host         string  `json:"host"`
	ServerAddr   string  `json:"server_addr"`
	UpstreamAddr string  `json:"upstream_addr"`
	User         string  `json:"user"`
	Time         string  `json:"time"`
	Method       string  `json:"method"`
	Code         int     `json:"code"`
	Size         int     `json:"size"`
	Referer      string  `json:"referer"`
	RequestTime  float64 `json:"request_time"`
	Agent        string  `json:"agent"`
	Scheme       string  `json:"scheme"`
	Protocol     string  `json:"protocol"`
	RequestId    string  `json:"request_id"`
	Path         string  `json:"path"`
}

type Accesslog struct {
	Uuid             string  `json:"uuid"`
	Host             string  `json:"host"`
	Method           string  `json:"method"`
	Requests         int     `json:"requests"`
	Time             int64   `json:"time"`
	RequestTimeMin   float64 `json:"request_time_min"`
	RequestTimeMax   float64 `json:"request_time_max"`
	RequestTimeAvg   float64 `json:"request_time_avg"`
	Code200          int     `json:"code_200"`
	Code500          int     `json:"code_500"`
	Code502          int     `json:"code_502"`
	Code503          int     `json:"code_503"`
	Code404          int     `json:"code_404"`
	Code2xx          int     `json:"code_2xx"`
	Code3xx          int     `json:"code_3xx"`
	Code4xx          int     `json:"code_4xx"`
	Code5xx          int     `json:"code_5xx"`
	ErrorNum         int     `json:"error_num"`
	SuccessRate      float64 `json:"success_rate"`
	Timestamp        string  `json:"@timestamp"`
	TrendError       float64 `json:"trend_error"`
	TrendRequestTime float64 `json:"trend_request_time"`
	TrendSuccessRate float64 `json:"trend_success_rate"`
	TrendRequests    float64 `json:"trend_requests"`
}

var (
	buffer utils.Buffer = utils.NewBufferList()
)

func parseDatetime(datetime string) string {
	t, _ := time.Parse(time.RFC3339, datetime)
	var second string
	if s := t.Second(); s <= 30 {
		second = "00"
	} else {
		second = "30"
	}
	return fmt.Sprintf("%d%02d%02d%02d%02d%s",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), second)
}

func datetimeUnix(datetime string) int64 {
	t, _ := time.Parse(time.RFC3339, datetime)
	if s := t.Second(); s <= 30 {
		return t.Unix() - int64(s)
	} else {
		return t.Unix() - int64(s) + 30
	}
}

func statusCode(accesslog *Accesslog, code int) {
	switch {
	case code == 200:
		accesslog.Code200 += 1
		accesslog.Code2xx += 1
	case code == 500:
		accesslog.Code500 += 1
		accesslog.Code5xx += 1
	case code == 502:
		accesslog.Code502 += 1
		accesslog.Code5xx += 1
	case code == 503:
		accesslog.Code503 += 1
		accesslog.Code5xx += 1
	case code == 404:
		accesslog.Code404 += 1
		accesslog.Code4xx += 1
	case code >= 200 && code < 300:
		accesslog.Code2xx += 1
	case code >= 300 && code < 400:
		accesslog.Code3xx += 1
	case code >= 400 && code < 500:
		accesslog.Code4xx += 1
	case code >= 500 && code < 600:
		accesslog.Code5xx += 1
	}
}

func convert2Accesslog(accesslogKafka AccesslogKafka, uuid string) Accesslog {
	accesslog := Accesslog{
		Uuid:           uuid,
		Host:           accesslogKafka.Host,
		Method:         accesslogKafka.Method,
		Requests:       1,
		Time:           datetimeUnix(accesslogKafka.Time),
		RequestTimeAvg: accesslogKafka.RequestTime,
		RequestTimeMin: accesslogKafka.RequestTime,
		RequestTimeMax: accesslogKafka.RequestTime,
		Code200:        0,
		Code500:        0,
		Code502:        0,
		Code503:        0,
		Code404:        0,
		Code2xx:        0,
		Code3xx:        0,
		Code4xx:        0,
		Code5xx:        0,
	}
	statusCode(&accesslog, accesslogKafka.Code)
	accesslog.ErrorNum = accesslog.Code5xx
	accesslog.SuccessRate = float64(accesslog.Code200 * 100)
	return accesslog
}

func Aggregate(in <-chan string) {
	for s := range in {
		var accesslogKafka AccesslogKafka
		json.NewDecoder(strings.NewReader(s)).Decode(&accesslogKafka)
		uuid := accesslogKafka.Host + "" + parseDatetime(accesslogKafka.Time)
		if el := buffer.Find(func(v interface{}) bool {
			var accesslog Accesslog = v.(Accesslog)
			return accesslog.Uuid == uuid
		}); nil != el {
			v := el.Value.(Accesslog)
			if v.RequestTimeMax < accesslogKafka.RequestTime {
				v.RequestTimeMax = accesslogKafka.RequestTime
			}
			if v.RequestTimeMin > accesslogKafka.RequestTime {
				v.RequestTimeMin = accesslogKafka.RequestTime
			}
			v.RequestTimeAvg = (v.RequestTimeAvg*float64(v.Requests) + accesslogKafka.RequestTime) / (float64(v.Requests) + 1)
			v.Requests += 1
			statusCode(&v, accesslogKafka.Code)
			v.SuccessRate = float64(v.Code200 * 100) / float64(v.Requests)
			el.Value = v
		} else {
			buffer.Push(convert2Accesslog(accesslogKafka, uuid))
		}
	}
}

func Resolve(data <-chan string) <-chan interface{} {
	out := make(chan interface{})
	go Aggregate(data)
	go func() {
		ticker := time.NewTicker(time.Second)
		for t := range ticker.C {
			if s := t.Second(); s == 5 || s == 35 {
				// filter 30 seconds ago
				items := buffer.Filter(func(v interface{}) bool {
					// var accesslog Accesslog = v.(Accesslog)
					return true
				})
				fmt.Println("length: ", len(items))
				for _, item := range items {
					out <- item
				}
			}
		}
	}()
	return out
}
