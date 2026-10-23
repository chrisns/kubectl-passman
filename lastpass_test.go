package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type lastpassCall struct {
	args  []string
	stdin string
}

func fakeLastpass(t *testing.T, fn func(args []string) (string, error)) *[]lastpassCall {
	calls := &[]lastpassCall{}
	orig := lastpassRun
	t.Cleanup(func() { lastpassRun = orig })
	lastpassRun = func(args []string, stdin string) (string, error) {
		*calls = append(*calls, lastpassCall{args, stdin})
		return fn(args)
	}
	return calls
}

func Test_lastpassGetterOK(t *testing.T) {
	calls := fakeLastpass(t, func([]string) (string, error) { return "{\"token\":\"a\"}\n", nil })
	v, err := lastpassGetter("item")
	require.Nil(t, err)
	require.Equal(t, `{"token":"a"}`, v)
	require.Equal(t, []string{"show", "--password", "--", "item"}, (*calls)[0].args)
}

func Test_lastpassGetterERROR(t *testing.T) {
	fakeLastpass(t, func([]string) (string, error) { return "", errors.New("not logged in") })
	_, err := lastpassGetter("item")
	require.EqualError(t, err, "not logged in")
}

func Test_lastpassSetterEdit(t *testing.T) {
	calls := fakeLastpass(t, func([]string) (string, error) { return "1234", nil })
	require.Nil(t, lastpassSetter("item", "secret"))
	require.Len(t, *calls, 2)
	require.Equal(t, []string{"edit", "--non-interactive", "--sync=now", "--password", "--", "item"}, (*calls)[1].args)
	require.Equal(t, "secret\n", (*calls)[1].stdin)
}

func Test_lastpassSetterAdd(t *testing.T) {
	calls := fakeLastpass(t, func(args []string) (string, error) {
		if args[0] == "show" {
			return "", errors.New("Could not find specified account")
		}
		return "", nil
	})
	require.Nil(t, lastpassSetter("item", "secret"))
	require.Len(t, *calls, 2)
	require.Equal(t, []string{"add", "--non-interactive", "--sync=now", "--", "item"}, (*calls)[1].args)
	require.Equal(t, "Password: secret\n", (*calls)[1].stdin)
}

func Test_lastpassSetterERROR(t *testing.T) {
	fakeLastpass(t, func(args []string) (string, error) {
		if args[0] == "show" {
			return "", errors.New(lastpassNotFound)
		}
		return "", errors.New("add failed")
	})
	err := lastpassSetter("item", "secret")
	require.EqualError(t, err, "add failed (looking up the entry failed first: "+lastpassNotFound+")")
}

func Test_lastpassSetterProbeFailureDoesNotAdd(t *testing.T) {
	calls := fakeLastpass(t, func([]string) (string, error) { return "", errors.New("Error: Could not sync") })
	require.EqualError(t, lastpassSetter("item", "secret"), "Error: Could not sync")
	require.Len(t, *calls, 1)
}

func Test_lastpassGetterEmptyOutput(t *testing.T) {
	fakeLastpass(t, func([]string) (string, error) { return "\n", nil })
	_, err := lastpassGetter("item")
	require.ErrorContains(t, err, "no password for item")
}

func Test_lastpassSetterEditERROR(t *testing.T) {
	fakeLastpass(t, func(args []string) (string, error) {
		if args[0] == "edit" {
			return "", errors.New("edit failed")
		}
		return "id", nil
	})
	require.EqualError(t, lastpassSetter("item", "secret"), "edit failed")
}

func Test_lastpassSetterRejectsLineBreak(t *testing.T) {
	calls := fakeLastpass(t, func([]string) (string, error) { return "", nil })
	require.NotNil(t, lastpassSetter("item", "a\nb"))
	require.Len(t, *calls, 0)
}
