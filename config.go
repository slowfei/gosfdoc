package main

import (
	"encoding/json"
	"errors"
	"github.com/slowfei/gosfcore/utils/filemanager"
	"io/ioutil"
	"path/filepath"
)

var (
	_defaultConfigJson = `
    {
        "CodeLang"         : "go",
        "Outdir"           : "doc",
        "CopyCode"         : false,
        "Languages"        : {"Default" : "default"},
        "FilterPaths"      : [],
        "FilePaths"        : {}
    }`
)

/**
 *  main config info
 *  output `gosfdoc.json` use
 */
type MainConfig struct {
	CodeLang    string              // code language
	Outdir      string              // output document directory
	CopyCode    bool                // whether output source code to document directory
	Languages   map[string]string   // document support the language. key is show text, value is lang dirctory name
	FilterPaths []string            // filter directory path
	FilePaths   map[string][]string // file path list. format:{"directory":["file1","file2"...]}
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
