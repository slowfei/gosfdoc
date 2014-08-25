package main

import (
	"flag"
	"fmt"
	"github.com/slowfei/gosfdoc"
	"strings"
)

var (
	_lang       *string = nil
	_configFile         = flag.String("config", "gosfdoc.json", "custom config file path.")
)

func init() {
	impls := make([]string, 0, 0)
	for k, _ := range gosfdoc.MapParser() {
		impls = append(impls, k)
	}
	implstr := strings.Join(impls, ",")
	_lang = flag.String("lang", "", "[\""+implstr+"\"] specify code language type ',' separated, default all language.")
}

/**
 *  print usage help
 */
func usage() {
	fmt.Println(gosfdoc.APPNAME, "v"+gosfdoc.VERSION)
	fmt.Println("")
	fmt.Println("usage help:")
	fmt.Println("'" + gosfdoc.APPNAME + " create' command init create default gosfdoc.json file, can be custom to modify file content.")
	fmt.Println("'" + gosfdoc.APPNAME + "' command by gosfdoc.json output document ")
	fmt.Println("")
	fmt.Println("other param:")
	flag.PrintDefaults()
	fmt.Println("")
}

func main() {
	flag.Usage = usage
	flag.Parse()

	arg := flag.Args()[0]

	if 0 != len(arg) {
		switch arg {
		case "help":
			flag.Usage()
		case "version":
			fmt.Println(gosfdoc.APPNAME, "v"+gosfdoc.VERSION)
		case "create":

		case "output":

		default:
			fmt.Println("invalid command parameter.")
		}
		return
	}

}
