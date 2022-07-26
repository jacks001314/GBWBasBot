package detect

import (
	luahelper "common/scripts/lua"

	glua "github.com/yuin/gopher-lua"
)

type DetectLuaModule struct {
	script  *DLuaScript
	dtarget *DTarget
}

var detectApis = map[string]glua.LGFunction{
	newDetectResultMethod: newDetectResultApi,
}

var detectScriptApis = map[string]glua.LGFunction{

	publishDetectResultMethod: publishApi,
}

func NewDetectLuaModule(script *DLuaScript, dtarget *DTarget) *DetectLuaModule {

	return &DetectLuaModule{
		script:  script,
		dtarget: dtarget,
	}

}

func (dlm *DetectLuaModule) Loader(L *glua.LState) int {

	mod := L.SetFuncs(L.NewTable(), detectApis)

	luahelper.RegisterApis(L, mod, dresultApis, detectResultApiName, detectResultUDName)

	luahelper.RegisterUserData(L, mod, detectScriptUDName, dlm.script, detectScriptApis)
	luahelper.RegisterUserData(L, mod, detectTargetUDName, dlm.dtarget, dtargetApis)

	L.Push(mod)

	return 1
}

func checkScript(L *glua.LState, n int) *DLuaScript {

	ud := L.CheckUserData(n)

	if v, ok := ud.Value.(*DLuaScript); ok {
		return v
	}

	L.ArgError(1, "DLuaScript expected")
	return nil
}

func publishApi(L *glua.LState) int {

	ds := checkScript(L, 1)
	dr := checkDResult(L, 2)

	ds.Publish(dr)

	return 0
}
