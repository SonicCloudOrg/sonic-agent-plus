package api

import (
	"github.com/SonicCloudOrg/sonic-agent-plus/entity"
	"github.com/tidwall/gjson"
	"io"
	"net"
	"os"
)

type Device interface {
	GetDeviceInfo() gjson.Result
	Proxy(localPort int, remotePort string) error
	RemoveProxy(localPort int, remotePort string) error
	InstallApp(apkFile *os.File) error
	UnInstallApp(packageName string) error
	StartApp(appName string) error
	StopApp(appName string) error
	Touch(touch entity.TouchData) error
	GetCurrentAppName() string
	KeyCode(keyCode int) error
	ExecuteCommand(cmd string) (net.Conn, error)
	PushFile(local *os.File, remotePath string) error
	Pull(remotePath string, dest io.Writer) error
}
