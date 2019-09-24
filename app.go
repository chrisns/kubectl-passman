package main

import (
	"encoding/json"
	"fmt"

	"github.com/creasty/defaults"
)

func main() {

	fmt.Println(returnResponse())
}

type responseStatus struct {
	Token string `default:"my-bearer-token" json:"token"`
}
type response struct {
	APIVersion string         `default:"client.authentication.k8s.io/v1beta1" json:"apiVersion"`
	Kind       string         `default:"ExecCredential" json:"kind"`
	Status     responseStatus `json:"status"`
}

func returnResponse() string {
	res1D := &response{}

	defaults.Set(res1D)
	res1B, _ := json.Marshal(res1D)
	return string(res1B)
}
