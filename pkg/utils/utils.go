package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func ReadToken() (string, error) {

	// read token from /var/run/argocd/token
	filePath := "/var/run/argo/token"

	// Read the token of the file
	token, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading token: %v\n", err)
		return "", err
	}

	// Convert the token to string
	tokenStr := string(token)

	strippedToken := strings.TrimSpace(tokenStr)

	return strippedToken, nil

}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func ParseBody(r *http.Request, x interface{}) {
	if body, err := ioutil.ReadAll(r.Body); err == nil {
		fmt.Println("body: ", string(body))
		if err := json.Unmarshal([]byte(body), x); err != nil {
			return
		}
	}
}
