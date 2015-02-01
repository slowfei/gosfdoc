package golang

import (
	// "strings"
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
	s := REGType.FindAllSubmatchIndex([]byte(testStr), -1)

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
	s := REGPackage.FindAllSubmatchIndex([]byte(testStr), -1)

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
	subBytes := REGPackageInfo.FindSubmatch([]byte(testStr))
	subLen := len(REGPackageInfo.FindAllString(testStr, -1))

	if 2 != len(subBytes) || 3 != subLen {
		t.Fail()
		return
	}

	t.Log(string(subBytes[1]))

}

func TestRegexpConst(t *testing.T) {
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

`
	result := REGConst.FindAllString(testStr, -1)

	if 3 != len(result) {
		t.Fail()
	}

	for i := 0; i < len(result); i++ {
		t.Log(result[i])
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
