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
