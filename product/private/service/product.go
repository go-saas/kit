package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-saas/kit/pkg/authz/authz"
	"github.com/go-saas/kit/pkg/price"
	"github.com/go-saas/kit/pkg/utils"
	"github.com/go-saas/kit/product/api"
	v1 "github.com/go-saas/kit/product/api/category/v1"
	v12 "github.com/go-saas/kit/product/api/price/v1"
	pb "github.com/go-saas/kit/product/api/product/v1"
	"github.com/go-saas/kit/product/private/biz"
	"github.com/google/uuid"
	"github.com/goxiaoy/vfs"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ProductService struct {
	repo biz.ProductRepo
	auth authz.Service
	blob vfs.Blob
	*UploadService
}

var _ pb.ProductServiceServer = (*ProductService)(nil)
var _ pb.ProductInternalServiceServer = (*ProductService)(nil)

func NewProductService(
	repo biz.ProductRepo,
	auth authz.Service,
	upload *UploadService,
	blob vfs.Blob,
) *ProductService {
	return &ProductService{repo: repo, auth: auth, UploadService: upload, blob: blob}
}

func (s *ProductService) ListProduct(ctx context.Context, req *pb.ListProductRequest) (*pb.ListProductReply, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceProduct, "*"), authz.ReadAction); err != nil {
		return nil, err
	}
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
		s.MapBizProduct2Pb(ctx, g, b)
		return b
	})

	ret.Items = rItems
	return ret, nil
}

func (s *ProductService) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.Product, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceProduct, req.GetId()), authz.ReadAction); err != nil {
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

func (s *ProductService) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.Product, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceProduct, "*"), authz.WriteAction); err != nil {
		return nil, err
	}
	e := &biz.Product{}
	MapCreatePbProduct2Biz(req, e)
	err := s.repo.Create(ctx, e)
	if err != nil {
		return nil, err
	}
	res := &pb.Product{}
	s.MapBizProduct2Pb(ctx, e, res)
	return res, nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.Product, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceProduct, req.GetProduct().GetId()), authz.WriteAction); err != nil {
		return nil, err
	}
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
	s.MapBizProduct2Pb(ctx, g, res)
	return res, nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductReply, error) {
	if _, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceProduct, req.GetId()), authz.WriteAction); err != nil {
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
	return &pb.DeleteProductReply{Id: g.ID.String()}, nil
}

func (s *ProductService) MapBizProduct2Pb(ctx context.Context, a *biz.Product, b *pb.Product) {
	b.Id = a.ID.String()

	b.CreatedAt = timestamppb.New(a.CreatedAt)
	b.UpdatedAt = timestamppb.New(a.UpdatedAt)
	b.Version = a.Version.String
	b.TenantId = a.TenantId.String

	b.Title = a.Title

	b.ShortDesc = a.ShortDesc
	b.Desc = a.Desc
	b.Content = utils.Map2Structpb(a.Content)

	b.MainPic = mapBizMedia2Pb(ctx, s.blob, a.MainPic)
	b.Medias = lo.Map(a.Medias, func(t biz.ProductMedia, _ int) *pb.ProductMedia {
		return mapBizMedia2Pb(ctx, s.blob, &t)
	})

	b.Badges = lo.Map(a.Badges, func(t biz.Badge, _ int) *pb.Badge {
		r := &pb.Badge{}
		mapBizBadge2Pb(&t, r)
		return r
	})
	b.VisibleFrom = utils.Time2Timepb(a.VisibleFrom)
	b.VisibleTo = utils.Time2Timepb(a.VisibleTo)
	b.IsNew = a.IsNew
	b.Categories = lo.Map(a.Categories, func(t biz.ProductCategory, _ int) *v1.ProductCategory {
		r := &v1.ProductCategory{}
		mapBizCategory2Pb(&t, r)
		return r
	})
	if a.MainCategory != nil {
		r := &v1.ProductCategory{}
		mapBizCategory2Pb(a.MainCategory, r)
		b.MainCategory = r
	}

	b.Keywords = lo.Map(a.Keywords, func(t biz.KeyWord, _ int) *pb.Keyword {
		r := &pb.Keyword{}
		mapBizKeyword2Pb(&t, r)
		return r
	})

	b.Model = a.Model
	b.BrandId = a.BrandId
	b.IsGiveaway = a.IsGiveaway

	b.Attributes = lo.Map(a.Attributes, func(t biz.ProductAttribute, _ int) *pb.ProductAttribute {
		r := &pb.ProductAttribute{}
		mapBizAttribute2Pb(&t, r)
		return r
	})

	b.MultiSku = a.MultiSku

	b.CampaignRules = lo.Map(a.CampaignRules, func(t biz.CampaignRule, _ int) *pb.CampaignRule {
		r := &pb.CampaignRule{}
		mapBizCampaignRule2Pb(&t, r)
		return r
	})

	b.NeedShipping = a.NeedShipping
	b.Stocks = lo.Map(a.Stocks, func(t biz.Stock, _ int) *pb.Stock {
		r := &pb.Stock{}
		mapBizStock2Pb(&t, r)
		return r
	})

	b.Prices = lo.Map(a.Prices, func(t biz.Price, _ int) *v12.Price {
		r := &v12.Price{}
		mapBizPrice2Pb(ctx, &t, r)
		return r
	})

	b.IsSaleable = a.IsSaleable
	b.SaleableFrom = utils.Time2Timepb(a.SaleableFrom)
	b.SaleableTo = utils.Time2Timepb(a.SaleableTo)
	b.Barcode = a.Barcode

	manageInfo := &pb.ProductManageInfo{}
	mapBizManageInfo2Pb(&a.ManageInfo, manageInfo)

	b.ManageInfo = manageInfo
}

