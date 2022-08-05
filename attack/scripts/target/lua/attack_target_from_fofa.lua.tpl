---
--- Created by jacks.
--- DateTime: 2022/8/5 10:43
--- Download Attack Targets from fofa.info

local http = require("http")
local source = require("source")
local json = require("json")
local base64 = require("base64")
local strings = require("strings")


local host = "fofa.info"
local port = 443
local isSSL = true
local timeOut =10000

local email = "{{.Email}}"
local key  = "{{.Key}}"

local query = base64.StdEncoding:encode_to_string('{{.Query}}')

local defaultProto = "{{.Proto}}"
local defaultPort = {{.Port}}

local pageSize = 100

local client = http.newHttpClient(host,port,isSSL,10000)

function fetch(page,size)
    local url = string.format("/api/v1/search/all?email=%s&key=%s&qbase64=%s&page=%d&size=%d",email,key,query,page,size)
    local request = http.newHttpRequest("get",url)
    request:addHeader("User-Agent","GOClient")

    local response = client:send(request)

    if response:getStatusCode() ~=200 then
        return nil
    end

    local content = response:getBodyAsString()

    if content == "" or not strings.contains(content,"results") then
        return nil
    end

    return json.decode(content)
end

function parseData(data)

    if strings.has_prefix(data[1],"https://") then
        return {"https",data[2],strings.parse_int(data[3],10,32)}
    end

    return {defaultProto,data[2],strings.parse_int(data[3],10,32)}
end

function putTarget(pdata)

    local proto = pdata[1]
    local ip = pdata[2]
    local port = pdata[3]
    local target = source.newTarget()
    target:ip(ip)
    target:host(ip)

    if port == 0 then
        port = defaultPort
    end

    target:port(port)

    target:app("{{.App}}")
    target:version("{{.Version}}")
    target:proto(proto)

    local ssl = {{.IsSSL}}
    if proto == "https" then
        ssl = true
    end

    target:isSSL(ssl)

    script:put(target)
end

function main()

    local page = 1

    while(true) do
        local jsonData = fetch(page,pageSize)
        if  jsonData == nil or jsonData["error"] then
            break
        end

        local results = jsonData["results"]

        for _,data in pairs(results) do

            local pdata = parseData(data)
            putTarget(pdata)
        end
        page = page+1
    end

end

main()
