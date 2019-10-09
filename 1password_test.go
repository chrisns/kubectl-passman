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
	require.Contains(t, opgetter("gabriel"), expected)
}

func Test_opgetter_op_fail(t *testing.T) {
	err := errors.New("test")
	defaultOpGet = func(itemName string) (*opResponse, error) {
		return nil, err
	}
	require.PanicsWithValue(t, err, func() { opgetter("mykubecreds") })
}

func Test_opgetter_password_not_found(t *testing.T) {
	var expected = "RSA"
	defaultOpGet = func(itemName string) (*opResponse, error) {
		return &opResponse{
			Details: opResponseDetails{
				Fields: []opResponseField{{
					Name:  "notpassword",
					Value: expected,
				}},
			},
		}, nil
	}
	require.Panics(t, func() { opgetter("test") }) // TODO: panics with index out of range; is this expected behavior?
}
