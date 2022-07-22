package ipv4

import (
	"common/util/netutils"

	luahelper "common/scripts/lua"

	glua "github.com/yuin/gopher-lua"
)

var ipv4genApis = map[string]glua.LGFunction{

	ipv4GenCurIPMethod:  getCurIP,
	ipv4GenNextIPMethod: getNextIP,
}

func newIPV4GenFromFile(L *glua.LState) int {

	ipgen, err := NewIPV4Generator(L.ToString(1), L.ToString(2), []string{}, []string{}, true)

	if err != nil {
		L.Push(glua.LNil)
		L.Push(glua.LString(err.Error()))
		return 2
	}

	luahelper.SetUserData(L, ipv4genUDName, ipgen)

	return 1
}

func newIPV4GenFromArray(L *glua.LState) int {

	wlist := luahelper.LuaTableToStringArray(L.ToTable(1))
	blist := luahelper.LuaTableToStringArray(L.ToTable(2))

	ipgen, err := NewIPV4Generator("", "", wlist, blist, true)

	if err != nil {
		L.Push(glua.LNil)
		L.Push(glua.LString(err.Error()))
		return 2
	}

	luahelper.SetUserData(L, ipv4genUDName, ipgen)

	return 1

}

func checkIPV4Gen(L *glua.LState) *IPV4Generator {

	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*IPV4Generator); ok {
		return v
	}

	L.ArgError(1, "ipv4.gen expected")
	return nil
}

func getCurIP(L *glua.LState) int {

	ipgen := checkIPV4Gen(L)

	ip := ipgen.GetCurIP()
	ipstr := ""

	if ip != 0 {

		ipstr = netutils.IPv4StrBig(ip)
	}

	L.Push(glua.LString(ipstr))

	return 1
}

func getNextIP(L *glua.LState) int {

	ipgen := checkIPV4Gen(L)

	ip := ipgen.GetNextIP()
	ipstr := ""

	if ip != 0 {

		ipstr = netutils.IPv4StrBig(ip)
	}

	L.Push(glua.LString(ipstr))

	return 1
}
