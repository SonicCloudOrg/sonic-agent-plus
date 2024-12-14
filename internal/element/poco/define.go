package poco

import (
	"encoding/binary"
	jsoniter "github.com/json-iterator/go"
)

type iPocoClient interface {
	SendAndReceive(message *PocoRequestData) (string, error)
	Dump() (string, error)
	Connect() error
	ReConnect() error
	IsConnClose() bool
	Disconnect()
}

type PocoRequestData struct {
	Jsonrpc string        `json:"jsonrpc"`
	Params  []interface{} `json:"params"`
	Id      string        `json:"id"`
	Method  string        `json:"method"`
}

func NewPocoMessage() *PocoRequestData {
	return &PocoRequestData{
		Jsonrpc: "2.0",
	}
}

func (m *PocoRequestData) ToJson() string {
	data, err := jsoniter.Marshal(m)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func intToByteArray(i int) []byte {
	result := make([]byte, 4)
	binary.LittleEndian.PutUint32(result, uint32(i))
	return result
}

func toInt(b []byte) int {
	return int(binary.LittleEndian.Uint32(b))
}

type PocoType string

const (
	Android            PocoType = "Android"
	UNITY_3D           PocoType = "UNITY_3D"
	UE4                PocoType = "UE4"
	COCOS_2DX_JS       PocoType = "COCOS_2DX_JS"
	COCOS_2DX_LUA      PocoType = "COCOS_2DX_LUA"
	COCOS_2DX_C_PLUS_1 PocoType = "COCOS_2DX_C_PLUS_1"
	COCOS_CREATOR      PocoType = "COCOS_CREATOR"
	EGRET              PocoType = "EGRET"
)

func getPocoDefaultPortByName(pocoType PocoType) int {
	switch pocoType {
	case Android:
		return 6790
	case UE4:
		return 5001
	case UNITY_3D:
		return 5001
	case EGRET:
		return 5003
	case COCOS_CREATOR:
		return 5003
	case COCOS_2DX_JS:
		return 5003
	case COCOS_2DX_LUA:
		return 15004
	case COCOS_2DX_C_PLUS_1:
		return 18888
	}
	return -1
}

func iSPoco(pocoType PocoType) bool {
	switch pocoType {
	case UNITY_3D:
		return true
	case UE4:
		return true
	case COCOS_2DX_JS:
		return true
	case COCOS_2DX_LUA:
		return true
	case COCOS_2DX_C_PLUS_1:
		return true
	case COCOS_CREATOR:
		return true
	case EGRET:
		return true
	}
	return false
}
