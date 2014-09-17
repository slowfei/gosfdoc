//  The MIT License (MIT) - http://opensource.org/licenses/MIT
//
//  Copyright (c) 2014 slowfei
//
//  Create on 2014-09-10
//  Update on 2014-09-18
//  Email  slowfei#foxmail.com
//  Home   http://www.slowfei.com

//
package gosfdoc

import (
	"bytes"
	"github.com/slowfei/gosfcore/utils/strings"
	"regexp"
	// "fmt"
	"os"
	"strings"
)

const (
	DOC_FILE_SUFFIX = ".doc"     // document file suffix
	NIL_DOC_NAME    = "document" // nilDocParser struct use
)

var (
	_defaultAbout = []byte(`
## About
------

gosfdoc document generator

More references: [https://github.com/slowfei/gosfdoc][0]<br/>
The MIT license (MIT) - [http://opensource.org/licenses/MIT][1]

Copyright (c) 2014 slowfei<br/>
Email: slowfei#foxmail.com

[0]:https://github.com/slowfei/gosfdoc
[1]:http://opensource.org/licenses/MIT
`)

	_defaultIntro = []byte(`
##Document Introduction

Sorry! Document author did not write any information.

----
This is a good tool, Can help you make beautiful documents.

More references: [https://github.com/slowfei/gosfdoc][0]<br/>

[0]:https://github.com/slowfei/gosfdoc
`)

	//  主要用于去除注释的前缀
	_prefixFilterTags = [][]byte{
		[]byte(" *\t"),
		[]byte(" *  "),
		[]byte(" * "),
		[]byte("//\t"),
		[]byte("//  "),
		[]byte("// "),
		[]byte("//"),
	}

	_tagStar   = []byte("*")  // comments (*)
	_tagDSlash = []byte("//") // double slash
)

/**
 *	nil document parser
 *	specifically for .doc file serve
 */
type nilDocParser struct {
}

/**
 *	see DocParser interface
 */
func (n *nilDocParser) Name() string {
	return NIL_DOC_NAME
}

/**
 *	see DocParser interface
 */
func (n *nilDocParser) CheckFile(path string, info os.FileInfo) bool {
	return strings.HasSuffix(path, DOC_FILE_SUFFIX)
}

/**
 *	see DocParser interface
 */
func (n *nilDocParser) EachIndexFile(filebuf *FileBuf) {

}

/**
 *	see DocParser interface
 */
func (n *nilDocParser) ParsePreview(fileCont *bytes.Buffer) []Preview {
	return nil
}

/**
 *	see DocParser interface
 */
func (n *nilDocParser) ParseCodeblock(fileCont *bytes.Buffer) []CodeBlock {
	return nil
}

/**
 *  parse public document content
 *
 *  @param `fileBuf`
 *  @return document array
 */
func ParseDocument(fileBuf *FileBuf) []Document {
	var resultDocs []Document = nil

	docsBuf := fileBuf.FinaAll(REXDocument)
	docsCount := len(docsBuf)

	if 0 == docsCount {
		return resultDocs
	}

	resultDocs = make([]Document, 0, docsCount)

	for i := 0; i < docsCount; i++ {
		docStruct := Document{}
		buf := docsBuf[i]

		lines := bytes.Split(buf, []byte("\n"))
		linesCount := len(lines)

		//  title and index parse
		indexTitleLine := lines[0]
		indexTitleMatch := REXDocIndexTitle.FindSubmatch(indexTitleLine)
		//  index 0 is source string
		//  index 1 is "///" || "/***"
		//  index 2 is "index-" index string
		//  index 3 is title

		if 4 == len(indexTitleMatch) {
			// extract title and z-index
			docStruct.SortTag = SFStringsUtil.ToInt(string(indexTitleMatch[2]))
			docStruct.Title = string(indexTitleMatch[3])
		}

		//  content parse
		contentBuf := bytes.NewBuffer(nil)
		var prefixTag []byte = nil
		prefixLen := 0

		for i := 1; i < linesCount-1; i++ {
			newLine := lines[i]

			if i == 1 {
				prefixTag = findPrefixFilterTag(newLine)
				prefixLen = len(prefixTag)
			}

			if nil != prefixTag {

				if 0 == bytes.Index(newLine, prefixTag) {
					contentBuf.Write(newLine[prefixLen:])
				} else {
					trimed := bytes.TrimSpace(newLine)
					// 有可能是空行，所需需要判断这行是否只有（ "*" || "//" ），如果不是则添加追加这一行内容
					if !bytes.Equal(trimed, _tagStar) && !bytes.Equal(trimed, _tagDSlash) {
						contentBuf.Write(newLine)
					}
				}

			} else {
				contentBuf.Write(newLine)
			}

			contentBuf.WriteByte('\n')
		}
		docStruct.Content = contentBuf.String()

		if 0 != len(docStruct.Content) {
			resultDocs = append(resultDocs, docStruct)
		}

	}

	return resultDocs
}

/**
 *  commons parse file about content
 *
 *  @param `fileBuf`
 *  @return about content
 */
func ParseAbout(fileBuf *FileBuf) []byte {
	return parseAboutAndIntro(fileBuf, REXAbout)
}

/**
 *  commons parse file introduction content
 *
 *  @param `fileBuf`
 *  @return introduction content
 */
func ParseIntro(fileBuf *FileBuf) []byte {
	return parseAboutAndIntro(fileBuf, REXIntro)
}

/**
 *  commons about intro parse
 *
 *  @param fileBuf
 *  @param rex
 */
func parseAboutAndIntro(fileBuf *FileBuf, rex *regexp.Regexp) []byte {
	var result []byte = nil
	var prefixTag []byte = nil
	prefixLen := 0

	buf := fileBuf.Find(rex)

	if 0 < len(buf) {
		appendLine := bytes.NewBuffer(nil)

		lines := bytes.Split(buf, []byte("\n"))
		linesCount := len(lines)

		for i := 1; i < linesCount-1; i++ {
			newLine := lines[i]

			if i == 1 {
				//  记录第一个前缀的标识，以第一个为准，后面的根据要求都要是相同的注释前缀。
				/**
				  (*)remove prefix tag
				  (*)
				  (*)
				*/
				prefixTag = findPrefixFilterTag(newLine)
				prefixLen = len(prefixTag)
			}

			if 0 != len(prefixTag) {

				if 0 == bytes.Index(newLine, prefixTag) {
					appendLine.Write(newLine[prefixLen:])
				} else {
					trimed := bytes.TrimSpace(newLine)
					// 有可能是空行，所需需要判断这行是否只有（ "*" || "//" ），如果不是则添加追加这一行内容
					if !bytes.Equal(trimed, _tagStar) && !bytes.Equal(trimed, _tagDSlash) {
						appendLine.Write(newLine)
					}
				}

			} else {
				appendLine.Write(newLine)
			}

			appendLine.WriteByte('\n')
		}

		if 0 < appendLine.Len() {
			result = appendLine.Bytes()
		}
	}

	return result
}

/**
 *  find prefix filter tag index
 *  //
 *  // content ("// ") is prefix tag
 *  //
 *  see var _prefixFilterTags
 */
func findPrefixFilterTag(src []byte) []byte {
	var pftCount = len(_prefixFilterTags)

	for i := 0; i < pftCount; i++ {
		prefix := _prefixFilterTags[i]
		if 0 == bytes.Index(src, prefix) {
			return prefix
		}
	}

	return nil
}
