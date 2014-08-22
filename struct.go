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
	"bytes"
)

/**
 *  file info contain parser
 */
type DocFile struct {
	FileCont *bytes.Buffer //
	Parser   DocParser     //
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
