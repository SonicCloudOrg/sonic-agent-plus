package touch_test

import (
	"github.com/stretchr/testify/assert"
	"sonic-agent-plus/entity"
	"sonic-agent-plus/internal/android/tool"
	"sonic-agent-plus/internal/android/touch"
	"sonic-agent-plus/pkg/gadb"
	"testing"
)

var (
	client gadb.Client
)

func SetClient() {
	client, _ = gadb.NewClient()
}

func TestTouchTap(t *testing.T) {
	SetClient()
	device, err := tool.GetDevice(client, "")
	if err != nil {
		panic(err)
	}

	androidTouch := touch.NewAndroidTouch(device)
	err = androidTouch.Start()
	assert.NoError(t, err)

	err = androidTouch.Tap(entity.NewAritestPoint(0.5, 0.5, 1), 100)
	assert.NoError(t, err)

	err = androidTouch.Stop()
	assert.NoError(t, err)
}
func TestTouchSwipe(t *testing.T) {
	SetClient()
	device, err := tool.GetDevice(client, "")
	if err != nil {
		panic(err)
	}

	androidTouch := touch.NewAndroidTouch(device)
	err = androidTouch.Start()
	assert.NoError(t, err)

	err = androidTouch.Swipe(
		entity.NewAritestPoint(0.1, 0.5, 1),
		entity.NewAritestPoint(0.9, 0.5, 1), 100,
		1000)
	assert.NoError(t, err)

	err = androidTouch.Stop()
	assert.NoError(t, err)
}
