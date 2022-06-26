package gorm

import (
	"fmt"
	"github.com/go-saas/kit/pkg/query"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"gorm.io/gorm"
)

func BuildStringFilter(field string, filter *query.StringFilterOperation) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		res := db
		if filter == nil {
			return res
		}
		if filter.Eq != nil {
			res = res.Where(fmt.Sprintf("%s = ?", field), filter.Eq.Value)
		}
		if filter.Neq != nil {
			res = res.Where(fmt.Sprintf("%s <> ?", field), filter.Neq.Value)
		}
		if filter.Contains != nil {
			res = res.Where(fmt.Sprintf("%s LIKE ?", field), fmt.Sprintf("%%%v%%", filter.Contains.Value))
		}
		if filter.StartsWith != nil {
			res = res.Where(fmt.Sprintf("%s LIKE ?", field), fmt.Sprintf("%%%v", filter.StartsWith.Value))
		}
		if filter.NstartsWith != nil {
			res = res.Where(fmt.Sprintf("%s NOT LIKE ?", field), fmt.Sprintf("%%%v", filter.StartsWith.Value))
		}
		if filter.EndsWith != nil {
			res = res.Where(fmt.Sprintf("%s LIKE ?", field), fmt.Sprintf("%v%%", filter.EndsWith.Value))
		}
		if filter.EndsWith != nil {
			res = res.Where(fmt.Sprintf("%s NOT LIKE ?", field), fmt.Sprintf("%v%%", filter.EndsWith.Value))
		}
		if filter.In != nil {
			res = res.Where(fmt.Sprintf("%s IN (?)", field), lo.Uniq(lo.Map(filter.In, func(t *wrapperspb.StringValue, _ int) string {
				return t.Value
			})))
		}
		if filter.Nin != nil {
			res = res.Where(fmt.Sprintf("%s NOT IN (?)", field), lo.Uniq(lo.Map(filter.Nin, func(t *wrapperspb.StringValue, _ int) string {
				return t.Value
			})))
		}
		if filter.Null != nil {
			res = res.Where(fmt.Sprintf("%s IS NULL ", field))
		}
		if filter.Nnull != nil {
			res = res.Where(fmt.Sprintf("%s IS NOT NULL ", field))
		}
		if filter.Empty != nil {
			res = res.Where(fmt.Sprintf("%s IS NULL OR %s = '' ", field, field))
		}
		if filter.Nempty != nil {
			res = res.Where(fmt.Sprintf("%s IS NOT NULL AND %s <> '' ", field, field))
		}
		if filter.Like != nil {
			res = res.Where(fmt.Sprintf("%s LIKE ?", field), fmt.Sprintf("%%%v%%", filter.Like.Value))
		}
		return res
	}
}

func BuildBooleanFilter(field string, filter *query.BooleanFilterOperators) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		res := db
		if filter == nil {
			return res
		}
		if filter.Eq != nil {
			res = res.Where(fmt.Sprintf("%s = ?", field), filter.Eq.Value)
		}
		if filter.Neq != nil {
			res = res.Where(fmt.Sprintf("%s <> ?", field), filter.Neq.Value)
		}
		if filter.Null != nil {
			res = res.Where(fmt.Sprintf("%s IS NULL ", field))
		}
		if filter.Nnull != nil {
			res = res.Where(fmt.Sprintf("%s IS NOT NULL ", field))
		}
		return res
	}
}

