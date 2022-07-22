package http

import (
	"fmt"

	stengo "common/scripts/tengo"

	"github.com/d5/tengo/objects"
	"github.com/d5/tengo/v2"
)

/*for create a http client function*/
type HttpClientTengo struct {
	stengo.TengoObj

	client *HttpClient
}

/*for http send  function*/
type HttpClientMethodTengo struct {
	stengo.TengoObj
	clientTengo *HttpClientTengo
}

func (hc *HttpClientTengo) IndexGet(index objects.Object) (value objects.Object, err error) {

	key, ok := objects.ToString(index)

	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "index",
			Expected: "string(compatible)",
			Found:    index.TypeName(),
		}
	}

	switch key {

	case sendMethod:
		return &HttpClientMethodTengo{
			TengoObj:    stengo.TengoObj{Name: sendMethod},
			clientTengo: hc,
		}, nil

	default:
		return nil, fmt.Errorf("undefine http client method:%s", key)
	}

}

func (hcm *HttpClientMethodTengo) sendRequest(args ...objects.Object) (objects.Object, error) {

	if len(args) != 1 {

		return nil, tengo.ErrWrongNumArguments
	}

	request, ok := args[0].(*HttpRequestTengo)

	if !ok {

		return nil, tengo.ErrInvalidArgumentType{
			Name:     "request",
			Expected: "httpRequestTengo(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	response, err := hcm.clientTengo.client.Send(request.req)

	if err != nil {
		return nil, err
	}

	return &HttpResponseTengo{
		TengoObj: stengo.TengoObj{Name: "response"},
		res:      response,
	}, nil
}

func (hcm *HttpClientMethodTengo) Call(args ...objects.Object) (objects.Object, error) {

	switch hcm.Name {
	case sendMethod:
		return hcm.sendRequest(args...)

	default:
		return nil, fmt.Errorf("unknown http client method:%s", hcm.Name)

	}

}
