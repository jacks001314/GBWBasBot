package main

import (
	"common/proto/http"
	luahelper "common/scripts/lua"
	"log"

	glua "github.com/yuin/gopher-lua"
)

const httpscript = `
        local http = require("http")
        local person =  require("person")
		local host = "www.baidu.com"
		local port = 80
        local timeout = 10000
        local isSSL = false
		client = http.newHttpClient(host,port,isSSL,timeout)
        request = http.newHttpRequest("get","/")

        request:addHeaders({User_Agent="LuaClient",Connection="Close"})
        response = client:send(request)
        print(response:getStatusCode())
        local headers = response:getHeaders()
       
		print(response:getHeader("Connection"))

		local pp = person.pp 

		print(pp:name())
		print(pp:age())
		pp:age(120)
	
		print(pp:name())
		print(pp:age())
    `

const gscript = `
   

`

type Person struct {
	name string
	age  int
}

var apis = map[string]glua.LGFunction{
	"name": getName,
	"age":  getAge,
}

func checkPerson(L *glua.LState, n int) *Person {

	ud := L.CheckUserData(n)

	if v, ok := ud.Value.(*Person); ok {
		return v
	}

	L.ArgError(1, "Person expected")
	return nil
}

func getName(L *glua.LState) int {

	p := checkPerson(L, 1)

	if L.GetTop() == 2 {

		v := L.CheckString(2)
		p.name = v
		return 0
	}

	L.Push(glua.LString(p.name))

	return 1
}

func getAge(L *glua.LState) int {

	p := checkPerson(L, 1)

	if L.GetTop() == 2 {

		v := L.CheckInt(2)
		p.age = v
		return 0
	}

	L.Push(glua.LNumber(p.age))

	return 1
}

func personLoader(L *glua.LState) int {

	p := &Person{
		name: "shajf",
		age:  40,
	}

	mod := L.NewTable()

	luahelper.RegisterUserData(L, mod, "pp", p, apis)

	L.Push(mod)

	return 1
}

func testLua(s, name string) {

	bcode, err := luahelper.CompileLuaScript([]byte(s), name)

	if err != nil {

		log.Fatalf("%v\n", err)
	}

	L := glua.NewState()
	defer L.Close()

	luahelper.RegisterModule(L, "http", http.Loader)
	luahelper.RegisterModule(L, "person", personLoader)

	//L.PreloadModule("http", http.Loader)

	err = luahelper.RunLua(L, bcode)
	if err != nil {

		log.Fatalf("%v\n", err)
	}

}

func main() {

	testLua(httpscript, "p.lua")
}
