package authorization

import "fmt"

type UserSubject struct {
	userId string
}

var _ Subject = (*UserSubject)(nil)

func NewUserSubject(userId string) *UserSubject {
	return &UserSubject{userId: userId}
}

func (u *UserSubject) GetName() string {
	return "user"
}
func (u *UserSubject) GetIdentity() string {
	return fmt.Sprintf("%s/%s", u.GetName(), u.userId)
}

func (u *UserSubject) GetUserId() string {
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
	return fmt.Sprintf("%s/%s", t.GetName(), t.token)
}

func (t *TokenSubject) GetName() string {
	return "token"
}

func (t *TokenSubject) GetToken() string {
	return t.token
}

type RoleSubject struct {
	id string
}

var _ Subject = (*RoleSubject)(nil)

func NewRoleSubject(id string) *RoleSubject {
	return &RoleSubject{id: id}
}

func (r *RoleSubject) GetIdentity() string {
	return fmt.Sprintf("%s/%s", r.GetName(), r.id)
}

func (r *RoleSubject) GetName() string {
	return "role"
}

func (r *RoleSubject) GetRoleId() string {
	return r.id
}

type ClientSubject struct {
	clientId string
}

var _ Subject = (*ClientSubject)(nil)

func NewClientSubject(clientId string) *ClientSubject {
	return &ClientSubject{clientId: clientId}
}
func (c *ClientSubject) GetIdentity() string {
	return fmt.Sprintf("%s/%s", "client", c.clientId)
}
func (c *ClientSubject) GetClientId() string {
	return c.clientId
}
