package golang

import (
	"strings"
	// "fmt"
	"github.com/slowfei/gosfcore/utils/filemanager"
	"github.com/slowfei/gosfdoc"
	"path/filepath"
	"testing"
)

func TestRegexpFunc(t *testing.T) {
	testStr := `
/**
 * new parser
 */
func NewParser() *GolangParser {
}

/**
 *	see DocParser interface
 */
func (g *GolangParser) ParseStart(config gosfdoc.MainConfig) {
}

/**
 *	see DocParser interface
 */
func (g *GolangParser) CheckFile(filePath string, info os.FileInfo) bool {
}

func TestReturn()(bool,int){

}

func testFilter(){

}

`
	s := REXFunc.FindAllSubmatchIndex([]byte(testStr), -1)

	if 4 != len(s) && 16 != len(s[0]) {
		t.Fatal()
		return
	}

	// [0-1:prototype] [2-3:comment] [4-5:func type]
	// [6-7:func name] [8-9:func params] [10-11:single return value]
	// [12-13:multi return value] [14-15:"{"]
	if "NewParser" != testStr[s[0][6]:s[0][7]] {
		t.Fatal()
	}

	if "g *GolangParser" != testStr[s[1][4]:s[1][5]] {
		t.Fatal(testStr[s[1][4]:s[1][5]])
	}

	if "filePath string, info os.FileInfo" != testStr[s[2][8]:s[2][9]] {
		t.Fatal()
	}

	if "bool" != testStr[s[2][10]:s[2][11]] {
		t.Fatal(testStr[s[2][10]:s[2][11]])
	}

	if "bool,int" != testStr[s[3][12]:s[3][13]] {
		t.Fatal(testStr[s[3][12]:s[3][13]])
	}

}

func TestRegexpType(t *testing.T) {
	testStr := `
type GolangParser struct {
    config  gosfdoc.MainConfig
    indexDB index.IndexDB
}

type OperateResult int

/**
 *	temp comt
 */
type Temp struct{
	temp string
}

    `
	s := REXType.FindAllSubmatchIndex([]byte(testStr), -1)

	if 3 != len(s) && 10 != len(s[0]) {
		t.Fatal()
		return
	}

	// 0-index: type GolangParser
	// 1-index: type OperateResult
	// 2-index: type Temp

	index := 2
	newType := s[index]

	// index 4-5 type (Temp) struct
	if "Temp" != testStr[newType[4]:newType[5]] {
		t.Fatal()
	}

	//	检查大括号
	// 下标1 OperateResult不存在大括号，所以等于-1
	if -1 != s[1][8] || -1 != s[1][9] {
		t.Fatal()
	}

	//	下标0 有大括号的
	if -1 == s[0][8] || -1 == s[0][9] {
		t.Fatal()
	}

	// 0-index 4-5 type (GolangParser) struct
	if "GolangParser" != testStr[s[0][4]:s[0][5]] {
		t.Fatal()
	}

	// 由于"type "是固定的，所以从获取的类型名 -5 个位数
	if "type " != testStr[s[0][4]-5:s[0][4]] {
		t.Fatal()
	}

}

func TestRegexpPackage(t *testing.T) {
	testStr := `
package main

`
	s := REXPackage.FindAllSubmatchIndex([]byte(testStr), -1)

	result := testStr[s[0][2]:s[0][3]]
	t.Log(result)

	if result != "main" {
		t.Fatal()
	}

}

func TestRegexpPackageInfo(t *testing.T) {
	testStr := `

/**
 *	test pacakage1
 */
 package main1

 /**
  *	test pacakage2
  */
 package main2

 // 
 // test package3
 package main3

`
	subBytes := REXPackageInfo.FindSubmatch([]byte(testStr))
	subLen := len(REXPackageInfo.FindAllString(testStr, -1))

	if 2 != len(subBytes) || 3 != subLen {
		t.Fatal()
		return
	}

	t.Log(string(subBytes[1]))

}

func TestRegexpDefine(t *testing.T) {
	testStr := `
//
// temp1 
const (   
	Test1 = "1"
)

/**
 * temp2
 */
const (
	Test2 = "2"
)

const Temp3 = "3"

// VTest1 cont
var VTest1  = "1"
var (  
	VTest2 = 3
	VTest3 = 4
)

`
	result := REXDefine.FindAllStringSubmatch(testStr, -1)

	if 5 != len(result) {
		t.Fatal()
		return
	}

	for i := 0; i < len(result); i++ {
		t.Log(strings.Replace(strings.Join(result[i], ", "), "\n", "<br>", -1))
	}

}

