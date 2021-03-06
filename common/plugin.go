package common

import (
	"github.com/xunyu/config"
	"github.com/xunyu/lib/log"
)

type Pluginer interface {
	Close()
	Start() <-chan DataInter
	Output(data DataStr) error
	Filter(out chan<- DataStr) error
	GetFilterChannel() chan<- DataInter
}

type PluginBuilder func(config *config.Config) (Pluginer, error)

type Plugin struct {
	Name   string
	Config *config.Config
	Plugin Pluginer
}

type Plugins struct {
	Inputs   []Plugin
	Outputs  []Plugin
	Channels []Plugin
}

var (
	inputsPlugins   = make(map[string]PluginBuilder)
	outputsPlugins  = make(map[string]PluginBuilder)
	channelsPlugins = make(map[string]PluginBuilder)

	catagorys = []string{"input", "output", "channel"}
)

func RegisterOutputPlugin(name string, builder PluginBuilder) {
	outputsPlugins[name] = builder
}

func RegisterInputPlugin(name string, builder PluginBuilder) {
	inputsPlugins[name] = builder
}

func RegisterChannelPlugin(name string, builder PluginBuilder) {
	channelsPlugins[name] = builder
}

func InitPlugin(
	configs map[string]map[string]*config.Config,
) (*Plugins, error) {
	plugins := &Plugins{}
	for _, name := range catagorys {
		pb := getPlugins(name)
		p, err := doInitPlugin(name, pb, configs[name])
		if nil != err {
			return nil, err
		}
		plugins.setPlugins(name, p)
	}
	return plugins, nil
}

func doInitPlugin(
	catagory string,
	pb map[string]PluginBuilder,
	cfg map[string]*config.Config,
) (plugins []Plugin, err error) {
	for name, builder := range pb {
		c, ok := cfg[name]
		if !ok {
			continue
		}
		log.Info("load %s plugin %s", catagory, name)
		p, err := builder(c)
		if nil != err {
			return nil, err
		}
		plugins = append(plugins, Plugin{Name: name, Config: c, Plugin: p})
	}

	return plugins, nil
}

func getPlugins(name string) map[string]PluginBuilder {
	switch name {
	case "input":
		return inputsPlugins
	case "output":
		return outputsPlugins
	case "channel":
		return channelsPlugins
	}
	return nil
}

func (plugins *Plugins) setPlugins(name string, p []Plugin) {
	switch name {
	case "input":
		plugins.Inputs = p
	case "output":
		plugins.Outputs = p
	case "channel":
		plugins.Channels = p
	}
}

type PluginPrototype struct{}

func (*PluginPrototype) Start() <-chan DataInter {
	return nil
}

func (*PluginPrototype) Close() {}

func (*PluginPrototype) Output(data DataStr) error {
	return nil
}

func (*PluginPrototype) Filter(out chan<- DataStr) error {
	return nil
}

func (*PluginPrototype) GetFilterChannel() chan<- DataInter {
	return nil
}
