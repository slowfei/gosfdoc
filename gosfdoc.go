//  The MIT License (MIT) - http://opensource.org/licenses/MIT
//
//  Copyright (c) 2014 slowfei
//
//  Create on 2014-08-16
//  Update on 2014-10-29
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
	"github.com/slowfei/gosfdoc/assets"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const (
	APPNAME = "gosfdoc"   //
	VERSION = "0.0.1.000" //

	DIR_NAME_MAIN_MARKDOWN    = "md"      // save markdown file main directory name
	DIR_NAME_MARKDOWN_DEFAULT = "default" // markdown default directory
	DIR_NAME_SOURCE_CODE      = "src"     // source code save directory
	DIR_NAME_ASSETS           = "assets"  // html use assets file directory

	FILE_SUFFIX_MARKDOWN = ".md"

	FILE_NAME_ABOUT_MD     = "about.md"
	FILE_NAME_INTRO_MD     = "intro.md"
	FILE_NAME_CONTENT_JSON = "content.json"

	FILE_NAME_GOSFDOC_MIN_CSS    = "gosfdoc.min.css"
	FILE_NAME_ASSETS_MIN_JS      = "assets.min.js"
	FILE_NAME_GOSFDOC_MIN_JS     = "gosfdoc.min.js"
	FILE_NAME_GOSFDOC_SRC_MIN_JS = "gosfdoc.src.min.js"

	FILE_NAME_HTML_INDEX = "index.html"
	FILE_NAME_HTML_SRC   = "src.html"
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
	ResultFileOutFail
	ResultDebugErr
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
	 *  @param `filebuf`    file content buffer
	 */
	EachIndexFile(filebuf *FileBuf)

	/**
	 *  parse file preview tag
	 *
	 *  @param `filebuf` file content buffer
	 *  @return slice
	 */
	ParsePreview(filebuf *FileBuf) []Preview

	/**
	 *  parse code block tag
	 *
	 *  @param `filebuf` file content buffer
	 *  @return slice
	 */
	ParseCodeblock(filebuf *FileBuf) []CodeBlock

	/**
	 *	parse directory package info
	 *	each file directory parse string join
	 *
	 *	@param `filebuf`
	 *	@return string file parse the only string
	 */
	ParsePackageInfo(filebuf *FileBuf) string
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
 *	create directory path
 *
 *	@param `path`
 */
func dirpathMkall(path string) {
	exists, isDir, err := SFFileManager.Exists(path)
	if !exists {
		err = os.MkdirAll(path, os.ModePerm)
		if nil != err {
			panic(fmt.Sprintln(path, err.Error()))
		}
	} else if !isDir {
		panic(fmt.Sprintln(path, "has been occupied."))
	} else if nil != err {
		panic(fmt.Sprintln(path, err.Error()))
	}
}

/**
 *  get default markdown directory save path
 *
 *  @param `config`
 *  @return full path
 */
func dirpathMarkdownDefault(config *MainConfig) string {
	path := filepath.Join(config.Outpath, DIR_NAME_MAIN_MARKDOWN, DIR_NAME_MARKDOWN_DEFAULT)

	dirpathMkall(path)

	return path
}

/**
 *	get assets directory save path
 */
func dirpathAssets(config *MainConfig) string {
	path := filepath.Join(config.Outpath, DIR_NAME_ASSETS)

	dirpathMkall(path)

	return path
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

		cmdDir := SFFileManager.GetCmdDir()
		appendPath := filepath.Base(cmdDir)

		//	考虑到如果没有项目名称，base获取的是src追加的路径则为空
		//	$GOPATH/src
		//	$GOPATH/src/projectname
		if "src" == appendPath {
			appendPath = ""
		}

		// 将指定的语言保存进默认配置信息中。
		// 默认初始值：
		//	ScanPath = command directory
		//	CodeLang = implement parser the code language
		//	OutAppendPath = command directory base name
		defaultConfigText := fmt.Sprintf(_gosfdocConfigJson, cmdDir, codeLangs, filepath.Base(cmdDir))

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
	config.setAbspath()
	scanPath := config.ScanPath

	isExists, isDir, _ := SFFileManager.Exists(scanPath)
	if !isExists || !isDir {
		return errors.New(fmt.Sprintf("invalid scan path path: %v", scanPath)), false
	}

	files, about, intro, scanErr := scanFiles(config, fileFunc)
	if nil != scanErr {
		return scanErr, false
	}

	// markdown defualt directory path
	mdDefaultPath := dirpathMarkdownDefault(config)

	//  output content.json
	contentPath := filepath.Join(mdDefaultPath, FILE_NAME_CONTENT_JSON)
	contentStruct := ContentJson{
		HtmlTitle: config.HtmlTitle,
		DocTitle:  config.DocTitle,
		MenuTitle: config.MenuTitle,
	}
	contentErr := contentStruct.WriteFilepath(contentPath)
	if nil != contentErr {
		fileFunc(contentPath, ResultFileOutFail)
	}

	//  output about.md and intro.md
	if nil == about {
		about = NewDefaultAbout()
	}
	if nil == intro {
		intro = NewDefaultIntro()
	}
	aboutPath := filepath.Join(mdDefaultPath, FILE_NAME_ABOUT_MD)
	aboutErr := about.WriteFilepath(aboutPath)
	if nil != aboutErr {
		fileFunc(aboutPath, ResultFileOutFail)
	}
	introPath := filepath.Join(mdDefaultPath, FILE_NAME_INTRO_MD)
	introErr := intro.WriteFilepath(introPath)
	if nil != introErr {
		fileFunc(introPath, ResultFileOutFail)
	}

	// output source code file and markdown document file
	outCodeFiles(config, files, fileFunc)

	// output html assets file
	outAssets(config, fileFunc)

	return nil, true
}

/**
 *	map[string]*CodeFiles
 *	source code and markdown output
 *
 *	@param `config`
 *	@param `files`
 *	@param `fileFunc`
 *	@return []PackageInfo
 */
func outCodeFiles(config *MainConfig, files map[string]*CodeFiles, fileFunc FileResultFunc) []PackageInfo {

	scanPath := config.ScanPath
	appendPath := config.OutAppendPath

	//	source code ouput path operation
	//	projectroot/doc/src/[appendpath/main.go]
	isLinkRoot := config.CodeLinkRoot
	outCodeDir := ""
	if config.CopyCode {
		outCodeDir = filepath.Join(config.Outpath, DIR_NAME_SOURCE_CODE)
	}

	//	markdown file save directory
	mdDir := dirpathMarkdownDefault(config)
	packInfos := make([]PackageInfo, 0, len(files))

	// 1.FOR Directory
	for dirPath, codefiles := range files {
		fileLen := codefiles.files.Len()
		relativeDirPath := ""

		if 0 == strings.Index(dirPath, scanPath) {
			relativeDirPath = dirPath[len(scanPath):]
		} else {
			if nil != fileFunc {
				fileFunc(dirPath, ResultDebugErr)
			}
			fmt.Println("map CodeFiles save path error.")
			fmt.Println("ScanPath:", scanPath)
			fmt.Println("CodeFiles Dirpath:", dirPath)
			continue
		}

		//	append custom path prefix
		relativeDirPath = filepath.Join(appendPath, relativeDirPath)

		previews := make([]Preview, 0, 0)
		blocks := make([]CodeBlock, 0, 0)
		documents := make([]Document, 0, 0)
		filesName := make([]string, 0, fileLen)
		packStrList := make([]string, 0, fileLen)

		// 2.FOR Files
		codefiles.Each(func(code CodeFile) bool {
			// 3. source code check
			if !code.PrivateCode {
				switch code.parser.(type) {
				case *nilDocParser:
				default:
					var outErr error = nil
					fileName := code.FileCont.FileInfo().Name()

					if 0 != len(outCodeDir) {
						outPath := filepath.Join(outCodeDir, relativeDirPath, fileName)
						outErr = code.FileCont.WriteFilepath(outPath)
						if nil == outErr && nil != fileFunc {
							fileFunc(outPath, ResultFileSuccess)
						}
					}

					if nil == outErr {
						if 0 != len(outCodeDir) || isLinkRoot {
							filesName = append(filesName, fileName)
						}
					} else if nil != fileFunc {
						fileFunc(code.FileCont.path, ResultFileOutFail)
					}
				}
			}

			// 4. parse Preview and CodeBlock and Document
			if !code.PrivateDoc {
				ps := code.parser.ParsePreview(code.FileCont)
				bs := code.parser.ParseCodeblock(code.FileCont)

				if 0 != len(ps) {
					previews = append(previews, ps...)
				}

				if 0 != len(bs) {
					blocks = append(blocks, bs...)
				}

				if 0 != len(code.docs) {
					documents = append(documents, code.docs...)
				}
			}

			// 5. parse package info
			packInfo := code.parser.ParsePackageInfo(code.FileCont)
			if 0 != len(packInfo) {
				packStrList = append(packStrList, packInfo)
			}

			return true
		})

		//
		sort.Sort(SortSet{previews: previews})
		sort.Sort(SortSet{codeBlocks: blocks})
		sort.Sort(SortSet{documents: documents})

		// 5.output markdown
		mdBytes := ParseMarkdown(documents, previews, blocks, filesName, relativeDirPath)
		if 0 != len(mdBytes) {
			//	markdown file name is directory base name + suffix
			mdFileName := filepath.Base(dirPath) + FILE_SUFFIX_MARKDOWN
			outPath := filepath.Join(mdDir, relativeDirPath, mdFileName)

			err := SFFileManager.WirteFilepath(outPath, mdBytes)
			result := ResultFileSuccess

			if nil != err {
				result = ResultFileOutFail
			} else {
				info := PackageInfo{}
				info.Name = path.Join(relativeDirPath, mdFileName)

				joinStr := strings.Join(packStrList, ";")
				newStr := strings.Replace(joinStr, "\n", ", ", -1)
				info.Desc = newStr

				packInfos = append(packInfos, info)
			}

			if nil != fileFunc {
				fileFunc(outPath, result)
			}
		}

	} // end for dirPath, codefiles := range files

	return packInfos
}

/**
 *	output assets file
 *
 *	@param `config`
 *	@param `fileFunc`
 */
func outAssets(config *MainConfig, fileFunc FileResultFunc) {

	dirpath := dirpathAssets(config)

	assetsPath := filepath.Join(dirpath, FILE_NAME_ASSETS_MIN_JS)
	err := SFFileManager.WirteFilepath(assetsPath, []byte(assets.ASSETS_MIN_JS))

	if nil != err {
		fileFunc(assetsPath, ResultFileOutFail)
	}

	gosfdocPath := filepath.Join(dirpath, FILE_NAME_GOSFDOC_MIN_JS)
	err = SFFileManager.WirteFilepath(gosfdocPath, []byte(assets.GOSFDOC_MIN_JS))

	if nil != err {
		fileFunc(gosfdocPath, ResultFileOutFail)
	}

	gosfdocsrcPath := filepath.Join(dirpath, FILE_NAME_GOSFDOC_SRC_MIN_JS)

	if config.CopyCode {
		err = SFFileManager.WirteFilepath(gosfdocsrcPath, []byte(assets.GOSFDOC_SRC_MIN_JS))
	} else {
		err = SFFileManager.WirteFilepath(gosfdocsrcPath, []byte(assets.GOSFDOC_SRC_MIN_JS_ROOT))
	}

	if nil != err {
		fileFunc(gosfdocsrcPath, ResultFileOutFail)
	}

	gosfdoccssPath := filepath.Join(dirpath, FILE_NAME_GOSFDOC_MIN_CSS)
	err = SFFileManager.WirteFilepath(gosfdoccssPath, []byte(assets.GOSFDOC_MIN_CSS))

	if nil != err {
		fileFunc(gosfdoccssPath, ResultFileOutFail)
	}

}

/**
 *	output html file, index.html src.html
 *
 *	@param `config`
 *	@param `fileFunc`
 */
func outHTML(config *MainConfig, fileFunc FileResultFunc) {

	indexPath := filepath.Join(config.Outpath, FILE_NAME_HTML_INDEX)
	err := SFFileManager.WirteFilepath(indexPath, []byte(assets.HTML_INDEX))

	if nil != err {
		fileFunc(indexPath, ResultFileOutFail)
	}

	srcPath := filepath.Join(config.Outpath, FILE_NAME_HTML_SRC)
	err = SFFileManager.WirteFilepath(srcPath, []byte(assets.HTML_SRC))

	if nil != err {
		fileFunc(srcPath, ResultFileOutFail)
	}

}

/**
 *	output html used config.json
 *
 *	@param `config`
 *	@param `fileFunc`
 */
func outHTMLConfig(config *MainConfig, fileFunc FileResultFunc, packageInfos []PackageInfo) {
	// type DocConfig struct {
	// ContentJson string                   // content json file
	// IntroMd     string                   // intro markdown file
	// AboutMd     string                   // about markdown file
	// Languages   map[string]string        // key is directory name, value is show text
	// Markdowns   map[string][]PackageInfo // markdown info list
	// Files       map[string][]string
	// }

	docConfig := DocConfig{}
	docConfig.ContentJson = FILE_NAME_CONTENT_JSON
	docConfig.IntroMd = FILE_NAME_INTRO_MD
	docConfig.AboutMd = FILE_NAME_ABOUT_MD
	docConfig.Languages = config.Languages

	//	TODO 目前先需要确定页面后才进行Markdowns和Files的输出问题。

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
		if 0 == strings.Index(path, config.Outpath) {
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
		if 0 >= strings.LastIndex(fileName, ".") {
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
		firstLineBuf := make([]byte, 1024) //4096
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
			// remove private tag
			firstLine = bytes.Replace(firstLine, privateTag, nil, 1)
		}
		if isPCode && isPDoc {
			return callFileFunc(path, ResultFileFilter)
		}

		//  handle file content bytes
		fileBytes, rFileErr := readFile(firstLine, file, info.Size())
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
 *  @param `beforeReadData` append before read data
 *  @param `r`
 *  @param `fileSize`
 */
func readFile(beforeReadData []byte, r io.Reader, fileSize int64) (b []byte, err error) {
	var capacity int64

	if fileSize < 1e9 {
		capacity = fileSize
	}

	buf := bytes.NewBuffer(make([]byte, 0, capacity+bytes.MinRead))
	buf.Write(beforeReadData)
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
