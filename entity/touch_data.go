package entity

type TouchType string

const (
	TOUCH_DOWN = "down"
	TOUCH_MOVE = "move"
	TOUCH_UP   = "up"
)

type TouchMode string

const (
	TouchAritestMode TouchMode = "airtest"
	TouchDefaultMode TouchMode = "default"
)

type TouchData struct {
	X         float32   `json:"x"`
	Y         float32   `json:"y"`
	TouchType TouchType `json:"touchType"`
	TouchMode TouchMode `json:"touchMode"`
	FingerID  int       `json:"fingerID"`
}
