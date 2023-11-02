package data

import (
	"errors"
	gorm2 "github.com/go-saas/kit/pkg/gorm"
	kitgorm "github.com/go-saas/kit/pkg/gorm"
	"github.com/go-saas/kit/pkg/query"
	v1 "github.com/go-saas/kit/user/api/user/v1"
	"github.com/go-saas/kit/user/private/biz"
	"github.com/go-saas/saas"
	"github.com/google/uuid"
	"github.com/goxiaoy/go-eventbus"
	concurrency "github.com/goxiaoy/gorm-concurrency/v2"
	"gorm.io/gorm"
)
import "context"

type UserRepo struct {
	*kitgorm.Repo[biz.User, string, *v1.ListUsersRequest]
}

func NewUserRepo(data *Data, eventbus *eventbus.EventBus) biz.UserRepo {
	res := &UserRepo{}
	res.Repo = kitgorm.NewRepo[biz.User, string, *v1.ListUsersRequest](data.DbProvider, eventbus, res)
	return res
}

var _ biz.UserRepo = (*UserRepo)(nil)

func (u *UserRepo) GetDb(ctx context.Context) *gorm.DB {
	return GetDb(ctx, u.DbProvider)
}

func buildUserScope(filter *v1.UserFilter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		ret := db
		if filter == nil {
			return ret
		}
		if len(filter.And) > 0 {
			for _, filter := range filter.And {
				ret = ret.Where(buildUserScope(filter)(db.Session(&gorm.Session{NewDB: true})))
			}
		}
		if len(filter.Or) > 0 {
			for _, filter := range filter.Or {
				ret = ret.Or(buildUserScope(filter)(db.Session(&gorm.Session{NewDB: true})))
			}
		}
		ret = ret.Scopes(gorm2.BuildStringFilter("`id`", filter.Id))
		ret = ret.Scopes(gorm2.BuildStringFilter("`gender`", filter.Gender))
		ret = ret.Scopes(gorm2.BuildDateFilter("`birthday`", filter.Birthday))
		return ret
	}
}

func buildUserTenantsScope() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		ti, _ := saas.FromCurrentTenant(db.Statement.Context)
		return db.Joins("INNER JOIN user_tenants on user_tenants.user_id = users.id").
			Where("user_tenants.tenant_id = ?", ti.GetId()).Group("users.id")
	}
}

func (u *UserRepo) BuildDetailScope(withDetail bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Preload("Roles").Preload("Tenants")
	}
}
func (u *UserRepo) BuildFilterScope(q *v1.ListUsersRequest) func(db *gorm.DB) *gorm.DB {
	return buildUserScope(q.Filter)
}

// DefaultSorting get default sorting
func (u *UserRepo) DefaultSorting() []string {
	return []string{"-created_at"}
}

func (u *UserRepo) List(ctx context.Context, query *v1.ListUsersRequest) ([]*biz.User, error) {
	var e biz.User
	db := u.GetDb(ctx).Model(&e)
	db = db.Scopes(u.BuildFilterScope(query), u.BuildDetailScope(false), buildUserTenantsScope(), kitgorm.SortScope(query, u.DefaultSorting()), kitgorm.PageScope(query))
	var items []*biz.User
	res := db.Find(&items)
	return items, res.Error
}

func (u *UserRepo) Count(ctx context.Context, query *v1.ListUsersRequest) (total int64, filtered int64, err error) {
	var e biz.User
	db := u.GetDb(ctx).Model(&e).Scopes(buildUserTenantsScope())
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	db = db.Scopes(u.BuildFilterScope(query))
	err = db.Count(&filtered).Error
	return
}

func (u *UserRepo) ListAdmin(ctx context.Context, query *v1.AdminListUsersRequest) ([]*biz.User, error) {
	var e biz.User
	db := u.GetDb(ctx).Model(&e)
	db = db.Scopes(buildUserScope(query.Filter), u.BuildDetailScope(false), kitgorm.SortScope(query, u.DefaultSorting()), kitgorm.PageScope(query))
	var items []*biz.User
	res := db.Find(&items)
	return items, res.Error
}

func (u *UserRepo) CountAdmin(ctx context.Context, query *v1.AdminListUsersRequest) (total int64, filtered int64, err error) {
	var e biz.User
	db := u.GetDb(ctx).Model(&e)
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	db = db.Scopes(buildUserScope(query.Filter))
	err = db.Count(&filtered).Error
	return
}

func (u *UserRepo) Create(ctx context.Context, user *biz.User) error {
	return u.GetDb(ctx).Create(user).Error
}

