package client

import (
	"bytes"
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type HTTPClientCallUtilsSuite struct {
	suite.Suite
	call *HTTPClientCall
}

func (suite *HTTPClientCallUtilsSuite) SetupTest() {
	client := &http.Client{}
	suite.call = NewHTTPClientCall("http://example.com", client)
}

func (suite *HTTPClientCallUtilsSuite) TestNewClientRequest() {
	suite.Run("creates a new HTTP request with the given context, method, and host", func() {
		ctx := context.Background()
		method := http.MethodGet
		host := "http://example.com"
		req, err := newClientRequest(ctx, method, host)

		require.NoError(suite.T(), err)
		suite.Equal(method, req.Method)
		suite.Equal(host, req.URL.String())
		suite.Equal(ctx, req.Context())
	})
}

func (suite *HTTPClientCallUtilsSuite) TestConstructURL() {
	suite.Run("constructs the full URL based on the host and path", func() {
		suite.call.Path("/test-path")
		fullURL := suite.call.constructURL()

		suite.Equal("http://example.com/test-path", fullURL)
	})
}

func (suite *HTTPClientCallUtilsSuite) TestConstructURLPath() {
	suite.Run("constructs the URL path including query parameters", func() {
		suite.call.Path("/test-path")
		params := url.Values{"key": []string{"value"}}
		suite.call.Params(params)
		pathWithParams := suite.call.constructURLPath()

		suite.Equal("/test-path?key=value", pathWithParams)
	})

	suite.Run("encodes URL parameters without escaping special characters", func() {
		suite.call.Path("/test-path")
		params := url.Values{"key": []string{"value with spaces"}}
		suite.call.Params(params)
		suite.call.IsEncodeURL(false)
		pathWithParams := suite.call.constructURLPath()

		suite.Equal("/test-path?key=value+with+spaces", pathWithParams)
	})
}

func (suite *HTTPClientCallUtilsSuite) TestSetRequestBody() {
	suite.Run("sets the request body with JSON encoding", func() {
		suite.call.Body(map[string]string{"key": "value"})
		req, _ := http.NewRequest(http.MethodPost, "http://example.com", nil)

		err := suite.call.setRequestBody(req)
		require.NoError(suite.T(), err)

		var buf bytes.Buffer
		buf.ReadFrom(req.Body)

		expectedBody := "{\"key\":\"value\"}\n"
		suite.JSONEq(expectedBody, buf.String())
		suite.Equal(int64(len(expectedBody)), req.ContentLength)
	})

	suite.Run("sets the request body with gzip compression", func() {
		suite.call.Body(map[string]string{"key": "value"})
		suite.call.UseGzipCompress(true)
		req, _ := http.NewRequest(http.MethodPost, "http://example.com", nil)

		err := suite.call.setRequestBody(req)
		require.NoError(suite.T(), err)

		var buf bytes.Buffer
		buf.ReadFrom(req.Body)

		suite.Equal("gzip", req.Header.Get("Content-Encoding"))
		suite.NotEqual(int64(len(`{"key":"value"}`)), req.ContentLength)
	})
}

func (suite *HTTPClientCallUtilsSuite) TestEncodeWithoutScapes() {
	suite.Run("encodes URL values without escaping special characters", func() {
		params := url.Values{"key": []string{"value with spaces", "another value"}}
		encoded := EncodeWithoutScapes(params)

		suite.Equal("key=value with spaces&key=another value", encoded)
	})
}

func TestHTTPClientCallUtilsSuite(t *testing.T) {
	suite.Run(t, new(HTTPClientCallUtilsSuite))
}