func MapUpdatePbProduct2Biz(a *pb.UpdateProduct, b *biz.Product) {
	//b.Name = a.Name
	//TODO
}
func MapCreatePbProduct2Biz(a *pb.CreateProductRequest, b *biz.Product) {
	//b.Name = a.Name
	//TODO
}

func mapBizMedia2Pb(ctx context.Context, v vfs.Blob, a *biz.ProductMedia) *pb.ProductMedia {
	if a == nil {
		return nil
	}
	b := &pb.ProductMedia{}
	mapMedia(ctx, v, a, b)
	return b
}

func mapPbMedia2Biz(a *pb.ProductMedia) *biz.ProductMedia {
	if a == nil {
		return nil
	}
	return &biz.ProductMedia{
		ID:       a.Id,
		Type:     a.Type,
		MimeType: a.MimeType,
		Name:     a.Name,
	}
}

func mapMedia(ctx context.Context, v vfs.Blob, a *biz.ProductMedia, b *pb.ProductMedia) {
	b.Id = a.ID
	b.Type = a.Type
	b.MimeType = a.MimeType
	b.Name = a.Name
	url, _ := v.PublicUrl(ctx, b.Id)
	b.Url = url.URL
}

func mapBizBadge2Pb(a *biz.Badge, b *pb.Badge) {
	b.Id = a.ID.String()
	b.Code = a.Code
	b.Label = a.Label
}

func mapPbBadge2Biz(a *pb.Badge, b *biz.Badge) {
	if len(a.Id) > 0 {
		b.ID = uuid.MustParse(a.Id)
	}
	b.Code = a.Code
	b.Label = a.Label
}

func mapBizKeyword2Pb(a *biz.KeyWord, b *pb.Keyword) {
	b.Text = a.Text
	b.Refer = a.Refer
}

func mapBizCategory2Pb(a *biz.ProductCategory, b *v1.ProductCategory) {
	b.Key = a.Key
	b.Name = a.Name
	b.Path = a.Path
	if a.ParentID != nil {
		b.Parent = *a.ParentID
	}
}

func mapBizAttribute2Pb(a *biz.ProductAttribute, b *pb.ProductAttribute) {
	b.Title = a.Title
}

func mapBizCampaignRule2Pb(a *biz.CampaignRule, b *pb.CampaignRule) {
	b.Rule = a.Rule
	b.Extra = utils.Map2Structpb(a.Extra)
}

func mapBizStock2Pb(a *biz.Stock, b *pb.Stock) {
	b.InStock = a.InStock
	b.Level = a.Level
	b.Amount = int32(a.Amount)
	b.DeliveryCode = a.DeliveryCode
}

func mapBizProductSku2Pb(ctx context.Context, blob vfs.Blob, a *biz.ProductSku, b *pb.ProductSku) {
	b.Id = a.ID.String()

	b.CreatedAt = timestamppb.New(a.CreatedAt)
	b.UpdatedAt = timestamppb.New(a.UpdatedAt)
	b.Version = a.Version.String
	b.TenantId = a.TenantId.String

	b.Title = a.Title

	b.MainPic = mapBizMedia2Pb(ctx, blob, a.MainPic)
	b.Medias = lo.Map(a.Medias, func(t biz.ProductMedia, _ int) *pb.ProductMedia {
		return mapBizMedia2Pb(ctx, blob, &t)
	})

	b.Prices = lo.Map(a.Prices, func(t biz.Price, _ int) *v12.Price {
		r := &v12.Price{}
		mapBizPrice2Pb(ctx, &t, r)
		return r
	})

	b.Stocks = lo.Map(a.Stocks, func(t biz.Stock, _ int) *pb.Stock {
		r := &pb.Stock{}
		mapBizStock2Pb(&t, r)
		return r
	})
	b.Keywords = lo.Map(a.Keywords, func(t biz.KeyWord, _ int) *pb.Keyword {
		r := &pb.Keyword{}
		mapBizKeyword2Pb(&t, r)
		return r
	})

	b.IsSaleable = a.IsSaleable
	b.SaleableFrom = utils.Time2Timepb(a.SaleableFrom)
	b.SaleableTo = utils.Time2Timepb(a.SaleableTo)
	b.Barcode = a.Barcode
}

