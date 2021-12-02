package authorization

type UserSubject struct {
	userId string
}

func (u *UserSubject) GetName() string {
	return "user"
}

var _ Subject = (*UserSubject)(nil)

func NewUserSubject(userId string) *UserSubject {
	return &UserSubject{userId: userId}
}

func (u *UserSubject) GetIdentity() string {
	return u.userId
}

type TokenSubject struct {
	token string
}

var _ Subject = (*TokenSubject)(nil)

func NewTokenSubject(token string) *TokenSubject {
	return &TokenSubject{token: token}
}

func (t *TokenSubject) GetIdentity() string {
	return t.token
}

func (t *TokenSubject) GetName() string {
	return "token"
}
