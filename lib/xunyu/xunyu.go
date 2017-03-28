package xunyu

import (
	"fmt"
	"sync"

	"github.com/xunyu/common"
	"github.com/xunyu/config"
	"github.com/xunyu/lib/plugins"
)

type Xunyu struct {
	Path    string
	Config  XunyuConfig
	Version string
	Plugins *common.Plugins
}

type XunyuConfig struct {
	Inputs   map[string]*config.Config `config:"inputs"`
	Outputs  map[string]*config.Config `config:"outputs"`
	Channels map[string]*config.Config `config:"channels"`
}

func Run() error {
	xy := newXunyu("")

	if err := xy.configure(); err != nil {
		return err
	}

	err := xy.init()
	if nil != err {
		return err
	}

	return xy.Run()
}

func newXunyu(version string) *Xunyu {
	if version == "" {
		version = defaultXunyuVersion
	}

	return &Xunyu{
		Version: version,
	}
}

func (xy *Xunyu) init() error {
	p, err := plugins.LoadPlugins(xy.Config.Inputs, xy.Config.Outputs, xy.Config.Channels)

	if nil != err {
		return err
	}

	xy.Plugins = p
	return nil
}

func (xy *Xunyu) configure() error {
	xy.Path = "/Users/michael/works/go/src/github.com/xunyu/config.json"
	cfg, err := config.Load(xy.Path)
	if nil != err {
		return err
	}

	if err := cfg.Assemble(&xy.Config); nil != err {
		return err
	}

	return nil
}

func (xy *Xunyu) Run() error {
	in := runInput(xy.Plugins.Inputs)
	ch := runChannel(xy.Plugins.Channels, in)
	runOutput(xy.Plugins.Outputs, ch)
	return nil
}

func runInput(inputs []common.Plugin) <-chan common.DataInter {
	fmt.Println("Starting Input")

	out := make(chan common.DataInter, 1)

	var wg sync.WaitGroup
	wg.Add(len(inputs))
	for _, p := range inputs {
		go func(p common.Plugin) {
			defer wg.Done()
			o := p.Plugin.Start()
			for data := range o {
				out <- data
			}
		}(p)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func runChannel(channels []common.Plugin, in <-chan common.DataInter) <-chan common.DataStr {
	fmt.Println("Starting Channel")

	out := make(chan common.DataStr, 1)
	var wg sync.WaitGroup

	filter := func(p common.Plugin) {
		defer wg.Done()
		p.Plugin.Filter(out)
		for data := range in {
			p.Plugin.GetFilterChannel() <- data
		}
	}

	wg.Add(len(channels))
	for _, p := range channels {
		go filter(p)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func runOutput(outputs []common.Plugin, cs <-chan common.DataStr) {
	fmt.Println("Starting Output")
	defer fmt.Println("Stopped Output")

	var wg sync.WaitGroup

	output := func(p common.Plugin) {
		defer wg.Done()
		for data := range cs {
			p.Plugin.Output(data)
		}
	}

	wg.Add(len(outputs))
	for _, p := range outputs {
		go output(p)
	}

	wg.Wait()
}
