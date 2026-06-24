package xjson

import "encoding/json"

func Null() *Value {
	return &Value{
		kind:  KindNull,
		value: (*struct{})(nil),
	}
}

func Bool(boolValue bool) *Value {
	return &Value{
		kind:  KindBool,
		value: boolValue,
	}
}

func Number(numberValue float64) *Value {
	return &Value{
		kind:  KindNumber,
		value: numberValue,
	}
}

func String(stringValue string) *Value {
	return &Value{
		kind:  KindString,
		value: stringValue,
	}
}

func Array(arrayValue []*Value) *Value {
	return &Value{
		kind:  KindArray,
		value: arrayValue,
	}
}

func Object(objectValue map[string]*Value) *Value {
	return &Value{
		kind:  KindObject,
		value: objectValue,
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
