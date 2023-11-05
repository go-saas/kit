package api

import (
	_ "embed"
	"github.com/go-saas/kit/pkg/authz/authz"
)

const (
	ResourceSubscription = "payment.subscription"
)

//go:embed permission.yaml
var permission []byte

func init() {
	authz.LoadFromYaml(permission)
}
