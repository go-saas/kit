package common

type UserSubject struct {
	userId string
}

func NewUserSubject(userId string) *UserSubject {
	return &UserSubject{userId: userId}
}

func (u *UserSubject) GetIdentity() string {
	return u.userId
}

var _ Subject = (*UserSubject)(nil)
