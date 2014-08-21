//	The MIT License (MIT) - http://opensource.org/licenses/MIT
//
//	Copyright (c) 2014 slowfei
//
//  Create on 2014-08-16
//  Update on 2014-08-22
//  Email  slowfei#foxmail.com
//  Home   http://www.slowfei.com

//
package gosfdoc

import (
	"bytes"
	"fmt"
	"github.com/slowfei/gosfdoc/lang/golang"
)

const (
	APPNAME = "gosfdoc"
	VERSION = "0.0.1.000"
)

var (
	//	document parser implement interface
	_mapParser = make(map[string]DocParser)
)

/**
 *	document parser
 *
 */
type DocParser interface {

	/**
	 *	parser name
	 *
	 *	@return
	 */
	Name() string

	/**
	 *	check file
	 *	detecting whether the file is a valid file
	 *
	 *	@return true is valid file
	 */
	CheckFilepath() bool

	/**
	 *	each file the content
	 *	can be create keyword index and other operations
	 *
	 *	@param `index` while file index
	 *	@param `fileCont` file content
	 */
	EachFile(index int, fileCont *bytes.Buffer)

	/**
	 *
	 *
	 */
	// ParseDoc(fileCont *bytes.Buffer)
}

/**
 *	init
 */
func init() {
	AddParser(golang.NewParser())
}

/**
 *	add parser
 *
 *	@param parser
 */
func AddParser(parser DocParser) {
	if nil != parser {
		_mapParser[parser.Name()] = parser
	}
}

/**
 *	get parsers
 *	key is parser name
 *	value is parser implement
 *
 *	@return
 */
func MapParser() map[string]DocParser {
	return _mapParser
}

/**
 *  create config file
 *
 *	@param `path` directory path
 *	@param `lang` specify code language, nil is all language, value is parser name.
 */
func CreateConfigFile(path string, lang string) {
	fmt.Println("create :", path)
}

/**
 *	build output document
 *
 *	@param `configPath` config file path
 */
func Output(configPath string) {

}
