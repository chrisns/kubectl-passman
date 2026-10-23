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

// bitwardenFind returns the one login item with exactly this name, or nil if there is none.
// The bw get command matches names as a search, so it can return a different item.
func bitwardenFind(itemName string) (map[string]interface{}, error) {
	out, err := bitwardenRun([]string{"list", "items", "--search=" + itemName}, "")
	if err != nil {
		return nil, err
	}
	var items []map[string]interface{}
	if err := json.Unmarshal([]byte(out), &items); err != nil {
		return nil, fmt.Errorf("cannot read the Bitwarden item list: %w", err)
	}
	var found []map[string]interface{}
	for _, item := range items {
		if name, _ := item["name"].(string); name != itemName {
			continue
		}
		if t, _ := item["type"].(float64); t != 1 {
			continue
		}
		found = append(found, item)
	}
	if len(found) > 1 {
		return nil, fmt.Errorf("%d Bitwarden login items are named %q, use a unique name", len(found), itemName)
	}
	if len(found) == 0 {
		return nil, nil
	}
	return found[0], nil
}

func bitwardenGetter(itemName string) (string, error) {
	item, err := bitwardenFind(itemName)
	if err != nil {
		return "", err
	}
	if item == nil {
		return "", fmt.Errorf("no Bitwarden login item is named %q", itemName)
	}
	login, _ := item["login"].(map[string]interface{})
	secret, _ := login["password"].(string)
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
	item, err := bitwardenFind(itemName)
	if err != nil {
		return err
	}
	if item == nil {
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

	id, _ := item["id"].(string)
	if id == "" {
		return errors.New("the Bitwarden item has no id")
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
