package api

type IPlugin interface {
	Start() error
	Stop() error
}
