package lolp

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	now := os.Getenv("LOLP_TOKEN")
	os.Unsetenv("LOLP_TOKEN")
	defer os.Setenv("LOLP_TOKEN", now)
	defer os.Unsetenv("LOLP_ENDPOINT")

	c := New()
	if c.URL.String() != "https://api.mc.lolipop.jp/" {
		t.Errorf("client URL is wrong: %s", c.URL)
	}
	if c.Token != "" {
		t.Errorf("API token expects empty, but got %s", c.Token)
	}
	ct := strings.Join(c.DefaultHeader["Content-Type"], "")
	cte := "application/json"
	if ct != cte {
		t.Errorf("Content-Type header expects %s, but got %s", cte, ct)
	}
	ua := strings.Join(c.DefaultHeader["User-Agent"], "")
	uae := "lolp/0.0.1 (+https://github.com/pepabo/golipop; go"
	if !strings.HasPrefix(ua, uae) {
		t.Errorf("User-Agent header expects %s, but got %s", uae, ua)
	}

	dummyEndpoint := "https://example.com/"
	os.Setenv("LOLP_ENDPOINT", dummyEndpoint)
	cc := New()
	if cc.URL.String() != dummyEndpoint {
		t.Errorf("client URL is wrong: %s", cc.URL)
	}
}

func TestNewClient(t *testing.T) {
	_, err := NewClient("")
	if err == nil {
		t.Errorf("empty string as argument expects error")
	}

	dummyEndpoint := "https://example.com/"
	c, err := NewClient(dummyEndpoint)
	if c.URL.String() != dummyEndpoint {
		t.Errorf("client URL is wrong: %s", c.URL)
	}
}

func TestClientInit(t *testing.T) {
	os.Setenv("LOLP_TLS_NOVERIFY", "true")
	defer os.Unsetenv("LOLP_TLS_NOVERIFY")

	c := &Client{
		DefaultHeader: make(http.Header),
	}
	c.init()
	if !c.HTTPClient.Transport.(*http.Transport).TLSClientConfig.InsecureSkipVerify {
		t.Errorf("skip verify expects true")
	}
}

func TestClientRequest(t *testing.T) {
	o := new(bytes.Buffer)
	log.SetOutput(o)

	c, err := NewClient("https://api.example.com/")
	if err != nil {
		t.Fatal(err)
	}
	c.Token = "secret"
	req, err := c.Request("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	if req.Method != "GET" {
		t.Errorf("HTTP method is wrong: %s", req.Method)
	}
	if req.URL.String() != "https://api.example.com/test" {
		t.Errorf("HTTP URL is wrong: %s", req.URL)
	}
	a := strings.Join(req.Header["Authorization"], "")
	if a != "Bearer secret" {
		t.Errorf("Authorization header is wrong: %s", a)
	}
}
