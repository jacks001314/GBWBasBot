package main

import (
	attack "attack/core"
	"attack/target"
	"common/scripts"
	"common/util/jsonutils"
	tplutils "common/util/tpl"
	"flag"
	"fmt"
	"log"
	"runtime"
	"strings"
	"sync"
)

func run(stype scripts.ScriptType, tfile string, gensf bool, data *target.AttackTargetData) {

	content, err := tplutils.MakeSourceScriptFromTemplate(tfile, data)

	if err != nil {

		log.Fatalf("%v", err)
	}

	if gensf {

		fmt.Println(content)
		return
	}

	tchan := make(chan *attack.AttackTarget, 1000)

	stask, err := target.NewSourceTask(tchan)

	if err != nil {

		log.Fatal(err)
	}

	wg := sync.WaitGroup{}

	wg.Add(2)

	go func() {

		defer func() {
			wg.Done()
		}()

		err := stask.Fetch(stype, content, []string{"target"})
		if err != nil {
			log.Fatal(err)
		}

	}()

	count := 0
	go func() {
		defer func() {

			if p := recover(); p != nil {

				var buf [4096]byte
				n := runtime.Stack(buf[:], false)

				log.Printf("Attack Source Exit from panic:%s", string(buf[:n]))
			}

			wg.Done()
		}()

		for {
			select {
			case t := <-tchan:
				fmt.Println(jsonutils.ToJsonString(t, true))
				//time.Sleep(2 * time.Second)
				count++

			}
		}
	}()

	wg.Wait()
}

func main() {

	stypes := flag.String("stype", "lua", "the attack target template script language:<lua/tengo>")
	tplFile := flag.String("tfile", "", "please specify attack targets template file path")
	wlist := flag.String("wlist", "", "attack target ip white list,for example:192.168.1.0/24,10.0.1.0/24")
	blist := flag.String("blist", "", "attack target ip white list,for example:192.168.1.1,10.0.1.1")
	email := flag.String("email", "", "download attack target from fofa ,need an email address")
	key := flag.String("key", "", "download attack targets from third part org by restfull api,need an api key")
	query := flag.String("query", "", "download attack targets from third part org by restfull api,specify query string")
	port := flag.Int("port", 0, "attack target default port")
	isSSL := flag.Bool("isSSL", false, "whether use ssl proto to connect to attack target or not")
	isGenScriptFile := flag.Bool("isGSF", false, "whether only generate a attack targets script file or not")

	flag.Parse()

	if *tplFile == "" {

		log.Fatal("please specify attack targets template file path!")
	}

	data := &target.AttackTargetData{
		WhiteLists: strings.Split(*wlist, ","),
		BlackLists: strings.Split(*blist, ","),
		Email:      *email,
		Key:        *key,
		Query:      *query,
		Port:       *port,
		IsSSL:      *isSSL,
		Version:    "1.0",
		Proto:      "http",
		App:        "web",
	}

	stype := scripts.ScriptLua

	if *stypes == "tengo" {

		stype = scripts.ScriptTengo
	}

	run(stype, *tplFile, *isGenScriptFile, data)
}
