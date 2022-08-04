package attack

import (
	"errors"

	stengo "common/scripts/tengo"

	"github.com/d5/tengo/objects"
	"github.com/d5/tengo/v2"
)

var attackProcessMethodMaps = map[string]*AttackProcessMethod{

	attackProcessIPMethod: &AttackProcessMethod{
		TengoObj: stengo.TengoObj{Name: attackProcessIPMethod},
		setFunc: func(m *AttackProcessMethod, v interface{}) {
			m.apTengo.ap.IP = v.(string)
		},
		getFunc: func(m *AttackProcessMethod) interface{} {
			return m.apTengo.ap.IP
		},
	},

	attackProcessHostMethod: &AttackProcessMethod{
		TengoObj: stengo.TengoObj{Name: attackProcessHostMethod},
		setFunc: func(m *AttackProcessMethod, v interface{}) {
			m.apTengo.ap.Host = v.(string)
		},
		getFunc: func(m *AttackProcessMethod) interface{} {
			return m.apTengo.ap.Host
		},
	},

	attackProcessPortMethod: &AttackProcessMethod{
		TengoObj: stengo.TengoObj{Name: attackProcessPortMethod},
		setFunc: func(m *AttackProcessMethod, v interface{}) {
			m.apTengo.ap.Port = v.(int)
		},
		getFunc: func(m *AttackProcessMethod) interface{} {
			return m.apTengo.ap.Port
		},
	},
	attackProcessProtoMethod: &AttackProcessMethod{
		TengoObj: stengo.TengoObj{Name: attackProcessProtoMethod},
		setFunc: func(m *AttackProcessMethod, v interface{}) {
			m.apTengo.ap.Proto = v.(string)
		},
		getFunc: func(m *AttackProcessMethod) interface{} {
			return m.apTengo.ap.Proto
		},
	},
	attackProcessAppMethod: &AttackProcessMethod{
		TengoObj: stengo.TengoObj{Name: attackProcessAppMethod},
		setFunc: func(m *AttackProcessMethod, v interface{}) {
			m.apTengo.ap.App = v.(string)
		},
		getFunc: func(m *AttackProcessMethod) interface{} {
			return m.apTengo.ap.App
		},
	},
	attackProcessOsMethod: &AttackProcessMethod{
		TengoObj: stengo.TengoObj{Name: attackProcessOsMethod},
		setFunc: func(m *AttackProcessMethod, v interface{}) {
			m.apTengo.ap.OS = v.(string)
		},
		getFunc: func(m *AttackProcessMethod) interface{} {
			return m.apTengo.ap.OS
		},
	},
	attackProcessVersionMethod: &AttackProcessMethod{
		TengoObj: stengo.TengoObj{Name: attackProcessVersionMethod},
		setFunc: func(m *AttackProcessMethod, v interface{}) {
			m.apTengo.ap.Version = v.(string)
		},
		getFunc: func(m *AttackProcessMethod) interface{} {
			return m.apTengo.ap.Version
		},
	},
	attackProcessIsSSLMethod: &AttackProcessMethod{
		TengoObj: stengo.TengoObj{Name: attackProcessIsSSLMethod},
		setFunc: func(m *AttackProcessMethod, v interface{}) {
			m.apTengo.ap.IsSSL = v.(bool)
		},
		getFunc: func(m *AttackProcessMethod) interface{} {
			return m.apTengo.ap.IsSSL
		},
	},
	attackProcessIdMethod: &AttackProcessMethod{
		TengoObj: stengo.TengoObj{Name: attackProcessIdMethod},
		setFunc: func(m *AttackProcessMethod, v interface{}) {
			m.apTengo.ap.Id = v.(int)
		},
		getFunc: func(m *AttackProcessMethod) interface{} {
			return m.apTengo.ap.Id
		},
	},
	attackProcessLanguageMethod: &AttackProcessMethod{
		TengoObj: stengo.TengoObj{Name: attackProcessLanguageMethod},
		setFunc: func(m *AttackProcessMethod, v interface{}) {
			m.apTengo.ap.Language = v.(string)
		},
		getFunc: func(m *AttackProcessMethod) interface{} {
			return m.apTengo.ap.Language
		},
	},
	attackProcessNameMethod: &AttackProcessMethod{
		TengoObj: stengo.TengoObj{Name: attackProcessNameMethod},
		setFunc: func(m *AttackProcessMethod, v interface{}) {
			m.apTengo.ap.Name = v.(string)
		},
		getFunc: func(m *AttackProcessMethod) interface{} {
			return m.apTengo.ap.Name
		},
	},
	attackProcessCVECodeMethod: &AttackProcessMethod{
		TengoObj: stengo.TengoObj{Name: attackProcessCVECodeMethod},
		setFunc: func(m *AttackProcessMethod, v interface{}) {
			m.apTengo.ap.CVECode = v.(string)
		},
		getFunc: func(m *AttackProcessMethod) interface{} {
			return m.apTengo.ap.CVECode
		},
	},

	attackProcessDescMethod: &AttackProcessMethod{
		TengoObj: stengo.TengoObj{Name: attackProcessDescMethod},
		setFunc: func(m *AttackProcessMethod, v interface{}) {
			m.apTengo.ap.Desc = v.(string)
		},
		getFunc: func(m *AttackProcessMethod) interface{} {
			return m.apTengo.ap.Desc
		},
	},
	attackProcessSuggestMethod: &AttackProcessMethod{
		TengoObj: stengo.TengoObj{Name: attackProcessSuggestMethod},
		setFunc: func(m *AttackProcessMethod, v interface{}) {
			m.apTengo.ap.Suggest = v.(string)
		},
		getFunc: func(m *AttackProcessMethod) interface{} {
			return m.apTengo.ap.Suggest
		},
	},
	attackProcessStatusMethod: &AttackProcessMethod{
		TengoObj: stengo.TengoObj{Name: attackProcessStatusMethod},
		setFunc: func(m *AttackProcessMethod, v interface{}) {
			m.apTengo.ap.Status = v.(int)
		},
		getFunc: func(m *AttackProcessMethod) interface{} {
			return m.apTengo.ap.Status
		},
	},
	attackProcessPayloadMethod: &AttackProcessMethod{
		TengoObj: stengo.TengoObj{Name: attackProcessPayloadMethod},
		setFunc: func(m *AttackProcessMethod, v interface{}) {
			m.apTengo.ap.Payload = v.(string)
		},
		getFunc: func(m *AttackProcessMethod) interface{} {
			return m.apTengo.ap.Payload
		},
	},

	attackProcessTypeMethod: &AttackProcessMethod{
		TengoObj: stengo.TengoObj{Name: attackProcessTypeMethod},
		setFunc: func(m *AttackProcessMethod, v interface{}) {
			m.apTengo.ap.Type = v.(string)
		},
		getFunc: func(m *AttackProcessMethod) interface{} {
			return m.apTengo.ap.Type
		},
	},

	attackProcessResultMethod: &AttackProcessMethod{
		TengoObj: stengo.TengoObj{Name: attackProcessResultMethod},
		setFunc: func(m *AttackProcessMethod, v interface{}) {
			m.apTengo.ap.Result = v.(string)
		},
		getFunc: func(m *AttackProcessMethod) interface{} {
			return m.apTengo.ap.Result
		},
	},

	attackProcessDetailsMehtod: &AttackProcessMethod{
		TengoObj: stengo.TengoObj{Name: attackProcessDetailsMehtod},
		setFunc: func(m *AttackProcessMethod, v interface{}) {
			m.apTengo.ap.Details = v.(string)
		},
		getFunc: func(m *AttackProcessMethod) interface{} {
			return m.apTengo.ap.Details
		},
	},
}

