package lolp

import (
	"fmt"
	"log"
	"strings"
)

// Login for authorization
func (c *Client) Login(u string, p string) (string, error) {
	if len(u) == 0 {
		return "", fmt.Errorf("client: missing username")
	}

	if len(p) == 0 {
		return "", fmt.Errorf("client: missing password")
	}

	json := fmt.Sprintf(`{"username":"%s","password":"%s"}`, u, p)
	res, err := c.HTTP("POST", "/v1/authorizations", &RequestOptions{
		Body: strings.NewReader(json),
	})
	if err != nil {
		return "", err
	}

	var t string
	if err := decodeJSON(response, &t); err != nil {
		return "", err
	}

	log.Printf("[DEBUG] setting token (%s)", mask(t))
	c.Token = t

	return t, nil
}

// mask for string
func mask(s string) string {
	if len(s) <= 3 {
		return "***[masked]"
	}

	return s[0:3] + "***[masked]"
}
