//  The MIT License (MIT) - http://opensource.org/licenses/MIT
//
//  Copyright (c) 2014 slowfei
//
//  Create on 2014-11-05
//  Update on 2015-06-12
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
	// e.g.: func (t type)funcname(params) return val{;
	// https://www.debuggex.com/r/Su6Ns1LhVxpfD_Di
	// [0-1:prototype] [2-3:comment or null] [4-5:func type or null]
	// [6-7:func name] [8-9:func params] [10-11:single return value or null]
	// [12-13:multi return value or null] [14-15:"{"]
	REXFunc = regexp.MustCompile(`(/\*\*[\s]*(?:[ ]+.*?\n)+[ ]*\*/[ ]*\n|(?:(?:[ ]*//.*?\n)+))?func[ ]*(?:\(([\w \*\n\r\.\[\]]*)\))?[ \n]*([A-Z]\w*)[ \n]*(?:\(([\w ,\*\n\r\.\{\}\[\]]*)\))+?[ \n]*(?:([\w\.\*\{\}\[\]]*)|(?:\(([\w ,\*\n\r\.\{\}\[\]]*)\)))?[ \n]*({)`)
	// e.g.: type Temp struct {; [0-1:prototype][2-3:comment][4-5:type define name][6-7:type name][8-9:"{"]
	REXType = regexp.MustCompile(`(/\*\*[\s]*(?:[ ]+.*?\n)+[ ]*\*/[ ]*\n|(?:(?:[ ]*//.*?\n)+))?type[ ]+([A-Z]\w*)[ ]+(\w+)[ ]*(\{)?`)
	// e.g.: package main
	REXPackage = regexp.MustCompile(`package (\w+)\s*`)
	// e.g.: /** ... */[\n]package main; //...[\n]package main
	REXPackageInfo = regexp.MustCompile(`(/\*\*[\s]*(?:[ ]+.*?\n)+[ ]*\*/[ ]*\n|(?:(?:[ ]*//.*?\n)+))[ ]*package \w+`)
	// e.g.: /** ... */[\n]const|var TConst = 1; //...[\n]const (
	// https://www.debuggex.com/r/2qeKD9vwnBjkgORT
	// [0-1:prototype] [2-3:comment or null] [4-5:const|var] [6-7: define name]
	REXDefine = regexp.MustCompile(`(/\*\*[\s]*(?:[ ]+.*?\n)+[ ]*\*/[ ]*\n|(?:(?:[ ]*//.*?\n)+))?[ ]*(const|var)\s+(?:\(|(?:([A-Z]\w*)\s*=.+))`)
	// e.g: rows data
	REXRows = regexp.MustCompile("\\w+.*")

	SNRoundBrackets = SFSubUtil.NewSubNest([]byte("("), []byte(")"))
	SNBraces        = SFSubUtil.NewSubNest([]byte("{"), []byte("}"))
	SNBetweens      = []*SFSubUtil.SubNest{
		SFSubUtil.NewSubNest([]byte(`"`), []byte(`"`)),
		SFSubUtil.NewSubNest([]byte(`'`), []byte(`'`)),
		SFSubUtil.NewSubNest([]byte("`"), []byte("`")),
		SFSubUtil.NewSubNest([]byte("/*"), []byte("*/")),
		SFSubUtil.NewSubNotNest([]byte("//"), []byte("\n")),
		SNBraces,
	}
)

var (
	// name text replacer
	_nameTextReplacer = strings.NewReplacer(
		"\t", "_",
		" ", "_",
		".", "-",
		"*", "+",
		"(", "_",
		")", "_",
		"[", "_",
		"]", "_",
		",", "",
		",", "",
		"\"", "_",
		"/", "_",
		"\\", "_",
	)
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
	nameIndexs   []int // define name (var DEFINDE_NAME) index
	bodyIndex    []int // define body index
	multiterm    bool  // true is multiterm definitions
}

