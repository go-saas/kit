package plugins

import (
	"fmt"
	pkgHTTP "github.com/apache/apisix-go-plugin-runner/pkg/http"
	"github.com/apache/apisix-go-plugin-runner/pkg/log"
	"github.com/apache/apisix-go-plugin-runner/pkg/plugin"
	kerrors "github.com/go-kratos/kratos/v2/errors"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-saas/kit/pkg/api"
	"github.com/go-saas/kit/pkg/authn"
	"github.com/go-saas/kit/pkg/authn/jwt"
	"github.com/go-saas/kit/pkg/authn/session"
	"github.com/go-saas/kit/pkg/authz/authz"
	conf2 "github.com/go-saas/kit/pkg/conf"
	errors2 "github.com/go-saas/kit/pkg/errors"
	v1 "github.com/go-saas/kit/saas/api/tenant/v1"
	uapi "github.com/go-saas/kit/user/api"
	v12 "github.com/go-saas/kit/user/api/auth/v1"
	"github.com/go-saas/saas"
	shttp "github.com/go-saas/saas/http"
	"github.com/go-saas/sessions"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"net/http"
	"strings"
	"time"
)

func init() {
	err := plugin.RegisterPlugin(&KitAuthn{})
	if err != nil {
		log.Fatalf("failed to register plugin kit_authn: %s", err)
	}
	err = plugin.RegisterPlugin(&KitAuthz{})
	if err != nil {
		log.Fatalf("failed to register plugin kit_authz: %s", err)
	}
}

type KitAuthn struct {
	plugin.DefaultPlugin
}

type KitAuthConf struct {
}

var (
	WhiteAuthList = []string{
		"/oauth2",
		"/userinfo",
	}
)

var (
	tokenizer           jwt.Tokenizer
	tokenManager        api.TokenManager
	apiClient           *conf2.Client
	apiOpt              *api.Option
	sessionInfoStore    sessions.Store
	rememberStore       sessions.Store
	securityCfg         *conf2.Security
	userTenantValidator *uapi.UserTenantContrib
	refreshProvider     session.RefreshTokenProvider
	ts                  saas.TenantStore
	authService         authz.Service
	subjectResolver     authz.SubjectResolver
	saasWebConfig       *shttp.WebMultiTenancyOption
)

func Init(
	t jwt.Tokenizer,
	tmr api.TokenManager,
	tenantConfig *shttp.WebMultiTenancyOption,
	clientName api.ClientName,
	services *conf2.Services,
	security *conf2.Security,

	userTenant *uapi.UserTenantContrib,
	tenantStore saas.TenantStore,
	refreshTokenProvider session.RefreshTokenProvider,
	as authz.Service,
	sr authz.SubjectResolver,
	logger klog.Logger,
) error {
	tokenizer = t
	tokenManager = tmr
	saasWebConfig = tenantConfig
	clientCfg := &conf2.Client{Timeout: durationpb.New(1 * time.Second)}
	if c, ok := services.Clients[string(clientName)]; ok {
		proto.Merge(clientCfg, c)
	}
	if len(clientCfg.ClientId) == 0 {
		clientCfg.ClientId = string(clientName)
	}
	apiClient = clientCfg
	apiOpt = api.NewOption(true, api.NewUserPropagator(logger), api.NewClientPropagator(false, logger)).WithInsecure()
	securityCfg = security
	sessionInfoStore = session.NewSessionInfoStore(security)
	rememberStore = session.NewRememberStore(security)
	userTenantValidator = userTenant
	refreshProvider = refreshTokenProvider
	ts = tenantStore
	subjectResolver = sr
	authService = as
	return nil
}

func (p *KitAuthn) Name() string {
	return "kit_authn"
}

func (p *KitAuthn) ParseConf(in []byte) (interface{}, error) {
	return nil, nil
}

func abortWithError(r pkgHTTP.Request, err error, w http.ResponseWriter) {
	//use error codec
	fr := kerrors.FromError(err)

	//check tenant error
	if !isAjax(r) {
		//TODO should redirect to error page
		key := shttp.KeyOrDefault("")
		if fr.Reason == v1.ErrorReason_TENANT_NOT_FOUND.String() || fr.Reason == v1.ErrorReason_TENANT_FORBIDDEN.String() {
			//delete cookie
			sessions.SetCookie(w.Header(), sessions.NewCookie(key, "", &sessions.Options{MaxAge: -1}))
		}
	}
	w.WriteHeader(int(fr.Code))
	khttp.DefaultErrorEncoder(w, &http.Request{}, err)
}

func isAjax(r pkgHTTP.Request) bool {
	h := r.Header().View()["X-Requested-With"]
	return len(h) > 0 && h[0] == "XMLHttpRequest"
}

