package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
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

	var res = &opResponse{
		Title:    itemName,
		Category: "API_CREDENTIAL",
		Fields: []opResponseField{
			{
				Id:    "credential",
				Label: "credential",
				Type:  "CONCEALED",
				Value: secret,
			},
			{
				Id:    "type",
				Label: "type",
				Type:  "MENU",
				Value: "json",
			},
		},
	}

	jsonResponse, err := json.Marshal(res)

	if err != nil {
		return err
	}

	// the best way to write a new value to 1password seems to be using a json file and crating an entry from that
	// we will create a temp file, and then delete it when this function returns to be as secure as possible

	f, err := os.CreateTemp("", "sample")
	if err != nil {
		panic(err)
	}
	fmt.Println("Temp file name:", f.Name())
	defer os.Remove(f.Name())
	_, err = f.Write(jsonResponse)
	if err != nil {
		panic(err)
	}
	f.Close()

	stdoutStderr, err := exec.Command("op", "item", "create", "--template", f.Name()).CombinedOutput()

	if err != nil {
		fmt.Printf("%s\n", stdoutStderr)
	}

	return err
}
