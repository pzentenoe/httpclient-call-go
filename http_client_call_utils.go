package client

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

func newClientRequest(ctx context.Context, method, host string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, host, nil)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (r *HTTPClientCall) constructURL() string {
	return fmt.Sprintf("%s%s", r.host, r.constructURLPath())
}

func (r *HTTPClientCall) constructURLPath() string {
	pathWithParams := r.path
	if len(r.params) > 0 {
		if r.isEncodeURL {
			pathWithParams += "?" + r.params.Encode()
		} else {
			pathWithParams += "?" + EncodeWithoutScapes(r.params)
			pathWithParams = strings.ReplaceAll(pathWithParams, " ", "+")
		}
	}
	return pathWithParams
}

func (r *HTTPClientCall) setHeaders(req *http.Request) {
	for key, values := range r.headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
}

func (r *HTTPClientCall) setRequestBody(req *http.Request) error {
	if r.body == nil {
		req.ContentLength = 0
		return nil
	}

	var serializedBody bytes.Buffer
	if err := json.NewEncoder(&serializedBody).Encode(r.body); err != nil {
		return err
	}

	if r.gzipCompress {
		var buf bytes.Buffer
		gz := gzip.NewWriter(&buf)
		if _, err := gz.Write(serializedBody.Bytes()); err != nil {
			return err
		}
		if err := gz.Close(); err != nil {
			return err
		}
		req.Header.Set("Content-Encoding", "gzip")
		req.Body = io.NopCloser(&buf)
		req.ContentLength = int64(buf.Len())
	} else {
		req.Body = io.NopCloser(&serializedBody)
		req.ContentLength = int64(serializedBody.Len())
	}

	return nil
}

func EncodeWithoutScapes(v url.Values) string {
	if v == nil {
		return ""
	}
	var buf strings.Builder
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := v[k]
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(k)
			buf.WriteByte('=')
			buf.WriteString(v)
		}
	}
	return buf.String()
}
