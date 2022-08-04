---
--- Created by jacks.
--- DateTime: 2022/8/4 15:48
---

local source = require("source")
local ipv4  = require("ipv4")

local wlists = {{toLStrArray .WhiteLists}}
local blists = {{toLStrArray .BlackLists}}

local ipg = ipv4.newFromArray(wlists,blists)

function putTarget(ip)

    local target = source.newTarget()
    target:ip(ip)
    target:host(ip)
    target:port({{.Port}})
    target:app("{{.App}}")
    target:version("{{.Version}}")
    target:proto("{{.Proto}}")
    target:isSSL({{.IsSSL}})

    script:put(target)

end

function main()

    local ip = ipg:curIP()
    while(ip ~= "") do
        putTarget(ip)
        ip = ipg:nextIP()
    end

end

main()
