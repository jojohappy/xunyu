package xunyu

import (
    "fmt"
    "os"
    "os/signal"
    "sync"
    "syscall"
)

func HandleSignals(stopFunc func()) {
    var once sync.Once

    sigc := make(chan os.Signal, 1)
    signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        <-sigc
        fmt.Println("received sigterm/sigin, stopping")
        once.Do(stopFunc)
    }()
}
