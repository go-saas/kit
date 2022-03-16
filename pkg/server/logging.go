package server

import (
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
)

func PatchFilter(logger klog.Logger, c *conf.Logging) klog.Logger {
	res := logger
	if c == nil {
		return res
	}
	for _, f := range c.Filter {
		if f == nil || f.By == nil {
			continue
		}
		switch by := f.By.(type) {
		case *conf.Logging_Filter_Level:
			if by.Level != conf.Logging_ALL {
				res = klog.NewFilter(res, klog.FilterLevel(klog.ParseLevel(by.Level.String())))
			}
		default:
		}
	}
	return res
}
