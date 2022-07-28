package detect

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/panjf2000/ants/v2"
)

const (
	defaultThreads        = 100
	defaultMaxWaitThreads = 100
	defaultThreadTimeout  = 60
)

type ScriptType int

const (
	ScriptLua ScriptType = iota + 1
	ScriptTengo
)

type DTask struct {
	threadsPool *ants.PoolWithFunc

	scripts sync.Map

	scriptNum int32

	targets chan *DTarget

	results chan *DResult
}

type DetectThreadContext struct {
	dt     Detect
	target *DTarget
}

func getValue(v, d int) int {

	if v <= 0 {

		return d
	}

	return v
}

func NewDetectTask(threads, maxWaitThreads, threadTimeout int, targets chan *DTarget, results chan *DResult) (*DTask, error) {

	tpool, err := ants.NewPoolWithFunc(getValue(threads, defaultThreads), func(ctx interface{}) {

		dctx := ctx.(*DetectThreadContext)

		dctx.dt.Run(dctx.target)

	},
		ants.WithNonblocking(false),
		ants.WithMaxBlockingTasks(getValue(maxWaitThreads, defaultMaxWaitThreads)),
		ants.WithExpiryDuration(time.Duration(getValue(threadTimeout, defaultThreadTimeout))*time.Second))

	if err != nil {

		return nil, err
	}

	return &DTask{
		threadsPool: tpool,
		scripts:     sync.Map{},
		scriptNum:   0,
		targets:     targets,
		results:     results,
	}, nil

}

func (d *DTask) Start() {

	go func() {

		defer func() {

			if p := recover(); p != nil {

				var buf [4096]byte
				n := runtime.Stack(buf[:], false)

				log.Errorf("Detect Task Exit from panic:%s", string(buf[:n]))
			}
		}()

		for {

			select {

			case target := <-d.targets:

				d.runDetect(target)
			}
		}
	}()
}

func (d *DTask) runDetect(target *DTarget) {

	d.scripts.Range(func(k interface{}, v interface{}) bool {

		dt := v.(Detect)

		dctx := &DetectThreadContext{
			dt:     dt,
			target: target,
		}

		for {

			if err := d.threadsPool.Invoke(dctx); err != nil {

				//wait more detect job ,sleep 1 second and again
				log.Warnf("Too many detect job to wait sleep 1 time again,info:%v", err)
				time.Sleep(time.Second)
				continue
			}

			break
		}

		return true
	})

}

func (d *DTask) Stop() {

	d.threadsPool.Release()
}

func getScriptTypeName(stype ScriptType) string {

	switch stype {

	case ScriptLua:
		return "lua"

	case ScriptTengo:
		return "tengo"

	default:
		return "unknown"
	}

}

func (d *DTask) createScript(stype ScriptType, key string, content []byte) (err error) {

	var dt Detect

	switch stype {

	case ScriptLua:

		if dt, err = LoadLuaScriptFromContent(d, content, key); err != nil {

			return
		}

	case ScriptTengo:

		if dt, err = LoadTengoScriptFromContent(d, content, key); err != nil {

			return
		}

	default:
		return fmt.Errorf("UnKnown Detect Script Type:%v", stype)

	}

	d.scripts.Store(key, dt)
	atomic.AddInt32(&d.scriptNum, 1)

	return nil
}

func (d *DTask) AddDetectScriptFromContent(stype ScriptType, key string, content []byte) (err error) {

	if _, ok := d.scripts.Load(key); ok {
		//existed
		return errors.New("The script existed:" + key)
	}

	log.Debugf("Add a detect %s script from content ok!", getScriptTypeName(stype))

	return d.createScript(stype, key, content)
}

func (d *DTask) AddDetectScriptFromFile(stype ScriptType, key string, fpath string) error {

	if _, ok := d.scripts.Load(key); ok {
		//existed
		return errors.New("The script existed:" + key)
	}

	content, err := ioutil.ReadFile(fpath)

	if err != nil {

		return err
	}

	log.Debugf("Add a detect %s script from filepath:%s ok!", getScriptTypeName(stype), fpath)

	return d.createScript(stype, key, content)
}

func (d *DTask) AddDetectScriptFromDir(stype ScriptType, fdir, extName string) {

	var count uint32

	filepath.Walk(fdir, func(path string, info os.FileInfo, err error) error {

		if !info.IsDir() {

			if strings.HasSuffix(path, extName) {

				//file name as a script key
				if err := d.AddDetectScriptFromFile(stype, info.Name(), path); err != nil {

					log.Errorf("add detect %s script from file:%s is failed:%v",
						getScriptTypeName(stype), path, err)

				} else {
					count++
				}

			}
		}

		return nil
	})

	log.Debugf("add detect %s script from dir:%s is ok,add script number:%d", getScriptTypeName(stype), fdir, count)

}

func (d *DTask) RemoveDetectScript(key string) {

	if _, ok := d.scripts.Load(key); ok {
		//existed
		d.scripts.Delete(key)

		atomic.AddInt32(&d.scriptNum, -1)

		log.Debugf("Remove a detect  script:%s ok!", key)
	}

}

func (d *DTask) Publish(result *DResult) {

	d.results <- result
}
