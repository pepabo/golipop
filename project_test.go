package lolp

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func projectCreateHandler(t *testing.T) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err.Error())
		}

		p := &Project{}
		if err := json.Unmarshal(body, &p); err != nil {
			panic(err.Error())
		}
		ctx := p.Kind

		expected := fixture(ctx+".request", r)
		actual := string(body)
		if expected != actual {
			t.Errorf("request body\nexpected: %s\nactual: %s", expected, actual)
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, fixture(ctx+".response", r))
	}
}

func TestCreateProject(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(projectCreateHandler(t)))
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

func projectsHandler(t *testing.T) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err.Error())
		}

		ctx := "all"
		expected := fixture(ctx+".request", r)
		actual := string(body)
		if expected != actual {
			t.Errorf("request body\nexpected: %s\nactual: %s", expected, actual)
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, fixture(ctx+".response", r))
	}
}

func TestProjects(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(projectsHandler(t)))
	defer s.Close()

	c, err := NewClient(s.URL)
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		ret []Project
	}{
		{[]Project{Project{Domain: "node-1.lolipop.io"}, Project{Domain: "php-1.lolipop.io"}, Project{Domain: "rails-1.lolipop.io"}}},
	}

	for _, cc := range cases {
		r, err := c.Projects()
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(&cc.ret, r) {
			t.Errorf("return object\nexpect: %#v\ngot: %#v", cc.ret, r)
		}
	}
}

func TestProject(t *testing.T) {
}

func TestDeleteProject(t *testing.T) {
}
