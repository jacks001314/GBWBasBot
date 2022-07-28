--[[
    find apache hadoop application
--]]

local http  = require("http")
local detect = require("detect")


local UA = "FindHadoopApplication"
local timeoutMS = 10000
local host = target:ip()
local port = target:port()

local key = "Hadoop IPC port"

function fetchResponse(isSSL)
    local client = http.newHttpClient(host,port,isSSL,timeoutMS)
    local request = http.newHttpRequest("get","/")
    request:addHeader("User-Agent",UA)
    local response = client:send(request)
    local content = response:getBodyAsString()
    return content
end

function isMatch(content)
    return string.match(content,key)
end

function pub(isSSL)

    local proto = "http"
    if isSSL then
        proto = "https"
    end

    local r = detect.newDetectResult()
    r:ip(host)
    r:port(port)
    r:app("hadoop")
    r:version("")
    r:proto(proto)

    r:isSSL(isSSL)

    script:publish(r)
end

function main()

    local content = fetchResponse(false)

    if isMatch(content) then
        pub(false)
        return true
    end

    content = fetchResponse(true)
    if isMatch(content) then
        pub(true)
        return true
    end

    return false

end

return main()