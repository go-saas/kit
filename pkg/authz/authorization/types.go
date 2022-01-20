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
}

type Effect int32

const (
	EffectUnknown Effect = iota
	EffectGrant
	EffectForbidden
)

type PermissionBean struct {
	Namespace string
	Resource  string
	Action    string
	Subject   string
	Effect    Effect
}

func NewPermissionBean(resource Resource, action Action, subject Subject, effect Effect) PermissionBean {
	return PermissionBean{
		Namespace: resource.GetNamespace(),
		Resource:  resource.GetIdentity(),
		Action:    action.GetIdentity(),
		Subject:   subject.GetIdentity(),
		Effect:    effect,
	}
}

type UpdateSubjectPermission struct {
	Resource Resource
	Action   Action
	Effect   Effect
}
