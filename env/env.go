package env

import (
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/go-playground/validator/v10"
)

const StructTag = "env"

func Bind(v interface{}) error {
	if v == nil {
		return fmt.Errorf("env: cannot bind nil")
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer {
		return fmt.Errorf("env: cannot bind non pointer %s", rv.Kind())
	}
	if rv.IsNil() {
		return fmt.Errorf("env: cannot bind nil pointer %s", rv.Kind())
	}
	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return fmt.Errorf("env: cannot bind non struct pointer %s", rv.Kind())
	}
	rt := rv.Type()
	for i := 0; i < rv.NumField(); i++ {
		sf := rt.Field(i)
		tag, ok := sf.Tag.Lookup(StructTag)
		if !ok {
			continue
		}
		if tag == "" {
			return fmt.Errorf("env: field %q has empty %s tag", sf.Name, StructTag)
		}
		field := rv.Field(i)
		if !field.CanSet() {
			return fmt.Errorf("env: field %q cannot be set", sf.Name)
		}
		value := os.Getenv(tag)
		if value == "" {
			continue
		}
		switch field.Kind() {
		case reflect.String:
			field.SetString(value)
		case reflect.Bool:
			b, err := strconv.ParseBool(value)
			if err != nil {
				return err
			}
			field.SetBool(b)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			i, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return err
			}
			field.SetInt(i)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			u, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return err
			}
			field.SetUint(u)
		case reflect.Float32, reflect.Float64:
			f, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return err
			}
			field.SetFloat(f)
		default:
			return fmt.Errorf("env: field %q must be of type (bool|string|number)", sf.Name)
		}
	}
	return validator.New().Struct(v)
}
