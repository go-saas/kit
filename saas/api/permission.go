package api

import (
	_ "embed"
	"github.com/goxiaoy/go-saas-kit/pkg/authz/authz"
)

const (
	ResourceTenant = "saas.tenant"
)

//go:embed permission.yaml
var permission []byte

func init() {
	authz.LoadFromYaml(permission)
}
