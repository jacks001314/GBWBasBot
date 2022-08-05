/*
*This Tengo Script Download Attack Target from shodan
**/

fmt  := import("fmt")
http := import("http")
source := import("source")
json  := import("json")

host := "api.shodan.io"
port :=443
timeoutMS := 10000
UA := "GoClient"

key := "{{.Key}}"
query := http.urlEncode(`{{.Query}}`)
defaultPort := {{.Port}}

client := http.newHttpClient(host,port,true,timeoutMS)

putTarget := func (ip,port) {

     target := source.newTarget()
     target.ip(ip)
     target.host(ip)
     if port == 0 {
        target.port(defaultPort)
     }else {
        target.port(port)
     }

     target.app("{{.App}}")
     target.version("{{.Version}}")
     target.proto("{{.Proto}}")
     target.isSSL({{.IsSSL}})

     script.put(target)

}

//parse fetch json data
parseData := func(content) {

     count := 0
     jsonData := json.decode(content)

     if is_error(jsonData)||len(jsonData["matches"])==0 {
            return 0
     }

     for entry in jsonData["matches"] {

        putTarget(entry["ip_str"],entry["port"])
        count++
     }

     return count
}

//get host data from shodan by restfull api
main := func() {

    page := 1

    for {

        url := fmt.sprintf("/shodan/host/search?key=%s&query=%s&page=%d",key,query,page)
        request := http.newHttpRequest("get",url)
        request.addHeader("User-Agent",UA)

        response := client.send(request)

        if response.getStatusCode() !=200 {
                break
         }

         content := response.getBodyAsString()
         if content == "" {
                break
         }

         if parseData(content) <=0 {
            break
         }
         page++
    }

}

main()
