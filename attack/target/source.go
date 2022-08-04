package target

import attack "attack/core"

type Source interface {
	AttackTypes() []string

	Put(target *attack.AttackTarget) error

	Start() error

	Stop()

	AtEnd()
}
