package main

import (
	"fmt"
	"io"
	"log"
	"os/exec"
)

var defaultGopassGet = func(itemName string) (string, error) {
	out, err := exec.Command("gopass", "show", "--password", itemName).Output()
	return string(out), err
}

func gopassGetter(itemName string) string {
	resp, err := defaultGopassGet(itemName)
	if err != nil {
		panic(err)
	}
	return resp
}

func gopassSetter(itemName, secret string) {
	cmd := exec.Command("gopass", "insert", "--force", itemName)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		defer stdin.Close()
		_, err := io.WriteString(stdin, secret)
		if err != nil {
			log.Fatal(err)
		}
	}()

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", out)
}
