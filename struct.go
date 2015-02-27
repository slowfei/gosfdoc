//  The MIT License (MIT) - http://opensource.org/licenses/MIT
//
//  Copyright (c) 2014 slowfei
//
//  Create on 2014-08-22
//  Update on 2015-02-27
//  Email  slowfei(#)foxmail.com
//  Home   http://www.slowfei.com

//
package gosfdoc

import (
	"container/list"
	"github.com/slowfei/gosfcore/encoding/json"
	"github.com/slowfei/gosfcore/utils/filemanager"
	"github.com/slowfei/gosfcore/utils/sub"
	"io/ioutil"
	"os"
	"regexp"
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
)

/**
 *  file content buffer
 */
type FileBuf struct {
	path       string
	fileInfo   os.FileInfo
	buf        []byte
	lineLenSum []int       // 记录每行长度的总和
	UserData   interface{} // 自定义存储数据
}

/**
 *  new file buffer
 *
 *  @param `fileContent`
 *  @param `path` file path
 *  @param `info` file info
 *  @param replace regexp, replace text to empty(''), call regexp.ReplaceAll func
 */
func NewFileBuf(fileContent []byte, path string, info os.FileInfo, filter *regexp.Regexp) *FileBuf {
	buf := new(FileBuf)
	if nil != filter {
		buf.buf = filter.ReplaceAll(fileContent, nil)
	} else {
		buf.buf = fileContent
	}

	if 0 != len(buf.buf) {
		buf.lineLenSum = make([]int, 0, 0)
		count := len(buf.buf)
		for i := 0; i < count; i++ {
			b := buf.buf[i]
			if '\n' == b {
				buf.lineLenSum = append(buf.lineLenSum, i)
			}
		}
		//	add end line
		buf.lineLenSum = append(buf.lineLenSum, count)
	}

	buf.fileInfo = info
	buf.path = path
	return buf
}

/**
 *	out file
 *
 *	@param `path` out path
 */
func (f *FileBuf) WriteFilepath(path string) error {
	return SFFileManager.WirteFilepath(path, f.buf)
}

/**
 *  regexp find bytes
 *
 *  @param `rex`
 *  @return
 */
func (f *FileBuf) Find(rex *regexp.Regexp) []byte {
	return rex.Find(f.buf)
}

/**
 *  regexp find submatch bytes
 *
 *  @param `rex`
 *  @return
 */
func (f *FileBuf) FindSubmatch(rex *regexp.Regexp) [][]byte {
	return rex.FindSubmatch(f.buf)
}

/**
 *  regexp find all bytes
 *
 *  @param `rex`
 *  @return
 */
func (f *FileBuf) FindAll(rex *regexp.Regexp) [][]byte {
	return rex.FindAll(f.buf, -1)
}

/**
 *	Regexp.FindAllSubmatch
 *
 *  @param `rex`
 *  @return
 */
func (f *FileBuf) FindAllSubmatch(rex *regexp.Regexp) [][][]byte {
	return rex.FindAllSubmatch(f.buf, -1)
}

/**
 *	Regexp.FindSubmatchIndex
 *
 *	@param `rex`
 *	@return
 */
func (f *FileBuf) FindSubmatchIndex(rex *regexp.Regexp) []int {
	return rex.FindSubmatchIndex(f.buf)
}

/**
 *	Regexp.FindAllSubmatchIndex
 *
 *	@param `rex`
 *	@return
 */
func (f *FileBuf) FindAllSubmatchIndex(rex *regexp.Regexp) [][]int {
	return rex.FindAllSubmatchIndex(f.buf, -1)
}

/**
 *  block subset return range index
 *
 *	@param `startIndex` buffer start index
 *	@param `subNest`
 *	@param `outBetweens` rule out between index
 *	@return buffer start and end index
 */
func (f *FileBuf) SubNestIndex(startIndex int, subNest *SFSubUtil.SubNest, outBetweens [][]int) []int {
	var result []int = nil

	if startIndex < len(f.buf) {
		indexs := subNest.BytesToIndex(startIndex, f.buf, outBetweens)
		if 2 == len(indexs) {
			result = []int{indexs[0] + startIndex, indexs[1] + startIndex}
		}
	}

	return result
}

/**
 *	all blocks subset
 *
 *	@param `subNest`
 *	@param `outBetweens` rule out between index
 *	@return buffer start and end index list
 */
func (f *FileBuf) SubNestAllIndex(subNest *SFSubUtil.SubNest, outBetweens [][]int) [][]int {
	return subNest.BytesToAllIndex(0, f.buf, outBetweens)
}

