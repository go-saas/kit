package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-saas/kit/pkg/authz/authz"
	pb "github.com/go-saas/kit/product/api/product/v1"
	"github.com/go-saas/kit/product/private/biz"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ProductService struct {
	pb.UnimplementedProductServiceServer
	repo biz.ProductRepo
	auth authz.Service
}

func NewProductService(repo biz.ProductRepo, auth authz.Service) *ProductService {
	return &ProductService{repo: repo, auth: auth}
}

func (s *ProductService) ListProduct(ctx context.Context, req *pb.ListProductRequest) (*pb.ListProductReply, error) {
	ret := &pb.ListProductReply{}

	totalCount, filterCount, err := s.repo.Count(ctx, req)
	ret.TotalSize = int32(totalCount)
	ret.FilterSize = int32(filterCount)

	if err != nil {
		return ret, err
	}
	items, err := s.repo.List(ctx, req)
	if err != nil {
		return ret, err
	}
	rItems := lo.Map(items, func(g *biz.Product, _ int) *pb.Product {
		b := &pb.Product{}
		MapBizProduct2Pb(g, b)
		return b
	})

	ret.Items = rItems
	return ret, nil
}
func (s *ProductService) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.Product, error) {
	g, err := s.repo.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	if g == nil {
		return nil, errors.NotFound("", "")
	}
	res := &pb.Product{}
	MapBizProduct2Pb(g, res)
	return res, nil
}
func (s *ProductService) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.Product, error) {
	e := &biz.Product{}
	MapCreatePbProduct2Biz(req, e)
	err := s.repo.Create(ctx, e)
	if err != nil {
		return nil, err
	}
	res := &pb.Product{}
	MapBizProduct2Pb(e, res)
	return res, nil
}
func (s *ProductService) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.Product, error) {
	g, err := s.repo.Get(ctx, req.Product.Id)
	if err != nil {
		return nil, err
	}
	if g == nil {
		return nil, errors.NotFound("", "")
	}

	MapUpdatePbProduct2Biz(req.Product, g)
	if err := s.repo.Update(ctx, g.ID.String(), g, nil); err != nil {
		return nil, err
	}
	res := &pb.Product{}
	MapBizProduct2Pb(g, res)
	return res, nil
}
func (s *ProductService) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductReply, error) {
	g, err := s.repo.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if g == nil {
		return nil, errors.NotFound("", "")
	}

	err = s.repo.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteProductReply{Id: g.ID.String()}, nil
}
func MapBizProduct2Pb(a *biz.Product, b *pb.Product) {
	b.Id = a.ID.String()
	//b.Name = a.Name
	b.CreatedAt = timestamppb.New(a.CreatedAt)
}

func MapUpdatePbProduct2Biz(a *pb.UpdateProduct, b *biz.Product) {
	//b.Name = a.Name
}
func MapCreatePbProduct2Biz(a *pb.CreateProductRequest, b *biz.Product) {
	//b.Name = a.Name
}
