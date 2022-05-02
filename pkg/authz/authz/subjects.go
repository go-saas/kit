package authz

import (
	"fmt"
	"strings"
)

type UserSubject struct {
	userId string
}

var _ Subject = (*UserSubject)(nil)

func NewUserSubject(userId string) *UserSubject {
	return &UserSubject{userId: userId}
}

func ParseUserSubject(subject Subject) (*UserSubject, bool) {
	if strings.HasPrefix(subject.GetIdentity(), "user/") {
		return NewUserSubject(strings.TrimPrefix(subject.GetIdentity(), "user/")), true
	}
	return nil, false

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
func ParseTokenSubject(subject Subject) (*TokenSubject, bool) {
	if strings.HasPrefix(subject.GetIdentity(), "token/") {
		return NewTokenSubject(strings.TrimPrefix(subject.GetIdentity(), "token/")), true
	}
	return nil, false
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

func ParseRoleSubject(subject Subject) (*RoleSubject, bool) {
	if strings.HasPrefix(subject.GetIdentity(), "role/") {
		return NewRoleSubject(strings.TrimPrefix(subject.GetIdentity(), "role/")), true
	}
	return nil, false
}

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

func ParseClientSubject(subject Subject) (*ClientSubject, bool) {
	if strings.HasPrefix(subject.GetIdentity(), "client/") {
		return NewClientSubject(strings.TrimPrefix(subject.GetIdentity(), "client/")), true
	}
	return nil, false
}

func NewClientSubject(clientId string) *ClientSubject {
	return &ClientSubject{clientId: clientId}
}
func (c *ClientSubject) GetIdentity() string {
	return fmt.Sprintf("%s/%s", "client", c.clientId)
}
func (c *ClientSubject) GetClientId() string {
	return c.clientId
}

type SubjectStr string

var _ Subject = (*SubjectStr)(nil)

func (s SubjectStr) GetIdentity() string {
	return string(s)
}

type TenantSubject struct {
	id string
}

var _ Subject = (*TenantSubject)(nil)

func ParseTenantSubject(subject Subject) (*TenantSubject, bool) {
	if strings.HasPrefix(subject.GetIdentity(), "tenant/") {
		return NewTenantSubject(strings.TrimPrefix(subject.GetIdentity(), "tenant/")), true
	}
	return nil, false
}

func NewTenantSubject(id string) *TenantSubject {
	return &TenantSubject{id: id}
}

func (r *TenantSubject) GetIdentity() string {
	return fmt.Sprintf("%s/%s", r.GetName(), r.id)
}

func (r *TenantSubject) GetName() string {
	return "tenant"
}

func (r *TenantSubject) GetTenantId() string {
	return r.id
}
