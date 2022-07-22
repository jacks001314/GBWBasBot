package http

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

type HttpResponse struct {
	resp *http.Response
}

// GetResp get net/http original response
func (r *HttpResponse) GetResp() *http.Response {
	return r.resp
}

// GetStatusCode returns http status code
// if Response is not returned from a Request
// the status code will be 0
func (r *HttpResponse) GetStatusCode() int {
	if r.resp == nil {
		return 0
	}

	return r.resp.StatusCode
}

// GetBody returns response body
// It is the caller's responsibility to close Body
func (r *HttpResponse) GetBody() io.ReadCloser {
	if r.resp == nil {
		return nil
	}

	return r.resp.Body
}

// GetBodyAsByte returns response body as byte
func (r *HttpResponse) GetBodyAsByte() ([]byte, error) {

	body := r.GetBody()
	if body == nil {
		return nil, nil
	}
	defer body.Close()

	byts, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}

	return byts, nil
}

// GetBodyAsString returns response body as string
func (r *HttpResponse) GetBodyAsString() (string, error) {
	body, err := r.GetBodyAsByte()
	if err != nil || body == nil {
		return "", err
	}

	return string(body), nil
}

// GetBodyAsJSONRawMessage returns response body as json.RawMessage
func (r *HttpResponse) GetBodyAsJSONRawMessage() (json.RawMessage, error) {
	body, err := r.GetBodyAsByte()
	if err != nil || body == nil {
		return nil, err
	}

	return json.RawMessage(body), nil
}

// UnmarshalBody unmarshal response body
func (r *HttpResponse) UnmarshalBody(v interface{}) error {
	body, err := r.GetBodyAsByte()
	if err != nil || body == nil {
		return err
	}

	return json.Unmarshal(body, &v)
}

//Protocol returns response proto
func (r *HttpResponse) Protocol() string {
	return r.resp.Proto
}

//URL returns response Location
func (r *HttpResponse) URL() (*url.URL, error) {
	return r.resp.Location()
}

/*get first header value with key*/
func (r *HttpResponse) GetHeaderValue(k string) string {
	return r.resp.Header.Get(k)
}

/*get all values with key*/
func (r *HttpResponse) GetHeaderValues(k string) []string {
	return r.resp.Header.Values(k)
}

func (r *HttpResponse) GetHader() http.Header {

	return r.resp.Header
}
