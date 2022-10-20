package utils

import (
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"time"
)

func Time2Timepb(time *time.Time) *timestamppb.Timestamp {
	if time == nil {
		return nil
	}
	return timestamppb.New(*time)
}

func Timepb2Time(t *timestamppb.Timestamp) *time.Time {
	if t == nil {
		return nil
	}
	ret := t.AsTime()
	return &ret
}

func Map2Structpb(m map[string]interface{}) *structpb.Struct {
	if m == nil {
		return nil
	}
	r, _ := structpb.NewStruct(m)
	return r
}

func Structpb2Map(m *structpb.Struct) map[string]interface{} {
	if m == nil {
		return nil
	}
	return m.AsMap()
}

func ToWrapString(s []string) []*wrapperspb.StringValue {
	return lo.Map(s, func(t string, _ int) *wrapperspb.StringValue {
		return &wrapperspb.StringValue{Value: t}
	})
}

func ToWrapFloat64(s []float64) []*wrapperspb.DoubleValue {
	return lo.Map(s, func(t float64, _ int) *wrapperspb.DoubleValue {
		return &wrapperspb.DoubleValue{Value: t}
	})
}

func ToWrapFlot32(s []float32) []*wrapperspb.FloatValue {
	return lo.Map(s, func(t float32, _ int) *wrapperspb.FloatValue {
		return &wrapperspb.FloatValue{Value: t}
	})
}
func ToWrapInt32(s []int32) []*wrapperspb.Int32Value {
	return lo.Map(s, func(t int32, _ int) *wrapperspb.Int32Value {
		return &wrapperspb.Int32Value{Value: t}
	})
}
func ToWrapInt64(s []int64) []*wrapperspb.Int64Value {
	return lo.Map(s, func(t int64, _ int) *wrapperspb.Int64Value {
		return &wrapperspb.Int64Value{Value: t}
	})
}
