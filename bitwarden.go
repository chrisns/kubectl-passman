package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// bitwardenRun runs the Bitwarden CLI with the given arguments and stdin.
// The CLI reads the session key from BW_SESSION.
var bitwardenRun = func(args []string, stdin string) (string, error) {
	cmd := exec.Command("bw", append([]string{"--nointeraction", "--raw"}, args...)...)
	cmd.Stdin = strings.NewReader(stdin)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()
	if err != nil {
		msg := strings.TrimSpace(errb.String())
		if msg == "" {
			return "", fmt.Errorf("%w: check that the vault is unlocked and BW_SESSION is set", err)
		}
		return "", fmt.Errorf("%w: %s", err, msg)
	}
	return outb.String(), nil
}

func bitwardenGetter(itemName string) (string, error) {
	out, err := bitwardenRun([]string{"get", "password", "--", itemName}, "")
	if err != nil {
		return "", err
	}
	secret := strings.TrimRight(out, "\r\n")
	if secret == "" {
		return "", errors.New("the Bitwarden item has no password")
	}
	return secret, nil
}

// bitwardenEncode gives the base64 JSON that the bw create and edit commands read from stdin.
func bitwardenEncode(item map[string]interface{}) (string, error) {
	data, err := json.Marshal(item)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

func bitwardenSetter(itemName, secret string) error {
	out, err := bitwardenRun([]string{"get", "item", "--", itemName}, "")
	if err != nil {
		if !strings.Contains(err.Error(), "Not found") {
			return err
		}
		encoded, err := bitwardenEncode(map[string]interface{}{
			"type":  1,
			"name":  itemName,
			"login": map[string]interface{}{"password": secret},
		})
		if err != nil {
			return err
		}
		_, err = bitwardenRun([]string{"create", "item"}, encoded)
		return err
	}

	item := map[string]interface{}{}
	if err := json.Unmarshal([]byte(out), &item); err != nil {
		return fmt.Errorf("cannot read the Bitwarden item: %w", err)
	}
	id, _ := item["id"].(string)
	if id == "" {
		return errors.New("the Bitwarden item has no id")
	}
	if t, _ := item["type"].(float64); t != 1 {
		return errors.New("the Bitwarden item is not a login item")
	}
	login, _ := item["login"].(map[string]interface{})
	if login == nil {
		login = map[string]interface{}{}
	}
	login["password"] = secret
	item["login"] = login
	encoded, err := bitwardenEncode(item)
	if err != nil {
		return err
	}
	_, err = bitwardenRun([]string{"edit", "item", id}, encoded)
	return err
}
