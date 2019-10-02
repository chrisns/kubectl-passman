package main

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_formatResponse(t *testing.T) {
	fixture := `{"apiVersion":"client.authentication.k8s.io/v1beta1","kind":"ExecCredential","status":{"token":"my-bearer-token"}}`
	require.Equal(t, fixture, formatResponse(&response{}))
}

func Test_formatResponse_is_json(t *testing.T) {
	require.True(t, json.Valid([]byte(formatResponse(&response{}))))
}

func Test_formatResponse_populate_defaults(t *testing.T) {
	require.Contains(t, formatResponse(&response{}), "apiVersion")
}
func Test_formatResponse_override_defaults(t *testing.T) {
	require.Contains(t, formatResponse(&response{Kind: "foo"}), `"kind":"foo"`)
}

func Test_opgetter_happy(t *testing.T) {
	var expected = "RSA"
	defaultOp = func(itemName string) (*opResponse, error) {
		return &opResponse{
			Details: opResponseDetails{
				Fields: []opResponseField{{
					Name:  "password",
					Value: expected,
				}},
			},
		}, nil
	}
	require.Contains(t, opgetter("gabriel"), expected)
}

func Test_opgetter_op_fail(t *testing.T) {
	err := errors.New("test")
	defaultOp = func(itemName string) (*opResponse, error) {
		return nil, err
	}
	require.PanicsWithValue(t, err, func() { opgetter("mykubecreds") })
}

func Test_opgetter_password_not_found(t *testing.T) {
	var expected = "RSA"
	defaultOp = func(itemName string) (*opResponse, error) {
		return &opResponse{
			Details: opResponseDetails{
				Fields: []opResponseField{{
					Name:  "notpassword",
					Value: expected,
				}},
			},
		}, nil
	}
	require.Panics(t, func() { opgetter("test") }) // TODO: panics with index out of range; is this expected behaviour?
}
