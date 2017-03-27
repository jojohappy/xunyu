package plugins

import (
    "github.com/xunyu/config"
    "github.com/xunyu/common"

	// load support channels
	_ "github.com/xunyu/lib/channels/json"

	// load support inputs
	_ "github.com/xunyu/lib/inputs/file"

	// load support outputs
	// _ "github.com/xunyu/outputs/elasticsearch"
	_ "github.com/xunyu/lib/outputs/console"
)

func LoadPlugins(
    inputsConfigs map[string]*config.Config,
    outputsConfigs map[string]*config.Config,
    channelsConfigs map[string]*config.Config,
) (*common.Plugins, error) {
    configs := map[string]map[string]*config.Config{
        "input": inputsConfigs,
        "output": outputsConfigs,
        "channel": channelsConfigs,
    }
    return common.InitPlugin(configs)
}
