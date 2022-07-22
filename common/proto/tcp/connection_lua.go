package tcp

import (
	"encoding/hex"
	"fmt"
	"time"

	luahelper "common/scripts/lua"

	glua "github.com/yuin/gopher-lua"
)

var connectionApis = map[string]glua.LGFunction{

	closeMethod:       closeApi,
	flushMethod:       flushApi,
	writeHexMethod:    writeHexApi,
	writeStringMethod: writeStringApi,
	readLineMethod:    readLineApi,
	readBytesMethod:   readBytesApi,
	readStringMethod:  readStringApi,
}

/*
*
  args[0] ---network
  args[1] ---host
  args[2] ---port
  args[3] ---isSSL
  args[4] ---connectionTimeout
  args[5] ---readTimeout
  args[6] ---writeTimeout
*/

func newConnectionApi(L *glua.LState) int {

	var conn *Connection
	var err error

	network := L.CheckString(1)
	host := L.CheckString(2)
	port := L.CheckInt(3)
	isSSL := L.CheckBool(4)
	ctimeout := L.CheckInt64(5)
	rtimeout := L.CheckInt64(6)
	wtimeout := L.CheckInt64(7)

	addr := fmt.Sprintf("%s:%d", host, port)

	if isSSL {

		conn, err = Dial(network, addr,
			DialConnectTimeout(time.Duration(ctimeout)*time.Millisecond),
			DialReadTimeout(time.Duration(rtimeout)*time.Millisecond),
			DialWriteTimeout(time.Duration(wtimeout)*time.Millisecond),
			DialTLSSkipVerify(true),
			DialUseTLS(true))

	} else {

		conn, err = Dial(network, addr,
			DialConnectTimeout(time.Duration(ctimeout)*time.Millisecond),
			DialReadTimeout(time.Duration(rtimeout)*time.Millisecond),
			DialWriteTimeout(time.Duration(wtimeout)*time.Millisecond))
	}

	if err != nil {
		L.Push(glua.LNil)
		L.Push(glua.LString(err.Error()))
		return 2
	}

	luahelper.SetUserData(L, tcpConnectionUDName, conn)

	return 1
}

func checkConnection(L *glua.LState, n int) *Connection {

	ud := L.CheckUserData(n)

	if v, ok := ud.Value.(*Connection); ok {
		return v
	}

	L.ArgError(1, "Connection expected")
	return nil
}

func closeApi(L *glua.LState) int {

	conn := checkConnection(L, 1)

	conn.Close()

	return 0
}

func flushApi(L *glua.LState) int {

	conn := checkConnection(L, 1)

	if err := conn.Flush(); err != nil {

		L.Push(glua.LString(err.Error()))
		return 1
	}

	return 0
}

func writeHexApi(L *glua.LState) int {

	conn := checkConnection(L, 1)
	content := L.CheckAny(2).String()

	if err := conn.WriteHex(content); err != nil {

		L.Push(glua.LString(err.Error()))
		return 1
	}

	return 0
}

func writeStringApi(L *glua.LState) int {

	conn := checkConnection(L, 1)
	content := L.CheckAny(2).String()

	if err := conn.WriteString(content); err != nil {

		L.Push(glua.LString(err.Error()))
		return 1
	}

	return 0
}

func readLineApi(L *glua.LState) int {

	conn := checkConnection(L, 1)
	content, err := conn.ReadLine()

	if err != nil {
		L.Push(glua.LString(""))
		L.Push(glua.LString(err.Error()))
		return 2
	}

	L.Push(glua.LString(string(content)))
	return 1
}

func readBytesApi(L *glua.LState) int {

	conn := checkConnection(L, 1)
	bytes := L.CheckInt(2)

	content, err := conn.ReadBytes(bytes)

	if err != nil {
		L.Push(glua.LString(""))
		L.Push(glua.LString(err.Error()))
		return 2
	}

	L.Push(glua.LString(hex.EncodeToString(content)))
	return 1
}

func readStringApi(L *glua.LState) int {

	conn := checkConnection(L, 1)
	bytes := L.CheckInt(2)

	content, err := conn.ReadBytes(bytes)

	if err != nil {
		L.Push(glua.LString(""))
		L.Push(glua.LString(err.Error()))
		return 2
	}

	L.Push(glua.LString(string(content)))
	return 1
}
