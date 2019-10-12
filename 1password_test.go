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
			Details: opResponseDetails{
				Fields: []opResponseField{{
					Name:        "password",
					Designation: "password",
					Value:       expected,
				}},
			},
		}, nil
	}
	actual, _ := opgetter("gabriel")
	require.Contains(t, actual, expected)
}

func Test_opgetter_op_fail(t *testing.T) {
	expected := errors.New("test")
	defaultOpGet = func(itemName string) (*opResponse, error) {
		return nil, expected
	}
	_, actual := opgetter("foo")
	require.Equal(t, expected, actual)
}

func Test_opgetter_password_not_found(t *testing.T) {
	var expected = "RSA"
	defaultOpGet = func(itemName string) (*opResponse, error) {
		return &opResponse{
			Details: opResponseDetails{
				Fields: []opResponseField{{
					Name:        "notpassword",
					Designation: "notpassword",
					Value:       expected,
				}},
			},
		}, nil
	}
	_, err := opgetter("test")
	require.Equal(t, err.Error(), "Unable to find password")
}