// golang function struct
type goFunc struct {
	commentIndex  []int // note content index
	funcTypeIndex []int // func (type index) funcname index
	funcNameIndex []int // func [type] (func name index)()
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

// goType goFunc goDefine struct memory
type goMemory struct {
	packageName string
	packagePath string
	defines     []goDefine
	funcs       []goFunc
	types       []goType
	outBetweens [][]int
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
	outBetweens := getOutBetweens(filebuf)

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
		fmt.Println("InvalidFile: is not a valid golang working environment file, may current project not set GOPATH. file path:", filebuf.Path())
		return
	}

	// 类型查询，并且保存到存储区域中
	goTypes := findType(filebuf, outBetweens)
	for i := 0; i < len(goTypes); i++ {
		goType := goTypes[i]

		tempType := index.TypeInfo{}
		tempType.PackageName = tempPackageName
		tempType.PackagePath = tempPackagePath

		lines := filebuf.LineNumberByIndex(goType.bodyIndex[0], goType.bodyIndex[1])
		tempType.LineStart = lines[0]
		tempType.LineEnd = lines[1]

		tempType.TypeName = string(filebuf.SubBytes(goType.typeNameIndex[0], goType.typeNameIndex[1]))
		if 0 != len(tempType.TypeName) {
			err := g.indexDB.SetType(tempType)
			if nil != err {
				fmt.Println("IndexError:", err.Error())
			}
		}
	}

	// 查询 定义的类型
	goDefines := findDefine(filebuf, outBetweens)

	// 查询定义开放的函数
	goFuncs := findFunc(filebuf, outBetweens)

	//	存储特定数据，以便实现使用
	if nil == filebuf.UserData {
		gomem := goMemory{}
		gomem.defines = goDefines
		gomem.funcs = goFuncs
		gomem.types = goTypes
		gomem.outBetweens = outBetweens
		gomem.packageName = tempPackageName
		gomem.packagePath = tempPackagePath
		filebuf.UserData = gomem
	}
}

/**
 *	ParsePreview 与 ParseCodeblock 排序的统一处理
 */
func sortPreviewAndCodeblock(gd *goDefine, gt *goType, gf *goFunc) string {
	// 1. const const 		 	 --sortTag: "0_0"
	// 2. const const( )  		 --sortTag: "0_1"
	// 3. var var   	  		 --sortTag: "1_0"
	// 4. var var ( )	  		 --sortTag: "1_1"
	// 5. func 					 --sortTag: "2_0"
	// 6. type 					 --sortTag: "3_0_structtype_1_"
	// 7. 	return func type 	 --sortTag: "3_0_structtype_2_funcname"
	// 8. 	type func 			 --sortTag: "3_0_structtype_3_funcname"
	sortTag := ""

	if nil != gd {
		if goDTypeConst == gd.dtype {
			sortTag = "0_"
		} else {
			sortTag = "1_"
		}

		if gd.multiterm {
			sortTag += "1"
		} else {
			sortTag += "0"
		}

	} else if nil != gt {
		sortTag = "3_0_"
	} else if nil != gf {
		if 0 == len(gf.funcTypeIndex) {
			sortTag = "2_0_func_"
		} else {
			// 7. return func type  --sortTag: "3_0_structtype_0_funcname" 这个排序处理不了，需要在自行处理
			sortTag = "3_0_type_"
		}
	}

	return sortTag

}

/**
 *	parse define private set
 *
 *	@return anchor and show text and sort tag string
 */
