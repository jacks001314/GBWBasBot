package detect

import (
	"errors"

	"github.com/d5/tengo/objects"
)

var moduleMap objects.Object = &objects.ImmutableMap{

	Value: map[string]objects.Object{

		newDetectResultMethod: &objects.UserFunction{
			Name:  newDetectResultMethod,
			Value: newDetectResultTengo,
		},
	},
}

func (DetectTengoScript) Import(moduleName string) (interface{}, error) {

	switch moduleName {

	case DetectModuleName:
		return moduleMap, nil
	default:
		return nil, errors.New("undefined module:" + moduleName)
	}
}
