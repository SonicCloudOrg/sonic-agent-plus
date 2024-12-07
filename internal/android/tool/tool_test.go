package tool_test

import (
	"fmt"
	"sonic-agent-plus/internal/android/tool"
	"sonic-agent-plus/pkg/gadb"
	"strings"
	"testing"
)

var (
	client gadb.Client
)

func SetClient() {
	client, _ = gadb.NewClient()
}

func TestGetDevice(t *testing.T) {
	SetClient()
	device, err := tool.GetDevice(client, "231341")
	if err != nil {
		panic(err)
	}
	fmt.Println(device.Serial())
	data, _ := device.RunShellCommand("ps")
	fmt.Println(data)
}

func TestGetPackageNameList(t *testing.T) {
	SetClient()
	device, err := tool.GetDevice(client, "")
	if err != nil {
		panic(err)
	}
	packageList, err := tool.GetPackageNameList(device)
	if err != nil {
		panic(err)
	}
	fmt.Println(strings.Join(packageList, "\n"))
}

func TestGetCurrentPackageName(t *testing.T) {
	SetClient()
	device, err := tool.GetDevice(client, "")
	if err != nil {
		panic(err)
	}
	packageName, pid, err := tool.GetCurrentPackageNameAndPid(device)
	if err != nil {
		panic(err)
	}
	fmt.Println(packageName, pid)
}
