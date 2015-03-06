//  The MIT License (MIT) - http://opensource.org/licenses/MIT
//
//  Copyright (c) 2014 slowfei
//
//  Create on 2014-11-05
//  Update on 2015-03-07
//  Email  slowfei#foxmail.com
//  Home   http://www.slowfei.com

//	golang implement parser
package golang

import (
	"bytes"
	"fmt"
	"github.com/slowfei/gosfcore/utils/filemanager"
	"github.com/slowfei/gosfcore/utils/sub"
	"github.com/slowfei/gosfdoc"
	"github.com/slowfei/gosfdoc/index"
	"os"
	"path"
	"regexp"
	"strings"
)

const (
	GO_NAME        = "go"
	GO_SUFFIX      = ".go"
	GO_TEST_SUFFIX = "_test.go"
)

var (
	// e.g.: type Temp struct {; [0-1:prototype][2-3:comment][4-5:type define name][6-7:type name][8-9:"{"]
	REXType = regexp.MustCompile("(/\\*\\*[\\S\\s]+?\\*/\n|(?:(?:[ ]*//.*?\n)+))?type[ ]+([A-Z]\\w*)[ ]+(\\w+)[ ]*(\\{)?")
	// e.g.: package main
	REXPackage = regexp.MustCompile("package (\\w+)\\s*")
	// e.g.: /** ... */[\n]package main; //...[\n]package main
	REXPackageInfo = regexp.MustCompile("(/\\*\\*[\\S\\s]+?\\*/\n|(?:(?:[ ]*//.*?\n)+))[ ]*package \\w+")
	// e.g.: /** ... */[\n]const|var TConst = 1; //...[\n]const (
	REXDefine = regexp.MustCompile("(/\\*\\*[\\S\\s]+?\\*/\n|(?:(?:[ ]*//.*?\n)+))?[ ]*((const|var)[\\s]+(?:[A-Z].*|\\()+?)")
	// e.g: rows data
	REXRows = regexp.MustCompile("\\w+.*")

	SNRoundBrackets = SFSubUtil.NewSubNest([]byte("("), []byte(")"))
	SNBraces        = SFSubUtil.NewSubNest([]byte("{"), []byte("}"))
	SNBetweens      = []*SFSubUtil.SubNest{
		SFSubUtil.NewSubNest([]byte(`"`), []byte(`"`)),
		SFSubUtil.NewSubNest([]byte("`"), []byte("`")),
		SFSubUtil.NewSubNest([]byte(`'`), []byte(`'`)),
		SNBraces,
		SFSubUtil.NewSubNest([]byte("/*"), []byte("*/")),
		SFSubUtil.NewSubNest([]byte("//"), []byte("\n")),
	}
)

func init() {
	gosfdoc.AddParser(NewParser())
}

// golang define type
type goDType int

// golang define const
// constant and variable
const (
	goDTypeVar goDType = iota
	goDTypeConst
	goDTypeInvalid
)

// golang define constant and variable struct
type goDefine struct {
	dtype        goDType
	commentIndex []int // note content index, 0 is buffer start index, 1 is buffer end index
	contIndex    []int // content start and end index
	multiterm    bool  // true is multiterm definitions
}

// golang function struct
type goFunc struct {
	commentIndex  []int // note content index
	funcTypeIndex []int // func (type index) funcname index
	funcNameIndex []int // func [type] (func name)
	paramIndex    []int // func [type] funcname (param index)
	returnIndex   []int // func [type] funcname () (return index)
	bodyIndex     []int // function body index
}

// golang type struct
type goType struct {
	commentIndex  []int //
	typeNameIndex []int // type (type name) struct
	typeIndex     []int // type name (struct)
	bodyIndex     []int // type body index
}

/**
 *	golang parser
 */
type GolangParser struct {
	config  gosfdoc.MainConfig
	indexDB index.IndexDB
}

/**
 *	new golang parser
 */
func NewParser() *GolangParser {
	gp := new(GolangParser)
	gp.indexDB = index.CreateIndexDB(GO_NAME, index.DBTypeFile)
	return gp
}

//#pragma mark github.com/slowfei/gosfdoc.DocParser interface ---------------------------------------------------------------------

func (g *GolangParser) Name() string {
	return GO_NAME
}

/**
 *	see DocParser interface
 */
func (g *GolangParser) ParseStart(config gosfdoc.MainConfig) {
	g.config = config
	g.indexDB.Open()
}

/**
 *	see DocParser interface
 */