func createFileBuf(fileCont string) *gosfdoc.FileBuf {
	testpackagepath := "github.com/slowfei/gosfdoc"
	testfilename := "testfile.go"
	testfilepath := filepath.Join(SFFileManager.GetGOPATHDirs()[0], "src", testpackagepath, testfilename)
	filebuf := gosfdoc.NewFileBuf([]byte(fileCont), testfilepath, nil, nil)
	return filebuf
}

func TestEachIndexFile(t *testing.T) {
	testFile := `
package main

type Test1 int

type TestStruct struct{
	v1 string
}
`
	parser := NewParser()
	parser.EachIndexFile(createFileBuf(testFile))

	result, bool := parser.indexDB.Type("main", "github.com/slowfei/gosfdoc", "Test1")
	t.Log(result)
	if !bool {
		t.Fatal()
	}
}

func TestFindType(t *testing.T) {
	testFile := `
package main

/**
 *	test1 comment
 */
type Test1 int

type TestStruct struct{
	v1 string
}

// test interface
type TestInterface interface{
	Temp() string
	temp2() interface{}
}

/*
	/*
		tempcomt
	*/
	type TestComt int
 */

{
	type TestOut struct{
		v2 string
	}
}

`
	buf := createFileBuf(testFile)
	outBetweens := getOutBetweens(buf)
	result := findType(buf, outBetweens)

	if 3 != len(result) {
		t.Fatal(len(result))
		return
	}

	for i := 0; i < len(result); i++ {
		gt := result[i]
		if 2 == len(gt.commentIndex) {
			t.Log("注释：", strings.Replace(testFile[gt.commentIndex[0]:gt.commentIndex[1]], "\n", "<br>", -1))
		}
		t.Log("类型名：", testFile[gt.typeNameIndex[0]:gt.typeNameIndex[1]])
		t.Log("类型：", testFile[gt.typeIndex[0]:gt.typeIndex[1]])
		t.Log("原型：", strings.Replace(testFile[gt.bodyIndex[0]:gt.bodyIndex[1]], "\n", "<br>", -1))
		t.Log("----------")
	}
}

func TestFindDefine(t *testing.T) {
	testFile := `
//
// temp1 
const (
	Test1 = "1"
	SNRoundBrackets = SFSubUtil.NewSubNest(
		[]byte("("),
	)
	SNBetweens      = []*SFSubUtil.SubNest{
		SNBraces,
		sub.NewSubNest([]byte("/*"), []byte("*/")),
	}
	SNTemp = "
		temp string
	"
)

/**
 * temp2
 */
const (
	Test2 = "2"
)

const Temp3 = "3"

// VTest1 cont
var VTest1  = "1"
var (
	VTest2 = 3
	VTest3 = 4
	temp = 4 // 由于是小写开头命名，所以整个var()都会被过滤
)

// 以下都是被过滤的
{
	const Temp4 = "4"
}

/*
	const Temp5 = "5"
*/

' const Temp6 = "6" '

// const Temp7 = "7"

`

	buf := createFileBuf(testFile)
	outBetweens := getOutBetweens(buf)
	result := findDefine(buf, outBetweens)

	if 4 != len(result) {
		t.Log(len(result))
		t.Fatal()
		return
	}

	for _, define := range result {

		note := ""
		if -1 != define.commentIndex[0] {
			note = strings.Replace(testFile[define.commentIndex[0]:define.commentIndex[1]], "\n", "<br>", -1)
		}
		t.Log("注释：", note)
		t.Log("内容：", strings.Replace(testFile[define.bodyIndex[0]:define.bodyIndex[1]], "\n", "<br>", -1))
		t.Log("是否多行：", define.multiterm)

		dtype := ""
		switch define.dtype {
		case goDTypeConst:
			dtype = "const"
		case goDTypeVar:
			dtype = "var"
		}
		t.Log("类型：", dtype)

		names := ""
		for i := 0; i < len(define.nameIndexs); i++ {
			names += testFile[define.nameIndexs[i]:define.nameIndexs[i+1]] + ","
			i++
		}
		t.Log("参数名：", names)

		t.Log("----------")
	}

}

