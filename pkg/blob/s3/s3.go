package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-saas/kit/pkg/blob"
	"github.com/goxiaoy/vfs"
	s32 "github.com/goxiaoy/vfs/s3"
	"net/url"
	"time"
)

func init() {
	blob.Register("s3", func(cfg *blob.Config) (vfs.Blob, error) {
		// You create a session
		sess, _ := session.NewSession(&aws.Config{
			Region:      aws.String(cfg.S3.Region),
			Credentials: credentials.NewStaticCredentials(cfg.S3.Key, cfg.S3.Secret, ""),
		})

		public, err := url.Parse(cfg.PublicUrl)
		if err != nil {
			return nil, err
		}
		internal, err := url.Parse(cfg.InternalUrl)
		if err != nil {
			return nil, err
		}

		return s32.NewBlob(sess, cfg.S3.Bucket, *public, *internal, time.Hour*24*15), nil

	})
}
