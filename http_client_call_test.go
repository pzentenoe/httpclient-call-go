package client

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestHTTPClientCall_Do(t *testing.T) {

	t.Run("when test connection GET is OK", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(http.StatusOK)
			_, _ = writer.Write([]byte("{}"))
		}))
		defer server.Close()

		params := url.Values{}
		params.Set("pageNumber", "1")
		params.Add("pageSize", "10")

		httpClient := NewHTTPClientCall(server.Client()).
			Host(server.URL).
			Path("/dummy").
			Params(params).
			Method(http.MethodGet)

		response, err := httpClient.Do()
		if err == nil {
			defer response.Body.Close()
		}

		data, _ := io.ReadAll(response.Body)

		assert.NoError(t, err)
		assert.Equal(t, "{}", string(data))
	})

	t.Run("when test connection POST with gzip compress return 201", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(http.StatusCreated)
			_, _ = writer.Write([]byte("{}"))
		}))
		defer server.Close()

		dummyBody := make(map[string]interface{})
		dummyBody["age"] = 30
		dummyBody["name"] = "test"

		headers := http.Header{
			HeaderContentType: []string{MIMEApplicationJSON},
		}

		httpClient := NewHTTPClientCall(server.Client()).
			Host(server.URL).
			Path("/dummypath").
			Method(http.MethodPost).
			Headers(headers).
			Body(dummyBody).
			UseGzipCompress(true)

		response, err := httpClient.Do()
		if err == nil {
			defer response.Body.Close()
		}

		data, _ := io.ReadAll(response.Body)

		assert.NoError(t, err)
		assert.Equal(t, "{}", string(data))
	})

	t.Run("when test connection POST return  OK", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(http.StatusCreated)
			_, _ = writer.Write([]byte("{}"))
		}))
		defer server.Close()

		dummyBody := make(map[string]interface{})
		dummyBody["age"] = 30
		dummyBody["name"] = "test"

		headers := http.Header{
			HeaderContentType: []string{MIMEApplicationJSON},
		}

		httpClient := NewHTTPClientCall(server.Client()).
			Host(server.URL).
			Path("/dummypath").
			Method(http.MethodPut).
			Headers(headers).
			Body(dummyBody)

		response, err := httpClient.Do()
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
		server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(http.StatusOK)
			_, _ = writer.Write([]byte(`{"name":"Pablo"}`))
		}))
		defer server.Close()
		type testBody struct {
			Name string `json:"name"`
		}
		var testBodyReponse *testBody

		params := url.Values{}
		params.Set("pageNumber", "1")
		params.Add("pageSize", "10")

		httpClient := NewHTTPClientCall(server.Client()).
			Host(server.URL).
			UseGzipCompress(true).
			Path("/dummy").
			Params(params).
			Method(http.MethodGet)

		response, err := httpClient.DoWithUnmarshal(&testBodyReponse)

		assert.NoError(t, err)
		assert.Equal(t, "Pablo", testBodyReponse.Name)
		assert.Equal(t, http.StatusOK, response.StatusCode)
	})
}
