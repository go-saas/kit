package query

import (
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

type Filter[TFilter any] interface {
	GetFilter() TFilter
}
