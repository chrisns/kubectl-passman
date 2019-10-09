package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"

	"github.com/creasty/defaults"
)

// VERSION populated at build time
var VERSION = "0.0.0"

var app = cli.NewApp()

func cliCommands() {
	app.Commands = []cli.Command{
		{
			Name:      "keychain",
			Usage:     "Use osx-keychain",
			ArgsUsage: "[item-name]",
			Action: func(c *cli.Context) error {
				var itemName = c.Args().Get(0)
				var secret = c.Args().Get(1)
				if itemName == "" {
					return cli.NewExitError("Please provide [item-name]", 1)
				}
				if secret != "" {
					return write("keychain", itemName, secret)
				}
				return read("keychain", itemName)
			},
		},
		{
			Name:      "1password",
			Usage:     "Use 1Password",
			Aliases:   []string{"1pass", "op"},
			ArgsUsage: "[item-name]",
			Action: func(c *cli.Context) error {
				var itemName = c.Args().Get(0)
				var secret = c.Args().Get(1)
				if itemName == "" {
					return cli.NewExitError("Please provide [item-name]", 1)
				}
				if secret != "" {
					return write("1password", itemName, secret)
				}
				return read("1password", itemName)
			},
		},
	}
}

func cliInfo() {
	app.Name = "kubectl-passman"
	app.Usage = "Store kubeconfig credentials in keychains or password managers"
	app.Authors = []cli.Author{
		{
			Name:  "Chris Nesbitt-Smith",
			Email: "chris@cns.me.uk",
		},
	}
	app.Copyright = "(c) 2019 Chris Nesbitt-Smith"
	app.UsageText = `kubectl-passman [command] [item-name] [new-value(optional)]
	 If new-value is provided it will write to the item, otherwise it will read`
	app.Version = VERSION
}

func main() {
	cliInfo()
	cliCommands()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func write(handler string, itemName string, secret string) error {
	if handler == "keychain" {
		keychainWriter(itemName, secret)
	}
	if handler == "1password" {
		opsetter(itemName, secret)
	}
	return nil
}

func read(handler string, itemName string) error {
	var secret string
	if handler == "keychain" {
		secret = keychainFetcher(itemName)
	}
	if handler == "1password" {
		secret = opgetter(itemName)
	}
	res := &response{}
	err := json.Unmarshal([]byte(secret), &res.Status)
	if err != nil {
		panic(err)
	}
	fmt.Println(formatResponse(res))
	return nil
}

type responseStatus struct {
	Token                 string `default:"my-bearer-token" json:"token,omitempty"`
	ClientCertificateData string `json:"clientCertificateData,omitempty"`
	ClientKeyData         string `json:"clientKeyData,omitempty"`
}
type response struct {
	APIVersion string         `default:"client.authentication.k8s.io/v1beta1" json:"apiVersion"`
	Kind       string         `default:"ExecCredential" json:"kind"`
	Status     responseStatus `json:"status"`
}

func formatResponse(res *response) string {
	err := defaults.Set(res)
	if err != nil {
		panic(err)
	}
	jsonResponse, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}
	return string(jsonResponse)
}
