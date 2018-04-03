package lolp

import (
	"fmt"
	"log"
	"strings"
)

type Login struct {
	Username string
	Password string
}

func (c *Client) Login(username, password string) (string, error) {
	log.Printf("[INFO] logging in user %s", username)

	if len(username) == 0 {
		return "", fmt.Errorf("client: missing username")
	}

	if len(password) == 0 {
		return "", fmt.Errorf("client: missing password")
	}

	jsonStr := fmt.Sprintf(`{"username":"%s","password":"%s"}`, username, password)
	request, err := c.Request("POST", "/v1/authorizations", &RequestOptions{
		Body: strings.NewReader(jsonStr),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	})

	if err != nil {
		return "", err
	}

	response, err := dispose(c.HTTPClient.Do(request))
	if err != nil {
		return "", err
	}

	var token string
	if err := decodeJSON(response, &token); err != nil {
		return "", nil
	}

	c.Token = token

	return c.Token, nil
}
