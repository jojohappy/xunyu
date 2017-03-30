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
		val, err := cfg.GetDictValue(fieldName)
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
		case reflect.Slice:
			sub, _ := val.toConfig()
			err := assembleSlice(vField, sub.arr)
			if nil != err {
				return err
			}
		default:
			t := vField.Type()
			v, err := assembleValue(t, val, vField)
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
	fields := cfg.dict
	if to.IsNil() {
		to.Set(reflect.MakeMap(to.Type()))
	}
	for key, value := range fields {
		k := reflect.ValueOf(key)
		old := to.MapIndex(k)
		var v reflect.Value
		v, err := assembleValue(to.Type().Elem(), value, old)
		if nil != err {
			return err
		}
		to.SetMapIndex(k, v)
	}
	return nil
}

func assembleValue(t reflect.Type, val value, old reflect.Value) (reflect.Value, error) {
	t = parseTypePointers(t)
	if tConfig.ConvertibleTo(t) {
		cfg, err := val.toConfig()
		if nil != err {
			return reflect.Value{}, err
		}

		v := reflect.ValueOf(cfg).Convert(reflect.PtrTo(t))
		return v, nil
	}
	old = parseValuePointers(old)

	switch t.Kind() {
	case reflect.Bool:
		b, err := val.toBool()
		if nil != err {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(b), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := val.toInt()
		if nil != err {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(i).Convert(t), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u, err := val.toUint()
		if nil != err {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(u).Convert(t), nil
	case reflect.Float32, reflect.Float64:
		f, err := val.toFloat()
		if nil != err {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(f).Convert(t), nil
	case reflect.String:
		s, err := val.toString()
		if nil != err {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(s), nil
	case reflect.Struct:
		sub, _ := val.toConfig()
		assembleStruct(old, sub)
		return old, nil
	}

	return reflect.ValueOf(nil), nil
}

func getFieldName(tagName string, fieldName string) string {
	if "" != tagName {
		return tagName
	}
	return strings.ToLower(fieldName)
}

func assembleSlice(to reflect.Value, arr []value) error {
	if to.IsNil() {
		to.Set(reflect.MakeSlice(to.Type(), len(arr), len(arr)))
	}

	for i, from := range arr {
		_, err := assembleValue(to.Type().Elem(), from, to.Index(i))
		if err != nil {
			return err
		}
	}
	return nil
}
