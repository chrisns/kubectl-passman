package main

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// lastpassRun runs lpass with the given arguments and stdin.
var lastpassRun = func(args []string, stdin string) (string, error) {
	cmd := exec.Command("lpass", args...)
	cmd.Stdin = strings.NewReader(stdin)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()
	if err != nil {
		msg := strings.TrimSpace(errb.String())
		if msg == "" {
			return "", fmt.Errorf("%w: check that you are logged in, run lpass login", err)
		}
		return "", fmt.Errorf("%w: %s", err, msg)
	}
	return outb.String(), nil
}

func lastpassGetter(itemName string) (string, error) {
	out, err := lastpassRun([]string{"show", "--password", "--", itemName}, "")
	if err != nil {
		return "", err
	}
	secret := strings.TrimRight(out, "\r\n")
	if secret == "" {
		return "", errors.New("lpass returned no password for " + itemName + ", check that exactly one entry has this name and that its password is set")
	}
	if strings.ContainsAny(secret, "\r\n") {
		return "", errors.New("lpass returned more than one line for " + itemName + ", the entry name must be unique")
	}
	return secret, nil
}

// lastpassNotFound is what lpass prints when no entry has the name.
const lastpassNotFound = "Could not find specified account"

func lastpassSetter(itemName, secret string) error {
	if strings.ContainsAny(secret, "\r\n") {
		return errors.New("secret must not contain a line break")
	}
	// lpass edit fails when the entry is missing, so look it up to choose edit or add.
	_, probeErr := lastpassRun([]string{"show", "--id", "--", itemName}, "")
	// Only add when the entry is missing, any other failure would create a duplicate entry.
	if probeErr != nil && !strings.Contains(probeErr.Error(), lastpassNotFound) {
		return probeErr
	}
	var err error
	if probeErr == nil {
		_, err = lastpassRun([]string{"edit", "--non-interactive", "--sync=now", "--password", "--", itemName}, secret+"\n")
	} else {
		_, err = lastpassRun([]string{"add", "--non-interactive", "--sync=now", "--", itemName}, "Password: "+secret+"\n")
	}
	if err != nil && probeErr != nil {
		return fmt.Errorf("%w (looking up the entry failed first: %v)", err, probeErr)
	}
	return err
}
