package authz

type ActionStr string

func (a ActionStr) GetIdentity() string {
	return string(a)
}

const (
	ListAction   ActionStr = "list"
	UpdateAction ActionStr = "update"
	CreateAction ActionStr = "create"
	DeleteAction ActionStr = "delete"
	GetAction    ActionStr = "get"
)
