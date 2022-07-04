package localize

import (
	"context"
	"encoding/json"
	"gopkg.in/yaml.v3"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

type (
	localizerKey struct{}

	FileBundle struct {
		Fs fs.FS
	}
)

var (
	globalFileBundles []FileBundle
	globalLock        sync.RWMutex
)

func RegisterFileBundle(files ...FileBundle) {
	globalLock.Lock()
	defer globalLock.Unlock()
	globalFileBundles = append(globalFileBundles, files...)
}

func I18N() middleware.Middleware {
	globalLock.RLock()
	defer globalLock.RUnlock()
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)
	for _, f := range globalFileBundles {
		fs.WalkDir(f.Fs, ".", func(path string, d fs.DirEntry, err error) error {
			if filepath.Ext(path) == ".toml" || filepath.Ext(path) == ".json" || filepath.Ext(path) == ".yaml" {
				f, err := f.Fs.Open(path)
				if err != nil {
					panic(f)
				}
				defer f.Close()
				b, err := ioutil.ReadAll(f)
				if err != nil {
					panic(f)
				}
				bundle.MustParseMessageFileBytes(b, path)
			}
			return nil
		})
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
