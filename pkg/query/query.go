package query

import (
	"encoding/json"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// Page handle pagination
type Page interface {
	GetPageOffset() int32
	GetPageSize() int32
}

// Search field
type Search interface {
	GetSearch() string
}

// Sort interface handle sorting like '+created_at','-created_at'
type Sort interface {
	GetSort() []string
}

// Select fields to query or update
type Select interface {
	GetFields() *fieldmaskpb.FieldMask
}

// Filter by json
type Filter interface {
	GetFilter() map[string]interface{}
}

type QueryWrapper struct {
	// source object
	source interface{}
	// filter object
	filter proto.Message
	// filter normalized
	v map[string]interface{}
}

var _ Page = (*QueryWrapper)(nil)
var _ Search = (*QueryWrapper)(nil)
var _ Sort = (*QueryWrapper)(nil)
var _ Select = (*QueryWrapper)(nil)
var _ Filter = (*QueryWrapper)(nil)

func NewQueryFromProto(source proto.Message, filter proto.Message) *QueryWrapper {
	return &QueryWrapper{
		source: source,
		filter: filter,
	}
}

func (f *QueryWrapper) GetPageOffset() int32 {
	if s, ok := f.source.(Page); ok {
		return s.GetPageOffset()
	}
	return 0
}

func (f *QueryWrapper) GetPageSize() int32 {
	if s, ok := f.source.(Page); ok {
		return s.GetPageSize()
	}
	return 0
}

func (f *QueryWrapper) GetFields() *fieldmaskpb.FieldMask {
	if s, ok := f.source.(Select); ok {
		return s.GetFields()
	}
	return nil
}

func (f *QueryWrapper) GetSort() []string {
	if s, ok := f.source.(Sort); ok {
		return s.GetSort()
	}
	return make([]string, 0)
}

func (f *QueryWrapper) GetSearch() string {
	if s, ok := f.source.(Search); ok {
		return s.GetSearch()
	}
	return ""
}

func (f *QueryWrapper) ParseAndValidate() error {
	if f.v == nil {
		f.v = make(map[string]interface{})
	}
	b, err := protojson.Marshal(f.filter)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(b, &f.v); err != nil {
		return err
	}
	return nil
}

func (f *QueryWrapper) GetFilter() map[string]interface{} {
	if f.v == nil {
		if err := f.ParseAndValidate(); err != nil {
			panic(err)
		}
	}
	return f.v
}
