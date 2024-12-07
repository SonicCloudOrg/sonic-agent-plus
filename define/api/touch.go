package api

import "sonic-agent-plus/entity"

type ITouch interface {
	Start() error
	Stop() error
	TouchEvent(data entity.TouchData) error
}
