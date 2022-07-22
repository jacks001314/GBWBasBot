package tcp

import (
	luahelper "common/scripts/lua"

	glua "github.com/yuin/gopher-lua"
)

func Loader(L *glua.LState) int {

	mod := L.SetFuncs(L.NewTable(), tcpApis)

	//register connection apis
	luahelper.RegisterApis(L, mod, connectionApis, tcpConnectionApiName, tcpConnectionUDName)

	L.Push(mod)

	return 1
}

var tcpApis = map[string]glua.LGFunction{
	newConnectionMethod: newConnectionApi,
}
