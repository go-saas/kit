package biz

import "time"

type TenantFeature struct {
	TenantId string
	//key of feature
	Key string
	//value
	Value     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
