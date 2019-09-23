package main

import "fmt"

func main() {
	fmt.Println(returnResponse())
}

func returnResponse() string {
	return `{"apiVersion": "client.authentication.k8s.io/v1beta1","kind": "ExecCredential","status": {"token": "my-bearer-token"}}`
}
