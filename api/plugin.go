package api

type Plugin interface {
	Start() error
	Stop() error
}
