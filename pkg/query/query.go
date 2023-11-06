package query

import (
	"github.com/iancoleman/strcase"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"strings"
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

func SelectContains(p Select, name string) bool {
	if p == nil || p.GetFields() == nil {
		return true
	}
	name = strcase.ToSnake(name)
	_, r := lo.Find(p.GetFields().Paths, func(s string) bool {
		return s == "*" || s == name || strings.HasPrefix(s, name+".")
	})
	return r
}

func SelectStrictContains(p Select, name string) bool {
	if p == nil || p.GetFields() == nil {
		return false
	}
	name = strcase.ToSnake(name)
	_, r := lo.Find(p.GetFields().Paths, func(s string) bool {
		return s == name || strings.HasPrefix(s, name+".")
	})
	return r
}

func SelectGetCurrentLevelPath(p Select) []string {
	if p == nil || p.GetFields() == nil {
		return []string{"*"}
	}
	return GetCurrentLevelPath(p.GetFields().Paths)
}

func GetCurrentLevelPath(paths []string) []string {
	ret := map[string]bool{}
	for _, path := range paths {
		p, _, _ := strings.Cut(path, ".")
		ret[p] = true
	}
	return lo.Keys(ret)
}

func GetNextLevelPath(paths []string) []string {
	ret := map[string]bool{}
	for _, path := range paths {
		_, p, found := strings.Cut(path, ".")
		if found {
			ret[p] = true
		}
	}
	return lo.Keys(ret)
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
