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
	ID            string                 `json:"id,omitempty"`
	UserID        string                 `json:"userID,omitempty"`
	Name          string                 `json:"name,omitempty"`
	Kind          string                 `json:"kind,omitempty"`
	Domain        string                 `json:"domain,omitempty"`
	SubDomain     string                 `json:"sub_domain,omitempty"`
	CustomDomains []string               `json:"custom_domains,omitempty"`
	Database      map[string]interface{} `json:"database,omitempty"`
	CreatedAt     time.Time              `json:"createdAt,omitempty"`
	UpdatedAt     time.Time              `json:"updatedAt,omitempty"`
}

// ProjectGet struct
type ProjectGet struct {
	ID                 int                    `json:"id,omitempty"`
	UUID               string                 `json:"uuid,omitempty"`
	AccountHumaneID    string                 `json:"account_humane_id,omitempty"`
	SVM                string                 `json:"svm,omitempty"`
	Volume             string                 `json:"volume,omitempty"`
	DatbaseHost        string                 `json:"database_host",omitempty"`
	CustomDomains      []string               `json:"custom_domains,omitempty"`
	ContainerTemplates []interface{}          `json:"container_templates",omitempty`
	Containers         []interface{}          `json:"containers",omitempty`
	BaseSpec           map[string]interface{} `json:"base_spec",omitempty`
	ManagedConfig      map[string]interface{} `json:"managed_config,omitempty"`
	SSL                map[string]interface{} `json:"ssl",omitempty`
	CreatedAt          time.Time              `json:"created_at,omitempty"`
	UpdatedAt          time.Time              `json:"updated_at,omitempty"`
}

// ProjectNew struct on create
type ProjectNew struct {
	Name          string                 `json:"name,omitempty"`
	Kind          string                 `json:"kind,omitempty""`
	SubDomain     string                 `json:"sub_domain,omitempty"`
	CustomDomains []string               `json:"custom_domains,omitempty"`
	Payload       map[string]interface{} `json:"payload,omitempty"`
	Database      map[string]interface{} `json:"database,omitempty"`
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
func (c *Client) Project(name string) (*ProjectGet, error) {
	res, err := c.HTTP("GET", `/v1/projects/`+name, nil)
	if err != nil {
		return nil, err
	}

	var p ProjectGet
	if err := decodeJSON(res, &p); err != nil {
		return nil, err
	}

	return &p, nil
}

// CreateProject creates project with kind
func (c *Client) CreateProject(p *ProjectNew) (*Project, error) {
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

	var pp Project
	if err := decodeJSON(res, &pp); err != nil {
		return nil, err
	}

	return &pp, nil
}

// DeleteProject deletes project by project sub-domain name
func (c *Client) DeleteProject(name string) error {
	_, err := c.HTTP("DELETE", `/v1/projects/`+name, nil)
	if err != nil {
		return err
	}

	return nil
}
