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

func keychainWriter(serviceLabel string, secret string) {
	item := keychain.NewItem()
	item.SetSecClass(keychain.SecClassGenericPassword)
	item.SetService(serviceLabel)
	item.SetAccount(serviceLabel)
	item.SetData([]byte(secret))
	err := keychain.AddItem(item)

	if err == keychain.ErrorDuplicateItem {
		keychainDeleter(serviceLabel)
		keychainWriter(serviceLabel, secret)
	}
}

func keychainDeleter(serviceLabel string) {
	query := keychain.NewItem()
	query.SetSecClass(keychain.SecClassGenericPassword)
	query.SetService(serviceLabel)
	query.SetMatchLimit(keychain.MatchLimitOne)
	err := keychain.DeleteItem(query)
	if err == keychain.ErrorDuplicateItem {
		panic("errored deleting")
	}
}
