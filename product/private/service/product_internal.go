package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	sapi "github.com/go-saas/kit/pkg/api"
	"github.com/go-saas/kit/pkg/utils"
	v12 "github.com/go-saas/kit/product/api/price/v1"
	pb "github.com/go-saas/kit/product/api/product/v1"
	"github.com/go-saas/kit/product/private/biz"
	"github.com/samber/lo"
)

func (s *ProductService) CreateInternalProduct(ctx context.Context, req *pb.CreateInternalProductRequest) (*pb.Product, error) {
	if err := sapi.ErrIfUntrusted(ctx, s.trusted); err != nil {
		return nil, err
	}
	//TODO implement me
	panic("implement me")
}

func (s *ProductService) GetInternalProduct(ctx context.Context, req *pb.GetInternalProductRequest) (*pb.Product, error) {
	if err := sapi.ErrIfUntrusted(ctx, s.trusted); err != nil {
		return nil, err
	}
	g, err := s.repo.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	if g == nil {
		return nil, errors.NotFound("", "")
	}
	res := &pb.Product{}
	s.MapBizProduct2Pb(ctx, g, res)
	return res, nil
}

func (s *ProductService) UpdateInternalProduct(ctx context.Context, req *pb.UpdateInternalProductRequest) (*pb.Product, error) {
	if err := sapi.ErrIfUntrusted(ctx, s.trusted); err != nil {
		return nil, err
	}
	g, err := s.repo.Get(ctx, req.Product.Id)
	if err != nil {
		return nil, err
	}
	if g == nil {
		return nil, errors.NotFound("", "")
	}

	if err := s.MapUpdatePbProduct2Biz(ctx, req.Product, g); err != nil {
		return nil, err
	}
	if err := s.repo.Update(ctx, g.ID.String(), g, nil); err != nil {
		return nil, err
	}
	res := &pb.Product{}
	s.MapBizProduct2Pb(ctx, g, res)
	return res, nil
}

func (s *ProductService) DeleteInternalProduct(ctx context.Context, req *pb.DeleteInternalProductRequest) (*pb.DeleteInternalProductReply, error) {
	if err := sapi.ErrIfUntrusted(ctx, s.trusted); err != nil {
		return nil, err
	}

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
	return &pb.DeleteInternalProductReply{Id: g.ID.String()}, nil
}

func (s *ProductService) MapCreateInternalPbProduct2Biz(ctx context.Context, a *pb.CreateInternalProductRequest, b *biz.Product) error {
	b.Title = a.Title
	b.ShortDesc = a.ShortDesc
	b.Desc = a.Desc
	b.Content = utils.Structpb2Map(a.Content)

	b.MainPic = mapPbMedia2Biz(a.MainPic)
	b.Medias = lo.Map(a.Medias, func(t *pb.ProductMedia, _ int) biz.ProductMedia {
		return *mapPbMedia2Biz(t)
	})

	b.Badges = lo.Map(a.Badges, func(t *pb.Badge, _ int) biz.Badge {
		r := &biz.Badge{}
		mapPbBadge2Biz(t, r)
		return *r
	})
	b.VisibleFrom = utils.Timepb2Time(a.VisibleFrom)
	b.VisibleTo = utils.Timepb2Time(a.VisibleTo)
	b.IsNew = a.IsNew

	if len(a.CategoryKeys) > 0 {
		c, err := s.categoryRepo.FindByKeys(ctx, a.CategoryKeys)
		if err != nil {
			return err
		}
		b.Categories = c
	}

	b.MainCategoryKey = a.MainCategoryKey

	b.Keywords = lo.Map(a.Keywords, func(t *pb.Keyword, _ int) biz.Keyword {
		r := &biz.Keyword{}
		mapPbKeyword2Biz(t, r)
		return *r
	})

	b.Model = a.Model
	b.BrandId = a.BrandId
	b.IsGiveaway = a.IsGiveaway

	b.Attributes = lo.Map(a.Attributes, func(t *pb.ProductAttribute, _ int) biz.ProductAttribute {
		r := &biz.ProductAttribute{}
		mapPbAttribute2Biz(t, r)
		return *r
	})

	b.MultiSku = a.MultiSku

	b.CampaignRules = lo.Map(a.CampaignRules, func(t *pb.CampaignRule, _ int) biz.CampaignRule {
		r := &biz.CampaignRule{}
		mapPbCampaignRule2Biz(t, r)
		return *r
	})

	b.NeedShipping = a.NeedShipping

	b.Stocks = lo.Map(a.Stocks, func(t *pb.Stock, _ int) biz.Stock {
		r := &biz.Stock{}
		mapPbStock2Biz(t, r)
		return *r
	})

	b.Prices = lo.Map(a.Prices, func(t *v12.CreatePriceRequest, _ int) biz.Price {
		r := &biz.Price{}
		mapPbCreatePrice2Biz(t, r)
		return *r
	})

	b.IsSaleable = a.IsSaleable
	b.SaleableFrom = utils.Timepb2Time(a.SaleableFrom)
	b.SaleableTo = utils.Timepb2Time(a.SaleableTo)
	b.Barcode = a.Barcode

	b.ManageInfo = biz.ProductManageInfo{}
	mapPbManageInfo2Biz(a.ManageInfo, &b.ManageInfo)

	return nil
}
