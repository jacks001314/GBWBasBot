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

func RegisterApis(L *glua.LState, module *glua.LTable,
	apis map[string]glua.LGFunction,
	fkey, udname string) {

	mt := L.NewTypeMetatable(udname)
	L.SetField(mt, MTIndexName, L.SetFuncs(L.NewTable(), apis))

	L.SetField(module, fkey, mt)

}

func SetUserData(L *glua.LState, name string, udata interface{}) {

	ud := L.NewUserData()
	ud.Value = udata
	L.SetMetatable(ud, L.GetTypeMetatable(name))
	L.Push(ud)
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
