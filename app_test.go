package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_returnResponse(t *testing.T) {
	assert.Equal(t, `{"apiVersion": "client.authentication.k8s.io/v1beta1","kind": "ExecCredential","status": {"token": "my-bearer-token"}}`, returnResponse())
}
