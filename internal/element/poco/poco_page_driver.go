package poco

import (
	"errors"
	"github.com/antchfx/jsonquery"
	"github.com/antchfx/xmlquery"
	"github.com/tidwall/gjson"
	"sonic-agent-plus/define/api"
	"strings"
)

var _ api.IPageDriver = (*PocoPageDriver)(nil)

func NewPocoPage(pocoType PocoType) *PocoPageDriver {
	return &PocoPageDriver{
		pType: pocoType,
	}
}
func NewPocoPageByPort(port int, pocoType PocoType) *PocoPageDriver {
	return &PocoPageDriver{
		port:  port,
		pType: pocoType,
	}
}

type PocoPageDriver struct {
	port       int
	pType      PocoType
	pocoClient iPocoClient
}

func (p *PocoPageDriver) Start() error {
	if !iSPoco(p.pType) {
		return errors.New("the specified dump driver is not a poco")
	}
	if p.port == 0 {
		p.port = getPocoDefaultPortByName(p.pType)
	}
	if p.pType == COCOS_2DX_JS || p.pType == COCOS_CREATOR || p.pType == EGRET {
		p.pocoClient = newWebSocketClientImpl(p.port, p.pType)
	} else {
		p.pocoClient = newSocketClientImpl(p.port, p.pType)
	}
	// todo 加入重连机制
	err := p.pocoClient.Connect()
	return err
}

func (p *PocoPageDriver) Stop() error {
	p.pocoClient.Disconnect()
	return nil
}

func (p *PocoPageDriver) DumpRootToJson() (*jsonquery.Node, error) {
	result, err := p.pocoClient.Dump()
	if err != nil {
		return nil, err
	}
	dumpJsonStr := gjson.Parse(result).Get("result").String()
	doc, err := jsonquery.Parse(strings.NewReader(dumpJsonStr))
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func (p *PocoPageDriver) DumpRootToXml() (*xmlquery.Node, error) {
	panic("")
}

func (p *PocoPageDriver) FindJsonNodeByXpath(xpath string) (*jsonquery.Node, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PocoPageDriver) FindXmlNodeByXpath(xpath string) (*xmlquery.Node, error) {
	//TODO implement me
	panic("implement me")
}
