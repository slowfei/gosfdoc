package gosfdoc

import (
	"bytes"
	"testing"
)

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
	bf := NewFileBuf([]byte(str), REXPrivateBlock)

	if bf.String() != strResult {
		t.Fatalf("private block replace failed.")
	}
}

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
//
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
 *
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

	newBuf := ParseAbout(NewFileBuf([]byte(strTest), nil))
	lineNumber := len(bytes.Split(newBuf, []byte("\n")))

	if 15 != lineNumber {
		// t.Log(string(newBuf))
		t.Fatal("About: Target Line 14, parse Line error:", lineNumber)
	}

	newBuf = ParseIntro(NewFileBuf([]byte(strTest), nil))
	lineNumber = len(bytes.Split(newBuf, []byte("\n")))
	if 5 != lineNumber {
		t.Log(string(newBuf))
		t.Fatal("Intro: Target Line 5, parse Line error:", lineNumber)
	}

}
