package ipv4

import (
	luahelper "common/scripts/lua"

	glua "github.com/yuin/gopher-lua"
)

func Loader(L *glua.LState) int {

	mod := L.SetFuncs(L.NewTable(), api)

	luahelper.RegisterApis(L, mod, ipv4genApis, ipv4genApiName, ipv4genUDName)

	L.Push(mod)
	return 1
}

var api = map[string]glua.LGFunction{

	newIPGenFromFile:  newIPV4GenFromFile,
	newIPGenFromArray: newIPV4GenFromArray,
}
