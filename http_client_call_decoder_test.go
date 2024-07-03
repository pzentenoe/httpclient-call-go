package client

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/suite"
)

type HTTPClientCallDecoderSuite struct {
	suite.Suite
}

func (suite *HTTPClientCallDecoderSuite) TestJSONResponseDecoder_Decode() {
	suite.Run("decodes JSON response successfully", func() {
		decoder := &JSONResponseDecoder{}
		body := `{"key":"value"}`
		r := io.NopCloser(bytes.NewBufferString(body))
		var result map[string]string

		err := decoder.Decode(r, &result)
		suite.NoError(err)
		suite.Equal("value", result["key"])
	})

	suite.Run("returns error for invalid JSON", func() {
		decoder := &JSONResponseDecoder{}
		body := `{"key":"value"`
		r := io.NopCloser(bytes.NewBufferString(body))
		var result map[string]string

		err := decoder.Decode(r, &result)
		suite.Error(err)
	})
}

func (suite *HTTPClientCallDecoderSuite) TestStringResponseDecoder_Decode() {
	suite.Run("decodes plain text response successfully", func() {
		decoder := &StringResponseDecoder{}
		body := "Hello, world!"
		r := io.NopCloser(bytes.NewBufferString(body))
		var result string

		err := decoder.Decode(r, &result)
		suite.NoError(err)
		suite.Equal("Hello, world!", result)
	})

	suite.Run("decodes HTML response successfully", func() {
		decoder := &StringResponseDecoder{}
		body := "<html><body>Hello, world!</body></html>"
		r := io.NopCloser(bytes.NewBufferString(body))
		var result string

		err := decoder.Decode(r, &result)
		suite.NoError(err)
		suite.Equal("<html><body>Hello, world!</body></html>", result)
	})

	suite.Run("returns error for unsupported type", func() {
		decoder := &StringResponseDecoder{}
		body := "Hello, world!"
		r := io.NopCloser(bytes.NewBufferString(body))
		var result int

		err := decoder.Decode(r, &result)
		suite.Error(err)
		suite.EqualError(err, "StringResponseDecoder: unsupported type *int")
	})
}

func (suite *HTTPClientCallDecoderSuite) TestSelectDecoder() {
	suite.Run("selects JSONResponseDecoder for application/json", func() {
		contentType := "application/json"
		decoder := selectDecoder(contentType)
		_, ok := decoder.(*JSONResponseDecoder)
		suite.True(ok)
	})

	suite.Run("selects StringResponseDecoder for text/plain", func() {
		contentType := "text/plain"
		decoder := selectDecoder(contentType)
		_, ok := decoder.(*StringResponseDecoder)
		suite.True(ok)
	})

	suite.Run("selects StringResponseDecoder for text/html", func() {
		contentType := "text/html"
		decoder := selectDecoder(contentType)
		_, ok := decoder.(*StringResponseDecoder)
		suite.True(ok)
	})

	suite.Run("returns nil for unsupported content type", func() {
		contentType := "application/xml"
		decoder := selectDecoder(contentType)
		suite.Nil(decoder)
	})
}

func TestHTTPClientCallDecoderSuite(t *testing.T) {
	suite.Run(t, new(HTTPClientCallDecoderSuite))
}
