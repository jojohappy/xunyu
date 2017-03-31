package apdex

import (
	"math/rand"
	"testing"
	"time"

	"github.com/xunyu/common"
	"github.com/xunyu/config"
)

type node map[string]interface{}

var (
	testRequest = []string{
		"{\"host\":\"example.a.com\",\"code\":200,\"request_time\":0.034",
		"{\"host\":\"example.b.com\",\"code\":200,\"request_time\":0.035",
		"{\"host\":\"example.c.com\",\"code\":200,\"request_time\":0.071",
		"{\"host\":\"example.a.com\",\"code\":200,\"request_time\":1.133",
		"{\"host\":\"example.b.com\",\"code\":200,\"request_time\":0.334",
		"{\"host\":\"example.a.com\",\"code\":200,\"request_time\":0.216",
		"{\"host\":\"example.b.com\",\"code\":200,\"request_time\":0.603",
		"{\"host\":\"example.a.com\",\"code\":200,\"request_time\":0.401",
		"{\"host\":\"example.b.com\",\"code\":200,\"request_time\":0.316",
		"{\"host\":\"example.a.com\",\"code\":200,\"request_time\":0.316",
		"{\"host\":\"example.b.com\",\"code\":200,\"request_time\":0.099",
		"{\"host\":\"example.a.com\",\"code\":200,\"request_time\":0.216",
		"{\"host\":\"example.b.com\",\"code\":200,\"request_time\":0.016",
		"{\"host\":\"example.b.com\",\"code\":200,\"request_time\":0.616",
		"{\"host\":\"example.a.com\",\"code\":200,\"request_time\":0.216",
		"{\"host\":\"example.b.com\",\"code\":200,\"request_time\":0.116",
		"{\"host\":\"example.a.com\",\"code\":200,\"request_time\":0.216",
		"{\"host\":\"example.b.com\",\"code\":200,\"request_time\":0.116",
		"{\"host\":\"example.a.com\",\"code\":200,\"request_time\":0.016",
	}
)

func genConfig() (*config.Config, error) {
	cfg, err := config.From(node{
		"rules": []node{
			node{
				"host":       "example.a.com",
				"satisfied":  200,
				"tolerating": 300,
			},
			node{
				"host":       "example.b.com",
				"satisfied":  100,
				"tolerating": 300,
			},
		},
	})

	return cfg, err
}

func genInput() <-chan common.DataInter {
	in := make(chan common.DataInter)
	go func() {
		t := time.Now()
		for j := 0; j < len(testRequest); j++ {
			rn := rand.Intn(45) * -1
			newT := t.Add(time.Second * time.Duration(rn))
			s := testRequest[j] + ",\"time\":\"" + newT.Format(time.RFC3339) + "\"}"
			in <- s
		}
		close(in)
	}()
	return in
}

func TestApdexConfig(t *testing.T) {
	cfg, err := genConfig()
	if nil != err {
		t.Error(err)
		t.FailNow()
	}

	ap, err := New(cfg)
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("instance apdex is %v", ap)
}

func TestApdexFilter(t *testing.T) {
	done := make(chan bool)
	in := genInput()
	out := filter(in)

	go func() {
		for data := range out {
			t.Logf("test apdex filter message: %v", data)
		}
		done <- true
	}()
	<-done
}

func filter(in <-chan common.DataInter) <-chan common.DataStr {
	cfg, _ := genConfig()
	ap, _ := New(cfg)

	out := make(chan common.DataStr, 1)

	go func(ap common.Pluginer) {
		ap.Filter(out)
		fc := ap.GetFilterChannel()
		for data := range in {
			fc <- data
		}
		close(fc)
	}(ap)

	return out
}
