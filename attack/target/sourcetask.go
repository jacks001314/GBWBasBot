package target

import (
	attack "attack/core"
	"common/scripts"
	"fmt"
	"io/ioutil"
	"sync"
)

type SourceTask struct {
	wg      *sync.WaitGroup
	targets chan *attack.AttackTarget
}

func NewSourceTask(targets chan *attack.AttackTarget) (*SourceTask, error) {

	return &SourceTask{
		targets: targets,
	}, nil

}

func (st *SourceTask) Fetch(stype scripts.ScriptType, content []byte, types []string) error {

	var s Source
	var err error

	switch stype {

	case scripts.ScriptLua:

		if s, err = LoadLuaScriptFromContent(st, content, types); err != nil {

			return err
		}

	case scripts.ScriptTengo:

		if s, err = LoadTengoScriptFromContent(st, content, types); err != nil {

			return err
		}

	default:

		return fmt.Errorf("UnKnown Attack Source Script Type:%v", stype)
	}

	s.Start()

	return nil
}

func (st *SourceTask) FetchFromFile(stype scripts.ScriptType, content []byte, types []string, fpath string) error {

	content, err := ioutil.ReadFile(fpath)

	if err != nil {

		return err
	}

	return st.Fetch(stype, content, types)
}

func (st *SourceTask) Put(target *attack.AttackTarget) {

	st.targets <- target
}

func (st *SourceTask) CloseSource(s Source) {

}
