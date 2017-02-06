package sender

import (
	"encoding/json"
	"fmt"
	"github.com/xunyu/compute"
)

func Send(in <-chan interface{}) {
	for {
		for item := range in {
			switch item.(type) {
			case compute.Accesslog:
				d, _ := json.Marshal(item)
				fmt.Println(string(d))
			default:
				fmt.Println("default")
			}
		}
	}
}
