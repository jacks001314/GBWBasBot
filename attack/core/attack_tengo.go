package attack

import (
	"common/proto/http"
	"common/proto/tcp"
	stengo "common/scripts/tengo"
	"errors"
	"io/ioutil"
	"strings"

	"github.com/d5/tengo/objects"
	"github.com/d5/tengo/script"
	"github.com/d5/tengo/stdlib"
	"github.com/d5/tengo/v2"
)

var attackTengoMethodMaps = map[string]*AttackTengoScriptMethod{

	publishAttackProcessMethod: &AttackTengoScriptMethod{
		TengoObj: stengo.TengoObj{Name: publishAttackProcessMethod},
	},
}

type AttackTengoScript struct {
	stengo.TengoObj

	task *AttackTask

	config *Config

	/*tengo script instanse Compiled*/
	attackTengo *script.Compiled
}

type AttackTengoScriptMethod struct {
	stengo.TengoObj
	scriptTengo *AttackTengoScript
}

/*compile tengo script*/
func scriptCompile(sdata []byte) (*script.Compiled, error) {

	script := script.New(sdata)

	script.Add(attackScriptUDName, nil)
	script.Add(attackTargetUDName, nil)

	mm := objects.NewModuleMap()

	/*add all stdlibs*/
	builtinMaps := objects.NewModuleMap()
	for name, im := range stdlib.BuiltinModules {
		builtinMaps.AddBuiltinModule(name, im)
	}

	mm.AddMap(builtinMaps)
	mm.Add(AttackModuleName, AttackTengoScript{})
	mm.Add(http.HTTPModuleName, http.HttpTengo{})
	mm.Add(tcp.TCPModName, tcp.TCPTengo{})

	script.SetImports(mm)
	return script.Compile()
}

func LoadTengoScriptFromContent(task *AttackTask, data []byte, config *Config) (*AttackTengoScript, error) {

	com, err := scriptCompile(data)

	if err != nil {

		return nil, err
	}

	return &AttackTengoScript{
		TengoObj:    stengo.TengoObj{Name: config.Name},
		config:      config,
		attackTengo: com,
		task:        task,
	}, nil

}

func LoadTengoScriptFromFile(task *AttackTask, config *Config) (*AttackTengoScript, error) {

	data, err := ioutil.ReadFile(config.FPath)

	if err != nil {
		return nil, err
	}

	return LoadTengoScriptFromContent(task, data, config)
}

func (ats *AttackTengoScript) Accept(target *AttackTarget) bool {

	atype := ats.config.Atype
	app := ats.config.App

	if target.Types != nil && len(target.Types) > 0 {

		if _, ok := target.Types[atype]; !ok {

			return false
		}
	}

	if target.App != "" && !strings.EqualFold(target.App, app) {

		return false
	}

	return true
}

func (ats *AttackTengoScript) Run(target *AttackTarget) error {

	att := NewAttackTargetTengo(target)

	ts := ats.attackTengo.Clone()

	ts.Set(attackTargetUDName, att)
	ts.Set(attackScriptUDName, ats)

	if err := ts.Run(); err != nil {

		return err

	}

	return nil
}

func (ats *AttackTengoScript) Publish(ap *AttackProcess) {

	ats.task.Publish(ap)

}

func (ats *AttackTengoScript) IndexGet(index objects.Object) (value objects.Object, err error) {

	key, ok := objects.ToString(index)

	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "index",
			Expected: "string(compatible)",
			Found:    index.TypeName(),
		}
	}

	if m, ok := attackTengoMethodMaps[key]; ok {

		m.scriptTengo = ats

		return m, nil
	}

	return nil, errors.New("Undefine attack tengo script function:" + key)

}

func (m *AttackTengoScriptMethod) Call(args ...objects.Object) (objects.Object, error) {

	switch m.Name {

	case publishAttackProcessMethod:
		apt := args[0].(*AttackProcessTengo)
		m.scriptTengo.Publish(apt.ap)
		return nil, nil

	default:
		return nil, errors.New("unknown detect script method:" + m.Name)

	}
}
