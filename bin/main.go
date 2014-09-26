package main

import (
	"flag"
	"fmt"
	"github.com/slowfei/gosfcore/utils/filemanager"
	"github.com/slowfei/gosfdoc"
	_ "github.com/slowfei/gosfdoc/lang/golang"
	"path/filepath"
	"strings"
)

const (
	DEFAULT_GOSFDOC_JSON = "gosfdoc.json"
)

var (
	//	command param
	_lang       *string = nil
	_configFile         = flag.String("config", DEFAULT_GOSFDOC_JSON, "custom config file path.")

	//	private param
	_currendCodelang = ""
)

func init() {
	impls := make([]string, 0, 0)
	for k, _ := range gosfdoc.MapParser() {
		if k != gosfdoc.NIL_DOC_NAME {
			impls = append(impls, k)
		}
	}
	_currendCodelang = strings.Join(impls, ",")
	_lang = flag.String("lang", "", "[\""+_currendCodelang+"\"] specify code language type ',' separated, default all language.")
}

/**
 *  print usage help
 */
func usage() {
	fmt.Println(gosfdoc.APPNAME, "v"+gosfdoc.VERSION)
	fmt.Println("")
	fmt.Println("usage help:")
	fmt.Println("'create' command init create default gosfdoc.json file, can be custom to modify file content.")
	fmt.Println("'output' command by gosfdoc.json output document ")
	fmt.Println("")
	fmt.Println("other param:")
	flag.PrintDefaults()
	fmt.Println("")
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if nil == flag.Args() || 0 == len(flag.Args()) {
		flag.Usage()
		return
	}

	arg := flag.Args()[0]

	if 0 != len(arg) {
		switch arg {
		case "help":
			flag.Usage()
		case "version":
			fmt.Println(gosfdoc.APPNAME, "v"+gosfdoc.VERSION)
		case "create":
			lang := ""
			if 0 != len(*_lang) {
				lang = *_lang
			} else {
				lang = _currendCodelang
			}

			err, isCre := gosfdoc.CreateConfigFile(SFFileManager.GetCmdDir(), strings.Split(lang, ","))

			if nil != err {
				fmt.Println("Warn:")
				fmt.Println(err.Error())
			}
			if isCre {
				fmt.Println("operation success.")
			} else {
				fmt.Println("error or warn message, check the config file.")
			}
		case "output":
			configPath := *_configFile

			if 0 == len(configPath) {
				configPath = DEFAULT_GOSFDOC_JSON
			}

			if !filepath.IsAbs(configPath) {
				configPath = filepath.Join(SFFileManager.GetCmdDir(), configPath)
			}

			exists, isDir, readErr := SFFileManager.Exists(configPath)

			if !exists || isDir || nil != readErr {
				fmt.Println(DEFAULT_GOSFDOC_JSON, "file invalid, please use 'create' command create config file.")
				return
			}

			outErr, outPass := gosfdoc.Output(configPath, func(path string, result gosfdoc.OperateResult) {
				resultStr := "Invalid:"
				switch result {
				case gosfdoc.ResultFileSuccess:
					resultStr = "Success:"
				case gosfdoc.ResultFileFilter:
					resultStr = "Filter:"
				case gosfdoc.ResultFileNotRead:
					resultStr = "NotRead:"
				case gosfdoc.ResultFileReadErr:
					resultStr = "ReadError:"
				case gosfdoc.ResultFileOutFail:
					resultStr = "OutputFail:"
				}
				fmt.Println(resultStr, path)
			})

			if nil != outErr {
				fmt.Println("Warn or Error:")
				fmt.Println(outErr.Error())
			}
			if outPass {
				fmt.Println("operation complete.")
			} else {
				fmt.Println("operation not success.")
			}

		default:
			fmt.Println("invalid command parameter.")
		}
		return
	}

}
