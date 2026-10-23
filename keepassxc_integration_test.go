package main

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// These tests run the real keepassxc-cli and skip when it is not installed.
func keepassxcRealEnv(t *testing.T) {
	t.Helper()
	cli, err := exec.LookPath("keepassxc-cli")
	if err != nil {
		t.Skip("keepassxc-cli is not installed")
	}
	db := filepath.Join(t.TempDir(), "test.kdbx")
	create := exec.Command(cli, "db-create", "--set-password", db)
	create.Stdin = strings.NewReader("dbpass\ndbpass\n")
	out, err := create.CombinedOutput()
	require.NoError(t, err, string(out))

	t.Setenv(keepassxcDatabaseEnv, db)
	t.Setenv(keepassxcKeyFileEnv, "")
	t.Setenv(keepassxcUnlockEnv, "dbpass")
	t.Setenv(keepassxcKeyOnlyEnv, "")
}

func Test_keepassxcRealRoundTrip(t *testing.T) {
	keepassxcRealEnv(t)
	require.NoError(t, keepassxcSetter("kube-prod", jsonToken))
	v, err := keepassxcGetter("kube-prod")
	require.NoError(t, err)
	require.Equal(t, jsonToken, v)
}

func Test_keepassxcRealOverwrite(t *testing.T) {
	keepassxcRealEnv(t)
	require.NoError(t, keepassxcSetter("kube-prod", jsonToken))
	require.NoError(t, keepassxcSetter("kube-prod", `{"token":"11111"}`))
	v, err := keepassxcGetter("kube-prod")
	require.NoError(t, err)
	require.Equal(t, `{"token":"11111"}`, v)
}

func Test_keepassxcRealWrongUnlock(t *testing.T) {
	keepassxcRealEnv(t)
	require.NoError(t, keepassxcSetter("kube-prod", jsonToken))
	t.Setenv(keepassxcUnlockEnv, "wrong")
	_, err := keepassxcGetter("kube-prod")
	require.Error(t, err)
}

func Test_keepassxcRealKeyFileOnly(t *testing.T) {
	keepassxcRealEnv(t)
	dir := t.TempDir()
	db, key := filepath.Join(dir, "k.kdbx"), filepath.Join(dir, "k.key")
	out, err := exec.Command("keepassxc-cli", "db-create", "--set-key-file", key, db).CombinedOutput()
	require.NoError(t, err, string(out))
	t.Setenv(keepassxcDatabaseEnv, db)
	t.Setenv(keepassxcKeyFileEnv, key)
	t.Setenv(keepassxcKeyOnlyEnv, "1")
	require.NoError(t, keepassxcSetter("kf", jsonToken))
	require.NoError(t, keepassxcSetter("kf", `{"token":"22222"}`))
	v, err := keepassxcGetter("kf")
	require.NoError(t, err)
	require.Equal(t, `{"token":"22222"}`, v)
}
