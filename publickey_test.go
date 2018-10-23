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
)

func publickeyAddHandler(t *testing.T) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err.Error())
		}

		p := &PublicKey{}
		if err := json.Unmarshal(body, &p); err != nil {
			panic(err.Error())
		}
		ctx := p.Name

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

func TestAddPublickey(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(publickeyAddHandler(t)))
	defer s.Close()

	c, err := NewClient(s.URL)
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		arg        PublicKey
		wantReturn PublicKey
	}{
		{
			PublicKey{Name: "dummy", Key: "dummy"},
			PublicKey{Name: "dummy", Key: "dummy"},
		},
	}

	for _, cc := range cases {
		r, err := c.AddPublicKey(&cc.arg)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(&cc.wantReturn, r) {
			t.Errorf("return object\nexpect: %#v\ngot: %#v", cc.wantReturn, r)
		}
	}
}

func publickeyDeleteHandler(t *testing.T) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err.Error())
		}
		var ctx string
		name := path.Base(r.RequestURI)
		if name == "dummy-ok" {
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
		} else if ctx == "ng" {
			w.WriteHeader(http.StatusNotFound)
		}
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, fixture(ctx+".response", r))
	}
}

func TestDeletePublickey(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(publickeyDeleteHandler(t)))
	defer s.Close()

	c, err := NewClient(s.URL)
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name    string
		wantErr bool
	}{
		{"dummy-ok", false},
		{"dummy-ng", true},
	}

	for _, cc := range cases {
		err := c.DeletePublicKey(cc.name)
		if cc.wantErr {
			if err == nil {
				t.Errorf("expect public key delete failure but succeeded")
			}
			return
		}
		if err != nil {
			t.Errorf("expect to succeed in public key delete, but failed: %s", err)
		}
	}
}
