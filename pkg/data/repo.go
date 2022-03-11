package data

import (
	"github.com/samber/lo"
	"strings"
)

//type Repo[T interface{}] interface {
//	List(ctx context.Context, query interface{}) ([]*T, error)
//	First(ctx context.Context, query interface{})(*T, error)
//	Count(ctx context.Context, query interface{}) (total int64, filtered int64, err error)
//	Get(ctx context.Context, id string) (*T, error)
//	Create(ctx context.Context, group *T) error
//	Update(ctx context.Context, id string, group *T, p rql.Select) error
//	Delete(ctx context.Context, id string) error
//}

var (
	sortDirection = map[byte]string{
		'+': "asc",
		'-': "desc",
	}
)

func ParseSort(fields []string) string {
	opts := ParseSortIntoOpt(fields)
	sortParams := lo.Map[*SortOpt, string](opts, func(s *SortOpt, _ int) string {
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