func (p *KitAuthn) RequestFilter(_ interface{}, w http.ResponseWriter, r pkgHTTP.Request) {
	currentPath := string(r.Path())
	_, shouldSkip := lo.Find(WhiteAuthList, func(s string) bool { return strings.HasPrefix(currentPath, s) })
	if shouldSkip {
		return
	}
	//clean internal headers
	for s, _ := range r.Header().View() {
		if strings.HasPrefix(strings.ToLower(s), api.InternalKeyPrefix) {
			//just clean headers
			log.Infof("clean untrusted internal header %s", s)
			r.Header().Set(s, "")
		}
	}

	var err error
	ctx := r.Context()
	propagator := propagation.NewCompositeTextMapPropagator(tracing.Metadata{}, propagation.Baggage{}, propagation.TraceContext{})
	ctx = propagator.Extract(ctx, propagation.HeaderCarrier(r.Header().View()))

	ctx, err = Saas(ctx, ts, saasWebConfig.DomainFormat, w, r)
	//format error
	if err != nil {
		abortWithError(r, err, w)
		//stop
		return
	}

	uid := ""
	clientId := ""

	//session auth
	header := r.Header().View()
	ctx = sessions.NewRegistryContext(ctx, header)

	s, _ := session.GetSession(ctx, header, sessionInfoStore, securityCfg)

	rs, _ := session.GetRememberSession(ctx, header, rememberStore, securityCfg)

	stateWriter := session.NewClientStateWriter(s, rs, r.RespHeader(), header)

	ctx = session.NewClientStateWriterContext(ctx, stateWriter)
	state := session.NewClientState(s, rs)
	ctx = session.NewClientStateContext(ctx, state)

	if len(state.GetUid()) > 0 {
		//set uid from cookie
		uid = state.GetUid()
	}

	rmToken := state.GetRememberToken()
	if len(state.GetUid()) == 0 && rmToken != nil {
		//call refresh
		log.Infof("call refresh token")
		_, err := refreshProvider.Refresh(ctx, rmToken.Token)
		if err != nil {
			err = kerrors.FromError(err)
			log.Errorf("refresh fail %v", err)
			if errors2.UnRecoverableError(err) {
				//abort with error
				abortWithError(r, err, w)
				return
			}
			if v12.IsRememberTokenUsed(err) {
				//for concurrent refresh, treat as logged in
				uid = rmToken.Uid
			}
		} else {
			uid = rmToken.Uid
		}

	}

	//extract token
	var t = ""
	if auth := r.Header().Get(jwt.AuthorizationHeader); len(auth) > 0 {
		t = jwt.ExtractHeaderToken(auth)
	}
	if len(t) == 0 {
		t = r.Args().Get(jwt.AuthorizationQuery)
	}
	//jwt auth
	if len(t) > 0 {
		if claims, err := jwt.ExtractAndValidate(tokenizer, t); err != nil {
			log.Errorf("fail to extract and validate token %s", err)
		} else {
			if claims.Subject != "" {
				uid = claims.Subject
			} else if claims.Subject != "" {
				uid = claims.Uid
			}
			clientId = claims.ClientId
		}
	}

	ctx = authn.NewUserContext(ctx, authn.NewUserInfo(uid))

	//check tenant and user mismatch
	ti, _ := saas.FromCurrentTenant(ctx)
	trCtx := saas.NewTenantResolveContext(ctx)
	trCtx.TenantIdOrName = ti.GetId()

	log.Infof("resolve user: %s client: %s tenantId: %s", uid, clientId, ti.GetId())
	err = userTenantValidator.Resolve(trCtx)
	if err != nil {
		// user can not in this tenant
		abortWithError(r, err, w)
		return
	}

	ctx = trCtx.Context()

	//keep previous client id
	ctx = authn.NewClientContext(ctx, clientId)
	ctx = authn.NewUserContext(ctx, authn.NewUserInfo(uid))

	//set auth token
	//use token mgr
	token, err := tokenManager.GetOrGenerateToken(ctx, &conf2.Client{
		ClientId:     apiClient.ClientId,
		ClientSecret: apiClient.ClientSecret,
	})
	if err != nil {
		log.Errorf("%s", err)
		abortWithError(r, err, w)
		return
	}
	//replace with internal token
	r.Header().Set(jwt.AuthorizationHeader, fmt.Sprintf("%s %s", jwt.BearerTokenType, token))

	headers := api.HeaderCarrier(map[string]string{})
	//inject header
	for _, contributor := range apiOpt.Propagators {
		//do not handle error
		contributor.Inject(ctx, headers)
	}
	for k, v := range headers {
		log.Infof("set header: %s value: %s", k, v)
		r.Header().Set(k, v)
	}

	return
}
