package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
)

type Config struct {
	content map[string]value
}

var (
	tConfig = reflect.TypeOf(Config{})
)

func New() *Config {
	return &Config{}
}

func Load(path string) (*Config, error) {
	var c interface{}
	in, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(in, &c)
	if nil != err {
		return nil, err
	}

	return parseConfig(c)
}

func From(from interface{}) (*Config, error) {
	return parseConfig(from)
}

func parseConfig(from interface{}) (*Config, error) {
	reflectFrom := reflect.ValueOf(from)
	switch t := reflectFrom.Type(); t {
	case reflect.TypeOf((map[string]interface{})(nil)):
		return parseMap(reflectFrom)
	default:
		switch reflectFrom.Kind() {
		case reflect.Map:
			return parseMap(reflectFrom)
		default:
			fmt.Println("error type")
		}
	}
	return nil, nil
}

func parseMap(from reflect.Value) (*Config, error) {
	cfg := New()
	for _, k := range from.MapKeys() {
		field := k.String()
		v, err := parseValue(from.MapIndex(k))
		if nil != err {
			return nil, err
		}
		if nil == cfg.content {
			cfg.content = map[string]value{}
		}
		cfg.content[field] = v
	}
	return cfg, nil
}

func parseMapValue(from reflect.Value) (value, error) {
	sub, err := parseMap(from)
	if err != nil {
		return nil, err
	}
	v := cfgSub{sub}
	return v, nil
}

func parseValueInterfaces(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Interface && !v.IsNil() {
		v = v.Elem()
	}
	return v
}

func parseValuePoniters(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	return v
}

func parseTypePointers(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

func parseValuePointers(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	return v
}

func parseValue(v reflect.Value) (value, error) {
	v = parseValueInterfaces(v)
	switch v.Kind() {
	case reflect.Bool:
		return newBool(v.Bool()), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i := v.Int()
		if i > 0 {
			return newUint(uint64(i)), nil
		}
		return newInt(i), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return newUint(v.Uint()), nil
	case reflect.Float32, reflect.Float64:
		f := v.Float()
		return newFloat(f), nil
	case reflect.String:
		return newString(v.String()), nil
	// case reflect.Array, reflect.Slice:
	//     return normalizeArray(v)
	case reflect.Map:
		return parseMapValue(v)
	}
	return nil, nil
}

func (c *Config) GetValue(name string) (value, error) {
	if val, ok := c.content[name]; ok {
		return val, nil
	} else {
		return nil, errors.New("missing field")
	}
}
