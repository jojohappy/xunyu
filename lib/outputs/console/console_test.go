package console

import (
    "bytes"
    "io"
    "os"
    "testing"

    "github.com/xunyu/common"
)

func withStdout(fn func()) (string, error) {
    stdout := os.Stdout

    r, w, err := os.Pipe()
    if err != nil {
        return "", err
    }

    os.Stdout = w
    defer func() {
        os.Stdout = stdout
    }()

    outC := make(chan string)
    go func() {
        var buf bytes.Buffer
        _, err = io.Copy(&buf, r)
        r.Close()
        outC <- buf.String()
    }()

    fn()
    w.Close()
    result := <-outC
    return result, err
}

func run(cfg Config, data ...common.DataStr) (string, error) {
    return withStdout(func() {
        c := &console{out: os.Stdout, config: cfg}
        for _, d := range data {
            c.Output(d)
        }
    })
}

func TestConsoleOneData(t *testing.T) {
    lines, err := run(defaultConfig, common.DataStr{"data":"console"})
    if nil != err {
        t.Error(err)
        t.Fail()
    }
    t.Logf("test console output: %v", lines)
}

func TestConsoleMultipleData(t *testing.T) {
    lines, err := run(Config{Pretty: true},
        common.DataStr{"data":"console1"},
        common.DataStr{"data":"console2"},
        common.DataStr{"data":"console3"},
    )
    if nil != err {
        t.Error(err)
        t.Fail()
    }
    t.Logf("test console output: %v", lines)
}
