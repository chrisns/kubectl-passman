package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"sort"
)

type opResponse struct {
	UUID    string            `json:"uuid"`
	Details opResponseDetails `json:"details"`
}
type opResponseDetails struct {
	Fields []opResponseField `json:"fields"`
	Title  string            `json:"title"`
}

type opResponseField struct {
	Name        string `json:"name"`
	Value       string `json:"value"`
	Type        string `json:"type"`
	Designation string `json:"designation"`
}

var defaultOpGet = func(itemName string) (*opResponse, error) {
	out, err := exec.Command("op", "get", "item", itemName).Output()
	if err != nil {
		return nil, err
	}
	var resp opResponse
	err = json.Unmarshal(out, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func opgetter(itemName string) string {
	resp, err := defaultOpGet(itemName)
	if err != nil {
		panic(err)
	}
	i := sort.Search(len(resp.Details.Fields), func(i int) bool {
		return resp.Details.Fields[i].Designation == "password"
	})
	return resp.Details.Fields[i].Value
}

func opsetter(itemName, secret string) {
	var res = &opResponseDetails{
		Fields: []opResponseField{
			{
				Name:        "password",
				Designation: "password",
				Type:        "P",
				Value:       secret,
			},
		},
	}

	jsonResponse, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}

	stdoutStderr, err := exec.Command("op", "create", "item", "login",
		base64.StdEncoding.EncodeToString(jsonResponse), "--title="+itemName).CombinedOutput()

	fmt.Printf("%s\n", stdoutStderr)

	if err != nil {
		log.Fatal(err)
	}
}
