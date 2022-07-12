package apisix

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	ErrOption = errors.New("apisix endpoint and apiKey are required")
)

type Option struct {
	Endpoint string
	ApiKey   string
	Services []string
	Timeout  time.Duration
	Log      klog.Logger
}

type WatchSyncAdmin struct {
	discovery registry.Discovery
	client    *http.Client
	opt       *Option
	canceler  context.CancelFunc
	stopWg    *sync.WaitGroup
}

type watcherSync struct {
	service    string
	w          registry.Watcher
	ctx        context.Context
	updateFunc func(ins []*registry.ServiceInstance) error
	wg         *sync.WaitGroup
}

func (r *watcherSync) watch() {
	defer r.wg.Done()
	for {
		select {
		case <-r.ctx.Done():
			return
		default:
		}
		ins, err := r.w.Next()
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			klog.Errorf("[apisix] Failed to watch discovery endpoint: %v", err)
			time.Sleep(time.Second)
			continue
		}
		err = r.updateFunc(ins)
		klog.Errorf("[apisix] Failed to update service %s : %v", r.service, err)
	}
}

type Node struct {
	Host   string `json:"host"`
	Port   uint64 `json:"port"`
	Weight int    `json:"weight"`
}

type Upstream struct {
	Nodes  []Node `json:"nodes"`
	Type   string `json:"type"`
	Schema string `json:"schema"`
}

func toUpstreams(serviceName string, srvs []*registry.ServiceInstance) (map[string]*Upstream, error) {
	var ret map[string]*Upstream
	//group by schemas
	endpoints := lo.FlatMap(srvs, func(t *registry.ServiceInstance, _ int) []string {
		return t.Endpoints
	})

	for _, endpoint := range endpoints {
		raw, err := url.Parse(endpoint)
		if err != nil {
			return nil, err
		}
		addr := raw.Hostname()
		port, _ := strconv.ParseUint(raw.Port(), 10, 16)
		srvName := serviceName + "-" + raw.Scheme
		srv, ok := ret[srvName]
		if !ok {
			srv = &Upstream{Schema: raw.Scheme, Type: "roundrobin"}
			ret[srvName] = srv
		}
		srv.Nodes = append(srv.Nodes, Node{
			Host:   addr,
			Port:   port,
			Weight: 1,
		})
	}

	return ret, nil
}

func putServices(client *http.Client, endPoint, apiKey, serviceName string, ins []*registry.ServiceInstance) error {
	if strings.HasSuffix(endPoint, "/") {
		endPoint = strings.TrimSuffix(endPoint, "/")
	}
	up, err := toUpstreams(serviceName, ins)
	if err != nil {
		return err
	}
	for srvName, upstream := range up {
		j, err := json.Marshal(upstream)
		if err != nil {
			return err
		}
		klog.Infof("[apisix]  update service %s : %v", srvName, j)
		req, err := http.NewRequest(http.MethodPut, endPoint+"/apisix/admin/upstreams/"+srvName, bytes.NewReader(j))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-API-KEY", apiKey)
		_, err = client.Do(req)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewWatchSyncAdmin(discovery registry.Discovery, opt *Option) *WatchSyncAdmin {
	return &WatchSyncAdmin{discovery: discovery, opt: opt}
}

var _ transport.Server = (*WatchSyncAdmin)(nil)

func (w *WatchSyncAdmin) Start(ctx context.Context) (err error) {
	if w.opt == nil || len(w.opt.Endpoint) == 0 || len(w.opt.ApiKey) == 0 {
		return ErrOption
	}
	if w.opt.Log == nil {
		w.opt.Log = klog.GetLogger()
	}
	if w.opt.Timeout == 0 {
		w.opt.Timeout = 10 * time.Second
	}

	ctx, cancel := context.WithCancel(ctx)
	w.canceler = cancel

	//generate admin client
	w.client = &http.Client{}
	if err != nil {
		return err
	}

	if len(w.opt.Services) == 0 {
		w.opt.Log.Log(klog.LevelWarn, klog.DefaultMessageKey, "empty service list. will not sync to apisix admin")
	}

	//generate watcher for all services
	g, ctx := errgroup.WithContext(ctx)

	//start watcher wait group
	wg := sync.WaitGroup{}
	stopWg := sync.WaitGroup{}
	w.stopWg = &stopWg
	for _, service := range w.opt.Services {
		g.Go(func() error {
			defer wg.Done()
			watcher, err := w.discovery.Watch(ctx, service)
			if err != nil {
				//add watch into group
				s := &watcherSync{
					service: service,
					w:       watcher,
					ctx:     ctx,
					updateFunc: func(ins []*registry.ServiceInstance) error {
						return putServices(w.client, w.opt.Endpoint, w.opt.ApiKey, service, ins)
					},
					wg: &stopWg,
				}
				stopWg.Add(1)
				go s.watch()
			}
			return err
		})
		wg.Add(1)
	}
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()

	select {
	case <-c:
	case <-time.After(w.opt.Timeout):
		err = errors.New("discovery create watcher overtime")
	}

	if err := g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}

func (w *WatchSyncAdmin) Stop(ctx context.Context) error {
	//cancel all watcher
	if w.canceler != nil {
		w.canceler()
	}
	//should wait all watcher exists
	if w.stopWg != nil {
		w.stopWg.Wait()
	}

	return nil
}
