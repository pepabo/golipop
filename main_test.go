package lolp

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	o := new(bytes.Buffer)
	log.SetOutput(o)
	code := m.Run()
	os.Exit(code)
}

func loginHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`eyJ0eXAiOiJKV1QiLCJ.eyJzdWIiOiIxOTg4M2VjZi0yYWUyLTRmZmUtYjg1MS03MzMTgwMDMsImV4cCI6MT2MzcxYyJ9.Jm6Pzn4bXFHLHy2HNqJY`))
	}
}

func projectCreateHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"domain": "foobar-baz-9999.lolipop.io"}`)
		//w.Write([]byte(`{"kind":"wordpress","payload":{"email":"foo@example.com","password":"Secret#Gopher123?","username":"foo"}}`))
	}
}
