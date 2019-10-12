// +build darwin,amd64

package main

import (
	"errors"
	"testing"

	"github.com/keybase/go-keychain"
	"github.com/stretchr/testify/require"
)

func Test_keychainFetcher_NoKeychainError(t *testing.T) {
	defaultKeychain = func(serviceLabel string) ([]keychain.QueryResult, error) {
		return nil, errors.New("error")
	}
	_, err := keychainFetcher("ff")
	require.Equal(t, "unable to connect to keychain", err.Error())
}

func Test_keychainFetcher_NoItemFoundError(t *testing.T) {
	defaultKeychain = func(serviceLabel string) ([]keychain.QueryResult, error) {
		return nil, nil
	}
	_, err := keychainFetcher("ff")
	require.Equal(t, "item doesn't exist", err.Error())
}

func Test_keychainFetcher_ToManyMatching(t *testing.T) {
	defaultKeychain = func(serviceLabel string) ([]keychain.QueryResult, error) {
		return []keychain.QueryResult{{Data: []byte("foo")}, {Data: []byte("foo")}}, nil
	}
	_, err := keychainFetcher("ff")
	require.Equal(t, "too many matching items", err.Error())
}

func Test_keychainFetcher_ItemFound(t *testing.T) {
	var expected = "RSA"
	defaultKeychain = func(serviceLabel string) ([]keychain.QueryResult, error) {
		return []keychain.QueryResult{{Data: []byte(expected)}}, nil
	}
	actual, err := keychainFetcher("ff")

	require.Contains(t, expected, actual)
	require.Nil(t, err)
}
