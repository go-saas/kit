package idp

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/miniprogram"
	miniConfig "github.com/silenceper/wechat/v2/miniprogram/config"
	"github.com/silenceper/wechat/v2/officialaccount"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/openplatform"
	openConfig "github.com/silenceper/wechat/v2/openplatform/config"
	"github.com/silenceper/wechat/v2/pay"
	payConfig "github.com/silenceper/wechat/v2/pay/config"
	"github.com/silenceper/wechat/v2/work"
	workConfig "github.com/silenceper/wechat/v2/work/config"
	"time"
)

func GetWeChatOpenPlatformOpenIDProvider(appID string) string {
	return fmt.Sprintf("wechat/openplatform/%s", appID)
}
func GetWeChatOfficialAccountOpenIDProvider(appID string) string {
	return fmt.Sprintf("wechat/officailaccount/%s", appID)
}
func GetWeChatMiniProgramOpenIDProvider(appID string) string {
	return fmt.Sprintf("wechat/miniprogoram/%s", appID)
}
func GetWeChatMiniGameOpenIDProvider(appID string) string {
	return fmt.Sprintf("wechat/minigame/%s", appID)
}
func GetWeChatPayOpenIDProvider(appID string) string {
	return fmt.Sprintf("wechat/pay/%s", appID)
}

func GetWeChatWorkOpenIDProvider(corpID string) string {
	return fmt.Sprintf("wechat/work/%s", corpID)
}

func GetWeChatUnionIDProvider() string {
	return "wechat/union"
}

func (x *WeChat) GetOpenPlatformByAppID(appID string) *WeChat_OpenPlatform {
	for _, m := range x.OpenPlatform {
		if m.AppId == appID {
			return m
		}
	}
	return nil
}

func (x *WeChat) GetOfficialAccountByAppID(appID string) *WeChat_OfficialAccount {
	for _, m := range x.OfficialAccount {
		if m.AppId == appID {
			return m
		}
	}
	return nil
}

func (x *WeChat) GetMiniProgramByAppID(appID string) *WeChat_MiniProgram {
	for _, m := range x.MiniProgram {
		if m.AppId == appID {
			return m
		}
	}
	return nil
}

func (x *WeChat) GetMiniGameByAppID(appID string) *WeChat_MiniGame {
	for _, m := range x.MiniGame {
		if m.AppId == appID {
			return m
		}
	}
	return nil
}

func (x *WeChat) GetPayByAppID(appID string) *WeChat_Pay {
	for _, m := range x.Pay {
		if m.AppId == appID {
			return m
		}
	}
	return nil
}

func (x *WeChat) GetWorkByCorpID(corpID string) *WeChat_Work {
	for _, m := range x.Work {
		if m.CorpId == corpID {
			return m
		}
	}
	return nil
}

const (
	cachePrefix = "wechat:"
)

type PrefixedCache struct {
	c cache.Cache
}

func NewPrefixedCache(c cache.Cache) *PrefixedCache {
	return &PrefixedCache{c: c}
}

func (p *PrefixedCache) Get(key string) interface{} {
	return p.Get(fmt.Sprintf("%s%s", cachePrefix, key))
}

func (p *PrefixedCache) Set(key string, val interface{}, timeout time.Duration) error {
	return p.Set(fmt.Sprintf("%s%s", cachePrefix, key), val, timeout)
}

func (p *PrefixedCache) IsExist(key string) bool {
	return p.IsExist(fmt.Sprintf("%s%s", cachePrefix, key))
}

func (p *PrefixedCache) Delete(key string) error {
	return p.Delete(fmt.Sprintf("%s%s", cachePrefix, key))
}

type WechatFactory struct {
	*wechat.Wechat
	cfg *WeChat
}

func (w *WechatFactory) GetOpenPlatformByAppID(ctx context.Context, appID string) (*openplatform.OpenPlatform, error) {
	cfg := w.cfg.GetOpenPlatformByAppID(appID)
	if cfg == nil {
		return nil, ErrorWechatConfigNotFoundLocalized(ctx, nil, nil)
	}
	return w.GetOpenPlatform(&openConfig.Config{
		AppID:          cfg.AppId,
		AppSecret:      cfg.AppSecret,
		Token:          cfg.Token,
		EncodingAESKey: cfg.EncodingAesKey,
	}), nil
}

func (w *WechatFactory) GetOfficialAccountByAppID(ctx context.Context, appID string) (*officialaccount.OfficialAccount, error) {
	cfg := w.cfg.GetOfficialAccountByAppID(appID)
	if cfg == nil {
		return nil, ErrorWechatConfigNotFoundLocalized(ctx, nil, nil)
	}
	return w.GetOfficialAccount(&offConfig.Config{
		AppID:          cfg.AppId,
		AppSecret:      cfg.AppSecret,
		Token:          cfg.Token,
		EncodingAESKey: cfg.EncodingAesKey,
	}), nil
}
func (w *WechatFactory) GetMiniProgramByAppID(ctx context.Context, appID string) (*miniprogram.MiniProgram, error) {
	cfg := w.cfg.GetMiniProgramByAppID(appID)
	if cfg == nil {
		return nil, ErrorWechatConfigNotFoundLocalized(ctx, nil, nil)
	}
	return w.GetMiniProgram(&miniConfig.Config{
		AppID:     cfg.AppId,
		AppSecret: cfg.AppSecret,
	}), nil
}

// func (w *WechatFactory) GetMiniGameByAppID(ctx context.Context, appID string) *WeChat_MiniGame {
//
// }

func (w *WechatFactory) GetPayByAppID(ctx context.Context, appID string) (*pay.Pay, error) {
	cfg := w.cfg.GetPayByAppID(appID)
	if cfg == nil {
		return nil, ErrorWechatConfigNotFoundLocalized(ctx, nil, nil)
	}
	return w.GetPay(&payConfig.Config{
		AppID:     cfg.AppId,
		MchID:     cfg.MchId,
		Key:       cfg.Key,
		NotifyURL: cfg.NotifyUrl,
	}), nil
}

func (w *WechatFactory) GetWorkByCorpID(ctx context.Context, corpID string) (*work.Work, error) {
	cfg := w.cfg.GetWorkByCorpID(corpID)
	if cfg == nil {
		return nil, ErrorWechatConfigNotFoundLocalized(ctx, nil, nil)
	}
	return w.GetWork(&workConfig.Config{
		CorpID:         cfg.CorpId,
		CorpSecret:     cfg.CorpSecret,
		AgentID:        cfg.AgentId,
		RasPrivateKey:  cfg.RasPrivateKey,
		Token:          cfg.Token,
		EncodingAESKey: cfg.EncodingAesKey,
	}), nil
}

func NewWeChat(client redis.UniversalClient, cfg *Config) *WechatFactory {
	wc := wechat.NewWechat()
	c := &cache.Redis{}
	c.SetConn(client)
	wc.SetCache(NewPrefixedCache(c))
	return &WechatFactory{
		Wechat: wc,
		cfg:    cfg.Wechat,
	}
}
