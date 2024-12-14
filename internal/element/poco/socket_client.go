package poco

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net"
	"sonic-agent-plus/internal/logger"
	"sync"
	"time"
)

var _ iPocoClient = (*socketClient)(nil)

type socketClient struct {
	lock       sync.Mutex
	poco       net.Conn
	driverType PocoType
	port       int
	host       string
	isClose    bool
}

func newSocketClientImpl(port int, driverType PocoType) *socketClient {
	return &socketClient{
		port:       port,
		host:       "127.0.0.1",
		driverType: driverType,
	}
}

func (s *socketClient) SendAndReceive(message *PocoRequestData) (string, error) {
	data := message.ToJson()
	header := intToByteArray(len(data))
	s.lock.Lock()
	defer s.lock.Unlock()

	_, err := s.poco.Write(header)
	if err != nil {
		return "", err
	}

	_, err = s.poco.Write([]byte(data))
	if err != nil {
		return "", err
	}

	head := make([]byte, 4)
	_, err = s.poco.Read(head)
	if err != nil {
		return "", err
	}

	headLen := toInt(head)
	rData := make([]byte, 0)
	for {
		buffer := make([]byte, 8192)
		realLen, err := s.poco.Read(buffer)
		if err != nil {
			return "", err
		}

		if realLen > 0 {
			if realLen < headLen {
				rData = append(rData, buffer[:realLen]...)
			} else {
				rData = append(rData, buffer[:headLen]...)
			}
		}
		//fmt.Println(string(rData))
		if len(rData) == headLen {
			return string(rData), nil
		}
	}
}

func (s *socketClient) Dump() (string, error) {
	pocoMessage := NewPocoMessage()
	pocoMessage.Params = []interface{}{true}
	pocoMessage.Id = uuid.New().String()
	pocoMessage.Method = "Dump"
	data, err := s.SendAndReceive(pocoMessage)
	if err != nil {
		if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
			// If it's a timeout error, the connection is probably still open.
			s.isClose = false
		}
		s.isClose = true
		return "", err
	}
	return data, nil
}

func (s *socketClient) Connect() error {
	for i := 0; i < 20; i++ {
		poco, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", s.host, s.port), 100*time.Millisecond)
		if err == nil {
			s.poco = poco
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	if s.poco != nil {
		logger.Debug("poco socket connected.")
		s.isClose = false
		return nil
	} else {
		return errors.New("poco socket disconnected")
	}
}

func (s *socketClient) ReConnect() error {
	s.Disconnect()
	return s.Connect()
}

func (s *socketClient) IsConnClose() bool {
	return s.isClose
}

func (s *socketClient) Disconnect() {
	if s.poco != nil {
		s.poco.Close()
	}
}
