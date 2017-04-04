package file

import (
	"fmt"

	"github.com/xunyu/common"
	"github.com/xunyu/config"

	"github.com/hpcloud/tail"
)

type FileConfig struct {
	Path   string `config:"path"`
	ReOpen bool   `config:"reOpen"`
	Follow bool   `config:"follow"`
}

var (
	defaultConfig = FileConfig{
		ReOpen: false,
		Follow: true,
	}
)

type file struct {
	common.PluginPrototype
	done   chan struct{}
	config FileConfig
}

func init() {
	common.RegisterInputPlugin("file", New)
}

func New(config *config.Config) (common.Pluginer, error) {
	f := &file{
		config: defaultConfig,
		done:   make(chan struct{}),
	}
	if err := f.init(config); nil != err {
		return nil, err
	}

	return f, nil
}

func (f *file) init(config *config.Config) error {
	if err := config.Assemble(&f.config); nil != err {
		return err
	}

	if _, err := f.newTailConfig(); nil != err {
		return err
	}

	return nil
}

func (f *file) newTailConfig() (tail.Config, error) {
	cfg := tail.Config{}
	cfg.Follow = f.config.Follow
	cfg.ReOpen = f.config.ReOpen

	return cfg, nil
}

func (f *file) Start() <-chan common.DataInter {
	out := make(chan common.DataInter, 1)
	cfg, err := f.newTailConfig()
	if nil != err {
		fmt.Println(err)
		close(out)
		return out
	}

	t, err := tail.TailFile(f.config.Path, cfg)
	if nil != err {
		fmt.Println(err)
		close(out)
		return out
	}

	go func() {
		defer close(out)
		for {
			select {
			case line := <-t.Lines:
				out <- line.Text
			case <-f.done:
				return
			}
		}
	}()
	return out
}

func (f *file) Close() {
	close(f.done)
}
