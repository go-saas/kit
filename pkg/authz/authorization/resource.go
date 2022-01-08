package authorization

type EntityResource struct {
	Namespace string
	Id        string
}

var _ Resource = (*EntityResource)(nil)

func NewEntityResource(namespace string, id string) *EntityResource {
	return &EntityResource{namespace, id}
}

func (r *EntityResource) GetNamespace() string {
	return r.Namespace
}

func (r *EntityResource) GetIdentity() string {
	return r.Id
}
