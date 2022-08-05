
/*
*This Tengo Script Download Attack target data from fofa.info
**/

fmt  := import("fmt")
http := import("http")
source := import("source")
json := import("json")
text := import("text")
base64 := import("base64")

host:= "fofa.info"
port:= 443
isSSL:= true
timeOut:=10000

email := "{{.Email}}"
key := "{{.Key}}"
query := base64.encode(`{{.Query}}`)
defaultProto := "{{.Proto}}"
defaultPort := {{.Port}}

pageSize := 100
client := http.newHttpClient(host,port,isSSL,10000)

fetch := func (page,size) {

    url := fmt.sprintf("/api/v1/search/all?email=%s&key=%s&qbase64=%s&page=%d&size=%d",email,key,query,page,size)
    request := http.newHttpRequest("get",url)
    request.addHeader("User-Agent","GOClient")
    response := client.send(request)

    if response.getStatusCode() !=200 {
        return {}
    }

    content := response.getBodyAsString()
    if content == "" || !text.contains(content,"results"){
        return {}
    }

    return json.decode(content)
}

parseData := func (data) {

    if text.has_prefix(data[0],"https://") {
        return ["https",data[1],text.parse_int(data[2],10,32)]
    }else {
        return [defaultProto,data[1],text.parse_int(data[2],10,32)]
    }
}

putTarget := func (pdata) {

     ip := pdata[1]
     port := pdata[2]
     proto := pdata[0]

     target := source.newTarget()
     target.ip(ip)
     target.host(ip)
     target.port(port==0?defaultPort:port)
     target.app("{{.App}}")
     target.version("{{.Version}}")
     target.proto(proto)

     target.isSSL(proto=="https"?true:{{.IsSSL}})


     script.put(target)

}

main := func () {

    page := 1

    for {

        jsonData := fetch(page,pageSize)

        if len(jsonData) == 0 || jsonData["error"] {
            break
        }

        results := jsonData["results"]

        if len(results) == 0 {
            break
        }

        for data in results {

            if len(data)==3 {
                pdata := parseData(data)
                putTarget(pdata)
            }
        }

        page = page+1
    }
}

main()



