package lolp

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path"
	"reflect"
	"testing"
	"time"
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
		wantReturn ProjectCreateResponse
	}{
		{
			ProjectNew{Kind: "wordpress", Payload: map[string]interface{}{"username": "foo", "email": "foo@example.com", "password": "Secret#Gopher123?"}},
			ProjectCreateResponse{Domain: "foobar-baz-9999.lolipop.io"},
		},
		{
			ProjectNew{Kind: "rails", DBPassword: "Secret#Gopher123?"},
			ProjectCreateResponse{Domain: "foobar-baz-9999.lolipop.io"},
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

func projectListHandler(t *testing.T) func(http.ResponseWriter, *http.Request) {
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
	s := httptest.NewServer(http.HandlerFunc(projectListHandler(t)))
	defer s.Close()

	c, err := NewClient(s.URL)
	if err != nil {
		t.Fatal(err)
	}

	tt, _ := time.Parse(time.RFC3339, "2018-02-13T08:36:06.380Z")
	cases := []struct {
		ret []Project
	}{
		{[]Project{
			Project{ID: "58b22c80-5c64-41ed-ac51-7ca0c695e592", Kind: "rails", Name: "rails-1.lolipop.io", CreatedAt: tt, UpdatedAt: tt},
			Project{ID: "507dbd34-d6af-49d5-9d3d-98933c02a019", Kind: "php", Name: "php-1.lolipop.io", CreatedAt: tt, UpdatedAt: tt},
			Project{ID: "b3585f32-2265-418e-a762-45894620e1e0", Kind: "wordpress", Name: "wordpress-1.lolipop.io", CreatedAt: tt, UpdatedAt: tt},
		}},
	}

	for _, cc := range cases {
		r, err := c.Projects()
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(&cc.ret, r) {
			t.Errorf("return object\nexpect: %#v\ngot: %#v", cc.ret, *r)
		}
	}
}

func projectHandler(t *testing.T) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err.Error())
		}

		ctx := "ok"
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

func TestProject(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(projectHandler(t)))
	defer s.Close()

	c, err := NewClient(s.URL)
	if err != nil {
		t.Fatal(err)
	}

	tt, _ := time.Parse(time.RFC3339, "2018-02-13T08:36:06.380Z")
	cases := []struct {
		name string
		ret  Project
	}{
		{
			"rails-1",
			Project{
				ID:        "58b22c80-5c64-41ed-ac51-7ca0c695e592",
				Kind:      "rails",
				Name:      "rails-1.lolipop.io",
				Domain:    "rails-1.lolipop.io",
				SubDomain: "rails-1",
				Database: Database{
					Host: "mysql-1.mc.lolipop.lan",
					Name: "7e7aef038f314742c064deb6e6e84714",
					User: "7e7aef038f314742c064deb6e6e84714",
				},
				SSH: &SSH{
					User: "sweet-ebino-9052",
					Host: "ssh-1.mc.lolipop.jp",
					Port: 12345,
				},
				CustomDomains: []string{},
				CreatedAt:     tt,
				UpdatedAt:     tt,
			},
		},
	}

	for _, cc := range cases {
		r, err := c.Project(cc.name)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(&cc.ret, r) {
			t.Errorf("return object\nexpect: %#v\ngot: %#v", cc.ret, *r)
		}
	}
}

func projectDeleteHandler(t *testing.T) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err.Error())
		}

		var ctx string
		name := path.Base(r.RequestURI)
		if name == "rails-1" {
			ctx = "ok"
		} else {
			ctx = "ng"
		}

		expected := ""
		actual := string(body)
		if expected != actual {
			t.Errorf("request body\nexpected: %s\nactual: %s", expected, actual)
		}

		if ctx == "ok" {
			w.WriteHeader(http.StatusNoContent)
			w.Header().Set("Content-Type", "application/json")
			io.WriteString(w, fixture(ctx+".response", r))
		}

		if ctx == "ng" {
			w.WriteHeader(http.StatusNotFound)
			w.Header().Set("Content-Type", "application/json")
			io.WriteString(w, fixture(ctx+".response", r))
		}
	}
}

func TestDeleteProject(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(projectDeleteHandler(t)))
	defer s.Close()

	c, err := NewClient(s.URL)
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name        string
		expectedErr bool
	}{
		{"rails-1", false},
		{"not-exist", true},
	}

	for _, cc := range cases {
		err := c.DeleteProject(cc.name)
		if cc.expectedErr {
			if err == nil {
				t.Errorf("expect project delete failure but succeeded")
			}
		} else {
			if err != nil {
				t.Errorf("expect to succeed in project delete, but failed: %s", err)
			}
		}
	}
}
