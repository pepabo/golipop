package lolp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"
)

// Project struct
type Project struct {
	ID            string         `json:"id,omitempty"`
	Name          string         `json:"name,omitempty"`
	Kind          string         `json:"kind,omitempty"`
	Domain        string         `json:"domain,omitempty"`
	SubDomain     string         `json:"subDomain,omitempty"`
	Autoscalable  bool           `json:"autoscalable,omitempty"`
	CustomDomains []CustomDomain `json:"customDomains,omitempty"`
	Database      Database       `json:"database,omitempty"`
	SSH           *SSH           `json:"ssh,omitempty"`
	CreatedAt     time.Time      `json:"createdAt,omitempty"`
	UpdatedAt     time.Time      `json:"updatedAt,omitempty"`
}

type CustomDomain struct {
	Name string `json:"name,omitempty"`
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

// EnableAutoscaling enable autoscaling by project sub-domain name
func (c *Client) EnableAutoscaling(name string) error {
	_, err := c.HTTP("PUT", `/v1/projects/`+name+`/autoscaling/enable`, nil)
	if err != nil {
		return err
	}

	return nil
}

// DisableAutoscaling disable autoscaling by project sub-domain name
func (c *Client) DisableAutoscaling(name string) error {
	_, err := c.HTTP("PUT", `/v1/projects/`+name+`/autoscaling/disable`, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetEnvironmentVariables(name string) (string, error) {
	res, err := c.HTTP("GET", `/v1/projects/`+name+`/environment-variables`, nil)
	if err != nil {
		return "", err
	}

	resbody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(resbody), nil
}

type UpdateEnvironmentVariablesParam struct {
	Method   string `json:"method"`
	Variable struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"variable"`
}

func (c *Client) UpdateEnvironmentVariables(name string, params []UpdateEnvironmentVariablesParam) error {
	body, err := json.Marshal(params)
	if err != nil {
		return err
	}

	_, err = c.HTTP("PUT", `/v1/projects/`+name+`/environment-variables`, &RequestOptions{
		Body: bytes.NewReader(body),
	})
	if err != nil {
		return err
	}

	return nil
}
