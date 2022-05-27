package query

import (
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func FromString(s []string) []*wrapperspb.StringValue {
	return lo.Map(s, func(t string, _ int) *wrapperspb.StringValue {
		return &wrapperspb.StringValue{Value: t}
	})
}

func FromFloat64(s []float64) []*wrapperspb.DoubleValue {
	return lo.Map(s, func(t float64, _ int) *wrapperspb.DoubleValue {
		return &wrapperspb.DoubleValue{Value: t}
	})
}

func FromFlot32(s []float32) []*wrapperspb.FloatValue {
	return lo.Map(s, func(t float32, _ int) *wrapperspb.FloatValue {
		return &wrapperspb.FloatValue{Value: t}
	})
}
func FromInt32(s []int32) []*wrapperspb.Int32Value {
	return lo.Map(s, func(t int32, _ int) *wrapperspb.Int32Value {
		return &wrapperspb.Int32Value{Value: t}
	})
}
func FromInt64(s []int64) []*wrapperspb.Int64Value {
	return lo.Map(s, func(t int64, _ int) *wrapperspb.Int64Value {
		return &wrapperspb.Int64Value{Value: t}
	})
}
