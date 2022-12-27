package localize

import (
	"context"
	"encoding/json"
	"gopkg.in/yaml.v3"
	"io"
	"io/fs"
	"path/filepath"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/go-saas/go-i18n/v2/i18n"
	"golang.org/x/text/language"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

type (
	localizerKey struct{}
	languagesKey struct{}
	FileBundle   struct {
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

// LanguageProvider resolve preferred languages settings from local or remote
type LanguageProvider interface {
	GetPrefer(ctx context.Context) ([]language.Tag, error)
}

type i18nOption struct {
	p           LanguageProvider
	defaultLang language.Tag
}

type Option func(*i18nOption)

// WithLanguageProvider set LanguageProvider
func WithLanguageProvider(p LanguageProvider) Option {
	return func(option *i18nOption) {
		option.p = p
	}
}

func WithDefaultLanguage(l language.Tag) Option {
	return func(option *i18nOption) {
		option.defaultLang = l
	}
}

func I18N(opts ...Option) middleware.Middleware {
	opt := &i18nOption{
		defaultLang: language.English,
	}
	for _, option := range opts {
		option(opt)
	}

	globalLock.RLock()
	defer globalLock.RUnlock()
	bundle := i18n.NewBundle(opt.defaultLang)
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
				b, err := io.ReadAll(f)
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
				tags := i18n.ParseTag([]string{accept})
				if opt.p != nil {
					settingTags, err := opt.p.GetPrefer(ctx)
					if err != nil {
						return nil, err
					}
					tags = append(settingTags, tags...)
				}
				tags = append(tags, opt.defaultLang)
				ctx = NewLanguageTagsContext(ctx, tags)
				localizer := i18n.NewLocalizerFromTags(bundle, tags...)
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

func LanguageTags(ctx context.Context) []language.Tag {
	if ret, ok := ctx.Value(languagesKey{}).([]language.Tag); ok {
		return ret
	}
	return nil
}

func NewLanguageTagsContext(ctx context.Context, tags []language.Tag) context.Context {
	return context.WithValue(ctx, languagesKey{}, tags)
}

func GetMsg(ctx context.Context, id, defaultMsg string, data map[string]interface{}, pluralCount interface{}) string {
	l := FromContext(ctx)
	if l == nil {
		return defaultMsg
	}
	msg, err := l.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: id,
		},
		TemplateData: data,
		PluralCount:  pluralCount,
	})
	if err != nil {
		return defaultMsg
	}
	return msg
}
