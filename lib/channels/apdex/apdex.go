package apdex

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/xunyu/common"
	"github.com/xunyu/config"
	"github.com/xunyu/lib/log"
)

type ApdexRule struct {
	Host       string  `config:"host"`
	Satisfied  float64 `config:"satisfied"`
	Tolerating float64 `config:"tolerating"`
}

type ApdexConfig struct {
	Rules []ApdexRule `config:"rules"`
}

type apdex struct {
	common.PluginPrototype
	c      chan common.DataInter
	config ApdexConfig
}

var (
	defaultConfig = ApdexConfig{}
)

func init() {
	common.RegisterChannelPlugin("apdex", New)
}

func New(cfg *config.Config) (common.Pluginer, error) {
	ap := &apdex{config: defaultConfig}
	if err := ap.init(cfg); nil != err {
		return nil, err
	}

	return ap, nil
}

func (ap *apdex) init(config *config.Config) error {
	if err := config.Assemble(&ap.config); nil != err {
		return err
	}
	initQueue()

	return nil
}

func (ap *apdex) GetFilterChannel() chan<- common.DataInter {
	return ap.c
}

func (ap *apdex) Filter(out chan<- common.DataStr) error {
	ap.c = make(chan common.DataInter, 1)

	go func() {
		for req := range ap.c {
			err := ap.filterRequest(req)
			if nil != err {
				log.Debug("error on filtering request: %s", err)
			}
		}
	}()
	go ap.setupTimer(out)

	return nil
}

func (ap *apdex) setupTimer(out chan<- common.DataStr) {
	t := time.Now()
	d := time.Second * 30
	nextTime := t.Truncate(d).Add(d)
	select {
	case <-time.After(nextTime.Sub(t)):
		ap.alert(out)
		ticker := time.NewTicker(d)
		for range ticker.C {
			ap.alert(out)
		}
	}
}

func (*apdex) alert(out chan<- common.DataStr) {
	queue := pop()
	for _, ar := range queue {
		ar.SetApdex()
		out <- ar.toDataStr()
	}
}

func (ap *apdex) filterRequest(req common.DataInter) error {
	var b bytes.Buffer
	al := Accesslog{}

	b.Write([]byte(req.(string)))
	json.Unmarshal(b.Bytes(), &al)

	host := al.Host
	requestTime := al.RequestTime * 1000

	rule, err := ap.fetchRule(host)
	if nil != err {
		return nil
	}

	idx := int(time.Since(al.Time.Truncate(time.Second*30)).Seconds()/30) ^ 1
	if idx < 0 || idx > 1 {
		return ErrInvalidIndex
	}

	ar, err := fetchResult(idx, host)
	if ErrHostNotExists == err {
		ar = NewApdexResult(host, al.Time.Truncate(time.Second*30))
	}

	ar = updateResult(rule, ar, requestTime)
	updateQueue(idx, ar)
	return nil
}

func (ap *apdex) fetchRule(host string) (ApdexRule, error) {
	for _, rule := range ap.config.Rules {
		if rule.Host == host {
			return rule, nil
		}
	}
	return ApdexRule{}, ErrRuleMismatch
}

func updateResult(rule ApdexRule, ar *ApdexResult, requestTime float64) *ApdexResult {
	switch {
	case requestTime <= rule.Satisfied:
		ar.SatisfiedCount++
	case rule.Satisfied < requestTime && requestTime <= rule.Tolerating:
		ar.ToleratingCount++
	case requestTime > rule.Tolerating:
		ar.FrustratedCount++
	}
	return ar
}
