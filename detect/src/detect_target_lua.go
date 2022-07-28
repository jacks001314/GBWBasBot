package detect

import (
	glua "github.com/yuin/gopher-lua"
)

var dtargetApis = map[string]glua.LGFunction{

	"ip":   dtargetIP,
	"port": dtargetPort,
}

func checkDTarget(L *glua.LState, n int) *DTarget {

	ud := L.CheckUserData(n)

	if v, ok := ud.Value.(*DTarget); ok {
		return v
	}

	L.ArgError(1, "DTarget expected")
	return nil
}

func dtargetIP(L *glua.LState) int {

	dt := checkDTarget(L, 1)

	if L.GetTop() == 2 {

		v := L.CheckString(2)
		dt.IP = v
		return 0
	}

	L.Push(glua.LString(dt.IP))

	return 1
}

func dtargetPort(L *glua.LState) int {

	dt := checkDTarget(L, 1)

	if L.GetTop() == 2 {

		v := L.CheckInt(2)
		dt.Port = v
		return 0
	}

	L.Push(glua.LNumber(dt.Port))

	return 1
}
