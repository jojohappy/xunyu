package zabbix

import (
	"fmt"
	"testing"
	"time"

	"github.com/xunyu/common"
	"github.com/xunyu/config"
)

type node map[string]interface{}

var (
	data = []common.DataStr{
		common.DataStr{
			"Key":   "item1",
			"Time":  time.Now().Truncate(time.Second * 30),
			"Value": 0.97,
		},
		common.DataStr{
			"Key":   "item2",
			"Time":  time.Now(),
			"Value": 0.95,
		},
	}
)

func genConfig(key string) (*config.Config, error) {
	cfg, err := config.From(node{
		"host":         "zabbix server",
		"key":          key,
		"server":       "localhost",
		"target_value": "Value",
		"target_clock": "Time",
		"target_key":   "Key",
	})

	return cfg, err
}

func TestZabbixSenderFaild(t *testing.T) {
	cfg, _ := genConfig("xunyu.test")
	testSender(cfg, t)
}

func TestZabbixSenderCorrect(t *testing.T) {
	cfg, _ := genConfig("xunyu.test[%s]")
	testSender(cfg, t)
}

func testSender(cfg *config.Config, t *testing.T) {
	z := &zabbix{config: defaultConfig}
	err := cfg.Assemble(&z.config)
	if nil != err {
		t.Error(err)
		t.FailNow()
	}

	for _, d := range data {
		items, err := z.makeItems(d)
		if nil != err {
			t.Error(err)
		}
		t.Log(items)

		server := fmt.Sprintf("%s:%d", z.config.Server, z.config.Port)
		res, err := send(server, items)
		if nil != err {
			t.Error(err)
		}
		t.Logf("response from zabbix: %s, processed: %d, failed: %d", res.Response, res.Processed, res.Failed)
	}
}
