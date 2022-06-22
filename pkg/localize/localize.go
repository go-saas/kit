package localize

import (
	"context"
	"encoding/json"
	"gopkg.in/yaml.v3"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

type (
	localizerKey struct{}

	FileBundle struct {
		Buf  []byte
		Path string
	}
)

func I18N(files ...FileBundle) middleware.Middleware {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)
	for _, f := range files {
		bundle.MustParseMessageFileBytes(f.Buf, f.Path)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				accept := tr.RequestHeader().Get("accept-language")
				localizer := i18n.NewLocalizer(bundle, accept)
				ctx = context.WithValue(ctx, localizerKey{}, localizer)
			}
			return handler(ctx, req)
		}
	}
}

// FromContext resolve *i18n.Localizer from context. return nil if not found
func FromContext(ctx context.Context) *i18n.Localizer {
	if ret, ok := ctx.Value(localizerKey{}).(*i18n.Localizer); ok {
		return ret
	}
	return nil
}
