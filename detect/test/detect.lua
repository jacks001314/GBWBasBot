local http = require("http")
local detect = require("detect")

local result = detect.newDetectResult()

result:ip(target:ip())
result:port(target:port())
result:app("web")
result:isSSL(true)

script:publish(result)









