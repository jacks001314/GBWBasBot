package http

import (
	"errors"

	stengo "common/scripts/tengo"

	"github.com/d5/tengo/objects"
	"github.com/d5/tengo/v2"
)

var resMethodMaps = map[string]*HttpResponseMethod{

	getStatusCodeMethod: &HttpResponseMethod{
		TengoObj: stengo.TengoObj{Name: getStatusCodeMethod},
	},

	getBodyAsByteMethod: &HttpResponseMethod{
		TengoObj: stengo.TengoObj{Name: getBodyAsByteMethod},
	},

	getBodyAsStringMethod: &HttpResponseMethod{
		TengoObj: stengo.TengoObj{Name: getBodyAsStringMethod},
	},

	getProtocolMethod: &HttpResponseMethod{
		TengoObj: stengo.TengoObj{Name: getProtocolMethod},
	},

	getHeaderMethod: &HttpResponseMethod{
		TengoObj: stengo.TengoObj{Name: getHeaderMethod},
	},

	getHeadersMethod: &HttpResponseMethod{
		TengoObj: stengo.TengoObj{Name: getHeadersMethod},
	},
}

/*for http response*/
type HttpResponseTengo struct {
	stengo.TengoObj
	res *HttpResponse
}

type HttpResponseMethod struct {
	stengo.TengoObj
	resTengo *HttpResponseTengo
}

func (res *HttpResponseTengo) IndexGet(index objects.Object) (value objects.Object, err error) {

	key, ok := objects.ToString(index)

	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "index",
			Expected: "string(compatible)",
			Found:    index.TypeName(),
		}
	}

	if m, ok := resMethodMaps[key]; ok {

		m.resTengo = res
		return m, nil
	}

	return nil, errors.New("Undefine http response function:" + key)

}

func (hrm *HttpResponseMethod) getStatusCode() (objects.Object, error) {

	return objects.FromInterface(hrm.resTengo.res.GetStatusCode())
}

func (hrm *HttpResponseMethod) getBodyAsByte() (objects.Object, error) {

	content, _ := hrm.resTengo.res.GetBodyAsByte()

	return objects.FromInterface(content)
}

func (hrm *HttpResponseMethod) getBodyAsString() (objects.Object, error) {

	content, _ := hrm.resTengo.res.GetBodyAsString()
	return objects.FromInterface(content)
}

func (hrm *HttpResponseMethod) getProtocol() (objects.Object, error) {

	content := hrm.resTengo.res.Protocol()
	return objects.FromInterface(content)
}

func (hrm *HttpResponseMethod) getHeader(args ...objects.Object) (objects.Object, error) {

	if len(args) != 1 {

		return nil, tengo.ErrWrongNumArguments
	}

	key, ok := objects.ToString(args[0])

	if !ok {

		return nil, tengo.ErrInvalidArgumentType{
			Name:     "name",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	return objects.FromInterface(hrm.resTengo.res.GetHeaderValue(key))
}

func (hrm *HttpResponseMethod) getHeaders(args ...objects.Object) (objects.Object, error) {

	if len(args) != 1 {

		return nil, tengo.ErrWrongNumArguments
	}

	key, ok := objects.ToString(args[0])

	if !ok {

		return nil, tengo.ErrInvalidArgumentType{
			Name:     "name",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	return objects.FromInterface(hrm.resTengo.res.GetHeaderValues(key))
}

func (hrm *HttpResponseMethod) Call(args ...objects.Object) (objects.Object, error) {

	switch hrm.Name {
	case getStatusCodeMethod:
		return hrm.getStatusCode()

	case getBodyAsByteMethod:
		return hrm.getBodyAsByte()

	case getBodyAsStringMethod:
		return hrm.getBodyAsString()

	case getProtocolMethod:
		return hrm.getProtocol()

	case getHeaderMethod:
		return hrm.getHeader(args...)

	case getHeadersMethod:
		return hrm.getHeaders(args...)

	default:
		return nil, errors.New("unknown http response method:" + hrm.Name)

	}
}
