package target

type AttackTargetIPRange struct {
	WhiteLists []string

	BlackLists []string

	Port int

	IsSSL bool

	Version string

	Proto string

	App string
}
