package touch

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"sonic-agent-plus/define/api"
	"sonic-agent-plus/define/common_error"
	"sonic-agent-plus/define/lib"
	"sonic-agent-plus/entity"
	"sonic-agent-plus/internal/logger"
	"sonic-agent-plus/pkg/gadb"
	"strings"
	"sync"
	"time"
)

var _ api.ITouch = (*AndroidTouch)(nil)

const (
	androidTouchAritestMode = "airtest"
	androidTouchDefaultMode = "touch"
)

func NewAndroidTouch(device *gadb.Device) *AndroidTouch {
	return &AndroidTouch{
		device: device,
	}
}

type AndroidTouch struct {
	device    *gadb.Device
	touchConn net.Conn
	touchLock sync.Mutex
}

func (t *AndroidTouch) Tap(data entity.TouchData, duration int64) error {
	var err error
	data.Type = entity.TOUCH_TYPE_DOWN
	err = t.TouchEvent(data)
	if err != nil {
		return err
	}
	time.Sleep(time.Millisecond * time.Duration(duration))
	data.Type = entity.TOUCH_TYPE_UP
	err = t.TouchEvent(data)
	if err != nil {
		return err
	}
	return nil
}

func (t *AndroidTouch) Swipe(startPoint entity.TouchData, endPoint entity.TouchData, waitTime, duration int64) error {
	var err error
	startPoint.Type = entity.TOUCH_TYPE_DOWN
	var touchMode = startPoint.Mode
	var fingerID = startPoint.FingerID

	err = t.TouchEvent(startPoint)
	if err != nil {
		return err
	}

	lastTime := time.Now().UnixMilli()

	for {
		nowTime := time.Now().UnixMilli()
		timeProgress := float32(nowTime-lastTime) / float32(duration)
		if timeProgress >= 1.0 {
			endPoint.Type = entity.TOUCH_TYPE_UP
			err = t.TouchEvent(endPoint)
			if err != nil {
				return err
			}
			break
		}
		transitionX := startPoint.X + (endPoint.X-startPoint.X)*timeProgress
		transitionY := startPoint.Y + (endPoint.Y-startPoint.Y)*timeProgress
		err = t.TouchEvent(entity.TouchData{
			X:        transitionX,
			Y:        transitionY,
			Type:     entity.TOUCH_TYPE_MOVE,
			Mode:     touchMode,
			FingerID: fingerID,
		})
		if err != nil {
			return err
		}
		time.Sleep(time.Duration(waitTime) * time.Millisecond)
	}

	return nil
}

func (t *AndroidTouch) Start() error {
	var err error
	if t.device == nil {
		return common_error.DeviceNullError
	}
	err = t.device.Push(bytes.NewReader(lib.AndroidTouchJarBytes), lib.RemoteTouchToolPath, time.Now())
	if err != nil {
		return err
	}
	t.touchConn, err = t.device.RunShellLoopCommandSock(fmt.Sprintf(
		"CLASSPATH=%s app_process / com.aoliaoaojiao.AndroidTouch.Run v2.2",
		lib.RemoteTouchToolPath))
	if err != nil {
		return err
	}
	return t.checkJarStart()
}

func (t *AndroidTouch) checkJarStart() error {
	var byteDatas = make([]byte, 1024)
	n, err := t.touchConn.Read(byteDatas)
	if err != nil {
		return errors.New("start android touch err: " + err.Error())

	}
	if !strings.Contains(string(byteDatas[:n]), "Device") {
		return errors.New("not start android touch: " + string(byteDatas[:n]))
	}
	return nil
}

func (t *AndroidTouch) Stop() error {
	return t.touchConn.Close()
}

func (t *AndroidTouch) TouchEvent(data entity.TouchData) error {
	cmd := ""
	if data.Mode == entity.TOUCH_MODE_BY_DEFAULT {
		cmd = t.genDefaultTouch(data)
	} else {
		cmd = t.genAirtestCmd(data)
	}
	return t.executeTouchCmd(cmd)
}

func (t *AndroidTouch) genDefaultTouch(data entity.TouchData) string {
	if data.Type != entity.TOUCH_TYPE_UP {
		return fmt.Sprintf(
			"%s %s %d %d %d\n",
			androidTouchDefaultMode,
			data.Type,
			int(data.X),
			int(data.Y),
			data.FingerID)
	} else {
		return fmt.Sprintf(
			"%s up %d\n",
			androidTouchDefaultMode,
			data.FingerID)
	}
}

func (t *AndroidTouch) genAirtestCmd(data entity.TouchData) string {
	if data.Type != entity.TOUCH_TYPE_UP {
		return fmt.Sprintf(
			"%s %s %f %f %d\n",
			androidTouchAritestMode,
			data.Type,
			data.X,
			data.Y,
			data.FingerID)
	} else {
		return fmt.Sprintf(
			"%s up %d\n",
			androidTouchAritestMode,
			data.FingerID)
	}
}

func (t *AndroidTouch) executeTouchCmd(cmd string) error {
	t.touchLock.Lock()
	defer t.touchLock.Unlock()
	logger.Debug(cmd)
	now := time.Now()
	t.touchConn.SetReadDeadline(now.Add(time.Second * 2))
	_, err := t.touchConn.Write([]byte(cmd))
	if err != nil {
		return err
	}
	var byteDatas = make([]byte, 1024)
	n, err := t.touchConn.Read(byteDatas)
	if err != nil {
		return errors.New("touch fail:" + cmd)
	}
	if !strings.Contains(string(byteDatas[:n]), "succeed") {
		return errors.New("touch result fail:" + cmd)
	}
	return nil
}
