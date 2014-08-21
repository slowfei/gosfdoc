package golang

import (
	"bytes"
)

const (
	GO_NAME = "go"
)

type GolangParser struct {
}

func NewParser() *GolangParser {
	return new(GolangParser)
}

//#pragma mark github.com/slowfei/gosfdoc.DocParser interface ---------------------------------------------------------------------

func (g *GolangParser) Name() string {
	return GO_NAME
}
func (g *GolangParser) CheckFilepath() bool {
	return false
}
func (g *GolangParser) EachFile(index int, fileCont *bytes.Buffer) {

}
