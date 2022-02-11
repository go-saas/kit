package blob

import (
	"context"
	"fmt"
	"github.com/google/wire"
	"github.com/spf13/afero"
	"regexp"
	"sync"
	"time"
)

var (
	providerFactory = make(map[string]func(cfg *BlobConfig) Blob)
)

func Register(name string, f func(cfg *BlobConfig) Blob) {
	if _, ok := providerFactory[name]; ok {
		panic(fmt.Sprintf("provider %s already registered", name))
	}
	providerFactory[name] = f
}

type Factory interface {
	Get(ctx context.Context, name string) Blob
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
		cfg: cfg,
	}
}

func (f *FactoryImpl) Get(ctx context.Context, name string) Blob {
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
	r := factory(cfg)
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
}

type FileBlob struct {
	*afero.Afero
	Prefix string
}

func (f *FileBlob) GetAfero() *afero.Afero {
	return f.Afero
}

func (f *FileBlob) GeneratePreSignedURL(name string, expire time.Duration) (string, error) {
	return fmt.Sprintf("%s%s", f.Prefix, name), nil
}

func PatchOpt(cfg *BlobConfig, fs afero.Fs) afero.Fs {
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
