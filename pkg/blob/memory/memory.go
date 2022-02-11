package os

import (
	"github.com/goxiaoy/go-saas-kit/pkg/blob"
	"github.com/spf13/afero"
)

func init() {
	blob.Register("memory", func(cfg *blob.BlobConfig) blob.Blob {
		// Initialize the file system
		mm := afero.NewMemMapFs()
		return &blob.FileBlob{
			Afero:  blob.NewAfs(blob.PatchOpt(cfg, mm)),
			Prefix: "",
		}
	})
}
