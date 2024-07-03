package client

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestHTTPClientCall_Do(t *testing.T) {

	t.Run("", func(t *testing.T) {
		client := &http.Client{}
		host := "http://example.com"

		call := NewHTTPClientCall(host, client)

		require.NotNil(t, call, "The HTTPClientCall instance should not be nil")
		assert.Equal(t, host, call.host, "The host should match the input")
	})

	t.Run("when method is not allowed", func(t *testing.T) {
		call := NewHTTPClientCall("http://example.com", &http.Client{})
		call.Method("INVALID")

		resp, err := call.Do(context.Background())

		assert.Nil(t, resp, "Expected no response for an invalid method")
		assert.Error(t, err, "Expected an error for an invalid method")
		assert.EqualError(t, err, errorMethodNotAllowed, "Error should indicate method not allowed")
	})

	t.Run("when do Get is Success", func(t *testing.T) {
		handler := newHandlerFunc(http.StatusOK, nil)
		server := httptest.NewServer(handler)
		defer server.Close()

		client := server.Client()
		call := NewHTTPClientCall(server.URL, client).
			Method(http.MethodGet)

		resp, err := call.Do(context.Background())

		require.NoError(t, err, "Expected no error for a successful request")
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected HTTP 200 OK status")
	})

	t.Run("when DoWithGzipCompression", func(t *testing.T) {
		expectedBody := "test body"
		handler := newHandlerFunc(http.StatusOK, []byte("test body"))
		server := httptest.NewServer(handler)
		defer server.Close()

		client := server.Client()
		call := NewHTTPClientCall(server.URL, client).
			Method(http.MethodPost).
			UseGzipCompress(true).
			Body(expectedBody)

		resp, err := call.Do(context.Background())

		require.NoError(t, err, "Expected no error for a successful request with gzip compression")
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected HTTP 200 OK status")
	})

	t.Run("when test connection POST return  OK", func(t *testing.T) {
		ctx := context.Background()
		handler := newHandlerFunc(http.StatusCreated, []byte("{}"))
		server := httptest.NewServer(handler)
		defer server.Close()

		dummyBody := make(map[string]any)
		dummyBody["age"] = 30
		dummyBody["name"] = "test"

		headers := http.Header{
			HeaderContentType: []string{MIMEApplicationJSON},
		}

		httpClient := NewHTTPClientCall(server.URL, server.Client()).
			Path("/dummypath").
			Method(http.MethodPut).
			Headers(headers).
			Body(dummyBody)

		response, err := httpClient.Do(ctx)
		if err == nil {
			defer func() {
				_ = response.Body.Close()
			}()
		}

		data, _ := io.ReadAll(response.Body)

		assert.NoError(t, err)
		assert.Equal(t, "{}", string(data))
	})

	t.Run("when test connection POST with ContentLength returns OK", func(t *testing.T) {
		ctx := context.Background()
		handler := newHandlerFunc(http.StatusCreated, []byte("{}"))
		server := httptest.NewServer(handler)

		defer server.Close()

		dummyBody := make(map[string]any)
		dummyBody["age"] = 30
		dummyBody["name"] = "test"

		headers := http.Header{
			HeaderContentType: []string{MIMEApplicationJSON},
		}

		httpClient := NewHTTPClientCall(server.URL, server.Client()).
			Path("/dummypath").
			Method(http.MethodPut).
			Headers(headers).
			Body(dummyBody)

		response, err := httpClient.Do(ctx)
		if err == nil {
			defer response.Body.Close()
		}

		data, _ := io.ReadAll(response.Body)

		assert.NoError(t, err)
		assert.Equal(t, "{}", string(data))
	})

	t.Run("when test connection POST without body and ContentLength equals 0 then returns OK", func(t *testing.T) {
		ctx := context.Background()
		handler := newHandlerFunc(http.StatusCreated, []byte("{}"))
		server := httptest.NewServer(handler)

		defer server.Close()

		headers := http.Header{
			HeaderContentType: []string{MIMEApplicationJSON},
		}

		httpClient := NewHTTPClientCall(server.URL, server.Client()).
			Path("/dummypath").
			Method(http.MethodPut).
			Headers(headers).
			Body(nil)

		response, err := httpClient.Do(ctx)
		if err == nil {
			defer response.Body.Close()
		}

		data, _ := io.ReadAll(response.Body)

		assert.NoError(t, err)
		assert.Equal(t, "{}", string(data))
	})
}

func TestHTTPClientCall_DoWithUnmarshal(t *testing.T) {
	t.Run("when server response OK", func(t *testing.T) {
		ctx := context.Background()
		handler := newHandlerFunc(http.StatusCreated, []byte(`{"name":"Pablo"}`))
		server := httptest.NewServer(handler)

		defer server.Close()

		type testBody struct {
			Name string `json:"name"`
		}

		var testBodyReponse *testBody

		params := url.Values{}
		params.Set("pageNumber", "1")
		params.Add("pageSize", "10")

		httpClient := NewHTTPClientCall(server.URL, server.Client()).
			UseGzipCompress(true).
			Path("/dummy").
			Params(params).
			Method(http.MethodGet)

		response, err := httpClient.DoWithUnmarshal(ctx, &testBodyReponse)

		assert.NoError(t, err)
		assert.Equal(t, "Pablo", testBodyReponse.Name)
		assert.Equal(t, http.StatusCreated, response.StatusCode)
	})
}

func newHandlerFunc(statusCode int, body []byte) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(statusCode)
		_, _ = writer.Write(body)
	}
}
