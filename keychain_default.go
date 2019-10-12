// +build !darwin AND !amd64

package main

import (
	"github.com/zalando/go-keyring"
)

func keychainFetcher(serviceLabel string) (string, error) {
	return keyring.Get(serviceLabel, serviceLabel)
}

func keychainWriter(serviceLabel, secret string) error {
	return keyring.Set(serviceLabel, serviceLabel, secret)
}

// func keychainDeleter(serviceLabel string) {
// 	err := keyring.Delete(serviceLabel, serviceLabel)
// 	if err != nil {
// 		panic(err)
// 	}
// }
