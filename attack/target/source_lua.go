package target

import (
	attack "attack/core"
	"common/ipgen/ipv4"
	"common/proto/http"
	"common/proto/tcp"
	luahelper "common/scripts/lua"
	"common/scripts/lua/base64"
	json "common/scripts/lua/json"
	"common/scripts/lua/strings"

	"sync"

	glua "github.com/yuin/gopher-lua"
)

var lpool = sync.Pool{

	New: func() interface{} {

		L := glua.NewState()

		luahelper.RegisterModule(L, http.HTTPModuleName, http.Loader)
		luahelper.RegisterModule(L, tcp.TCPModName, tcp.Loader)
		luahelper.RegisterModule(L, ipv4.IPV4ModName, ipv4.Loader)
		luahelper.RegisterModule(L, json.JsonModuleName, json.Loader)
		luahelper.RegisterModule(L, base64.Base64ModuleName, base64.Loader)
		luahelper.RegisterModule(L, strings.StringsModuleName, strings.Loader)
		luahelper.RegisterModule(L, SourceModuleName, Loader)
		luahelper.RegisterGlobalType(L, SourceUDName, nil, sourceScriptApis)

		return L
	},
}

type SourceLuaScript struct {
	task *SourceTask

	types []string

	bcode *glua.FunctionProto //lua攻击脚本编译字节码
}

//从内存加载编译lua源脚本
func LoadLuaScriptFromContent(task *SourceTask, content []byte, types []string) (*SourceLuaScript, error) {

	bcode, err := luahelper.CompileLuaScript(content, SourceModuleName)

	if err != nil {

		return nil, err
	}

	return &SourceLuaScript{
		task:  task,
		types: types,
		bcode: bcode,
	}, nil
}

//从本地文件系统加载并编译lua源脚本
func LoadLuaScriptFromFile(task *SourceTask, fpath string, types []string) (*SourceLuaScript, error) {

	bcode, err := luahelper.CompileLuaScriptFromFile(fpath)

	if err != nil {

		return nil, err
	}

	return &SourceLuaScript{
		task:  task,
		types: types,
		bcode: bcode,
	}, nil
}

func (sls *SourceLuaScript) Stop() {

}

func (sls *SourceLuaScript) Put(target *attack.AttackTarget) error {

	if len(sls.types) > 0 {

		target.AddTypes(sls.types)
	}

	sls.task.Put(target)

	return nil
}

func (sls *SourceLuaScript) AtEnd() {

	sls.task.CloseSource(sls)
}

func (sls *SourceLuaScript) AttackTypes() []string {

	return sls.types
}

func (sls *SourceLuaScript) Start() error {

	L := lpool.Get().(*glua.LState)
	defer lpool.Put(L)

	luahelper.SetGlobalUserData(L, SourceUDName, sls)

	if err := luahelper.RunLua(L, sls.bcode); err != nil {

		return err
	}

	return nil
}
