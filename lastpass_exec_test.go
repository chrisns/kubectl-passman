package main

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

// fakeLpass puts a shell script called lpass first in PATH.
func fakeLpass(t *testing.T, script string) {
	if runtime.GOOS == "windows" {
		t.Skip("needs a shell")
	}
	dir := t.TempDir()
	require.Nil(t, os.WriteFile(filepath.Join(dir, "lpass"), []byte("#!/bin/sh\n"+script), 0o755))
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
}

func Test_lastpassRunOutput(t *testing.T) {
	fakeLpass(t, "cat\n")
	out, err := lastpassRun([]string{"x"}, "hello\n")
	require.Nil(t, err)
	require.Equal(t, "hello\n", out)
}

func Test_lastpassRunStderrInError(t *testing.T) {
	fakeLpass(t, "echo 'Error: Could not find specified account' >&2; exit 1\n")
	_, err := lastpassRun([]string{"x"}, "")
	require.ErrorContains(t, err, "Could not find specified account")
}

func Test_lastpassRunHintWithoutStderr(t *testing.T) {
	fakeLpass(t, "exit 1\n")
	_, err := lastpassRun([]string{"x"}, "")
	require.ErrorContains(t, err, "lpass login")
}

func Test_lastpassCommandRegistered(t *testing.T) {
	cliCommands()
	found := false
	for _, c := range app.Commands {
		if c.Name == "lastpass" {
			found = true
			require.Contains(t, c.Aliases, "lpass")
		}
	}
	require.True(t, found)
}
