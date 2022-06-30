package service

import (
	"context"
	"github.com/dtm-labs/dtmcli"
	"github.com/dtm-labs/dtmgrpc"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-saas/kit/dtm/data"
	"github.com/go-saas/kit/pkg/dal"

	pb "github.com/go-saas/kit/dtm/api/dtm/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type MsgServiceService struct {
	pb.UnimplementedMsgServiceServer
	provider dal.ConstDbProvider
	connName dal.ConnName
}

func NewMsgService(provider dal.ConstDbProvider, connName dal.ConnName) *MsgServiceService {
	return &MsgServiceService{
		provider: provider,
		connName: connName,
	}
}

func (s *MsgServiceService) QueryPrepared(ctx context.Context, req *pb.QueryPreparedRequest) (*emptypb.Empty, error) {

	var ba *dtmcli.BranchBarrier
	var err error
	if t, ok := transport.FromServerContext(ctx); ok {
		if ht, ok := t.(*khttp.Transport); ok {
			ba, err = dtmcli.BarrierFromQuery(ht.Request().URL.Query())
			if err != nil {
				return nil, errors.BadRequest("BARRIER_INVALID", err.Error())
			}
		}
		if _, ok := t.(*grpc.Transport); ok {
			ba, err = dtmgrpc.BarrierFromGrpc(ctx)
			if err != nil {
				return nil, errors.BadRequest("BARRIER_INVALID", err.Error())
			}
		}
	} else {
		panic("can not resolve server context")
	}

	db := data.GetDb(ctx, s.provider, s.connName)
	err = ba.QueryPrepared(data.ToSQLDB(db))
	return &emptypb.Empty{}, dtmgrpc.DtmError2GrpcError(err)
}
