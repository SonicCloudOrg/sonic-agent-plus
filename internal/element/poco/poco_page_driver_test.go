package poco_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"sonic-agent-plus/internal/android/tool"
	"sonic-agent-plus/internal/element/poco"
	"sonic-agent-plus/pkg/gadb"
	"testing"
)

var (
	client gadb.Client
)

func SetClient() {
	client, _ = gadb.NewClient()
}

func TestPocoElement_Start(t *testing.T) {
	SetClient()

	dev, err := tool.GetDevice(client, "")
	assert.NoError(t, err)

	dev.ForwardLocalAbstract(5001, "5001")

	pocoPageDriver := poco.NewPocoPage(poco.UNITY_3D)

	err = pocoPageDriver.Start()
	assert.NoError(t, err)
}

func TestPocoElement_Dump(t *testing.T) {
	SetClient()

	dev, err := tool.GetDevice(client, "")
	assert.NoError(t, err)

	dev.FrowardTcp(5001, "5001")

	pocoPageDriver := poco.NewPocoPage(poco.UNITY_3D)

	err = pocoPageDriver.Start()
	assert.NoError(t, err)

	res, err := pocoPageDriver.DumpRootToJson()
	assert.NoError(t, err)
	fmt.Println(res.OutputXML())
}
