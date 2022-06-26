package api

import (
	_ "embed"
	"github.com/go-saas/kit/pkg/authz/authz"
)

const (
	ResourceMenu = "sys.menu"
	ResourceDev  = "dev"

	ResourceDevJaeger = "dev.jaeger"

	ResourceDevJob = "dev.jobs"
)

//go:embed permission.yaml
var permission []byte

func init() {
	authz.LoadFromYaml(permission)
}
