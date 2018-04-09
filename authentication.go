package lolp

import (
	"fmt"
	"strings"
)

// Authenticate for authorization
func (c *Client) Authenticate(u string, p string) (string, error) {
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
	if err := decodeJSON(res, &t); err != nil {
		return "", err
	}
	c.Token = t

	return t, nil
}
