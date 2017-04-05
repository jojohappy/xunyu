package xunyu

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/xunyu/lib/log"
)

func HandleSignals(stopFunc func()) {
	var once sync.Once

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigc
		log.Info("received sigterm/sigin, stopping")
		once.Do(stopFunc)
	}()
}
