package detect

import (
	"common/proto/http"
	"common/proto/tcp"
	luahelper "common/scripts/lua"
	"fmt"
	"sync"

	glua "github.com/yuin/gopher-lua"
)

var lpool = sync.Pool{

	New: func() interface{} {

		L := glua.NewState()

		luahelper.RegisterModule(L, http.HTTPModuleName, http.Loader)
		luahelper.RegisterModule(L, tcp.TCPModName, tcp.Loader)
		luahelper.RegisterModule(L, DetectModuleName, Loader)
		luahelper.RegisterGlobalType(L, detectScriptUDName, nil, detectScriptApis)
		luahelper.RegisterGlobalType(L, detectTargetUDName, nil, dtargetApis)

		return L
	},
}

//对应一个lua探测脚本，要实现detect接口
type DLuaScript struct {
	task *DTask

	Key   string              //脚本key,用来唯一的标识一个脚本
	bcode *glua.FunctionProto //lua探测脚本编译字节码

}

//从内存加载编译lua探测脚本
func LoadLuaScriptFromContent(content []byte, key string) (*DLuaScript, error) {

	bcode, err := luahelper.CompileLuaScript(content, key)

	if err != nil {

		return nil, err
	}

	return &DLuaScript{
		Key:   key,
		bcode: bcode,
	}, nil
}

//从本地文件系统加载并编译lua探测脚本
func LoadLuaScriptFromFile(fpath string, key string) (*DLuaScript, error) {

	bcode, err := luahelper.CompileLuaScriptFromFile(fpath)

	if err != nil {

		return nil, err
	}

	return &DLuaScript{
		Key:   key,
		bcode: bcode,
	}, nil

}

//运行lua探测脚本
func (dl *DLuaScript) Run(target *DTarget) error {

	L := lpool.Get().(*glua.LState)
	defer lpool.Put(L)

	fmt.Println(L)

	luahelper.SetGlobalUserData(L, detectScriptUDName, dl)
	luahelper.SetGlobalUserData(L, detectTargetUDName, target)

	if err := luahelper.RunLua(L, dl.bcode); err != nil {

		return err
	}

	return nil
}

func (dl *DLuaScript) Publish(result *DResult) {

	//fmt.Println(result)
}
