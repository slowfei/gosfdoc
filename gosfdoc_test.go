package gosfdoc

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"testing"
)

func TestParseMarkdown(t *testing.T) {
	count := 3

	relPath := ""
	filesName := []string{"main.go"}
	docs := make([]Document, 0, count)
	pres := make([]Preview, 0, count)
	blocks := make([]CodeBlock, 0, count)

	for i := 1; i <= count; i++ {
		doc := Document{}
		doc.SortTag = i
		doc.Title = fmt.Sprintf("Document_title_%d", i)
		doc.Content = fmt.Sprintf("\ndocument markdown syntax content %d\n", i)
		docs = append(docs, doc)

		pre := Preview{}
		pre.Level = 0
		pre.ShowText = fmt.Sprintf("type TestStruct%d struct", i)
		pre.Anchor = fmt.Sprintf("temp.TestStruct%d", i)
		pre.DescText = "implements: [Test](javascript:;)"
		pres = append(pres, pre)

		pre2 := Preview{}
		pre2.Level = 1
		pre2.ShowText = fmt.Sprintf("func (* TestStruct) Hello%d(str string) string", i)
		pre2.Anchor = fmt.Sprintf("temp.TestStruct.Hello%d", i)
		pres = append(pres, pre2)

		pre3 := Preview{}
		pre3.Level = 1
		pre3.ShowText = fmt.Sprintf("func (* TestStruct) Hello2_%d(str string) string", i)
		pre3.Anchor = fmt.Sprintf("temp.TestStruct.Hello2_%d", i)
		pres = append(pres, pre3)

		block := CodeBlock{}
		block.Title = fmt.Sprintf("type TestStruct%d struct", i)
		block.Anchor = fmt.Sprintf("temp.TestStruct%d", i)
		block.Desc = `
函数介绍描述

@param 'v1' 字符串传递
@param 'v2' 字符串传递
@return
		`
		block.Code = `
type TestStruct struct{
	temp string
}
`
		block.CodeLang = "go"
		block.MenuTitle = "Func Details"
		block.SourceFileName = "main.go"
		block.FileLines = []int{10, 20}
		blocks = append(blocks, block)
	}

	body := ParseMarkdown(docs, pres, blocks, filesName, "1.0", relPath)
	t.Log(string(body))
}

/**
 *
 */
func TestSortSet(t *testing.T) {
	pres := make([]Preview, 5, 5)

	pres[0] = Preview{SortTag: "4"}
	pres[1] = Preview{SortTag: "2"}
	pres[2] = Preview{SortTag: "5"}
	pres[3] = Preview{SortTag: "3"}
	pres[4] = Preview{SortTag: "1"}

	sort.Sort(SortSet{previews: pres})

	if "1" != pres[0].SortTag {
		t.Fail()
	}
}

/**
 *
 */
func TestParseDocument(t *testing.T) {
	strTest := `
/***1-title1
 *  content_1
 */

/***title2
 *  content_2
 */



func test2(){
}

///3-title3
//  content_3
//End

func test2(){
}

`

	documents := ParseDocument(NewFileBuf([]byte(strTest), "", nil, nil))

	if 3 != len(documents) {
		t.Fatalf("3 != len(documents)")
		return
	}

	if documents[0].SortTag != 1 || documents[0].Title != "title1" || -1 == strings.Index(documents[0].Content, "content_1") {
		t.Fatalf("ErrorParse: %v %v \n%v", documents[0].SortTag, documents[0].Title, documents[0].Content)
	}

	if documents[1].SortTag != 0 || documents[1].Title != "title2" || -1 == strings.Index(documents[1].Content, "content_2") {
		t.Fatalf("ErrorParse: %v %v \n%v", documents[1].SortTag, documents[1].Title, documents[1].Content)
	}

	if documents[2].SortTag != 3 || documents[2].Title != "title3" || -1 == strings.Index(documents[2].Content, "content_3") {
		t.Fatalf("ErrorParse: %v %v \n%v", documents[2].SortTag, documents[2].Title, documents[2].Content)
	}

}

func TestREXDCPackageInfo(t *testing.T) {
	str := `
//#package-info use helper     

`

	result := REXDCPackageInfo.FindSubmatchIndex([]byte(str))

	if 4 != len(result) {
		t.Fatal()
		return
	}
	t.Log(result)

	restr := str[result[2]:result[3]]
	restr = strings.Trim(restr, " ")
	t.Log(restr)
}

/**
 *
 */
func TestREXDocument(t *testing.T) {
	str := `
/***1-title1
 *  content_1
 */

 /***2-
  *  content_2
  */

////3-title3
// content_3
//End

///4-
// content_4
//End

`
	results := REXDocument.FindAllString(str, -1)

	if 4 != len(results) {
		t.Fatal("find result num fatal.")
	}

	str2 := `///1234-title`
	subResult := REXDocIndexTitle.FindStringSubmatch(str2)
	if subResult[2] != "1234-" {
		t.Fatal("substring index fatal.")
	}

}

/**
 *
 */
func TestREXPrivateBlock(t *testing.T) {
	str := `
func test1(){
}
//#private
func test2(){
}
//#private-end
func test3(){
}
//#private
func test4(){
}
//#private-end
func test5(){
}
    `
	strResult := `
func test1(){
}
func test3(){
}
func test5(){
}
    `
	bf := NewFileBuf([]byte(str), "", nil, REXPrivateBlock)

	if bf.String() != strResult {
		t.Fatalf("private block replace failed.")
	}
}

/**
 *
 */
func TestREXParseAboutAndIntro(t *testing.T) {
	strTest := `

func func_name() {
    
}

//About
//  ## About
//  ------
//  
//  gosfdoc document generator
// 
//  More references: [https://github.com/slowfei/gosfdoc][0]<br/>
//  The MIT license (MIT) - [http://opensource.org/licenses/MIT][1]
temp
//  Copyright (c) 2014 slowfei<br/>
//  Email: slowfei#foxmail.com
//
//  [0]:https://github.com/slowfei/gosfdoc
//  [1]:http://opensource.org/licenses/MIT
//  [2]:http://opensource.org/licenses/MIT3
//End

/**About
 *  ## About
 *  ------
 *  
 *  gosfdoc document generator
 * 
 *  More references: [https://github.com/slowfei/gosfdoc][0]<br/>
 *  The MIT license (MIT) - [http://opensource.org/licenses/MIT][1]
temp
 *  Copyright (c) 2014 slowfei<br/>
 *  Email: slowfei#foxmail.com
 *
 *  [0]:https://github.com/slowfei/gosfdoc
 *  [1]:http://opensource.org/licenses/MIT
 *  [2]:http://opensource.org/licenses/MIT2
 */


/**Intro
 *  ## Intro
 *  ------
 *  
 *  Intro content
 */

 func func_name() {
     
 }
`

	newBuf := ParseAbout(NewFileBuf([]byte(strTest), "", nil, nil))
	lineNumber := len(bytes.Split(newBuf.Content, []byte("\n")))

	if 15 != lineNumber {
		t.Log(string(newBuf.Content))
		t.Fatal("About: Target Line 14, parse Line error:", lineNumber)
	}

	newBuf2 := ParseIntro(NewFileBuf([]byte(strTest), "", nil, nil))
	lineNumber = len(bytes.Split(newBuf2.Content, []byte("\n")))
	if 5 != lineNumber {
		t.Log(string(newBuf2.Content))
		t.Fatal("Intro: Target Line 5, parse Line error:", lineNumber)
	}

}
