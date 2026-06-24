package xjson

import (
	"bytes"
	"encoding/json"
	"fmt"
)

var _ interface {
	json.Marshaler
	json.Unmarshaler
} = (*Value)(nil)

type Value struct {
	kind  Kind
	value any
}

func (v *Value) GetKind() Kind {
	return v.kind
}

func (v *Value) GetNull() *struct{} {
	return v.value.(*struct{})
}

func (v *Value) GetBool() bool {
	return v.value.(bool)
}

func (v *Value) GetNumber() float64 {
	return v.value.(float64)
}

func (v *Value) GetString() string {
	return v.value.(string)
}

func (v *Value) GetArray() []*Value {
	return v.value.([]*Value)
}

func (v *Value) GetObject() map[string]*Value {
	return v.value.(map[string]*Value)
}

func (v *Value) IsNull() bool {
	return v.kind == KindNull
}

func (v *Value) IsBool() bool {
	return v.kind == KindBool
}

func (v *Value) IsNumber() bool {
	return v.kind == KindNumber
}

func (v *Value) IsString() bool {
	return v.kind == KindString
}

func (v *Value) IsArray() bool {
	return v.kind == KindArray
}

func (v *Value) IsObject() bool {
	return v.kind == KindObject
}

func (v *Value) AsNull() (*struct{}, bool) {
	if v.kind == KindNull {
		return v.value.(*struct{}), true
	}
	return nil, false
}

func (v *Value) AsBool() (bool, bool) {
	if v.kind == KindBool {
		return v.value.(bool), true
	}
	return false, false
}

func (v *Value) AsNumber() (float64, bool) {
	if v.kind == KindNumber {
		return v.value.(float64), true
	}
	return 0, false
}

func (v *Value) AsString() (string, bool) {
	if v.kind == KindString {
		return v.value.(string), true
	}
	return "", false
}

func (v *Value) AsArray() ([]*Value, bool) {
	if v.kind == KindArray {
		return v.value.([]*Value), true
	}
	return nil, false
}

func (v *Value) AsObject() (map[string]*Value, bool) {
	if v.kind == KindObject {
		return v.value.(map[string]*Value), true
	}
	return nil, false
}

func (v *Value) MarshalJSON() ([]byte, error) {
	switch v.kind {
	case KindNull:
		return json.Marshal(v.value.(*struct{}))

	case KindBool:
		return json.Marshal(v.value.(bool))

	case KindNumber:
		return json.Marshal(v.value.(float64))

	case KindString:
		return json.Marshal(v.value.(string))

	case KindArray:
		return json.Marshal(v.value.([]*Value))

	case KindObject:
		return json.Marshal(v.value.(map[string]*Value))
	}
	return nil, nil
}

func (v *Value) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)

	if len(data) == 0 {
		return fmt.Errorf("xjson: call of xjson.Value.UnmarshalJSON with empty data")
	}

	switch {
	case string(data) == "null":
		v.kind = KindNull
		var nullValue *struct{}
		v.value = nullValue

	case data[0] == 't' || data[0] == 'f':
		v.kind = KindBool
		var boolValue bool
		if err := json.Unmarshal(data, &boolValue); err != nil {
			return err
		}
		v.value = boolValue

	case data[0] == '"':
		v.kind = KindString
		var stringValue string
		if err := json.Unmarshal(data, &stringValue); err != nil {
			return err
		}
		v.value = stringValue

	case data[0] == '[':
		v.kind = KindArray
		var messages []json.RawMessage
		if err := json.Unmarshal(data, &messages); err != nil {
			return err
		}
		arrayValue := make([]*Value, len(messages))
		for index, message := range messages {
			child := &Value{}
			if err := child.UnmarshalJSON(message); err != nil {
				return err
			}
			arrayValue[index] = child
		}
		v.value = arrayValue

	case data[0] == '{':
		v.kind = KindObject
		var messages map[string]json.RawMessage
		if err := json.Unmarshal(data, &messages); err != nil {
			return err
		}
		objectValue := make(map[string]*Value, len(messages))
		for key, message := range messages {
			child := &Value{}
			if err := child.UnmarshalJSON(message); err != nil {
				return err
			}
			objectValue[key] = child
		}
		v.value = objectValue

	default:
		v.kind = KindNumber
		var numberValue float64
		if err := json.Unmarshal(data, &numberValue); err != nil {
			return err
		}
		v.value = numberValue
	}

	return nil
}
