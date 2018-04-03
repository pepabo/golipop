package lolp

import (
	"bytes"
	"encoding/json"
	"log"
)

type ProjectTemplate int

const (
	WORDPRESS ProjectTemplate = iota
	PHP
	RAILS
	NODE
)

func (t ProjectTemplate) String() string {
	switch t {
	case WORDPRESS:
		return "wordpress"
	case PHP:
		return "php"
	case RAILS:
		return "rails"
	case NODE:
		return "node"
	default:
		return "unknown"
	}
}

func GetProjectTemplate(t string) ProjectTemplate {
	switch t {
	case "wordpress":
		return WORDPRESS
	case "php":
		return PHP
	case "rails":
		return RAILS
	case "node":
		return NODE
	default:
		panic(`unknown template: ` + t)
	}
}

type Project struct {
	Type          string                 `json:"type"`
	Domain        string                 `json:"domain"`
	SubDomain     string                 `json:"sub_domain"`
	DBPassword    string                 `json:"db_password,omitempty"`
	CustomDomains []string               `json:"custom_domains,omitempty"`
	Payload       map[string]interface{} `json:"payload,omitempty"`
}

func (c *Client) CreateProject(t ProjectTemplate, p map[string]interface{}) (*Project, error) {
	log.Printf("[INFO] creating project (type: %s)", t)
	body, err := json.Marshal(&Project{
		Type:    t.String(),
		Payload: p,
	})
	if err != nil {
		return nil, err
	}

	request, err := c.Request("POST", "/v1/projects", &RequestOptions{
		Body: bytes.NewReader(body),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	})
	if err != nil {
		return nil, err
	}

	response, err := dispose(c.HTTPClient.Do(request))
	if err != nil {
		return nil, err
	}

	var project Project
	if err := decodeJSON(response, &project); err != nil {
		return nil, err
	}

	return &project, nil
}

func (c *Client) DeleteProject(ID string) error {
	request, err := c.Request("DELETE", `/v1/projects/`+ID, &RequestOptions{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	})
	if err != nil {
		return err
	}

	_, err = dispose(c.HTTPClient.Do(request))
	if err != nil {
		return err
	}

	return nil
}
