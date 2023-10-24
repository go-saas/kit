package menu

import (
	_ "embed"
	"github.com/go-saas/kit/sys/menu"
)

//go:embed menu.yaml
var menuData []byte

func init() {
	menu.LoadFromYaml(menuData)
}
