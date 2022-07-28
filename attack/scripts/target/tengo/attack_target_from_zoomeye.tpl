/*
*This Tengo Script Download Attack Target from zoomeye
**/

fmt  := import("fmt")
http := import("http")
source := import("source")
json  := import("json")

host := "api.zoomeye.org"
port :=443
timeoutMS := 10000
UA := "GoClient"

key := "{{.Key}}"
query := http.urlEncode(`{{.Query}}`)
useDefaultPort := {{.UseDefaultPort}}

client := http.newHttpClient(host,port,true,timeoutMS)

setEntry := func (ip,port) {

     entry := source.newEntry()
     entry.setIP(ip)
     entry.setHost(ip)
     if useDefaultPort {
        entry.setPort(0)

     }else {
        entry.setPort(port)
     }

     entry.setProto("{{.Proto}}")
     entry.setApp("{{.App}}")

     scriptSource.put(entry)
}

//parse fetch json data
parseData := func(content) {

     count := 0
     jsonData := json.decode(content)

     if is_error(jsonData)||len(jsonData["matches"])==0 {
            return 0
     }

     for entry in jsonData["matches"] {

        setEntry(entry["ip"],entry["portinfo"]["port"])
        count++
     }

     return count
}

//get host data from zoomeye by restfull api
fetchData := func() {

    page := 1

    for {

        url := fmt.sprintf("/host/search?query=%s&page=%d",query,page)
        request := http.newHttpRequest("get",url).addHeader("User-Agent",UA).addHeader("API-KEY",key)

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

main := func () {
    fetchData()
}

main()
