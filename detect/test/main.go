package main

import (
	detect "detect/src"
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/panjf2000/ants/v2"
)

type DetectContext struct {
	ds *detect.DLuaScript
	t  *detect.DTarget
}

func testLua() {

	fpath := `D:\shajf_dev\self\GBWBasBot\detect\test\detect.lua`

	ds, err := detect.LoadLuaScriptFromFile(fpath, "detect.lua")

	if err != nil {

		log.Fatal(err)

	}

	wg := sync.WaitGroup{}

	p, _ := ants.NewPoolWithFunc(100, func(dcxt interface{}) {

		cxt := dcxt.(*DetectContext)
		cxt.ds.Run(cxt.t)
		wg.Done()
	})
	defer p.Release()

	for i := 0; i < 10; i++ {

		wg.Add(1)
		ip := "192.168.1." + strconv.Itoa(i+1)

		p.Invoke(&DetectContext{
			ds: ds,
			t: &detect.DTarget{
				IP:   ip,
				Port: uint16(i + 80),
			},
		})
	}

	fmt.Println(p.Running())
	wg.Wait()
}

func main() {

	testLua()
}
