package data

import (
	"errors"
	gorm2 "github.com/goxiaoy/go-saas-kit/pkg/gorm"
	"github.com/goxiaoy/go-saas-kit/user/internal_/biz"
	concurrency "github.com/goxiaoy/gorm-concurrency"
	"gorm.io/gorm"
)
import "context"

type UserRepo struct {
	Repo
}

func NewUserRepo(data *Data) biz.UserRepo {
	return &UserRepo{
		Repo{
			DbProvider: data.DbProvider,
		},
	}
}

var _ biz.UserRepo = (*UserRepo)(nil)

func (u *UserRepo) GetDb(ctx context.Context) *gorm.DB {
	return GetDb(ctx, u.DbProvider)
}


func (u *UserRepo) Create(ctx context.Context, user *biz.User) error {
	return u.GetDb(ctx).Create(user).Error
}

func (u *UserRepo) Update(ctx context.Context, user *biz.User) error {
	return concurrency.ConcurrentUpdates(GetDb(ctx, u.DbProvider), user).Error
}

func (u *UserRepo) Delete(ctx context.Context, user *biz.User) error {
	return u.GetDb(ctx).Delete(user).Error
}

func (u *UserRepo) FindByID(ctx context.Context, id string) (*biz.User, error) {
	user := &biz.User{}
	err := u.GetDb(ctx).Model(&biz.User{}).Preload("Roles").First(user,"id=?", id).Error
	if err!=nil{
		if  errors.Is(err,gorm.ErrRecordNotFound){
			return nil,nil
		}
		return nil,err
	}
	return user, nil
}

func (u *UserRepo) FindByName(ctx context.Context, name string) (*biz.User, error) {
	user := &biz.User{}
	err := u.GetDb(ctx).Model(&biz.User{}).Preload("Roles").First(user, "normalized_username = ?", name).Error
	if err!=nil{
		if  errors.Is(err,gorm.ErrRecordNotFound){
			return nil,nil
		}
		return nil,err
	}
	return user, nil
}

func (u *UserRepo) FindByPhone(ctx context.Context, phone string) (*biz.User, error) {
	user := &biz.User{}
	err := u.GetDb(ctx).Model(&biz.User{}).Preload("Roles").First(user, "phone = ?", phone).Error
	if err!=nil{
		if  errors.Is(err,gorm.ErrRecordNotFound){
			return nil,nil
		}
		return nil,err
	}
	return user, nil
}

func (u *UserRepo) AddLogin(ctx context.Context, user *biz.User, userLogin *biz.UserLogin) error {
	userLogin.UserId = user.ID
	err := u.GetDb(ctx).Create(userLogin).Error
	return err
}

func (u *UserRepo) RemoveLogin(ctx context.Context, user *biz.User, loginProvider string, providerKey string) error {
	err := u.GetDb(ctx).Scopes(func(db *gorm.DB) *gorm.DB {
		return gorm2.WhereUserId(db, user.ID)
	}).Where("login_provider =?", loginProvider).Where("provider_key =?", providerKey).Delete(&biz.UserLogin{}).Error
	return err
}

func (u *UserRepo) ListLogin(ctx context.Context, user *biz.User) (userLogins []*biz.UserLogin, err error) {
	err = u.GetDb(ctx).Scopes(func(db *gorm.DB) *gorm.DB {
		return gorm2.WhereUserId(db, user.ID)
	}).Model(&biz.UserLogin{}).Find(userLogins).Error
	return
}

func (u *UserRepo) FindByLogin(ctx context.Context, loginProvider string, providerKey string) (*biz.User, error) {
	user := &biz.User{}
	err := u.GetDb(ctx).Model(&biz.User{}).Joins("left join user_logins on user_logins.user_id = users.id").Where("user_logins.login_provider=? and user_logins.provider_key=?", loginProvider, providerKey).First(user).Error
	return user,err
}

func (u *UserRepo) FindByEmail(ctx context.Context, email string) (*biz.User,error) {
	user := &biz.User{}
	err := u.GetDb(ctx).Model(&biz.User{}).Preload("Roles").First(user, "normalized_email = ?", email).Error
	if err!=nil{
		if  errors.Is(err,gorm.ErrRecordNotFound){
			return nil,nil
		}
		return nil,err
	}
	return user,nil
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
	err = u.GetDb(ctx).Scopes(func(db *gorm.DB) *gorm.DB {
		return gorm2.WhereUserId(db, user.ID)
	}).Where("login_provider =?", loginProvider).Where("name =?", name).Delete(&biz.UserToken{}).Error
	return
}

func (u *UserRepo) GetToken(ctx context.Context, user *biz.User, loginProvider string, name string) (token *string, err error) {
	var t biz.UserToken
	err = u.GetDb(ctx).Scopes(func(db *gorm.DB) *gorm.DB {
		return gorm2.WhereUserId(db, user.ID)
	}).Where("login_provider =?", loginProvider).Where("name =?", name).First(&t).Error
	if err!=nil{
		if  errors.Is(err,gorm.ErrRecordNotFound){
			return nil,nil
		}
		return nil,err
	}
	return &t.Value,nil
}

func (u *UserRepo) GetRoles(ctx context.Context, user *biz.User) ([]*biz.Role, error) {
	db :=u.GetDb(ctx).Preload("Roles")
	dbUser := &biz.User{}
	if err := db.Model(&biz.User{}).Preload("Roles").Where("id=?",user.ID).Find(dbUser).Error;err!=nil{
		return nil,err
	}
	var ret []*biz.Role
	for _,i :=range dbUser.Roles{
		ret = append(ret, &i)
	}
	return ret,nil
}

func (u *UserRepo) UpdateRoles(ctx context.Context, user *biz.User, roles []*biz.Role) error {
	//delete all previous
	db :=u.GetDb(ctx)
	if err :=db.Where("user_id=?",user.ID).Delete(biz.UserRole{}).Error;err!=nil{
		return err
	}
	var ur []*biz.UserRole
	for _, role := range roles {
		ur = append(ur, &biz.UserRole{
			UserID:       user.ID,
			RoleID: role.ID,
		})
	}
	if err:=db.CreateInBatches(ur,100).Error;err!=nil{
		return err
	}
	return nil

}

func (u *UserRepo) AddToRole(ctx context.Context, user *biz.User, role *biz.Role) error {
	db :=u.GetDb(ctx)
	ur := biz.UserRole{UserID:user.ID,RoleID:role.ID}
	err :=db.Model(&biz.UserRole{}).Where("user_id=? and role_id=?",user.ID,role.ID).FirstOrCreate(&ur).Error
	return err
}

func (u *UserRepo) RemoveFromRole(ctx context.Context, user *biz.User, role *biz.Role) error {
	db :=u.GetDb(ctx)
	err :=db.Where("user_id=? and role_id=?",user.ID,role.ID).Delete(&biz.UserRole{}).Error
	return err
}
