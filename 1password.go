package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
)

type opResponse struct {
	UUID     string            `json:"uuid"`
	Title    string            `json:"title"`
	Category string            `json:"category"`
	Fields   []opResponseField `json:"fields"`
}

type opResponseField struct {
	Id    string `json:"id"`
	Value string `json:"value"`
	Type  string `json:"type"`
	Label string `json:"label"`
}

var defaultOpGet = func(itemName string) (*opResponse, error) {
	cmd := exec.Command("op", "item", "get", itemName, "--format=json", "--debug")
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	var resp opResponse
	err = json.Unmarshal(outb.Bytes(), &resp)
	if err != nil {
		fmt.Println("Error unmarshaling data")
		return nil, err
	}
	return &resp, nil
}

func opgetter(itemName string) (string, error) {
	resp, err := defaultOpGet(itemName)
	if err != nil {
		return "", err
	}
	for _, v := range resp.Fields {
		if v.Label == "credential" {
			return v.Value, nil
		}
	}
	return "", errors.New("unable to find credential")
}

func opsetter(itemName, secret string) error {

	// create a string that will be passed to the op command to store our secret
	credString := "credential[concealed]=" + secret

	stdoutStderr, err := exec.Command("op", "item", "create", "--category",
		"API Credential", "--title", itemName, credString).CombinedOutput()

	if err != nil {
		fmt.Printf("%s\n", stdoutStderr)
	}

	return err
}
