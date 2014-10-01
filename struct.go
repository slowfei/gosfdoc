//  The MIT License (MIT) - http://opensource.org/licenses/MIT
//
//  Copyright (c) 2014 slowfei
//
//  Create on 2014-08-22
//  Update on 2014-09-27
//  Email  slowfei#foxmail.com
//  Home   http://www.slowfei.com

//
package gosfdoc

import (
	"container/list"
	"github.com/slowfei/gosfcore/encoding/json"
	"github.com/slowfei/gosfcore/utils/filemanager"
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
	path     string
	fileInfo os.FileInfo
	buf      []byte
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
 *  regexp find all bytes
 *
 *  @param `rex`
 *  @return
 */
func (f *FileBuf) FinaAll(rex *regexp.Regexp) [][]byte {
	return rex.FindAll(f.buf, -1)
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
	Name string // package name plain text
	Desc string // description plain text
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
