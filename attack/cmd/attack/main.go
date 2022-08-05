package main

import (
	attack "attack/core"
	"attack/target"
	"common/scripts"
	"flag"
	"io/ioutil"
	"log"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

var (
	wg         = &sync.WaitGroup{}
	atargets   = make(chan *attack.AttackTarget, 1000)
	aprocesses = make(chan *attack.AttackProcess, 1000)
)

func main() {

	sstype := flag.String("sstype", "lua", "Attack source script type:[lua/tengo]")
	astype := flag.String("astype", "lua", "Attack script type:[lua/tengo]")
	ssfile := flag.String("ssfile", "", "Attack source script file path")
	asfile := flag.String("asfile", "", "Attack script file path")

	targets := flag.String("targets", "", "attack targets format:[ip:port:isSSL,ip:port:isSSL],for example:192.168.1.2:8080:true,192.168.1.3:80:false")

	flag.Parse()

	if *asfile == "" {

		log.Fatal("Must Specify an attack script file path")
	}

	ssType := scripts.ScriptLua
	if *sstype == "tengo" {
		ssType = scripts.ScriptTengo
	}

	asType := scripts.ScriptLua

	if *astype == "tengo" {

		asType = scripts.ScriptTengo
	}

	attackSource, err := target.NewSourceTask(atargets)

	if err != nil {

		log.Fatal(err)
	}

	attackTask, err := attack.NewAttackTask(0, 0, 0, atargets, aprocesses)

	if err != nil {

		log.Fatal(err)
	}

	config := &attack.Config{
		Enable:       true,
		Name:         "attack",
		Author:       "jacks",
		Atype:        "Test",
		Language:     *astype,
		App:          "Test",
		Id:           0,
		CVECode:      "Test",
		Desc:         "Test",
		Suggest:      "Test",
		DefaultPort:  0,
		DefaultProto: "Test",
		FPath:        *asfile,
	}

	attackTask.AddAttackScriptFromFile(asType, config)

	wg.Add(2)

	go func() {
		defer func() {

			if p := recover(); p != nil {

				var buf [4096]byte
				n := runtime.Stack(buf[:], false)

				log.Printf("Attack Task Exit from panic:%s", string(buf[:n]))
			}

			wg.Done()
		}()

		attackTask.Start()
	}()

	go func() {
		defer func() {
			wg.Done()
		}()

		if *ssfile != "" {

			content, err := ioutil.ReadFile(*ssfile)

			if err != nil {

				log.Fatal(err)
			}

			attackSource.Fetch(ssType, content, []string{"target"})
		}

		if *targets != "" {

			arr := strings.Split(*targets, ",")

			for _, t := range arr {

				tts := strings.Split(t, ":")

				port, err := strconv.ParseInt(tts[1], 10, 32)

				if err != nil {
					continue
				}
				isSSL, err := strconv.ParseBool(tts[2])
				if err != nil {
					continue
				}

				target := &attack.AttackTarget{
					IP:      tts[0],
					Host:    tts[0],
					Port:    int(port),
					IsSSL:   isSSL,
					Version: "1.0",
					Proto:   "Test",
					App:     "Test",
					Types:   nil,
				}

				atargets <- target
			}
		}
	}()

	wg.Wait()
}
