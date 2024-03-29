package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_opgetter_happy(t *testing.T) {
	var expected = "RSA"
	defaultOpGet = func(itemName string) (*opResponse, error) {
		return &opResponse{
			Fields: []opResponseField{{
				Id:    "password",
				Label: "credential",
				Value: expected,
			}},
		}, nil
	}
	actual, err := opgetter("fakecred")
	require.Contains(t, actual, expected)
	require.Nil(t, err)
}

func Test_opgetter_op_fail(t *testing.T) {
	expected := errors.New("test")
	defaultOpGet = func(itemName string) (*opResponse, error) {
		return nil, expected
	}
	actual, err := opgetter("foo")
	require.Equal(t, actual, "")
	require.Equal(t, expected, err)
}

func Test_opgetter_password_not_found(t *testing.T) {
	var expected = "RSA"
	defaultOpGet = func(itemName string) (*opResponse, error) {
		return &opResponse{

			Fields: []opResponseField{{
				Id:    "notpassword",
				Label: "notpassword",
				Value: expected,
			}},
		}, nil
	}
	actual, err := opgetter("test")
	require.Equal(t, err.Error(), "unable to find credential")
	require.Equal(t, actual, "")
}
