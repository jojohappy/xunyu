package main

import (
	"os"

	"github.com/xunyu/lib/xunyu"
)

func main() {
	if err := xunyu.Run(); nil != err {
		os.Exit(1)
	}
}
