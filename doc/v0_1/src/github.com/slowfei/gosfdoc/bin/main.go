package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/slowfei/gosfcore/utils/filemanager"
	"github.com/slowfei/gosfdoc"
	_ "github.com/slowfei/gosfdoc/lang/golang"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	DEFAULT_GOSFDOC_JSON = "gosfdoc.json"
)

var (
	//	command param
	_lang       *string = nil
	_configFile         = flag.String("config", DEFAULT_GOSFDOC_JSON, "custom config file path.")
	_version            = flag.String("v", "", "output the document version string.")

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
	fmt.Println("'web [8080]' run web server can specify port, default 8080, command by gosfdoc.json ")
	fmt.Println("")
	fmt.Println("other param:")
	flag.PrintDefaults()
	fmt.Println("")
}

/**
 *	parse commond params
 *
 *	@param `args`
 */
func parseCommond(args []string) {
	for _, str := range args {
		commond := "-config="
		if 0 == strings.Index(str, commond) {
			*_configFile = str[len(commond):len(str)]
		}

		commond = "-lang="
		if 0 == strings.Index(str, commond) {
			*_lang = str[len(commond):len(str)]
		}

		commond = "-v="
		if 0 == strings.Index(str, commond) {
			*_version = str[len(commond):len(str)]
		}

		commond = "-version="
		if 0 == strings.Index(str, commond) {
			*_version = str[len(commond):len(str)]
		}
	}
}

func checkConfigPath() string {

	configPath := *_configFile
	if 0 == len(configPath) {
		configPath = DEFAULT_GOSFDOC_JSON
	}

	if !filepath.IsAbs(configPath) {
		configPath = filepath.Join(SFFileManager.GetCmdDir(), configPath)
	}

	exists, isDir, readErr := SFFileManager.Exists(configPath)

	if !exists || isDir || nil != readErr {
		return ""
	}

	return configPath
}

func main() {
	flag.Usage = usage
	flag.Parse()
	args := flag.Args()
	parseCommond(args)

	if nil == args || 0 == len(args) {
		flag.Usage()
		return
	}

	arg := args[0]

	if 0 != len(arg) {
		switch arg {
		case "help":
			flag.Usage()
		case "version":
			fmt.Println(gosfdoc.APPNAME, "v"+gosfdoc.VERSION)
		case "web":
			port := 8080
			if 2 == len(args) {
				i, err := strconv.Atoi(args[1])
				if nil != err {
					fmt.Println("will be used port 8080, error port command input.")
				} else {
					port = i
				}
			}

			configPath := checkConfigPath()
			if 0 == len(configPath) {
				fmt.Println(DEFAULT_GOSFDOC_JSON, "file can not be found.")
				return
			}

			config, err, pass := gosfdoc.ReadConfigFile(configPath)
			if !pass {
				fmt.Println(err.Error())
				return
			}
			path := ""
			expandPath := ""
			if filepath.IsAbs(config.Outpath) {
				path = config.Outpath
			} else {
				expandPath = config.Outpath
				path = config.ScanPath
			}

			fmt.Println("\nweb server run path:", path)
			if 0 == len(expandPath) {
				fmt.Printf("Browser path: http://localhost:%v%v\n", port, "/index.html")
			} else {
				fmt.Printf("Browser path: http://localhost:%v%v%v\n", port, "/"+expandPath, "/index.html")
			}
			startWeb(path, expandPath, port)
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

			version := *_version
			if 0 == len(version) {
				fmt.Println("Please input the document version string info.\n")
				return
			}

			configPath := checkConfigPath()
			if 0 == len(configPath) {
				fmt.Println(DEFAULT_GOSFDOC_JSON, "file invalid, please use 'create' command create config file.")
				return
			}

			if gosfdoc.CheckExistVersion(configPath, version) {
				fmt.Println("Current output version of the document already exists!")
				fmt.Println("Whether to overwrite existing files? (yes/no, y/n)")

				reader := bufio.NewReader(os.Stdin)
				data, _, _ := reader.ReadLine()
				command := string(data)

				if "yes" != command && "y" != command {
					fmt.Println("End output.")
					return
				}
				fmt.Println("")
			}

			outErr, outPass := gosfdoc.Output(configPath, version, func(path string, result gosfdoc.OperateResult, err error) {
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
				if nil != err {
					fmt.Printf("(Error Info: %v)\n", err)
				}
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
