package json

import (
	"fmt"
	"testing"

	"github.com/xunyu/common"
	"github.com/xunyu/config"
)

func TestJsonFilter(t *testing.T) {
	done := make(chan bool)
	in := genInput()
	out := filter(in)

	go func() {
		for data := range out {
			t.Logf("test json filter message: %v", data)
		}
		done <- true
	}()
	<-done
}

func genInput() <-chan common.DataInter {
	in := make(chan common.DataInter)
	go func() {
		for j := 1; j <= 3; j++ {
			in <- "{\"real_ip\":\"-\",\"remote\":\"218.26.54.2\",\"port\":\"6232\"}"
		}
		close(in)
	}()
	return in
}

func genBenchmarkInput(size int) <-chan common.DataInter {
	in := make(chan common.DataInter)
	go func() {
		for j := 0; j < size; j++ {
			msg := fmt.Sprintf("message%d", j)
			in <- "{\"json\":\"" + msg + "\"}"
		}
		close(in)
	}()
	return in
}

func filter(in <-chan common.DataInter) <-chan common.DataStr {
	js, _ := New(&config.Config{})
	out := make(chan common.DataStr)

	go func(j common.Pluginer) {
		j.Filter(out)
		fc := j.GetFilterChannel()
		for data := range in {
			fc <- data
		}
		close(fc)
	}(js)

	return out
}

func runBenchmarkLoops(size int, b *testing.B) {
	for i := 0; i < b.N; i++ {
		runBenchmark(size, b)
	}
}

func runBenchmark(size int, b *testing.B) {
	count := 0

	done := make(chan bool)
	in := genInput()
	out := filter(in)

	go func() {
		for data := range out {
			data["json"] = nil
			count++
		}
		done <- true
	}()
	<-done
}

func BenchmarkJsonFilter1000000(b *testing.B) {
	runBenchmarkLoops(1000000, b)
}

func BenchmarkJsonFilter1000(b *testing.B) {
	runBenchmarkLoops(1000, b)
}
