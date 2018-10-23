package lolp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
)

type PublicKey struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

// AddPublicKey add OpenSSH public key
func (c *Client) AddPublicKey(p *PublicKey) (*PublicKey, error) {
	if len(p.Name) == 0 {
		return nil, fmt.Errorf("client: missing name")
	}
	if len(p.Key) == 0 {
		return nil, fmt.Errorf("client: missing key")
	}

	body, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] request body: %s", body)

	res, err := c.HTTP("POST", "/v1/pubkeys", &RequestOptions{
		Body: bytes.NewReader(body),
	})
	if err != nil {
		return nil, err
	}

	var pubKey PublicKey
	if err := decodeJSON(res, &pubKey); err != nil {
		return nil, err
	}

	return &pubKey, nil
}

// DeletePublicKey delete OpenSSH public key
func (c *Client) DeletePublicKey(name string) error {
	if len(name) == 0 {
		return fmt.Errorf("client: missing name")
	}

	_, err := c.HTTP("DELETE", "/v1/pubkeys/"+name, nil)
	if err != nil {
		return err
	}

	return nil
}
