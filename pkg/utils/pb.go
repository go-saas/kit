package utils

import (
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
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
