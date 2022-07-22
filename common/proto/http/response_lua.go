package http

import (
	glua "github.com/yuin/gopher-lua"
)

var httpResponseApis = map[string]glua.LGFunction{

	getStatusCodeMethod:   getStatusCodeApi,
	getBodyAsByteMethod:   getBodyAsByteApi,
	getBodyAsStringMethod: getBodyAsStringApi,
	getProtocolMethod:     getProtocolApi,
	getHeaderMethod:       getHeaderApi,
	getHeadersMethod:      getHeadersApi,
}

func checkResponse(L *glua.LState, n int) *HttpResponse {

	ud := L.CheckUserData(n)

	if v, ok := ud.Value.(*HttpResponse); ok {
		return v
	}

	L.ArgError(1, "HttpResponse expected")
	return nil
}

func getStatusCodeApi(L *glua.LState) int {

	res := checkResponse(L, 1)

	L.Push(glua.LNumber(res.GetStatusCode()))

	return 1
}

func getBodyAsByteApi(L *glua.LState) int {

	return getBodyAsStringApi(L)
}

func getBodyAsStringApi(L *glua.LState) int {

	res := checkResponse(L, 1)

	content, _ := res.GetBodyAsString()

	L.Push(glua.LString(content))

	return 1
}

func getProtocolApi(L *glua.LState) int {

	res := checkResponse(L, 1)
	L.Push(glua.LString(res.Protocol()))

	return 1
}

func getHeaderApi(L *glua.LState) int {

	res := checkResponse(L, 1)

	key := L.CheckString(2)

	L.Push(glua.LString(res.GetHeaderValue(key)))

	return 1
}

func getHeadersApi(L *glua.LState) int {

	res := checkResponse(L, 1)

	headers := L.NewTable()

	resHeader := res.GetHader()

	if resHeader != nil {

		for k, v := range res.resp.Header {

			if len(v) > 0 {

				headers.RawSetString(k, glua.LString(v[0]))
			}
		}

	}

	L.Push(headers)

	return 1
}
