package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
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
			Usage:     "Use keychain/keyring",
			Aliases:   []string{"keyring"},
			ArgsUsage: "[item-name]",
			Action:    cliHandler,
		},
		{
			Name:      "1password",
			Usage:     "Use 1Password",
			Aliases:   []string{"1pass", "op"},
			ArgsUsage: "[item-name]",
			Action:    cliHandler,
		},
		{
			Name:      "gopass",
			Usage:     "Use gopass",
			ArgsUsage: "[item-name]",
			Action:    cliHandler,
		},
	}
}

func cliHandler(c *cli.Context) error {
	var handler string = c.Command.Name
	var itemName = c.Args().Get(0)
	var secret = c.Args().Get(1)
	if itemName == "" {
		return cli.NewExitError("Please provide [item-name]", 1)
	}
	if secret != "" {
		return write(handler, itemName, secret)
	}
	return read(handler, itemName)
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

func write(handler, itemName, secret string) error {

	validSecret, err := formatValidator(secret)
	if err != nil {
		return err
	}

	switch handler {
	case "keychain":
		return keychainWriter(itemName, validSecret)
	case "1password":
		return opsetter(itemName, validSecret)
	case "gopass":
		return gopassSetter(itemName, validSecret)
	}
	return nil
}

func read(handler, itemName string) error {
	var secret string
	var err error
	var out string

	switch handler {
	case "keychain":
		secret, err = keychainFetcher(itemName)
	case "1password":
		secret, err = opgetter(itemName)
	case "gopass":
		secret, err = gopassGetter(itemName)
	}

	if err != nil {
		return err
	}
	res := &response{}
	err = json.Unmarshal([]byte(secret), &res.Status)
	if err != nil {
		return err
	}
	out, err = formatResponse(res)
	fmt.Println(out)
	return err
}

type responseStatus struct {
	Token                  string `default:"my-bearer-token" json:"token,omitempty"`
	ClientCertificateData  string `json:"clientCertificateData,omitempty"`
	ClientCertificateDataD string `json:"client-certificate-data,omitempty"`
	ClientKeyData          string `json:"clientKeyData,omitempty"`
	ClientKeyDataD         string `json:"client-key-data,omitempty"`
}
type response struct {
	APIVersion string         `default:"client.authentication.k8s.io/v1beta1" json:"apiVersion"`
	Kind       string         `default:"ExecCredential" json:"kind"`
	Status     responseStatus `json:"status"`
}

func formatValidator(secret string) (string, error) {
	s := &responseStatus{}
	data := []byte(secret)

	err := json.Unmarshal(data, s)
	if err != nil {
		return "", err
	}

	if len(s.ClientCertificateDataD) > 0 && len(s.ClientKeyDataD) > 0 {
		dataCrt, errCrt := base64.StdEncoding.DecodeString(s.ClientCertificateDataD)
		dataKey, errKey := base64.StdEncoding.DecodeString(s.ClientKeyDataD)
		if errCrt == nil && errKey == nil {
			s.ClientCertificateData = string(dataCrt)
			s.ClientKeyData = string(dataKey)
		} else if errCrt != nil {
			return "", errCrt
		} else if errKey != nil {
			return "", errKey
		}
		s.ClientCertificateDataD = ""
		s.ClientKeyDataD = ""
		s.Token = ""
	} else if len(s.ClientCertificateData) > 0 && len(s.ClientKeyData) > 0 {
		s.ClientCertificateDataD = ""
		s.ClientKeyDataD = ""
		s.Token = ""
	} else if len(s.Token) > 0 {
		s.ClientCertificateDataD = ""
		s.ClientKeyDataD = ""
		s.ClientCertificateData = ""
		s.ClientKeyData = ""
	} else {
		return "", errors.New("Cannot define valid secret format")
	}

	secretByte, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	return string(secretByte), nil
}

func formatResponse(res *response) (string, error) {
	err := defaults.Set(res)
	if err != nil {
		return "", err
	}
	jsonResponse, err := json.Marshal(res)
	if err != nil {
		return "", err
	}
	return string(jsonResponse), nil
}
