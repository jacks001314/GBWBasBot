package source

import (
	"common/proto/http"
	"common/proto/tcp"
	luahelper "common/scripts/lua"
	"strings"
	"sync"

	glua "github.com/yuin/gopher-lua"
)

var lpool = sync.Pool{

	New: func() interface{} {

		L := glua.NewState()

		luahelper.RegisterModule(L, http.HTTPModuleName, http.Loader)
		luahelper.RegisterModule(L, tcp.TCPModName, tcp.Loader)
		luahelper.RegisterModule(L, AttackModuleName, Loader)
		luahelper.RegisterGlobalType(L, attackScriptUDName, nil, attackScriptApis)
		luahelper.RegisterGlobalType(L, attackTargetUDName, nil, attackTargetApis)

		return L
	},
}

type AttackLuaScript struct {
	task *AttackTask

	config *Config

	bcode *glua.FunctionProto //lua攻击脚本编译字节码
}

//从内存加载编译lua攻击脚本
func LoadLuaScriptFromContent(task *AttackTask, content []byte, config *Config) (*AttackLuaScript, error) {

	bcode, err := luahelper.CompileLuaScript(content, config.Name)

	if err != nil {

		return nil, err
	}

	return &AttackLuaScript{
		task:   task,
		bcode:  bcode,
		config: config,
	}, nil
}

//从本地文件系统加载并编译lua攻击脚本
func LoadLuaScriptFromFile(task *AttackTask, config *Config) (*AttackLuaScript, error) {

	bcode, err := luahelper.CompileLuaScriptFromFile(config.FPath)

	if err != nil {

		return nil, err
	}

	return &AttackLuaScript{
		task:   task,
		bcode:  bcode,
		config: config,
	}, nil

}

func (als *AttackLuaScript) Accept(target *AttackTarget) bool {

	atype := als.config.Atype
	app := als.config.App

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

func (als *AttackLuaScript) Run(target *AttackTarget) error {

	L := lpool.Get().(*glua.LState)
	defer lpool.Put(L)

	luahelper.SetGlobalUserData(L, attackScriptUDName, als)
	luahelper.SetGlobalUserData(L, attackTargetUDName, target)

	if err := luahelper.RunLua(L, als.bcode); err != nil {

		return err
	}

	return nil
}

func (als *AttackLuaScript) Publish(process *AttackProcess) {

	als.task.Publish(process)
}
