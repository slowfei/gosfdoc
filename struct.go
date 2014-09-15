//  The MIT License (MIT) - http://opensource.org/licenses/MIT
//
//  Copyright (c) 2014 slowfei
//
//  Create on 2014-08-22
//  Update on 2014-08-22
//  Email  slowfei#foxmail.com
//  Home   http://www.slowfei.com

//
package gosfdoc

import (
	"container/list"
	"os"
	"regexp"
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
 *  @param fileContent
 *  @param replace regexp, replace text to empty(''), call regexp.ReplaceAll func
 */
func NewFileBuf(fileContent []byte, filter *regexp.Regexp) *FileBuf {
	buf := new(FileBuf)
	if nil != filter {
		buf.buf = filter.ReplaceAll(fileContent, nil)
	} else {
		buf.buf = fileContent
	}
	return buf
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
	SortTag  int       // sort tag
	ShowText string    // show plain text
	Anchor   string    // preferably unique, with the func link
	DescText string    // markdown brief description or implement objects, can empty.
	Children []Preview //
}

/**
 *  body code block struct
 */
type CodeBlock struct {
	SortTag    int    // sort tag
	Title      string // function name or custom title
	Desc       string // description markdown text or plain text
	Code       string // show code text
	SourceLink string // source code link text
	Anchor     string // function anchor text. with the Preview link
}
