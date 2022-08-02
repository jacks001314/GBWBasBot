package attack

import glua "github.com/yuin/gopher-lua"

var attackTargetApis = map[string]glua.LGFunction{

	attackTargetIPMethod:      attackTargetIP,
	attackTargetPortMethod:    attackTargetPort,
	attackTargetHostMethod:    attackTargetHost,
	attackTargetAppMethod:     atttackTargetApp,
	attackTargetVersionMethod: attackTargetVersion,
	attackTargetProtoMethod:   attackTargetProto,
	attackTargetIsSSLMethod:   attackTargetIsSSL,
}

func checkAttaclTarget(L *glua.LState, n int) *AttackTarget {

	ud := L.CheckUserData(n)

	if v, ok := ud.Value.(*AttackTarget); ok {
		return v
	}

	L.ArgError(1, "AttackTarget expected")
	return nil
}

func attackTargetIP(L *glua.LState) int {

	at := checkAttaclTarget(L, 1)

	if L.GetTop() == 2 {

		v := L.CheckString(2)
		at.IP = v
		return 0
	}

	L.Push(glua.LString(at.IP))

	return 1
}

func attackTargetHost(L *glua.LState) int {

	at := checkAttaclTarget(L, 1)

	if L.GetTop() == 2 {

		v := L.CheckString(2)
		at.Host = v
		return 0
	}

	L.Push(glua.LString(at.Host))

	return 1
}

func attackTargetPort(L *glua.LState) int {

	at := checkAttaclTarget(L, 1)

	if L.GetTop() == 2 {

		v := L.CheckInt(2)
		at.Port = v
		return 0
	}

	L.Push(glua.LNumber(at.Port))

	return 1
}

func atttackTargetApp(L *glua.LState) int {

	at := checkAttaclTarget(L, 1)

	if L.GetTop() == 2 {

		v := L.CheckString(2)
		at.App = v
		return 0
	}

	L.Push(glua.LString(at.App))

	return 1
}

func attackTargetProto(L *glua.LState) int {

	at := checkAttaclTarget(L, 1)

	if L.GetTop() == 2 {

		v := L.CheckString(2)
		at.Proto = v
		return 0
	}

	L.Push(glua.LString(at.Proto))

	return 1
}

func attackTargetVersion(L *glua.LState) int {

	at := checkAttaclTarget(L, 1)

	if L.GetTop() == 2 {

		v := L.CheckString(2)
		at.Version = v
		return 0
	}

	L.Push(glua.LString(at.Version))

	return 1
}

func attackTargetIsSSL(L *glua.LState) int {

	at := checkAttaclTarget(L, 1)

	if L.GetTop() == 2 {

		v := L.CheckBool(2)
		at.IsSSL = v
		return 0
	}

	L.Push(glua.LBool(at.IsSSL))

	return 1
}
