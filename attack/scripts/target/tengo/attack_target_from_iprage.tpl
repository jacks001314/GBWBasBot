fmt  := import("fmt")
source := import("source")
ipv4 := import("ipv4")

wlists := {{toTStrArray .WhiteLists}}
blists  := {{toTStrArray .BlackLists}}

ipg := ipv4.newFromArray(wlists,blists)


putTarget := func (ip) {

     target := source.newTarget()
     target.ip(ip)
     target.host(ip)
     target.port({{.Port}})
     target.app("{{.App}}")
     target.version("{{.Version}}")
     target.proto("{{.Proto}}")
     target.isSSL({{.IsSSL}})

     script.put(target)

}

main := func () {

    for ip:= ipg.curIP();ip!="";ip= ipg.nextIP() {

        putTarget(ip)
    }
}

main()
