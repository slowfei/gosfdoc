package gosfdoc

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/slowfei/gosfcore/utils/filemanager"
	"io/ioutil"
	"path/filepath"
)

const (
	DEFAULT_CONFIG_FILE_NAME = "gosfdoc.json"
)

var (
	_gosfdocConfigJson = `{
    "Path"             : %#v,
    "CodeLang"         : [%v],
    "Outdir"           : "doc",
    "CopyCode"         : false,
    "HtmlTitle"        : "Document",
    "DocTitle"         : "<b>Document:</b>",
    "MenuTitle"        : "<center><b>package</b></center>",
    "Languages"        : {"default" : "Default"},
    "FilterPaths"      : []
}`
)

/**
 *  main config info
 *  output `gosfdoc.json` use
 */
type MainConfig struct {
	Path        string            // operate absolute path
	CodeLang    []string          // code languages
	Outdir      string            // output document directory
	CopyCode    bool              // whether output source code to document directory
	HtmlTitle   string            // document html show title
	DocTitle    string            // html top tabbar show title
	MenuTitle   string            // html left menu show title
	Languages   map[string]string // document support the language. key is lang dirctory name, value is show text
	FilterPaths []string          // filter directory relative path
}

/**
 *  check config param value
 *  error value will update default.
 *
 *  @return error
 *  @return bool    fatal error is false, pass is true. (pass does not mean that there are no errors)
 */
func (mc *MainConfig) Check() (error, bool) {
	errBuf := bytes.NewBufferString("")
	pass := true

	if 0 == len(mc.Path) {
		errBuf.WriteString("Path: please set the absolute path need to operate.\n")
		pass = false
	}

	if 0 == len(mc.CodeLang) {
		errBuf.WriteString("CodeLang: specify code language type nil.\n")
		pass = false
	} else {
		count := len(mc.CodeLang)
		for i := 0; i < count; i++ {
			lang := mc.CodeLang[i]
			if _, ok := _mapParser[lang]; !ok {
				errBuf.WriteString("CodeLang: not " + lang + " Parser.\n")
			}
		}
	}

	if 0 == len(mc.Outdir) {
		errBuf.WriteString("Outdir: output directory is nil, will use 'doc' default directory.\n")
		mc.Outdir = "doc"
	}

	if 0 == len(mc.HtmlTitle) {
		mc.HtmlTitle = "Document"
		errBuf.WriteString("HtmlTitle: to set the html title.\n")
	}

	if 0 == len(mc.DocTitle) {
		mc.DocTitle = "<b>Document:</b>"
		errBuf.WriteString("DocTitle: to set the doc title.\n")
	}

	if 0 == len(mc.MenuTitle) {
		mc.MenuTitle = "<center><b>package</b></center>"
		errBuf.WriteString("MenuTitle: to set the menu title.\n")
	}

	if 0 == len(mc.Languages) {
		mc.Languages = map[string]string{"Default": "default"}
		errBuf.WriteString("Languages: to set the default html text language.\n")
	} else {
		if _, ok := mc.Languages["default"]; !ok {
			mc.Languages["default"] = "Default"
			errBuf.WriteString("Languages: to set the default html text language.\n")
		}
	}

	var err error = nil
	if 0 != errBuf.Len() {
		err = errors.New(errBuf.String())
	}
	return err, pass
}

/**
 *  document directory html javascript use config
 *
 *  output `config.json`
 */
type DocConfig struct {
	ContentJson string                   // content json file
	IntroMd     string                   // intro markdown file
	AboutMd     string                   // about markdown file
	Languages   map[string]string        // key is directory name, value is show text
	Markdowns   map[string][]PackageInfo // markdown info list
}

/**
 *  conifg load
 *
 *  @param jsonData
 *  @return load error info
 */
func configLoadByJson(jsonData []byte, c *MainConfig) error {
	e2 := json.Unmarshal(jsonData, c)
	if nil != e2 {
		return e2
	}
	return nil
}

/**
 *  conifg load
 *
 *  @param configPath
 *  @return error info
 */
func configLoadByFilepath(configPath string, c *MainConfig) error {

	var path string
	if filepath.IsAbs(configPath) {
		path = configPath
	} else {
		path = filepath.Join(SFFileManager.GetExecDir(), configPath)
	}

	isExists, isDir, _ := SFFileManager.Exists(path)
	if !isExists || isDir {
		return errors.New("failed to load configuration file:" + configPath)
	}

	jsonData, e1 := ioutil.ReadFile(path)
	if nil != e1 {
		return e1
	}

	return configLoadByJson(jsonData, c)
}
