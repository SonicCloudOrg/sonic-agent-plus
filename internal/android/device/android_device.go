package device

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/SonicCloudOrg/sonic-agent-plus/api"
	"github.com/SonicCloudOrg/sonic-agent-plus/pkg/gadb"
	json "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
)

var _ api.IDevice = (*AndroidDevice)(nil)

func NewAndroidDevice(device *gadb.Device) *AndroidDevice {
	return &AndroidDevice{
		dev: device,
	}
}

type AndroidDevice struct {
	dev *gadb.Device
}

func (d *AndroidDevice) GetDeviceInfo() (gjson.Result, error) {
	// todo 参数不全，待重写
	data := d.dev.DeviceInfo()
	bytes, err := json.Marshal(&data)
	if err != nil {
		return gjson.Parse(""), err
	}
	return gjson.Parse(string(bytes)), nil
}

func (d *AndroidDevice) Proxy(localPort int, remotePort string) error {
	return d.dev.FrowardTcp(localPort, remotePort)
}

func (d *AndroidDevice) RemoveProxy(localPort int, remotePort string) error {
	return d.dev.ForwardKill(localPort)
}

func (d *AndroidDevice) InstallApp(apkFile *os.File) error {
	return d.installApk(apkFile)
}

func (d *AndroidDevice) UnInstallApp(packageName string) error {
	return d.unInstallApk(packageName)
}

func (d *AndroidDevice) installApk(apkFile *os.File) error {
	remotePath := fmt.Sprintf("/data/local/tmp/%s", apkFile.Name())
	err := d.dev.PushFile(apkFile, remotePath, time.Now())
	if err != nil {
		return err
	}

	res, err := d.dev.RunShellCommand("pm install " + remotePath)
	if err != nil {
		return err
	}
	if strings.Contains(strings.ToLower(res), "success") {
		return nil
	}
	return fmt.Errorf("install %s error:%s", apkFile.Name(), res)
}

func (d *AndroidDevice) unInstallApk(packageName string) error {
	res, err := d.dev.RunShellCommand("pm uninstall " + packageName)
	if err != nil {
		return err
	}
	if strings.Contains(strings.ToLower(res), "success") {
		return nil
	}
	return fmt.Errorf("uninstall %s error:%s", packageName, res)
}

func (d *AndroidDevice) StartApp(appName string) error {
	cmd := fmt.Sprintf("monkey -p %s -c android.intent.category.LAUNCHER 1", appName)
	_, err := d.dev.RunShellCommand(cmd)
	return err
}

func (d *AndroidDevice) StopApp(appName string) error {
	cmd := fmt.Sprintf("adb shell am force-stop %s", appName)
	_, err := d.dev.RunShellCommand(cmd)
	return err
}

func (d *AndroidDevice) GetCurrentAppName() (string, error) {
	output, err := d.dev.RunShellCommand("dumpsys window windows")
	if err != nil {
		return "", errors.New("exec command error : dumpsys window windows")
	}
	return d.searchForCurrentPackage(output)
}

func (d *AndroidDevice) searchForCurrentPackage(output string) (string, error) {
	packageRE := regexp.MustCompile(`\s*mCurrentFocus=Window{.* ([A-Za-z0-9_.]+)/[A-Za-z0-9_.]+}`)
	matches := packageRE.FindStringSubmatch(output)
	if len(matches) > 1 {
		return matches[len(matches)-1], nil
	} else {
		appName, _, err := d.GetCurrentPackageNameAndPid()
		return appName, err
	}
}

func (d *AndroidDevice) GetCurrentPackageNameAndPid() (packageName string, pid string, err error) {
	data, err := d.dev.RunShellCommand("dumpsys activity top | grep ACTIVITY")
	if err != nil {
		return "", "", fmt.Errorf("exec command error : " + "dumpsys activity top | grep ACTIVITY")
	}

	var dataSplit []string

	dataSplitTemp := strings.Split(data, "\n")

	for _, lineStr := range dataSplitTemp {
		if lineStr != "" {
			dataSplit = append(dataSplit, lineStr)
		}
	}

	currentActivityLineStr := strings.TrimLeft(dataSplit[len(dataSplit)-1], " ")

	dataSplit = strings.Split(currentActivityLineStr, " ")
	if len(dataSplit) > 1 {
		packageNameSplit := strings.Split(dataSplit[1], "/")
		for _, param := range dataSplit {
			if strings.Contains(param, "pid=") {
				return packageNameSplit[0], strings.ReplaceAll(param, "pid=", ""), nil
			}
		}
		return packageNameSplit[0], "", nil
	} else {
		return "", "", errors.New("not get current package name")
	}
}

func (d *AndroidDevice) KeyCode(keyCode int) error {
	_, err := d.dev.RunShellCommand(fmt.Sprintf("input keyevent %d", keyCode))
	return err
}

func (d *AndroidDevice) ExecuteCommand(cmd string) (net.Conn, error) {
	return d.dev.RunShellLoopCommandSock(cmd)
}

func (d *AndroidDevice) PushFile(local *os.File, remotePath string) error {
	return d.dev.PushFile(local, remotePath)
}

func (d *AndroidDevice) Pull(remotePath string, dest io.Writer) error {
	return d.dev.Pull(remotePath, dest)
}
