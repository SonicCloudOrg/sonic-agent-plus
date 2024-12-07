package api

import "sonic-agent-plus/entity"

type ITouch interface {
	Start() error
	Stop() error
	//Tap param duration unit ms
	Tap(data entity.TouchData, duration int64) error
	//Swipe param waitTime:touch wait time,duration: unit ms
	Swipe(startPoint entity.TouchData, endPoint entity.TouchData, waitTime, duration int64) error
	TouchEvent(data entity.TouchData) error
}
