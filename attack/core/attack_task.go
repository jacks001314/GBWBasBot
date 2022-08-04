package attack

import (
	"common/scripts"
	"common/util/jsonutils"
	"errors"
	"fmt"
	"io/ioutil"
	"runtime"
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

type AttackTask struct {
	threadsPool *ants.PoolWithFunc

	scripts sync.Map

	scriptNum int32

	attackProcesses chan *AttackProcess
	attackTargets   chan *AttackTarget
}

type AttackThreadContext struct {
	script Attack
	target *AttackTarget
}

func NewAttackTask(threads, maxWaitThreads, threadTimeout int, targets chan *AttackTarget, attackProcesses chan *AttackProcess) (*AttackTask, error) {

	tpool, err := ants.NewPoolWithFunc(getValue(threads, defaultThreads), func(ctx interface{}) {

		actx := ctx.(*AttackThreadContext)
		actx.script.Run(actx.target)
	},
		ants.WithNonblocking(false),
		ants.WithMaxBlockingTasks(getValue(maxWaitThreads, defaultMaxWaitThreads)),
		ants.WithExpiryDuration(time.Duration(getValue(threadTimeout, defaultThreadTimeout))*time.Second))

	if err != nil {

		return nil, err
	}

	return &AttackTask{

		threadsPool:     tpool,
		scripts:         sync.Map{},
		scriptNum:       0,
		attackTargets:   targets,
		attackProcesses: attackProcesses,
	}, nil

}

func (at *AttackTask) Publish(ap *AttackProcess) {

	at.attackProcesses <- ap
}

func getValue(v, d int) int {

	if v <= 0 {

		return d
	}

	return v
}

func (at *AttackTask) Start() {

	go func() {

		defer func() {

			if p := recover(); p != nil {

				var buf [4096]byte
				n := runtime.Stack(buf[:], false)

				log.Errorf("Attack Task Exit from panic:%s", string(buf[:n]))
			}
		}()

		for {

			select {

			case target := <-at.attackTargets:
				at.doAttack(target)

			}
		}
	}()
}

func (at *AttackTask) doAttack(target *AttackTarget) {

	at.scripts.Range(func(k interface{}, v interface{}) bool {

		ats := v.(Attack)

		actx := &AttackThreadContext{
			script: ats,
			target: target,
		}

		for {

			if err := at.threadsPool.Invoke(actx); err != nil {

				//wait more detect job ,sleep 1 second and again
				log.Warnf("Too many attack job to wait sleep 1 time again,info:%v", err)
				time.Sleep(time.Second)
				continue
			}

			break
		}

		return true
	})

}

func (at *AttackTask) Stop() {

	at.threadsPool.Release()
}

func getScriptTypeName(stype scripts.ScriptType) string {

	switch stype {

	case scripts.ScriptLua:
		return "lua"

	case scripts.ScriptTengo:
		return "tengo"

	default:
		return "unknown"
	}

}

func makeKey(config *Config) string {

	return fmt.Sprintf("%s_%d", config.Name, config.Id)
}

func (at *AttackTask) createScript(stype scripts.ScriptType, content []byte, config *Config, key string) (err error) {

	var attack Attack

	switch stype {

	case scripts.ScriptLua:

		if attack, err = LoadLuaScriptFromContent(at, content, config); err != nil {

			return
		}

	case scripts.ScriptTengo:

		if attack, err = LoadTengoScriptFromContent(at, content, config); err != nil {

			return
		}

	default:
		return fmt.Errorf("UnKnown Attack Script Type:%v", stype)

	}

	at.scripts.Store(key, attack)

	atomic.AddInt32(&at.scriptNum, 1)

	return nil
}

func (at *AttackTask) AddAttackScriptFromContent(stype scripts.ScriptType, content []byte, config *Config) (err error) {

	key := makeKey(config)

	if _, ok := at.scripts.Load(key); ok {
		//existed
		return errors.New("The script existed:" + key)
	}

	log.Debugf("Add a attack %s script from content ok!", getScriptTypeName(stype))

	return at.createScript(stype, content, config, key)
}

func (at *AttackTask) AddAttackScriptFromFile(stype scripts.ScriptType, config *Config) error {

	key := makeKey(config)

	if _, ok := at.scripts.Load(key); ok {
		//existed
		return errors.New("The script existed:" + key)
	}

	content, err := ioutil.ReadFile(config.FPath)

	if err != nil {

		return err
	}

	log.Debugf("Add a Attack %s script from filepath:%s ok!", getScriptTypeName(stype), config.FPath)

	return at.createScript(stype, content, config, key)
}

func (at *AttackTask) AddAttackScriptFromConfig(stype scripts.ScriptType, cfile string) error {

	var count uint32
	var cfg AttackScriptConfig

	if err := jsonutils.UNMarshalFromFile(&cfg, cfile); err != nil {

		log.Errorf("Load Attack Script Config failed:%v from file:%s", err, cfile)
		return err
	}

	for _, config := range cfg.scripts {

		if !config.Enable {

			continue
		}

		//file name as a script key
		if err := at.AddAttackScriptFromFile(stype, config); err != nil {

			log.Errorf("add Attack %s script from file:%s is failed:%v",
				getScriptTypeName(stype), config.FPath, err)

			continue
		}

		count++
	}

	log.Debugf("add attack %s script from config:%s is ok,add script number:%d", getScriptTypeName(stype), cfile, count)

	return nil
}

func (at *AttackTask) RemoveAttackScript(key string) {

	if _, ok := at.scripts.Load(key); ok {
		//existed
		at.scripts.Delete(key)

		atomic.AddInt32(&at.scriptNum, -1)

		log.Debugf("Remove a attack  script:%s ok!", key)
	}

}
