package golang

import (
	// "strings"
	"testing"
)

func TestRegexpType(t *testing.T) {
	testStr := `
type GolangParser struct {
    config  gosfdoc.MainConfig
    indexDB index.IndexDB
}

type OperateResult int
    `
	s := REXType.FindAllSubmatchIndex([]byte(testStr), -1)

	for i := 0; i < len(s); i++ {
		indexs := s[i]
		t.Log(indexs)
		t.Log(testStr[indexs[0]:indexs[1]])
		t.Log(testStr[indexs[2]:indexs[3]])
	}
}
