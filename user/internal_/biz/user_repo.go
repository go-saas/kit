package biz

import "context"

type UserRepo interface {
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id string) (*User, error)
	FindByName(ctx context.Context, name string) (*User, error)
	FindByPhone(ctx context.Context, phone string) (*User, error)
	AddLogin(ctx context.Context, user *User, userLogin *UserLogin) error
	RemoveLogin(ctx context.Context, user *User, loginProvider string, providerKey string) error
	ListLogin(ctx context.Context, user *User) ([]*UserLogin, error)
	FindByLogin(ctx context.Context, loginProvider string, providerKey string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	SetToken(ctx context.Context, user *User, loginProvider string, name string, value string) error
	RemoveToken(ctx context.Context, user *User, loginProvider string, name string) error
	GetToken(ctx context.Context, user *User, loginProvider string, name string) (*string, error)
	GetRoles(ctx context.Context, user *User, )([]*Role,error)
	UpdateRoles(ctx context.Context, user *User,roles []*Role)error
	AddToRole(ctx context.Context, user *User,role *Role )error
	RemoveFromRole(ctx context.Context, user *User,role *Role )error

}
