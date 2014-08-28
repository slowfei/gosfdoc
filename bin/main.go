package main

import (
	"flag"
	"fmt"
	"github.com/slowfei/gosfcore/utils/filemanager"
	"github.com/slowfei/gosfdoc"
	_ "github.com/slowfei/gosfdoc/lang/golang"
	"strings"
)

var (
	//	command param
	_lang       *string = nil
	_configFile         = flag.String("config", "gosfdoc.json", "custom config file path.")

	//	private param
	_currendCodelang = ""
)

func init() {
	impls := make([]string, 0, 0)
	for k, _ := range gosfdoc.MapParser() {
		impls = append(impls, k)
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

		default:
			fmt.Println("invalid command parameter.")
		}
		return
	}

}
