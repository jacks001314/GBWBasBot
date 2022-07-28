package detect

import (
	"common/proto/http"
	"common/proto/tcp"
	stengo "common/scripts/tengo"
	"errors"
	"io/ioutil"

	"github.com/d5/tengo/v2"

	"github.com/d5/tengo/objects"
	"github.com/d5/tengo/script"
	"github.com/d5/tengo/stdlib"
)

var dtsMethodMaps = map[string]*DetectTengoScriptMethod{

	publishDetectResultMethod: &DetectTengoScriptMethod{
		TengoObj: stengo.TengoObj{Name: publishDetectResultMethod},
	},
}

type DetectTengoScript struct {
	stengo.TengoObj

	task *DTask

	Key string

	/*tengo script instanse Compiled*/
	detectTengo *script.Compiled
}

type DetectTengoScriptMethod struct {
	stengo.TengoObj
	scriptTengo *DetectTengoScript
}

/*compile tengo script*/
func scriptCompile(sdata []byte) (*script.Compiled, error) {

	script := script.New(sdata)

	script.Add(detectScriptUDName, nil)
	script.Add(detectTargetUDName, nil)

	mm := objects.NewModuleMap()

	/*add all stdlibs*/
	builtinMaps := objects.NewModuleMap()
	for name, im := range stdlib.BuiltinModules {
		builtinMaps.AddBuiltinModule(name, im)
	}

	mm.AddMap(builtinMaps)
	mm.Add(DetectModuleName, DetectTengoScript{})
	mm.Add(http.HTTPModuleName, http.HttpTengo{})
	mm.Add(tcp.TCPModName, tcp.TCPTengo{})

	script.SetImports(mm)
	return script.Compile()
}

func LoadTengoScriptFromContent(task *DTask, data []byte, key string) (*DetectTengoScript, error) {

	com, err := scriptCompile(data)

	if err != nil {

		return nil, err
	}

	return &DetectTengoScript{
		TengoObj:    stengo.TengoObj{Name: key},
		Key:         key,
		detectTengo: com,
		task:        task,
	}, nil

}

func LoadTengoScriptFromFile(task *DTask, fpath, key string) (*DetectTengoScript, error) {

	data, err := ioutil.ReadFile(fpath)

	if err != nil {
		return nil, err
	}

	return LoadTengoScriptFromContent(task, data, key)

}

func (dt *DetectTengoScript) Run(target *DTarget) error {

	dtarget := newDetectTargetTengo(target)

	ts := dt.detectTengo.Clone()

	ts.Set(detectTargetUDName, dtarget)
	ts.Set(detectScriptUDName, dt)

	if err := ts.Run(); err != nil {

		return err

	}

	return nil
}

func (dt *DetectTengoScript) Publish(result *DResult) {

	dt.task.Publish(result)

}

func (dt *DetectTengoScript) IndexGet(index objects.Object) (value objects.Object, err error) {

	key, ok := objects.ToString(index)

	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "index",
			Expected: "string(compatible)",
			Found:    index.TypeName(),
		}
	}

	if m, ok := dtsMethodMaps[key]; ok {

		m.scriptTengo = dt

		return m, nil
	}

	return nil, errors.New("Undefine detect script function:" + key)

}

func (m *DetectTengoScriptMethod) Call(args ...objects.Object) (objects.Object, error) {

	switch m.Name {

	case publishDetectResultMethod:
		dr := args[0].(*DetectResultTengo)
		m.scriptTengo.Publish(dr.dr)
		return nil, nil

	default:
		return nil, errors.New("unknown detect script method:" + m.Name)

	}
}
