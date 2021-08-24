package gorm

import (
	"github.com/a8m/rql"
	g "gorm.io/gorm"
)

type Repo struct {

}

func (r *Repo) BuildQuery(db *g.DB, model interface{}, query interface{}) (*g.DB, error) {
	if query == nil {
		return db, nil
	}
	queryParser := rql.MustNewParser(rql.Config{
		Model:    model,
		FieldSep: ".",
		OpPrefix: "",
	})
	q := rql.Query{}
	if page, ok := query.(rql.Page); ok {
		q.Limit = int(page.GetPageSize())
		q.Offset = int(page.GetPageOffset())
	}
	if sort, ok := query.(rql.Sort); ok {
		q.Sort = sort.GetSort()
	}
	if filter, ok := query.(rql.Filter); ok {
		q.Filter = filter.GetFilter()
	}
	if sel, ok := query.(rql.Select); ok {
		if f := sel.GetFields(); f != nil {
			q.Select = f.GetPaths()
		}
	}
	p, err := queryParser.ParseQuery(&q)
	if err != nil {
		return db, err
	}

	ret :=db.Model(model)
	if p.FilterExp!=""{
		if len(p.FilterArgs)>0{
			ret = ret.Where(p.FilterExp, p.FilterArgs)
		}else{
			ret = ret.Where(p.FilterExp)
		}
	}
	if p.Sort !=""{
		ret =ret.Order(p.Sort)
	}
	return ret.
		Offset(p.Offset).
		Limit(p.Limit), nil
}

func (r *Repo) BuildFilter(db *g.DB, model interface{}, query interface{}) (*g.DB, error) {
	if query == nil {
		return db, nil
	}
	queryParser := rql.MustNewParser(rql.Config{
		Model:    model,
		FieldSep: ".",
		OpPrefix: "",
	})
	q := rql.Query{}
	if filter, ok := query.(rql.Filter); ok {
		q.Filter = filter.GetFilter()
	}
	p, err := queryParser.ParseQuery(&q)
	if err != nil {
		return db, err
	}
	ret :=db.Model(model)
	if p.FilterExp!=""{
		if len(p.FilterArgs)>0{
			ret = ret.Where(p.FilterExp, p.FilterArgs)
		}else{
			ret = ret.Where(p.FilterExp)
		}
	}
	return ret, nil
}
