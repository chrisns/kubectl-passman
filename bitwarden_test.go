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

func listReply(list string) func([]string) (string, error) {
	return func(args []string) (string, error) {
		if args[0] == "list" {
			return list, nil
		}
		return "", nil
	}
}

func Test_bitwardenGetter(t *testing.T) {
	calls := fakeBitwarden(t, listReply(`[{"id":"1","type":1,"name":"my item","login":{"password":"{\"token\":\"a\"}"}}]`))
	actual, err := bitwardenGetter("my item")
	require.Nil(t, err)
	require.Equal(t, `{"token":"a"}`, actual)
	require.Equal(t, []string{"list", "items", "--search=my item"}, (*calls)[0].args)
}

func Test_bitwardenGetterIgnoresSimilarNames(t *testing.T) {
	fakeBitwarden(t, listReply(`[
		{"id":"1","type":1,"name":"kube-prod-old","login":{"password":"wrong"}},
		{"id":"2","type":1,"name":"Kube","login":{"password":"wrong"}},
		{"id":"3","type":2,"name":"kube","login":{"password":"wrong"}},
		{"id":"4","type":1,"name":"kube","login":{"password":"right"}}]`))
	actual, err := bitwardenGetter("kube")
	require.Nil(t, err)
	require.Equal(t, "right", actual)
}

func Test_bitwardenGetterNotFound(t *testing.T) {
	fakeBitwarden(t, listReply(`[{"id":"1","type":1,"name":"kube-prod","login":{"password":"wrong"}}]`))
	_, err := bitwardenGetter("kube")
	require.ErrorContains(t, err, "no Bitwarden login item")
}

func Test_bitwardenGetterDuplicates(t *testing.T) {
	fakeBitwarden(t, listReply(`[{"id":"1","type":1,"name":"kube","login":{"password":"a"}},{"id":"2","type":1,"name":"kube","login":{"password":"b"}}]`))
	_, err := bitwardenGetter("kube")
	require.ErrorContains(t, err, "unique name")
}

func Test_bitwardenGetterError(t *testing.T) {
	expected := errors.New("Vault is locked")
	fakeBitwarden(t, func([]string) (string, error) { return "", expected })
	_, err := bitwardenGetter("x")
	require.Equal(t, expected, err)
}

func Test_bitwardenGetterEmptyPassword(t *testing.T) {
	fakeBitwarden(t, listReply(`[{"id":"1","type":1,"name":"x","login":{"password":""}}]`))
	_, err := bitwardenGetter("x")
	require.ErrorContains(t, err, "no password")
}

func Test_bitwardenGetterNoLogin(t *testing.T) {
	fakeBitwarden(t, listReply(`[{"id":"1","type":1,"name":"x"}]`))
	_, err := bitwardenGetter("x")
	require.ErrorContains(t, err, "no password")
}

func Test_bitwardenGetterBadJSON(t *testing.T) {
	fakeBitwarden(t, listReply(`nope`))
	_, err := bitwardenGetter("x")
	require.ErrorContains(t, err, "cannot read")
}

func Test_bitwardenSetterCreates(t *testing.T) {
	calls := fakeBitwarden(t, listReply(`[]`))
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

func Test_bitwardenSetterCreatesWhenOnlySimilarNamesExist(t *testing.T) {
	calls := fakeBitwarden(t, listReply(`[{"id":"1","type":1,"name":"kube-prod","login":{"password":"keep"}}]`))
	require.Nil(t, bitwardenSetter("kube", "new"))
	require.Equal(t, []string{"create", "item"}, (*calls)[1].args)
}

func Test_bitwardenSetterEdits(t *testing.T) {
	calls := fakeBitwarden(t, listReply(`[
		{"id":"other","type":1,"name":"kube-prod","login":{"password":"keep"}},
		{"id":"abc-1","type":1,"name":"kube","notes":"keep me","login":{"username":"u","password":"old"}}]`))
	require.Nil(t, bitwardenSetter("kube", "new"))
	require.Equal(t, []string{"edit", "item", "abc-1"}, (*calls)[1].args)
	item := decodeBitwardenStdin(t, (*calls)[1].stdin)
	require.Equal(t, "keep me", item["notes"])
	login := item["login"].(map[string]interface{})
	require.Equal(t, "new", login["password"])
	require.Equal(t, "u", login["username"])
}

func Test_bitwardenSetterEditsItemWithoutLogin(t *testing.T) {
	calls := fakeBitwarden(t, listReply(`[{"id":"a","type":1,"name":"n"}]`))
	require.Nil(t, bitwardenSetter("n", "s"))
	item := decodeBitwardenStdin(t, (*calls)[1].stdin)
	require.Equal(t, "s", item["login"].(map[string]interface{})["password"])
}

func Test_bitwardenSetterLookupFails(t *testing.T) {
	calls := fakeBitwarden(t, func([]string) (string, error) { return "", errors.New("exit status 1: Vault is locked.") })
	err := bitwardenSetter("kube", "new")
	require.ErrorContains(t, err, "Vault is locked")
	require.Len(t, *calls, 1)
}

func Test_bitwardenSetterDuplicates(t *testing.T) {
	calls := fakeBitwarden(t, listReply(`[{"id":"1","type":1,"name":"kube"},{"id":"2","type":1,"name":"kube"}]`))
	require.ErrorContains(t, bitwardenSetter("kube", "s"), "unique name")
	require.Len(t, *calls, 1)
}

func Test_bitwardenSetterNoID(t *testing.T) {
	fakeBitwarden(t, listReply(`[{"type":1,"name":"n"}]`))
	require.ErrorContains(t, bitwardenSetter("n", "s"), "no id")
}

func Test_bitwardenSetterWriteFails(t *testing.T) {
	expected := errors.New("create failed")
	fakeBitwarden(t, func(args []string) (string, error) {
		if args[0] == "list" {
			return "[]", nil
		}
		return "", expected
	})
	require.Equal(t, expected, bitwardenSetter("n", "s"))
}

func Test_bitwardenRunReportsMissingCli(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	_, err := bitwardenRun([]string{"list", "items"}, "")
	require.NotNil(t, err)
	require.ErrorContains(t, err, "BW_SESSION")
}
