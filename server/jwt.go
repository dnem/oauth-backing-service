package server

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

func parseToken(token string) (t *jwt.Token, err error) {
	tokenKey, err := getTokenKey()
	if err != nil {
		return nil, err
	}

	keyFunc := func(t *jwt.Token) (interface{}, error) {
		return []byte(tokenKey), nil
	}

	t, err = jwt.Parse(token, keyFunc)
	if err != nil {
		return nil, err
	}

	if !t.Valid {
		err = errors.New("Token is not valid")
		return nil, err
	}

	return t, nil
}

func hasScope(token *jwt.Token, desiredScopes ...string) bool {
	scopeFound := false

	scopes := token.Claims["scope"]
	a := scopes.([]interface{})

	for _, scope := range a {
		for _, desiredScope := range desiredScopes {
			if scope.(string) == desiredScope {
				scopeFound = true
				break
			}
		}
	}
	return scopeFound
}

type keyObject struct {
	// Alg is the encryption algorithm
	Alg string `json:"alg"`
	// Value is the actual pem-encoded key used to parse JWT tokens
	Value string `json:"value"`
	// Kty
	Kty string `json:"kty,omitempty"`
	// Use
	Use string `json:"use,omitempty"`
	// N
	N string `json:"n,omitempty"`
	// E
	E string `json:"e,omitempty"`
}

func getTokenKey() (key string, err error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get("https://sso-internal.login.system.pcf.local/token_key")
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}

	payload, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ERROR RETRIEVING TOKEN_KEY!")
		return "", err
	}

	ko := &keyObject{}
	err = json.Unmarshal(payload, ko)
	if err != nil {
		fmt.Println("ERROR PARSING TOKEN_KEY!")
		return "", err
	}

	fmt.Printf("RETRIEVED KEY: %s\n", ko.Value)
	if len(ko.Value) == 0 {
		fmt.Println("RETRIEVED TOKEN KEY IS EMPTY!")
		return "", err
	}

	return ko.Value, nil
}
