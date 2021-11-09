package common

type ActionStr string

func (a ActionStr) GetIdentity() string {
	return string(a)
}

type NamespaceStr string

func (n NamespaceStr) GetIdentity() string {
	return string(n)
}

type ResourceStr string

func (r ResourceStr) GetIdentity() string {
	return string(r)
}

type HasIdentity interface {
	GetIdentity() string
}

type Action interface {
	HasIdentity
}

type Namespace interface {
	HasIdentity
}

type Resource interface {
	HasIdentity
}

type Subject interface {
	HasIdentity
}
