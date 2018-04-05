package lolp

import (
	"fmt"
	"log"
	"strings"
)

// Login for authorization
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

	log.Printf("[DEBUG] setting token (%s)", mask(token))
	c.Token = token

	return c.Token, nil
}

// mask for string
func mask(s string) string {
	if len(s) <= 3 {
		return "***[masked]"
	}

	return s[0:3] + "***[masked]"
}
