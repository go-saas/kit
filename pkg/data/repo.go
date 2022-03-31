package data

import (
	"context"
	"github.com/goxiaoy/go-saas-kit/pkg/query"
	"github.com/samber/lo"
	"strings"
)

type Repo[TEntity any, TKey any, TQuery any] interface {
	List(ctx context.Context, query *TQuery) ([]*TEntity, error)
	First(ctx context.Context, query *TQuery) (*TEntity, error)
	Count(ctx context.Context, query *TQuery) (total int64, filtered int64, err error)
	Get(ctx context.Context, id TKey) (*TEntity, error)
	Create(ctx context.Context, entity *TEntity) error
	BatchCreate(ctx context.Context, entity []*TEntity, batchSize int) error
	Update(ctx context.Context, id TKey, entity *TEntity, p query.Select) error
	Delete(ctx context.Context, id TKey) error
}

type BeforeCreate[TEntity any] struct {
	Entity TEntity
}

func NewBeforeCreate[TEntity any](entity TEntity) *BeforeCreate[TEntity] {
	return &BeforeCreate[TEntity]{
		Entity: entity,
	}
}

type AfterCreate[TEntity any] struct {
	Entity TEntity
}

func NewAfterCreate[TEntity any](entity TEntity) *AfterCreate[TEntity] {
	return &AfterCreate[TEntity]{
		Entity: entity,
	}
}

type BeforeUpdate[TEntity any] struct {
	Entity TEntity
	P      query.Select
}

func NewBeforeUpdate[TEntity any](entity TEntity) *BeforeUpdate[TEntity] {
	return &BeforeUpdate[TEntity]{
		Entity: entity,
	}
}

type AfterUpdate[TEntity any] struct {
	Entity TEntity
}

func NewAfterUpdate[TEntity any](entity TEntity) *AfterUpdate[TEntity] {
	return &AfterUpdate[TEntity]{
		Entity: entity,
	}
}

type BeforeDelete[TEntity any] struct {
	Entity TEntity
}

func NewBeforeDelete[TEntity any](entity TEntity) *BeforeDelete[TEntity] {
	return &BeforeDelete[TEntity]{
		Entity: entity,
	}
}

type AfterDelete[TEntity any] struct {
	Entity TEntity
}

func NewAfterDelete[TEntity any](entity TEntity) *AfterDelete[TEntity] {
	return &AfterDelete[TEntity]{
		Entity: entity,
	}
}

var (
	sortDirection = map[byte]string{
		'+': "asc",
		'-': "desc",
	}
)

func ParseSort(fields []string) string {
	opts := ParseSortIntoOpt(fields)
	sortParams := lo.Map(opts, func(s *SortOpt, _ int) string {
		colName := s.Field
		if s.IsDesc {
			colName += " " + "desc"
		} else {
			colName += " " + "asc"
		}
		return colName
	})
	return strings.Join(sortParams, ", ")
}

func ParseSortIntoOpt(fields []string) []*SortOpt {
	sortParams := make([]*SortOpt, len(fields))
	for i, field := range fields {
		var orderBy string
		if order, ok := sortDirection[field[0]]; ok {
			orderBy = order
			field = field[1:]
		}
		opt := &SortOpt{
			Field:  field,
			IsDesc: orderBy == "desc",
		}
		sortParams[i] = opt
	}
	return sortParams
}

type SortOpt struct {
	Field  string
	IsDesc bool
}
