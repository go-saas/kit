package api

import (
	_ "embed"
	"github.com/go-saas/kit/pkg/authz/authz"
)

const (
	ResourcePermission = "permission"

	ResourceUser = "user.user"
	ResourceRole = "user.role"
)

//go:embed permission.yaml
var permission []byte

func init() {
	authz.LoadFromYaml(permission)
}
