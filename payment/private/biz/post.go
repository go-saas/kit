package biz

import (
	v1 "github.com/goxiaoy/go-saas-kit/payment/api/post/v1"
	"github.com/goxiaoy/go-saas-kit/pkg/data"
	"github.com/goxiaoy/go-saas-kit/pkg/gorm"
)

type Post struct {
	gorm.UIDBase
	gorm.AuditedModel
	Name string
}

type PostRepo interface {
	data.Repo[Post, string, v1.ListPostRequest]
}
