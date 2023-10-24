package service

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-saas/kit/pkg/blob"
	"github.com/go-saas/kit/product/private/biz"
	"github.com/google/uuid"
	"github.com/goxiaoy/vfs"
	"io"
	"mime"
	"os"
	"path/filepath"
	"strings"
)

type UploadService struct {
	mediaRepo biz.ProductMediaRepo
	blob      vfs.Blob
}

func NewUploadService(mediaRepo biz.ProductMediaRepo, blob vfs.Blob) *UploadService {
	return &UploadService{mediaRepo: mediaRepo, blob: blob}
}

func (s *UploadService) upload(ctx http.Context, basePath string, beforeUpload func(ctx context.Context) error) error {
	req := ctx.Request()
	//TODO do not know why should read form file first ...
	if _, _, err := req.FormFile("file"); err != nil {
		return err
	}

	h := ctx.Middleware(func(ctx context.Context, _ interface{}) (interface{}, error) {
		if beforeUpload != nil {
			err := beforeUpload(ctx)
			if err != nil {
				return nil, err
			}
		}
		file, handle, err := req.FormFile("file")
		if err != nil {
			return nil, err
		}
		defer file.Close()
		fileName := handle.Filename
		ext := filepath.Ext(fileName)
		normalizedName := fmt.Sprintf("%s/%s%s", basePath, uuid.New().String(), ext)

		err = s.blob.MkdirAll(basePath, 0755)
		if err != nil {
			return nil, err
		}
		f, err := s.blob.OpenFile(normalizedName, os.O_WRONLY|os.O_CREATE, 0o666)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		_, err = io.Copy(f, file)
		if err != nil {
			return nil, err
		}
		err = s.mediaRepo.Create(ctx, &biz.ProductMedia{
			ID:       normalizedName,
			MimeType: mime.TypeByExtension(ext),
			Name:     strings.TrimSuffix(fileName, ext),
		})
		if err != nil {
			return nil, err
		}

		url, _ := s.blob.PublicUrl(ctx, normalizedName)
		return &blob.BlobFile{
			Id:   normalizedName,
			Name: strings.TrimSuffix(fileName, ext),
			Url:  url.URL,
		}, nil
	})
	out, err := h(ctx, nil)
	if err != nil {
		return err
	}
	return ctx.Result(201, out)
}
