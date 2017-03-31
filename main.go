package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/xunyu/lib/xunyu"
)

func main() {
	configArgu := flag.String("config", "config.json", "path of config file")
	flag.Parse()
	if err := xunyu.Run(*configArgu); nil != err {
		fmt.Println(err)
		os.Exit(1)
	}
}
