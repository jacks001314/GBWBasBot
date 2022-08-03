package scripts

type ScriptType int

const (
	ScriptLua ScriptType = iota + 1
	ScriptTengo
)
