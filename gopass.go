package main

import (
	"io"
	"log"
	"os/exec"
)

var defaultGopassGet = func(itemName string) (string, error) {
	out, err := exec.Command("gopass", "show", "--password", itemName).Output()
	return string(out), err
}

func gopassGetter(itemName string) (string, error) {
	return defaultGopassGet(itemName)
}

var gopassWriteCmd = func(itemName string) *exec.Cmd {
	return exec.Command("gopass", "insert", "--force", itemName)
}

var defaultGopassSet = func(itemName, secret string) error {
	var stdin io.WriteCloser
	var err error
	var out []byte

	cmd := gopassWriteCmd(itemName)

	stdin, err = cmd.StdinPipe()

	if err != nil {
		log.Fatal(err)
		return err
	}

	gopassWriteSecret(stdin, secret)

	out, err = cmd.CombinedOutput()

	if err != nil {
		log.Fatal(err)
		return err
	}
	log.Print(string(out))
	return nil
}

var gopassWriteSecret = func(stdin io.WriteCloser, secret string) error {
	defer stdin.Close()
	_, err := io.WriteString(stdin, secret)
	return err
}

func gopassSetter(itemName, secret string) error {
	return defaultGopassSet(itemName, secret)
}
