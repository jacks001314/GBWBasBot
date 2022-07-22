package http

import (
	"net/url"

	glua "github.com/yuin/gopher-lua"
)

func urlEncodeApi(L *glua.LState) int {

	urlRaw := L.CheckString(1)

	L.Push(glua.LString(url.QueryEscape(urlRaw)))

	return 1
}

func urlDecodeApi(L *glua.LState) int {

	urlRaw := L.CheckString(1)

	durl, err := url.QueryUnescape(urlRaw)

	if err != nil {
		L.Push(glua.LString(err.Error()))
		return 2
	}

	L.Push(glua.LString(durl))

	return 1
}
