package zabbix

import (
	"fmt"
	"net"
	"reflect"
	"time"

	zbx "github.com/AlekSi/zabbix-sender"
	"github.com/xunyu/common"
	"github.com/xunyu/config"
	"github.com/xunyu/lib/log"
)

type ZabbixConfig struct {
	Host        string `config:"host"`
	Key         string `config:"key"`
	Server      string `config:"server"`
	Port        int    `config:"port"`
	TargetValue string `config:"target_value"`
	TargetKey   string `config:"target_key"`
	TargetClock string `config:"target_clock"`
}

type zabbix struct {
	common.PluginPrototype
	config ZabbixConfig
}

var (
	defaultConfig = ZabbixConfig{
		Host:        "localhost",
		Key:         "xunyu.ping",
		Server:      "localhost",
		Port:        10051,
		TargetValue: "message",
	}
)

func init() {
	common.RegisterOutputPlugin("zabbix", New)
}

func New(config *config.Config) (common.Pluginer, error) {
	z := &zabbix{config: defaultConfig}
	if err := config.Assemble(&z.config); nil != err {
		return nil, err
	}

	return z, nil
}

func (z *zabbix) Output(data common.DataStr) error {
	items, err := z.makeItems(data)
	if nil != err {
		return err
	}
	server := fmt.Sprintf("%s:%d", z.config.Server, z.config.Port)
	res, err := send(server, items)
	if nil != err {
		return err
	}
	log.Info("response from zabbix: %s, processed: %d, failed: %d", res.Response, res.Processed, res.Failed)
	if res.Response != "success" {
		return fmt.Errorf("send to zabbix failed: %s", res.Info)
	}
	return nil
}

func (z *zabbix) makeItems(data map[string]interface{}) (zbx.DataItems, error) {
	var key_ string
	items := make(zbx.DataItems, 1)

	cfg := z.config

	if v, ok := data[cfg.TargetKey]; cfg.TargetKey != "" && ok {
		key_ = fmt.Sprintf(cfg.Key, v)
	} else {
		key_ = cfg.Key
	}

	clock, err := convertUnixTime(data, cfg.TargetClock)
	if nil != err {
		return nil, err
	}

	val, err := convertValue(data, cfg.TargetValue)
	if nil != err {
		return nil, err
	}

	items[0] = zbx.DataItem{
		Hostname:  cfg.Host,
		Key:       key_,
		Timestamp: clock,
		Value:     val,
	}
	log.Debug("data item: %v", items[0])
	return items, nil
}

func convertUnixTime(data map[string]interface{}, target string) (int64, error) {
	if target != "" {
		if v, ok := data[target]; ok {
			val := reflect.ValueOf(v)
			switch val.Type() {
			case reflect.TypeOf(time.Time{}):
				return val.Interface().(time.Time).Unix(), nil
			}

			switch val.Type().Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				return val.Int(), nil
			case reflect.String:
				t, err := time.Parse(time.RFC3339, val.String())
				if nil != err {
					return -1, err
				}
				return t.Unix(), nil
			default:
				return time.Now().Unix(), nil
			}
		} else {
			return -1, fmt.Errorf("field %s does not exists", target)
		}
	} else {
		return time.Now().Unix(), nil
	}
}

func convertValue(data map[string]interface{}, target string) (string, error) {
	if v, ok := data[target]; ok {
		val := reflect.ValueOf(v)
		t := val.Type()
		switch t {
		case reflect.TypeOf(time.Time{}):
			return val.Interface().(time.Time).Format(time.RFC3339), nil
		}

		switch t.Kind() {
		case reflect.Bool:
			return fmt.Sprintf("%t", val.Bool()), nil
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return fmt.Sprintf("%d", val.Int()), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return fmt.Sprintf("%d", val.Uint()), nil
		case reflect.Float32, reflect.Float64:
			return fmt.Sprintf("%.6f", val.Float()), nil
		case reflect.String:
			return fmt.Sprintf("%s", val.String()), nil
		default:
			return "", fmt.Errorf("unsupport type of value %s", t.Kind())
		}
	} else {
		return "", fmt.Errorf("field %s does not exists", target)
	}
}

func send(server string, items zbx.DataItems) (*zbx.Response, error) {
	addr, err := net.ResolveTCPAddr("tcp", server)
	if nil != err {
		return nil, err
	}

	return zbx.Send(addr, items)
}
