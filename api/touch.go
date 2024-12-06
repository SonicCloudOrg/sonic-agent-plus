package api

import "github.com/SonicCloudOrg/sonic-agent-plus/entity"

type ITouch interface {
	Start() error
	Stop() error
	TouchEvent(data entity.TouchData) error
}
