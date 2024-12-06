package api

type IInputText interface {
	Start() error
	Stop() error
	InputEvent(data string) error
}
