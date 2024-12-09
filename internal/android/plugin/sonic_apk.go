package plugin

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"os"
	"sonic-agent-plus/define/api"
	"sonic-agent-plus/internal/android/device"
	"sonic-agent-plus/internal/logger"
	"strings"
)

var _ api.IInputText = (*SonicAndroidPlugin)(nil)

func NewSonicAndroidPlugin(pluginPath string, sonicApkVersion string, dev *device.AndroidDevice) *SonicAndroidPlugin {
	return &SonicAndroidPlugin{
		pluginPath:      pluginPath,
		sonicApkVersion: sonicApkVersion,
		dev:             dev,
	}
}

type SonicAndroidPlugin struct {
	pluginPath      string
	pluginFile      *os.File
	sonicApkVersion string
	dev             *device.AndroidDevice
}

func (s *SonicAndroidPlugin) Start() error {
	var err error
	s.pluginFile, err = os.Open(s.pluginPath)
	if err != nil {
		return err
	}
	if !s.installSonicApk() {
		return errors.New("sonic apk start fail")
	}

	_, err = s.dev.ExecuteCommand("am start -n org.cloud.sonic.android/.SonicServiceActivity")
	if err != nil {
		return err
	}
	currentIme, _ := s.dev.ExecuteCommand("settings get secure default_input_method")
	if !strings.Contains(currentIme, "org.cloud.sonic.android/.keyboard.SonicKeyboard") {
		_, err = s.dev.ExecuteCommand("ime enable org.cloud.sonic.android/.keyboard.SonicKeyboard")
		if err != nil {
			return err
		}
		_, err = s.dev.ExecuteCommand("ime set org.cloud.sonic.android/.keyboard.SonicKeyboard")
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SonicAndroidPlugin) checkVersion() bool {
	res, err := s.dev.ExecuteCommand("dumpsys package org.cloud.sonic.android")
	if err != nil {
		return false
	}
	return strings.Contains(res, "versionName="+s.sonicApkVersion)
}

func (s *SonicAndroidPlugin) installSonicApk() bool {
	path, err := s.dev.ExecuteCommand("pm path org.cloud.sonic.android")
	if err != nil {
		logger.Error("Error executing command: pm path org.cloud.sonic.android", zap.Error(err))
		return false
	}
	path = strings.TrimSpace(strings.ReplaceAll(path, "package:", ""))
	path = strings.ReplaceAll(path, "\n", "")
	path = strings.ReplaceAll(path, "\t", "")

	logger.Info("Check Sonic Apk version and status pass...")

	if len(path) > 0 && s.checkVersion() {
		logger.Info("Sonic apk check passed")
		return true
	} else {
		logger.Info("Sonic Apk version not newest or not installed, starting install...")

		err = s.dev.UnInstallApp("org.cloud.sonic.android")
		if err != nil {
			logger.Error("Uninstall sonic Apk error", zap.Error(err))
		}

		err = s.dev.InstallApp(s.pluginFile)
		if err != nil {
			logger.Error("Sonic Apk install failed", zap.Error(err))
			return false
		}
		_, err = s.dev.ExecuteCommand("appops set org.cloud.sonic.android POST_NOTIFICATION allow")
		if err != nil {
			logger.Error("Setting appops POST_NOTIFICATION failed", zap.Error(err))
		}
		_, err = s.dev.ExecuteCommand("appops set org.cloud.sonic.android RUN_IN_BACKGROUND allow")
		if err != nil {
			logger.Error("Setting appops RUN_IN_BACKGROUND failed", zap.Error(err))
		}
		_, err = s.dev.ExecuteCommand("dumpsys deviceidle whitelist +org.cloud.sonic.android")
		if err != nil {
			logger.Error("Adding Sonic APK to deviceidle whitelist failed", zap.Error(err))
		}
		logger.Debug("Sonic Apk install successful.")
		return true
	}
}

func (s *SonicAndroidPlugin) Stop() error {
	_, err := s.dev.ExecuteCommand("ime reset")
	return err
}

func (s *SonicAndroidPlugin) InputEvent(data string) error {
	//_, err := s.dev.ExecuteCommand(fmt.Sprintf("am broadcast -a SONIC_KEYBOARD --es msg \"%s\"", "CODE_AC_CLEAN"))
	//if err != nil {
	//	return err
	//}
	_, err := s.dev.ExecuteCommand(fmt.Sprintf("am broadcast -a SONIC_KEYBOARD --es msg \"%s\"", data))
	if err != nil {
		return err
	}
	_, err = s.dev.ExecuteCommand(fmt.Sprintf("am broadcast -a SONIC_KEYBOARD --es msg \"%s\"", "CODE_AC_ENTER"))
	return err
}
