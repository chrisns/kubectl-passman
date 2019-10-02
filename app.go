package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"sort"

	"github.com/creasty/defaults"
	"github.com/keybase/go-keychain"
)

func main() {
	var secret string
	if os.Args[1:][0] == "keychain" {
		secret = keychainFetcher(os.Args[1:][1])
	}
	if os.Args[1:][0] == "1password" {
		secret = opgetter(os.Args[1:][1])
	}
	res := &response{}
	json.Unmarshal([]byte(secret), &res.Status)

	fmt.Println(string(formatResponse(res)))
}

func opgetter(itemName string) string {
	out, err := exec.Command("op", "get", "item", itemName).Output()
	if err != nil {
		panic(err)
	}
	dat := opResponse{}
	if err := json.Unmarshal(out, &dat); err != nil {
		panic(err)
	}
	i := sort.Search(len(dat.Details.Fields), func(i int) bool { return dat.Details.Fields[i].Name == "password" })
	return dat.Details.Fields[i].Value

}

var defaultKeychain = func(serviceLabel string) ([]keychain.QueryResult, error) {
	query := keychain.NewItem()
	query.SetSecClass(keychain.SecClassGenericPassword)
	query.SetService(serviceLabel)
	query.SetMatchLimit(keychain.MatchLimitOne)
	query.SetReturnData(true)
	return keychain.QueryItem(query)
}

func keychainFetcher(serviceLabel string) string {
	results, err := defaultKeychain(serviceLabel)
	if err != nil {
		fmt.Println("err", err)
		panic("unable to connect to keychain")
	} else if len(results) != 1 {
		panic("item doesn't exist")
	}
	password := string(results[0].Data)
	return password
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
	defaults.Set(res)
	jsonResponse, _ := json.Marshal(res)
	return string(jsonResponse)
}
