package target

type AttackTargetZoomEye struct {
	Key   string
	Query string

	Port int

	IsSSL bool

	Version string

	Proto string

	App string
}
