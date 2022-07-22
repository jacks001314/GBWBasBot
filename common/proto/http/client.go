package http

import (
	"net/http"
	"time"
)

type HttpClient struct {
	client *http.Client
	proto  string
	host   string
	port   int
}

func NewHttpClient(host string, port int, isSSL bool, timeout int64) (httpClient *HttpClient) {

	tr := http.DefaultTransport.(*http.Transport)
	proto := "http"

	if isSSL {
		proto = "https"
		if port == 443 {
			port = 0
		}
	} else {

		if port == 80 {
			port = 0
		}
	}

	client := http.Client{
		Transport:     tr,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       time.Duration(timeout) * time.Millisecond,
	}

	return &HttpClient{
		client: &client,
		proto:  proto,
		host:   host,
		port:   port,
	}

}

func (c *HttpClient) Send(req *HttpRequest) (*HttpResponse, error) {

	/*make http request*/
	request, err := req.Build(c.proto, c.host, c.port)

	if err != nil {

		return nil, err
	}

	res, err := c.client.Do(request)

	if err != nil {

		return nil, err
	}

	return &HttpResponse{resp: res}, nil
}
