package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type bitwardenCall struct {
	args  []string
	stdin string
}

func fakeBitwarden(t *testing.T, reply func(args []string) (string, error)) *[]bitwardenCall {
	original := bitwardenRun
	t.Cleanup(func() { bitwardenRun = original })
	calls := &[]bitwardenCall{}
	bitwardenRun = func(args []string, stdin string) (string, error) {
		*calls = append(*calls, bitwardenCall{args, stdin})
		return reply(args)
	}
	return calls
}

func decodeBitwardenStdin(t *testing.T, stdin string) map[string]interface{} {
	raw, err := base64.StdEncoding.DecodeString(stdin)
	require.Nil(t, err)
	item := map[string]interface{}{}
	require.Nil(t, json.Unmarshal(raw, &item))
	return item
}

func Test_bitwardenGetter(t *testing.T) {
	calls := fakeBitwarden(t, func([]string) (string, error) { return `{"token":"a"}` + "\n", nil })
	actual, err := bitwardenGetter("my item")
	require.Nil(t, err)
	require.Equal(t, `{"token":"a"}`, actual)
	require.Equal(t, []string{"get", "password", "--", "my item"}, (*calls)[0].args)
}

func Test_bitwardenGetterError(t *testing.T) {
	expected := errors.New("Vault is locked")
	fakeBitwarden(t, func([]string) (string, error) { return "", expected })
	_, err := bitwardenGetter("x")
	require.Equal(t, expected, err)
}

func Test_bitwardenGetterEmpty(t *testing.T) {
	fakeBitwarden(t, func([]string) (string, error) { return "\n", nil })
	_, err := bitwardenGetter("x")
	require.NotNil(t, err)
}

func Test_bitwardenSetterCreates(t *testing.T) {
	calls := fakeBitwarden(t, func(args []string) (string, error) {
		if args[0] == "get" {
			return "", errors.New("exit status 1: Not found.")
		}
		return "", nil
	})
	secret := "{\"token\":\"a\"}\nwith newline"
	require.Nil(t, bitwardenSetter("kube", secret))
	require.Len(t, *calls, 2)
	require.Equal(t, []string{"create", "item"}, (*calls)[1].args)
	item := decodeBitwardenStdin(t, (*calls)[1].stdin)
	require.Equal(t, "kube", item["name"])
	require.EqualValues(t, 1, item["type"])
	require.Equal(t, secret, item["login"].(map[string]interface{})["password"])
	for _, c := range *calls {
		require.NotContains(t, strings.Join(c.args, " "), "token")
	}
}

func Test_bitwardenSetterEdits(t *testing.T) {
	calls := fakeBitwarden(t, func(args []string) (string, error) {
		if args[0] == "get" {
			return `{"id":"abc-1","type":1,"name":"kube","notes":"keep me","login":{"username":"u","password":"old"}}`, nil
		}
		return "", nil
	})
	require.Nil(t, bitwardenSetter("kube", "new"))
	require.Equal(t, []string{"edit", "item", "abc-1"}, (*calls)[1].args)
	item := decodeBitwardenStdin(t, (*calls)[1].stdin)
	require.Equal(t, "keep me", item["notes"])
	login := item["login"].(map[string]interface{})
	require.Equal(t, "new", login["password"])
	require.Equal(t, "u", login["username"])
}

func Test_bitwardenSetterLookupFails(t *testing.T) {
	calls := fakeBitwarden(t, func([]string) (string, error) { return "", errors.New("exit status 1: Vault is locked.") })
	err := bitwardenSetter("kube", "new")
	require.ErrorContains(t, err, "Vault is locked")
	require.Len(t, *calls, 1)
}

func Test_bitwardenSetterNotLogin(t *testing.T) {
	fakeBitwarden(t, func([]string) (string, error) { return `{"id":"a","type":2,"name":"n"}`, nil })
	require.ErrorContains(t, bitwardenSetter("n", "s"), "not a login item")
}

func Test_bitwardenSetterBadJSON(t *testing.T) {
	fakeBitwarden(t, func([]string) (string, error) { return "nope", nil })
	require.ErrorContains(t, bitwardenSetter("n", "s"), "cannot read")
}

func Test_bitwardenSetterNoID(t *testing.T) {
	fakeBitwarden(t, func([]string) (string, error) { return `{"type":1}`, nil })
	require.ErrorContains(t, bitwardenSetter("n", "s"), "no id")
}

func Test_bitwardenSetterWriteFails(t *testing.T) {
	expected := errors.New("create failed")
	fakeBitwarden(t, func(args []string) (string, error) {
		if args[0] == "get" {
			return "", errors.New("Not found.")
		}
		return "", expected
	})
	require.Equal(t, expected, bitwardenSetter("n", "s"))
}
