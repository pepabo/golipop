package lolp

import (
	"bytes"
	"encoding/json"
	"log"
	"time"
)

// Project struct
type Project struct {
	ID            string                 `json:"id,omitempty`
	UserID        string                 `json:"userID,omitempty"`
	Name          string                 `json:"name,omitempty"`
	Kind          string                 `json:"kind"`
	Domain        string                 `json:"domain,omitempty"`
	SubDomain     string                 `json:"sub_domain,omitempty"`
	CustomDomains []string               `json:"custom_domains,omitempty"`
	Database      map[string]interface{} `json:"database,omitempty"`
	CreatedAt     time.Time              `json:"createdAt,omitempty"`
	UpdatedAt     time.Time              `json:"updatedAt,omitempty"`
}

// ProjectCreate struct on create
type ProjectCreate struct {
	UserID        string                 `json:"userID,omitempty"`
	Name          string                 `json:"name,omitempty"`
	Kind          string                 `json:"kind"`
	Domain        string                 `json:"domain,omitempty"`
	SubDomain     string                 `json:"sub_domain,omitempty"`
	CustomDomains []string               `json:"custom_domains,omitempty"`
	Payload       map[string]interface{} `json:"payload,omitempty"`
	Database      map[string]interface{} `json:"database,omitempty"`
}

// Projects returns project list
func (c *Client) Projects(opts map[string]interface{}) (*[]Project, error) {
	log.Printf("[INFO] listing project")

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

// CreateProject creates project with kind
func (c *Client) CreateProject(values map[string]interface{}) (*Project, error) {
	log.Printf("[INFO] creating project")

	data := &ProjectCreate{}
	for k, v := range values {
		err := SetField(data, k, v)
		if err != nil {
			return nil, err
		}
	}

	body, err := json.Marshal(data)
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