func parseDefineAnchorShowTextSortTag(define goDefine, filebuf *gosfdoc.FileBuf) (anchor, showText, sortTag string) {

	if define.multiterm {
		switch define.dtype {
		case goDTypeConst:
			showText = "const ( "
		case goDTypeVar:
			showText = "var ( "
		}

		firstName := filebuf.SubBytes(define.nameIndexs[0], define.nameIndexs[1])
		nameIndexsLen := len(define.nameIndexs)
		endName := filebuf.SubBytes(define.nameIndexs[nameIndexsLen-2], define.nameIndexs[nameIndexsLen-1])

		// 超出截取
		if 32 < len(firstName) {
			showText += string(firstName[:32]) + "... ..."
		} else {
			showText += string(firstName) + "... ..."
		}

		if 32 < len(endName) {
			showText += string(endName[len(endName)-32:]) + " )"
		} else {
			showText += string(endName) + " )"
		}

	} else {

		switch define.dtype {
		case goDTypeConst:
			showText = "const "
		case goDTypeVar:
			showText = "var "
		}

		dName := filebuf.SubBytes(define.nameIndexs[0], define.nameIndexs[1])

		// 大于64个字符的处理
		if 64 < len(dName) {
			showText += string(dName[:64]) + "..."
		} else {
			showText += string(dName)
		}

	}

	anchor = _nameTextReplacer.Replace(showText)
	sortTag = sortPreviewAndCodeblock(&define, nil, nil) + "_" + anchor

	return
}

/**
 *	parse type private set
 *
 *	@return anchor and show text and sort tag string
 */
func parseTypeAnchorShowTextSortTag(got goType, filebuf *gosfdoc.FileBuf) (anchor, showText, sortTag string) {
	typeName := string(filebuf.SubBytes(got.typeNameIndex[0], got.typeNameIndex[1]))

	showText = "type " + typeName + " " + string(filebuf.SubBytes(got.typeIndex[0], got.typeIndex[1]))
	anchor = _nameTextReplacer.Replace(showText)

	sortText := "type " + typeName
	sortText = _nameTextReplacer.Replace(sortText)
	sortTag = sortPreviewAndCodeblock(nil, &got, nil) + sortText + "_1_"

	return
}

/**
 *	parse func private set
 *
 *	@return `anchor`
 *	@return `showText` 返回函数主体名称，包括函数类型、函数名、参数、返回值
 *	@return `showText` 返回函数名称，包括函数类型、函数名
 *	@return `sortTag`
 *	@return `level`
 */
