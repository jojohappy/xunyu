package plugins

import (
	"testing"

	"github.com/xunyu/common"
	"github.com/xunyu/config"
)

type node map[string]interface{}

type testPluginer struct {
	common.PluginPrototype
}

type pluginConfig struct {
	Inputs   map[string]*config.Config `config:"inputs"`
	Outputs  map[string]*config.Config `config:"outputs"`
	Channels map[string]*config.Config `config:"channels"`
}

func createTestInputPluginer(config *config.Config) (common.Pluginer, error) {
	return &testPluginer{}, nil
}

func createTestOutputPluginer(config *config.Config) (common.Pluginer, error) {
	return &testPluginer{}, nil
}

func TestRegisterPlugin(t *testing.T) {
	common.RegisterInputPlugin("test-input-plugin", createTestInputPluginer)
	common.RegisterOutputPlugin("test-output-plugin", createTestOutputPluginer)

	cfg, _ := config.From(node{
		"inputs": node{
			"test-input-plugin": node{},
		},
		"outputs": node{
			"test-output-plugin": node{},
		},
	})

	pConfig := &pluginConfig{}
	cfg.Assemble(&pConfig)

	ps, err := LoadPlugins(pConfig.Inputs, pConfig.Outputs, pConfig.Channels)
	if nil != err {
		t.Fatalf("failed to load plugin: %v", err)
	}

	t.Logf("number of out plugin: %v", len(ps.Outputs))
	t.Logf("number of in plugin: %v", len(ps.Inputs))
}
