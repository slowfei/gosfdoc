package golang

import (
	"strings"
	// "fmt"
	"github.com/slowfei/gosfcore/utils/filemanager"
	"github.com/slowfei/gosfdoc"
	"path/filepath"
	"testing"
)

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
	//	TODO EachIndexFile已经改变，需要重新测试
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

func TestFindDefine(t *testing.T) {
	testFile := `
//
// temp1 
const (
	Test1 = "1"
	SNRoundBrackets = SFSubUtil.NewSubNest([]byte("("), []byte(")"))
	SNBetweens      = []*SFSubUtil.SubNest{
		SNBraces,
		sub.NewSubNest([]byte("/*"), []byte("*/")),
		"temo",
	}
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
		t.Log("内容：", strings.Replace(testFile[define.contIndex[0]:define.contIndex[1]], "\n", "<br>", -1))
		t.Log("是否多行：", define.multiterm)

		dtype := ""
		switch define.dtype {
		case goDTypeConst:
			dtype = "const"
		case goDTypeVar:
			dtype = "var"
		}
		t.Log("类型：", dtype)
		t.Log("----------")
	}

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
