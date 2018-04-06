package lolp

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"runtime"
	"strings"

	cleanhttp "github.com/hashicorp/go-cleanhttp"
)

const (
	// defaultEndpoint
	defaultEndpoint = "https://api.mc.lolipop.jp/"

	// EndpointEnvVar for endpoint
	EndpointEnvVar = "GOLIPOP_ENDPOINT"

	// TLSNoVerifyEnvVar for TLS verify skip flag
	TLSNoVerifyEnvVar = "GOLIPOP_TLS_NOVERIFY"

	// TokenEnvVar for authentication
	TokenEnvVar = "GOLIPOP_TOKEN"
)

// projectURL for this
var projectURL = "https://github.com/pepabo/golipop"

// userAgent for request
var userAgent = fmt.Sprintf("lolp/%s (+%s; %s)", Version, projectURL, runtime.Version())

// Client struct
type Client struct {
	URL           *url.URL
	HTTPClient    *http.Client
	DefaultHeader http.Header
	Token         string
}

// DefaultClient returns client struct pointer
func DefaultClient() *Client {
	endpoint := os.Getenv(EndpointEnvVar)
	if endpoint == "" {
		endpoint = defaultEndpoint
	}

	client, err := NewClient(endpoint)
	if err != nil {
		panic(err)
	}

	token := os.Getenv(TokenEnvVar)
	if token != "" {
		client.Token = token
	}

	return client
}

// NewClient returns clean client struct pointer
func NewClient(u string) (*Client, error) {
	if len(u) == 0 {
		return nil, fmt.Errorf("client: missing url")
	}

	parsedURL, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	c := &Client{
		URL:           parsedURL,
		DefaultHeader: make(http.Header),
	}

	if err := c.init(); err != nil {
		return nil, err
	}

	return c, nil
}

// init initializes for client
func (c *Client) init() error {
	c.DefaultHeader.Set("User-Agent", userAgent)
	c.DefaultHeader.Set("Content-Type", "application/json")
	c.HTTPClient = cleanhttp.DefaultClient()

	tlsConfig := &tls.Config{}
	if os.Getenv(TLSNoVerifyEnvVar) != "" {
		tlsConfig.InsecureSkipVerify = true
	}
	t := cleanhttp.DefaultTransport()
	t.TLSClientConfig = tlsConfig
	c.HTTPClient.Transport = t

	return nil
}

// RequestOptions struct
type RequestOptions struct {
	Params     map[string]string
	Headers    map[string]string
	Body       io.Reader
	BodyLength int64
}

// HTTP returns http.Response with dispose
func (c *Client) HTTP(verb, spath string, ro *RequestOptions) (*http.Response, error) {
	req, err := c.Request(verb, spath, ro)
	if err != nil {
		return nil, err
	}

	res, err := dispose(c.HTTPClient.Do(req))
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Request returns http.Request pointer with error
func (c *Client) Request(verb, spath string, ro *RequestOptions) (*http.Request, error) {
	log.Printf("[INFO] request: %s %s", verb, spath)

	if ro == nil {
		ro = new(RequestOptions)
	}

	u := *c.URL
	u.Path = path.Join(c.URL.Path, spath)

	if c.Token != "" {
		if ro.Headers == nil {
			ro.Headers = make(map[string]string)
		}
		ro.Headers["Authorization"] = fmt.Sprintf("Bearer %s", c.Token)
	}

	return c.rawRequest(verb, &u, ro)
}

// rawRequest returns http.Request pointer with error
func (c *Client) rawRequest(verb string, u *url.URL, ro *RequestOptions) (*http.Request, error) {
	if verb == "" {
		return nil, fmt.Errorf("client: missing verb")
	}

	if u == nil {
		return nil, fmt.Errorf("client: missing URL.url")
	}

	if ro == nil {
		return nil, fmt.Errorf("client: missing RequestOptions")
	}

	var params = make(url.Values)
	for k, v := range ro.Params {
		params.Add(k, v)
	}
	u.RawQuery = params.Encode()

	request, err := http.NewRequest(verb, u.String(), ro.Body)
	if err != nil {
		return nil, err
	}

	for k, v := range c.DefaultHeader {
		request.Header[k] = v
	}

	for k, v := range ro.Headers {
		request.Header.Add(k, v)
	}

	if ro.BodyLength > 0 {
		request.ContentLength = ro.BodyLength
	}

	log.Printf("[DEBUG] raw request: %#v", request)

	return request, nil
}

// dispose returns http.Request pointer with error
func dispose(res *http.Response, err error) (*http.Response, error) {
	if err != nil {
		return res, err
	}

	log.Printf("[INFO] response: %d (%s)", res.StatusCode, res.Status)
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, res.Body); err != nil {
		log.Printf("[ERR] response: error copying response body")
	} else {
		log.Printf("[DEBUG] response: %s", buf.String())
		res.Body.Close()
		res.Body = &bytesReadCloser{&buf}
	}

	switch res.StatusCode {
	case 200:
		return res, nil
	case 201:
		return res, nil
	case 202:
		return res, nil
	case 204:
		return res, nil
	case 400:
		return nil, parseErr(res)
	case 401:
		return nil, fmt.Errorf("authentication failed")
	case 404:
		return nil, fmt.Errorf("resource not found")
	case 422:
		return nil, parseErr(res)
	default:
		return nil, fmt.Errorf("client: %s", res.Status)
	}
}

// parseErr parses for error response
func parseErr(r *http.Response) error {
	re := &AppError{}

	if err := decodeJSON(r, &re); err != nil {
		return fmt.Errorf("error decoding JSON body: %s", err)
	}

	return re
}

// decodeJSON decodes for response
func decodeJSON(res *http.Response, out interface{}) error {
	defer res.Body.Close()
	dec := json.NewDecoder(res.Body)
	return dec.Decode(out)
}

// bytesReadCloser struct
type bytesReadCloser struct {
	*bytes.Buffer
}

// Close returns nil
func (nrc *bytesReadCloser) Close() error {
	return nil
}

// AppError struct
type AppError struct {
	Errors []string `json:"errors"`
}

// Error returns error by string
func (re *AppError) Error() string {
	return strings.Join(re.Errors, ", ")
}
