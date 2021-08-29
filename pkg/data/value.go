package data

type Value struct {
	// which kind of data, can be "null","int","long","float","double","string","bool"
	Kind        string
	IntValue    int32
	LongValue   int64
	FloatValue  float32
	DoubleValue float64
	StringValue string
	BoolValue   bool
}

const NullKind = "null"
const IntKind = "int"
const LongKind = "long"
const FloatKind = "float"
const DoubleKind = "double"
const StringKind = "string"
const BoolKind = "bool"

func (v *Value) GetValue() interface{} {
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

func (v *Value) RunIfInt(f func()) {
	if v.Kind == IntKind {
		f()
	}
}

func (v *Value) SetAsLong(value int64) {
	v.Kind = LongKind
	v.LongValue = value
}

func (v *Value) RunIfLong(f func()) {
	if v.Kind == LongKind {
		f()
	}
}

func (v *Value) SetAsFloat(value float32) {
	v.Kind = FloatKind
	v.FloatValue = value
}
func (v *Value) RunIfFloat(f func()) {
	if v.Kind == FloatKind {
		f()
	}
}
func (v *Value) SetAsDouble(value float64) {
	v.Kind = DoubleKind
	v.DoubleValue = value
}
func (v *Value) RunIfDouble(f func()) {
	if v.Kind == DoubleKind {
		f()
	}
}
func (v *Value) SetAsString(value string) {
	v.Kind = StringKind
	v.StringValue = value
}
func (v *Value) RunIfString(f func()) {
	if v.Kind == StringKind {
		f()
	}
}
func (v *Value) SetAsBool(value bool) {
	v.Kind = BoolKind
	v.BoolValue = value
}
func (v *Value) RunIfBool(f func()) {
	if v.Kind == BoolKind {
		f()
	}
}
