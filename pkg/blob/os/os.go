package os

import (
	"github.com/goxiaoy/go-saas-kit/pkg/blob"
	"github.com/spf13/afero"
)

func init() {
	blob.Register("os", func(cfg blob.BlobConfig) blob.Blob {
		// Initialize the file system
		appfs := afero.NewOsFs()
		return blob.NewFileBlob(blob.NewAfs(blob.PatchOpt(cfg, appfs)), cfg.BasePath, cfg.PublicUrl)
	})
}
