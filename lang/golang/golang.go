//  The MIT License (MIT) - http://opensource.org/licenses/MIT
//
//  Copyright (c) 2014 slowfei
//
//  Create on 2014-11-05
//  Update on 2014-11-26
//  Email  slowfei#foxmail.com
//  Home   http://www.slowfei.com

//	golang implement parser
package golang

import (
	"github.com/slowfei/gosfdoc"
	"os"
	// "regexp"
	"strings"
)

const (
	GO_NAME   = "go"
	GO_SUFFIX = ".go"
)

var (
//	\/\*\*[\s]*\n(\s|.)*?\*/\nfunc\s[a-zA-z_].*\(.*\)\s*{\s*\n(\s|.)*?}
)

func init() {
	gosfdoc.AddParser(NewParser())
}

/**
 *	golang parser
 */
type GolangParser struct {
}

/**
 *	new golang parser
 */
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
func (g *GolangParser) CheckFile(filePath string, info os.FileInfo) bool {
	result := false

	if 0 != len(filePath) && nil != info && !info.IsDir() {
		result = strings.HasSuffix(filePath, GO_SUFFIX)
	}
	return result
}

/**
 *	see DocParser interface
 */
func (g *GolangParser) EachIndexFile(filebuf *gosfdoc.FileBuf) {

}

/**
 *	see DocParser interface
 */
func (g *GolangParser) ParsePreview(filebuf *gosfdoc.FileBuf) []gosfdoc.Preview {
	return nil
}

/**
 *	see DocParser interface
 */
func (g *GolangParser) ParseCodeblock(filebuf *gosfdoc.FileBuf) []gosfdoc.CodeBlock {
	return nil
}

/**
 *	see DocParser interface
 */
func (n *GolangParser) ParsePackageInfo(filebuf *gosfdoc.FileBuf) string {
	return ""
}
