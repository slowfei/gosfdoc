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

/**
 *	see DocParser interface
 */
func (g *GolangParser) CheckFile(path string, info os.FileInfo) bool {
	return false
}

/**
 *	see DocParser interface
 */
func (g *GolangParser) EachIndexFile(filebuf *gosfdoc.FileBuf) {

}

/**
 *	see DocParser interface
 */
func (g *GolangParser) ParsePreview(fileCont *bytes.Buffer) []gosfdoc.Preview {
	return nil
}

/**
 *	see DocParser interface
 */
func (g *GolangParser) ParseCodeblock(fileCont *bytes.Buffer) []gosfdoc.CodeBlock {
	return nil
}
