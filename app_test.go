package main

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_formatResponse(t *testing.T) {
	fixture := `{"apiVersion":"client.authentication.k8s.io/v1beta1","kind":"ExecCredential","status":{"token":"my-bearer-token"}}`
	actual, err := formatResponse(&response{})
	require.Equal(t, fixture, actual)
	require.Nil(t, err)
}

func Test_formatResponse_is_json(t *testing.T) {
	actual, err := formatResponse(&response{})
	require.True(t, json.Valid([]byte(actual)))
	require.Nil(t, err)
}

func Test_formatResponse_populate_defaults(t *testing.T) {
	actual, err := formatResponse(&response{})
	require.Contains(t, actual, "apiVersion")
	require.Nil(t, err)
}
func Test_formatResponse_override_defaults(t *testing.T) {
	actual, _ := formatResponse(&response{Kind: "foo"})
	require.Contains(t, actual, `"kind":"foo"`)
}

func Test_cli_info(t *testing.T) {
	cliInfo()
	require.Equal(t, "kubectl-passman", app.Name)
	require.Equal(t, "0.0.0", app.Version)
}

func Test_commands(t *testing.T) {
	cliCommands()
	require.Equal(t, "keychain", app.Commands[0].Name)
	require.Equal(t, "1password", app.Commands[1].Name)
	require.Equal(t, "gopass", app.Commands[2].Name)
}
