package client

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type HTTPClientCallHeadersSuite struct {
	suite.Suite
	call *HTTPClientCall
	req  *http.Request
}

func (suite *HTTPClientCallHeadersSuite) SetupTest() {
	client := &http.Client{}
	suite.call = NewHTTPClientCall("http://example.com", client)
	suite.req, _ = http.NewRequest("GET", "http://example.com", nil)
}

func (suite *HTTPClientCallHeadersSuite) TestSetHeaders() {
	suite.Run("sets single header correctly", func() {
		headers := http.Header{
			HeaderContentType: []string{"application/json"},
		}
		suite.call.Headers(headers)
		suite.call.setHeaders(suite.req)

		assert.Equal(suite.T(), "application/json", suite.req.Header.Get(HeaderContentType))
	})

	suite.Run("sets multiple headers correctly", func() {
		headers := http.Header{
			HeaderContentType: []string{"application/json"},
			HeaderAccept:      []string{"application/json"},
		}
		suite.call.Headers(headers)
		suite.call.setHeaders(suite.req)

		assert.Equal(suite.T(), "application/json", suite.req.Header.Get(HeaderContentType))
		assert.Equal(suite.T(), "application/json", suite.req.Header.Get(HeaderAccept))
	})

	suite.Run("sets multiple values for a single header correctly", func() {
		headers := http.Header{
			HeaderAcceptEncoding: []string{"gzip", "deflate"},
		}
		suite.call.Headers(headers)
		suite.call.setHeaders(suite.req)

		assert.ElementsMatch(suite.T(), []string{"gzip", "deflate"}, suite.req.Header[HeaderAcceptEncoding])
	})
}

func TestHTTPClientCallHeadersSuite(t *testing.T) {
	suite.Run(t, new(HTTPClientCallHeadersSuite))
}