func (u *UserRepo) Update(ctx context.Context, id string, user *biz.User, p query.Select) error {
	return concurrency.ConcurrentUpdates(u.GetDb(ctx).Model(user).Omit("Roles", "Tenants"), *user).Error
}

func (u *UserRepo) FindByID(ctx context.Context, id string) (*biz.User, error) {
	return u.Get(ctx, id)
}

func (u *UserRepo) FindByName(ctx context.Context, name string) (*biz.User, error) {
	user := &biz.User{}
	err := u.GetDb(ctx).Model(&biz.User{}).Scopes(u.BuildDetailScope(true)).First(user, "normalized_username = ?", name).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (u *UserRepo) FindByPhone(ctx context.Context, phone string) (*biz.User, error) {
	user := &biz.User{}
	err := u.GetDb(ctx).Model(&biz.User{}).Scopes(u.BuildDetailScope(true)).First(user, "phone = ?", phone).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (u *UserRepo) FindByEmail(ctx context.Context, email string) (*biz.User, error) {
	user := &biz.User{}
	err := u.GetDb(ctx).Model(&biz.User{}).Scopes(u.BuildDetailScope(true)).First(user, "normalized_email = ?", email).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (u *UserRepo) AddLogin(ctx context.Context, user *biz.User, userLogin *biz.UserLogin) error {
	userLogin.UserId = user.ID
	err := u.GetDb(ctx).Create(userLogin).Error
	return err
}

func (u *UserRepo) RemoveLogin(ctx context.Context, user *biz.User, loginProvider string, providerKey string) error {
	err := u.GetDb(ctx).Scopes(gorm2.WhereUserId(user.ID.String())).Where("login_provider =?", loginProvider).Where("provider_key =?", providerKey).Delete(&biz.UserLogin{}).Error
	return err
}

func (u *UserRepo) ListLogin(ctx context.Context, user *biz.User) (userLogins []*biz.UserLogin, err error) {
	err = u.GetDb(ctx).Model(&biz.UserLogin{}).Scopes(gorm2.WhereUserId(user.ID.String())).Find(&userLogins).Error
	return
}

func (u *UserRepo) FindByLogin(ctx context.Context, loginProvider string, providerKey string) (*biz.User, error) {
	user := &biz.User{}
	err := u.GetDb(ctx).Model(&biz.User{}).
		Joins("left join user_logins on user_logins.user_id = users.id").
		Where("user_logins.login_provider=? and user_logins.provider_key=?", loginProvider, providerKey).First(user).Error
	return user, err
}

func (u *UserRepo) SetToken(ctx context.Context, user *biz.User, loginProvider string, name string, value string) (err error) {
	err = u.RemoveToken(ctx, user, loginProvider, name)
	if err != nil {
		return
	}
	err = u.GetDb(ctx).Model(&biz.UserToken{}).Create(&biz.UserToken{UserId: user.ID, LoginProvider: loginProvider, Name: name, Value: value}).Error
	return
}

func (u *UserRepo) RemoveToken(ctx context.Context, user *biz.User, loginProvider string, name string) (err error) {
	err = u.GetDb(ctx).Scopes(gorm2.WhereUserId(user.ID.String())).Where("login_provider =?", loginProvider).Where("name =?", name).Delete(&biz.UserToken{}).Error
	return
}

func (u *UserRepo) GetToken(ctx context.Context, user *biz.User, loginProvider string, name string) (token *string, err error) {
	var t biz.UserToken
	err = u.GetDb(ctx).Scopes(gorm2.WhereUserId(user.ID.String())).Where("login_provider =?", loginProvider).Where("name =?", name).First(&t).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &t.Value, nil
}

func (u *UserRepo) GetRoles(ctx context.Context, userId string) ([]biz.Role, error) {
	db := u.GetDb(ctx)
	var ret []biz.Role
	user := &biz.User{UIDBase: gorm2.UIDBase{ID: uuid.MustParse(userId)}}
	if err := db.Model(user).Association("Roles").Find(&ret); err != nil {
		return nil, err
	} else {
		return ret, nil
	}
}

func (u *UserRepo) UpdateRoles(ctx context.Context, user *biz.User, roles []biz.Role) error {
	db := u.GetDb(ctx)
	if err := db.Model(user).Association("Roles").Replace(roles); err != nil {
		return err
	}
	return nil
}

func (u *UserRepo) AddToRole(ctx context.Context, user *biz.User, role *biz.Role) error {
	db := u.GetDb(ctx)
	return db.Model(user).Association("Roles").Append(role)
}

func (u *UserRepo) RemoveFromRole(ctx context.Context, user *biz.User, role *biz.Role) error {
	db := u.GetDb(ctx)
	return db.Model(user).Association("Roles").Delete(role)
}