func parseFuncAnchorShowTextSortTag(gof goFunc, filebuf *gosfdoc.FileBuf, g *GolangParser, gomen goMemory) (anchor, showText, showText2, sortTag string, level int) {

	// get func string value
	funcType := ""
	funcName := ""
	funcParam := ""
	funcReturn := ""
	if 0 != len(gof.funcTypeIndex) {
		funcType = string(filebuf.SubBytes(gof.funcTypeIndex[0], gof.funcTypeIndex[1]))

		//	除去类型的命名 func (t *TStruct) func 将"t "去除
		starIndex := strings.Index(funcType, "*")
		if -1 != starIndex {
			funcType = funcType[starIndex:]
		} else {
			blankIndex := strings.Index(funcType, " ")
			if -1 != blankIndex {
				funcType = funcType[blankIndex+1:]
			}
		}

	}
	if 0 != len(gof.paramIndex) {
		funcParam = string(filebuf.SubBytes(gof.paramIndex[0], gof.paramIndex[1]))
	}
	if 0 != len(gof.returnIndex) {
		funcReturn = string(filebuf.SubBytes(gof.returnIndex[0], gof.returnIndex[1]))
	}
	funcName = string(filebuf.SubBytes(gof.funcNameIndex[0], gof.funcNameIndex[1]))

	// preview struct value
	showText = "func "
	showText2 = "func "
	anchor = ""
	sortTag = ""
	level = 0

	// set showText
	if 0 != len(funcType) {
		showText += "(" + funcType + ") "
		showText2 += "(" + funcType + ") "
	}
	showText += funcName
	showText2 += funcName

	if 0 != len(funcParam) {
		showText += "(" + funcParam + ") "
	} else {
		showText += "() "
	}
	if 0 != len(funcReturn) {
		if -1 != strings.Index(funcReturn, ",") {
			showText += "(" + funcReturn + ")"
		} else {
			showText += funcReturn
		}
	}

	// set anchor
	anchor = _nameTextReplacer.Replace(showText)

	// set level and sortTag
	// sortTag主要参考 sortPreviewAndCodeblock
	// 5. func 					 --sortTag: "2_0_funcname"
	// 6. type 					 --sortTag: "3_0_0_structtype_1_"
	// 7. 	return func type 	 --sortTag: "3_0_structtype_2_funcname"
	// 8. 	type func 			 --sortTag: "3_0_structtype_3_funcname"

	if 0 != len(gof.funcTypeIndex) {

		if 0 == strings.Index(funcType, "*") {
			funcType = funcType[1:]
		}

		level = 1
		sortTag = sortPreviewAndCodeblock(nil, nil, &gof) + funcType + "_3_" + funcName
	} else if 0 != len(funcReturn) && -1 == strings.Index(funcReturn, ",") && -1 == strings.Index(funcReturn, ".") {
		// 判断返回值为一个参数

		//	除去类型的命名 func fname() (t string) 将"t "去除或指针符号
		if index := strings.Index(funcReturn, "*"); -1 != index {
			funcReturn = funcReturn[index+1:]
		} else if index := strings.Index(funcReturn, "]"); -1 != index {
			funcReturn = funcReturn[index+1:]
		} else if index := strings.Index(funcReturn, " "); -1 != index {
			funcReturn = funcReturn[index+1:]
		}

		// 查询方法的返回值是否在当前包和路径中查询到，查询到的话就归类到该类型下显示。
		if _, ok := g.indexDB.Type(gomen.packageName, gomen.packagePath, funcReturn); ok {
			level = 1

			gof.funcTypeIndex = make([]int, 1) // 由于排除逻辑处理有些困难，所以做个处理。看一下sortPreviewAndCodeblock处理排序的要求就明白了。
			sortTag = sortPreviewAndCodeblock(nil, nil, &gof) + funcReturn + "_2_" + funcName
			gof.funcTypeIndex = nil
		} else {
			sortTag = sortPreviewAndCodeblock(nil, nil, &gof) + funcName
		}
	} else {
		sortTag = sortPreviewAndCodeblock(nil, nil, &gof) + funcName
	}

	return
}

/**
 *	see DocParser interface
 */
func (g *GolangParser) ParsePreview(filebuf *gosfdoc.FileBuf) []gosfdoc.Preview {

	var result []gosfdoc.Preview = nil

	switch gomen := filebuf.UserData.(type) {
	case goMemory:
		{
			definesLen := len(gomen.defines)
			typesLen := len(gomen.types)
			funcsLen := len(gomen.funcs)
			result = make([]gosfdoc.Preview, 0, definesLen+typesLen+funcsLen)

			// add const and var
			for i := 0; i < definesLen; i++ {
				define := gomen.defines[i]
				if 0 == len(define.bodyIndex) {
					continue
				}

				pre := gosfdoc.Preview{}
				pre.Level = 0

				anchor, showText, sortTag := parseDefineAnchorShowTextSortTag(define, filebuf)

				pre.Anchor = anchor
				pre.ShowText = showText
				pre.SortTag = sortTag

				result = append(result, pre)
			}

			// add struct func
			for i := 0; i < typesLen; i++ {
				got := gomen.types[i]
				if 0 == len(got.bodyIndex) {
					continue
				}

				anchor, showText, sortTag := parseTypeAnchorShowTextSortTag(got, filebuf)

				pre := gosfdoc.Preview{}
				pre.ShowText = showText
				pre.Anchor = anchor
				pre.SortTag = sortTag
				pre.Level = 0
				result = append(result, pre)
			}

			// add func
			for i := 0; i < funcsLen; i++ {
				gof := gomen.funcs[i]
				if 0 == len(gof.bodyIndex) {
					continue
				}

				anchor, showText, _, sortTag, level := parseFuncAnchorShowTextSortTag(gof, filebuf, g, gomen)

				pre := gosfdoc.Preview{}
				pre.ShowText = strings.Replace(showText, "*", "\\*", -1)
				pre.Anchor = anchor
				pre.SortTag = sortTag
				pre.Level = level
				result = append(result, pre)

			}
		}
	}

	return result
}

