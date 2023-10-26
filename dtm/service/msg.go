package service

import (
	"context"
	"github.com/dtm-labs/dtm/client/dtmgrpc"
	"github.com/go-saas/kit/dtm/data"
	"github.com/go-saas/kit/dtm/utils"
	sapi "github.com/go-saas/kit/pkg/api"
	"github.com/go-saas/kit/pkg/dal"

	pb "github.com/go-saas/kit/dtm/api/dtm/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type MsgService struct {
	pb.UnimplementedMsgServiceServer
	provider dal.ConstDbProvider
	connName dal.ConnName
	trusted  sapi.TrustedContextValidator
}

func NewMsgService(provider dal.ConstDbProvider, connName dal.ConnName, trusted sapi.TrustedContextValidator) *MsgService {
	return &MsgService{
		provider: provider,
		connName: connName,
		trusted:  trusted,
	}
}

func (s *MsgService) QueryPrepared(ctx context.Context, req *pb.QueryPreparedRequest) (*emptypb.Empty, error) {
	if err := sapi.ErrIfUntrusted(ctx, s.trusted); err != nil {
		return nil, err
	}
	ba, err := utils.BarrierFromContext(ctx)
	if err != nil {
		return nil, err
	}
	db := data.GetDb(ctx, s.provider, s.connName)
	err = ba.QueryPrepared(utils.ToSQLDB(db))
	return &emptypb.Empty{}, dtmgrpc.DtmError2GrpcError(err)
}
