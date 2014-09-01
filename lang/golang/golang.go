package golang

import (
	"bytes"
	"github.com/slowfei/gosfdoc"
	"os"
)

const (
	GO_NAME = "go"
)

func init() {
	gosfdoc.AddParser(NewParser())
}

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

func (g *GolangParser) EachFile(index int, fileCont *bytes.Buffer, info os.FileInfo) {

}

func (g *GolangParser) ParseDoc(fileCont *bytes.Buffer) []gosfdoc.Document {
	return nil
}

func (g *GolangParser) ParsePreview(fileCont *bytes.Buffer) []gosfdoc.Preview {
	return nil
}

func (g *GolangParser) ParseCodeblock(fileCont *bytes.Buffer) []gosfdoc.CodeBlock {
	return nil
}
