package kratos

import (
	"context"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
	http2 "net/http"
)

func ResolveHttpRequest(ctx context.Context) (*http2.Request, bool) {
	if t, ok := transport.FromServerContext(ctx); ok {
		if ht, ok := t.(*http.Transport); ok {
			return ht.Request(), true
		}
	}
	return nil, false
}
