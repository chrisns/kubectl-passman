package main

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	jsonCertBase64 = `{"client-certificate-data":"MDAwMDA=","client-key-data":"MDAwMDA="}`
	jsonCert       = `{"clientCertificateData":"00000","clientKeyData":"00000"}`
	jsonToken      = `{"token":"00000"}`
)

func Test_formatValidatorCertBase64(t *testing.T) {
	actual, err := formatValidator(jsonCertBase64)
	require.Equal(t, jsonCert, actual)
	require.Nil(t, err)
}

func Test_formatValidatorCertBase64ErrorDecode(t *testing.T) {
	actual, err := formatValidator(`{"client-certificate-data":"BAD-DATA","client-key-data":"MDAwMDA="}`)
	require.Equal(t, "", actual)
	require.Equal(t, "illegal base64 data at input byte 3", err.Error())
}

func Test_formatValidatorCertBase64ErrorMisKey(t *testing.T) {
	actual, err := formatValidator(`{"clientCertificateData":"BAD-DATA","client-certificate-data":"MDAwMDA="}`)
	require.Equal(t, "", actual)
	require.Equal(t, "cannot define valid secret format", err.Error())
}

func Test_formatValidatorCertRaw(t *testing.T) {
	actual, err := formatValidator(jsonCert)
	require.Equal(t, jsonCert, actual)
	require.Nil(t, err)
}

func Test_formatValidatorCertRawError(t *testing.T) {
	actual, err := formatValidator(`{"clientCertificateData":"00000"}`)
	require.Equal(t, "", actual)
	require.Equal(t, "cannot define valid secret format", err.Error())
}

func Test_formatValidatorCertRawErrorMisKey(t *testing.T) {
	actual, err := formatValidator(`{"clientCertificateData":"00000"}`)
	require.Equal(t, "", actual)
	require.Equal(t, "cannot define valid secret format", err.Error())
}

func Test_formatValidatorToken(t *testing.T) {
	actual, err := formatValidator(jsonToken)
	require.Equal(t, jsonToken, actual)
	require.Nil(t, err)
}

func Test_formatResponse(t *testing.T) {
	fixture := `{"apiVersion":"client.authentication.k8s.io/v1beta1","kind":"ExecCredential","status":{}}`
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
	actual, err := formatResponse(&response{Kind: "foo"})
	require.Contains(t, actual, `"kind":"foo"`)
	require.Nil(t, err)
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
