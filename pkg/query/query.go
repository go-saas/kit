package query

import (
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

type HasPageSize interface {
	GetPageSize() int32
}

// Page handle pagination
type Page interface {
	HasPageSize
	GetPageOffset() int32
}

// Search full text search field
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

type Filter[TFilter any] interface {
	GetFilter() TFilter
}

type Field struct {
	*fieldmaskpb.FieldMask
}

func NewField(f *fieldmaskpb.FieldMask) *Field {
	return &Field{
		f,
	}
}
func (f *Field) GetFields() *fieldmaskpb.FieldMask {
	return f.FieldMask
}

type CursorAfterPage interface {
	GetAfterPageToken() string
}

type CursorBeforePage interface {
	GetBeforePageToken() string
}
