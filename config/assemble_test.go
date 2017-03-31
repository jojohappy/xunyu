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
		t.Logf("test unpack config(%v) into: %v", i, out)
		err := c.Assemble(out)
		if err != nil {
			t.Fatalf("failed to unpack: %v", err)
		}
		t.Logf("test unpack config(%v) into: %v", i, out)
	}
}

func TestAssembleConfigNested(t *testing.T) {
	tests := []interface{}{
		&struct {
			B bool
			S struct {
				I int
				U uint
			}
		}{},
		&struct {
			B *bool
			S *struct {
				I *int
				U *uint
			}
		}{S: &struct {
			I *int
			U *uint
		}{}},
	}

	c, _ := parseConfig(node{
		"b": true,
		"s": node{
			"i": 42,
			"u": 23,
		},
	})

	for i, out := range tests {
		t.Logf("test unpack config nested (%v) into: %v", i, out)
		err := c.Assemble(out)
		if err != nil {
			t.Fatalf("failed to unpack: %v", err)
		}
		t.Logf("test unpack config nested (%v) into: %v", i, out)
	}
}

func TestAssembleConfigSlice(t *testing.T) {
	tests := []interface{}{
		&struct {
			S []struct {
				I int
				U uint
			}
		}{},
	}

	c, _ := parseConfig(node{
		"s": []node{
			node{
				"i": 42,
				"u": 23,
			},
			node{
				"i": 41,
				"u": 24,
			},
		},
	})

	for i, out := range tests {
		t.Logf("test unpack config array (%v) into: %v", i, out)
		err := c.Assemble(out)
		if err != nil {
			t.Fatalf("failed to unpack: %v", err)
		}
		t.Logf("test unpack config array (%v) into: %v", i, out)
	}
}
