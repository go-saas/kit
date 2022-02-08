package session

import (
	"context"
	"github.com/gorilla/sessions"
	"net/http"
)

const (
	sessionNamePrefix              = "identity"
	sessionNameUserId              = sessionNamePrefix + ".uid"
	sessionNameSecurityStamp       = sessionNamePrefix + ".security_stamp"
	sessionNameTwoFactorUserId     = sessionNamePrefix + ".tfa_uid"
	sessionNameTwoFactorProvider   = sessionNamePrefix + ".tfa_provider"
	sessionNameTwoFactorRememberMe = sessionNamePrefix + ".tfa_rm"
	sessionNameRememberToken       = sessionNamePrefix + ".rm"
	sessionNameExternal            = sessionNamePrefix + ".external"
)

type ClientState interface {
	//GetUid return user id
	GetUid() string
	//GetSecurityStamp return SecurityStamp
	GetSecurityStamp() string
	GetTwoFactorClientRemembered() bool
	GetTFAInfo() TFAInfo
	GetRememberToken() string
}

type TFAInfo struct {
	UserId        string
	LoginProvider string
}

type ClientStateImpl struct {
	s *sessions.Session
}

func NewClientState(s *sessions.Session) ClientState {
	return &ClientStateImpl{s: s}
}

func (c *ClientStateImpl) GetUid() string {
	if v, ok := c.s.Values[sessionNameUserId].(string); ok {
		return v
	} else {
		return ""
	}
}

func (c *ClientStateImpl) GetSecurityStamp() string {
	if v, ok := c.s.Values[sessionNameSecurityStamp].(string); ok {
		return v
	} else {
		return ""
	}
}

func (c *ClientStateImpl) GetTwoFactorClientRemembered() bool {
	if v, ok := c.s.Values[sessionNameTwoFactorRememberMe].(bool); ok {
		return v
	} else {
		return false
	}
}

func (c *ClientStateImpl) GetTFAInfo() TFAInfo {
	if tfaId, ok := c.s.Values[sessionNameTwoFactorUserId].(string); ok {
		if tfaProvider, ok := c.s.Values[sessionNameTwoFactorProvider].(string); ok {
			return TFAInfo{
				UserId:        tfaId,
				LoginProvider: tfaProvider,
			}
		}
	}
	return TFAInfo{}
}

func (c *ClientStateImpl) GetRememberToken() string {
	if v, ok := c.s.Values[sessionNameRememberToken].(string); ok {
		return v
	} else {
		return ""
	}
}

type ClientStateWriter interface {
	SetUid(ctx context.Context, uid string) error
	//SignOutUid clear uid in client state
	SignOutUid(ctx context.Context) error

	SetSecurityStamp(ctx context.Context, s string) error
	SignOutSecurityStamp(ctx context.Context) error

	SetTwoFactorClientRemembered(ctx context.Context) error
	SignOutTwoFactorClientRemembered(ctx context.Context) error

	SetTFAInfo(ctx context.Context, t TFAInfo) error
	SignOutTFAInfo(ctx context.Context) error

	SetRememberToken(ctx context.Context, r string) error
	SignOutRememberToken(ctx context.Context) error

	Clear(ctx context.Context) error
	//Save TODO should we keep this?
	Save(ctx context.Context) error
}

type ClientStateWriterImpl struct {
	s       *sessions.Session
	w       http.ResponseWriter
	r       *http.Request
	changed bool
}

var _ ClientStateWriter = (*ClientStateWriterImpl)(nil)

func NewClientStateWriter(s *sessions.Session, w http.ResponseWriter, r *http.Request) ClientStateWriter {
	return &ClientStateWriterImpl{s: s, w: w, r: r}
}

func (c *ClientStateWriterImpl) SetUid(ctx context.Context, uid string) error {
	c.s.Values[sessionNameUserId] = uid
	c.changed = true
	return nil
}

func (c *ClientStateWriterImpl) SetSecurityStamp(ctx context.Context, s string) error {
	c.s.Values[sessionNameSecurityStamp] = s
	c.changed = true
	return nil
}

func (c *ClientStateWriterImpl) SetTwoFactorClientRemembered(ctx context.Context) error {
	c.s.Values[sessionNameTwoFactorRememberMe] = true
	c.changed = true
	return nil
}

func (c *ClientStateWriterImpl) SetTFAInfo(ctx context.Context, t TFAInfo) error {
	c.s.Values[sessionNameTwoFactorUserId] = t.UserId
	c.s.Values[sessionNameTwoFactorProvider] = t.LoginProvider
	c.changed = true
	return nil
}

func (c *ClientStateWriterImpl) SetRememberToken(ctx context.Context, r string) error {
	c.s.Values[sessionNameRememberToken] = r
	c.changed = true
	return nil
}

func (c *ClientStateWriterImpl) SignOutUid(ctx context.Context) error {
	delete(c.s.Values, sessionNameUserId)
	c.changed = true
	return nil
}

func (c *ClientStateWriterImpl) SignOutSecurityStamp(ctx context.Context) error {
	delete(c.s.Values, sessionNameSecurityStamp)
	c.changed = true
	return nil
}

func (c *ClientStateWriterImpl) SignOutTwoFactorClientRemembered(ctx context.Context) error {
	delete(c.s.Values, sessionNameTwoFactorRememberMe)
	c.changed = true
	return nil
}

func (c *ClientStateWriterImpl) SignOutTFAInfo(ctx context.Context) error {
	delete(c.s.Values, sessionNameTwoFactorUserId)
	delete(c.s.Values, sessionNameTwoFactorProvider)
	c.changed = true
	return nil
}

func (c *ClientStateWriterImpl) SignOutRememberToken(ctx context.Context) error {
	delete(c.s.Values, sessionNameRememberToken)
	c.changed = true
	return nil
}

func (c *ClientStateWriterImpl) Clear(ctx context.Context) error {
	keys := make([]interface{}, len(c.s.Values))
	i := 0
	for k := range c.s.Values {
		keys[i] = k
		i++
	}
	for k := range keys {
		delete(c.s.Values, k)
	}
	c.changed = true
	return nil
}

func (c *ClientStateWriterImpl) Save(ctx context.Context) error {
	if c.changed {
		return c.s.Save(c.r, c.w)
	}
	return nil
}