func mapBizManageInfo2Pb(a *biz.ProductManageInfo, b *pb.ProductManageInfo) {
	b.Managed = a.Managed
	b.ManagedBy = a.ManagedBy
}

func mapBizPrice2Pb(ctx context.Context, a *biz.Price, b *v12.Price) {
	b.Id = a.ID.String()

	b.CreatedAt = timestamppb.New(a.CreatedAt)
	b.UpdatedAt = timestamppb.New(a.UpdatedAt)
	b.TenantId = a.TenantId.String

	b.Default = price.MustNewFromInt64(a.DefaultAmount, a.CurrencyCode).ToPricePb(ctx)
	b.Discounted = price.MustNewFromInt64(a.DiscountedAmount, a.CurrencyCode).ToPricePb(ctx)
	b.DiscountText = a.DiscountText
	b.DenyMoreDiscounts = a.DenyMoreDiscounts

	b.BillingSchema = string(a.BillingScheme)
	b.CurrencyOptions = lo.Map(a.CurrencyOptions, func(t biz.PriceCurrencyOption, i int) *v12.PriceCurrencyOption {
		r := &v12.PriceCurrencyOption{}
		mapBizCurrencyOption2Pb(ctx, &t, r)
		return r
	})

	if a.Recurring != nil {
		b.Recurring = &v12.PriceRecurring{}
		maoBizPriceRecurring2Pb(a.Recurring, b.Recurring)
	}
	b.Tiers = lo.Map(a.Tiers, func(t biz.PriceTier, i int) *v12.PriceTier {
		r := &v12.PriceTier{}
		mapBizPriceTier2Pb(&t, r)
		return r
	})
	b.TiersMode = string(a.TiersMode)
	b.TransformQuantity = &v12.PriceTransformQuantity{}
	mapBizPriceTransformQuantity2Pb(&a.TransformQuantity, b.TransformQuantity)
	b.Type = string(a.Type)

}

func mapBizCurrencyOption2Pb(ctx context.Context, a *biz.PriceCurrencyOption, b *v12.PriceCurrencyOption) {
	b.Default = price.MustNewFromInt64(a.DefaultAmount, a.CurrencyCode).ToPricePb(ctx)
	b.Discounted = price.MustNewFromInt64(a.DiscountedAmount, a.CurrencyCode).ToPricePb(ctx)
	b.DiscountText = a.DiscountText
	b.DenyMoreDiscounts = a.DenyMoreDiscounts
	b.CurrencyCode = a.CurrencyCode
	b.Tiers = lo.Map(a.Tiers, func(t biz.PriceCurrencyOptionTier, i int) *v12.PriceCurrencyOptionTier {
		r := &v12.PriceCurrencyOptionTier{}
		mapBizPriceCurrencyOptionTier2Pb(ctx, &t, r)
		return r
	})

}

func mapBizPriceCurrencyOptionTier2Pb(ctx context.Context, a *biz.PriceCurrencyOptionTier, b *v12.PriceCurrencyOptionTier) {
	b.FlatAmount = a.FlatAmount
	b.UnitAmount = a.UnitAmount
	b.UpTo = a.UpTo
}

func maoBizPriceRecurring2Pb(a *biz.PriceRecurring, b *v12.PriceRecurring) {
	b.Interval = string(a.Interval)
	b.IntervalCount = a.IntervalCount
	b.TrialPeriodDays = a.TrialPeriodDays
	b.AggregateUsage = string(a.AggregateUsage)
	b.UsageType = string(a.UsageType)
}

func mapBizPriceTier2Pb(a *biz.PriceTier, b *v12.PriceTier) {
	b.FlatAmount = a.FlatAmount
	b.UnitAmount = a.UnitAmount
	b.UpTo = a.UpTo
}

func mapBizPriceTransformQuantity2Pb(a *biz.PriceTransformQuantity, b *v12.PriceTransformQuantity) {
	b.DivideBy = a.DivideBy
	b.Round = string(a.Round)
}
func (s *ProductService) UploadMedias(ctx http.Context) error {
	return s.upload(ctx, biz.ProductMediaPath, func(ctx context.Context) error {
		_, err := s.auth.Check(ctx, authz.NewEntityResource(api.ResourceProduct, "*"), authz.WriteAction)
		return err
	})
}
