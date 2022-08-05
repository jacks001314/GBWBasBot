package strings

import lua "github.com/yuin/gopher-lua"

const (
	StringsModuleName = "strings"
)

var api = map[string]lua.LGFunction{
	"split":       Split,
	"trim":        Trim,
	"trim_space":  TrimSpace,
	"trim_prefix": TrimPrefix,
	"trim_suffix": TrimSuffix,
	"has_prefix":  HasPrefix,
	"has_suffix":  HasSuffix,
	"parse_int":   ParseInt,
	"contains":    Contains,
	"new_reader":  newStringsReader,
	"new_builder": newStringsBuilder,
	"fields":      Fields,
}

// Loader is the module loader function.
func Loader(L *lua.LState) int {

	registerStringsReader(L)
	registerStringsBuilder(L)

	t := L.NewTable()
	L.SetFuncs(t, api)
	L.Push(t)
	return 1
}
