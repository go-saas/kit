package apisix

import (
	"context"
	"errors"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
	"net/url"
	"strconv"
	"sync"
	"time"
)

var (
	validSchemas = []string{"http", "https", "grpc", "grpcs", "tcp", "udp", "tls"}
)

type Option struct {
	Services []string
	Timeout  time.Duration
}

type WatchSyncAdmin struct {
	discovery registry.Discovery
	client    *AdminClient
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
		if err != nil {
			klog.Errorf("[apisix] Failed to update service %s : %v", r.service, err)
		}

	}
}

func toUpstreams(serviceName string, srvs []*registry.ServiceInstance) (map[string]*Upstream, error) {
	var ret = map[string]*Upstream{}
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
			schema := raw.Scheme
			if !lo.Contains(validSchemas, schema) {
				schema = ""
			}
			srv = &Upstream{Scheme: schema, Type: "roundrobin"}
			ret[srvName] = srv
		}
		srv.Nodes = append(srv.Nodes, &Node{
			Host:   addr,
			Port:   port,
			Weight: 1,
		})
	}

	return ret, nil
}

func putServices(client *AdminClient, serviceName string, ins []*registry.ServiceInstance) error {

	up, err := toUpstreams(serviceName, ins)
	if err != nil {
		return err
	}
	if len(up) == 0 {
		klog.Warnf("[apisix] Skipped to put service %s empty upstream", serviceName)
	}
	for srvName, upstream := range up {
		err = client.PutUpstream(srvName, upstream)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewWatchSyncAdmin(discovery registry.Discovery, adminClient *AdminClient, opt *Option) *WatchSyncAdmin {
	return &WatchSyncAdmin{discovery: discovery, client: adminClient, opt: opt}
}

var _ transport.Server = (*WatchSyncAdmin)(nil)

func (w *WatchSyncAdmin) Start(ctx context.Context) (err error) {

	if w.opt.Timeout == 0 {
		w.opt.Timeout = 10 * time.Second
	}

	ctx, cancel := context.WithCancel(ctx)
	w.canceler = cancel

	if err != nil {
		return err
	}

	if len(w.opt.Services) == 0 {
		klog.Warn("empty service list. will not sync to apisix admin")
	}

	//generate watcher for all services
	g := errgroup.Group{}

	//start watcher wait group
	wg := sync.WaitGroup{}
	stopWg := sync.WaitGroup{}
	w.stopWg = &stopWg
	for _, service := range w.opt.Services {
		service := service
		g.Go(func() error {
			defer wg.Done()
			watcher, err := w.discovery.Watch(ctx, service)
			if err == nil {
				//add watch into group
				s := &watcherSync{
					service: service,
					w:       watcher,
					ctx:     ctx,
					updateFunc: func(ins []*registry.ServiceInstance) error {
						return putServices(w.client, service, ins)
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
