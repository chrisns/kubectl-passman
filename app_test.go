package main

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

var fixture = `{"apiVersion":"client.authentication.k8s.io/v1beta1","kind":"ExecCredential","status":{"token":"my-bearer-token"}}`

func Test_returnResponse(t *testing.T) {
	t.Parallel()
	assert.Equal(t, returnResponse(), fixture)
}

func Test_returnResponse_is_json(t *testing.T) {
	assert.True(t, json.Valid([]byte(returnResponse())))
}

func Test_formatResponse_populate_defaults(t *testing.T) {
	t.Parallel()
	dummyResponse := &response{}
	assert.Contains(t, formatResponse(dummyResponse), "apiVersion")
}
func Test_formatResponse_override_defaults(t *testing.T) {
	t.Parallel()
	dummyResponse := &response{
		Kind: "foo",
	}
	assert.Contains(t, formatResponse(dummyResponse), `"kind":"foo"`)
}
