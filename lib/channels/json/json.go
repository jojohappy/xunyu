package json

import (
	"bytes"
	"encoding/json"
	"sync"

	"github.com/xunyu/common"
	"github.com/xunyu/config"
)

type jsonFilter struct {
	common.PluginPrototype
}

func init() {
	common.RegisterChannelPlugin("json", New)
}

func New(_ *config.Config) (common.Pluginer, error) {
	return &jsonFilter{}, nil
}

func (j *jsonFilter) Filter(cs ...<-chan common.DataInter) <-chan common.DataStr {
	out := make(chan common.DataStr, 1)

	var wg sync.WaitGroup
	filter := func(c <-chan common.DataInter) {
		defer wg.Done()
		j.filterJson(c, out)
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go filter(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func (j *jsonFilter) filterJson(c <-chan common.DataInter, out chan<- common.DataStr) {
	for data := range c {
		var b bytes.Buffer
		var ds map[string]interface{}

		b.Write([]byte(data.(string)))
		json.Unmarshal(b.Bytes(), &ds)

		out <- common.DataStr(ds)
	}
}
