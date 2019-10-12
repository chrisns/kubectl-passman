package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_gopassGetterOK(t *testing.T) {
	defaultGopassGet = func(itemName string) (string, error) {
		return "foo", nil
	}
	v, e := gopassGetter("something")
	require.Equal(t, "foo", v)
	require.Nil(t, e)
}

func Test_gopassGetterERROR(t *testing.T) {
	defaultGopassGet = func(itemName string) (string, error) {
		return "bar", errors.New("foobar")
	}
	v, e := gopassGetter("something")
	require.Equal(t, "foobar", e.Error())
	require.Equal(t, "bar", v)
}

func Test_gopassSetterOK(t *testing.T) {
	defaultGopassSet = func(itemName, secret string) error {
		return nil
	}
	e := gopassSetter("foo", "bar")
	require.Nil(t, e)
}

func Test_gopassSetterERROR(t *testing.T) {
	defaultGopassSet = func(itemName, secret string) error {
		return errors.New("foofoo")
	}
	e := gopassSetter("foo", "bar")
	require.Equal(t, "foofoo", e.Error())
}
