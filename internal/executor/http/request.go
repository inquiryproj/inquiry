package http

import (
	"bytes"
	"context"
	"io"
	"net/http"
)

func (r Request) toHTTPRequest() (*http.Request, error) {
	var reader io.Reader
	if r.Body != "" {
		reader = bytes.NewReader([]byte(r.Body))
	}

	req, err := http.NewRequestWithContext(context.Background(), r.Method, r.URL, reader)
	if err != nil {
		return nil, err
	}

	for _, h := range r.Headers {
		req.Header.Set(h.Name, h.Value)
	}
	return req, nil
}
