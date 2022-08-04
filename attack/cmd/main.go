package main

import (
	attack "attack/core"
	"attack/target"
	"common/scripts"
	"common/util/jsonutils"
	tplutils "common/util/tpl"
	"fmt"
	"log"
	"sync"
)

func testSource() {

	//tfpath := `D:\shajf_dev\self\GBWBasBot\attack\scripts\target\tengo\attack_target_from_zoomeye.tpl`
	//tfname := "attack_target_from_zoomeye.tpl"
	lfpath := `D:\shajf_dev\self\GBWBasBot\attack\scripts\target\lua\attack_target_from_zoomeye.lua.tpl`
	lfname := "attack_target_from_zoomeye.lua.tpl"
	stype := scripts.ScriptLua

	data := &target.AttackTargetZoomEye{
		Key:     "97109Fa5-e22C-C30b3-265A-F99E70e4F33",
		Query:   `app:"Hadoop IPC"`,
		Port:    80,
		IsSSL:   false,
		Version: "1.0",
		Proto:   "http",
		App:     "weblogic",
	}

	content, err := tplutils.MakeSourceScriptFromTemplate(lfpath, lfname, data)

	if err != nil {

		log.Fatalf("%v", err)
	}

	tchan := make(chan *attack.AttackTarget, 1000)

	stask, err := target.NewSourceTask(tchan)

	if err != nil {

		log.Fatal(err)
	}

	wg := sync.WaitGroup{}

	wg.Add(2)

	go func() {

		err := stask.Fetch(stype, content, []string{"hadoop"})
		if err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		for {
			select {
			case t := <-tchan:
				fmt.Println(jsonutils.ToJsonString(t, true))
				//time.Sleep(2 * time.Second)

			}
		}
	}()

	wg.Wait()
}

func main() {

	testSource()
}
