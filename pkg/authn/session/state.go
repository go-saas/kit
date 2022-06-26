package session

import (
	"context"
	"github.com/go-saas/sessions"
)

const (
	sessionNamePrefix              = "identity"
	sessionNameUserId              = sessionNamePrefix + ".uid"
	sessionNameSecurityStamp       = sessionNamePrefix + ".security_stamp"
	sessionNameTwoFactorUserId     = sessionNamePrefix + ".tfa_uid"
	sessionNameTwoFactorProvider   = sessionNamePrefix + ".tfa_provider"
	sessionNameTwoFactorRememberMe = sessionNamePrefix + ".tfa_rm"
	sessionNameRememberToken       = sessionNamePrefix + ".rm"
	sessionNameRememberTokenUserId = sessionNameRememberToken + ".uid"
	sessionNameExternal            = sessionNamePrefix + ".external"
)

type ClientState interface {
	//GetUid return user id
	GetUid() string
	//GetSecurityStamp return SecurityStamp
	GetSecurityStamp() string
	GetTwoFactorClientRemembered() bool
	GetTFAInfo() *TFAInfo
	GetRememberToken() *RememberTokenInfo
}

type TFAInfo struct {
	UserId        string
	LoginProvider string
}
type RememberTokenInfo struct {
	Token string
	Uid   string
}

type ClientStateImpl struct {
	//s user info session
	s *sessions.Session
	//rs remember session
	rs *sessions.Session
}

func NewClientState(s *sessions.Session, rs *sessions.Session) ClientState {
	return &ClientStateImpl{s: s, rs: rs}
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

func (c *ClientStateImpl) GetTFAInfo() *TFAInfo {
	if tfaId, ok := c.s.Values[sessionNameTwoFactorUserId].(string); ok {
		if tfaProvider, ok := c.s.Values[sessionNameTwoFactorProvider].(string); ok {
			return &TFAInfo{
				UserId:        tfaId,
				LoginProvider: tfaProvider,
			}
		}
	}
	return nil
}

func (c *ClientStateImpl) GetRememberToken() *RememberTokenInfo {
	if v, ok := c.rs.Values[sessionNameRememberToken].(string); ok {
		if u, ok := c.rs.Values[sessionNameRememberTokenUserId].(string); ok {
			return &RememberTokenInfo{
				Token: v,
				Uid:   u,
			}
		}
	}
	return nil
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

	SetRememberToken(ctx context.Context, token, uid string) error
	SignOutRememberToken(ctx context.Context) error

	Clear(ctx context.Context) error

	Save(ctx context.Context) error
}

type ClientStateWriterImpl struct {
	//s user info session
	s *sessions.Session
	//rs remember session
	rs *sessions.Session

	w         sessions.Header
	h         sessions.Header
	sChanged  bool
	rsChanged bool
}

var _ ClientStateWriter = (*ClientStateWriterImpl)(nil)

func NewClientStateWriter(s *sessions.Session, rs *sessions.Session, w sessions.Header, h sessions.Header) ClientStateWriter {
	return &ClientStateWriterImpl{s: s, rs: rs, w: w, h: h}
}

func (c *ClientStateWriterImpl) SetUid(ctx context.Context, uid string) error {
	c.s.Values[sessionNameUserId] = uid
	c.sChanged = true
	return nil
}

func (c *ClientStateWriterImpl) SetSecurityStamp(ctx context.Context, s string) error {
	c.s.Values[sessionNameSecurityStamp] = s
	c.sChanged = true
	return nil
}

func (c *ClientStateWriterImpl) SetTwoFactorClientRemembered(ctx context.Context) error {
	c.s.Values[sessionNameTwoFactorRememberMe] = true
	c.sChanged = true
	return nil
}

func (c *ClientStateWriterImpl) SetTFAInfo(ctx context.Context, t TFAInfo) error {
	c.s.Values[sessionNameTwoFactorUserId] = t.UserId
	c.s.Values[sessionNameTwoFactorProvider] = t.LoginProvider
	c.sChanged = true
	return nil
}

func (c *ClientStateWriterImpl) SetRememberToken(ctx context.Context, token, uid string) error {
	if len(token) == 0 {
		return c.SignOutRememberToken(ctx)
	}
	c.rs.Values[sessionNameRememberToken] = token
	c.rs.Values[sessionNameRememberTokenUserId] = uid
	c.rsChanged = true
	return nil
}

func (c *ClientStateWriterImpl) SignOutUid(ctx context.Context) error {
	delete(c.s.Values, sessionNameUserId)
	c.sChanged = true
	return nil
}

func (c *ClientStateWriterImpl) SignOutSecurityStamp(ctx context.Context) error {
	delete(c.s.Values, sessionNameSecurityStamp)
	c.sChanged = true
	return nil
}

func (c *ClientStateWriterImpl) SignOutTwoFactorClientRemembered(ctx context.Context) error {
	delete(c.s.Values, sessionNameTwoFactorRememberMe)
	c.sChanged = true
	return nil
}

func (c *ClientStateWriterImpl) SignOutTFAInfo(ctx context.Context) error {
	delete(c.s.Values, sessionNameTwoFactorUserId)
	delete(c.s.Values, sessionNameTwoFactorProvider)
	c.sChanged = true
	return nil
}

func (c *ClientStateWriterImpl) SignOutRememberToken(ctx context.Context) error {
	c.rs.Options.MaxAge = -1
	c.rsChanged = true
	return nil
}

func (c *ClientStateWriterImpl) Clear(ctx context.Context) error {
	c.s.Options.MaxAge = -1
	c.rs.Options.MaxAge = -1
	c.rsChanged = true
	c.sChanged = true
	return nil
}

func (c *ClientStateWriterImpl) Save(ctx context.Context) error {
	if c.sChanged {
		if err := c.s.Save(c.h, c.w); err != nil {
			return err
		}
	}
	if c.rsChanged {
		if err := c.rs.Save(c.h, c.w); err != nil {
			return err
		}
	}
	return nil
}
