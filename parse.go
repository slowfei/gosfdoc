//  The MIT License (MIT) - http://opensource.org/licenses/MIT
//
//  Copyright (c) 2014 slowfei
//
//  Create on 2014-09-10
//  Update on 2014-10-07
//  Email  slowfei#foxmail.com
//  Home   http://www.slowfei.com

//
package gosfdoc

import (
	"bytes"
	"fmt"
	"github.com/slowfei/gosfcore/utils/strings"
	"os"
	"regexp"
	"strings"
)

const (
	DOC_FILE_SUFFIX = ".dc"      // document file suffix(document comments)
	NIL_DOC_NAME    = "document" // nilDocParser struct use

)

var (

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
func (n *nilDocParser) ParsePreview(filebuf *FileBuf) []Preview {
	return nil
}

/**
 *	see DocParser interface
 */
func (n *nilDocParser) ParseCodeblock(filebuf *FileBuf) []CodeBlock {
	return nil
}

/**
 *	see DocParser interface
 */
func (n *nilDocParser) ParsePackageInfo(filebuf *FileBuf) string {
	return ""
}

/**
 *	prese Preview Document CodeBlock array to markdown
 *
 *	@param `documents` after sorting
 *	@param `previews`  after sorting
 *	@param `blocks`	   after sorting
 *	@param `filesName` file names
 *	@param `relPath`   before code file name join path
 *	@return bytes
 */
func ParseMarkdown(documents []Document, previews []Preview, blocks []CodeBlock,
	filesName []string, relPath string) []byte {

	relPath = strings.TrimPrefix(relPath, "/")
	relPath = strings.TrimSuffix(relPath, "/")
	joinSymbol := ""
	if 0 != len(relPath) {
		joinSymbol = "/"
	}

	buf := bytes.NewBuffer([]byte{'\n'})

	for _, doc := range documents {
		// ## Welcome to gosfdoc
		// ------
		//
		// markdown syntax content
		//
		buf.WriteString("## " + doc.Title + "\n------\n")
		buf.WriteString(doc.Content)
		buf.WriteByte('\n')
	}

	if 0 != len(previews) {
		// ## Preview
		// ------
		// > [func Main()][#]<br/>
		// > [type TestStruct struct][#]<br/>
		// > implements: [Test][#]<br/>
		// >>[func (* TestStruct) hello(str string) string](#func_TestStruct.hello)<a name="preview_TestStruct.hello"><a/><br/>
		// >>[func (* TestStruct) hello2() string][#]<br/>
		buf.WriteString("## Preview\n------\n")
		for _, pre := range previews {
			buf.WriteByte('\n')
			angle := ">"
			for i := 0; i < pre.Level; i++ {
				angle += ">"
			}

			anchor := ""
			if 0 == len(pre.Anchor) {
				anchor = "(javascript:;)"
				// [show text](javascript:;)
			} else {
				anchor = fmt.Sprintf("(#f_%s)<a name=\"p_%s\"><a/>", pre.Anchor, pre.Anchor)
				// [show test](#f_anchor)<a name="p_anchor"><a/><br/>
			}
			buf.WriteString(angle + " [" + pre.ShowText + "]" + anchor + "\n")

			if 0 != len(pre.DescText) {
				buf.WriteByte('\n')
				buf.WriteString(angle + " " + pre.DescText + "\n")
			}
		}
		buf.WriteByte('\n')
	}

	// out associate files
	if 0 != len(filesName) {
		// ###Package files
		// [a.go](#) [b.go](#) [c.go](#)
		buf.WriteString("<br/>\n### Directory files\n")
		for _, name := range filesName {
			joinPath := relPath + joinSymbol + name
			buf.WriteString(fmt.Sprintf("[%s](%s) ", name, joinPath))
		}
		buf.WriteByte('\n')
	}

	if 0 != len(blocks) {
		buf.WriteByte('\n')
		isLinkCode := 0 != len(filesName)
		// ## Func Details
		// ------
		// ###[func (*TestStruct) hello](src.html?f=gosfdoc.go#L17) <a name="func_TestStruct.hello"><a/> [↩](#preview_TestStruct.hello)|[#](#func_TestStruct.hello)
		// > 函数介绍描述<br/>
		// > <br/>
		// > @param `str` 字符串传递<br/>
		// > @return `v1` 返回参数v1<br/>
		// > @return v2 返回参数v2<br/>
		//
		// ```go
		// func (* TestStruct) hello(str string) (v1,v2 string)
		// ```
		currentMenuTitle := ""

		for _, block := range blocks {

			if 0 != len(block.MenuTitle) && currentMenuTitle != block.MenuTitle {
				buf.WriteString("## " + block.MenuTitle + "\n------\n")
				currentMenuTitle = block.MenuTitle
			}

			if 0 != len(block.Title) {
				joinPath := "javascript:;"
				if isLinkCode && 0 != len(block.SourceFileName) {
					lineStr := ""

					lineLen := len(block.FileLines)
					if 1 == lineLen {
						lineStr = fmt.Sprintf("#L%d", block.FileLines[0])
					} else if 2 == lineLen {
						lineStr = fmt.Sprintf("#L%d-L%d", block.FileLines[0], block.FileLines[1])
					}

					joinPath = "src.html?f=" + relPath + joinSymbol + block.SourceFileName + lineStr
				}

				anchor := ""
				if 0 != len(block.Anchor) {
					anchor = fmt.Sprintf("<a name=\"f_%s\"><a/> [↩](#p_%s) | [#](#f_%s)", block.Anchor, block.Anchor, block.Anchor)
				}

				buf.WriteString(fmt.Sprintf("### [%s](%s) %s\n", block.Title, joinPath, anchor))
			}

			if 0 != len(block.Desc) {
				//	content description
				descLines := strings.Split(block.Desc, "\n")
				for _, desc := range descLines {
					buf.WriteString(fmt.Sprintf("> %s<br/>\n", desc))
				}
				buf.WriteByte('\n')
			}

			// code
			if 0 != len(block.Code) {
				buf.WriteString(fmt.Sprintf("\n```%s\n%s\n```\n\n", block.CodeLang, block.Code))
			}

			buf.WriteByte('\n')
		}
	}
	return buf.Bytes()
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
func ParseAbout(fileBuf *FileBuf) *About {
	data := parseAboutAndIntro(fileBuf, REXAbout)

	var result *About = nil
	if 0 != len(data) {
		result = &About{Content: data}
	}

	return result
}

/**
 *  commons parse file introduction content
 *
 *  @param `fileBuf`
 *  @return introduction content
 */
func ParseIntro(fileBuf *FileBuf) *Intro {
	data := parseAboutAndIntro(fileBuf, REXIntro)

	var result *Intro = nil
	if 0 != len(data) {
		result = &Intro{Content: data}
	}

	return result
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
