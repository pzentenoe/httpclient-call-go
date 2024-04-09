package client

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
)

const (
	ErrEmptyHost        = "empty host"
	ErrEmptyMethod      = "empty method"
	ErrMethodNotAllowed = "method not allowed"
)

type HttpClientDoer interface {
	Do() (*http.Response, error)
}

type HTTPClientCall struct {
	client       *http.Client
	method       string
	host         string
	path         string
	params       url.Values
	headers      http.Header
	body         any
	isEncodeURL  bool
	gzipCompress bool
}

func NewHTTPClientCall(host string, client *http.Client) *HTTPClientCall {
	if client == nil {
		panic("You must create client")
	}
	if host == "" {
		panic("empty host")
	}
	return &HTTPClientCall{
		client:       client,
		host:         host,
		path:         "",
		params:       nil,
		isEncodeURL:  true,
		method:       "",
		headers:      nil,
		body:         nil,
		gzipCompress: false,
	}
}

func (r *HTTPClientCall) Path(path string) *HTTPClientCall {
	r.path = path
	return r
}

func (r *HTTPClientCall) Params(params url.Values) *HTTPClientCall {
	r.params = params
	return r
}

func (r *HTTPClientCall) IsEncodeURL(isEncodeURL bool) *HTTPClientCall {
	r.isEncodeURL = isEncodeURL
	return r
}

func (r *HTTPClientCall) Method(method string) *HTTPClientCall {
	r.method = method
	return r
}

func (r *HTTPClientCall) Headers(headers http.Header) *HTTPClientCall {
	r.headers = headers
	return r
}

func (r *HTTPClientCall) Body(body any) *HTTPClientCall {
	r.body = body
	return r
}

func (r *HTTPClientCall) UseGzipCompress(gzipCompress bool) *HTTPClientCall {
	r.gzipCompress = gzipCompress
	return r
}

func (r *HTTPClientCall) Do(ctx context.Context) (*http.Response, error) {
	if r.host == "" {
		return nil, errors.New(ErrEmptyHost)
	}

	if err := r.validateHTTPMethod(); err != nil {
		return nil, err
	}
	fullURL := r.constructURL()
	req, err := newClientRequest(ctx, r.method, fullURL)
	if err != nil {
		return nil, err
	}

	if err = r.setRequestBody(req); err != nil {
		return nil, err
	}
	r.setHeaders(req)

	resp, err := r.client.Do(req)
	r.params = nil
	r.body = nil
	return resp, err
}

type HTTPClientCallResponse struct {
	StatusCode int `json:"status_code"`
}

func (r *HTTPClientCall) DoWithUnmarshal(ctx context.Context, responseBody any) (*HTTPClientCallResponse, error) {
	resp, err := r.Do(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &responseBody)
	if err != nil {
		return nil, err
	}
	httpClientCallResponse := &HTTPClientCallResponse{
		StatusCode: resp.StatusCode,
	}
	return httpClientCallResponse, nil
}

func (r *HTTPClientCall) validateHTTPMethod() error {
	if r.method == "" {
		return errors.New(ErrEmptyMethod)
	}
	switch r.method {
	case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch:
		return nil
	default:
		return errors.New(ErrMethodNotAllowed)
	}
}
