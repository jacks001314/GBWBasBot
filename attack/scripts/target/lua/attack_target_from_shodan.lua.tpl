---
--- Created by jacks.
--- DateTime: 2022/8/5 13:47
---Download Attack Targets from shadan

local http = require("http")
local source = require("source")
local json  = require("json")

local host = "api.shodan.io"
local port =443
local timeoutMS = 10000
local UA = "GoClient"

local key = "{{.Key}}"
local query = http.urlEncode('{{.Query}}')
local defaultPort = {{.Port}}

local client = http.newHttpClient(host,port,true,timeoutMS)

function putTarget(ip,port)

    local target = source.newTarget()
    target:ip(ip)
    target:host(ip)

    if port == 0 then
        port = defaultPort
    end
    target:port(port)

    target:app("{{.App}}")
    target:version("{{.Version}}")
    target:proto("{{.Proto}}")
    target:isSSL({{.IsSSL}})
    script:put(target)
end

---parse fetch json data
function parseData(content)

    local count = 0
    local jsonData = json.decode(content)

    if jsonData["matches"] == nil then
        return 0
    end

    for _,entry in pairs(jsonData["matches"]) do
        putTarget(entry["ip_str"],entry["port"])
        count = count+1
    end

    return count

end

---get host data from shodan by restfull api
function main()

    local page = 1

    while(true) do
        local url = string.format("/shodan/host/search?key=%s&query=%s&page=%d",key,query,page)
        local request = http.newHttpRequest("get",url)
        request:addHeader("User-Agent",UA)
        local response = client:send(request)
        if response:getStatusCode() ~=200 then
            break
        end
        local content = response:getBodyAsString()
        if content == "" then
            break
        end

        if parseData(content) <= 0 then
            break
        end
        page = page+1
    end
end

main()