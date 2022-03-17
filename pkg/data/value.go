package data

import "google.golang.org/protobuf/types/known/structpb"

//Value represents dynamic field which is friendly for database
type Value struct {
	// which kind of data, can be "null","int","long","float","double","string","bool"
	Kind        string
	IntValue    int32
	LongValue   int64
	FloatValue  float32
	DoubleValue float64
	StringValue string
	BoolValue   bool
	JsonValue   JSONMap
}

const NullKind = "null"
const IntKind = "int"
const LongKind = "long"
const FloatKind = "float"
const DoubleKind = "double"
const StringKind = "string"
const BoolKind = "bool"
const JsonKind = "json"

func (v *Value) GetNativeValue() interface{} {
	switch v.Kind {
	case NullKind:
		return nil
	case IntKind:
		return v.IntValue
	case LongKind:
		return v.LongValue
	case FloatKind:
		return v.FloatValue
	case DoubleKind:
		return v.DoubleValue
	case StringKind:
		return v.StringValue
	case BoolKind:
		return v.BoolValue
	case JsonKind:
		return map[string]interface{}(v.JsonValue)
	default:
		return nil
	}
}

func (v *Value) SetAsNull() {
	v.Kind = NullKind
}

func (v *Value) RunIfNull(f func()) {
	if v.Kind == NullKind {
		f()
	}
}

func (v *Value) SetAsInt(value int32) {
	v.Kind = IntKind
	v.IntValue = value
}

func (v *Value) RunIfInt(f func(v int32)) {
	if v.Kind == IntKind {
		f(v.IntValue)
	}
}

func (v *Value) SetAsLong(value int64) {
	v.Kind = LongKind
	v.LongValue = value
}

func (v *Value) RunIfLong(f func(v int64)) {
	if v.Kind == LongKind {
		f(v.LongValue)
	}
}

func (v *Value) SetAsFloat(value float32) {
	v.Kind = FloatKind
	v.FloatValue = value
}
func (v *Value) RunIfFloat(f func(v float32)) {
	if v.Kind == FloatKind {
		f(v.FloatValue)
	}
}
func (v *Value) SetAsDouble(value float64) {
	v.Kind = DoubleKind
	v.DoubleValue = value
}
func (v *Value) RunIfDouble(f func(v float64)) {
	if v.Kind == DoubleKind {
		f(v.DoubleValue)
	}
}
func (v *Value) SetAsString(value string) {
	v.Kind = StringKind
	v.StringValue = value
}
func (v *Value) RunIfString(f func(v string)) {
	if v.Kind == StringKind {
		f(v.StringValue)
	}
}
func (v *Value) SetAsBool(value bool) {
	v.Kind = BoolKind
	v.BoolValue = value
}
func (v *Value) RunIfBool(f func(v bool)) {
	if v.Kind == BoolKind {
		f(v.BoolValue)
	}
}
func (v *Value) SetAsJson(value map[string]interface{}) {
	v.Kind = JsonKind
	v.JsonValue = value
}
func (v *Value) RunIfJson(f func(v map[string]interface{})) {
	if v.Kind == JsonKind {
		f(v.JsonValue)
	}
}

func (v *Value) ToStructPb() *structpb.Value {
	p, _ := structpb.NewValue(v.GetNativeValue())
	return p
}

func (v *Value) ToDynamicValue() *DynamicValue {
	res := &DynamicValue{}
	v.RunIfInt(func(v int32) {
		res.Value = &DynamicValue_IntValue{IntValue: v}
	})
	v.RunIfLong(func(v int64) {
		res.Value = &DynamicValue_LongValue{LongValue: v}
	})
	v.RunIfFloat(func(v float32) {
		res.Value = &DynamicValue_FloatValue{FloatValue: v}
	})
	v.RunIfDouble(func(v float64) {
		res.Value = &DynamicValue_DoubleValue{DoubleValue: v}
	})
	v.RunIfString(func(v string) {
		res.Value = &DynamicValue_StringValue{StringValue: v}
	})
	v.RunIfBool(func(v bool) {
		res.Value = &DynamicValue_BoolValue{BoolValue: v}
	})
	v.RunIfNull(func() {
		res.Value = &DynamicValue_NullValue{}
	})
	v.RunIfJson(func(v map[string]interface{}) {
		if v == nil {
			res.Value = &DynamicValue_JsonValue{}
		} else {
			j, _ := structpb.NewStruct(v)
			res.Value = &DynamicValue_JsonValue{JsonValue: j}
		}
	})
	return res
}

func NewFromDynamicValue(v *DynamicValue) *Value {
	switch vv := v.Value.(type) {
	case *DynamicValue_IntValue:
		return &Value{
			Kind:     IntKind,
			IntValue: vv.IntValue,
		}
	case *DynamicValue_LongValue:
		return &Value{
			Kind:      LongKind,
			LongValue: vv.LongValue,
		}
	case *DynamicValue_FloatValue:
		return &Value{
			Kind:       FloatKind,
			FloatValue: vv.FloatValue,
		}
	case *DynamicValue_DoubleValue:
		return &Value{
			Kind:        DoubleKind,
			DoubleValue: vv.DoubleValue,
		}
	case *DynamicValue_StringValue:
		return &Value{
			Kind:        StringKind,
			StringValue: vv.StringValue,
		}
	case *DynamicValue_BoolValue:
		return &Value{
			Kind:      BoolKind,
			BoolValue: vv.BoolValue,
		}
	case *DynamicValue_NullValue:
		return &Value{
			Kind: NullKind,
		}
	case *DynamicValue_JsonValue:
		if vv.JsonValue != nil {
			return &Value{
				Kind:      JsonKind,
				JsonValue: vv.JsonValue.AsMap(),
			}
		}
		return &Value{
			Kind:      JsonKind,
			JsonValue: map[string]interface{}{},
		}
	}
	return &Value{}
}
