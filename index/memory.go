//  The MIT License (MIT) - http://opensource.org/licenses/MIT
//
//  Copyright (c) 2014 slowfei
//
//  Create on 2014-11-28
//  Update on 2014-12-04
//  Email  slowfei(#)foxmail.com
//  Home   http://www.slowfei.com

//
package index

/**
 *  type index data memory systems
 *  implement: IndexDB
 */
type mendb struct {
	langName string
	data     map[string]TypeInfo
}

/**
 *  init menory data
 */
func initMenDB(langName string) *mendb {
	return &mendb{langName: langName, data: make(map[string]TypeInfo)}
}

func (m *mendb) Open() error {
	return nil
}

func (m *mendb) Close() {
}

func (m *mendb) SetType(t TypeInfo) error {
	if 0 == len(t.Package) || 0 == len(t.Name) {
		return ErrInvalidIndex
	}

	key := t.Package + "+" + t.Name
	m.data[key] = t

	return nil
}

func (m *mendb) Type(packageName, typeName string) (TypeInfo, bool) {
	key := packageName + "+" + typeName
	t, ok := m.data[key]
	return t, ok
}
