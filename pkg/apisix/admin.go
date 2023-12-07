package apisix

import (
	"bytes"
	"encoding/json"
	"errors"
	klog "github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	ErrOption = errors.New("apisix endpoint and apiKey are required")
)

type AdminClient struct {
	endpoint string
	apiKey   string
	client   *http.Client
}

func NewAdminClient(endpoint string, apiKey string) (*AdminClient, error) {
	if len(endpoint) == 0 || len(apiKey) == 0 {
		return nil, ErrOption
	}
	if strings.HasSuffix(endpoint, "/") {
		endpoint = strings.TrimSuffix(endpoint, "/")
	}
	return &AdminClient{
		endpoint: endpoint,
		apiKey:   apiKey,
		client:   &http.Client{},
	}, nil
}

// PutUpstream https://apisix.apache.org/docs/apisix/admin-api/#upstream-api
func (a *AdminClient) PutUpstream(id string, upstream *Upstream) error {
	if len(upstream.Name) == 0 {
		upstream.Name = id
	}
	j, err := json.Marshal(upstream)
	if err != nil {
		return err
	}
	klog.Infof("[apisix]  update upstream %s : %s", id, j)
	req, err := http.NewRequest(http.MethodPut, a.endpoint+"/apisix/admin/upstreams/"+id, bytes.NewReader(j))
	if err != nil {
		return err
	}
	body, err := a.do(req)
	if err != nil {
		return err
	}
	klog.Infof("[apisix]  update upstream response %s : %s", id, body)
	return nil
}

// PutUpstreamStruct https://apisix.apache.org/docs/apisix/admin-api/#upstream-api
func (a *AdminClient) PutUpstreamStruct(id string, upstream *structpb.Struct) error {
	if _, ok := upstream.Fields["name"]; !ok {
		upstream.Fields["name"] = structpb.NewStringValue(id)
	}
	j, err := protojson.Marshal(upstream)
	if err != nil {
		return err
	}
	klog.Infof("[apisix]  update upstream %s : %s", id, j)
	req, err := http.NewRequest(http.MethodPut, a.endpoint+"/apisix/admin/upstreams/"+id, bytes.NewReader(j))
	if err != nil {
		return err
	}
	body, err := a.do(req)
	if err != nil {
		return err
	}
	klog.Infof("[apisix]  update upstream response %s : %s", id, body)
	return nil
}

// PutRoute https://apisix.apache.org/docs/apisix/admin-api/#route-api
func (a *AdminClient) PutRoute(id string, route *structpb.Struct) error {
	if _, ok := route.Fields["name"]; !ok {
		route.Fields["name"] = structpb.NewStringValue(id)
	}
	j, err := protojson.Marshal(route)
	if err != nil {
		return err
	}
	klog.Infof("[apisix]  update route %s : %s", id, j)
	req, err := http.NewRequest(http.MethodPut, a.endpoint+"/apisix/admin/routes/"+id, bytes.NewReader(j))
	if err != nil {
		return err
	}
	body, err := a.do(req)
	if err != nil {
		return err
	}
	klog.Infof("[apisix]  update route response %s : %s", id, body)
	return nil
}

// PutStreamRoute https://apisix.apache.org/docs/apisix/admin-api/#stream-route-api
func (a *AdminClient) PutStreamRoute(id string, route *structpb.Struct) error {
	j, err := protojson.Marshal(route)
	if err != nil {
		return err
	}
	klog.Infof("[apisix]  update stream_routes %s : %s", id, j)
	req, err := http.NewRequest(http.MethodPut, a.endpoint+"/apisix/admin/stream_routes/"+id, bytes.NewReader(j))
	if err != nil {
		return err
	}
	body, err := a.do(req)
	if err != nil {
		return err
	}
	klog.Infof("[apisix]  update stream_routes response %s : %s", id, body)
	return nil
}

// PutGlobalRules https://apisix.apache.org/docs/apisix/admin-api/#global-rule-api
func (a *AdminClient) PutGlobalRules(id string, route *structpb.Struct) error {
	j, err := protojson.Marshal(route)
	if err != nil {
		return err
	}
	klog.Infof("[apisix]  update global rules %s : %s", id, j)
	req, err := http.NewRequest(http.MethodPut, a.endpoint+"/apisix/admin/global_rules/"+id, bytes.NewReader(j))
	if err != nil {
		return err
	}
	body, err := a.do(req)
	if err != nil {
		return err
	}
	klog.Infof("[apisix]  update route response %s : %s", id, body)
	return nil
}

func (a *AdminClient) do(req *http.Request) ([]byte, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", a.apiKey)
	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return body, errors.New(string(body))
	}
	return body, nil
}
