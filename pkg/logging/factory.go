package logging

import (
	"encoding/json"
	"fmt"
	kzap "github.com/go-kratos/kratos/contrib/log/zap/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-saas/kit/pkg/conf"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type writerWrapper struct {
	*lumberjack.Logger
}

func (w writerWrapper) Sync() error {
	return nil
}

func init() {
	zap.RegisterSink("rolling", func(url *url.URL) (zap.Sink, error) {

		l := &lumberjack.Logger{
			Filename: url.Host + url.Path,
		}
		q := queryOptions{
			q: url.Query(),
		}
		l.MaxSize = q.int("max_size")
		l.MaxBackups = q.int("max_backups")
		l.MaxAge = q.int("max_age")
		l.LocalTime = q.bool("local_time")
		l.Compress = q.bool("compress")
		if q.err != nil {
			return nil, q.err
		}
		// any parameters left?
		if r := q.remaining(); len(r) > 0 {
			return nil, fmt.Errorf("logging:: unexpected option: %s", strings.Join(r, ", "))
		}
		s := writerWrapper{l}
		return s, nil
	})
}
func NewLogger(cfg *conf.Logging) (log.Logger, func(), error) {
	if cfg != nil {
		if cfg.Zap != nil {
			var zapCfg zap.Config
			jsonStr, _ := cfg.Zap.MarshalJSON()
			if err := json.Unmarshal(jsonStr, &zapCfg); err != nil {
				return nil, func() {}, err
			}
			l, err := zapCfg.Build()
			if err != nil {
				return nil, func() {}, err
			}
			return kzap.NewLogger(l), func() {
				l.Sync()
			}, nil
		}
	}
	l := log.NewStdLogger(os.Stdout)
	return l, func() {}, nil
}

type queryOptions struct {
	q   url.Values
	err error
}

func (o *queryOptions) string(name string) string {
	vs := o.q[name]
	if len(vs) == 0 {
		return ""
	}
	delete(o.q, name) // enable detection of unknown parameters
	return vs[len(vs)-1]
}

func (o *queryOptions) int(name string) int {
	s := o.string(name)
	if s == "" {
		return 0
	}
	i, err := strconv.Atoi(s)
	if err == nil {
		return i
	}
	if o.err == nil {
		o.err = fmt.Errorf("logging: invalid %s number: %s", name, err)
	}
	return 0
}

func (o *queryOptions) duration(name string) time.Duration {
	s := o.string(name)
	if s == "" {
		return 0
	}
	// try plain number first
	if i, err := strconv.Atoi(s); err == nil {
		if i <= 0 {
			// disable timeouts
			return -1
		}
		return time.Duration(i) * time.Second
	}
	dur, err := time.ParseDuration(s)
	if err == nil {
		return dur
	}
	if o.err == nil {
		o.err = fmt.Errorf("logging: invalid %s duration: %w", name, err)
	}
	return 0
}

func (o *queryOptions) bool(name string) bool {
	switch s := o.string(name); s {
	case "true", "1":
		return true
	case "false", "0", "":
		return false
	default:
		if o.err == nil {
			o.err = fmt.Errorf("logging: invalid %s boolean: expected true/false/1/0 or an empty string, got %q", name, s)
		}
		return false
	}
}

func (o *queryOptions) remaining() []string {
	if len(o.q) == 0 {
		return nil
	}
	keys := make([]string, 0, len(o.q))
	for k := range o.q {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
