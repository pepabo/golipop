package lolp

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestCreateProject(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(projectCreateHandler()))
	defer s.Close()

	c, err := NewClient(s.URL)
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		arg        ProjectNew
		wantReturn Project
	}{
		{
			ProjectNew{Kind: "wordpress", Payload: map[string]interface{}{"username": "foo", "email": "foo@example.com", "password": "Secret#Gopher123?"}},
			Project{Domain: "foobar-baz-9999.lolipop.io"},
		},
		{
			ProjectNew{Kind: "rails", Database: map[string]interface{}{"password": "Secret#Gopher123?"}},
			Project{Domain: "foobar-baz-9999.lolipop.io"},
		},
	}

	for _, cc := range cases {
		r, err := c.CreateProject(&cc.arg)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(&cc.wantReturn, r) {
			t.Errorf("return object\nexpect: %#v\ngot: %#v", cc.wantReturn, r)
		}
	}
}

func TestProjects(t *testing.T) {
}

func TestProject(t *testing.T) {
}

func TestDeleteProject(t *testing.T) {
}
