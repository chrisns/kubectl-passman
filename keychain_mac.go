// +build darwin,amd64

package main

import (
	"fmt"

	"github.com/keybase/go-keychain"
)

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
