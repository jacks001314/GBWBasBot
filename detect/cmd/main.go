package main

import (
	"common/scripts"
	"common/util/fileutils"
	"common/util/jsonutils"
	detect "detect/src"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
)

var (
	targetsCh = make(chan *detect.DTarget, 100)
	resultsCh = make(chan *detect.DResult, 100)
)

func makeTargets(wg *sync.WaitGroup, targets string, targetsF string) {

	defer wg.Done()

	mt := func(lines []string) {

		for _, line := range lines {

			ipport := strings.Split(strings.TrimSpace(line), ":")

			if port, err := strconv.Atoi(ipport[1]); err == nil {

				targetsCh <- &detect.DTarget{
					IP:   ipport[0],
					Port: port,
				}
			}
		}

	}

	if targets != "" {

		mt(strings.Split(targets, ","))
	}

	if targetsF != "" {

		if lines, err := fileutils.ReadAllLines(targetsF); err == nil {

			mt(lines)
		}
	}

}

func handleResults(wg *sync.WaitGroup) {

	defer wg.Done()

	for {

		select {
		case r := <-resultsCh:

			fmt.Println(jsonutils.ToJsonString(r, true))
		}
	}

}

func main() {

	threads := flag.Int("threads", 0, "set the thread number for detecting")

	waitnum := flag.Int("wnum", 0, "can be blocked and wait max number when too many detect threads are running")

	cleanTimeout := flag.Int("ctimeout", 0, "clean some timeout death threads when threads not working")

	stype := flag.String("stype", "lua", "the detect script type(lua/tengo)")

	spath := flag.String("spath", "", "the detect script path/dir")

	targets := flag.String("targets", "", "the detect targets,example:192.168.1.1:8080,192.168.1.2:8081,...")

	targetsF := flag.String("targetsF", "", "the detect targets from file")

	flag.Parse()

	if *spath == "" {

		flag.Usage()
		os.Exit(-1)
	}

	extName := ".lua"
	sstype := scripts.ScriptLua

	if *stype == "tengo" {

		extName = ".tengo"
		sstype = scripts.ScriptTengo
	}

	dtask, err := detect.NewDetectTask(*threads, *waitnum, *cleanTimeout, targetsCh, resultsCh)

	if err != nil {

		log.Fatalf("%v\n", err)

	}

	if strings.HasSuffix(*spath, extName) {

		//script file
		if err = dtask.AddDetectScriptFromFile(sstype, path.Base(*spath), *spath); err != nil {
			log.Fatalf("%v\n", err)
		}

	} else {

		//script dir
		dtask.AddDetectScriptFromDir(sstype, *spath, extName)
	}

	dtask.Start()
	defer dtask.Stop()

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go makeTargets(wg, *targets, *targetsF)

	go handleResults(wg)

	wg.Wait()
}
