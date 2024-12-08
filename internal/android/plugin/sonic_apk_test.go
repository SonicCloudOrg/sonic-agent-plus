package plugin_test

import (
	"github.com/stretchr/testify/assert"
	"sonic-agent-plus/internal/android/device"
	"sonic-agent-plus/internal/android/plugin"
	"sonic-agent-plus/internal/android/tool"
	"sonic-agent-plus/pkg/gadb"
	"testing"
	"time"
)

const (
	version = "2.0.8"
)

var (
	client gadb.Client
)

func SetClient() {
	client, _ = gadb.NewClient()
}

func TestSonicAndroidPlugin(t *testing.T) {
	SetClient()

	dev, err := tool.GetDevice(client, "")
	assert.NoError(t, err)

	sonicAndroidPlugin := plugin.NewSonicAndroidPlugin(
		"H:\\CodeProject\\GoProject\\sonic-agent-plus\\plugins\\sonic-android-apk.apk",
		version,
		device.NewAndroidDevice(dev),
	)
	err = sonicAndroidPlugin.Start()
	assert.NoError(t, err)

	time.Sleep(3 * time.Second)

	err = sonicAndroidPlugin.InputEvent("测试中文\n")
	assert.NoError(t, err)

	err = sonicAndroidPlugin.InputEvent("test abc\n")
	assert.NoError(t, err)

	err = sonicAndroidPlugin.InputEvent("123\n")
	assert.NoError(t, err)

	err = sonicAndroidPlugin.InputEvent("'`!@#$%^\t\n")
	assert.NoError(t, err)

	err = sonicAndroidPlugin.Stop()
	assert.NoError(t, err)
}
