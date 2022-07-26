package detect

import (
	stengo "common/scripts/tengo"
	"errors"

	"github.com/d5/tengo/objects"
	"github.com/d5/tengo/v2"
)

var dtMethodMaps = map[string]*DetectTargetMethod{

	dtargetIPMethod: &DetectTargetMethod{
		TengoObj: stengo.TengoObj{Name: dtargetIPMethod},
	},
	dtargetPortMethod: &DetectTargetMethod{
		TengoObj: stengo.TengoObj{Name: dtargetPortMethod},
	},
}

type DetectTargetTengo struct {
	stengo.TengoObj
	dt *DTarget
}

type DetectTargetMethod struct {
	stengo.TengoObj
	dtTengo *DetectTargetTengo
}

func newDetectTargetTengo(dt *DTarget) *DetectTargetTengo {

	return &DetectTargetTengo{
		TengoObj: stengo.TengoObj{Name: "dtarget"},
		dt:       dt,
	}

}

func (dt *DetectTargetTengo) IndexGet(index objects.Object) (value objects.Object, err error) {

	key, ok := objects.ToString(index)

	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "index",
			Expected: "string(compatible)",
			Found:    index.TypeName(),
		}
	}

	if m, ok := dtMethodMaps[key]; ok {

		m.dtTengo = dt

		return m, nil
	}

	return nil, errors.New("Undefine detect target function:" + key)

}

func (m *DetectTargetMethod) dtIP(args ...objects.Object) (objects.Object, error) {

	if len(args) == 1 {

		//set
		v, ok := objects.ToString(args[0])
		if !ok {
			return nil, tengo.ErrInvalidArgumentType{
				Name:     dtargetIPMethod,
				Expected: "string(compatible)",
				Found:    args[0].TypeName(),
			}
		}

		m.dtTengo.dt.IP = v

		return m.dtTengo, nil
	}

	return objects.FromInterface(m.dtTengo.dt.IP)
}

func (m *DetectTargetMethod) dtPort(args ...objects.Object) (objects.Object, error) {

	if len(args) == 1 {

		//set
		v, ok := objects.ToInt(args[0])
		if !ok {
			return nil, tengo.ErrInvalidArgumentType{
				Name:     dtargetPortMethod,
				Expected: "int(compatible)",
				Found:    args[0].TypeName(),
			}
		}

		m.dtTengo.dt.Port = uint16(v)

		return m.dtTengo, nil
	}

	return objects.FromInterface(m.dtTengo.dt.Port)
}

func (m *DetectTargetMethod) Call(args ...objects.Object) (objects.Object, error) {

	switch m.Name {

	case dtargetIPMethod:
		return m.dtIP(args...)

	case dtargetPortMethod:
		return m.dtPort(args...)

	default:
		return nil, errors.New("unknown detect target method:" + m.Name)

	}
}
