//  The MIT License (MIT) - http://opensource.org/licenses/MIT
//
//  Copyright (c) 2014 slowfei
//
//  Create on 2014-11-05
//  Update on 2014-12-09
//  Email  slowfei#foxmail.com
//  Home   http://www.slowfei.com

//	golang implement parser
package golang

import (
	"github.com/slowfei/gosfcore/utils/filemanager"
	"github.com/slowfei/gosfcore/utils/sub"
	"github.com/slowfei/gosfdoc"
	"github.com/slowfei/gosfdoc/index"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	GO_NAME        = "go"
	GO_SUFFIX      = ".go"
	GO_TEST_SUFFIX = "_test.go"
)

var (
	//	\/\*\*[\s]*\n(\s|.)*?\*/\nfunc\s[a-zA-z_].*\(.*\)\s*\{\s*\n(?:\{\s*|.*\}|\s|.)*?\}
	//	type ([A-Z]\w*) \w+(\s?\{(\s*|.*)*?\})*
	REXType    = regexp.MustCompile("type ([A-Z]\\w*) \\w+[ ]*(\\{)?")
	SNBraces   = SFSubUtil.NewSubNest([]byte("{"), []byte("}"))
	SNBetweens = []*SFSubUtil.SubNest{
		SNBraces,
		SFSubUtil.NewSubNest([]byte("`"), []byte("`")),
		SFSubUtil.NewSubNest([]byte(`"`), []byte(`"`)),
		SFSubUtil.NewSubNest([]byte(`'`), []byte(`'`)),
	}
)

func init() {
	gosfdoc.AddParser(NewParser())
}

/**
 *	golang parser
 */
type GolangParser struct {
	config  gosfdoc.MainConfig
	indexDB index.IndexDB
}

/**
 *	new golang parser
 */
func NewParser() *GolangParser {
	gp := new(GolangParser)
	gp.indexDB = index.CreateIndexDB(GO_NAME, index.DBTypeFile)
	return gp
}

//#pragma mark github.com/slowfei/gosfdoc.DocParser interface ---------------------------------------------------------------------

func (g *GolangParser) Name() string {
	return GO_NAME
}

/**
 *	see DocParser interface
 */
func (g *GolangParser) ParseStart(config gosfdoc.MainConfig) {
	g.config = config
	g.indexDB.Open()
}

/**
 *	see DocParser interface
 */
func (g *GolangParser) ParseEnd() {
	g.indexDB.Close()
}

/**
 *	see DocParser interface
 */
func (g *GolangParser) CheckFile(filePath string, info os.FileInfo) bool {
	result := false

	if 0 != len(filePath) && nil != info && !info.IsDir() {
		result = strings.HasSuffix(filePath, GO_SUFFIX)

		if result {
			result = !strings.HasSuffix(filePath, GO_TEST_SUFFIX)
		}
	}
	return result
}

/**
 *	see DocParser interface
 */
func (g *GolangParser) EachIndexFile(filebuf *gosfdoc.FileBuf) {
	// find type (XXXX)
	var outBetweens [][]int
	for i := 0; i < len(SNBetweens); i++ {
		tempIndexs := filebuf.SubNestAllIndex(SNBetweens[i], nil)
		if 0 != len(tempIndexs) {
			outBetweens = append(outBetweens, tempIndexs...)
		}
	}
	var typeInfos []index.TypeInfo
	tempPackage := ""

	//	包查询 TODO， 解决如何查询包
	gopaths := SFFileManager.GetGOPATHDirs()
	for i := 0; i < len(gopaths); i++ {
		gopath := path.Join(gopaths[i], "src")
		filebufPath := path.Dir(filebuf.Path())
		if strings.HasPrefix(filebufPath, gopath) {
			tempPackage = filebufPath[len(gopath)+1 : len(filebufPath)]
			//	TODO 需要考虑查询的情况，直接使用这样的包名是否可以查找。
		}
	}

	//	类型查询
	typeIndexs := filebuf.FindAllSubmatchIndex(REXType)
	for i := 0; i < len(typeIndexs); i++ {
		indexs := typeIndexs[i]
		startIndex := indexs[0]
		endIndex := indexs[1]
		tempType := index.TypeInfo{}

		// type GolangParser struct { [1 27 6 18 26 27]
		// type OperateResult int [88 110 93 106 -1 -1]
		if 6 == len(indexs) && !isRuleOutIndex(startIndex, outBetweens) {

			leftBraces := indexs[4]
			rightBraces := indexs[5]
			if -1 != leftBraces && -1 != rightBraces {
				bracesIndexs := filebuf.SubNestIndex(leftBraces, SNBraces, outBetweens)
				if 2 == len(bracesIndexs) && -1 != bracesIndexs[0] && -1 != bracesIndexs[1] {
					endIndex = bracesIndexs[1]
				}
			}

			lines := filebuf.LineNumberByIndex(startIndex, endIndex)
			if -1 != lines[0] && -1 != lines[1] {
				tempType.LineStart = lines[0]
				tempType.LineEnd = lines[1]
			}

			startName := indexs[2]
			endName := indexs[3]
			if -1 != startName && -1 != endName {
				tempType.Name = string(filebuf.SubBytes(startName, endName))
			} else {
				tempType.Name = ""
			}
		}
	}

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

/**
 *	判断是否是排除坐标
 *
 *	@return 在坐标范围内返回 true
 */
func isRuleOutIndex(index int, outBetweens [][]int) bool {
	result := false

	for i := 0; i < len(outBetweens); i++ {
		indexs := outBetweens[i]
		if 2 == len(indexs) {
			s := indexs[0]
			e := indexs[1]
			if index > s && index < e {
				result = true
				break
			}
		}
	}
	return result
}
