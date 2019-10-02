// +build darwin,amd64

package main

import (
	"errors"
	"testing"

	"github.com/keybase/go-keychain"
	"github.com/stretchr/testify/require"
)

func Test_keychainFetcher_NoKeychainError(t *testing.T) {
	panicker := func() {
		defaultKeychain = func(serviceLabel string) ([]keychain.QueryResult, error) {
			return nil, errors.New("error")
		}
		keychainFetcher("error")
	}
	require.PanicsWithValue(t, "unable to connect to keychain", panicker)
}
func Test_keychainFetcher_NoItemFoundError(t *testing.T) {
	panicker := func() {
		// TODO: MOCK keychain.QueryItem(query) returns empty array
		defaultKeychain = func(serviceLabel string) ([]keychain.QueryResult, error) {
			return nil, nil
		}
		keychainFetcher("doesn't exist")
	}
	require.PanicsWithValue(t, "item doesn't exist", panicker)
}

func Test_keychainFetcher_ItemFound(t *testing.T) {
	var expected = "RSA"
	defaultKeychain = func(serviceLabel string) ([]keychain.QueryResult, error) {
		return []keychain.QueryResult{{Data: []byte(expected)}}, nil
	}
	require.Contains(t, expected, keychainFetcher("gabriel"))
}
