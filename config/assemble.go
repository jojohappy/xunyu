package config

import (
	"reflect"
	"strings"
)

func (c *Config) Assemble(to interface{}) error {
	v := reflect.ValueOf(to)
	v = parseValuePoniters(v)

	switch v.Kind() {
	case reflect.Map:
		return assembleMap(v, c)
	case reflect.Struct:
		return assembleStruct(v, c)
	}
	return nil
}

func assembleStruct(to reflect.Value, cfg *Config) error {
	numFields := to.NumField()
	for i := 0; i < numFields; i++ {
		stField := to.Type().Field(i)
		vField := to.Field(i)
		name := stField.Tag.Get("config")
		fieldName := getFieldName(name, stField.Name)
		val, err := cfg.GetValue(fieldName)
		if nil != err {
			continue
		}

		switch vField.Kind() {
		case reflect.Map:
			sub, _ := val.toConfig()
			assembleMap(vField, sub)
		case reflect.Struct:
			sub, _ := val.toConfig()
			assembleStruct(vField, sub)
		default:
			t := vField.Type()
			v, err := assembleValue(t, val)
			if nil != err {
				return err
			}
			for t != v.Type() {
				if !v.CanAddr() {
					tmp := reflect.New(v.Type())
					tmp.Elem().Set(v)
					v = tmp
				} else {
					v = v.Addr()
				}
			}
			vField.Set(v)
		}
	}
	return nil
}

func assembleMap(to reflect.Value, cfg *Config) error {
	fields := cfg.content
	if to.IsNil() {
		to.Set(reflect.MakeMap(to.Type()))
	}
	for key, value := range fields {
		k := reflect.ValueOf(key)
		var v reflect.Value
		v, err := assembleValue(to.Type().Elem(), value)
		if nil != err {
			return err
		}
		to.SetMapIndex(k, v)
	}
	return nil
}

func assembleValue(t reflect.Type, val value) (reflect.Value, error) {
	t = parseTypePointers(t)
	if tConfig.ConvertibleTo(t) {
		cfg, err := val.toConfig()
		if nil != err {
			return reflect.Value{}, err
		}

		v := reflect.ValueOf(cfg).Convert(reflect.PtrTo(t))
		return v, nil
	}

	switch t.Kind() {
	case reflect.Bool:
		b, err := val.toBool()
		if nil != err {
			return reflect.Value{}, nil
		}
		return reflect.ValueOf(b), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := val.toInt()
		if nil != err {
			return reflect.Value{}, nil
		}
		return reflect.ValueOf(i).Convert(t), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u, err := val.toUint()
		if nil != err {
			return reflect.Value{}, nil
		}
		return reflect.ValueOf(u).Convert(t), nil
	case reflect.Float32, reflect.Float64:
		f, err := val.toFloat()
		if nil != err {
			return reflect.Value{}, nil
		}
		return reflect.ValueOf(f).Convert(t), nil
	case reflect.String:
		s, err := val.toString()
		if nil != err {
			return reflect.Value{}, nil
		}
		return reflect.ValueOf(s), nil
	}
	return reflect.Value{}, nil
}

func getFieldName(tagName string, fieldName string) string {
	if "" != tagName {
		return tagName
	}
	return strings.ToLower(fieldName)
}
