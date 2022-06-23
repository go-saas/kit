package blob

import (
	"context"
	"fmt"
	"github.com/google/wire"
	"github.com/goxiaoy/go-saas"

	"github.com/spf13/afero"
	"net/url"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

var (
	providerFactory = make(map[string]func(cfg BlobConfig) Blob)
)

func Register(name string, f func(cfg BlobConfig) Blob) {
	if _, ok := providerFactory[name]; ok {
		panic(fmt.Sprintf("provider %s already registered", name))
	}
	providerFactory[name] = f
}

type Factory interface {
	Get(ctx context.Context, name string, tenancy bool) Blob
}

type FactoryImpl struct {
	m   map[string]Blob
	mtx sync.Mutex
	cfg Config
}

type Config map[string]*BlobConfig

var _ Factory = (*FactoryImpl)(nil)

func NewFactory(cfg Config) Factory {
	return &FactoryImpl{
		m:   map[string]Blob{},
		cfg: cfg,
	}
}

func (f *FactoryImpl) Get(ctx context.Context, name string, tenancy bool) Blob {
	f.mtx.Lock()
	defer f.mtx.Unlock()
	b, ok := f.m[name]
	if ok {
		return b
	}
	//resolve cfg
	cfg, ok := f.cfg[name]
	if !ok {
		panic(fmt.Sprintf("blob %s config  not found", name))
	}
	factory, ok := providerFactory[cfg.Provider]
	if !ok {
		panic(fmt.Sprintf("blob provider %s not registered", cfg.Provider))
	}
	opt := *cfg
	if tenancy {
		ti, _ := saas.FromCurrentTenant(ctx)
		t := ti.GetId()
		if t == "" {
			t = "_"
		}
		opt.BasePath = filepath.Join(t, opt.BasePath)
	}
	opt.BasePath = strings.Trim(opt.BasePath, "/")
	opt.PublicUrl = strings.TrimSuffix(opt.PublicUrl, "/")
	opt.InternalUrl = strings.TrimSuffix(opt.InternalUrl, "/")
	r := factory(opt)
	f.m[name] = r

	return r
}

func NewAfs(fs afero.Fs) *afero.Afero {
	afs := &afero.Afero{Fs: fs}
	return afs
}

type Blob interface {
	GetAfero() *afero.Afero
	GeneratePreSignedURL(name string, expire time.Duration) (string, error)
	GeneratePublicUrl(name string) (string, error)
	GenerateInternalUrl(name string) (string, error)
}

type FileBlob struct {
	*afero.Afero
	BasePath  string
	PublicUrl string
}

var _ Blob = (*FileBlob)(nil)

func NewFileBlob(a *afero.Afero, basePath, publicUrl string) *FileBlob {
	return &FileBlob{
		Afero:     a,
		BasePath:  basePath,
		PublicUrl: publicUrl,
	}
}

func (f *FileBlob) GetAfero() *afero.Afero {
	return f.Afero
}

func (f *FileBlob) GeneratePreSignedURL(name string, expire time.Duration) (string, error) {
	u, err := url.Parse(f.PublicUrl)
	if err != nil {
		return "", err
	}
	u.Path = path.Join(u.Path, f.BasePath, name)
	return u.String(), nil
}

func (f *FileBlob) GeneratePublicUrl(name string) (string, error) {
	u, err := url.Parse(f.PublicUrl)
	if err != nil {
		return "", err
	}
	u.Path = path.Join(u.Path, f.BasePath, name)
	return u.String(), nil
}

func (f *FileBlob) GenerateInternalUrl(name string) (string, error) {
	return path.Join(f.BasePath, name), nil
}

func PatchOpt(cfg BlobConfig, fs afero.Fs) afero.Fs {
	r := fs
	if cfg.BasePath != "" {
		r = afero.NewBasePathFs(r, cfg.BasePath)
	}
	if cfg.ReadOnly {
		r = afero.NewReadOnlyFs(r)
	}
	if cfg.RegexFilter != "" {
		r = afero.NewRegexpFs(r, regexp.MustCompile(cfg.RegexFilter))
	}
	return r
}

var ProviderSet = wire.NewSet(NewFactory)
