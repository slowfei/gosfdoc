package main

import (
	"flag"
	"fmt"
)

const (
	APPNAME = "gosfdoc"
	VERSION = "0.0.1.000"
)

var (
	_lang       = flag.String("lang", "go", "specify code language type. 'go|java|obj-c|javascript'")
	_configFile = flag.String("config", "gosfdoc.json", "custom config file path.")
)

/**
 *  create config file
 */
func createConfigFile() {
	// TODO
}

/**
 *  print usage help
 */
func usage() {
	fmt.Println(APPNAME, "v"+VERSION)
	fmt.Println("")
	fmt.Println("usage help:")
	fmt.Println("'" + APPNAME + " create' command init create default gosfdoc.json file, can be custom to modify file content.")
	fmt.Println("'" + APPNAME + "' command by gosfdoc.json output document ")
	fmt.Println("")
	fmt.Println("command:")
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
			fmt.Println(APPNAME, "v"+VERSION)
		case "create":
			createConfigFile()
		default:
			fmt.Println("invalid command parameter.")
		}
		return
	}

}
