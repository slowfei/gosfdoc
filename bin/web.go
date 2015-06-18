//  Copyright 2013 slowfei And The Contributors All rights reserved.
//
//  Software Source Code License Agreement (BSD License)
//
//  Create on 2013-11-30
//  Update on 2014-06-01
//  Email  slowfei@foxmail.com
//  Home   http://www.slowfei.com

//  gosfdoc web server handle
package main

import (
	"fmt"
	"github.com/slowfei/gosfcore/utils/filemanager"
	"github.com/slowfei/leafveingo/template"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	_template   = LVTemplate.NewTemplate()
	_suffix     = "html"
	_charset    = "utf-8"
	_expandPath = ""
)

func htmlOut(rw http.ResponseWriter, req *http.Request) {

	reqPath := req.URL.Path

	if '/' == reqPath[len(reqPath)-1] {
		reqPath += "index." + _suffix
	}

	filePath := filepath.Join(_template.BaseDir(), _expandPath, reqPath)

	isExists, isDir, _ := SFFileManager.Exists(filePath)
	if !isExists || isDir {
		http.NotFound(rw, req)
		return
	}

	if strings.HasSuffix(reqPath, _suffix) {
		e := req.ParseForm()
		if nil != e {
			fmt.Println("parse form error:", e)
			return
		}
		rw.Header().Set("Content-Type", "text/html; charset="+_charset)
		err := _template.Execute(rw, LVTemplate.NewTemplateValue(reqPath, req.Form))
		if nil != err {
			fmt.Println("template error: ", err)
		}
	} else {
		http.ServeFile(rw, req, filePath)
	}

}

func startWeb(path, expandPath string, port int) {

	http.HandleFunc("/", htmlOut)

	_template.SetBaseDir(path)
	_template.SetCache(false)
	_template.SetCompactHTML(false)

	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		fmt.Println("gosfdoc web run error:", err)
	}
}