/**
 *	line number by begin and end index
 *
 *	@param `beginIndex` buffer byte begin index
 *	@param `endIndex`	end index
 *	@return []int [start line,end line], line number 1 start.
 */
func (f *FileBuf) LineNumberByIndex(beginIndex, endIndex int) []int {
	result := []int{-1, -1}

	if endIndex < beginIndex {
		beginIndex, endIndex = endIndex, beginIndex
	}

	isBegin := false
	isEnd := false

	/*
		text:
		abcde  [a=0,b=1,c=2,d=3,e=4,\n=5]
		fghij  [f=6,g=8,h=9,i=10,j=11]

	*/

	lineLen := len(f.lineLenSum)
	for i := 0; i < lineLen; i++ {
		lineSum := f.lineLenSum[i]
		//[0 13 14 29 30 54 65 67 68]
		if !isBegin && lineSum >= beginIndex {
			result[0] = i + 1
			isBegin = true
		}

		if !isEnd && lineSum >= endIndex {
			result[1] = i + 1
			isEnd = true
		}

		if isBegin && isEnd {
			break
		}
	}

	if !isBegin || !isEnd {
		result[0] = -1
		result[1] = -1
	}

	return result
}

/**
 *	get row content by line number 1 start.
 *
 *	@param `lineNumber` line number
 *	@param	content of the specified line number
 */
func (f *FileBuf) RowByIndex(lineNumber int) []byte {
	var result []byte = nil
	lineLen := len(f.lineLenSum)

	if 0 >= lineNumber {
		lineNumber = 1
	}

	if 0 < lineLen && lineLen >= lineNumber {
		realIndex := lineNumber - 1
		upIndex := realIndex - 1 // 上一行下标，

		/*
			text:
			12345
			67890

			line length sum:
			[0] = 5 = len("12345\n")
			[1] = 11 = len("67890")

			lineNumber = 2; 获取第二行数据 buf[5+1,11]; (5+1)除去上一节点的换行
		*/

		startIndex := 0
		endIndex := f.lineLenSum[realIndex]

		if 0 <= upIndex {
			startIndex = f.lineLenSum[upIndex] + 1
		}

		result = f.SubBytes(startIndex, endIndex)
	}

	return result
}

/**
 *	extracts the file buffer from a bytes
 *
 *	@param `beginIndex`
 *	@param `endIndex`
 *	@return bytes
 */
func (f *FileBuf) SubBytes(beginIndex, endIndex int) []byte {
	if 0 > beginIndex || 0 > endIndex || beginIndex >= endIndex {
		return nil
	}

	bufLen := len(f.buf)
	if bufLen < endIndex {
		return nil
	}

	result := f.buf[beginIndex:endIndex]
	return result
}

/**
 *	by index get file buffer byte
 *
 *	@param `index` buffer index
 *	@return `byte`
 *	@return `bool` success return true
 */
func (f *FileBuf) Byte(index int) (byte, bool) {
	var result byte
	ok := false

	if -1 != index {
		bufLen := len(f.buf)
		if index < bufLen {
			result = f.buf[index]
		}
	}

	return result, ok
}

/**
 *	get line length
 *
 *	@return int
 */
func (f *FileBuf) LineLen() int {
	return len(f.lineLenSum)
}

/**
 *  get file path
 *
 *  @return
 */
func (f *FileBuf) Path() string {
	return f.path
}

/**
 *  get file info
 *
 *  @return
 */
func (f *FileBuf) FileInfo() os.FileInfo {
	return f.fileInfo
}

/**
 *  buffer to string
 *
 *  @return
 */
func (f *FileBuf) String() string {
	return string(f.buf)
}

/**
 *  source code file
 */
type CodeFile struct {
	parser      DocParser  // file parser
	docs        []Document // current file public documents
	FileCont    *FileBuf   // file buffer content
	PrivateDoc  bool       // if private document not output
	PrivateCode bool       // if private source code not output
}

/**
 *  source code file list
 */
type CodeFiles struct {
	files *list.List
}

/**
 *  new CodeFiles
 */
func NewCodeFiles() *CodeFiles {
	cf := new(CodeFiles)
	cf.files = list.New()
	return cf
}

/**
 *  add file
 *
 *  @param file
 */
func (c *CodeFiles) addFile(file CodeFile) {
	if nil == c.files {
		c.files = list.New()
	}
	c.files.PushBack(file)
}

/**
 *  each CodeFile
 *
 *  @param `f` func return true continue
 */
