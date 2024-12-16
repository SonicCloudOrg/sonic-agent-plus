package entity

type TouchType string

const (
	TOUCH_TYPE_DOWN TouchType = "down"
	TOUCH_TYPE_MOVE TouchType = "move"
	TOUCH_TYPE_UP   TouchType = "up"
)

type TouchMode string

const (
	TOUCH_MODE_BY_AIRTEST TouchMode = "airtest"
	TOUCH_MODE_BY_DEFAULT TouchMode = "default"
)

func NewAritestPoint(x, y float32, fingerID int) TouchData {
	return TouchData{
		X:        x,
		Y:        y,
		FingerID: fingerID,
		Mode:     TOUCH_MODE_BY_AIRTEST,
	}
}

func NewDefaultPoint(x, y, fingerID int) TouchData {
	return TouchData{
		X:        float32(x),
		Y:        float32(y),
		FingerID: fingerID,
		Mode:     TOUCH_MODE_BY_DEFAULT,
	}
}

func NewAritestTouchData(x, y float32, fingerID int, touchType TouchType) TouchData {
	return TouchData{
		X:        x,
		Y:        y,
		Type:     touchType,
		FingerID: fingerID,
		Mode:     TOUCH_MODE_BY_AIRTEST,
	}
}

func NewDefalutTouchData(x, y, fingerID int, touchType TouchType) TouchData {
	return TouchData{
		X:        float32(x),
		Y:        float32(y),
		Type:     touchType,
		FingerID: fingerID,
		Mode:     TOUCH_MODE_BY_DEFAULT,
	}
}

type TouchData struct {
	X        float32   `json:"x"`
	Y        float32   `json:"y"`
	Type     TouchType `json:"type,omitempty"`
	Mode     TouchMode `json:"mode,omitempty"`
	FingerID int       `json:"fingerID"`
}
