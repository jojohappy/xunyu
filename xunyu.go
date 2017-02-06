package main

import (
	"bufio"
	"os"

	"github.com/xunyu/compute"
	"github.com/xunyu/sender"
)

func gen() <-chan string {
	out := make(chan string, 0)
	pwd, _ := os.Getwd()
	go func() {
		f, _ := os.Open(pwd + "/data.json")
		r := bufio.NewReader(f)
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			out <- scanner.Text()
		}
		close(out)
	}()
	return out
}

func main() {
	in := gen()
	out := compute.Resolve(in)
	sender.Send(out)
}
