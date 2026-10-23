package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/term"
)

const (
	keepassxcDatabaseEnv   = "KUBECTL_PASSMAN_KEEPASSXC_DB"
	keepassxcKeyFileEnv    = "KUBECTL_PASSMAN_KEEPASSXC_KEYFILE"
	keepassxcPasswordEnv   = "KUBECTL_PASSMAN_KEEPASSXC_PASSWORD"
	keepassxcNoPasswordEnv = "KUBECTL_PASSMAN_KEEPASSXC_NO_PASSWORD"
)

// keepassxcRun runs keepassxc-cli with the given arguments and stdin.
var keepassxcRun = func(args []string, stdin string) (string, error) {
	cmd := exec.Command("keepassxc-cli", args...)
	cmd.Stdin = strings.NewReader(stdin)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()
	if err != nil {
		msg := strings.TrimSpace(errb.String())
		if msg == "" {
			return "", err
		}
		return "", fmt.Errorf("%w: %s", err, msg)
	}
	return outb.String(), nil
}

// keepassxcPrompt reads the database password from the terminal without echo.
var keepassxcPrompt = func() (string, error) {
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		return "", errors.New("no terminal to prompt for the database password, set " + keepassxcPasswordEnv)
	}
	defer tty.Close()
	fmt.Fprint(tty, "Enter password to unlock KeePassXC database: ")
	pw, err := term.ReadPassword(int(tty.Fd()))
	fmt.Fprintln(tty)
	return string(pw), err
}

type keepassxcConfig struct {
	database   string
	keyFile    string
	noPassword bool
	password   string
}

func keepassxcLoadConfig() (*keepassxcConfig, error) {
	cfg := &keepassxcConfig{
		database:   os.Getenv(keepassxcDatabaseEnv),
		keyFile:    os.Getenv(keepassxcKeyFileEnv),
		noPassword: os.Getenv(keepassxcNoPasswordEnv) != "",
		password:   os.Getenv(keepassxcPasswordEnv),
	}
	if cfg.database == "" {
		return nil, errors.New("please set " + keepassxcDatabaseEnv + " to the path of the KeePassXC database")
	}
	if cfg.noPassword && cfg.keyFile == "" {
		return nil, errors.New(keepassxcNoPasswordEnv + " needs " + keepassxcKeyFileEnv)
	}
	if !cfg.noPassword && cfg.password == "" {
		pw, err := keepassxcPrompt()
		if err != nil {
			return nil, err
		}
		cfg.password = pw
	}
	return cfg, nil
}

// args builds the keepassxc-cli arguments, the database and entry come last.
func (c *keepassxcConfig) args(command, itemName string, extra ...string) []string {
	args := []string{command, "--quiet"}
	if c.noPassword {
		args = append(args, "--no-password")
	}
	if c.keyFile != "" {
		args = append(args, "--key-file", c.keyFile)
	}
	args = append(args, extra...)
	return append(args, "--", c.database, itemName)
}

// stdin gives the input that keepassxc-cli reads after its own arguments.
func (c *keepassxcConfig) stdin(lines ...string) string {
	var in []string
	if !c.noPassword {
		in = append(in, c.password)
	}
	in = append(in, lines...)
	return strings.Join(in, "\n") + "\n"
}

func keepassxcGetter(itemName string) (string, error) {
	cfg, err := keepassxcLoadConfig()
	if err != nil {
		return "", err
	}
	out, err := keepassxcRun(cfg.args("show", itemName, "--attributes", "Password"), cfg.stdin())
	if err != nil {
		return "", err
	}
	return strings.TrimRight(out, "\r\n"), nil
}

func keepassxcSetter(itemName, secret string) error {
	if strings.ContainsAny(secret, "\r\n") {
		return errors.New("secret must not contain a line break")
	}
	cfg, err := keepassxcLoadConfig()
	if err != nil {
		return err
	}
	command := "add"
	if _, err := keepassxcRun(cfg.args("show", itemName, "--attributes", "Title"), cfg.stdin()); err == nil {
		command = "edit"
	}
	_, err = keepassxcRun(cfg.args(command, itemName, "--password-prompt"), cfg.stdin(secret))
	return err
}
