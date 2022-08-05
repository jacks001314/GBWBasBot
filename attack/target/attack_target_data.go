package target

type AttackTargetData struct {
	WhiteLists []string

	BlackLists []string

	Email string

	Key string

	Query string

	Port int

	IsSSL bool

	Version string

	Proto string

	App string
}
