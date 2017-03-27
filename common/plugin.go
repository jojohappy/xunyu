package common

import (
    "github.com/xunyu/config"
)

type Pluginer interface {
    Close() error
    Start() <-chan DataStr
    Output(data DataStr) error
    Filter(in ...<-chan DataStr) <-chan DataStr
}

type PluginBuilder func(config *config.Config) (Pluginer, error)

type Plugin struct {
    Name   string
    Config *config.Config
    Plugin Pluginer
}

type Plugins struct {
    Inputs []Plugin
    Outputs []Plugin
    Channels []Plugin
}

var (
    inputsPlugins = make(map[string]PluginBuilder)
    outputsPlugins = make(map[string]PluginBuilder)
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
        p ,err := doInitPlugin(pb, configs[name])
        if nil != err {
            return nil, err
        }
        plugins.setPlugins(name, p)
    }
    return plugins, nil
}

func doInitPlugin(
    pb map[string]PluginBuilder,
    cfg map[string]*config.Config,
) (plugins []Plugin, err error) {
    for name, builder := range pb {
        c, ok := cfg[name]
        if !ok {
            continue
        }
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
