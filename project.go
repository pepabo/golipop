package lolp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// Project struct
type Project struct {
	ID            string    `json:"id,omitempty"`
	Name          string    `json:"name,omitempty"`
	Kind          string    `json:"kind,omitempty"`
	Domain        string    `json:"domain,omitempty"`
	SubDomain     string    `json:"subDomain,omitempty"`
	CustomDomains []string  `json:"customDomains,omitempty"`
	Database      Database  `json:"database,omitempty"`
	SSH           *SSH      `json:"ssh,omitempty"`
	CreatedAt     time.Time `json:"createdAt,omitempty"`
	UpdatedAt     time.Time `json:"updatedAt,omitempty"`
}

type Database struct {
	Host string `json:"host,omitempty"`
	User string `json:"user,omitempty"`
	Name string `json:"name,omitempty"`
}

type SSH struct {
	User string `json:"user,omitempty"`
	Host string `json:"host,omitempty"`
	Port int    `json:"port,omitempty"`
}

// ProjectNew struct on create
type ProjectNew struct {
	Name          string                 `json:"name,omitempty"`
	Kind          string                 `json:"kind,omitempty""`
	SubDomain     string                 `json:"sub_domain,omitempty"`
	CustomDomains []string               `json:"custom_domains,omitempty"`
	Payload       map[string]interface{} `json:"payload,omitempty"`
	DBPassword    string                 `json:"db_password,omitempty"`
}

type ProjectCreateResponse struct {
	ID     string `json:"id"`
	Domain string `json:"domain"`
}

// Projects returns project list
func (c *Client) Projects() (*[]Project, error) {
	res, err := c.HTTP("GET", "/v1/projects", nil)
	if err != nil {
		return nil, err
	}

	var ps []Project
	if err := decodeJSON(res, &ps); err != nil {
		return nil, err
	}

	return &ps, nil
}

// Project returns a project by sub-domain name
func (c *Client) Project(name string) (*Project, error) {
	res, err := c.HTTP("GET", `/v1/projects/`+name, nil)
	if err != nil {
		return nil, err
	}

	var p Project
	if err := decodeJSON(res, &p); err != nil {
		return nil, err
	}

	return &p, nil
}

// CreateProject creates project with kind
func (c *Client) CreateProject(p *ProjectNew) (*ProjectCreateResponse, error) {
	if len(p.Kind) == 0 {
		return nil, fmt.Errorf("client: missing kind")
	}

	body, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] request body: %s", body)

	res, err := c.HTTP("POST", "/v1/projects", &RequestOptions{
		Body: bytes.NewReader(body),
	})
	if err != nil {
		return nil, err
	}

	var r ProjectCreateResponse
	if err := decodeJSON(res, &r); err != nil {
		return nil, err
	}

	return &r, nil
}

// DeleteProject deletes project by project sub-domain name
func (c *Client) DeleteProject(name string) error {
	_, err := c.HTTP("DELETE", `/v1/projects/`+name, nil)
	if err != nil {
		return err
	}

	return nil
}
