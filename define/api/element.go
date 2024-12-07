package api

import (
	"github.com/antchfx/jsonquery"
	"github.com/antchfx/xmlquery"
)

type IElement interface {
	Start() error
	Stop() error
	DumpRootToJson() (jsonquery.Node, error)
	DumpRootToXml() (xmlquery.Node, error)
	FindJsonNodeByXpath(xpath string) (jsonquery.Node, error)
	FindXmlNodeByXpath(xpath string) (xmlquery.Node, error)
}
