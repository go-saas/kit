package authorization

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

type Effect int32

const (
	EffectUnknown Effect = iota
	EffectGrant
	EffectForbidden
)
