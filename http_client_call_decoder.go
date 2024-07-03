package client

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// ResponseDecoder is an interface for decoding an HTTP response body.
type ResponseDecoder interface {
	Decode(io.Reader, any) error
}

// JSONResponseDecoder decodes JSON-encoded HTTP response bodies.
type JSONResponseDecoder struct{}

// Decode decodes a JSON-encoded response body into the provided interface.
func (d *JSONResponseDecoder) Decode(r io.Reader, v any) error {
	return json.NewDecoder(r).Decode(v)
}

// StringResponseDecoder decodes plain text or HTML-encoded HTTP response bodies.
type StringResponseDecoder struct{}

// Decode decodes a plain text or HTML-encoded response body into the provided string.
func (d *StringResponseDecoder) Decode(r io.Reader, v any) error {
	if str, ok := v.(*string); ok {
		data, err := io.ReadAll(r)
		if err != nil {
			return err
		}
		*str = string(data)
		return nil
	}
	return fmt.Errorf("StringResponseDecoder: unsupported type %T", v)
}

// selectDecoder selects the appropriate ResponseDecoder based on the Content-Type of the response.
func selectDecoder(contentType string) ResponseDecoder {
	switch {
	case strings.Contains(contentType, MIMEApplicationJSON):
		return &JSONResponseDecoder{}
	case strings.Contains(contentType, MIMETextPlain), strings.Contains(contentType, MIMETextHTML):
		return &StringResponseDecoder{}
	// Add more cases as needed
	default:
		return nil
	}
}
