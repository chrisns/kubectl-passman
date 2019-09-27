package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/creasty/defaults"
	"github.com/keybase/go-keychain"
)

func main() {
	var secret string
	if os.Args[1:][0] == "keychain" {
		secret = keychainFetcher(os.Args[1:][1])
	}
	res := &response{}
	json.Unmarshal([]byte(secret), &res.Status)

	fmt.Println(string(formatResponse(res)))
}

func keychainFetcher(serviceLabel string) string {
	query := keychain.NewItem()
	query.SetSecClass(keychain.SecClassGenericPassword)
	query.SetService(serviceLabel)
	query.SetMatchLimit(keychain.MatchLimitOne)
	query.SetReturnData(true)
	results, err := keychain.QueryItem(query)
	if err != nil {
		fmt.Println("err", err)
		os.Exit(1)
	} else if len(results) != 1 {
		fmt.Println("no found")
		os.Exit(1)
	}
	password := string(results[0].Data)
	return password
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

func returnResponse() string {
	res1D := &response{}

	return formatResponse(res1D)
}

func formatResponse(res *response) string {
	defaults.Set(res)
	jsonResponse, _ := json.Marshal(res)
	return string(jsonResponse)
}
