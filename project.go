package lolp

import (
	"bytes"
	"encoding/json"
	"log"
	"time"
)

// Project struct
type Project struct {
	ID            int                    `json:"id,omitempty"`
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

type ManagedConfig struct {
	DBName string `json:"db_name",omitempty`
	DBUser string `json:"db_user",omitempty`
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
	UserID        string                 `json:"userID,omitempty"`
	Name          string                 `json:"name,omitempty"`
	Kind          string                 `json:"kind,omitempty""`
	Domain        string                 `json:"domain,omitempty"`
	SubDomain     string                 `json:"sub_domain,omitempty"`
	CustomDomains []string               `json:"custom_domains,omitempty"`
	Payload       map[string]interface{} `json:"payload,omitempty"`
	Database      map[string]interface{} `json:"database,omitempty"`
}

// Projects returns project list
func (c *Client) Projects() (*[]Project, error) {
	log.Printf("[INFO] listing project list")

	request, err := c.Request("GET", "/v1/projects", &RequestOptions{})
	if err != nil {
		return nil, err
	}

	response, err := dispose(c.HTTPClient.Do(request))
	if err != nil {
		return nil, err
	}

	var projects []Project
	if err := decodeJSON(response, &projects); err != nil {
		return nil, err
	}

	return &projects, nil
}

// Project returns a project by sub-domain name
func (c *Client) Project(name string) (*ProjectGet, error) {
	log.Printf("[INFO] listing a project")

	request, err := c.Request("GET", `/v1/projects/`+name, &RequestOptions{})
	if err != nil {
		return nil, err
	}

	response, err := dispose(c.HTTPClient.Do(request))
	if err != nil {
		return nil, err
	}

	var project ProjectGet
	if err := decodeJSON(response, &project); err != nil {
		return nil, err
	}

	return &project, nil
}

// CreateProject creates project with kind
func (c *Client) CreateProject(p *ProjectNew) (*Project, error) {
	log.Printf("[INFO] creating project")

	body, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] request body: %s", body)

	request, err := c.Request("POST", "/v1/projects", &RequestOptions{
		Body: bytes.NewReader(body),
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

// DeleteProject deletes project by project sub-domain name
func (c *Client) DeleteProject(name string) error {
	request, err := c.Request("DELETE", `/v1/projects/`+name, &RequestOptions{})
	if err != nil {
		return err
	}

	_, err = dispose(c.HTTPClient.Do(request))
	if err != nil {
		return err
	}

	return nil
}
