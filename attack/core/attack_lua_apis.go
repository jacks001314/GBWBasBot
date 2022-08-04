package attack

import (
	luahelper "common/scripts/lua"

	glua "github.com/yuin/gopher-lua"
)

var attackApis = map[string]glua.LGFunction{
	newAttackProcessMethod: newAttackProcessApi,
}

var attackScriptApis = map[string]glua.LGFunction{

	publishAttackProcessMethod: publishApi,
}

func Loader(L *glua.LState) int {

	mod := L.SetFuncs(L.NewTable(), attackApis)

	luahelper.RegisterApis(L, mod, attackProcessApis, attackProcessApiName, attackProcessUDName)

	L.Push(mod)

	return 1
}

func checkScript(L *glua.LState, n int) *AttackLuaScript {

	ud := L.CheckUserData(n)

	if v, ok := ud.Value.(*AttackLuaScript); ok {
		return v
	}

	L.ArgError(1, "AttackLuaScript expected")
	return nil
}

func publishApi(L *glua.LState) int {

	as := checkScript(L, 1)
	ap := checkAttackProcess(L, 2)

	as.Publish(ap)

	return 0
}
