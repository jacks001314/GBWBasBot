package target

import (
	attack "attack/core"
	luahelper "common/scripts/lua"

	glua "github.com/yuin/gopher-lua"
)

var sourceModApis = map[string]glua.LGFunction{

	newTargetMethod: newTargetApi,
}

var sourceScriptApis = map[string]glua.LGFunction{
	putMethod:   putApi,
	atEndMethod: atEndApi,
}

func Loader(L *glua.LState) int {

	mod := L.SetFuncs(L.NewTable(), sourceModApis)

	luahelper.RegisterApis(L, mod, attack.AttackTargetApis, targetApiName, targetUDName)

	L.Push(mod)

	return 1
}

func checkScript(L *glua.LState, n int) *SourceLuaScript {

	ud := L.CheckUserData(n)

	if v, ok := ud.Value.(*SourceLuaScript); ok {
		return v
	}

	L.ArgError(1, "SourceLuaScript expected")
	return nil
}

func checkTarget(L *glua.LState, n int) *attack.AttackTarget {

	ud := L.CheckUserData(n)

	if v, ok := ud.Value.(*attack.AttackTarget); ok {
		return v
	}

	L.ArgError(1, "AttackTarget expected")
	return nil
}

func newTargetApi(L *glua.LState) int {

	at := &attack.AttackTarget{}

	luahelper.SetUserData(L, targetUDName, at)

	return 1
}

func atEndApi(L *glua.LState) int {

	s := checkScript(L, 1)

	s.AtEnd()

	return 0
}

func putApi(L *glua.LState) int {

	s := checkScript(L, 1)
	t := checkTarget(L, 2)

	s.Put(t)

	return 0
}
