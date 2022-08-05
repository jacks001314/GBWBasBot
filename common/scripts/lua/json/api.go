package json

import lua "github.com/yuin/gopher-lua"

const (
	JsonModuleName = "json"
	decodeApi      = "decode"
	encodeApi      = "encode"
)

var api = map[string]lua.LGFunction{
	decodeApi: apiDecode,
	encodeApi: apiEncode,
}

// Loader is the module loader function.
func Loader(L *lua.LState) int {

	t := L.NewTable()
	L.SetFuncs(t, api)
	L.Push(t)
	return 1
}
