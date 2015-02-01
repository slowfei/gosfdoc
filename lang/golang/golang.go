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
	"bytes"
	"fmt"
	"github.com/slowfei/gosfcore/utils/filemanager"
	"github.com/slowfei/gosfcore/utils/sub"
	"github.com/slowfei/gosfdoc"
	"github.com/slowfei/gosfdoc/index"
	"os"
	"path"
	"regexp"
	"strings"
)

const (
	GO_NAME        = "go"
	GO_SUFFIX      = ".go"
	GO_TEST_SUFFIX = "_test.go"
)

var (
	// e.g.: type Temp struct {
	REGType = regexp.MustCompile("type ([A-Z]\\w*) \\w+[ ]*(\\{)?")
	// e.g.: package main
	REGPackage = regexp.MustCompile("package (\\w+)\\s*")
	// e.g.: /** ... */[\n]package main; //...[\n]package main
	REGPackageInfo = regexp.MustCompile("(/\\*\\*[\\S\\s]+?\\*/\n|(?:(?:[ ]*//.*?\n)+))[ ]*package \\w+")
	// e.g.: /** ... */[\n]const TConst = 1; //...[\n]const (
	REGConst = regexp.MustCompile("(/\\*\\*[\\S\\s]+?\\*/\n|(?:(?:[ ]*//.*?\n)+))?[ ]*(const[\\s]+(?:[A-Z].*|\\()+?)")

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

type goConst struct {
}

type goVar struct {
}

type goFunc struct {
}

type goType struct {
	funcs []goTypeFunc
}

type goTypeFunc struct {
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
	if nil == filebuf.UserData {
		outBetweens = getOutBetweens(filebuf)
		filebuf.UserData = outBetweens
	}

	tempPackagePath := ""
	tempPackageName := ""

	// find package name
	packageIndexs := filebuf.FindAllSubmatchIndex(REGPackage)
	for i := 0; i < len(packageIndexs); i++ {
		indexs := packageIndexs[i]
		if 4 == len(indexs) && !isRuleOutIndex(indexs[0], outBetweens) {
			tempPackageName = string(filebuf.SubBytes(indexs[2], indexs[3]))
			break
		}
	}

	//	查询不到package 证明是无效的文件
	if 0 == len(tempPackageName) {
		fmt.Println("InvalidFile: find less than the package name. file path:", filebuf.Path())
		return
	}

	// find package path
	gopaths := SFFileManager.GetGOPATHDirs()
	for i := 0; i < len(gopaths); i++ {
		gopath := path.Join(gopaths[i], "src")
		filebufPath := path.Dir(filebuf.Path())
		if strings.HasPrefix(filebufPath, gopath) {
			tempPackagePath = filebufPath[len(gopath)+1 : len(filebufPath)]
		}
	}

	//	无效文件提示
	if 0 == len(tempPackagePath) {
		fmt.Println("InvalidFile: is not a valid golang working environment file. file path:", filebuf.Path())
		return
	}

	//	类型查询
	typeIndexs := filebuf.FindAllSubmatchIndex(REGType)
	for i := 0; i < len(typeIndexs); i++ {
		indexs := typeIndexs[i]
		startIndex := indexs[0]
		endIndex := indexs[1]
		tempType := index.TypeInfo{}
		tempType.PackageName = tempPackageName
		tempType.PackagePath = tempPackagePath

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

			tempType.TypeName = string(filebuf.SubBytes(indexs[2], indexs[3]))
			if 0 != len(tempType.TypeName) {
				err := g.indexDB.SetType(tempType)
				if nil != err {
					fmt.Println("IndexError:", err.Error())
				}
			}

		}
	} // End for i := 0; i < len(typeIndexs); i++ {

}

/**
 *	see DocParser interface
 */
func (g *GolangParser) ParsePreview(filebuf *gosfdoc.FileBuf) []gosfdoc.Preview {
	//	TODO
	// 1. const const{ }
	// 2. var var { }
	// 3. func
	// 4. type
	// 5. 	return func type
	// 6. 	type func

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
	result := bytes.NewBuffer(nil)

	subBytes := filebuf.FindSubmatch(REGPackageInfo)
	if 2 != len(subBytes) {
		return ""
	}

	infoLines := bytes.Split(subBytes[1], []byte("\n"))
	reCount := 0
	var prefixTag []byte = nil
	prefixLen := 0

	//	判断是否存在 /* */ 如果存在则去除首行和尾行的扫描
	if 0 != len(infoLines) &&
		0 <= bytes.Index(infoLines[0], []byte("/*")) {
		reCount = 1
	}

	//	 len(infoLines)-reCount (-1) 由于正则截取的规则中会包含一个\n符号，所以需要去除
	for i := reCount; i < len(infoLines)-reCount-1; i++ {
		infoBytes := infoLines[i]

		if i == reCount {
			prefixTag = gosfdoc.FindPrefixFilterTag(infoBytes)
			prefixLen = len(prefixTag)
		}

		if nil != prefixTag {

			if 0 == bytes.Index(infoBytes, prefixTag) {
				result.Write(infoBytes[prefixLen:])
			} else {
				trimed := bytes.TrimSpace(infoBytes)
				// 有可能是空行，所需需要判断这行是否只有（ "*" || "//" ），如果不是则添加追加这一行内容
				if !bytes.Equal(trimed, []byte("*")) && !bytes.Equal(trimed, []byte("//")) {
					result.Write(infoBytes)
				} else {
					result.WriteByte('\n')
				}
			}

		} else {
			result.Write(infoBytes)
		}

		result.WriteByte('\n')
	}

	return result.String()
}

/**
 *	find constant
 *
 *	e.g:
 *	const xxx
 *	const ( ... )
 */
func findConst(filebuf *gosfdoc.FileBuf) {

}

/**
 *	find variable
 *
 *	e.g:
 *	var xxx
 *	var ( ... )
 */
func findVar() {

}

/**
 *	find function
 *
 *	e.g:
 *	func funcName()
 */
func findFunc() {

}

/**
 *	find type and function
 *
 *	e.g:
 *	type xxx
 *		func NewType() xxx
 *		func (type) funcName() xxx
 */
func findTypeAndFunc() {

}

/**
 *	获取文件排除范围的坐标范围
 *
 *	@param `filebuf`
 *	@return
 */
func getOutBetweens(filebuf *gosfdoc.FileBuf) [][]int {

	outBetweens := make([][]int, 0, 0)

	for i := 0; i < len(SNBetweens); i++ {
		tempIndexs := filebuf.SubNestAllIndex(SNBetweens[i], nil)
		if 0 != len(tempIndexs) {
			outBetweens = append(outBetweens, tempIndexs...)
		}
	}

	return outBetweens
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
