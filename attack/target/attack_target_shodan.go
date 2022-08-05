package target

type AttackTargetShodan struct {
	Key   string
	Query string

	Port int

	IsSSL bool

	Version string

	Proto string

	App string
}
