package attack

import (
	luahelper "common/scripts/lua"

	glua "github.com/yuin/gopher-lua"
)

var attackProcessApis = map[string]glua.LGFunction{

	attackProcessIPMethod:       attackProcessIP,
	attackProcessHostMethod:     attackProcessHost,
	attackProcessPortMethod:     attackProcessPort,
	attackProcessProtoMethod:    attackProcessProto,
	attackProcessAppMethod:      attackProcessApp,
	attackProcessOsMethod:       attackProcessOs,
	attackProcessVersionMethod:  attackProcessVersion,
	attackProcessIsSSLMethod:    attackProcessIsSSL,
	attackProcessIdMethod:       attackProcessId,
	attackProcessLanguageMethod: attackProcessLanguage,
	attackProcessNameMethod:     attackProcessName,
	attackProcessTypeMethod:     attackProcessType,
	attackProcessCVECodeMethod:  attackProcessCVECode,
	attackProcessDescMethod:     attackProcessDesc,
	attackProcessSuggestMethod:  attackProcessSuggest,
	attackProcessStatusMethod:   attackProcessStatus,
	attackProcessPayloadMethod:  attackProcessPayload,
	attackProcessResultMethod:   attackProcessResult,
	attackProcessDetailsMehtod:  attackProcessDetails,
}

func newAttackProcessApi(L *glua.LState) int {

	ap := &AttackProcess{}
	luahelper.SetUserData(L, attackProcessUDName, ap)

	return 1
}

func checkAttackProcess(L *glua.LState, n int) *AttackProcess {

	ud := L.CheckUserData(n)

	if v, ok := ud.Value.(*AttackProcess); ok {
		return v
	}

	L.ArgError(1, "AttackProcess expected")
	return nil
}

func attackProcessIP(L *glua.LState) int {

	ap := checkAttackProcess(L, 1)

	if L.GetTop() == 2 {

		v := L.CheckString(2)
		ap.IP = v
		return 0
	}

	L.Push(glua.LString(ap.IP))

	return 1
}

func attackProcessHost(L *glua.LState) int {

	ap := checkAttackProcess(L, 1)

	if L.GetTop() == 2 {

		v := L.CheckString(2)
		ap.Host = v
		return 0
	}

	L.Push(glua.LString(ap.Host))

	return 1
}

func attackProcessPort(L *glua.LState) int {

	ap := checkAttackProcess(L, 1)

	if L.GetTop() == 2 {

		v := L.CheckInt(2)
		ap.Port = v
		return 0
	}

	L.Push(glua.LNumber(ap.Port))

	return 1
}

func attackProcessProto(L *glua.LState) int {

	ap := checkAttackProcess(L, 1)

	if L.GetTop() == 2 {

		v := L.CheckString(2)
		ap.Proto = v
		return 0
	}

	L.Push(glua.LString(ap.Proto))

	return 1
}

func attackProcessApp(L *glua.LState) int {

	ap := checkAttackProcess(L, 1)

	if L.GetTop() == 2 {

		v := L.CheckString(2)
		ap.App = v
		return 0
	}

	L.Push(glua.LString(ap.App))

	return 1
}

func attackProcessOs(L *glua.LState) int {

	ap := checkAttackProcess(L, 1)

	if L.GetTop() == 2 {

		v := L.CheckString(2)
		ap.OS = v
		return 0
	}

	L.Push(glua.LString(ap.OS))

	return 1
}

func attackProcessVersion(L *glua.LState) int {

	ap := checkAttackProcess(L, 1)

	if L.GetTop() == 2 {

		v := L.CheckString(2)
		ap.Version = v
		return 0
	}

	L.Push(glua.LString(ap.Version))

	return 1
}

func attackProcessIsSSL(L *glua.LState) int {

	ap := checkAttackProcess(L, 1)

	if L.GetTop() == 2 {

		v := L.CheckBool(2)
		ap.IsSSL = v
		return 0
	}

	L.Push(glua.LBool(ap.IsSSL))

	return 1
}

func attackProcessId(L *glua.LState) int {

	ap := checkAttackProcess(L, 1)

	if L.GetTop() == 2 {

		v := L.CheckInt(2)
		ap.Id = v
		return 0
	}

	L.Push(glua.LNumber(ap.Id))

	return 1
}

func attackProcessLanguage(L *glua.LState) int {

	ap := checkAttackProcess(L, 1)

	if L.GetTop() == 2 {

		v := L.CheckString(2)
		ap.Language = v
		return 0
	}

	L.Push(glua.LString(ap.Language))

	return 1
}

func attackProcessName(L *glua.LState) int {

	ap := checkAttackProcess(L, 1)

	if L.GetTop() == 2 {

		v := L.CheckString(2)
		ap.Name = v
		return 0
	}

	L.Push(glua.LString(ap.Name))

	return 1
}

func attackProcessType(L *glua.LState) int {

	ap := checkAttackProcess(L, 1)

	if L.GetTop() == 2 {

		v := L.CheckString(2)
		ap.Type = v
		return 0
	}

	L.Push(glua.LString(ap.Type))

	return 1
}

func attackProcessCVECode(L *glua.LState) int {

	ap := checkAttackProcess(L, 1)

	if L.GetTop() == 2 {

		v := L.CheckString(2)
		ap.CVECode = v
		return 0
	}

	L.Push(glua.LString(ap.CVECode))

	return 1
}

func attackProcessDesc(L *glua.LState) int {

	ap := checkAttackProcess(L, 1)

	if L.GetTop() == 2 {

		v := L.CheckString(2)
		ap.Desc = v
		return 0
	}

	L.Push(glua.LString(ap.Desc))

	return 1
}

func attackProcessDetails(L *glua.LState) int {

	ap := checkAttackProcess(L, 1)

	if L.GetTop() == 2 {

		v := L.CheckString(2)
		ap.Details = v
		return 0
	}

	L.Push(glua.LString(ap.Details))

	return 1
}

func attackProcessStatus(L *glua.LState) int {

	ap := checkAttackProcess(L, 1)

	if L.GetTop() == 2 {

		v := L.CheckInt(2)
		ap.Status = v
		return 0
	}

	L.Push(glua.LNumber(ap.Status))

	return 1
}

func attackProcessSuggest(L *glua.LState) int {

	ap := checkAttackProcess(L, 1)

	if L.GetTop() == 2 {

		v := L.CheckString(2)
		ap.Suggest = v
		return 0
	}

	L.Push(glua.LString(ap.Suggest))

	return 1
}

func attackProcessPayload(L *glua.LState) int {

	ap := checkAttackProcess(L, 1)

	if L.GetTop() == 2 {

		v := L.CheckString(2)
		ap.Payload = v
		return 0
	}

	L.Push(glua.LString(ap.Payload))

	return 1
}

func attackProcessResult(L *glua.LState) int {

	ap := checkAttackProcess(L, 1)

	if L.GetTop() == 2 {

		v := L.CheckString(2)
		ap.Result = v
		return 0
	}

	L.Push(glua.LString(ap.Result))

	return 1
}
