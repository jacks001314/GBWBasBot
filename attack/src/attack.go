package attack

type Attack interface {
	Accept(target *AttackTarget) bool

	Run(target *AttackTarget) error

	Publish(process *AttackProcess)
}
