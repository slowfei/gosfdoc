//  The MIT License (MIT) - http://opensource.org/licenses/MIT
//
//  Copyright (c) 2014 slowfei
//
//  Create on 2014-08-16
//  Update on 2014-09-19
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
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	APPNAME = "gosfdoc"   //
	VERSION = "0.0.1.000" //

	DIR_NAME_MAIN_MARKDOWN    = "md"      // save markdown file main directory name
	DIR_NAME_MARKDOWN_DEFAULT = "default" // markdown default directory
	DIR_NAME_SOURCE_CODE      = "src"     // source code save directory

	FILE_NAME_ABOUT_MD     = "about.md"
	FILE_NAME_INTRO_MD     = "intro.md"
	FILE_NAME_CONTENT_JSON = "content.json"
)

var (
	//  document parser implement interface
	_mapParser = make(map[string]DocParser)
	//  system filters
	_sysFilters = []string{DEFAULT_CONFIG_FILE_NAME, "."}

	//  error info
	ErrConfigNotRead      = errors.New("Can not read config file.")
	ErrSpecifyCodeLangNil = errors.New("Specify code language nil.")
	ErrDirNotExist        = errors.New("Specified directory does not exist.")
	ErrDirIsFilePath      = errors.New("This is a file path.")
	ErrFilePathOccupied   = errors.New("(gosfdoc.json) Config file path has been occupied.")
	// ErrPathInvalid        = errors.New("invalid operate path.")
)

/**
 *  regex compile variable
 */
var (
	// private file tag ( //#private-doc-code )
	REXPrivateFile = regexp.MustCompile("#private-(doc|code){1}(-doc|-code)?")
	TagPrivateCode = []byte("code")
	TagPrivateDoc  = []byte("doc")
	// private block tag ( //#private * //#private-end)
	REXPrivateBlock = regexp.MustCompile("[^\\n]?//#private(\\s|.)*?//#private-end[\\s]?")

	// parse about and intro block
	/**[About|Intro]
	 *  content text or markdown text
	 */
	//[About|Intro]
	// content text or markdown text
	//End
	REXAbout = regexp.MustCompile("(/\\*\\*About[\\s]+(\\s|.)*?[\\s]+\\*/)|(//About[\\s]?([\\s]|.)*?//[Ee][Nn][Dd])")
	REXIntro = regexp.MustCompile("(/\\*\\*Intro[\\s]+(\\s|.)*?[\\s]+\\*/)|(//Intro[\\s]?([\\s]|.)*?//[Ee][Nn][Dd])")

	// parse public document content
	/***[z-index-][title]
	 *  document text or markdown text
	 */
	///[z-index-][title]
	//  document text or markdown text
	//End
	REXDocument      = regexp.MustCompile("(/\\*\\*\\*[^\\*\\s](.+)\\n(\\s|.)*?\\*/)|(///[^/\\s](.+)\\n(\\s|.)*?//[Ee][Nn][Dd])")
	REXDocIndexTitle = regexp.MustCompile("(/\\*\\*\\*|///)(\\d*-)?(.*)?")
)

/**
 *  operate result
 */
type OperateResult int

const (
	ResultFileSuccess OperateResult = iota
	ResultFileInvalid
	ResultFileNotRead
	ResultFileReadErr
	ResultFileFilter
)

/**
 *  file scan result func
 *
 *  @param `path`
 *  @param `result`
 */
type FileResultFunc func(path string, result OperateResult)

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
	 *  @param `parh` file path
	 *  @param `info` file info
	 *  @return true is valid file
	 */
	CheckFile(path string, info os.FileInfo) bool

	/**
	 *  each file the content
	 *  can be create keyword index and other operations
	 *
	 *  @param `filebuf`    file content buf
	 */
	EachIndexFile(filebuf *FileBuf)

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
	AddParser(new(nilDocParser))
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
func readConfigFile(filepath string) (config *MainConfig, err error, pass bool) {
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
	config = mainConfig

	return
}

