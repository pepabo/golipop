package lolp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"time"
)

func SetField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return errors.New("Provided value type didn't match obj field type")
	}

	if reflect.TypeOf(val).String() != "" {
		structFieldValue.Set(val)
	}
	return nil
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
