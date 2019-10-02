// +build !darwin

package main

import (
	"github.com/zalando/go-keyring"
)

func keychainFetcher(serviceLabel string) string {
	secret, err := keyring.Get(serviceLabel, serviceLabel)
	if err != nil {
		panic(err)
	}
	return secret
}