/**
 *  get source code directory svae path
 *
 *  @param `config`
 *  @return full path
 */
func dirpathSourceCode(config *MainConfig) string {
	return filepath.Join(config.Outdir, DIR_NAME_SOURCE_CODE)
}

/**
 *  get default markdown directory save path
 *
 *  @param `config`
 *  @return full path
 */
func dirpathMarkdownDefault(config *MainConfig) string {
	return filepath.Join(config.Outdir, DIR_NAME_MAIN_MARKDOWN, DIR_NAME_MARKDOWN_DEFAULT)
}

/**
 *  create config file
 *
 *  @param `dirPath` directory path
 *  @param `langs`   specify code language, nil is all language, value is parser name.
 *  @return `error`  warn or error message
 *  @return `bool`   true is operation success
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
		defaultConfigText := fmt.Sprintf(_gosfdocConfigJson, SFFileManager.GetCmdDir(), codeLangs)

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
 *  @return `error` warn or error message
 *  @return `bool`  true is operation success
 */
func Output(configPath string, fileFunc FileResultFunc) (error, bool) {
	config, err, pass := readConfigFile(configPath)
	if !pass {
		return err, pass
	}
	return OutputWithConfig(config, fileFunc)
}

/**
 *  build output document with config content
 *
 *  @param `config`
 *  @return `error` warn or error message
 *  @return `bool`  true is operation success
 */
func OutputWithConfig(config *MainConfig, fileFunc FileResultFunc) (error, bool) {
	err, pass := config.Check()
	if !pass {
		return err, pass
	}
	scanPath := config.ScanPath

	isExists, isDir, _ := SFFileManager.Exists(scanPath)
	if !isExists || !isDir {
		return errors.New(fmt.Sprintf("invalid scan path path: %v", scanPath)), false
	}

	files, about, intro, scanErr := scanFiles(config, fileFunc)
	if nil != scanErr {
		return scanErr, false
	}

	// output markdown defualt directory path
	mdDefaultPath := dirpathMarkdownDefault(config)

	//  output content.json
	contentPath := filepath.Join(mdDefaultPath, FILE_NAME_CONTENT_JSON)
	contentStruct := ContentJson{
		HtmlTitle: config.HtmlTitle,
		DocTitle:  config.DocTitle,
		MenuTitle: config.MenuTitle,
	}
	contentStruct.WriteFilepath(contentPath)

	//  output about.md and intro.md
	about.WriteFilepath(filepath.Join(mdDefaultPath, FILE_NAME_ABOUT_MD))
	intro.WriteFilepath(filepath.Join(mdDefaultPath, FILE_NAME_INTRO_MD))

	//  TODO 等待测试
	if nil != files {

	}

	return nil, true
}

/**
 *  scan files
 *
 *  @param `config`
 *  @param `fileFunc`
 *  @return `resultFiles` map[string]*CodeFiles
 *  @return `aboutBuf`
 *  @return `introBuf`
 *  @return `resultErr`
 */
