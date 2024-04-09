package client

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestHTTPClientCall_constructURLPath(t *testing.T) {
	t.Run("when constructURL Path with multiple Params", func(t *testing.T) {
		params := url.Values{}
		params.Add("key1", "value1")
		params.Add("key2", "value2")

		call := &HTTPClientCall{
			path:   "/test",
			params: params,
		}

		assert.Contains(t, call.constructURLPath(), "key1=value1", "El path debe contener el primer parámetro")
		assert.Contains(t, call.constructURLPath(), "key2=value2", "El path debe contener el segundo parámetro")
	})
}
