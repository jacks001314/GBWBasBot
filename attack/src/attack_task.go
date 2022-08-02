package attack

type AttackTask struct {
	attackProcesses chan *AttackProcess
}

func (at *AttackTask) Publish(ap *AttackProcess) {

	at.attackProcesses <- ap
}
