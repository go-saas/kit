package gorm

import (
	"github.com/a8m/rql"
	"github.com/goxiaoy/go-saas-kit/pkg/data"
	"gorm.io/gorm"
)



type Repo struct {
}

func PageScope(page rql.Page) func(db *gorm.DB) *gorm.DB {
	return func (db *gorm.DB) *gorm.DB {
		if page==nil{
			return db
		}
		ret := db
		if page.GetPageOffset()>0{
			ret=db.Offset(int(page.GetPageOffset()))
		}
		if page.GetPageSize()>0{
			ret=db.Limit(int(page.GetPageSize()))
		}
		return ret
	}
}

// SortScope build sorting by sort and default d
func SortScope(sort rql.Sort,d []string) func(db *gorm.DB) *gorm.DB {
	return func (db *gorm.DB) *gorm.DB  {
		var s []string
		if sort!=nil{
			s = sort.GetSort()
		}
		if len(s)==0{
			s = d
		}
		parsed := data.ParseSort(s)
		ret := db
		if parsed!=""{
			ret = ret.Order(parsed)
		}
		return ret
	}
}




//func (r *Repo) BuildQuery(db *g.DB, model interface{}, query interface{}) (*g.DB, error) {
//	if query == nil {
//		return db, nil
//	}
//	queryParser := rql.MustNewParser(rql.Config{
//		Model:    model,
//		FieldSep: ".",
//		OpPrefix: "",
//	})
//	q := rql.Query{}
//	if page, ok := query.(rql.Page); ok {
//		q.Limit = int(page.GetPageSize())
//		q.Offset = int(page.GetPageOffset())
//	}
//	if sort, ok := query.(rql.Sort); ok {
//		q.Sort = sort.GetSort()
//	}
//	if filter, ok := query.(rql.Filter); ok {
//		q.Filter = filter.GetFilter()
//	}
//	if sel, ok := query.(rql.Select); ok {
//		if f := sel.GetFields(); f != nil {
//			q.Select = f.GetPaths()
//		}
//	}
//	p, err := queryParser.ParseQuery(&q)
//	if err != nil {
//		return db, err
//	}
//
//	ret := db.Model(model)
//	if p.FilterExp != "" {
//		if len(p.FilterArgs) > 0 {
//			ret = ret.Where(p.FilterExp, p.FilterArgs)
//		} else {
//			ret = ret.Where(p.FilterExp)
//		}
//	}
//	if p.Sort != "" {
//		ret = ret.Order(p.Sort)
//	}
//	return ret.
//		Offset(p.Offset).
//		Limit(p.Limit), nil
//}
//
//func (r *Repo) BuildFilter(db *g.DB, model interface{}, query interface{}) (*g.DB, error) {
//	if query == nil {
//		return db, nil
//	}
//	queryParser := rql.MustNewParser(rql.Config{
//		Model:    model,
//		FieldSep: ".",
//		OpPrefix: "",
//	})
//	q := rql.Query{}
//	if filter, ok := query.(rql.Filter); ok {
//		q.Filter = filter.GetFilter()
//	}
//	p, err := queryParser.ParseQuery(&q)
//	if err != nil {
//		return db, err
//	}
//	ret := db.Model(model)
//	if p.FilterExp != "" {
//		if len(p.FilterArgs) > 0 {
//			ret = ret.Where(p.FilterExp, p.FilterArgs)
//		} else {
//			ret = ret.Where(p.FilterExp)
//		}
//	}
//	return ret, nil
//}


