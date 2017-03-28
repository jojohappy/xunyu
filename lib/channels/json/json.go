package json

import (
	"bytes"
	"encoding/json"

	"github.com/xunyu/common"
	"github.com/xunyu/config"
)

type jsonFilter struct {
	common.PluginPrototype
	c chan common.DataInter
}

func init() {
	common.RegisterChannelPlugin("json", New)
}

func New(_ *config.Config) (common.Pluginer, error) {
	return &jsonFilter{}, nil
}

func (j *jsonFilter) Filter(out chan<- common.DataStr) error {
	j.c = make(chan common.DataInter, 1)
	go func() {
		for {
			select {
			case data := <-j.c:
				out <- j.filterJson(data)
			}
		}
	}()
	return nil
}

func (j *jsonFilter) filterJson(data common.DataInter) common.DataStr {
	var b bytes.Buffer
	var ds map[string]interface{}

	b.Write([]byte(data.(string)))
	json.Unmarshal(b.Bytes(), &ds)

	return common.DataStr(ds)
}

func (j *jsonFilter) GetFilterChannel() chan<- common.DataInter {
	return j.c
}
