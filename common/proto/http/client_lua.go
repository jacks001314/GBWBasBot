package http

import (
	luahelper "common/scripts/lua"

	glua "github.com/yuin/gopher-lua"
)

var httpClientApis = map[string]glua.LGFunction{

	sendMethod: sendApi,
}

func newHttpClientApi(L *glua.LState) int {

	host := L.CheckString(1)
	port := L.CheckInt(2)
	isSSL := L.CheckBool(3)
	timeOut := L.CheckInt64(4)

	httpClient := NewHttpClient(host, port, isSSL, timeOut)

	luahelper.SetUserData(L, httpClientUDName, httpClient)

	return 1
}

func checkClient(L *glua.LState) *HttpClient {

	ud := L.CheckUserData(1)

	if v, ok := ud.Value.(*HttpClient); ok {
		return v
	}

	L.ArgError(1, "HttpClient expected")
	return nil
}

func sendApi(L *glua.LState) int {

	client := checkClient(L)
	request := checkRequest(L, 2)

	res, err := client.Send(request)
	if err != nil {
		L.Push(glua.LNil)
		L.Push(glua.LString(err.Error()))
		return 2
	}

	luahelper.SetUserData(L, httpResponseUDName, res)

	return 1
}
