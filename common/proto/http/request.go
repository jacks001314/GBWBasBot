package http

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

const (
	getMethod    = "get"
	postMethod   = "post"
	putMethod    = "put"
	deleteMethod = "delete"
	headMethod   = "head"
)

type HttpRequest struct {
	method                 string
	uri                    string
	headers                map[string]string
	basicUser, basicPasswd string
	multipartBuffer        bytes.Buffer
	writer                 *multipart.Writer
	body                   *bytes.Buffer
}

func NewHttpRequest(method string, uri string) (r *HttpRequest) {

	return &HttpRequest{
		method:      method,
		uri:         uri,
		headers:     make(map[string]string),
		basicUser:   "",
		basicPasswd: "",
		writer:      nil,
		body:        nil,
	}
}

func (r *HttpRequest) BasicAuth(user string, passwd string) *HttpRequest {

	r.basicUser = user
	r.basicPasswd = passwd

	return r
}

func (r *HttpRequest) AddHeader(k string, v string) *HttpRequest {

	r.headers[k] = v
	return r
}

func (r *HttpRequest) Headers(headers map[string]string) *HttpRequest {

	r.headers = headers

	return r
}

func (r *HttpRequest) BodyString(content string, fromFile bool) *HttpRequest {

	var data []byte
	var err error

	if fromFile {

		if data, err = ioutil.ReadFile(content); err != nil {

			data = []byte(``)
		}

	} else {

		data = []byte(content)
	}

	r.body = bytes.NewBuffer(data)

	return r
}

func (r *HttpRequest) BodyHex(content string) *HttpRequest {

	var data []byte
	var err error

	if data, err = hex.DecodeString(content); err != nil {

		data = []byte(``)
	}

	r.body = bytes.NewBuffer(data)

	return r
}

func (r *HttpRequest) UPload(name string, fpath string, boundary string) *HttpRequest {

	if r.writer == nil {
		r.writer = multipart.NewWriter(&r.multipartBuffer)
	}

	if boundary != "" {

		r.writer.SetBoundary(boundary)
	}

	f, err := os.Open(fpath)
	if err != nil {
		return r
	}
	defer f.Close()

	// Add file
	fw, err := r.writer.CreateFormFile(name, name)
	if err != nil {
		return r
	}

	if _, err = io.Copy(fw, f); err != nil {
		return r
	}

	r.body = &r.multipartBuffer

	return r
}

func (r *HttpRequest) makeURL(proto string, host string, port int) string {

	uri := r.uri

	if uri == "" {
		uri = "/"
	}

	if !strings.HasPrefix(uri, "/") {
		uri = "/" + uri
	}

	if port == 0 {
		return proto + "://" + host + uri
	}

	return fmt.Sprintf("%s://%s:%d%s", proto, host, port, uri)

}

/*build http.request */
func (r *HttpRequest) Build(proto string, host string, port int) (req *http.Request, err error) {

	body := r.body

	if !strings.EqualFold(r.method, getMethod) && body == nil {
		body = bytes.NewBuffer([]byte(``))
	}

	url := r.makeURL(proto, host, port)

	if r.writer != nil {

		r.writer.Close()
	}

	switch {
	case strings.EqualFold(r.method, getMethod):
		req, err = http.NewRequest(http.MethodGet, url, nil)

	case strings.EqualFold(r.method, postMethod):
		req, err = http.NewRequest(http.MethodPost, url, body)
	case strings.EqualFold(r.method, putMethod):
		req, err = http.NewRequest(http.MethodPut, url, body)

	case strings.EqualFold(r.method, postMethod):
		req, err = http.NewRequest(http.MethodPut, url, body)

	case strings.EqualFold(r.method, deleteMethod):
		req, err = http.NewRequest(http.MethodDelete, url, body)

	case strings.EqualFold(r.method, headMethod):
		req, err = http.NewRequest(http.MethodHead, url, body)

	default:
		req, err = http.NewRequest(strings.ToUpper(r.method), url, r.body)
	}

	if err != nil {
		return nil, fmt.Errorf("Cannot create http request,method:%s,url:%s,err:%v", r.method, r.uri, err)
	}

	if r.basicUser != "" && r.basicPasswd != "" {
		req.SetBasicAuth(r.basicUser, r.basicPasswd)
	}

	for k, v := range r.headers {
		req.Header[k] = []string{v}
		//req.Header.Set(k,v)
	}

	return
}
