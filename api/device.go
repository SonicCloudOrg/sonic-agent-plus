package api

import "github.com/SonicCloudOrg/sonic-agent-plus/entity"

type Device interface {
	Proxy(localPort int, remotePort string) error
	RemoveProxy(localPort int, remotePort string) error
	StartApp(appName string) error
	StopApp(appName string) error
	Touch(touch entity.TouchData) error
	GetCurrentAppName() string
	KeyCode(keyCode int) error
}
