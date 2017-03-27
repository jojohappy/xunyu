package main

import (
    "fmt"
    "os"

    "github.com/xunyu/lib/xunyu"
)

func main() {
    if err := xunyu.Run(); nil != err {
        fmt.Println(err)
        os.Exit(1)
    }
}
