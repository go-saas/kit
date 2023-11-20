package api

import (
	_ "embed"
	"github.com/go-saas/kit/pkg/authz/authz"
)

const (
	ResourceUser = "user.user"
	ResourceRole = "user.role"

	ResourceAdminUser = "user.admin.user"
)

//go:embed permission.yaml
var permission []byte

func init() {
	authz.LoadFromYaml(permission)
}
