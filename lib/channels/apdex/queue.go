package apdex

import (
    "fmt"
    "strconv"
	"time"

    "github.com/xunyu/common"
)

/*
access_log
{
  "real_ip": "-",
  "remote": "218.26.54.2",
  "port": "6232",
  "host": "example.a.com",
  "user": "-",
  "time": "2017-01-03T17:13:27+08:00",
  "method": "GET",
  "code": 200,
  "size": 4658,
  "referer": "-",
  "request_time": 0.034,
  "agent": "Mozilla/5.0 (iPhone; CPU iPhone OS 7_1_2 like Mac OS X) AppleWebKit/537.51.2 (KHTML, like Gecko) Mobile/11D257 isp/460.01 network/3G prokanqiu/7.0.17 Verizon iPhone 4",
  "scheme": "https",
  "protocol": "HTTP/1.1",
  "request_id": "-",
  "path": "/3/7.0.17/index?client=b5766cc0b78111c7481a071883320002c20ae6c2&night=0&crt=1483434807&time_zone=Asia/Shanghai&advId=5D8D230D-40F9-42DF-AAB5-804101B9426A&direc=next&preload=1"
}
*/

type Accesslog struct {
	Host         string    `json:"host"`
	Time         time.Time `json:"time"`
	Code         int       `json:"code"`
	RequestTime  float64   `json:"request_time"`
}

type ApdexResult struct {
	Host            string
	SatisfiedCount  int64
	ToleratingCount int64
	FrustratedCount int64
	Apdex           float64
	Time            time.Time
}

var queue []map[string]*ApdexResult

func NewApdexResult(
	host string,
	t time.Time,
) *ApdexResult {
	return &ApdexResult{
		Host:            host,
		SatisfiedCount:  0,
		ToleratingCount: 0,
		FrustratedCount: 0,
		Time:            t,
	}
}

func initQueue() {
	queue = make([]map[string]*ApdexResult, 2, 2)
	queue[0] = make(map[string]*ApdexResult)
	queue[1] = make(map[string]*ApdexResult)
}

func updateQueue(idx int, ar *ApdexResult) {
	queue[idx][ar.Host] = ar
}

func fetchResult(idx int, host string) (*ApdexResult, error) {
	val, ok := queue[idx][host]
	if ok {
		return val, nil
	} else {
		return nil, ErrHostNotExists
	}
}

func pop() map[string]*ApdexResult {
    q := queue[0]
    queue = append(queue[:0], queue[1:]...)
    queue = append(queue, make(map[string]*ApdexResult))
    return q
}

func (ar *ApdexResult) SetApdex() {
    total := float64(ar.SatisfiedCount + ar.ToleratingCount + ar.FrustratedCount)
    if total == 0 {
        ar.Apdex = 0
        return
    }
    apdex := float64(ar.SatisfiedCount + ar.ToleratingCount / 2) / total
    ar.Apdex, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", apdex), 64)
}

func (ar *ApdexResult) toDataStr() common.DataStr {
    out := common.DataStr{
        "Host": ar.Host,
        "Time": ar.Time,
        "Apdex": ar.Apdex,
    }
    return out
}
