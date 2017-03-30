package config

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	ErrTypeMismatch = errors.New("type mismatch")
)

type value interface {
	toConfig() (*Config, error)

	toBool() (bool, error)
	toString() (string, error)
	toInt() (int64, error)
	toUint() (uint64, error)
	toFloat() (float64, error)
}

type cfgBool struct {
	b bool
}

type cfgInt struct {
	i int64
}

type cfgUint struct {
	u uint64
}

type cfgFloat struct {
	f float64
}

type cfgString struct {
	s string
}

type cfgSub struct {
	cfg *Config
}

func newBool(b bool) *cfgBool {
	return &cfgBool{b}
}

func newInt(i int64) *cfgInt {
	return &cfgInt{i}
}

func newUint(u uint64) *cfgUint {
	return &cfgUint{u}
}

func newFloat(f float64) *cfgFloat {
	return &cfgFloat{f}
}

func newString(s string) *cfgString {
	return &cfgString{s}
}

func newCfgSub(c *Config) *cfgSub {
	return &cfgSub{c}
}

func (cfgBool) toConfig() (*Config, error)   { return nil, ErrTypeMismatch }
func (c *cfgBool) toBool() (bool, error)     { return c.b, nil }
func (c *cfgBool) toString() (string, error) { return fmt.Sprintf("%t", c.b), nil }
func (c *cfgBool) toInt() (int64, error)     { return 0, ErrTypeMismatch }
func (c *cfgBool) toUint() (uint64, error)   { return 0, ErrTypeMismatch }
func (c *cfgBool) toFloat() (float64, error) { return 0, ErrTypeMismatch }

func (cfgInt) toConfig() (*Config, error)   { return nil, ErrTypeMismatch }
func (cfgInt) toBool() (bool, error)        { return false, ErrTypeMismatch }
func (c *cfgInt) toString() (string, error) { return fmt.Sprintf("%d", c.i), nil }
func (c *cfgInt) toInt() (int64, error)     { return c.i, nil }
func (c *cfgInt) toUint() (uint64, error)   { return uint64(c.i), nil }
func (c *cfgInt) toFloat() (float64, error) { return float64(c.i), nil }

func (cfgUint) toConfig() (*Config, error)   { return nil, ErrTypeMismatch }
func (cfgUint) toBool() (bool, error)        { return false, ErrTypeMismatch }
func (c *cfgUint) toString() (string, error) { return fmt.Sprintf("%d", c.u), nil }
func (c *cfgUint) toInt() (int64, error)     { return int64(c.u), nil }
func (c *cfgUint) toUint() (uint64, error)   { return c.u, nil }
func (c *cfgUint) toFloat() (float64, error) { return float64(c.u), nil }

func (cfgFloat) toConfig() (*Config, error)   { return nil, ErrTypeMismatch }
func (cfgFloat) toBool() (bool, error)        { return false, ErrTypeMismatch }
func (c *cfgFloat) toString() (string, error) { return fmt.Sprintf("%v", c.f), nil }
func (c *cfgFloat) toInt() (int64, error)     { return int64(c.f), nil }
func (c *cfgFloat) toUint() (uint64, error)   { return uint64(c.f), nil }
func (c *cfgFloat) toFloat() (float64, error) { return c.f, nil }

func (cfgString) toConfig() (*Config, error)   { return nil, ErrTypeMismatch }
func (c *cfgString) toBool() (bool, error)     { return strconv.ParseBool(c.s) }
func (c *cfgString) toString() (string, error) { return c.s, nil }
func (c *cfgString) toInt() (int64, error)     { return strconv.ParseInt(c.s, 0, 64) }
func (c *cfgString) toUint() (uint64, error)   { return strconv.ParseUint(c.s, 0, 64) }
func (c *cfgString) toFloat() (float64, error) { return strconv.ParseFloat(c.s, 64) }

func (c cfgSub) toConfig() (*Config, error) { return c.cfg, nil }
func (cfgSub) toBool() (bool, error)        { return false, ErrTypeMismatch }
func (cfgSub) toString() (string, error)    { return "", nil }
func (cfgSub) toInt() (int64, error)        { return 0, ErrTypeMismatch }
func (cfgSub) toUint() (uint64, error)      { return 0, ErrTypeMismatch }
func (cfgSub) toFloat() (float64, error)    { return 0, ErrTypeMismatch }
