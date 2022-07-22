package http

import (
	luahelper "common/scripts/lua"

	glua "github.com/yuin/gopher-lua"
)

var httpRequestApis = map[string]glua.LGFunction{

	authMethod:       authApi,
	addHeaderMethod:  addHeaderApi,
	addHeadersMethod: addHeadersApi,
	putStringMethod:  putStringApi,
	putHexMethod:     putHexApi,
	uploadMethod:     uploadApi,
}

func checkRequest(L *glua.LState, n int) *HttpRequest {

	ud := L.CheckUserData(n)

	if v, ok := ud.Value.(*HttpRequest); ok {
		return v
	}

	L.ArgError(1, "HttpRequest expected")
	return nil
}

func newHttpRequestApi(L *glua.LState) int {

	method := L.CheckString(1)

	uri := L.CheckString(2)

	request := NewHttpRequest(method, uri)

	luahelper.SetUserData(L, httpRequestUDName, request)

	return 1
}

//授权函数
func authApi(L *glua.LState) int {

	req := checkRequest(L, 1)
	user := L.CheckString(2)
	pass := L.CheckString(3)

	req.BasicAuth(user, pass)

	return 0
}

//添加请求头部
func addHeaderApi(L *glua.LState) int {

	req := checkRequest(L, 1)
	key := L.CheckString(2)
	value := L.CheckString(3)

	req.AddHeader(key, value)

	return 0
}

//添加请求头部
func addHeadersApi(L *glua.LState) int {

	req := checkRequest(L, 1)
	values := L.CheckTable(2)

	values.ForEach(func(k glua.LValue, v glua.LValue) {

		req.AddHeader(k.String(), v.String())
	})

	return 0
}

//通过字符串格式设置post/put等方法请求体
func putStringApi(L *glua.LState) int {

	req := checkRequest(L, 1)

	content := L.CheckString(2)
	isFromFile := L.CheckBool(3)

	req.BodyString(content, isFromFile)

	return 0
}

// 通过十六进制格式设置post/put等方法请求体
func putHexApi(L *glua.LState) int {

	req := checkRequest(L, 1)

	content := L.CheckString(2)

	req.BodyHex(content)

	return 0
}

//文件上传
func uploadApi(L *glua.LState) int {

	req := checkRequest(L, 1)

	fname := L.CheckString(2)

	fpath := L.CheckString(3)

	boundary := L.CheckString(4)

	req.UPload(fname, fpath, boundary)

	return 0
}
