package lolp

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func authenticateHandler(t *testing.T) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err.Error())
		}

		l := struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}{}
		if err := json.Unmarshal(body, &l); err != nil {
			panic(err.Error())
		}

		var ctx string
		if l.Username == "foo@example.com" && l.Password == "Secret#Gopher123?" {
			ctx = "ok"
			w.WriteHeader(http.StatusOK)
		} else {
			ctx = "ng"
			w.WriteHeader(http.StatusUnauthorized)
		}
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, fixture(ctx+".response", r))
	}
}

func TestAuthenticate(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(authenticateHandler(t)))
	defer s.Close()

	c, err := NewClient(s.URL)
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		username    string
		password    string
		expectedErr bool
	}{
		{"foo@example.com", "Secret#Gopher123?", false},
		{"foo@example.com", "Secret#Gopher999?", true},
	}

	for _, cc := range cases {
		token, err := c.Authenticate(cc.username, cc.password)
		if cc.expectedErr {
			if err == nil {
				t.Errorf("expect authentication failure but succeeded")
			}
			if token != "" {
				t.Errorf("expect token empty but it returns as (%s)", token)
			}
		} else {
			if err != nil {
				t.Errorf("expect to succeed in authentication, but failed: %s", err)
			}
			if token == "" {
				t.Errorf("expect token but it is empty")
			}
		}
	}
}
