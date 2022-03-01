package authz

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
	TenantID  string
	Effect    Effect
}

func NewPermissionBean(resource Resource, action Action, subject Subject, tenantID string, effect Effect) PermissionBean {
	return PermissionBean{
		Namespace: resource.GetNamespace(),
		Resource:  resource.GetIdentity(),
		Action:    action.GetIdentity(),
		Subject:   subject.GetIdentity(),
		TenantID:  tenantID,
		Effect:    effect,
	}
}

type UpdateSubjectPermission struct {
	Resource Resource
	Action   Action
	Effect   Effect
	TenantID string
}

func NewUpdateSubjectPermission(resource Resource, action Action, tenantID string, effect Effect) *UpdateSubjectPermission {
	return &UpdateSubjectPermission{
		Resource: resource,
		Action:   action,
		TenantID: tenantID,
		Effect:   effect,
	}
}

type PermissionRequirement struct {
	Resource Resource
	Action   Action
}