type AttackProcessTengo struct {
	stengo.TengoObj
	ap *AttackProcess
}

type AttackProcessMethod struct {
	stengo.TengoObj
	apTengo *AttackProcessTengo

	setFunc func(m *AttackProcessMethod, v interface{})

	getFunc func(m *AttackProcessMethod) interface{}
}

func newAttackProcessTengo(args ...objects.Object) (objects.Object, error) {

	return &AttackProcessTengo{
		TengoObj: stengo.TengoObj{Name: attackProcessUDName},
		ap:       &AttackProcess{},
	}, nil
}

func (apt *AttackProcessTengo) IndexGet(index objects.Object) (value objects.Object, err error) {

	key, ok := objects.ToString(index)

	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "index",
			Expected: "string(compatible)",
			Found:    index.TypeName(),
		}
	}

	if m, ok := attackProcessMethodMaps[key]; ok {

		m.apTengo = apt

		return m, nil
	}

	return nil, errors.New("Undefine AttackProcess function:" + key)

}

func (m *AttackProcessMethod) Call(args ...objects.Object) (objects.Object, error) {

	if len(args) == 1 {

		//set
		if m.Name == attackProcessPortMethod || m.Name == attackProcessIdMethod || m.Name == attackProcessStatusMethod {

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

		} else if m.Name == attackProcessIsSSLMethod {

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

		return m.apTengo, nil
	}

	return objects.FromInterface(m.getFunc(m))
}
