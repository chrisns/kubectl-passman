package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type keepassxcCall struct {
	args  []string
	stdin string
}

func keepassxcFake(t *testing.T, calls *[]keepassxcCall, fn func(args []string) (string, error)) {
	orig := keepassxcRun
	t.Cleanup(func() { keepassxcRun = orig })
	keepassxcRun = func(args []string, stdin string) (string, error) {
		*calls = append(*calls, keepassxcCall{args, stdin})
		return fn(args)
	}
}

func keepassxcEnv(t *testing.T) {
	t.Setenv(keepassxcDatabaseEnv, "/db.kdbx")
	t.Setenv(keepassxcKeyFileEnv, "")
	t.Setenv(keepassxcPasswordEnv, "dbpass")
	t.Setenv(keepassxcNoPasswordEnv, "")
}

func Test_keepassxcGetterOK(t *testing.T) {
	keepassxcEnv(t)
	var calls []keepassxcCall
	keepassxcFake(t, &calls, func([]string) (string, error) { return jsonToken + "\n", nil })
	v, err := keepassxcGetter("kube/prod")
	require.Nil(t, err)
	require.Equal(t, jsonToken, v)
	require.Equal(t, []keepassxcCall{{
		args:  []string{"show", "--quiet", "--attributes", "Password", "--", "/db.kdbx", "kube/prod"},
		stdin: "dbpass\n",
	}}, calls)
}

func Test_keepassxcGetterKeyFileNoPassword(t *testing.T) {
	keepassxcEnv(t)
	t.Setenv(keepassxcPasswordEnv, "")
	t.Setenv(keepassxcNoPasswordEnv, "1")
	t.Setenv(keepassxcKeyFileEnv, "/db.key")
	var calls []keepassxcCall
	keepassxcFake(t, &calls, func([]string) (string, error) { return "x", nil })
	_, err := keepassxcGetter("item")
	require.Nil(t, err)
	require.Equal(t, []string{"show", "--quiet", "--no-password", "--key-file", "/db.key", "--attributes", "Password", "--", "/db.kdbx", "item"}, calls[0].args)
	require.Equal(t, "\n", calls[0].stdin)
}

func Test_keepassxcGetterNoDatabase(t *testing.T) {
	keepassxcEnv(t)
	t.Setenv(keepassxcDatabaseEnv, "")
	_, err := keepassxcGetter("item")
	require.ErrorContains(t, err, keepassxcDatabaseEnv)
}

func Test_keepassxcNoPasswordNeedsKeyFile(t *testing.T) {
	keepassxcEnv(t)
	t.Setenv(keepassxcNoPasswordEnv, "1")
	_, err := keepassxcGetter("item")
	require.ErrorContains(t, err, keepassxcKeyFileEnv)
}

func Test_keepassxcGetterPrompt(t *testing.T) {
	keepassxcEnv(t)
	t.Setenv(keepassxcPasswordEnv, "")
	orig := keepassxcPrompt
	t.Cleanup(func() { keepassxcPrompt = orig })
	keepassxcPrompt = func() (string, error) { return "typed", nil }
	var calls []keepassxcCall
	keepassxcFake(t, &calls, func([]string) (string, error) { return "x", nil })
	_, err := keepassxcGetter("item")
	require.Nil(t, err)
	require.Equal(t, "typed\n", calls[0].stdin)
}

func Test_keepassxcGetterERROR(t *testing.T) {
	keepassxcEnv(t)
	var calls []keepassxcCall
	keepassxcFake(t, &calls, func([]string) (string, error) { return "", errors.New("boom") })
	_, err := keepassxcGetter("item")
	require.Equal(t, "boom", err.Error())
}

func Test_keepassxcSetterAdd(t *testing.T) {
	keepassxcEnv(t)
	var calls []keepassxcCall
	keepassxcFake(t, &calls, func(args []string) (string, error) {
		if args[0] == "show" {
			return "", errors.New("could not find entry")
		}
		return "", nil
	})
	require.Nil(t, keepassxcSetter("item", jsonToken))
	require.Len(t, calls, 2)
	require.Equal(t, []string{"add", "--quiet", "--password-prompt", "--", "/db.kdbx", "item"}, calls[1].args)
	require.Equal(t, "dbpass\n"+jsonToken+"\n", calls[1].stdin)
}

func Test_keepassxcSetterEdit(t *testing.T) {
	keepassxcEnv(t)
	var calls []keepassxcCall
	keepassxcFake(t, &calls, func([]string) (string, error) { return "item", nil })
	require.Nil(t, keepassxcSetter("item", jsonToken))
	require.Equal(t, "edit", calls[1].args[0])
}

func Test_keepassxcSetterERROR(t *testing.T) {
	keepassxcEnv(t)
	var calls []keepassxcCall
	keepassxcFake(t, &calls, func(args []string) (string, error) {
		if args[0] == "show" {
			return "", errors.New("missing")
		}
		return "", errors.New("add failed")
	})
	require.Equal(t, "add failed", keepassxcSetter("item", jsonToken).Error())
}

func Test_keepassxcSetterRejectsLineBreak(t *testing.T) {
	keepassxcEnv(t)
	require.ErrorContains(t, keepassxcSetter("item", "a\nb"), "line break")
}
