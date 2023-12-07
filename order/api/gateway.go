package api

import (
	_ "embed"
	"github.com/go-saas/kit/pkg/apisix"
)

//go:embed gateway.yaml
var gateway []byte

func init() {
	apisix.LoadFromYaml(gateway)
}
