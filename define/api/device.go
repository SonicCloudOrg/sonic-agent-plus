package api

import (
	"github.com/tidwall/gjson"
	"io"
	"net"
	"os"
)

type IDevice interface {
	GetDeviceInfo() (gjson.Result, error)
	Proxy(localPort int, remotePort string) error
	RemoveProxy(localPort int, remotePort string) error
	InstallApp(apkFile io.Reader) error
	UnInstallApp(packageName string) error
	StartApp(appName string) error
	StopApp(appName string) error
	GetCurrentAppName() (string, error)
	KeyCode(keyCode int) error
	ExecuteCommand(cmd string) (string, error)
	ExecuteNohupCommand(cmd string) (net.Conn, error)
	PushFile(local *os.File, remotePath string) error
	Pull(remotePath string, dest io.Writer) error
}
