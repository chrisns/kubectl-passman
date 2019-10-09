package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"

	"github.com/urfave/cli"

	"github.com/creasty/defaults"
)

// VERSION populated at build time
var VERSION = "0.0.0"

var app = cli.NewApp()

func commands() {
	app.Commands = []cli.Command{
		{
			Name:      "keychain",
			Usage:     "Use osx-keychain",
			ArgsUsage: "[item-name]",
			Action: func(c *cli.Context) error {
				return read("keychain", c)
			},
		},
		{
			Name:      "1password",
			Usage:     "Use 1Password",
			Aliases:   []string{"1pass", "op"},
			ArgsUsage: "[item-name]",
			Action: func(c *cli.Context) error {
				return read("1password", c)
			},
		},
	}
}

func info() {
	app.Name = "kubectl-passman"
	app.Usage = "Store kubeconfig credentials in keychains or password managers"
	app.Authors = []cli.Author{
		{
			Name:  "Chris Nesbitt-Smith",
			Email: "chris@cns.me.uk",
		},
	}
	app.Copyright = "(c) 2019 Chris Nesbitt-Smith"
	app.UsageText = `kubectl-passman [command] [item-name]
	 If stdin is provided it will write to the item, otherwise it will read`
	app.Version = VERSION
}

func main() {
	info()
	commands()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func read(handler string, c *cli.Context) error {
	var itemName = c.Args().Get(0)
	if itemName == "" {
		return cli.NewExitError("Please provide [item-name]", 1)
	}
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

var defaultOp = func(itemName string) (*opResponse, error) {
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
	resp, err := defaultOp(itemName)
	if err != nil {
		panic(err)
	}
	i := sort.Search(len(resp.Details.Fields), func(i int) bool { return resp.Details.Fields[i].Name == "password" })
	return resp.Details.Fields[i].Value
}

type opResponse struct {
	UUID    string            `json:"uuid"`
	Details opResponseDetails `json:"details"`
}
type opResponseDetails struct {
	Fields []opResponseField `json:"fields"`
	Title  string            `json:"title"`
}

type opResponseField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
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