func (c *CodeFiles) Each(f func(file CodeFile) bool) {
	if nil == f {
		return
	}
	for e := c.files.Front(); e != nil; e = e.Next() {
		if !f(e.Value.(CodeFile)) {
			break
		}
	}
}

/**
 *	file list storage length
 *
 *	@return file number
 */
func (c *CodeFiles) FilesLen() int {
	return c.files.Len()
}

/**
 *  output `content.json`
 */
type ContentJson struct {
	HtmlTitle string // html document title
	DocTitle  string // html top show title
	MenuTitle string // html left menu title
}

/**
 *	output write file path
 */
func (c ContentJson) WriteFilepath(path string) error {
	json, err := SFJson.NewJson(c, "", "")
	if nil != err {
		return err
	}
	return json.WriteFilepath(path, true)
}

/**
 *	markdown about
 */
type About struct {
	Content []byte
}

/**
 *	new default about
 *
 *	@return pointer type
 */
func NewDefaultAbout() *About {
	return &About{Content: _defaultAbout}
}

/**
 *	output file
 *
 *	@param `path` output full path
 *	@return
 */
func (a *About) WriteFilepath(path string) error {
	if 0 == len(a.Content) {
		a.Content = _defaultAbout
	}
	return ioutil.WriteFile(path, a.Content, 0660)
}

/**
 *	markdown intro
 */
type Intro struct {
	Content []byte
}

/**
 *	new default intro
 *
 *	@return pointer type
 */
func NewDefaultIntro() *Intro {
	return &Intro{Content: _defaultIntro}
}

/**
 *	output file
 *
 *	@param `path` output full path
 *	@return
 */
func (c *Intro) WriteFilepath(path string) error {
	if 0 == len(c.Content) {
		c.Content = _defaultIntro
	}
	return ioutil.WriteFile(path, c.Content, 0660)
}

/**
 *  package info
 */
type PackageInfo struct {
	menuName string `json:"-"` // type belongs
	Name     string // package name plain text
	Desc     string // description plain text
}

type FileLink struct {
	menuName string `json:"-"` // type belongs
	Filename string // a tag show text
	Link     string // a tag link

}

/**
 *  document struct info
 */
type Document struct {
	SortTag int    // sort tag
	Title   string // title plain text
	Content string // markdown text or plain text
}

/**
 *  preview struct info
 */
type Preview struct {
	SortTag  string // sort tag
	Level    int    // hierarchy level show. 0 is >, 1 is >>, 3 is >>> ...(markdown syntax)
	ShowText string // show plain text
	Anchor   string // preferably unique, with the func link
	DescText string // markdown brief description or implement objects, can empty.
}

/**
 *  body code block struct
 */
type CodeBlock struct {
	SortTag        string // sort tag
	MenuTitle      string // left navigation menu title
	Title          string // function name or custom title
	Anchor         string // function anchor text.
	Desc           string // description markdown text or plain text
	Code           string // show code text
	CodeLang       string // source code lang type string
	SourceFileName string // source code file name
	FileLines      []int  // block where the file line [5,10] is L5-L10
}

/**
 *	Preview,CodeBlock,Document sort implement
 */
type SortSet struct {
	previews   []Preview
	documents  []Document
	codeBlocks []CodeBlock
}

/**
 *	sort Len() implement
 */
func (s SortSet) Len() int {

	if 0 != len(s.previews) {
		return len(s.previews)
	} else if 0 != len(s.documents) {
		return len(s.documents)
	} else if 0 != len(s.codeBlocks) {
		return len(s.codeBlocks)
	} else {
		return 0
	}

}

/**
 *	sort Less(...) implement
 */
func (s SortSet) Less(i, j int) bool {

	if 0 != len(s.previews) {
		return s.previews[i].SortTag < s.previews[j].SortTag
	} else if 0 != len(s.documents) {
		return s.documents[i].SortTag < s.documents[j].SortTag
	} else if 0 != len(s.codeBlocks) {
		return s.codeBlocks[i].SortTag < s.codeBlocks[j].SortTag
	} else {
		return false
	}

}

/**
 *	sort Swap(...) implement
 */
func (s SortSet) Swap(i, j int) {

	if 0 != len(s.previews) {
		s.previews[i], s.previews[j] = s.previews[j], s.previews[i]
	} else if 0 != len(s.documents) {
		s.documents[i], s.documents[j] = s.documents[j], s.documents[i]
	} else if 0 != len(s.codeBlocks) {
		s.codeBlocks[i], s.codeBlocks[j] = s.codeBlocks[j], s.codeBlocks[i]
	}

}
