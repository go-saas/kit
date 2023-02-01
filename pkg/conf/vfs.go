package conf

import (
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/env"
	vvfs "github.com/goxiaoy/vfs"
	"github.com/spf13/afero"
	"io"
	"io/fs"
	"strings"
)

var _ config.Source = (*vfs)(nil)

type vfs struct {
	v    vvfs.Blob
	path string
}

func NewVfs(v vvfs.Blob, path string) *vfs {
	return &vfs{
		v:    v,
		path: path,
	}
}

func (v *vfs) loadFile(path string) (*config.KeyValue, error) {
	file, err := v.v.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	return &config.KeyValue{
		Key:    info.Name(),
		Format: format(info.Name()),
		Value:  data,
	}, nil
}

func (v *vfs) loadDir(path string) (kvs []*config.KeyValue, err error) {
	err = afero.Walk(v.v, path, func(path string, file fs.FileInfo, err error) error {
		// ignore hidden files
		if file.IsDir() || strings.HasPrefix(file.Name(), ".") {
			return nil
		}
		kv, err := v.loadFile(path)
		if err != nil {
			return err
		}
		kvs = append(kvs, kv)
		return nil
	})
	return
}

func (v *vfs) Load() (kvs []*config.KeyValue, err error) {
	fi, err := v.v.Stat(v.path)
	if err != nil {
		return nil, err
	}
	if fi.IsDir() {
		return v.loadDir(v.path)
	}
	kv, err := v.loadFile(v.path)
	if err != nil {
		return nil, err
	}
	return []*config.KeyValue{kv}, nil
}

func (v *vfs) Watch() (config.Watcher, error) {
	//return empty watcher
	w, err := env.NewWatcher()
	if err != nil {
		return nil, err
	}
	return w, nil
}

func format(name string) string {
	if p := strings.Split(name, "."); len(p) > 1 {
		return p[len(p)-1]
	}
	return ""
}