func BuildNullFilter(field string, filter *query.NullFilterOperators) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		res := db
		if filter == nil {
			return res
		}
		if filter.Null != nil {
			res = res.Where(fmt.Sprintf("%s IS NULL ", field))
		}
		if filter.Nnull != nil {
			res = res.Where(fmt.Sprintf("%s IS NOT NULL ", field))
		}
		return res
	}
}
func BuildDateFilter(field string, filter *query.DateFilterOperators) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		res := db
		if filter == nil {
			return res
		}
		if filter.Eq != nil {
			res = res.Where(fmt.Sprintf("%s = ?", field), filter.Eq.AsTime())
		}
		if filter.Neq != nil {
			res = res.Where(fmt.Sprintf("%s <> ?", field), filter.Neq.AsTime())
		}
		if filter.Gt != nil {
			res = res.Where(fmt.Sprintf("%s > ?", field), filter.Gt.AsTime())
		}
		if filter.Gte != nil {
			res = res.Where(fmt.Sprintf("%s >= ?", field), filter.Gte.AsTime())
		}
		if filter.Lt != nil {
			res = res.Where(fmt.Sprintf("%s < ?", field), filter.Lt.AsTime())
		}
		if filter.Lte != nil {
			res = res.Where(fmt.Sprintf("%s <= ?", field), filter.Lte.AsTime())
		}

		if filter.Null != nil {
			res = res.Where(fmt.Sprintf("%s IS NULL ", field))
		}
		if filter.Nnull != nil {
			res = res.Where(fmt.Sprintf("%s IS NOT NULL ", field))
		}

		return res
	}
}

func BuildDoubleFilter(field string, filter *query.DoubleFilterOperators) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		res := db
		if filter == nil {
			return res
		}
		if filter.Eq != nil {
			res = res.Where(fmt.Sprintf("%s = ?", field), filter.Eq.Value)
		}
		if filter.Neq != nil {
			res = res.Where(fmt.Sprintf("%s <> ?", field), filter.Neq.Value)
		}
		if filter.In != nil {
			res = res.Where(fmt.Sprintf("%s IN (?)", field), lo.Uniq(lo.Map(filter.In, func(t *wrapperspb.DoubleValue, _ int) float64 {
				return t.Value
			})))
		}
		if filter.Nin != nil {
			res = res.Where(fmt.Sprintf("%s NOT IN (?)", field), lo.Uniq(lo.Map(filter.Nin, func(t *wrapperspb.DoubleValue, _ int) float64 {
				return t.Value
			})))
		}
		if filter.Gt != nil {
			res = res.Where(fmt.Sprintf("%s > ?", field), filter.Gt.Value)
		}
		if filter.Gte != nil {
			res = res.Where(fmt.Sprintf("%s >= ?", field), filter.Gte.Value)
		}
		if filter.Lt != nil {
			res = res.Where(fmt.Sprintf("%s < ?", field), filter.Lt.Value)
		}
		if filter.Lte != nil {
			res = res.Where(fmt.Sprintf("%s <= ?", field), filter.Lte.Value)
		}

		if filter.Null != nil {
			res = res.Where(fmt.Sprintf("%s IS NULL ", field))
		}
		if filter.Nnull != nil {
			res = res.Where(fmt.Sprintf("%s IS NOT NULL ", field))
		}

		return res
	}
}

func BuildFloatFilter(field string, filter *query.FloatFilterOperators) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		res := db
		if filter == nil {
			return res
		}
		if filter.Eq != nil {
			res = res.Where(fmt.Sprintf("%s = ?", field), filter.Eq.Value)
		}
		if filter.Neq != nil {
			res = res.Where(fmt.Sprintf("%s <> ?", field), filter.Neq.Value)
		}
		if filter.In != nil {
			res = res.Where(fmt.Sprintf("%s IN (?)", field), lo.Uniq(lo.Map(filter.In, func(t *wrapperspb.FloatValue, _ int) float32 {
				return t.Value
			})))
		}
		if filter.Nin != nil {
			res = res.Where(fmt.Sprintf("%s NOT IN (?)", field), lo.Uniq(lo.Map(filter.Nin, func(t *wrapperspb.FloatValue, _ int) float32 {
				return t.Value
			})))
		}
		if filter.Gt != nil {
			res = res.Where(fmt.Sprintf("%s > ?", field), filter.Gt.Value)
		}
		if filter.Gte != nil {
			res = res.Where(fmt.Sprintf("%s >= ?", field), filter.Gte.Value)
		}
		if filter.Lt != nil {
			res = res.Where(fmt.Sprintf("%s < ?", field), filter.Lt.Value)
		}
		if filter.Lte != nil {
			res = res.Where(fmt.Sprintf("%s <= ?", field), filter.Lte.Value)
		}

		if filter.Null != nil {
			res = res.Where(fmt.Sprintf("%s IS NULL ", field))
		}
		if filter.Nnull != nil {
			res = res.Where(fmt.Sprintf("%s IS NOT NULL ", field))
		}
		return res
	}
}

