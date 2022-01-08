package authn

type UserInfo interface {
	GetId() string
}

type DefaultUserInfo struct {
	id string
}

var _ UserInfo = (*DefaultUserInfo)(nil)

func NewUserInfo(id string) *DefaultUserInfo {
	return &DefaultUserInfo{id: id}
}

func (d *DefaultUserInfo) GetId() string {
	return d.id
}
