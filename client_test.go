package lolp

import (
	"os"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	now := os.Getenv("GOLIPOP_TOKEN")
	os.Unsetenv("GOLIPOP_TOKEN")
	defer os.Setenv("GOLIPOP_TOKEN", now)
	defer os.Unsetenv("GOLIPOP_ENDPOINT")

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
	uae := "lolp/0.0.1 (+https://github.com/pepabo/golipop; go1.9.3)"
	if ua != uae {
		t.Errorf("User-Agent header expects %s, but got %s", uae, ua)
	}

	os.Setenv("GOLIPOP_ENDPOINT", "https://example.com/")
	cc := New()
	if cc.URL.String() != "https://example.com/" {
		t.Errorf("client URL is wrong: %s", cc.URL)
	}
}