func BuildInt32Filter(field string, filter *query.Int32FilterOperators) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		res := db
		if filter == nil {
			return res
		}
		if filter.Eq != nil {
			res = res.Where(fmt.Sprintf("%s = ?", field), filter.Eq.Value)
		}
		if filter.Neq != nil {
			res = res.Where(fmt.Sprintf("%s <> ?", field), filter.Neq.Value)
		}
		if filter.In != nil {
			res = res.Where(fmt.Sprintf("%s IN (?)", field), lo.Uniq(lo.Map(filter.In, func(t *wrapperspb.Int32Value, _ int) int32 {
				return t.Value
			})))
		}
		if filter.Nin != nil {
			res = res.Where(fmt.Sprintf("%s NOT IN (?)", field), lo.Uniq(lo.Map(filter.Nin, func(t *wrapperspb.Int32Value, _ int) int32 {
				return t.Value
			})))
		}
		if filter.Gt != nil {
			res = res.Where(fmt.Sprintf("%s > ?", field), filter.Gt.Value)
		}
		if filter.Gte != nil {
			res = res.Where(fmt.Sprintf("%s >= ?", field), filter.Gte.Value)
		}
		if filter.Lt != nil {
			res = res.Where(fmt.Sprintf("%s < ?", field), filter.Lt.Value)
		}
		if filter.Lte != nil {
			res = res.Where(fmt.Sprintf("%s <= ?", field), filter.Lte.Value)
		}

		if filter.Null != nil {
			res = res.Where(fmt.Sprintf("%s IS NULL ", field))
		}
		if filter.Nnull != nil {
			res = res.Where(fmt.Sprintf("%s IS NOT NULL ", field))
		}
		return res
	}
}

func BuildInt64Filter(field string, filter *query.Int64FilterOperators) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		res := db
		if filter == nil {
			return res
		}
		if filter.Eq != nil {
			res = res.Where(fmt.Sprintf("%s = ?", field), filter.Eq.Value)
		}
		if filter.Neq != nil {
			res = res.Where(fmt.Sprintf("%s <> ?", field), filter.Neq.Value)
		}
		if filter.In != nil {
			res = res.Where(fmt.Sprintf("%s IN (?)", field), lo.Uniq(lo.Map(filter.In, func(t *wrapperspb.Int64Value, _ int) int64 {
				return t.Value
			})))
		}
		if filter.Nin != nil {
			res = res.Where(fmt.Sprintf("%s NOT IN (?)", field), lo.Uniq(lo.Map(filter.Nin, func(t *wrapperspb.Int64Value, _ int) int64 {
				return t.Value
			})))
		}
		if filter.Gt != nil {
			res = res.Where(fmt.Sprintf("%s > ?", field), filter.Gt.Value)
		}
		if filter.Gte != nil {
			res = res.Where(fmt.Sprintf("%s >= ?", field), filter.Gte.Value)
		}
		if filter.Lt != nil {
			res = res.Where(fmt.Sprintf("%s < ?", field), filter.Lt.Value)
		}
		if filter.Lte != nil {
			res = res.Where(fmt.Sprintf("%s <= ?", field), filter.Lte.Value)
		}

		if filter.Null != nil {
			res = res.Where(fmt.Sprintf("%s IS NULL ", field))
		}
		if filter.Nnull != nil {
			res = res.Where(fmt.Sprintf("%s IS NOT NULL ", field))
		}
		return res
	}
}