func scanFiles(config *MainConfig, fileFunc FileResultFunc) (
	resultFiles map[string]*CodeFiles,
	about *About,
	intro *Intro,
	resultErr error) {

	resultFiles = make(map[string]*CodeFiles)

	callFileFunc := func(p string, r OperateResult) error {
		if nil != fileFunc {
			fileFunc(p, r)
		}
		return nil
	}

	resultErr = filepath.Walk(config.ScanPath, func(path string, info os.FileInfo, err error) error {

		if nil != err || nil == info {
			return callFileFunc(path, ResultFileNotRead)
		}

		fileName := info.Name()

		// 1. system file filter
		for i := 0; i < len(_sysFilters); i++ {
			sysFileName := _sysFilters[i]
			if 0 == strings.Index(fileName, sysFileName) {
				return callFileFunc(path, ResultFileFilter)
			}
		}

		// 2. filter custom path
		for i := 0; i < len(config.FilterPaths); i++ {
			fpath := config.FilterPaths[i]
			if 0 == strings.Index(path, fpath) {
				return callFileFunc(path, ResultFileFilter)
			}
		}

		// filter document output dir
		if 0 == strings.Index(path, config.Outdir) {
			return callFileFunc(path, ResultFileFilter)
		}

		// 目录检测
		if info.IsDir() {
			if _, ok := resultFiles[path]; !ok {
				resultFiles[path] = NewCodeFiles()
			}
			return nil
		}

		//  无法找到后缀视为无效文件
		if 0 >= strings.LastIndex(".", fileName) {
			return callFileFunc(path, ResultFileInvalid)
		}

		// 3. check file and find parser
		var parser DocParser = nil
		for _, vp := range _mapParser {
			if vp.CheckFile(path, info) {
				parser = vp
				break
			}
		}
		if nil == parser {
			return callFileFunc(path, ResultFileInvalid)
		}

		file, openErr := os.Open(path)
		if openErr != nil {
			if nil != fileFunc {
				fileFunc(path, ResultFileNotRead)
			}
			return nil
		}
		defer file.Close()

		//  在特定的字节数查询换行符号，如果未查询到换行符就判定为无效的文件
		firstLineBuf := make([]byte, 4096*2)
		rn, readErr := file.Read(firstLineBuf)

		if -1 >= rn || nil != readErr {
			return callFileFunc(path, ResultFileReadErr)
		}

		firstLine := firstLineBuf[:rn]
		rnIndex := bytes.IndexByte(firstLine, '\n')
		if -1 == rnIndex {
			return callFileFunc(path, ResultFileInvalid)
		}

		// 4. check file private tag //#private-doc //#private-code //#private-doc-code
		privateTag := REXPrivateFile.Find(firstLine)
		isPCode := false
		isPDoc := false
		if nil != privateTag && 0 != len(privateTag) {
			if 0 < bytes.Index(privateTag, TagPrivateCode) {
				isPCode = true
			}
			if 0 < bytes.Index(privateTag, TagPrivateDoc) {
				isPDoc = true
			}
		}
		if isPCode && isPDoc {
			return callFileFunc(path, ResultFileFilter)
		}

		//  handle file content bytes
		fileBytes, rFileErr := readFile(file, info.Size())
		if nil != rFileErr {
			return callFileFunc(path, ResultFileReadErr)
		}

		// 5. filter private block and create file buffer
		fileBuf := NewFileBuf(fileBytes, path, info, REXPrivateBlock)

		// 6. parse about and intro
		if nil == about {
			about = ParseAbout(fileBuf)
		}
		if nil == intro {
			intro = ParseIntro(fileBuf)
		}

		// 7. parse documents
		var documents []Document = nil
		if !isPDoc {
			documents = ParseDocument(fileBuf)
		}

		// 8. create file index
		parser.EachIndexFile(fileBuf)

		//  pack CodeFile
		var files *CodeFiles = nil
		var ok bool = false
		pathDir := filepath.Dir(path)

		if files, ok = resultFiles[pathDir]; !ok {
			files = NewCodeFiles()
			resultFiles[pathDir] = files
		}

		codeFile := CodeFile{}
		codeFile.FileCont = fileBuf
		codeFile.PrivateCode = isPCode
		codeFile.PrivateDoc = isPDoc
		codeFile.parser = parser
		codeFile.docs = documents

		files.addFile(codeFile)

		return nil
	}) // end Walk file

	return
}

/**
 *  read file bytes
 *
 *  @param `r`
 *  @param `fileSize`
 */
func readFile(r io.Reader, fileSize int64) (b []byte, err error) {
	var capacity int64

	if fileSize < 1e9 {
		capacity = fileSize
	}

	buf := bytes.NewBuffer(make([]byte, 0, capacity+bytes.MinRead))
	defer func() {
		e := recover()
		if e == nil {
			return
		}
		if panicErr, ok := e.(error); ok && panicErr == bytes.ErrTooLarge {
			err = panicErr
		} else {
			panic(e)
		}
	}()
	_, err = buf.ReadFrom(r)
	return buf.Bytes(), err
}
