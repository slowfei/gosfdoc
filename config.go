//  The MIT License (MIT) - http://opensource.org/licenses/MIT
//
//  Copyright (c) 2014 slowfei
//
//  Create on 2014-08-16
//  Update on 2014-10-10
//  Email  slowfei#foxmail.com
//  Home   http://www.slowfei.com

//
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
	DEFAULT_OUTPATH          = "doc"
)

var (
	_gosfdocConfigJson = `{
    "ScanPath"         : %#v,
    "CodeLang"         : [%v],
    "Outpath"          : "doc",
    "OutAppendPath"    : %#v,
    "CopyCode"         : false,
    "CodeLinkRoot"     : true,
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
	path          string            // private handle path, save console command path.
	ScanPath      string            // scan document info file path, relative or absolute path, is "/" scan current console path.
	CodeLang      []string          // code languages
	Outpath       string            // output document path, relative or absolute path.
	OutAppendPath string            // append output source code and markdown relative path(scan path join). defalut ""
	CopyCode      bool              // copy source code to document directory. default false
	CodeLinkRoot  bool              // source code link to root directory, 'CopyCode' is true was invalid default true
	HtmlTitle     string            // document html show title
	DocTitle      string            // html top tabbar show title
	MenuTitle     string            // html left menu show title
	Languages     map[string]string // document support the language. key is directory name, value is show text.
	FilterPaths   []string          // filter path, relative or absolute path
}

/**
 *	set absolute path
 */
func (mc *MainConfig) setAbspath() {

	mc.path = SFFileManager.GetCmdDir()

	if 0 == len(mc.ScanPath) || "/" == mc.ScanPath {
		mc.ScanPath = mc.path
	} else if !filepath.IsAbs(mc.ScanPath) {
		mc.ScanPath = filepath.Join(mc.path, mc.ScanPath)
	}

	if 0 == len(mc.Outpath) {
		mc.Outpath = filepath.Join(mc.ScanPath, DEFAULT_OUTPATH)
	} else if !filepath.IsAbs(mc.Outpath) {
		mc.Outpath = filepath.Join(mc.ScanPath, mc.Outpath)
	}

	for i := 0; i < len(mc.FilterPaths); i++ {
		p := mc.FilterPaths[i]
		if !filepath.IsAbs(p) {
			mc.FilterPaths[i] = filepath.Join(mc.ScanPath, p)
		}
	}
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

	mc.path = SFFileManager.GetCmdDir()

	if 0 == len(mc.ScanPath) {
		errBuf.WriteString("ScanPath: please set document scan path.\n")
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

	if 0 != len(mc.OutAppendPath) && filepath.IsAbs(mc.OutAppendPath) {
		errBuf.WriteString("OutAppendPath: please use relative path.\n")

		//	这里主要怕效验不通过后强制执行，所以强行修改默认的OutAppendPath
		tempPath := filepath.Base(mc.path)
		if "src" == tempPath {
			tempPath = ""
		}
		mc.OutAppendPath = tempPath

		pass = false
	}

	if 0 == len(mc.Outpath) {
		errBuf.WriteString("Outpath: output directory is nil, will use 'doc' default directory.\n")
		mc.Outpath = DEFAULT_OUTPATH
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
