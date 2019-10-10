package main

import (
	"fmt"
	"io"
	"log"
	"os/exec"
)

var defaultGopassGet = func(itemName string) (error, string) {
	out, err := exec.Command("gopass", "show", "--password", itemName).Output()
	return err, string(out)
}

func gopassGetter(itemName string) string {
	err, resp := defaultGopassGet(itemName)
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
		err := io.WriteString(stdin, secret)
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
