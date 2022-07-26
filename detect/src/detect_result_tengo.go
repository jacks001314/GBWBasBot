package detect

import (
	"errors"

	stengo "common/scripts/tengo"

	"github.com/d5/tengo/objects"
	"github.com/d5/tengo/v2"
)

var drMethodMaps = map[string]*DetectResultMethod{

	dresultIPMethod: &DetectResultMethod{
		TengoObj: stengo.TengoObj{Name: dresultIPMethod},
	},
	dresultPortMethod: &DetectResultMethod{
		TengoObj: stengo.TengoObj{Name: dresultPortMethod},
	},
	dresultAppMethod: &DetectResultMethod{
		TengoObj: stengo.TengoObj{Name: dresultAppMethod},
	},
	dresultVersionMethod: &DetectResultMethod{
		TengoObj: stengo.TengoObj{Name: dresultVersionMethod},
	},
	dresultProtoMethod: &DetectResultMethod{
		TengoObj: stengo.TengoObj{Name: dresultProtoMethod},
	},
	dresultIsSSLMethod: &DetectResultMethod{
		TengoObj: stengo.TengoObj{Name: dresultIsSSLMethod},
	},
}

type DetectResultTengo struct {
	stengo.TengoObj
	dr *DResult
}

type DetectResultMethod struct {
	stengo.TengoObj
	drTengo *DetectResultTengo
}

func newDetectResultTengo(args ...objects.Object) (objects.Object, error) {

	return &DetectResultTengo{
		TengoObj: stengo.TengoObj{Name: "detect.result"},
		dr:       &DResult{},
	}, nil
}

func (dr *DetectResultTengo) IndexGet(index objects.Object) (value objects.Object, err error) {

	key, ok := objects.ToString(index)

	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "index",
			Expected: "string(compatible)",
			Found:    index.TypeName(),
		}
	}

	if m, ok := drMethodMaps[key]; ok {

		m.drTengo = dr

		return m, nil
	}

	return nil, errors.New("Undefine detect result function:" + key)

}

func (m *DetectResultMethod) drIP(args ...objects.Object) (objects.Object, error) {

	if len(args) == 1 {

		//set
		v, ok := objects.ToString(args[0])
		if !ok {
			return nil, tengo.ErrInvalidArgumentType{
				Name:     dresultIPMethod,
				Expected: "string(compatible)",
				Found:    args[0].TypeName(),
			}
		}

		m.drTengo.dr.IP = v

		return m.drTengo, nil
	}

	return objects.FromInterface(m.drTengo.dr.IP)
}

func (m *DetectResultMethod) drPort(args ...objects.Object) (objects.Object, error) {

	if len(args) == 1 {

		//set
		v, ok := objects.ToInt(args[0])
		if !ok {
			return nil, tengo.ErrInvalidArgumentType{
				Name:     dresultPortMethod,
				Expected: "int(compatible)",
				Found:    args[0].TypeName(),
			}
		}

		m.drTengo.dr.Port = uint16(v)

		return m.drTengo, nil
	}

	return objects.FromInterface(m.drTengo.dr.Port)
}

func (m *DetectResultMethod) drApp(args ...objects.Object) (objects.Object, error) {

	if len(args) == 1 {

		//set
		v, ok := objects.ToString(args[0])
		if !ok {
			return nil, tengo.ErrInvalidArgumentType{
				Name:     dresultAppMethod,
				Expected: "string(compatible)",
				Found:    args[0].TypeName(),
			}
		}

		m.drTengo.dr.App = v

		return m.drTengo, nil
	}

	return objects.FromInterface(m.drTengo.dr.App)
}

func (m *DetectResultMethod) drVersion(args ...objects.Object) (objects.Object, error) {

	if len(args) == 1 {

		//set
		v, ok := objects.ToString(args[0])
		if !ok {
			return nil, tengo.ErrInvalidArgumentType{
				Name:     dresultVersionMethod,
				Expected: "string(compatible)",
				Found:    args[0].TypeName(),
			}
		}

		m.drTengo.dr.Version = v

		return m.drTengo, nil
	}

	return objects.FromInterface(m.drTengo.dr.Version)
}

func (m *DetectResultMethod) drProto(args ...objects.Object) (objects.Object, error) {

	if len(args) == 1 {

		//set
		v, ok := objects.ToString(args[0])
		if !ok {
			return nil, tengo.ErrInvalidArgumentType{
				Name:     dresultProtoMethod,
				Expected: "string(compatible)",
				Found:    args[0].TypeName(),
			}
		}

		m.drTengo.dr.Proto = v

		return m.drTengo, nil
	}

	return objects.FromInterface(m.drTengo.dr.Proto)
}

func (m *DetectResultMethod) drIsSSL(args ...objects.Object) (objects.Object, error) {

	if len(args) == 1 {

		//set
		v, ok := objects.ToBool(args[0])
		if !ok {
			return nil, tengo.ErrInvalidArgumentType{
				Name:     dresultIsSSLMethod,
				Expected: "bool(compatible)",
				Found:    args[0].TypeName(),
			}
		}

		m.drTengo.dr.IsSSL = v

		return m.drTengo, nil
	}

	return objects.FromInterface(m.drTengo.dr.IsSSL)
}

func (m *DetectResultMethod) Call(args ...objects.Object) (objects.Object, error) {

	switch m.Name {

	case dresultIPMethod:
		return m.drIP(args...)

	case dresultPortMethod:
		return m.drPort(args...)

	case dresultAppMethod:
		return m.drApp(args...)

	case dresultVersionMethod:
		return m.drVersion(args...)

	case dresultProtoMethod:
		return m.drProto(args...)

	case dresultIsSSLMethod:
		return m.drIsSSL(args...)

	default:
		return nil, errors.New("unknown detect result method:" + m.Name)

	}
}
