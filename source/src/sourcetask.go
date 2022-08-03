package source

import (
	"common/scripts"
	"runtime"
	"time"

	"github.com/panjf2000/ants/v2"
)

const (
	defaultThreads         = 100
	defaultMaxWaitThreads  = 100
	defaultThreadTimeout   = 60
	defaultTargetsChanSize = 1000
)

type SourceTask struct {
	threadsPool *ants.PoolWithFunc
	targets     chan interface{}
}

func NewSourceTask(threads, maxWaitThreads, threadTimeout, targetsChanSize int) (*SourceTask, error) {

	tpool, err := ants.NewPoolWithFunc(getValue(threads, defaultThreads), func(s interface{}) {

		ss := s.(Source)
		ss.Start()
	},
		ants.WithNonblocking(false),
		ants.WithMaxBlockingTasks(getValue(maxWaitThreads, defaultMaxWaitThreads)),
		ants.WithExpiryDuration(time.Duration(getValue(threadTimeout, defaultThreadTimeout))*time.Second))

	if err != nil {

		return nil, err
	}

	return &SourceTask{
		threadsPool: tpool,
		targets:     make(chan interface{}, getValue(targetsChanSize, defaultTargetsChanSize)),
	}, nil

}

func getValue(v, d int) int {

	if v <= 0 {

		return d
	}

	return v
}

func (st *SourceTask) Fetch(stype scripts.ScriptType, content []byte) {

	go func() {

		defer func() {

			if p := recover(); p != nil {

				var buf [4096]byte
				n := runtime.Stack(buf[:], false)

				log.Errorf("Source Task Exit from panic:%s", string(buf[:n]))
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
