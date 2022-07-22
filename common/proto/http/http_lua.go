package http

import (
	luahelper "common/scripts/lua"

	glua "github.com/yuin/gopher-lua"
)

func Loader(L *glua.LState) int {

	mod := L.SetFuncs(L.NewTable(), httpApis)

	//register http client apis
	luahelper.RegisterApis(L, mod, httpClientApis, httpClientApiName, httpClientUDName)

	//register http request apis
	luahelper.RegisterApis(L, mod, httpRequestApis, httpRequestApiName, httpRequestUDName)

	//register http response apis
	luahelper.RegisterApis(L, mod, httpResponseApis, httpResponseApiName, httpResponseUDName)

	L.Push(mod)

	return 1
}

var httpApis = map[string]glua.LGFunction{

	newHttpClientMethod:  newHttpClientApi,
	newHttpRequestMethod: newHttpRequestApi,

	urlEncodeMethod: urlEncodeApi,

	urlDecodeMethod: urlDecodeApi,
}
