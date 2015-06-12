//  The MIT License (MIT) - http://opensource.org/licenses/MIT
//
//  Copyright (c) 2014 slowfei
//
//  Create on 2014-11-28
//  Update on 2015-01-21
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
	if 0 == len(t.PackageName) || 0 == len(t.TypeName) {
		return ErrInvalidIndex
	}

	key := t.PackagePath + "+" + t.PackageName + "+" + t.TypeName
	m.data[key] = t

	return nil
}

func (m *mendb) Type(packageName, packagePath, typeName string) (TypeInfo, bool) {
	var result TypeInfo
	var ok bool

	if 0 == len(typeName) {
		return result, false
	}

	key := packagePath + "+" + packageName + "+" + typeName
	if 0 != len(packagePath) && 0 != len(packageName) && 0 != len(typeName) {
		result, ok = m.data[key]
	} else {
		result, ok = m.data[key]
		if !ok {
			for _, v := range m.data {
				if 0 != len(packagePath) {
					if packagePath == v.PackagePath && typeName == v.TypeName {
						result = v
						ok = true
						break
					}
				} else if 0 != len(packageName) {
					if packageName == v.PackageName && typeName == v.TypeName {
						result = v
						ok = true
						break
					}
				}
			}
		}
	}

	return result, ok

}
