package detect

import (
	"common/proto/http"
	"common/proto/tcp"
	"common/scripts/tengo"
	"io/ioutil"

	"github.com/d5/tengo/objects"
	"github.com/d5/tengo/script"
	"github.com/d5/tengo/stdlib"
)

type DetectTengoScript struct {
	stengo tengo.TengoObj

	Key string

	/*tengo script instanse Compiled*/
	detectTengo *script.Compiled
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

func LoadTengoScriptFromContent(data []byte, key string) (*DetectTengoScript, error) {

	com, err := scriptCompile(data)

	if err != nil {

		return nil, err
	}

	return &DetectTengoScript{
		stengo:      tengo.TengoObj{Name: key},
		Key:         key,
		detectTengo: com,
	}, nil

}

func LoadTengoScriptFromFile(fpath, key string) (*DetectTengoScript, error) {

	data, err := ioutil.ReadFile(fpath)

	if err != nil {
		return nil, err
	}

	return LoadTengoScriptFromContent(data, key)

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

}
