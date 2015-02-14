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
    `
	s := REXType.FindAllSubmatchIndex([]byte(testStr), -1)

	for i := 0; i < len(s); i++ {
		indexs := s[i]
		t.Log(indexs)
		t.Log(testStr[indexs[0]:indexs[1]])
		t.Log(testStr[indexs[2]:indexs[3]])
	}

	if 2 != len(s) || testStr[s[0][2]:s[0][3]] != "GolangParser" {
		t.Fail()
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
		t.Fail()
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
		t.Fail()
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
		t.Fail()
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
		t.Fail()
	}
}

func TestFindDefine(t *testing.T) {
	testFile := `
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
		t.Fail()
		return
	}

	for _, define := range result {

		note := ""
		if -1 != define.noteIndex[0] {
			note = strings.Replace(testFile[define.noteIndex[0]:define.noteIndex[1]], "\n", "<br>", -1)
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
		t.Fail()
	}
}
