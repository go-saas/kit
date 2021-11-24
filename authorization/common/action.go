package common

type ActionStr string

func (a ActionStr) GetIdentity() string {
	return string(a)
}
