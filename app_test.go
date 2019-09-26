package main

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
)

type MyTestSuite struct {
	suite.Suite
}

var fixture = `{"apiVersion":"client.authentication.k8s.io/v1beta1","kind":"ExecCredential","status":{"token":"my-bearer-token"}}`

func (suite *MyTestSuite) Test_formatResponse() {
	suite.Equal(fixture, formatResponse(&response{}))
}

func (suite *MyTestSuite) Test_formatResponse_is_json() {
	suite.True(json.Valid([]byte(formatResponse(&response{}))))
}

func (suite *MyTestSuite) Test_formatResponse_populate_defaults() {
	suite.Contains(formatResponse(&response{}), "apiVersion")
}
func (suite *MyTestSuite) Test_formatResponse_override_defaults() {
	suite.Contains(formatResponse(&response{Kind: "foo"}), `"kind":"foo"`)
}

func (suite *MyTestSuite) Test_keychainFetcher_NoKeychainError() {
	// panicker := func() {
	// 	// TODO: MOCK keychain.QueryItem(query) returns err=1
	// 	keychainFetcher("error")
	// }
	// suite.PanicsWithValue("unable to connect to keychain", panicker)
}
func (suite *MyTestSuite) Test_keychainFetcher_NoItemFoundError() {
	panicker := func() {
		// TODO: MOCK keychain.QueryItem(query) returns empty array
		keychainFetcher("doesn't exist")
	}
	suite.PanicsWithValue("item doesn't exist", panicker)
}

func (suite *MyTestSuite) Test_keychainFetcher_ItemFound() {
	var expected = "RSA"
	// TODO: MOCK keychain.QueryItem(query) returns "RSARSA"
	suite.Contains(keychainFetcher("gabriel"), expected)
}

func (suite *MyTestSuite) Test_opgetter_happy() {
	var expected = "RSA"
	// TODO: MOCK exec.Command() returns {"details": {"fields":[{"name": "password", "value": "RSARSA"}]}}
	suite.Contains(opgetter("gabriel"), expected)
}

func (suite *MyTestSuite) Test_opgetter_op_fail() {
	// TODO: MOCK exec.Command() returns (ERROR)  item mykubecreds not found
}

func (suite *MyTestSuite) Test_opgetter_password_not_found() {
	// TODO: MOCK exec.Command() returns {"details": {"fields":[{"name": "notpassword", "value": "RSARSA"}]}}
}

func TestMyTestSuite(t *testing.T) {
	suite.Run(t, new(MyTestSuite))
}
