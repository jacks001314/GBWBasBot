fmt  := import("fmt")
source := import("source")
ipgen := import("ipgen")

wlists := {{tostring .WhiteLists}}
blists  := {{tostring .BlackLists}}

ipg := ipgen.newIPGenFromArray(wlists,blists)

setEntry := func (ip) {

     entry := source.newEntry()
     entry.setIP(ip)
     entry.setHost(ip)
     entry.setPort(0)
     entry.setProto("{{.Proto}}")
     entry.setApp("{{.App}}")

     scriptSource.put(entry)

}

main := func () {

    for ip:= ipg.curIP();ip!="";ip= ipg.nextIP() {

        setEntry(ip)
    }
}

main()
