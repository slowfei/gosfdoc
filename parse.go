package gosfdoc

import (
	"bytes"
	"regexp"
	// "fmt"
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

	_tagStar   = []byte("*")
	_tagDSlash = []byte("//")
)

/**
 *	parse public document content
 *
 *	@param `fileBuf`
 *	@return document array
 */
func ParseDocument(fileBuf *FileBuf) []Document {
	return nil
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

	aboutBuf := fileBuf.Find(rex)

	if 0 < len(aboutBuf) {
		appendLine := bytes.NewBuffer(nil)

		lines := bytes.Split(aboutBuf, []byte("\n"))
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
					//  有可能存在空行，所以需要过滤只存在 "*" || "//"
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
