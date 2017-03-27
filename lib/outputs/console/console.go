package console

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/xunyu/common"
	"github.com/xunyu/config"
)

type Config struct {
	Pretty bool `config:"pretty"`
}

type console struct {
	common.PluginPrototype
	out    *os.File
	config Config
}

var (
	defaultConfig = Config{
		Pretty: false,
	}

	nl = []byte{'\n'}
)

func init() {
	common.RegisterOutputPlugin("console", New)
}

func New(config *config.Config) (common.Pluginer, error) {
	c := &console{out: os.Stdout, config: defaultConfig}
	if err := config.Assemble(&c.config); nil != err {
		return nil, err
	}

	return c, nil
}

func (c *console) Output(data common.DataStr) error {
	var err error
	var buf []byte
	if c.config.Pretty {
		buf, err = json.MarshalIndent(data, "", "  ")
	} else {
		buf, err = json.Marshal(data)
	}

	if err != nil {
		return fmt.Errorf("Fail to convert the event to JSON (%v): %#v", err, data)
	}

	if err = c.writeBuffer(buf); err != nil {
		return err
	}
	if err = c.writeBuffer(nl); err != nil {
		return err
	}

	return nil
}

func (c *console) writeBuffer(buf []byte) error {
	written := 0
	for written < len(buf) {
		n, err := c.out.Write(buf[written:])
		if err != nil {
			return err
		}

		written += n
	}
	return nil
}
