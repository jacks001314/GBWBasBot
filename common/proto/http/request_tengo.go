package http

import (
	"errors"

	"github.com/d5/tengo/objects"
	"github.com/d5/tengo/v2"

	stengo "common/scripts/tengo"
)

var reqMethodMaps = map[string]*HttpRequestMethod{

	authMethod: &HttpRequestMethod{
		TengoObj: stengo.TengoObj{Name: authMethod},
	},
	addHeaderMethod: &HttpRequestMethod{
		TengoObj: stengo.TengoObj{Name: addHeaderMethod},
	},
	putStringMethod: &HttpRequestMethod{
		TengoObj: stengo.TengoObj{Name: putStringMethod},
	},
	putHexMethod: &HttpRequestMethod{
		TengoObj: stengo.TengoObj{Name: putHexMethod},
	},
	uploadMethod: &HttpRequestMethod{
		TengoObj: stengo.TengoObj{Name: uploadMethod},
	},
}

/*for new http request */
type HttpRequestTengo struct {
	stengo.TengoObj
	req *HttpRequest
}

type HttpRequestMethod struct {
	stengo.TengoObj
	reqTengo *HttpRequestTengo
}

func (req *HttpRequestTengo) IndexGet(index objects.Object) (value objects.Object, err error) {

	key, ok := objects.ToString(index)

	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "index",
			Expected: "string(compatible)",
			Found:    index.TypeName(),
		}
	}

	if m, ok := reqMethodMaps[key]; ok {

		m.reqTengo = req
		return m, nil
	}

	return nil, errors.New("undefine http request method:" + key)

}

func (hrm *HttpRequestMethod) makeRequestAuth(args ...objects.Object) (objects.Object, error) {

	if len(args) != 2 {

		return nil, tengo.ErrWrongNumArguments
	}

	user, ok := objects.ToString(args[0])
	if !ok {

		return nil, tengo.ErrInvalidArgumentType{
			Name:     "user",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	passwd, ok := objects.ToString(args[1])
	if !ok {

		return nil, tengo.ErrInvalidArgumentType{
			Name:     "passwd",
			Expected: "string(compatible)",
			Found:    args[1].TypeName(),
		}
	}

	hrm.reqTengo.req.BasicAuth(user, passwd)

	return hrm.reqTengo, nil
}

func (hrm *HttpRequestMethod) addHeader(args ...objects.Object) (objects.Object, error) {

	if len(args) != 2 {

		return nil, tengo.ErrWrongNumArguments
	}

	name, ok := objects.ToString(args[0])
	if !ok {

		return nil, tengo.ErrInvalidArgumentType{
			Name:     "name",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	value, ok := objects.ToString(args[1])
	if !ok {

		return nil, tengo.ErrInvalidArgumentType{
			Name:     "value",
			Expected: "string(compatible)",
			Found:    args[1].TypeName(),
		}
	}

	hrm.reqTengo.req.AddHeader(name, value)

	return hrm.reqTengo, nil
}

func (hrm *HttpRequestMethod) putString(args ...objects.Object) (objects.Object, error) {

	if len(args) != 2 {

		return nil, tengo.ErrWrongNumArguments
	}

	content, ok := objects.ToString(args[0])
	if !ok {

		return nil, tengo.ErrInvalidArgumentType{
			Name:     "content",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	fromFile, ok := objects.ToBool(args[1])
	if !ok {

		return nil, tengo.ErrInvalidArgumentType{
			Name:     "fromFile",
			Expected: "bool(compatible)",
			Found:    args[1].TypeName(),
		}
	}

	hrm.reqTengo.req.BodyString(content, fromFile)

	return hrm.reqTengo, nil
}

func (hrm *HttpRequestMethod) putHex(args ...objects.Object) (objects.Object, error) {

	if len(args) != 1 {

		return nil, tengo.ErrWrongNumArguments
	}

	content, ok := objects.ToString(args[0])
	if !ok {

		return nil, tengo.ErrInvalidArgumentType{
			Name:     "content",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	hrm.reqTengo.req.BodyHex(content)

	return hrm.reqTengo, nil
}

func (hrm *HttpRequestMethod) upload(args ...objects.Object) (objects.Object, error) {

	if len(args) != 3 {

		return nil, tengo.ErrWrongNumArguments
	}

	name, ok := objects.ToString(args[0])
	if !ok {

		return nil, tengo.ErrInvalidArgumentType{
			Name:     "name",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	fpath, ok := objects.ToString(args[1])
	if !ok {

		return nil, tengo.ErrInvalidArgumentType{
			Name:     "fpath",
			Expected: "string(compatible)",
			Found:    args[1].TypeName(),
		}
	}

	boundary, ok := objects.ToString(args[2])
	if !ok {

		return nil, tengo.ErrInvalidArgumentType{
			Name:     "boundary",
			Expected: "string(compatible)",
			Found:    args[2].TypeName(),
		}
	}

	hrm.reqTengo.req.UPload(name, fpath, boundary)

	return hrm.reqTengo, nil
}

func (hrm *HttpRequestMethod) Call(args ...objects.Object) (objects.Object, error) {

	switch hrm.Name {

	case authMethod:
		return hrm.makeRequestAuth(args...)

	case addHeaderMethod:
		return hrm.addHeader(args...)

	case putStringMethod:
		return hrm.putString(args...)

	case putHexMethod:
		return hrm.putHex(args...)

	case uploadMethod:
		return hrm.upload(args...)

	default:
		return nil, errors.New("unknown http request method:" + hrm.Name)
	}
}
