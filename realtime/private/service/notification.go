package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-saas/kit/pkg/authn"
	"github.com/go-saas/kit/pkg/utils"
	"github.com/go-saas/kit/realtime/private/biz"
	"github.com/samber/lo"

	pb "github.com/go-saas/kit/realtime/api/notification/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type NotificationService struct {
	pb.UnimplementedNotificationServiceServer
	repo biz.NotificationRepo
}

func NewNotificationService(repo biz.NotificationRepo) *NotificationService {
	return &NotificationService{repo: repo}
}

func (s *NotificationService) ListNotification(ctx context.Context, req *pb.ListNotificationRequest) (*pb.ListNotificationReply, error) {
	_, err := authn.ErrIfUnauthenticated(ctx)
	if err != nil {
		return nil, err
	}
	ret := &pb.ListNotificationReply{}
	req.Sort = []string{"-ID"}

	totalCount, filterCount, err := s.repo.Count(ctx, req)

	ret.TotalSize = int32(totalCount)
	ret.FilterSize = int32(filterCount)
	if err != nil {
		return ret, err
	}

	unread, err := s.repo.MyUnreadCount(ctx)
	if err != nil {
		return ret, err
	}
	ret.UnreadSize = unread

	cursorRet, err := s.repo.ListCursor(ctx, req)
	if err != nil {
		return nil, err
	}
	ret.NextBeforePageToken = cursorRet.Before
	ret.NextAfterPageToken = cursorRet.After
	ret.Items = lo.Map(cursorRet.Items, MapBizNotification2Pb)
	return ret, nil
}

func (s *NotificationService) GetNotification(ctx context.Context, req *pb.GetNotificationRequest) (*pb.Notification, error) {
	ui, err := authn.ErrIfUnauthenticated(ctx)
	if err != nil {
		return nil, err
	}
	entity, err := s.repo.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	if entity == nil || entity.UserId != ui.GetId() {
		return nil, errors.NotFound("", "")
	}

	return MapBizNotification2Pb(entity, 0), nil
}

func (s *NotificationService) ReadNotification(ctx context.Context, req *pb.ReadNotificationRequest) (*emptypb.Empty, error) {
	_, err := authn.ErrIfUnauthenticated(ctx)
	if err != nil {
		return nil, err
	}

	idFilter := req.Id
	if idFilter == "-" {
		idFilter = ""
	}
	err = s.repo.SetMyRead(ctx, idFilter)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *NotificationService) DeleteNotification(ctx context.Context, req *pb.DeleteNotificationRequest) (*pb.DeleteNotificationReply, error) {
	_, err := authn.ErrIfUnauthenticated(ctx)
	if err != nil {
		return nil, err
	}
	err = s.repo.DeleteMy(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteNotificationReply{}, nil
}

func MapBizNotification2Pb(entity *biz.Notification, _ int) *pb.Notification {
	return &pb.Notification{
		Id:       entity.ID,
		TenantId: entity.TenantId.String,
		Group:    entity.Group,
		Title:    entity.Title,
		Desc:     entity.Desc,
		Image:    entity.Image,
		Link:     entity.Link,
		Source:   entity.Source,
		UserId:   entity.UserId,
		Extra:    utils.Map2Structpb(entity.Extra),
		Level:    entity.Level,
		HasRead:  entity.HasRead,
	}
}
