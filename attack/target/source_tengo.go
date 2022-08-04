package target

import (
	attack "attack/core"
	"common/ipgen/ipv4"
	"common/proto/http"
	"common/proto/tcp"
	stengo "common/scripts/tengo"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/d5/tengo/objects"
	"github.com/d5/tengo/script"
	"github.com/d5/tengo/stdlib"
	"github.com/d5/tengo/v2"
)

var sourceTengoMethodMaps = map[string]*SourceTengoScriptMethod{

	putMethod: &SourceTengoScriptMethod{
		TengoObj: stengo.TengoObj{Name: putMethod},
	},

	atEndMethod: &SourceTengoScriptMethod{
		TengoObj: stengo.TengoObj{Name: atEndMethod},
	},
}

type SourceTengoScript struct {
	stengo.TengoObj

	task *SourceTask

	types []string

	/*tengo script instanse Compiled*/
	sourceTengo *script.Compiled
}

type SourceTengoScriptMethod struct {
	stengo.TengoObj
	scriptTengo *SourceTengoScript
}

/*compile tengo script*/
func scriptCompile(sdata []byte) (*script.Compiled, error) {

	script := script.New(sdata)

	script.Add(SourceUDName, nil)

	mm := objects.NewModuleMap()

	/*add all stdlibs*/
	builtinMaps := objects.NewModuleMap()
	for name, im := range stdlib.BuiltinModules {
		builtinMaps.AddBuiltinModule(name, im)
	}

	mm.AddMap(builtinMaps)
	mm.Add(SourceModuleName, SourceTengoScript{})
	mm.Add(http.HTTPModuleName, http.HttpTengo{})
	mm.Add(tcp.TCPModName, tcp.TCPTengo{})
	mm.Add(ipv4.IPV4ModName, ipv4.IPV4Tengo{})

	script.SetImports(mm)

	return script.Compile()
}

func LoadTengoScriptFromContent(task *SourceTask, data []byte, types []string) (*SourceTengoScript, error) {

	com, err := scriptCompile(data)

	if err != nil {

		return nil, err
	}

	return &SourceTengoScript{
		TengoObj:    stengo.TengoObj{Name: SourceModuleName},
		types:       types,
		task:        task,
		sourceTengo: com,
	}, nil

}

func LoadTengoScriptFromFile(task *SourceTask, fpath string, types []string) (*SourceTengoScript, error) {

	data, err := ioutil.ReadFile(fpath)

	if err != nil {
		return nil, err
	}

	return LoadTengoScriptFromContent(task, data, types)
}

func (sts *SourceTengoScript) Start() error {

	ts := sts.sourceTengo.Clone()

	ts.Set(SourceUDName, sts)

	return ts.Run()

}

func (sts *SourceTengoScript) Stop() {

}

func (sts *SourceTengoScript) Put(target *attack.AttackTarget) error {

	if len(sts.types) > 0 {

		target.AddTypes(sts.types)
	}

	sts.task.Put(target)
	return nil
}

func (sts *SourceTengoScript) AtEnd() {

	sts.task.CloseSource(sts)
}

func (sts *SourceTengoScript) AttackTypes() []string {

	return sts.types
}

func (sts *SourceTengoScript) IndexGet(index objects.Object) (value objects.Object, err error) {

	key, ok := objects.ToString(index)

	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "index",
			Expected: "string(compatible)",
			Found:    index.TypeName(),
		}
	}

	if m, ok := sourceTengoMethodMaps[key]; ok {

		m.scriptTengo = sts

		return m, nil
	}

	return nil, errors.New("Undefine Source tengo script function:" + key)

}

func (m *SourceTengoScriptMethod) Call(args ...objects.Object) (objects.Object, error) {

	switch m.Name {

	case putMethod:
		if len(args) != 1 {

			return nil, tengo.ErrWrongNumArguments
		}

		target := args[0].(*attack.AttackTargetTengo)

		return nil, m.scriptTengo.Put(target.GetAttackTarget())

	case atEndMethod:

		m.scriptTengo.AtEnd()
		return nil, nil

	default:
		return nil, errors.New("unknown source script method:" + m.Name)

	}
}

func newAttackTarget(args ...objects.Object) (objects.Object, error) {

	return attack.NewAttackTargetTengo(&attack.AttackTarget{}), nil
}

var moduleMap objects.Object = &objects.ImmutableMap{

	Value: map[string]objects.Object{

		newTargetMethod: &objects.UserFunction{
			Name:  newTargetMethod,
			Value: newAttackTarget,
		},
	},
}

func (SourceTengoScript) Import(moduleName string) (interface{}, error) {

	switch moduleName {
	case SourceModuleName:
		return moduleMap, nil
	default:
		return nil, fmt.Errorf("undefined module %q", moduleName)
	}
}
