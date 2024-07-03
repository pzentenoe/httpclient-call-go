package client

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type MockHTTPClient struct {
	Response *http.Response
	Err      error
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.Response, m.Err
}

type HTTPClientCallSuite struct {
	suite.Suite
	client *MockHTTPClient
	host   string
}

func (suite *HTTPClientCallSuite) SetupTest() {
	suite.client = &MockHTTPClient{}
	suite.host = "http://example.com"
}

func (suite *HTTPClientCallSuite) TestNewHTTPClientCall() {
	suite.Run("creates new HTTPClientCall with valid inputs", func() {
		call := NewHTTPClientCall(suite.host, suite.client)
		suite.Equal(suite.host, call.host)
		suite.Equal(suite.client, call.client)
		suite.Empty(call.path)
		suite.Nil(call.params)
		suite.True(call.isEncodeURL)
		suite.Empty(call.method)
		suite.Nil(call.headers)
		suite.Nil(call.body)
		suite.False(call.gzipCompress)
	})

	suite.Run("panics on nil client", func() {
		suite.PanicsWithValue("You must create client", func() {
			NewHTTPClientCall(suite.host, nil)
		})
	})

	suite.Run("panics on empty host", func() {
		suite.PanicsWithValue("empty host", func() {
			NewHTTPClientCall("", suite.client)
		})
	})
}

func (suite *HTTPClientCallSuite) TestHTTPClientCall_Methods() {
	call := NewHTTPClientCall(suite.host, suite.client)

	suite.Run("sets path", func() {
		call.Path("/test-path")
		suite.Equal("/test-path", call.path)
	})

	suite.Run("sets params", func() {
		params := url.Values{"key": []string{"value"}}
		call.Params(params)
		suite.Equal(params, call.params)
	})

	suite.Run("sets isEncodeURL", func() {
		call.IsEncodeURL(false)
		suite.False(call.isEncodeURL)
	})

	suite.Run("sets method", func() {
		call.Method(http.MethodPost)
		suite.Equal(http.MethodPost, call.method)
	})

	suite.Run("sets headers", func() {
		headers := http.Header{"Header": []string{"value"}}
		call.Headers(headers)
		suite.Equal(headers, call.headers)
	})

	suite.Run("sets body", func() {
		body := map[string]string{"key": "value"}
		call.Body(body)
		suite.Equal(body, call.body)
	})

	suite.Run("sets gzipCompress", func() {
		call.UseGzipCompress(true)
		suite.True(call.gzipCompress)
	})
}

func (suite *HTTPClientCallSuite) TestHTTPClientCall_Do() {
	suite.Run("executes HTTP request successfully", func() {
		mockResponse := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(`{"key":"value"}`)),
		}
		suite.client.Response = mockResponse
		suite.client.Err = nil

		call := &HTTPClientCall{
			client: suite.client,
			host:   suite.host,
			method: http.MethodGet,
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		resp, err := call.Do(ctx)
		require.NoError(suite.T(), err)
		suite.Equal(mockResponse.StatusCode, resp.StatusCode)
	})
}

func (suite *HTTPClientCallSuite) TestHTTPClientCall_DoWithUnmarshal_JSON() {
	suite.Run("executes HTTP request and unmarshals JSON response", func() {
		mockResponse := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(`{"key":"value"}`)),
			Header:     http.Header{"Content-Type": []string{"application/json"}},
		}
		suite.client.Response = mockResponse
		suite.client.Err = nil

		call := &HTTPClientCall{
			client: suite.client,
			host:   suite.host,
			method: http.MethodGet,
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		var responseBody map[string]string
		resp, err := call.DoWithUnmarshal(ctx, &responseBody)
		require.NoError(suite.T(), err)
		suite.Equal(mockResponse.StatusCode, resp.StatusCode)
		suite.Equal("value", responseBody["key"])
	})
}

func (suite *HTTPClientCallSuite) TestHTTPClientCall_DoWithUnmarshal_HTML() {
	suite.Run("executes HTTP request and unmarshals HTML response", func() {
		mockResponse := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(`<html><body>Hello, world!</body></html>`)),
			Header:     http.Header{"Content-Type": []string{"text/html"}},
		}
		suite.client.Response = mockResponse
		suite.client.Err = nil

		call := &HTTPClientCall{
			client: suite.client,
			host:   suite.host,
			method: http.MethodGet,
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		var responseBody string
		resp, err := call.DoWithUnmarshal(ctx, &responseBody)
		require.NoError(suite.T(), err)
		suite.Equal(mockResponse.StatusCode, resp.StatusCode)
		suite.Equal("<html><body>Hello, world!</body></html>", responseBody)
	})
}

func (suite *HTTPClientCallSuite) TestHTTPClientCall_DoWithUnmarshal_String() {
	suite.Run("executes HTTP request and unmarshals plain text response", func() {
		mockResponse := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString("Hello, world!")),
			Header:     http.Header{"Content-Type": []string{"text/plain"}},
		}
		suite.client.Response = mockResponse
		suite.client.Err = nil

		call := &HTTPClientCall{
			client: suite.client,
			host:   suite.host,
			method: http.MethodGet,
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		var responseBody string
		resp, err := call.DoWithUnmarshal(ctx, &responseBody)
		require.NoError(suite.T(), err)
		suite.Equal(mockResponse.StatusCode, resp.StatusCode)
		suite.Equal("Hello, world!", responseBody)
	})
}

func (suite *HTTPClientCallSuite) TestHTTPClientCall_validateHTTPMethod() {
	call := &HTTPClientCall{}

	suite.Run("returns error for empty method", func() {
		err := call.validateHTTPMethod()
		suite.EqualError(err, errorEmptyMethod)
	})

	suite.Run("validates allowed methods", func() {
		validMethods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch}
		for _, method := range validMethods {
			suite.Run(method, func() {
				call.Method(method)
				err := call.validateHTTPMethod()
				suite.NoError(err)
			})
		}
	})

	suite.Run("returns error for invalid method", func() {
		call.Method("INVALID")
		err := call.validateHTTPMethod()
		suite.EqualError(err, errorMethodNotAllowed)
	})
}

func (suite *HTTPClientCallSuite) TestHTTPClientCall_Do_ErrorCases() {
	call := NewHTTPClientCall(suite.host, suite.client)

	suite.Run("returns error for empty host", func() {
		call.host = ""
		_, err := call.Do(context.Background())
		suite.EqualError(err, errorEmptyHost)
	})

	suite.Run("returns error for invalid method", func() {
		call.host = suite.host
		call.method = "INVALID"
		_, err := call.Do(context.Background())
		suite.EqualError(err, errorMethodNotAllowed)
	})
}

func TestHTTPClientCallSuite(t *testing.T) {
	suite.Run(t, new(HTTPClientCallSuite))
}
