package gorm

import (
	"fmt"
	kitgorm "github.com/go-saas/kit/pkg/gorm"
	"github.com/go-saas/kit/pkg/localize"
	"github.com/go-saas/kit/pkg/multilingual"
	"github.com/samber/lo"
	"golang.org/x/text/language"
	"gorm.io/gorm"
)

func PreloadCurrentLanguage() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if _, ok := kitgorm.IsModel[multilingual.Multilingual](db); !ok {
			return db
		}
		if l := localize.LanguageTags(db.Statement.Context); l != nil {
			if len(l) > 0 {
				rootQuery := make(map[string]bool)
				for _, tag := range l {
					for {
						if tag.Parent() == language.Und {
							rootQuery[fmt.Sprintf("%s%%", tag.String())] = true
							break
						}
						tag = tag.Parent()
					}
				}
				queryStr := ""
				allLan := lo.Keys(rootQuery)
				for i, _ := range allLan {
					if i != 0 {
						queryStr = fmt.Sprintf("%s OR ", queryStr)
					}
					queryStr = fmt.Sprintf("%slanguage_code LIKE ?", queryStr)

				}
				queryArg := []interface{}{queryStr}
				for _, s := range allLan {
					queryArg = append(queryArg, s)
				}
				return db.Preload("Trans", queryArg...)
			}

		}
		return db
	}
}
