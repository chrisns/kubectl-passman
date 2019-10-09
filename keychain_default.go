// +build !darwin AND !amd64

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
func keychainWriter(serviceLabel, secret string) {
	err := keyring.Set(serviceLabel, serviceLabel, secret)
	if err != nil {
		panic(err)
	}
}

// func keychainDeleter(serviceLabel string) {
// 	err := keyring.Delete(serviceLabel, serviceLabel)
// 	if err != nil {
// 		panic(err)
// 	}
// }
