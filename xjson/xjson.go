package xjson

import "encoding/json"

func Null() *Value {
	return &Value{
		kind:  KindNull,
		value: (*struct{})(nil),
	}
}

func Bool(value bool) *Value {
	return &Value{
		kind:  KindBool,
		value: value,
	}
}

func Number(value float64) *Value {
	return &Value{
		kind:  KindNumber,
		value: value,
	}
}

func String(value string) *Value {
	return &Value{
		kind:  KindString,
		value: value,
	}
}

func Array(value []*Value) *Value {
	return &Value{
		kind:  KindArray,
		value: value,
	}
}

func Object(value map[string]*Value) *Value {
	return &Value{
		kind:  KindObject,
		value: value,
	}
}

func Marshal(v *Value) ([]byte, error) {
	return json.Marshal(v)
}

func MarshalIndent(v *Value, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

func Unmarshal(data []byte) (*Value, error) {
	v := &Value{}
	if err := v.UnmarshalJSON(data); err != nil {
		return nil, err
	}
	return v, nil
}
