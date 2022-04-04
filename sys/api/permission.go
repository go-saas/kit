package api

import (
	_ "embed"
	"github.com/goxiaoy/go-saas-kit/pkg/authz/authz"
)

const (
	ResourceMenu = "sys.menu"
	ResourceDev  = "dev"

	ResourceDevJaeger = "dev.jaeger"
)

//go:embed permission.yaml
var permission []byte

func init() {
	authz.LoadFromYaml(permission)
}
