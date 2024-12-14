package poco

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"sonic-agent-plus/internal/logger"
	"time"
)

var _ iPocoClient = (*webSocketClient)(nil)

type webSocketClient struct {
	port            int
	webSocketClient *websocket.Conn
	result          chan string
	driverType      PocoType
	isClose         bool
}

func newWebSocketClientImpl(port int, driverType PocoType) *webSocketClient {
	return &webSocketClient{
		port:       port,
		result:     make(chan string),
		driverType: driverType,
	}
}

func (w *webSocketClient) SendAndReceive(message *PocoRequestData) (string, error) {
	err := w.webSocketClient.WriteJSON(message)
	if err != nil {
		var closeErr *websocket.CloseError
		if errors.As(err, &closeErr) {
			w.isClose = true
		}
		return "", errors.New("the conn is close")
	}
	wait := 0
	for {
		select {
		case data := <-w.result:
			return data, nil
		default:
			time.Sleep(500 * time.Millisecond)
			wait++
			if wait >= 20 {
				return "", errors.New("poco ws not get result")
			}
		}
	}

}

func (w *webSocketClient) IsConnClose() bool {
	return w.isClose
}

func (w *webSocketClient) Dump() (string, error) {
	pocoMessage := NewPocoMessage()
	pocoMessage.Params = []interface{}{true}
	pocoMessage.Id = uuid.New().String()
	pocoMessage.Method = "Dump"
	if w.driverType == COCOS_2DX_JS || w.driverType == COCOS_CREATOR {
		pocoMessage.Method = "dump"
	}
	data, err := w.SendAndReceive(pocoMessage)
	if err != nil {
		return "", err
	}
	return data, nil
}

func (w *webSocketClient) Connect() error {
	url := fmt.Sprintf("ws://127.0.0.1:%d", w.port)
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}
	w.webSocketClient = c

	logger.Debug("poco ws connected.")

	w.isClose = false

	go func() {
		for {
			_, message, err := w.webSocketClient.ReadMessage()
			if err != nil {
				var closeErr *websocket.CloseError
				if errors.As(err, &closeErr) {
					w.isClose = true
				}
				break
			}
			w.result <- string(message)
		}
	}()
	return nil
}
func (w *webSocketClient) ReConnect() error {
	w.Disconnect()
	return w.Connect()
}

func (w *webSocketClient) Disconnect() {
	if w.webSocketClient != nil {
		w.webSocketClient.Close()
	}
	logger.Debug("poco ws closed.")
}
