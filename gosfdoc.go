//  The MIT License (MIT) - http://opensource.org/licenses/MIT
//
//  Copyright (c) 2014 slowfei
//
//  Create on 2014-08-16
//  Update on 2014-08-22
//  Email  slowfei#foxmail.com
//  Home   http://www.slowfei.com

//
package gosfdoc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/slowfei/gosfcore/utils/filemanager"
	"io/ioutil"
	"path/filepath"
	// "strings"
)

const (
	APPNAME = "gosfdoc"
	VERSION = "0.0.1.000"
)

var (
	//  document parser implement interface
	_mapParser = make(map[string]DocParser)

	//  error info
	ErrConfigNotRead      = errors.New("Can not read config file.")
	ErrSpecifyCodeLangNil = errors.New("Specify code language nil.")
	ErrDirNotExist        = errors.New("Specified directory does not exist.")
	ErrDirIsFilePath      = errors.New("This is a file path.")
	ErrFilePathOccupied   = errors.New("(gosfdoc.json) Config file path has been occupied.")
)

/**
 *  document parser
 *
 */
type DocParser interface {

	/**
	 *  parser name
	 *
	 *  @return
	 */
	Name() string

	/**
	 *  check file
	 *  detecting whether the file is a valid file
	 *
	 *  @return true is valid file
	 */
	CheckFilepath() bool

	/**
	 *  each file the content
	 *  can be create keyword index and other operations
	 *
	 *  @param `index` while file index
	 *  @param `fileCont` file content
	 */
	EachFile(index int, fileCont *bytes.Buffer)

	/**
	 *  parse file document tag
	 *
	 *  @param `fileCont` file content
	 *  @return slice
	 */
	ParseDoc(fileCont *bytes.Buffer) []Document

	/**
	 *  parse file preview tag
	 *
	 *  @param `fileCont` file content
	 *  @return slice
	 */
	ParsePreview(fileCont *bytes.Buffer) []Preview

	/**
	 *  parse code block tag
	 *
	 *  @param `fileCont` file content
	 *  @return slice
	 */
	ParseCodeblock(fileCont *bytes.Buffer) []CodeBlock
}

/**
 *  init
 */
func init() {

}

/**
 *  add parser
 *
 *  @param parser
 */
func AddParser(parser DocParser) {
	if nil != parser {
		_mapParser[parser.Name()] = parser
	}
}

/**
 *  get parsers
 *  key is parser name
 *  value is parser implement
 *
 *  @return
 */
func MapParser() map[string]DocParser {
	return _mapParser
}

/**
 *  read config file
 *
 *  @param `filepath`
 *  @return `config`
 *  @return `err`   contains warn info
 *  @return `pass`  true is valid file (pass does not mean that there are no errors)
 */
func readConfigFile(filepath string) (config MainConfig, err error, pass bool) {
	result := false

	isExists, isDir, _ := SFFileManager.Exists(filepath)
	if !isExists || isDir {
		err = ErrConfigNotRead
		pass = result
		return
	}

	jsonData, readErr := ioutil.ReadFile(filepath)
	if nil != readErr {
		err = ErrConfigNotRead
		pass = result
		return
	}

	mainConfig := new(MainConfig)
	json.Unmarshal(jsonData, mainConfig)

	err, pass = mainConfig.Check()
	config = *mainConfig

	return
}

/**
 *  create config file
 *
 *  @param `dirPath` directory path
 *  @param `langs`   specify code language, nil is all language, value is parser name.
 *  @return `error`  warn or error info
 *  @return `bool`   true is success create file
 */
func CreateConfigFile(dirPath string, langs []string) (error, bool) {
	if nil == langs || 0 == len(langs) {
		return ErrSpecifyCodeLangNil, false
	}
	isCreateFile := true
	errBuf := bytes.NewBufferString("")

	//  检测目录操作
	isExists, isDir, _ := SFFileManager.Exists(dirPath)
	if !isExists {
		return ErrDirNotExist, false
	}
	if !isDir {
		return ErrDirIsFilePath, false
	}

	//  检测配置文件操作
	filePath := filepath.Join(dirPath, DEFAULT_CONFIG_FILE_NAME)
	isExists, isDir, _ = SFFileManager.Exists(filePath)

	if !isExists {
		//  配置文件不存在，直接创建配置文件

		codeLangs := ""
		langCount := len(langs)

		for i := 0; i < langCount; i++ {
			lang := langs[i]
			if _, ok := _mapParser[lang]; !ok {
				errBuf.WriteString("Language: not " + lang + " Parser.\n")
			} else {
				codeLangs += "\"" + lang + "\","
			}
		}

		//  如果相等表示没有全部没有找到语言的解析器则直接返回
		if 0 == len(codeLangs) {
			return errors.New(errBuf.String()), false
		}

		if ',' == codeLangs[len(codeLangs)-1] {
			codeLangs = codeLangs[:len(codeLangs)-1]
		}

		// 将指定的语言保存进默认配置信息中。
		defaultConfigText := fmt.Sprintf(_gosfdocConfigJson, codeLangs)

		fileErr := ioutil.WriteFile(filePath, []byte(defaultConfigText), 0660)
		if nil != fileErr {
			isCreateFile = false
			errBuf.WriteString(fileErr.Error())
		}

	} else {
		if isDir {
			return ErrFilePathOccupied, false
		}

		_, err, _ := readConfigFile(filePath)
		if nil != err {
			isCreateFile = false
			errBuf.WriteString(err.Error())
		}

	}

	var resErr error = nil
	if 0 != errBuf.Len() {
		resErr = errors.New(errBuf.String())
	}

	return resErr, isCreateFile
}

/**
 *  build output document
 *
 *  @param `configPath` config file path
 */
func Output(configPath string) {

}

/**
 *  build output document with config content
 *
 *  @param `config`
 */
func OutputWithConfig(config []byte) {

}
