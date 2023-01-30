package blob

import (
	"context"
	"fmt"
	"github.com/go-saas/saas"
	"github.com/goxiaoy/vfs"

	"path"
	"sync"
)

var (
	mutex           = sync.RWMutex{}
	providerFactory = make(map[string]func(cfg *Config) (vfs.Blob, error))
)

func Register(name string, f func(cfg *Config) (vfs.Blob, error)) {
	mutex.Lock()
	defer mutex.Unlock()
	if _, ok := providerFactory[name]; ok {
		panic(fmt.Sprintf("provider %s already registered", name))
	}
	providerFactory[name] = f
}

func GetKey(ctx context.Context, name string, tenancy bool) string {
	key := name
	if tenancy {
		ti, _ := saas.FromCurrentTenant(ctx)
		t := ti.GetId()
		if t == "" {
			t = "_"
		}
		key = path.Join(t, key)
	}
	return key
}

func (x *Config) Normalize() {
	if len(x.Provider) == 0 {
		if x.S3 != nil {
			x.Provider = "s3"
		}
		if x.Os != nil {
			x.Provider = "os"
		}
	}
}

func New(mounts ...*Config) (v *vfs.Vfs, err error) {
	for _, mount := range mounts {
		mount.Normalize()
		err = mount.ValidateAll()
		if err != nil {
			return nil, err
		}
	}
	v = vfs.New()
	mutex.RLock()
	defer mutex.RUnlock()
	for _, mount := range mounts {
		factory, ok := providerFactory[mount.Provider]
		if !ok {
			panic(fmt.Sprintf("blob provider %s not registered", mount.Provider))
		}
		var b vfs.Blob
		b, err = factory(mount)
		if err != nil {
			return nil, err
		}
		err = v.Mount(mount.MountPath, b)
		if err != nil {
			return nil, err
		}
	}
	return
}
