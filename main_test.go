package lolp

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	o := new(bytes.Buffer)
	log.SetOutput(o)
	code := m.Run()
	os.Exit(code)
}

func fixture(ctx string, r *http.Request) string {
	f := fmt.Sprintf("%s--%s.json", r.Method, ctx)
	p := filepath.Join("testdata", r.URL.String(), f)
	b, err := ioutil.ReadFile(p)
	if err != nil {
		panic("fixture not found: " + p)
	}
	return strings.TrimSuffix(string(b), "\n")
}
