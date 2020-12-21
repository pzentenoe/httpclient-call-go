package client

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

type hTTPClientRequest http.Request

func newHTTPClientRequest(method, host string) (*hTTPClientRequest, error) {
	req, err := http.NewRequest(method, host, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add(HeaderAccept, MIMEApplicationJSON)
	req.Header.Set(HeaderContentType, MIMEApplicationJSON)
	return (*hTTPClientRequest)(req), nil
}

func (r *hTTPClientRequest) setBody(body interface{}, gzipCompress bool) error {
	switch b := body.(type) {
	case string:
		if gzipCompress {
			return r.setBodyGzip(b)
		}
		return r.setBodyString(b)
	default:
		if gzipCompress {
			return r.setBodyGzip(body)
		}
		return r.setBodyJSON(body)
	}
}

func (r *hTTPClientRequest) setBodyJSON(data interface{}) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}
	r.Header.Set(HeaderContentType, MIMEApplicationJSON)
	if err := r.setBodyReader(bytes.NewReader(body)); err != nil {
		return err
	}
	return nil
}

func (r *hTTPClientRequest) setBodyGzip(body interface{}) error {
	switch b := body.(type) {
	case string:
		buf := new(bytes.Buffer)
		w := gzip.NewWriter(buf)
		if _, err := w.Write([]byte(b)); err != nil {
			return err
		}
		if err := w.Close(); err != nil {
			return err
		}
		r.Header.Add(HeaderContentEncoding, "gzip")
		r.Header.Add(HeaderVary, HeaderAcceptEncoding)
		return r.setBodyReader(bytes.NewReader(buf.Bytes()))
	default:
		data, err := json.Marshal(b)
		if err != nil {
			return err
		}
		buf := new(bytes.Buffer)
		w := gzip.NewWriter(buf)
		if _, err := w.Write(data); err != nil {
			return err
		}
		if err := w.Close(); err != nil {
			return err
		}
		r.Header.Add(HeaderContentEncoding, "gzip")
		r.Header.Add(HeaderVary, HeaderAcceptEncoding)
		r.Header.Set(HeaderContentType, MIMEApplicationJSON)
		return r.setBodyReader(bytes.NewReader(buf.Bytes()))
	}
}

func (r *hTTPClientRequest) setBodyString(body string) error {
	return r.setBodyReader(strings.NewReader(body))
}

func (r *hTTPClientRequest) setBodyReader(body io.Reader) error {
	rc, ok := body.(io.ReadCloser)
	if !ok && body != nil {
		rc = ioutil.NopCloser(body)
	}
	r.Body = rc
	if body != nil {
		switch v := body.(type) {
		case *strings.Reader:
			r.ContentLength = int64(v.Len())
		case *bytes.Buffer:
			r.ContentLength = int64(v.Len())
		}
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
