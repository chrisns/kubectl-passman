package main

import (
	"encoding/json"
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

func Test_cli_info(t *testing.T) {
	cliInfo()
	require.Equal(t, app.Name, "kubectl-passman")
	require.Equal(t, app.Version, "0.0.0")
}

func Test_commands(t *testing.T) {
	cliCommands()
	require.Equal(t, app.Commands[0].Name, "keychain")
	require.Equal(t, app.Commands[1].Name, "1password")
}
