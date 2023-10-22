package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func httpRequestForStruct(t *testing.T, body any) *http.Request {
	bodyBytes, err := json.Marshal(body)
	assert.NoError(t, err)
	return &http.Request{
		Body: io.NopCloser(bytes.NewBuffer(bodyBytes)),
	}
}
