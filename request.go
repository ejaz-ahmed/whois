package whois

import (
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

// DefaultTimeout for whois queries.
const DefaultTimeout = 10 * time.Second

// Request represents a whois request.
type Request struct {
	Query   string
	Host    string
	URL     string
	Body    string
	Timeout time.Duration
}

// NewRequest returns a request ready to fetch.
func NewRequest(q string) *Request {
	return &Request{Query: q, Timeout: DefaultTimeout}
}

// Fetch queries a whois server via whois protocol or by HTTP if URL is set.
func (req *Request) Fetch() (*Response, error) {
	if req.URL != "" {
		return req.fetchURL()
	}
	return req.fetchWhois()
}

func (req *Request) fetchWhois() (*Response, error) {
	res := &Response{Request: req, FetchedAt: time.Now()}

	c, err := net.DialTimeout("tcp", req.Host+":43", req.Timeout)
	if err != nil {
		return res, err
	}
	defer c.Close()
	c.SetDeadline(time.Now().Add(req.Timeout))
	if _, err = io.WriteString(c, req.Body); err != nil {
		return res, err
	}
	if res.Body, err = ioutil.ReadAll(c); err != nil {
		return res, err
	}

	res.ContentType = http.DetectContentType(res.Body)

	return res, nil
}

func (req *Request) fetchURL() (*Response, error) {
	res := &Response{Request: req, FetchedAt: time.Now()}

	getResp, err := http.Get(req.URL)
	if err != nil {
		return res, err
	}
	defer getResp.Body.Close()
	if res.Body, err = ioutil.ReadAll(getResp.Body); err != nil {
		return res, err
	}

	res.ContentType = http.DetectContentType(res.Body)

	return res, nil
}
