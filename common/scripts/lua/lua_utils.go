package luahelper

import (
	"bytes"
	"io/ioutil"

	glua "github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/parse"
)

const (
	MTIndexName = "__index"
)

func RegisterModule(L *glua.LState, name string, loader glua.LGFunction) {
	L.PreloadModule(name, loader)
}

func RegisterApis(L *glua.LState, module *glua.LTable,
	apis map[string]glua.LGFunction,
	fkey, udname string) {

	mt := L.NewTypeMetatable(udname)
	L.SetField(mt, MTIndexName, L.SetFuncs(L.NewTable(), apis))

	L.SetField(module, fkey, mt)

}

func RegisterUserData(L *glua.LState, module *glua.LTable, udname string, userdata interface{},
	apis map[string]glua.LGFunction) {

	mt := L.NewTypeMetatable(udname)
	// methods
	L.SetField(mt, MTIndexName, L.SetFuncs(L.NewTable(), apis))
	ud := L.NewUserData()
	ud.Value = userdata
	L.SetMetatable(ud, mt)

	L.SetField(module, udname, ud)
}

func SetUserData(L *glua.LState, name string, udata interface{}) {

	ud := L.NewUserData()
	ud.Value = udata
	L.SetMetatable(ud, L.GetTypeMetatable(name))
	L.Push(ud)
}

func RegisterGlobalType(L *glua.LState, name string, userdata interface{},
	apis map[string]glua.LGFunction) {

	mt := L.NewTypeMetatable(name)
	// methods
	L.SetField(mt, MTIndexName, L.SetFuncs(L.NewTable(), apis))
	ud := L.NewUserData()
	ud.Value = userdata
	L.SetMetatable(ud, mt)
	L.SetGlobal(name, ud)
}

func LuaTableToStringArray(t *glua.LTable) []string {

	arr := make([]string, 0)

	t.ForEach(func(k glua.LValue, v glua.LValue) {

		arr = append(arr, v.String())
	})

	return arr
}

func CompileLuaScript(content []byte, name string) (*glua.FunctionProto, error) {

	reader := bytes.NewReader(content)

	chunk, err := parse.Parse(reader, name)

	if err != nil {

		return nil, err
	}

	proto, err := glua.Compile(chunk, name)

	if err != nil {
		return nil, err
	}

	return proto, nil
}

func CompileLuaScriptFromFile(fpath string) (*glua.FunctionProto, error) {

	content, err := ioutil.ReadFile(fpath)

	if err != nil {

		return nil, err
	}

	return CompileLuaScript(content, fpath)
}

func RunLua(L *glua.LState, bcode *glua.FunctionProto) error {

	lfunc := L.NewFunctionFromProto(bcode)
	L.Push(lfunc)

	return L.PCall(0, glua.MultRet, nil)
}