func TestFindFunc(t *testing.T) {
	testFile := `

/**
 * new parser
 */
func NewParser() *GolangParser {
}

/**
 *	see DocParser interface
 */
func (g *GolangParser) ParseStart(config gosfdoc.MainConfig) {
}

/**
 *	see DocParser interface
 */
func (g *GolangParser) CheckFile(filePath string, info os.FileInfo) bool {
}

//注意单引号，别过滤的
'func TestReturn()(bool,int){

}'

func testFilter(){

}

`
	buf := createFileBuf(testFile)
	outBetweens := getOutBetweens(buf)
	result := findFunc(buf, outBetweens)

	if 3 != len(result) {
		t.Fatal(len(result))
		return
	}

	for i := 0; i < len(result); i++ {
		gf := result[i]
		if 2 == len(gf.commentIndex) {
			t.Log("注释：", strings.Replace(testFile[gf.commentIndex[0]:gf.commentIndex[1]], "\n", "<br>", -1))
		}

		ft := ""
		if 0 != len(gf.funcTypeIndex) {
			ft = testFile[gf.funcTypeIndex[0]:gf.funcTypeIndex[1]]
		}
		t.Log("函数类型：", ft)

		t.Log("函数名：", testFile[gf.funcNameIndex[0]:gf.funcNameIndex[1]])

		params := ""
		if 0 != len(gf.paramIndex) {
			params = testFile[gf.paramIndex[0]:gf.paramIndex[1]]
		}
		t.Log("参数：", params)

		returnVal := ""
		if 0 != len(gf.returnIndex) {
			returnVal = testFile[gf.returnIndex[0]:gf.returnIndex[1]]
		}
		t.Log("返回值：", returnVal)

		t.Log("原型：", strings.Replace(testFile[gf.bodyIndex[0]:gf.bodyIndex[1]], "\n", "<br>", -1))

		t.Log("----------")
	}
}

func TestParsePreview(t *testing.T) {
	testFile := `
package main

//
// temp1 
const (
	Test1 = "1"
	SNRoundBrackets = SFSubUtil.NewSubNest(
		[]byte("("),
	)
	SNTemp = "
		temp string
	"
)

type TestStruct struct{
	v1 string
}

func (t *TestStruct) ParseStart() {
}

func NewTestStruct() []TestStruct {

}

func NewParser() {

}
`

	parser := NewParser()

	buf := createFileBuf(testFile)
	parser.EachIndexFile(buf)
	previews := parser.ParsePreview(buf)

	if 5 != len(previews) {
		t.Fatal()
		return
	}

	t.Log("Preview:")
	for i := 0; i < len(previews); i++ {
		pre := previews[i]
		t.Log("pre.Anchor = ", pre.Anchor)
		t.Log("pre.DescText = ", pre.DescText)
		t.Log("pre.Level = ", pre.Level)
		t.Log("pre.ShowText = ", pre.ShowText)
		t.Log("pre.SortTag = ", pre.SortTag)
		t.Log("----------")

	}
	t.Log("--------------------------------")

}

func TestParseCodeblock(t *testing.T) {
	testFile := `
package main

/**
 *	const comment
 */
const (
	Test1 = "1"
	SNRoundBrackets = SFSubUtil.NewSubNest(
		[]byte("("),
	)
	SNTemp = "
		temp string
	"
)

/**
 *	struct comment
 */
type TestStruct struct{
	v1 string
}

func (t *TestStruct) ParseStart() {
}

func NewTestStruct() []TestStruct {

}

func NewParser() {

}
`
	parser := NewParser()

	buf := createFileBuf(testFile)
	parser.EachIndexFile(buf)
	blockCodes := parser.ParseCodeblock(buf)

	t.Log("Codeblock:")
	for i := 0; i < len(blockCodes); i++ {
		block := blockCodes[i]

		t.Log("block.Anchor = ", block.Anchor)
		t.Log("block.Code = ", block.Code)
		t.Log("block.CodeLang = ", block.CodeLang)
		t.Log("block.Desc = ", block.Desc)
		t.Log("block.FileLines = ", block.FileLines)
		t.Log("block.MenuTitle = ", block.MenuTitle)
		t.Log("block.SortTag = ", block.SortTag)
		t.Log("block.SourceFileName = ", block.SourceFileName)

		t.Log("----------")

	}
	t.Log("--------------------------------")

}

func TestHandleComment(t *testing.T) {
	// 	src := []byte(`/**
	//  *  temp705
	//  * 	templet
	// */`)

	src2 := []byte(`// ntelan
// temp
//`)
	result := handleComment(src2)

	t.Log(string(result))

}

func TestParsePackageInfo(t *testing.T) {
	testFile := `


/**
 *	test pacakage
 *	test line2
 *	temp late3
 */
 package main1

// test package
// temp line2
// temp3
package main2

`

	parser := NewParser()
	result := parser.ParsePackageInfo(createFileBuf(testFile))
	t.Log(result)
	if 0 == len(result) {
		t.Fatal()
	}
}
