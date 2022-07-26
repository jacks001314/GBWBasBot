package detect

import (
	luahelper "common/scripts/lua"

	glua "github.com/yuin/gopher-lua"
)

var dresultApis = map[string]glua.LGFunction{

	dresultIPMethod:      dresultIP,
	dresultPortMethod:    dresultPort,
	dresultAppMethod:     dresultApp,
	dresultVersionMethod: dresultVersion,
	dresultProtoMethod:   dresultProto,
	dresultIsSSLMethod:   dresultIsSSL,
}

func newDetectResultApi(L *glua.LState) int {

	dresult := &DResult{}

	luahelper.SetUserData(L, detectResultUDName, dresult)

	return 1
}

func checkDResult(L *glua.LState, n int) *DResult {

	ud := L.CheckUserData(n)

	if v, ok := ud.Value.(*DResult); ok {
		return v
	}

	L.ArgError(1, "DResult expected")
	return nil
}

func dresultIP(L *glua.LState) int {

	dr := checkDResult(L, 1)

	if L.GetTop() == 2 {

		v := L.CheckString(2)
		dr.IP = v
		return 0
	}

	L.Push(glua.LString(dr.IP))

	return 1
}

func dresultPort(L *glua.LState) int {

	dr := checkDResult(L, 1)

	if L.GetTop() == 2 {

		v := L.CheckInt(2)
		dr.Port = uint16(v)
		return 0
	}

	L.Push(glua.LNumber(dr.Port))

	return 1
}

func dresultApp(L *glua.LState) int {

	dr := checkDResult(L, 1)

	if L.GetTop() == 2 {

		v := L.CheckString(2)
		dr.App = v
		return 0
	}

	L.Push(glua.LString(dr.App))

	return 1
}

func dresultVersion(L *glua.LState) int {

	dr := checkDResult(L, 1)

	if L.GetTop() == 2 {

		v := L.CheckString(2)
		dr.Version = v
		return 0
	}

	L.Push(glua.LString(dr.Version))

	return 1
}

func dresultProto(L *glua.LState) int {

	dr := checkDResult(L, 1)

	if L.GetTop() == 2 {

		v := L.CheckString(2)
		dr.Proto = v
		return 0
	}

	L.Push(glua.LString(dr.Proto))

	return 1
}

func dresultIsSSL(L *glua.LState) int {

	dr := checkDResult(L, 1)

	if L.GetTop() == 2 {

		v := L.CheckBool(2)
		dr.IsSSL = v
		return 0
	}

	L.Push(glua.LBool(dr.IsSSL))

	return 1
}
