package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	as3 "github.com/fclairamb/afero-s3"
	"github.com/goxiaoy/go-saas-kit/pkg/blob"
	"github.com/spf13/afero"
	"time"
)

func init() {
	blob.Register("s3", func(cfg blob.BlobConfig) blob.Blob {
		// You create a session
		sess, _ := session.NewSession(&aws.Config{
			Region:      aws.String(cfg.S3.Region),
			Credentials: credentials.NewStaticCredentials(cfg.S3.Key, cfg.S3.Secret, ""),
		})

		// Initialize the file system
		s3Fs := as3.NewFs(cfg.S3.Bucket, sess)

		return &Blob{
			Afero:  blob.NewAfs(blob.PatchOpt(cfg, s3Fs)),
			s3Api:  s3.New(sess),
			bucket: cfg.S3.Bucket,
		}
	})
}

type Blob struct {
	*afero.Afero
	s3Api  *s3.S3
	bucket string
}

func (b *Blob) GetAfero() *afero.Afero {
	return b.Afero
}

func (b *Blob) GeneratePreSignedURL(name string, expire time.Duration) (string, error) {
	r, _ := b.s3Api.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(b.bucket),
		Key:    aws.String(name),
	})

	// Create the pre-signed url with an expiry
	return r.Presign(expire)

}
