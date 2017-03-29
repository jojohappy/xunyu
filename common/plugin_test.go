package common

import (
    "testing"

    "github.com/xunyu/config"
)

type node map[string]interface{}

type testPluginer struct {
    PluginPrototype
}

type pluginConfig struct {
    Inputs   map[string]*config.Config `config:"inputs"`
    Outputs  map[string]*config.Config `config:"outputs"`
    Channels map[string]*config.Config `config:"channels"`
}

func createTestPluginer(config *config.Config) (Pluginer, error) {
    return &testPluginer{}, nil
}

func TestRegisterPlugin(t *testing.T) {
    RegisterOutputPlugin("test-plugin", createTestPluginer)

    cfg, _ := config.From(node{
        "outputs": node{
            "test-plugin": node{},
        },
    })

    pConfig := &pluginConfig{}
    cfg.Assemble(&pConfig)
    configs := map[string]map[string]*config.Config{
        "input":   pConfig.Inputs,
        "output":  pConfig.Outputs,
        "channel": pConfig.Channels,
    }
    ps, err := InitPlugin(configs)
    if nil != err {
        t.Fatalf("failed to load plugin: %v", err)
    }

    t.Logf("number of out plugin: %v", len(ps.Outputs))
}
