// +build darwin,amd64

package main

import (
	"errors"
	"log"

	"github.com/keybase/go-keychain"
)

var defaultKeychain = func(serviceLabel string) ([]keychain.QueryResult, error) {
	query := getKeychain(serviceLabel)
	query.SetReturnData(true)
	return keychain.QueryItem(query)
}

func getKeychain(serviceLabel string) keychain.Item {
	query := keychain.NewItem()
	query.SetSecClass(keychain.SecClassGenericPassword)
	query.SetService(serviceLabel)
	return query
}

func keychainFetcher(serviceLabel string) (string, error) {
	results, err := defaultKeychain(serviceLabel)
	if err != nil {
		return "", errors.New("unable to connect to keychain")
	} else if len(results) > 1 {
		return "", errors.New("too many matching items")
	} else if len(results) != 1 {
		return "", errors.New("item doesn't exist")
	}
	return string(results[0].Data), nil
}

func keychainWriter(serviceLabel string, secret string) error {
	query := getKeychain(serviceLabel)
	query.SetData([]byte(secret))
	err := keychain.AddItem(query)

	if err == keychain.ErrorDuplicateItem {
		log.Print("Item already exists, deleting it first")
		err = keychainDeleter(serviceLabel)
		if err != nil {
			return err
		}
		return keychainWriter(serviceLabel, secret)
	}
	return err
}

func keychainDeleter(serviceLabel string) error {
	query := getKeychain(serviceLabel)
	query.SetMatchLimit(keychain.MatchLimitOne)
	return keychain.DeleteItem(query)
}
