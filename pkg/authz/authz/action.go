package authz

type ActionStr string

func (a ActionStr) GetIdentity() string {
	return string(a)
}

const (
	AnyAction ActionStr = "*"

	CreateAction ActionStr = "create"
	UpdateAction ActionStr = "update"
	DeleteAction ActionStr = "delete"

	ReadAction  ActionStr = "read"
	WriteAction ActionStr = "write"
)
