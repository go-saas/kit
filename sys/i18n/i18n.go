package i18n

import (
	"embed"
	_ "github.com/go-saas/kit/oidc/i18n"
	"github.com/go-saas/kit/pkg/localize"
)

var (
	//go:embed  embed
	f embed.FS
)

func init() {
	localize.RegisterFileBundle(localize.FileBundle{Fs: f})
}