func (g *GolangParser) ParseEnd() {
	g.indexDB.Close()
}

/**
 *	see DocParser interface
 */
func (g *GolangParser) CheckFile(filePath string, info os.FileInfo) bool {
	result := false

	if 0 != len(filePath) && nil != info && !info.IsDir() {
		result = strings.HasSuffix(filePath, GO_SUFFIX)

		if result {
			result = !strings.HasSuffix(filePath, GO_TEST_SUFFIX)
		}
	}
	return result
}

/**
 *	see DocParser interface
 */
func (g *GolangParser) EachIndexFile(filebuf *gosfdoc.FileBuf) {
	// find type (XXXX)
	var outBetweens [][]int
	if nil == filebuf.UserData {
		outBetweens = getOutBetweens(filebuf)
		filebuf.UserData = outBetweens
	}

	tempPackagePath := ""
	tempPackageName := ""

	// find package name
	packageIndexs := filebuf.FindAllSubmatchIndex(REXPackage)
	for i := 0; i < len(packageIndexs); i++ {
		indexs := packageIndexs[i]
		if 4 == len(indexs) && !isRuleOutIndex(indexs[0], outBetweens) {
			tempPackageName = string(filebuf.SubBytes(indexs[2], indexs[3]))
			break
		}
	}

	//	查询不到package 证明是无效的文件
	if 0 == len(tempPackageName) {
		fmt.Println("InvalidFile: find less than the package name. file path:", filebuf.Path())
		return
	}

	// find package path
	gopaths := SFFileManager.GetGOPATHDirs()
	for i := 0; i < len(gopaths); i++ {
		gopath := path.Join(gopaths[i], "src")
		filebufPath := path.Dir(filebuf.Path())
		if strings.HasPrefix(filebufPath, gopath) {
			tempPackagePath = filebufPath[len(gopath)+1 : len(filebufPath)]
		}
	}

	//	无效文件提示
	if 0 == len(tempPackagePath) {
		fmt.Println("InvalidFile: is not a valid golang working environment file. file path:", filebuf.Path())
		return
	}

	//	类型查询
	typeIndexs := filebuf.FindAllSubmatchIndex(REXType)
	for i := 0; i < len(typeIndexs); i++ {
		indexs := typeIndexs[i]
		startIndex := indexs[0]
		endIndex := indexs[1]
		tempType := index.TypeInfo{}
		tempType.PackageName = tempPackageName
		tempType.PackagePath = tempPackagePath

		// type GolangParser struct { [1 27 6 18 26 27]
		// type OperateResult int [88 110 93 106 -1 -1]
		if 6 == len(indexs) && !isRuleOutIndex(startIndex, outBetweens) {

			leftBraces := indexs[4]
			rightBraces := indexs[5]
			if -1 != leftBraces && -1 != rightBraces {
				bracesIndexs := filebuf.SubNestIndex(leftBraces, SNBraces, outBetweens)
				if 2 == len(bracesIndexs) && -1 != bracesIndexs[0] && -1 != bracesIndexs[1] {
					endIndex = bracesIndexs[1]
				}
			}

			lines := filebuf.LineNumberByIndex(startIndex, endIndex)
			if -1 != lines[0] && -1 != lines[1] {
				tempType.LineStart = lines[0]
				tempType.LineEnd = lines[1]
			}

			tempType.TypeName = string(filebuf.SubBytes(indexs[2], indexs[3]))
			if 0 != len(tempType.TypeName) {
				err := g.indexDB.SetType(tempType)
				if nil != err {
					fmt.Println("IndexError:", err.Error())
				}
			}

		}
	} // End for i := 0; i < len(typeIndexs); i++ {

}

/**
 *	see DocParser interface
 */
func (g *GolangParser) ParsePreview(filebuf *gosfdoc.FileBuf) []gosfdoc.Preview {
	//	TODO
	// 1. const const{ }
	// 2. var var { }
	// 3. func
	// 4. type
	// 5. 	return func type
	// 6. 	type func

	return nil
}

/**
 *	see DocParser interface
 */
func (g *GolangParser) ParseCodeblock(filebuf *gosfdoc.FileBuf) []gosfdoc.CodeBlock {
	return nil
}

/**
 *	see DocParser interface
 */
