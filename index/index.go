//  The MIT License (MIT) - http://opensource.org/licenses/MIT
//
//  Copyright (c) 2014 slowfei
//
//  Create on 2014-11-28
//  Update on 2014-12-03
//  Email  slowfei(#)foxmail.com
//  Home   http://www.slowfei.com

//
//  type index data storage systems
//
package index

import (
	"errors"
)

// database type
type DBType int

// database type definit
const (
	DBTypeMemory DBType = iota // disposable memory cache
	DBTypeFile                 // file type storage
)

// error definit
var (
	ErrInvalidIndex = errors.New("gosfdoc/index: Invalid unique index, package or type name nil.")
)

/**
 *  data storage interface
 */
type IndexDB interface {

	/**
	 *  Open (operating data)-> Close -> Open (operating data)-> Close...
	 *
	 *  @return `error`
	 */
	Open() error

	/**
	 *  all finished operating data can close
	 */
	Close()

	/**
	 *  save as type info, the same data is overwritten
	 *  package and name identifies a unique index
	 *
	 *  @param `t`
	 *  @return `error`
	 */
	SetType(t TypeInfo) error

	/**
	 *  by package and type name get type info
	 *
	 *  @param `packageName` pacaage
	 *  @param `typeName` type name
	 *  @return `TypeInfo`
	 *  @return `bool`
	 */
	Type(packageName, typeName string) (TypeInfo, bool)
}

/**
 *  open or create IndexDB
 *
 *  @param `langName` language name string
 *  @param `dbT`      storage type, default DBTypeMemory
 */
func CreateIndexDB(langName string, dbT DBType) IndexDB {
	//  TODO 暂时只实现了内存存储，一次性的。
	return initMenDB(langName)
}

/**
 *  index type struct
 *  type info in various languages
 */
type TypeInfo struct {
	DocHttpUrl string // document http url e.g.: http://slowfei.github.io/gosfdoc
	Package    string // package and name identifies a unique index
	Name       string // package and name identifies a unique index
	LineStart  int    // line number start
	LineEnd    int    // line number end
}
