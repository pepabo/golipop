package lolp

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"testing"
)

type mcServer struct {
	URL    *url.URL
	t      *testing.T
	ln     net.Listener
	server *http.Server
}

type clientTestRes struct {
	RawPath string
	Host    string
	Header  http.Header
	Body    string
}

func newTestmcServer(t *testing.T) *mcServer {
	hs := &mcServer{t: t}

	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}
	hs.ln = ln

	hs.URL = &url.URL{
		Scheme: "http",
		Host:   ln.Addr().String(),
	}

	mux := http.NewServeMux()
	hs.setupRoutes(mux)

	server := &http.Server{}
	server.Handler = mux
	hs.server = server
	go server.Serve(ln)

	return hs
}

func (hs *mcServer) Stop() {
	hs.ln.Close()
}

func (hs *mcServer) setupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/_json", hs.appHandler)
}

func (hs *mcServer) appHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(422)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"errors": ["this is an error", "this is another error"]}`)
}

func (hs *mcServer) authenticationHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		hs.t.Fatal(err)
	}
	username, password := r.Form["username"][0], r.Form["password"][0]
	if username == "foo@example.com" && password == "Secret#Gopher123" {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"pX4AQ5vO7T-xJrxsnvlB0cfeF-tGUX-A-280LPxoryhDAbwmox7PKinMgA1F6R3BKaT"}`)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}
