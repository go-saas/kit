package biz

import (
	v1 "github.com/go-saas/kit/saas/event/v1"
)

func NewUserMigrationTaskFromTenantEvent(t *v1.TenantCreatedEvent) *UserMigrationTask {
	return &UserMigrationTask{
		Id:            t.Id,
		AdminEmail:    t.AdminEmail,
		AdminUsername: t.AdminUsername,
		AdminPassword: t.AdminPassword,
		AdminUserId:   t.AdminUserId,
	}
}