func (n *GolangParser) ParsePackageInfo(filebuf *gosfdoc.FileBuf) string {
	result := bytes.NewBuffer(nil)

	subBytes := filebuf.FindSubmatch(REXPackageInfo)
	if 2 != len(subBytes) {
		return ""
	}

	infoLines := bytes.Split(subBytes[1], []byte("\n"))
	reCount := 0
	var prefixTag []byte = nil
	prefixLen := 0

	//	判断是否存在 /* */ 如果存在则去除首行和尾行的扫描
	if 0 != len(infoLines) &&
		0 <= bytes.Index(infoLines[0], []byte("/*")) {
		reCount = 1
	}

	//	 len(infoLines)-reCount (-1) 由于正则截取的规则中会包含一个\n符号，所以需要去除
	for i := reCount; i < len(infoLines)-reCount-1; i++ {
		infoBytes := infoLines[i]

		if i == reCount {
			prefixTag = gosfdoc.FindPrefixFilterTag(infoBytes)
			prefixLen = len(prefixTag)
		}

		if nil != prefixTag {

			if 0 == bytes.Index(infoBytes, prefixTag) {
				result.Write(infoBytes[prefixLen:])
			} else {
				trimed := bytes.TrimSpace(infoBytes)
				// 有可能是空行，所需需要判断这行是否只有（ "*" || "//" ），如果不是则添加追加这一行内容
				if !bytes.Equal(trimed, []byte("*")) && !bytes.Equal(trimed, []byte("//")) {
					result.Write(infoBytes)
				} else {
					result.WriteByte('\n')
				}
			}

		} else {
			result.Write(infoBytes)
		}

		result.WriteByte('\n')
	}

	return result.String()
}

/**
 *	find constant and variable
 *
 *	e.g:
 *	const xxx | var xxx
 *	const ( ... ) | var ( ... )
 *
 *	@param `filebuf`
 *	@param `outBetweens`
 */
func findDefine(filebuf *gosfdoc.FileBuf, outBetweens [][]int) []goDefine {
	/*
		<br> is \n

		//<br>// temp1 <br>const (, //<br>// temp1 <br>, const (, const
		const Temp3 = "3", , const Temp3 = "3", const
		// VTest1 cont<br>var VTest1  = "1", // VTest1 cont<br>, var VTest1  = "1", var
		var (, , var (, var
	*/
	var (
		CTempConst = []byte("const")
		CTempVar   = []byte("var")
	)
	var result []goDefine

	subindexs := filebuf.FindAllSubmatchIndex(REXDefine)

	for i := 0; i < len(subindexs); i++ {
		indexs := subindexs[i]

		if 8 == len(indexs) {

			contIndex_1 := indexs[4]
			contIndex_2 := indexs[5]

			// 由于index[0]首位存在注释，所以需要从主体内容开始判断排除，即是index[2]
			if -1 != contIndex_1 && -1 != contIndex_2 &&
				!isRuleOutIndex(contIndex_1, outBetweens) {

				// sub const or var byte
				dtype := goDTypeInvalid
				defineTagByte := filebuf.SubBytes(indexs[6], indexs[7])
				if bytes.Equal(defineTagByte, CTempConst) {
					dtype = goDTypeConst
				} else if bytes.Equal(defineTagByte, CTempVar) {
					dtype = goDTypeVar
				}

				if goDTypeInvalid != dtype {
					multiterm := false
					isAppend := false
					// 需要注意的是下标 -1 和 +1 的处理， 关键在于regexp截取的下标数是从1开始
					bufByte, _ := filebuf.Byte(contIndex_2 - 1)

					if '(' == bufByte {
						//	如果判断是"("则表示是多行，这时就需要查询下一个")"的目标，然后还需要判断是否全部的参数为大写开头
						contNewIndexs := filebuf.SubNestIndex(contIndex_2-1, SNRoundBrackets, outBetweens)

						if 2 == len(contNewIndexs) {
							isAllUpper := false
							subBeginIndex := contIndex_2

							// 得到 var ( "双引号里这里的命名参数" )
							bracketsBytes := filebuf.SubBytes(subBeginIndex, contNewIndexs[1]-1)
							// 由于var ( ) 是包括在圆括号里的，而outBetweens不包含"()"的排除，这里就需要重新再次查找()进行排除
							roundBracketsBetweens := filebuf.SubNestAllIndexByBetween(subBeginIndex, contNewIndexs[1]-1, SNRoundBrackets, outBetweens)

							// 获取每行的命名函数进行首字母的判断
							rowsIndexs := REXRows.FindAllIndex(bracketsBytes, -1)

							for i := 0; i < len(rowsIndexs); i++ {
								rowStartIndex := rowsIndexs[i][0]
								sourceIndex := rowStartIndex + subBeginIndex

								//	由于isRuleOutIndex判断的是fileBuf里完整索引的信息，所以需要累加上开始截取的下标数
								if -1 != rowStartIndex &&
									!isRuleOutIndex(sourceIndex, outBetweens) &&
									!isRuleOutIndex(sourceIndex, roundBracketsBetweens) {

									firstByte := bracketsBytes[rowStartIndex]

									if 'A' <= firstByte && 'Z' >= firstByte {
										isAllUpper = true
									} else {
										isAllUpper = false
									}

									if !isAllUpper {
										//	只要出现一行不为大写开头的命名就表示不通过
										break
									}
								}
							}

							if isAllUpper {
								contIndex_2 = contNewIndexs[1]
								multiterm = true
								isAppend = true
							}

						} // end if 2 == len(contNewIndexs)

					} else {
						isAppend = true
					}

					if isAppend {
						tempDefine := goDefine{}
						tempDefine.dtype = dtype
						tempDefine.contIndex = []int{contIndex_1, contIndex_2}
						tempDefine.commentIndex = []int{indexs[2], indexs[3]}
						tempDefine.multiterm = multiterm
						result = append(result, tempDefine)
					}
				}

			}

		} // end 8 == len(indexs) {

	} // end for for i := 0; i < len(subindexs); i++ {

	return result
}

