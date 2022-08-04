package attack

import (
	"errors"

	"github.com/d5/tengo/objects"
)

var moduleMap objects.Object = &objects.ImmutableMap{

	Value: map[string]objects.Object{

		newAttackProcessMethod: &objects.UserFunction{
			Name:  newAttackProcessMethod,
			Value: newAttackProcessTengo,
		},
	},
}

func (AttackTengoScript) Import(moduleName string) (interface{}, error) {

	switch moduleName {

	case AttackModuleName:
		return moduleMap, nil
	default:
		return nil, errors.New("undefined module:" + moduleName)
	}
}
