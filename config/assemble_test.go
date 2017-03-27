package config

import (
    "testing"
)

type node map[string]interface{}

func TestAssembleConfig(t *testing.T) {
    tests := []interface{}{
        &struct {
            B bool
            I int
            U uint
            F float64
            S string
        }{},
        &struct {
            B *bool
            I *int
            U *uint
            F *float64
            S *string
        }{},
    }

    c, _ := parseConfig(node{
        "b": true,
        "i": 42,
        "u": 23,
        "f": 3.14,
        "s": "string",
    })

    for i, out := range tests {
        t.Logf("test unpack primitives(%v) into: %v", i, out)
        err := c.Assemble(out)
        if err != nil {
            t.Fatalf("failed to unpack: %v", err)
        }
    }
}
