package os

import (
	"github.com/go-saas/kit/pkg/blob"
	"github.com/spf13/afero"
)

func init() {
	blob.Register("memory", func(cfg blob.BlobConfig) blob.Blob {
		// Initialize the file system
		mm := afero.NewMemMapFs()
		return blob.NewFileBlob(blob.NewAfs(blob.PatchOpt(cfg, mm)), cfg.BasePath, cfg.PublicUrl)
	})
}