/**
 *	see DocParser interface
 */
func (g *GolangParser) ParseCodeblock(filebuf *gosfdoc.FileBuf) []gosfdoc.CodeBlock {

	var result []gosfdoc.CodeBlock = nil

	// SortTag        string // sort tag
	// MenuTitle      string // left navigation menu title; Constants、Variables、Func Details
	// Title          string // function name or custom title
	// Anchor         string // function anchor text.
	// Desc           string // description markdown text or plain text
	// Code           string // show code text
	// CodeLang       string // source code lang type string
	// SourceFileName string // source code file name
	// FileLines      []int  // block where the file line [5,10] is L5-L10

	switch gomen := filebuf.UserData.(type) {
	case goMemory:
		{
			definesLen := len(gomen.defines)
			typesLen := len(gomen.types)
			funcsLen := len(gomen.funcs)
			result = make([]gosfdoc.CodeBlock, 0, definesLen+typesLen+funcsLen)

			for i := 0; i < definesLen; i++ {
				define := gomen.defines[i]
				if 0 == len(define.bodyIndex) {
					continue
				}

				menuTitle := ""
				title := "source code"
				desc := ""
				anchor := ""
				sortTag := ""
				code := ""
				codeLang := "go"
				var fileLines []int = nil
				sourceFileName := ""

				if nil != filebuf.FileInfo() {
					sourceFileName = filebuf.FileInfo().Name()
				}

				// set menuTitle
				switch define.dtype {
				case goDTypeVar:
					menuTitle = "Variables"
				case goDTypeConst:
					menuTitle = "Constants"
				default:
					menuTitle = "Error Title"
				}

				// set desc
				if 0 != len(define.commentIndex) {
					comt := handleComment(filebuf.SubBytes(define.commentIndex[0], define.commentIndex[1]))
					desc = string(comt)
				}

				// set anchor and sortTag
				anchor, _, sortTag = parseDefineAnchorShowTextSortTag(define, filebuf)

				// set code
				code = string(filebuf.SubBytes(define.bodyIndex[0], define.bodyIndex[1]))
				// set fileLines
				fileLines = filebuf.LineNumberByIndex(define.bodyIndex[0], define.bodyIndex[1])

				codeBlock := gosfdoc.CodeBlock{}
				codeBlock.SortTag = sortTag
				codeBlock.MenuTitle = menuTitle
				codeBlock.Title = title
				codeBlock.Anchor = anchor
				codeBlock.Desc = desc
				codeBlock.Code = code
				codeBlock.CodeLang = codeLang
				codeBlock.SourceFileName = sourceFileName
				codeBlock.FileLines = fileLines

				result = append(result, codeBlock)
			}

			for i := 0; i < typesLen; i++ {
				got := gomen.types[i]
				if 0 == len(got.bodyIndex) {
					continue
				}

				menuTitle := "Func Details"
				title := ""
				desc := ""
				anchor := ""
				sortTag := ""
				code := ""
				codeLang := "go"
				var fileLines []int = nil
				sourceFileName := ""

				if nil != filebuf.FileInfo() {
					sourceFileName = filebuf.FileInfo().Name()
				}

				anchor, title, sortTag = parseTypeAnchorShowTextSortTag(got, filebuf)

				// set comment
				if 0 != len(got.commentIndex) {
					comt := handleComment(filebuf.SubBytes(got.commentIndex[0], got.commentIndex[1]))
					desc = string(comt)
				}

				// set code
				code = string(filebuf.SubBytes(got.bodyIndex[0], got.bodyIndex[1]))
				// set fileLines
				fileLines = filebuf.LineNumberByIndex(got.bodyIndex[0], got.bodyIndex[1])

				codeBlock := gosfdoc.CodeBlock{}
				codeBlock.SortTag = sortTag
				codeBlock.MenuTitle = menuTitle
				codeBlock.Title = strings.Replace(title, "*", "\\*", -1)
				codeBlock.Anchor = anchor
				codeBlock.Desc = desc
				codeBlock.Code = code
				codeBlock.CodeLang = codeLang
				codeBlock.SourceFileName = sourceFileName
				codeBlock.FileLines = fileLines

				result = append(result, codeBlock)
			}

			for i := 0; i < funcsLen; i++ {
				gof := gomen.funcs[i]
				if 0 == len(gof.bodyIndex) {
					continue
				}

				menuTitle := "Func Details"
				codeTitle := ""
				showTitle := ""
				desc := ""
				anchor := ""
				sortTag := ""
				code := ""
				codeLang := "go"
				var fileLines []int = nil
				sourceFileName := ""

				if nil != filebuf.FileInfo() {
					sourceFileName = filebuf.FileInfo().Name()
				}

				// set title anchor sotrTag
				anchor, codeTitle, showTitle, sortTag, _ = parseFuncAnchorShowTextSortTag(gof, filebuf, g, gomen)

				// set comment
				if 0 != len(gof.commentIndex) {
					comt := handleComment(filebuf.SubBytes(gof.commentIndex[0], gof.commentIndex[1]))
					desc = string(comt)
				}

				// set code
				code = codeTitle + " { ...... }"

				// set fileLines
				fileLines = filebuf.LineNumberByIndex(gof.bodyIndex[0], gof.bodyIndex[1])

				codeBlock := gosfdoc.CodeBlock{}
				codeBlock.SortTag = sortTag
				codeBlock.MenuTitle = menuTitle
				codeBlock.Title = strings.Replace(showTitle, "*", "\\*", -1)
				codeBlock.Anchor = anchor
				codeBlock.Desc = desc
				codeBlock.Code = code
				codeBlock.CodeLang = codeLang
				codeBlock.SourceFileName = sourceFileName
				codeBlock.FileLines = fileLines
				result = append(result, codeBlock)
			}
		}

	}

	return result
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
		//<br>// temp1 <br>const (, //<br>// temp1 <br>, const,
		const Temp3 = "3", , const, Temp3
		// VTest1 cont<br>var VTest1  = "1", // VTest1 cont<br>, var, VTest1
		var (, , var,

		// [0-1:prototype] [2-3:comment or null] [4-5:const|var] [6-7: define name]
	*/
	var (
		CTempConst = []byte("const")
		CTempVar   = []byte("var")
	)
	var result []goDefine

	subindexs := filebuf.FindAllSubmatchIndex(REXDefine)

	for i := 0; i < len(subindexs); i++ {
		indexs := subindexs[i]
		// 由于index[0]首位存在注释，所以需要从主体内容开始判断排除，即是index[2]

		if 8 == len(indexs) {

			pteIndex_2 := indexs[1]
			corvIndex_1 := indexs[4]
			corvIndex_2 := indexs[5]

			// 判断 [4-5:const|var] ,是否是注释内存或则存在
			if -1 != corvIndex_1 && -1 != corvIndex_2 &&
				!isRuleOutIndex(corvIndex_1, outBetweens) {

				// sub const or var byte
				dtype := goDTypeInvalid
				defineTagByte := filebuf.SubBytes(indexs[4], indexs[5])
				if bytes.Equal(defineTagByte, CTempConst) {
					dtype = goDTypeConst
				} else if bytes.Equal(defineTagByte, CTempVar) {
					dtype = goDTypeVar
				}

				if goDTypeInvalid != dtype {
					multiterm := false
					isAppend := false
					var dNameIndexs []int = nil // 定义类型名的下标存储

					// 截取 [0-1:prototype]原型的末尾最后一个字符，判断"var ("是否是括号
					bufByte, _ := filebuf.Byte(pteIndex_2 - 1)

					if '(' == bufByte {
						//	如果判断是"("则表示是多行，这时就需要查询下一个")"的目标，然后还需要判断是否全部的参数为大写开头
						contNewIndexs := filebuf.SubNestIndex(pteIndex_2-1, SNRoundBrackets, outBetweens)

						if 2 == len(contNewIndexs) {
							isAllUpper := false
							subBeginIndex := pteIndex_2

							// 得到 var ( "双引号里这里的命名参数" )
							bracketsBytes := filebuf.SubBytes(subBeginIndex, contNewIndexs[1]-1)
							// 由于var ( ) 是包括在圆括号里的，而outBetweens不包含"()"的排除，这里就需要重新再次查找()进行排除
							roundBracketsBetweens := filebuf.SubNestAllIndexByBetween(subBeginIndex, contNewIndexs[1]-1, SNRoundBrackets, outBetweens)

							// 获取每行的命名函数进行首字母的判断
							rowsIndexs := REXRows.FindAllIndex(bracketsBytes, -1)
							rowsLen := len(rowsIndexs)
							if 0 != rowsLen {
								//	定义存储容量，由于存储的是下标，需要开始和结尾，所以在查询得到的行数乘以2
								dNameIndexs = make([]int, 0, rowsLen*2)
							}

							for i := 0; i < rowsLen; i++ {
								rowStartIndex := rowsIndexs[i][0]
								sourceIndex := rowStartIndex + subBeginIndex

								//	由于isRuleOutIndex判断的是fileBuf里完整索引的信息，所以需要累加上开始截取的下标数
								if -1 != rowStartIndex &&
									!isRuleOutIndex(sourceIndex, outBetweens) &&
									!isRuleOutIndex(sourceIndex, roundBracketsBetweens) {

									firstByte := bracketsBytes[rowStartIndex]

									if 'A' <= firstByte && 'Z' >= firstByte {
										isAllUpper = true
										// 获取定义名称,名称就控制在64个字符内
										dNameIndexs = append(dNameIndexs, sourceIndex)
										dNameIndexs = append(dNameIndexs, sourceIndex+1)

										for j := rowStartIndex + 1; j < rowStartIndex+63; j++ {
											b := bracketsBytes[j]
											if '\n' == b || '\r' == b || ' ' == b || '=' == b {
												// 判断名称的结尾
												dNameIndexs[len(dNameIndexs)-1] = subBeginIndex + j
												break
											}

										}

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
								pteIndex_2 = contNewIndexs[1]
								multiterm = true
								isAppend = true
							}

						} // end if 2 == len(contNewIndexs)

					} else {
						isAppend = true
						dNameIndexs = []int{indexs[6], indexs[7]}
					}

					if isAppend {
						tempDefine := goDefine{}
						tempDefine.dtype = dtype
						// [4-5:const|var] [4]开始 至 [0-1:prototype] [1]结尾字符
						tempDefine.bodyIndex = []int{corvIndex_1, pteIndex_2}
						tempDefine.nameIndexs = dNameIndexs
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

	indexs := filebuf.FindAllSubmatchIndex(REXFunc)
	indexsLen := len(indexs)

	// [0-1:prototype] [2-3:comment or null] [4-5:func type or null]
	// [6-7:func name] [8-9:func params] [10-11:single return value or null]
	// [12-13:multi return value or null] [14-15:"{"]
	if 0 != indexsLen && 16 == len(indexs[0]) {
		for i := 0; i < indexsLen; i++ {
			funcIndexs := indexs[i]

			if isRuleOutIndex(funcIndexs[0], outBetweens) {
				continue
			}

			bodyStartIndex := -1
			bodyEndIndex := funcIndexs[15]

			// 查询"}" 检测是否是有效的函数体
			bodyNewIndexs := filebuf.SubNestIndex(bodyEndIndex-1, SNBraces, outBetweens)

			if 2 != len(bodyNewIndexs) {
				// 无效函数体则跳过
				continue
			} else {
				bodyEndIndex = bodyNewIndexs[1]
			}

			//
			gofunc := goFunc{}

			if -1 != funcIndexs[2] && -1 != funcIndexs[3] {
				gofunc.commentIndex = []int{funcIndexs[2], funcIndexs[3]}
				//	如果注释存在，则body从注释后开始
				bodyStartIndex = funcIndexs[3]
			} else {
				bodyStartIndex = funcIndexs[0]
			}

			gofunc.bodyIndex = []int{bodyStartIndex, bodyEndIndex}
			gofunc.funcNameIndex = []int{funcIndexs[6], funcIndexs[7]}

			// 如果参数为空则下标会相同
			if -1 != funcIndexs[8] && funcIndexs[8] != funcIndexs[9] {
				gofunc.paramIndex = []int{funcIndexs[8], funcIndexs[9]}
			}

			if -1 != funcIndexs[4] && -1 != funcIndexs[5] {
				gofunc.funcTypeIndex = []int{funcIndexs[4], funcIndexs[5]}
			}

			if -1 != funcIndexs[10] && -1 != funcIndexs[11] {
				gofunc.returnIndex = []int{funcIndexs[10], funcIndexs[11]}
			} else if -1 != funcIndexs[12] && -1 != funcIndexs[13] {
				gofunc.returnIndex = []int{funcIndexs[12], funcIndexs[13]}
			}

			result = append(result, gofunc)
		}
	}

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

			if isRuleOutIndex(typeIndexs[0], outBetweens) {
				continue
			}

			gt := goType{}
			bodyStartIndex := -1
			bodyEndIndex := typeIndexs[7]

			if -1 != typeIndexs[2] && -1 != typeIndexs[3] {
				gt.commentIndex = []int{typeIndexs[2], typeIndexs[3]}
				//	如果注释存在，则body从注释后开始
				bodyStartIndex = typeIndexs[3]
			} else {
				bodyStartIndex = typeIndexs[0]
			}

			// 正则确定存在的值
			gt.typeNameIndex = []int{typeIndexs[4], typeIndexs[5]}
			gt.typeIndex = []int{typeIndexs[6], typeIndexs[7]}

			//	判断大括号"{"
			symbolIndex := typeIndexs[8]
			bufByte, _ := filebuf.Byte(symbolIndex)

			if -1 != symbolIndex && '{' == bufByte {
				// 寻找下一个"}"
				bodyNewIndexs := filebuf.SubNestIndex(symbolIndex-1, SNBraces, outBetweens)

				if 2 == len(bodyNewIndexs) {
					bodyEndIndex = bodyNewIndexs[1]
				}
			}

			gt.bodyIndex = []int{bodyStartIndex, bodyEndIndex}

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
	return filebuf.SubNestGetOutBetweens(SNBetweens...)
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

/**
 *	handle comment symbol
 */
func handleComment(src []byte) []byte {
	result := bytes.NewBuffer(nil)
	tempSrc := bytes.Split(src, []byte("\n"))

	if 1 >= len(tempSrc) {
		return src
	}

	cmtStyle1 := 0 // 注释的样式判断, 0 表示是 // 样式, 1表示 /* */ 样式
	if -1 != bytes.Index(tempSrc[0], []byte("/*")) {
		cmtStyle1 = 1
	}

	//	获取注释行前缀的字符下标
	symbolIndex := gosfdoc.FindPrefixFilterTag(tempSrc[cmtStyle1])
	symbolIndexLen := len(symbolIndex)

	for i := cmtStyle1; i < len(tempSrc)-cmtStyle1; i++ {
		b := tempSrc[i]
		if len(b) > symbolIndexLen {
			result.Write(b[symbolIndexLen:])
			result.WriteByte('\n')
		}
	}

	return result.Bytes()
}
