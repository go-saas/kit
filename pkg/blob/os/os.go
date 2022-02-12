package os

import (
	"github.com/goxiaoy/go-saas-kit/pkg/blob"
	"github.com/spf13/afero"
)

func init() {
	blob.Register("os", func(cfg blob.BlobConfig) blob.Blob {
		// Initialize the file system
		appfs := afero.NewOsFs()
		if cfg.Os != nil && cfg.Os.Dir != nil {
			appfs = afero.NewBasePathFs(appfs, cfg.Os.Dir.Value)
		}
		appfs = blob.PatchOpt(cfg, appfs)
		return blob.NewFileBlob(blob.NewAfs(appfs), cfg.BasePath, cfg.PublicUrl)
	})
}
