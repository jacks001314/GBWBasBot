package attack

import (
	stengo "common/scripts/tengo"
	"errors"

	"github.com/d5/tengo/objects"
	"github.com/d5/tengo/v2"
)

var attackTargetMethodMap = map[string]*AttackTargetMethod{

	attackTargetIPMethod: &AttackTargetMethod{
		TengoObj: stengo.TengoObj{Name: attackTargetIPMethod},
		setFunc: func(m *AttackTargetMethod, v interface{}) {
			m.atTengo.at.IP = v.(string)
		},
		getFunc: func(m *AttackTargetMethod) interface{} {
			return m.atTengo.at.IP
		},
	},

	attackTargetHostMethod: &AttackTargetMethod{
		TengoObj: stengo.TengoObj{Name: attackTargetHostMethod},
		setFunc: func(m *AttackTargetMethod, v interface{}) {
			m.atTengo.at.Host = v.(string)
		},
		getFunc: func(m *AttackTargetMethod) interface{} {
			return m.atTengo.at.Host
		},
	},

	attackTargetPortMethod: &AttackTargetMethod{
		TengoObj: stengo.TengoObj{Name: attackTargetPortMethod},
		setFunc: func(m *AttackTargetMethod, v interface{}) {
			m.atTengo.at.Port = v.(int)
		},
		getFunc: func(m *AttackTargetMethod) interface{} {
			return m.atTengo.at.Port
		},
	},

	attackTargetAppMethod: &AttackTargetMethod{
		TengoObj: stengo.TengoObj{Name: attackTargetAppMethod},
		setFunc: func(m *AttackTargetMethod, v interface{}) {
			m.atTengo.at.App = v.(string)
		},
		getFunc: func(m *AttackTargetMethod) interface{} {
			return m.atTengo.at.App
		},
	},

	attackTargetVersionMethod: &AttackTargetMethod{
		TengoObj: stengo.TengoObj{Name: attackTargetVersionMethod},
		setFunc: func(m *AttackTargetMethod, v interface{}) {
			m.atTengo.at.Version = v.(string)
		},
		getFunc: func(m *AttackTargetMethod) interface{} {
			return m.atTengo.at.Version
		},
	},
	attackTargetProtoMethod: &AttackTargetMethod{
		TengoObj: stengo.TengoObj{Name: attackTargetProtoMethod},
		setFunc: func(m *AttackTargetMethod, v interface{}) {
			m.atTengo.at.Proto = v.(string)
		},
		getFunc: func(m *AttackTargetMethod) interface{} {
			return m.atTengo.at.Proto
		},
	},
	attackTargetIsSSLMethod: &AttackTargetMethod{
		TengoObj: stengo.TengoObj{Name: attackTargetIsSSLMethod},
		setFunc: func(m *AttackTargetMethod, v interface{}) {
			m.atTengo.at.IsSSL = v.(bool)
		},
		getFunc: func(m *AttackTargetMethod) interface{} {
			return m.atTengo.at.IsSSL
		},
	},
}

type AttackTargetTengo struct {
	stengo.TengoObj
	at *AttackTarget
}

type AttackTargetMethod struct {
	stengo.TengoObj
	atTengo *AttackTargetTengo

	setFunc func(m *AttackTargetMethod, v interface{})

	getFunc func(m *AttackTargetMethod) interface{}
}

func NewAttackTargetTengo(at *AttackTarget) *AttackTargetTengo {

	return &AttackTargetTengo{
		TengoObj: stengo.TengoObj{Name: attackTargetUDName},
		at:       at,
	}
}

func (att *AttackTargetTengo) GetAttackTarget() *AttackTarget {

	return att.at
}

func (att *AttackTargetTengo) IndexGet(index objects.Object) (value objects.Object, err error) {

	key, ok := objects.ToString(index)

	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "index",
			Expected: "string(compatible)",
			Found:    index.TypeName(),
		}
	}

	if m, ok := attackTargetMethodMap[key]; ok {

		m.atTengo = att

		return m, nil
	}

	return nil, errors.New("Undefine AttackTarget function:" + key)

}

func (m *AttackTargetMethod) Call(args ...objects.Object) (objects.Object, error) {

	if len(args) == 1 {

		//set
		if m.Name == attackTargetPortMethod {

			//set
			v, ok := objects.ToInt(args[0])
			if !ok {
				return nil, tengo.ErrInvalidArgumentType{
					Name:     m.Name,
					Expected: "int(compatible)",
					Found:    args[0].TypeName(),
				}
			}

			m.setFunc(m, v)

		} else if m.Name == attackTargetIsSSLMethod {

			//set
			v, ok := objects.ToBool(args[0])
			if !ok {
				return nil, tengo.ErrInvalidArgumentType{
					Name:     m.Name,
					Expected: "bool(compatible)",
					Found:    args[0].TypeName(),
				}
			}

			m.setFunc(m, v)
		} else {

			v, ok := objects.ToString(args[0])
			if !ok {
				return nil, tengo.ErrInvalidArgumentType{
					Name:     m.Name,
					Expected: "string(compatible)",
					Found:    args[0].TypeName(),
				}
			}

			m.setFunc(m, v)
		}

		return m.atTengo, nil
	}

	return objects.FromInterface(m.getFunc(m))
}
