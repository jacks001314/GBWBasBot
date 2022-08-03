package source

type Source interface {
	Put(target interface{}) error

	Start() error

	Stop()

	AtEnd()
}