/**
 *	find function
 *
 *	e.g:
 *	func [type] funcName()
 */
func findFunc(filebuf *gosfdoc.FileBuf, outBetweens [][]int) []goFunc {
	var result []goFunc = nil

	return result
}

/**
 *	find type and function
 *
 *	e.g:
 *	type xxx
 */
func findType(filebuf *gosfdoc.FileBuf, outBetweens [][]int) []goType {
	var result []goType = nil

	indexs := filebuf.FindAllSubmatchIndex(REXType)
	indexsLen := len(indexs)

	// e.g.: type Temp struct {; [0-1:prototype][2-3:comment][4-5:type define name][6-7:type name][8-9:"{"]
	if 0 != indexsLen && 10 == len(indexs[0]) {
		for i := 0; i < indexsLen; i++ {
			typeIndexs := indexs[i]

			gt := goType{}
			if -1 != typeIndexs[2] && -1 != typeIndexs[3] {
				gt.commentIndex = []int{typeIndexs[2], typeIndexs[3]}
			}

			// 正则确定存在的值
			gt.typeNameIndex = []int{typeIndexs[4], typeIndexs[5]}
			gt.typeIndex = []int{typeIndexs[6], typeIndexs[7]}

			//	判断大括号"{"
			symbolIndex := typeIndexs[8]
			bufByte, _ := filebuf.Byte(symbolIndex)
			isBraces := false

			if -1 != symbolIndex && '{' == bufByte {
				// 寻找下一个"}"
				bodyNewIndexs := filebuf.SubNestIndex(symbolIndex-1, SNBraces, outBetweens)

				if 2 == len(bodyNewIndexs) {
					isBraces = true
					gt.bodyIndex = []int{typeIndexs[4] - 5, bodyNewIndexs[1]}
				}
			}

			if !isBraces {
				//	typeIndexs[4] - 5 = "type "
				//	typeIndexs[7] = 类型名的结尾截取下标
				gt.bodyIndex = []int{typeIndexs[4] - 5, typeIndexs[7]}
			}

			result = append(result, gt)
		}
	}

	return result
}

/**
 *	获取文件排除范围的坐标范围
 *
 *	@param `filebuf`
 *	@return
 */
func getOutBetweens(filebuf *gosfdoc.FileBuf) [][]int {

	outBetweens := make([][]int, 0, 0)

	for i := 0; i < len(SNBetweens); i++ {
		tempIndexs := filebuf.SubNestAllIndex(SNBetweens[i], outBetweens)
		if 0 != len(tempIndexs) {
			outBetweens = append(outBetweens, tempIndexs...)
		}
	}

	return outBetweens
}

/**
 *	判断是否是排除坐标
 *
 *	@return 在坐标范围内返回 true
 */
func isRuleOutIndex(index int, outBetweens [][]int) bool {
	result := false

	for i := 0; i < len(outBetweens); i++ {
		indexs := outBetweens[i]
		if 2 == len(indexs) {
			s := indexs[0]
			e := indexs[1]
			if index > s && index < e {
				result = true
				break
			}
		}
	}
	return result
}
