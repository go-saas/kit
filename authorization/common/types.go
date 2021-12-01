package common

type HasIdentity interface {
	GetIdentity() string
}

type Action interface {
	HasIdentity
}

type Resource interface {
	HasIdentity
	GetNamespace() string
}

type Subject interface {
	HasIdentity
	GetName() string
}

type GrantType int32

const (
	GrantTypeUnknown GrantType = iota
	GrantTypeAllow
	GrantTypeDisallow
)
