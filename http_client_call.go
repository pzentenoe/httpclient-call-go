package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Doer interface {
	Do() (*http.Response, error)
}

type HTTPClientCall struct {
	client       *http.Client
	host         string
	path         string
	params       url.Values
	isEncodeURL  bool
	method       string
	headers      http.Header
	body         interface{}
	gzipCompress bool
	contentType  string
}

func NewHTTPClientCall(host string, client *http.Client) *HTTPClientCall {
	if client == nil {
		panic("You must create client")
	}
	if host == "" {
		panic("empty host")
	}
	return &HTTPClientCall{
		host:        host,
		client:      client,
		isEncodeURL: true,
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

func (r *HTTPClientCall) Body(body interface{}) *HTTPClientCall {
	r.body = body
	return r
}

func (r *HTTPClientCall) UseGzipCompress(gzipCompress bool) *HTTPClientCall {
	r.gzipCompress = gzipCompress
	return r
}

func (r *HTTPClientCall) ContentType(contentType string) *HTTPClientCall {
	r.contentType = contentType
	return r
}

func (r *HTTPClientCall) Do() (*http.Response, error) {
	if r.host == "" {
		return nil, errors.New("empty host")
	}

	if err := r.validateHTTPMethod(); err != nil {
		return nil, err
	}

	pathWithParams := r.path
	if len(r.params) > 0 {
		if r.isEncodeURL {
			pathWithParams += "?" + r.params.Encode()
		} else {
			pathWithParams += "?" + EncodeWithoutScapes(r.params)
			pathWithParams = strings.ReplaceAll(pathWithParams, " ", "+")
		}
	}

	req, err := newHTTPClientRequest(r.method, fmt.Sprintf("%s%s", r.host, pathWithParams))
	if err != nil {
		return nil, err
	}

	if r.contentType != "" {
		req.Header.Set(HeaderContentType, r.contentType)
	}
	if r.body != nil {
		err = req.setBody(r.body, r.gzipCompress)
		if err != nil {
			return nil, err
		}
	}

	if len(r.headers) > 0 {
		for key, value := range r.headers {
			for _, v := range value {
				req.Header.Add(key, v)
			}
		}
	}

	resp, err := r.client.Do((*http.Request)(req).WithContext(context.Background()))
	r.params = nil
	r.body = nil
	return resp, err
}

func (r *HTTPClientCall) validateHTTPMethod() error {
	if r.method == "" {
		return errors.New("empty Method")
	}
	switch r.method {
	case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch:
		return nil
	default:
		return errors.New("method not allowed")
	}
}
