package file

import (
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
	config FileConfig
}

func init() {
	common.RegisterInputPlugin("file", New)
}

func New(config *config.Config) (common.Pluginer, error) {
	f := &file{config: defaultConfig}
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
		close(out)
		return out
	}

	t, err := tail.TailFile(f.config.Path, cfg)

	go func() {
		for line := range t.Lines {
			out <- line.Text
		}

		close(out)
	}()
	return out
}
