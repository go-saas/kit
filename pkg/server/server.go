package server

import (
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"
)

var DefaultProviderSet = wire.NewSet(wire.Value(ReqDecode), wire.Value(ResEncoder), wire.Value(ErrEncoder))

var (
	ReqDecode  http.DecodeRequestFunc  = http.DefaultRequestDecoder
	ResEncoder http.EncodeResponseFunc = http.DefaultResponseEncoder
	ErrEncoder http.EncodeErrorFunc    = http.DefaultErrorEncoder
)
