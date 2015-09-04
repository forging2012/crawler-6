package crawler

import (
	"net/http"
	"net/url"
	"time"
)

type URL struct {
	*url.URL
	Str string
}

type Page struct {
	URL           URL
	Content       []byte
	ContentType   string
	ContentLength int64
}

type Doc struct {
	baseURL URL
	HTML    []byte
}

type Request struct {
	method, url string
	body        []byte
	client      *http.Client
	config      func(*http.Request)
}

type Response struct {
	*http.Response
	Locations    URL // distinguish with method Location
	ContentType  string
	Content      []byte
	Time         time.Time
	LastModified time.Time
	Expires      time.Time
}

type Pool struct {
	size    int
	workers []Worker
	free    chan *Worker
}

type Worker struct {
	req  chan *Request
	resp chan *Response
	err  chan error
	pool *Pool
}
