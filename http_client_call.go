package client

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
)

// Constants for error messages.
const (
	errorEmptyHost        = "empty host"
	errorEmptyMethod      = "empty method"
	errorMethodNotAllowed = "method not allowed"
)

// HTTPClientDoer is an interface for executing an HTTP request.
type HTTPClientDoer interface {
	Do() (*http.Response, error)
}

// HTTPClientCall encapsulates the configuration and execution of an HTTP request.
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

// NewHTTPClientCall creates a new HTTPClientCall with the specified host and HTTP client.
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

// Path sets the path for the HTTP request.
func (r *HTTPClientCall) Path(path string) *HTTPClientCall {
	r.path = path
	return r
}

// Params sets the URL parameters for the HTTP request.
func (r *HTTPClientCall) Params(params url.Values) *HTTPClientCall {
	r.params = params
	return r
}

// IsEncodeURL sets whether the URL should be encoded.
func (r *HTTPClientCall) IsEncodeURL(isEncodeURL bool) *HTTPClientCall {
	r.isEncodeURL = isEncodeURL
	return r
}

// Method sets the HTTP method for the request.
func (r *HTTPClientCall) Method(method string) *HTTPClientCall {
	r.method = method
	return r
}

// Headers sets the HTTP headers for the request.
func (r *HTTPClientCall) Headers(headers http.Header) *HTTPClientCall {
	r.headers = headers
	return r
}

// Body sets the body for the HTTP request.
func (r *HTTPClientCall) Body(body any) *HTTPClientCall {
	r.body = body
	return r
}

// UseGzipCompress sets whether the request body should be gzip compressed.
func (r *HTTPClientCall) UseGzipCompress(gzipCompress bool) *HTTPClientCall {
	r.gzipCompress = gzipCompress
	return r
}

// Do executes the HTTP request with the configured settings.
func (r *HTTPClientCall) Do(ctx context.Context) (*http.Response, error) {
	if r.host == "" {
		return nil, errors.New(errorEmptyHost)
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

// HTTPClientCallResponse encapsulates the response status code from an HTTP request.
type HTTPClientCallResponse struct {
	StatusCode int `json:"status_code"`
}

// DoWithUnmarshal executes the HTTP request and unmarshals the response body into the provided interface.
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

// validateHTTPMethod checks if the HTTP method is valid and allowed.
func (r *HTTPClientCall) validateHTTPMethod() error {
	if r.method == "" {
		return errors.New(errorEmptyMethod)
	}
	switch r.method {
	case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch:
		return nil
	default:
		return errors.New(errorMethodNotAllowed)
	}
}
