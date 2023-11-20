package service

import (
	"context"
	"fmt"
	"github.com/centrifugal/centrifuge"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/go-saas/kit/pkg/authn"
	"net/http"
)

func NewCentrifugeNode(redisOpt *redis.UniversalOptions, logger klog.Logger) (*centrifuge.Node, error) {
	l := klog.NewHelper(logger)
	node, err := centrifuge.New(centrifuge.Config{
		LogLevel: centrifuge.LogLevelTrace,
		LogHandler: func(entry centrifuge.LogEntry) {
			level := klog.LevelDebug
			switch entry.Level {
			case centrifuge.LogLevelNone:
			case centrifuge.LogLevelTrace:
			case centrifuge.LogLevelDebug:
				break
			case centrifuge.LogLevelInfo:
				level = klog.LevelInfo
				break
			case centrifuge.LogLevelWarn:
				level = klog.LevelWarn
				break
			case centrifuge.LogLevelError:
				level = klog.LevelError
				break
			default:
				break
			}
			var kvs []interface{}
			kvs = append(kvs, klog.DefaultMessageKey, entry.Message)
			if entry.Fields != nil {
				for k, v := range entry.Fields {
					kvs = append(kvs, k, v)
				}
			}
			klog.Log(level, kvs...)
		},
	})
	if err != nil {
		return nil, err
	}
	if redisOpt != nil {
		redisShards, err := redisShardsFromOpt(node, redisOpt)
		if err != nil {
			return nil, err
		}
		// Using Redis  Broker here to scale nodes.
		broker, err := centrifuge.NewRedisBroker(node, centrifuge.RedisBrokerConfig{
			Shards: redisShards,
		})
		if err != nil {
			return nil, err
		}
		node.SetBroker(broker)

		presenceManager, err := centrifuge.NewRedisPresenceManager(node, centrifuge.RedisPresenceManagerConfig{
			Shards: redisShards,
		})
		if err != nil {
			return nil, err
		}
		node.SetPresenceManager(presenceManager)
	}

	node.OnConnecting(func(ctx context.Context, e centrifuge.ConnectEvent) (centrifuge.ConnectReply, error) {
		uid := ""
		if c, ok := centrifuge.GetCredentials(ctx); ok {
			uid = c.UserID
		}
		return centrifuge.ConnectReply{
			Subscriptions: map[string]centrifuge.SubscribeOptions{
				//notification channel
				fmt.Sprintf("notification#%s", uid): {},
			},
		}, nil
	})

	node.OnConnect(func(client *centrifuge.Client) {
		client.OnUnsubscribe(func(e centrifuge.UnsubscribeEvent) {
			l.Debugf("user %s unsubscribed from %s", client.UserID(), e.Channel)
		})
		client.OnDisconnect(func(e centrifuge.DisconnectEvent) {
			l.Debugf("user %s disconnected, disconnect: %s", client.UserID(), e.Disconnect)
		})
		transport := client.Transport()
		l.Debugf("user %s connected via %s", client.UserID(), transport.Name())
	})

	return node, nil
}

func redisShardsFromOpt(node *centrifuge.Node, redisOpt *redis.UniversalOptions) ([]*centrifuge.RedisShard, error) {
	var redisShardConfigs []centrifuge.RedisShardConfig
	for _, addr := range redisOpt.Addrs {
		redisShardConfigs = append(redisShardConfigs, centrifuge.RedisShardConfig{
			Address:        addr,
			User:           redisOpt.Username,
			Password:       redisOpt.Password,
			DB:             redisOpt.DB,
			ConnectTimeout: redisOpt.DialTimeout,
		})
	}

	var redisShards []*centrifuge.RedisShard
	for _, redisConf := range redisShardConfigs {
		redisShard, err := centrifuge.NewRedisShard(node, redisConf)
		if err != nil {
			return nil, err
		}
		redisShards = append(redisShards, redisShard)
	}
	return redisShards, nil
}

// Authentication middleware example. Centrifuge expects Credentials
// with current user ID set. Without provided Credentials client
// connection won't be accepted. Another way to authenticate connection
// is reacting to node.OnConnecting event where you may authenticate
// connection based on a custom token sent by a client in first protocol
// frame. See _examples folder in repo to find real-life auth samples
// (OAuth2, Gin sessions, JWT etc).
func auth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ui, _ := authn.FromUserContext(ctx)
		// Put authentication Credentials into request Context.
		// Since we don't have any session backend here we simply
		// set user ID as empty string. Users with empty ID called
		// anonymous users, in real app you should decide whether
		// anonymous users allowed to connect to your server or not.
		cred := &centrifuge.Credentials{
			UserID: ui.GetId(),
		}
		newCtx := centrifuge.SetCredentials(ctx, cred)
		r = r.WithContext(newCtx)
		h.ServeHTTP(w, r)
	})
}
