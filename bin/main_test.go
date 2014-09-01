package main

import (
	"github.com/slowfei/gosfdoc"
	"testing"
)

func TestCreateConfigFile(t *testing.T) {
	err, sess := gosfdoc.CreateConfigFile("/Users/slowfei/Downloads/test/leafveingo", []string{"go", "java", "js"})
	t.Logf("error:%v \n", err)
	t.Log(sess)
}

func TestOutputWithConfig(t *testing.T) {
	config := new(gosfdoc.MainConfig)
	config.Path = "/Users/slowfei/Downloads/test/leafveingo"
	config.CodeLang = []string{"go"}
	config.Outdir = "doc"
	config.CopyCode = false
	config.HtmlTitle = "Document"
	config.DocTitle = "<b>Document:</b>"
	config.MenuTitle = "<center><b>package</b></center>"
	config.Languages = map[string]string{"default": "Default"}
	config.FilterPaths = make([]string, 0, 0)

	err, pass := gosfdoc.OutputWithConfig(config, func(path string, result gosfdoc.OperateResult) {
		if result == gosfdoc.ResultFileFilter {
			t.Log(path)
		}
	})

	if nil != err {
		t.Error(err)
	}
	t.Log(pass)
}
