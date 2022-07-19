package service

import (
	"context"
	"github.com/go-saas/kit/pkg/localize"
	pb "github.com/go-saas/kit/sys/api/locale/v1"
)

type LocaleService struct {
	pb.UnimplementedLocaleServiceServer
}

func NewLocaleService() *LocaleService {
	return &LocaleService{}
}

func (s *LocaleService) ListMessages(ctx context.Context, req *pb.ListMessageRequest) (*pb.ListMessageReply, error) {
	loc := localize.FromContext(ctx)
	allMsg := loc.GetBundle().GetMessageTemplates()
	var items []*pb.LocaleLanguage
	for tag, m := range allMsg {
		var msg []*pb.LocaleMessage
		for _, template := range m {
			msg = append(msg, &pb.LocaleMessage{Id: template.ID, Other: template.Other})
		}
		items = append(items, &pb.LocaleLanguage{
			Name: tag.String(),
			Msg:  msg,
		})
	}
	return &pb.ListMessageReply{Items: items}, nil
}
