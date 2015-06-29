
## Preview
------

> [const ( APPNAME... ...FILE_NAME_HTML_CONFIG_JSON )](#f_const___APPNAME---_---FILE_NAME_HTML_CONFIG_JSON__)<a name="p_const___APPNAME---_---FILE_NAME_HTML_CONFIG_JSON__"><a/>

> [const ( DEFAULT_CONFIG_FILE_NAME... ...DEFAULT_OUTPATH )](#f_const___DEFAULT_CONFIG_FILE_NAME---_---DEFAULT_OUTPATH__)<a name="p_const___DEFAULT_CONFIG_FILE_NAME---_---DEFAULT_OUTPATH__"><a/>

> [const ( DOC_FILE_SUFFIX... ...NIL_DOC_NAME )](#f_const___DOC_FILE_SUFFIX---_---NIL_DOC_NAME__)<a name="p_const___DOC_FILE_SUFFIX---_---NIL_DOC_NAME__"><a/>

> [const ( ResultFileSuccess... ...ResultDebugErr )](#f_const___ResultFileSuccess---_---ResultDebugErr__)<a name="p_const___ResultFileSuccess---_---ResultDebugErr__"><a/>

> [var ( REXPrivateFile... ...REXDocIndexTitle )](#f_var___REXPrivateFile---_---REXDocIndexTitle__)<a name="p_var___REXPrivateFile---_---REXDocIndexTitle__"><a/>

> [func AddParser(parser DocParser) ](#f_func_AddParser_parser_DocParser__)<a name="p_func_AddParser_parser_DocParser__"><a/>

> [func CheckExistVersion(configPath, version string) bool](#f_func_CheckExistVersion_configPath_version_string__bool)<a name="p_func_CheckExistVersion_configPath_version_string__bool"><a/>

> [func ConverToVersionPath(version string) string](#f_func_ConverToVersionPath_version_string__string)<a name="p_func_ConverToVersionPath_version_string__string"><a/>

> [func CreateConfigFile(dirPath string, langs []string) (error, bool)](#f_func_CreateConfigFile_dirPath_string_langs___string___error_bool_)<a name="p_func_CreateConfigFile_dirPath_string_langs___string___error_bool_"><a/>

> [func FindPrefixFilterTag(src []byte) []byte](#f_func_FindPrefixFilterTag_src___byte____byte)<a name="p_func_FindPrefixFilterTag_src___byte____byte"><a/>

> [func Output(configPath, version string, fileFunc FileResultFunc) (error, bool)](#f_func_Output_configPath_version_string_fileFunc_FileResultFunc___error_bool_)<a name="p_func_Output_configPath_version_string_fileFunc_FileResultFunc___error_bool_"><a/>

> [func OutputWithConfig(config \*MainConfig, version string, fileFunc FileResultFunc) (error, bool)](#f_func_OutputWithConfig_config_+MainConfig_version_string_fileFunc_FileResultFunc___error_bool_)<a name="p_func_OutputWithConfig_config_+MainConfig_version_string_fileFunc_FileResultFunc___error_bool_"><a/>

> [func ReadConfigFile(configFilePath string) (config \*MainConfig, err error, pass bool)](#f_func_ReadConfigFile_configFilePath_string___config_+MainConfig_err_error_pass_bool_)<a name="p_func_ReadConfigFile_configFilePath_string___config_+MainConfig_err_error_pass_bool_"><a/>

> [type About struct](#f_type_About_struct)<a name="p_type_About_struct"><a/>

>> [func NewDefaultAbout() \*About](#f_func_NewDefaultAbout___+About)<a name="p_func_NewDefaultAbout___+About"><a/>

>> [func ParseAbout(fileBuf \*FileBuf) \*About](#f_func_ParseAbout_fileBuf_+FileBuf__+About)<a name="p_func_ParseAbout_fileBuf_+FileBuf__+About"><a/>

>> [func (\*About) WriteFilepath(path string) error](#f_func__+About__WriteFilepath_path_string__error)<a name="p_func__+About__WriteFilepath_path_string__error"><a/>

> [type CodeBlock struct](#f_type_CodeBlock_struct)<a name="p_type_CodeBlock_struct"><a/>

> [type CodeFile struct](#f_type_CodeFile_struct)<a name="p_type_CodeFile_struct"><a/>

> [type CodeFiles struct](#f_type_CodeFiles_struct)<a name="p_type_CodeFiles_struct"><a/>

>> [func NewCodeFiles() \*CodeFiles](#f_func_NewCodeFiles___+CodeFiles)<a name="p_func_NewCodeFiles___+CodeFiles"><a/>

>> [func (\*CodeFiles) FilesLen() int](#f_func__+CodeFiles__FilesLen___int)<a name="p_func__+CodeFiles__FilesLen___int"><a/>

>> [func (\*CodeFiles) IsAllDocFile() bool](#f_func__+CodeFiles__IsAllDocFile___bool)<a name="p_func__+CodeFiles__IsAllDocFile___bool"><a/>

> [type ContentJson struct](#f_type_ContentJson_struct)<a name="p_type_ContentJson_struct"><a/>

>> [func (ContentJson) WriteFilepath(path string) error](#f_func__ContentJson__WriteFilepath_path_string__error)<a name="p_func__ContentJson__WriteFilepath_path_string__error"><a/>

> [type DocConfig struct](#f_type_DocConfig_struct)<a name="p_type_DocConfig_struct"><a/>

> [type DocParser interface](#f_type_DocParser_interface)<a name="p_type_DocParser_interface"><a/>

>> [func MapParser() map[string]DocParser](#f_func_MapParser___map_string_DocParser)<a name="p_func_MapParser___map_string_DocParser"><a/>

> [type Document struct](#f_type_Document_struct)<a name="p_type_Document_struct"><a/>

>> [func ParseDocument(fileBuf \*FileBuf) []Document](#f_func_ParseDocument_fileBuf_+FileBuf____Document)<a name="p_func_ParseDocument_fileBuf_+FileBuf____Document"><a/>

> [type FileBuf struct](#f_type_FileBuf_struct)<a name="p_type_FileBuf_struct"><a/>

>> [func NewFileBuf(fileContent []byte, path string, info os.FileInfo, filter \*regexp.Regexp) \*FileBuf](#f_func_NewFileBuf_fileContent___byte_path_string_info_os-FileInfo_filter_+regexp-Regexp__+FileBuf)<a name="p_func_NewFileBuf_fileContent___byte_path_string_info_os-FileInfo_filter_+regexp-Regexp__+FileBuf"><a/>

>> [func (\*FileBuf) Byte(index int) (byte, bool)](#f_func__+FileBuf__Byte_index_int___byte_bool_)<a name="p_func__+FileBuf__Byte_index_int___byte_bool_"><a/>

>> [func (\*FileBuf) FileInfo() os.FileInfo](#f_func__+FileBuf__FileInfo___os-FileInfo)<a name="p_func__+FileBuf__FileInfo___os-FileInfo"><a/>

>> [func (\*FileBuf) Find(rex \*regexp.Regexp) []byte](#f_func__+FileBuf__Find_rex_+regexp-Regexp____byte)<a name="p_func__+FileBuf__Find_rex_+regexp-Regexp____byte"><a/>

>> [func (\*FileBuf) FindAll(rex \*regexp.Regexp) [][]byte](#f_func__+FileBuf__FindAll_rex_+regexp-Regexp______byte)<a name="p_func__+FileBuf__FindAll_rex_+regexp-Regexp______byte"><a/>

>> [func (\*FileBuf) FindAllSubmatch(rex \*regexp.Regexp) [][][]byte](#f_func__+FileBuf__FindAllSubmatch_rex_+regexp-Regexp________byte)<a name="p_func__+FileBuf__FindAllSubmatch_rex_+regexp-Regexp________byte"><a/>

>> [func (\*FileBuf) FindAllSubmatchIndex(rex \*regexp.Regexp) [][]int](#f_func__+FileBuf__FindAllSubmatchIndex_rex_+regexp-Regexp______int)<a name="p_func__+FileBuf__FindAllSubmatchIndex_rex_+regexp-Regexp______int"><a/>

>> [func (\*FileBuf) FindSubmatch(rex \*regexp.Regexp) [][]byte](#f_func__+FileBuf__FindSubmatch_rex_+regexp-Regexp______byte)<a name="p_func__+FileBuf__FindSubmatch_rex_+regexp-Regexp______byte"><a/>

>> [func (\*FileBuf) FindSubmatchIndex(rex \*regexp.Regexp) []int](#f_func__+FileBuf__FindSubmatchIndex_rex_+regexp-Regexp____int)<a name="p_func__+FileBuf__FindSubmatchIndex_rex_+regexp-Regexp____int"><a/>

>> [func (\*FileBuf) LineLen() int](#f_func__+FileBuf__LineLen___int)<a name="p_func__+FileBuf__LineLen___int"><a/>

>> [func (\*FileBuf) LineNumberByIndex(beginIndex, endIndex int) []int](#f_func__+FileBuf__LineNumberByIndex_beginIndex_endIndex_int____int)<a name="p_func__+FileBuf__LineNumberByIndex_beginIndex_endIndex_int____int"><a/>

>> [func (\*FileBuf) Path() string](#f_func__+FileBuf__Path___string)<a name="p_func__+FileBuf__Path___string"><a/>

>> [func (\*FileBuf) RowByIndex(lineNumber int) []byte](#f_func__+FileBuf__RowByIndex_lineNumber_int____byte)<a name="p_func__+FileBuf__RowByIndex_lineNumber_int____byte"><a/>

>> [func (\*FileBuf) String() string](#f_func__+FileBuf__String___string)<a name="p_func__+FileBuf__String___string"><a/>

>> [func (\*FileBuf) SubBytes(beginIndex, endIndex int) []byte](#f_func__+FileBuf__SubBytes_beginIndex_endIndex_int____byte)<a name="p_func__+FileBuf__SubBytes_beginIndex_endIndex_int____byte"><a/>

>> [func (\*FileBuf) SubNestAllIndex(subNest \*SFSubUtil.SubNest, outBetweens [][]int) [][]int](#f_func__+FileBuf__SubNestAllIndex_subNest_+SFSubUtil-SubNest_outBetweens_____int______int)<a name="p_func__+FileBuf__SubNestAllIndex_subNest_+SFSubUtil-SubNest_outBetweens_____int______int"><a/>

>> [func (\*FileBuf) SubNestAllIndexByBetween(startIndex, endIndex int, subNest \*SFSubUtil.SubNest, outBetweens [][]int) [][]int](#f_func__+FileBuf__SubNestAllIndexByBetween_startIndex_endIndex_int_subNest_+SFSubUtil-SubNest_outBetweens_____int______int)<a name="p_func__+FileBuf__SubNestAllIndexByBetween_startIndex_endIndex_int_subNest_+SFSubUtil-SubNest_outBetweens_____int______int"><a/>

>> [func (\*FileBuf) SubNestGetOutBetweens(nests ...\*SFSubUtil.SubNest) [][]int](#f_func__+FileBuf__SubNestGetOutBetweens_nests_---+SFSubUtil-SubNest______int)<a name="p_func__+FileBuf__SubNestGetOutBetweens_nests_---+SFSubUtil-SubNest______int"><a/>

>> [func (\*FileBuf) SubNestIndex(startIndex int, subNest \*SFSubUtil.SubNest, outBetweens [][]int) []int](#f_func__+FileBuf__SubNestIndex_startIndex_int_subNest_+SFSubUtil-SubNest_outBetweens_____int____int)<a name="p_func__+FileBuf__SubNestIndex_startIndex_int_subNest_+SFSubUtil-SubNest_outBetweens_____int____int"><a/>

>> [func (\*FileBuf) WriteFilepath(path string) error](#f_func__+FileBuf__WriteFilepath_path_string__error)<a name="p_func__+FileBuf__WriteFilepath_path_string__error"><a/>

> [type FileLink struct](#f_type_FileLink_struct)<a name="p_type_FileLink_struct"><a/>

> [type FileResultFunc func](#f_type_FileResultFunc_func)<a name="p_type_FileResultFunc_func"><a/>

> [type Intro struct](#f_type_Intro_struct)<a name="p_type_Intro_struct"><a/>

>> [func NewDefaultIntro() \*Intro](#f_func_NewDefaultIntro___+Intro)<a name="p_func_NewDefaultIntro___+Intro"><a/>

>> [func ParseIntro(fileBuf \*FileBuf) \*Intro](#f_func_ParseIntro_fileBuf_+FileBuf__+Intro)<a name="p_func_ParseIntro_fileBuf_+FileBuf__+Intro"><a/>

>> [func (\*Intro) WriteFilepath(path string) error](#f_func__+Intro__WriteFilepath_path_string__error)<a name="p_func__+Intro__WriteFilepath_path_string__error"><a/>

> [type MainConfig struct](#f_type_MainConfig_struct)<a name="p_type_MainConfig_struct"><a/>

>> [func (\*MainConfig) Check() (error, bool)](#f_func__+MainConfig__Check____error_bool_)<a name="p_func__+MainConfig__Check____error_bool_"><a/>

>> [func (MainConfig) GithubLink(relMDPath string, isToMarkdown bool) string](#f_func__MainConfig__GithubLink_relMDPath_string_isToMarkdown_bool__string)<a name="p_func__MainConfig__GithubLink_relMDPath_string_isToMarkdown_bool__string"><a/>

> [type MenuFile struct](#f_type_MenuFile_struct)<a name="p_type_MenuFile_struct"><a/>

> [type MenuMarkdown struct](#f_type_MenuMarkdown_struct)<a name="p_type_MenuMarkdown_struct"><a/>

> [type OperateResult int](#f_type_OperateResult_int)<a name="p_type_OperateResult_int"><a/>

> [type PackageInfo struct](#f_type_PackageInfo_struct)<a name="p_type_PackageInfo_struct"><a/>

> [type Preview struct](#f_type_Preview_struct)<a name="p_type_Preview_struct"><a/>

> [type SortSet struct](#f_type_SortSet_struct)<a name="p_type_SortSet_struct"><a/>

>> [func (SortSet) Len() int](#f_func__SortSet__Len___int)<a name="p_func__SortSet__Len___int"><a/>

>> [func (SortSet) Less(i, j int) bool](#f_func__SortSet__Less_i_j_int__bool)<a name="p_func__SortSet__Less_i_j_int__bool"><a/>

>> [func (SortSet) Swap(i, j int) ](#f_func__SortSet__Swap_i_j_int__)<a name="p_func__SortSet__Swap_i_j_int__"><a/>

<br/>
### Directory files
[config.go ](../../../../../src/github.com/slowfei/gosfdoc/config.go)[gosfdoc.go ](../../../../../src/github.com/slowfei/gosfdoc/gosfdoc.go)[parse.go ](../../../../../src/github.com/slowfei/gosfdoc/parse.go)[struct.go ](../../../../../src/github.com/slowfei/gosfdoc/struct.go)

## Constants
------
### [source code](../../../../../src/github.com/slowfei/gosfdoc/gosfdoc.go#L30-L53) <a name="f_const___APPNAME---_---FILE_NAME_HTML_CONFIG_JSON__"><a/> [↩](#p_const___APPNAME---_---FILE_NAME_HTML_CONFIG_JSON__) | [#](#f_const___APPNAME---_---FILE_NAME_HTML_CONFIG_JSON__)

<pre><code class='go custom'>const (
	APPNAME = "gosfdoc" //
	VERSION = "0.1.000" //

	DIR_NAME_MAIN_MARKDOWN    = "md"      // save markdown file main directory name
	DIR_NAME_MARKDOWN_DEFAULT = "default" // markdown default directory
	DIR_NAME_SOURCE_CODE      = "src"     // source code save directory
	DIR_NAME_ASSETS           = "assets"  // html use assets file directory

	FILE_SUFFIX_MARKDOWN = ".md"

	FILE_NAME_ABOUT_MD     = "about.md"
	FILE_NAME_INTRO_MD     = "intro.md"
	FILE_NAME_CONTENT_JSON = "content.json"

	FILE_NAME_GOSFDOC_MIN_CSS    = "gosfdoc.min.css"
	FILE_NAME_ASSETS_MIN_JS      = "assets.min.js"
	FILE_NAME_GOSFDOC_MIN_JS     = "gosfdoc.min.js"
	FILE_NAME_GOSFDOC_SRC_MIN_JS = "gosfdoc.src.min.js"

	FILE_NAME_HTML_INDEX       = "index.html"
	FILE_NAME_HTML_SRC         = "src.html"
	FILE_NAME_HTML_CONFIG_JSON = "config.json"
)</code></pre>


### [source code](../../../../../src/github.com/slowfei/gosfdoc/config.go#L24-L27) <a name="f_const___DEFAULT_CONFIG_FILE_NAME---_---DEFAULT_OUTPATH__"><a/> [↩](#p_const___DEFAULT_CONFIG_FILE_NAME---_---DEFAULT_OUTPATH__) | [#](#f_const___DEFAULT_CONFIG_FILE_NAME---_---DEFAULT_OUTPATH__)

<pre><code class='go custom'>const (
	DEFAULT_CONFIG_FILE_NAME = "gosfdoc.json"
	DEFAULT_OUTPATH          = "doc"
)</code></pre>


### [source code](../../../../../src/github.com/slowfei/gosfdoc/parse.go#L22-L26) <a name="f_const___DOC_FILE_SUFFIX---_---NIL_DOC_NAME__"><a/> [↩](#p_const___DOC_FILE_SUFFIX---_---NIL_DOC_NAME__) | [#](#f_const___DOC_FILE_SUFFIX---_---NIL_DOC_NAME__)

<pre><code class='go custom'>const (
	DOC_FILE_SUFFIX = ".dc"      // document file suffix(document comments)
	NIL_DOC_NAME    = "document" // nilDocParser struct use

)</code></pre>


### [source code](../../../../../src/github.com/slowfei/gosfdoc/gosfdoc.go#L111-L119) <a name="f_const___ResultFileSuccess---_---ResultDebugErr__"><a/> [↩](#p_const___ResultFileSuccess---_---ResultDebugErr__) | [#](#f_const___ResultFileSuccess---_---ResultDebugErr__)

<pre><code class='go custom'>const (
	ResultFileSuccess OperateResult = iota
	ResultFileInvalid
	ResultFileNotRead
	ResultFileReadErr
	ResultFileFilter
	ResultFileOutFail
	ResultDebugErr
)</code></pre>


## Variables
------
### [source code](../../../../../src/github.com/slowfei/gosfdoc/gosfdoc.go#L74-L104) <a name="f_var___REXPrivateFile---_---REXDocIndexTitle__"><a/> [↩](#p_var___REXPrivateFile---_---REXDocIndexTitle__) | [#](#f_var___REXPrivateFile---_---REXDocIndexTitle__)
> regex compile variable<br/>
> <br/>


<pre><code class='go custom'>var (
	// private file tag ( //# private-doc-code )
	REXPrivateFile = regexp.MustCompile("#private-(doc|code){1}(-doc|-code)?")
	TagPrivateCode = []byte("code")
	TagPrivateDoc  = []byte("doc")
	// private block tag ( //# private * //# private-end)
	REXPrivateBlock = regexp.MustCompile("[^\\n][\\s]?")

	// parse separate document file package info( //# package-info brief intro row)
	REXDCPackageInfo = regexp.MustCompile("#package-info (.+)")

	// parse about and intro block
	/* * [About|Intro]
	 *  content text or markdown text
	 */
	// [About|Intro]
	// content text or markdown text
	// End
	REXAbout = regexp.MustCompile("(/\\*\\*About[\\s]+(\\s|.)*?[\\s]+\\*/)|(//About[\\s]?([\\s]|.)*?//[Ee][Nn][Dd])")
	REXIntro = regexp.MustCompile("(/\\*\\*Intro[\\s]+(\\s|.)*?[\\s]+\\*/)|(//Intro[\\s]?([\\s]|.)*?//[Ee][Nn][Dd])")

	// parse public document content
	/* * *[z-index-][title]
	 *  document text or markdown text
	 */
	// /[z-index-][title]
	//  document text or markdown text
	// End
	REXDocument      = regexp.MustCompile("(/\\*\\*\\*[^\\*\\s](.+)\\n(\\s|.)*?\\*/)|(///[^/\\s](.+)\\n(\\s|.)*?//[Ee][Nn][Dd])")
	REXDocIndexTitle = regexp.MustCompile("(/\\*\\*\\*|///)(\\d*-)?(.*)?")
)</code></pre>


## Func Details
------
### [func AddParser](../../../../../src/github.com/slowfei/gosfdoc/gosfdoc.go#L208-L212) <a name="f_func_AddParser_parser_DocParser__"><a/> [↩](#p_func_AddParser_parser_DocParser__) | [#](#f_func_AddParser_parser_DocParser__)
> add parser<br/>
> @param parser<br/>
> <br/>


<pre><code class='go custom'>func AddParser(parser DocParser)  { ...... }</code></pre>


### [func CheckExistVersion](../../../../../src/github.com/slowfei/gosfdoc/gosfdoc.go#L329-L339) <a name="f_func_CheckExistVersion_configPath_version_string__bool"><a/> [↩](#p_func_CheckExistVersion_configPath_version_string__bool) | [#](#f_func_CheckExistVersion_configPath_version_string__bool)
> check whether there are version info<br/>
> @param `configPath` config path<br/>
> @param `version` check version string<br/>
> <br/>


<pre><code class='go custom'>func CheckExistVersion(configPath, version string) bool { ...... }</code></pre>


### [func ConverToVersionPath](../../../../../src/github.com/slowfei/gosfdoc/gosfdoc.go#L344-L349) <a name="f_func_ConverToVersionPath_version_string__string"><a/> [↩](#p_func_ConverToVersionPath_version_string__string) | [#](#f_func_ConverToVersionPath_version_string__string)
> conver version to use the path info<br/>
> <br/>


<pre><code class='go custom'>func ConverToVersionPath(version string) string { ...... }</code></pre>


### [func CreateConfigFile](../../../../../src/github.com/slowfei/gosfdoc/gosfdoc.go#L359-L445) <a name="f_func_CreateConfigFile_dirPath_string_langs___string___error_bool_"><a/> [↩](#p_func_CreateConfigFile_dirPath_string_langs___string___error_bool_) | [#](#f_func_CreateConfigFile_dirPath_string_langs___string___error_bool_)
> create config file<br/>
> @param `dirPath` directory path<br/>
> @param `langs`   specify code language, nil is all language, value is parser name.<br/>
> @return `error`  warn or error message<br/>
> @return `bool`   true is operation success<br/>
> <br/>


<pre><code class='go custom'>func CreateConfigFile(dirPath string, langs []string) (error, bool) { ...... }</code></pre>


### [func FindPrefixFilterTag](../../../../../src/github.com/slowfei/gosfdoc/parse.go#L466-L477) <a name="f_func_FindPrefixFilterTag_src___byte____byte"><a/> [↩](#p_func_FindPrefixFilterTag_src___byte____byte) | [#](#f_func_FindPrefixFilterTag_src___byte____byte)
> find prefix filter tag index<br/>
> //<br/>
> // content ("// ") is prefix tag<br/>
> //<br/>
> see var _prefixFilterTags<br/>
> <br/>


<pre><code class='go custom'>func FindPrefixFilterTag(src []byte) []byte { ...... }</code></pre>


### [func Output](../../../../../src/github.com/slowfei/gosfdoc/gosfdoc.go#L456-L462) <a name="f_func_Output_configPath_version_string_fileFunc_FileResultFunc___error_bool_"><a/> [↩](#p_func_Output_configPath_version_string_fileFunc_FileResultFunc___error_bool_) | [#](#f_func_Output_configPath_version_string_fileFunc_FileResultFunc___error_bool_)
> build output document<br/>
> @param `configPath` config file path<br/>
> @param `version`    output document version<br/>
> @param `fileFunc`<br/>
> @return `error` warn or error message<br/>
> @return `bool`  true is operation success<br/>
> <br/>


<pre><code class='go custom'>func Output(configPath, version string, fileFunc FileResultFunc) (error, bool) { ...... }</code></pre>


### [func OutputWithConfig](../../../../../src/github.com/slowfei/gosfdoc/gosfdoc.go#L472-L549) <a name="f_func_OutputWithConfig_config_+MainConfig_version_string_fileFunc_FileResultFunc___error_bool_"><a/> [↩](#p_func_OutputWithConfig_config_+MainConfig_version_string_fileFunc_FileResultFunc___error_bool_) | [#](#f_func_OutputWithConfig_config_+MainConfig_version_string_fileFunc_FileResultFunc___error_bool_)
> build output document with config content<br/>
> @param `config`<br/>
> @param `version` e.g: "v=1.0"<br/>
> @return `error` warn or error message<br/>
> @return `bool`  true is operation success<br/>
> <br/>


<pre><code class='go custom'>func OutputWithConfig(config *MainConfig, version string, fileFunc FileResultFunc) (error, bool) { ...... }</code></pre>


### [func ReadConfigFile](../../../../../src/github.com/slowfei/gosfdoc/gosfdoc.go#L233-L258) <a name="f_func_ReadConfigFile_configFilePath_string___config_+MainConfig_err_error_pass_bool_"><a/> [↩](#p_func_ReadConfigFile_configFilePath_string___config_+MainConfig_err_error_pass_bool_) | [#](#f_func_ReadConfigFile_configFilePath_string___config_+MainConfig_err_error_pass_bool_)
> read config file<br/>
> @param `configFilePath`<br/>
> @return `config`<br/>
> @return `err`   contains warn info<br/>
> @return `pass`  true is valid file (pass does not mean that there are no errors)<br/>
> <br/>


<pre><code class='go custom'>func ReadConfigFile(configFilePath string) (config *MainConfig, err error, pass bool) { ...... }</code></pre>


### [type About struct](../../../../../src/github.com/slowfei/gosfdoc/struct.go#L508-L510) <a name="f_type_About_struct"><a/> [↩](#p_type_About_struct) | [#](#f_type_About_struct)
> markdown about<br/>
> <br/>


<pre><code class='go custom'>type About struct {
	Content []byte
}</code></pre>


### [func NewDefaultAbout](../../../../../src/github.com/slowfei/gosfdoc/struct.go#L517-L519) <a name="f_func_NewDefaultAbout___+About"><a/> [↩](#p_func_NewDefaultAbout___+About) | [#](#f_func_NewDefaultAbout___+About)
> new default about<br/>
> @return pointer type<br/>
> <br/>


<pre><code class='go custom'>func NewDefaultAbout() *About { ...... }</code></pre>


### [func ParseAbout](../../../../../src/github.com/slowfei/gosfdoc/parse.go#L371-L380) <a name="f_func_ParseAbout_fileBuf_+FileBuf__+About"><a/> [↩](#p_func_ParseAbout_fileBuf_+FileBuf__+About) | [#](#f_func_ParseAbout_fileBuf_+FileBuf__+About)
> commons parse file about content<br/>
> @param `fileBuf`<br/>
> @return about content<br/>
> <br/>


<pre><code class='go custom'>func ParseAbout(fileBuf *FileBuf) *About { ...... }</code></pre>


### [func (\*About) WriteFilepath](../../../../../src/github.com/slowfei/gosfdoc/struct.go#L527-L532) <a name="f_func__+About__WriteFilepath_path_string__error"><a/> [↩](#p_func__+About__WriteFilepath_path_string__error) | [#](#f_func__+About__WriteFilepath_path_string__error)
> output file<br/>
> @param `path` output full path<br/>
> @return<br/>
> <br/>


<pre><code class='go custom'>func (*About) WriteFilepath(path string) error { ...... }</code></pre>


### [type CodeBlock struct](../../../../../src/github.com/slowfei/gosfdoc/struct.go#L603-L613) <a name="f_type_CodeBlock_struct"><a/> [↩](#p_type_CodeBlock_struct) | [#](#f_type_CodeBlock_struct)
> body code block struct<br/>
> <br/>


<pre><code class='go custom'>type CodeBlock struct {
	SortTag        string // sort tag
	MenuTitle      string // left navigation menu title
	Title          string // function name or custom title
	Anchor         string // function anchor text.
	Desc           string // description markdown text or plain text
	Code           string // show code text
	CodeLang       string // source code lang type string
	SourceFileName string // source code file name
	FileLines      []int  // block where the file line [5,10] is L5-L10
}</code></pre>


### [type CodeFile struct](../../../../../src/github.com/slowfei/gosfdoc/struct.go#L401-L407) <a name="f_type_CodeFile_struct"><a/> [↩](#p_type_CodeFile_struct) | [#](#f_type_CodeFile_struct)
> source code file<br/>
> <br/>


<pre><code class='go custom'>type CodeFile struct {
	parser      DocParser  // file parser
	docs        []Document // current file public documents
	FileCont    *FileBuf   // file buffer content
	PrivateDoc  bool       // if private document not output
	PrivateCode bool       // if private source code not output
}</code></pre>


### [type CodeFiles struct](../../../../../src/github.com/slowfei/gosfdoc/struct.go#L412-L414) <a name="f_type_CodeFiles_struct"><a/> [↩](#p_type_CodeFiles_struct) | [#](#f_type_CodeFiles_struct)
> source code file list<br/>
> <br/>


<pre><code class='go custom'>type CodeFiles struct {
	files *list.List
}</code></pre>


### [func NewCodeFiles](../../../../../src/github.com/slowfei/gosfdoc/struct.go#L419-L423) <a name="f_func_NewCodeFiles___+CodeFiles"><a/> [↩](#p_func_NewCodeFiles___+CodeFiles) | [#](#f_func_NewCodeFiles___+CodeFiles)
> new CodeFiles<br/>
> <br/>


<pre><code class='go custom'>func NewCodeFiles() *CodeFiles { ...... }</code></pre>


### [func (\*CodeFiles) FilesLen](../../../../../src/github.com/slowfei/gosfdoc/struct.go#L481-L483) <a name="f_func__+CodeFiles__FilesLen___int"><a/> [↩](#p_func__+CodeFiles__FilesLen___int) | [#](#f_func__+CodeFiles__FilesLen___int)
> file list storage length<br/>
> @return file number<br/>
> <br/>


<pre><code class='go custom'>func (*CodeFiles) FilesLen() int { ...... }</code></pre>


### [func (\*CodeFiles) IsAllDocFile](../../../../../src/github.com/slowfei/gosfdoc/struct.go#L442-L458) <a name="f_func__+CodeFiles__IsAllDocFile___bool"><a/> [↩](#p_func__+CodeFiles__IsAllDocFile___bool) | [#](#f_func__+CodeFiles__IsAllDocFile___bool)
> is all document files<br/>
> @return all document is true<br/>
> <br/>


<pre><code class='go custom'>func (*CodeFiles) IsAllDocFile() bool { ...... }</code></pre>


### [type ContentJson struct](../../../../../src/github.com/slowfei/gosfdoc/struct.go#L488-L492) <a name="f_type_ContentJson_struct"><a/> [↩](#p_type_ContentJson_struct) | [#](#f_type_ContentJson_struct)
> output `content.json`<br/>
> <br/>


<pre><code class='go custom'>type ContentJson struct {
	HtmlTitle string // html document title
	DocTitle  string // html top show title
	MenuTitle string // html left menu title
}</code></pre>


### [func (ContentJson) WriteFilepath](../../../../../src/github.com/slowfei/gosfdoc/struct.go#L497-L503) <a name="f_func__ContentJson__WriteFilepath_path_string__error"><a/> [↩](#p_func__ContentJson__WriteFilepath_path_string__error) | [#](#f_func__ContentJson__WriteFilepath_path_string__error)
> output write file path<br/>
> <br/>


<pre><code class='go custom'>func (ContentJson) WriteFilepath(path string) error { ...... }</code></pre>


### [type DocConfig struct](../../../../../src/github.com/slowfei/gosfdoc/config.go#L268-L278) <a name="f_type_DocConfig_struct"><a/> [↩](#p_type_DocConfig_struct) | [#](#f_type_DocConfig_struct)
> document directory html javascript use config<br/>
> output `config.json`<br/>
> <br/>


<pre><code class='go custom'>type DocConfig struct {
	ContentJson string              // content json file
	IntroMd     string              // intro markdown file
	AboutMd     string              // about markdown file
	Languages   []map[string]string // key is directory name, value is show text
	LinkRoot    bool                // is link root directory
	AppendPath  string              // append output source code and markdown relative path(scan path join)
	Versions    []string            // output document versions
	Markdowns   []MenuMarkdown      // markdown info list
	Files       []MenuFile          // source code file links
}</code></pre>


### [type DocParser interface](../../../../../src/github.com/slowfei/gosfdoc/gosfdoc.go#L133-L194) <a name="f_type_DocParser_interface"><a/> [↩](#p_type_DocParser_interface) | [#](#f_type_DocParser_interface)
> document parser<br/>
> <br/>


<pre><code class='go custom'>type DocParser interface {

	/**
	 *  parser name
	 *
	 *  @return
	 */
	Name() string

	/**
	 *  check file
	 *  detecting whether the file is a valid file
	 *
	 *  @param `parh` file path
	 *  @param `info` file info
	 *  @return true is valid file
	 */
	CheckFile(path string, info os.FileInfo) bool

	/**
	 *  each file the content
	 *  can be create keyword index and other operations
	 *
	 *  @param `filebuf`    file content buffer
	 */
	EachIndexFile(filebuf *FileBuf)

	/**
	 *  parse file preview tag
	 *
	 *  @param `filebuf` file content buffer
	 *  @return slice
	 */
	ParsePreview(filebuf *FileBuf) []Preview

	/**
	 *  parse code block tag
	 *
	 *  @param `filebuf` file content buffer
	 *  @return slice
	 */
	ParseCodeblock(filebuf *FileBuf) []CodeBlock

	/**
	 *  parse directory package info
	 *  each file directory parse string join
	 *
	 *  @param `filebuf`
	 *  @return string file parse the only string
	 */
	ParsePackageInfo(filebuf *FileBuf) string

	/**
	 *  parse start
	 */
	ParseStart(config MainConfig)

	/**
	 *  parse end
	 */
	ParseEnd()
}</code></pre>


### [func MapParser](../../../../../src/github.com/slowfei/gosfdoc/gosfdoc.go#L221-L223) <a name="f_func_MapParser___map_string_DocParser"><a/> [↩](#p_func_MapParser___map_string_DocParser) | [#](#f_func_MapParser___map_string_DocParser)
> get parsers<br/>
> key is parser name<br/>
> value is parser implement<br/>
> @return<br/>
> <br/>


<pre><code class='go custom'>func MapParser() map[string]DocParser { ...... }</code></pre>


### [type Document struct](../../../../../src/github.com/slowfei/gosfdoc/struct.go#L583-L587) <a name="f_type_Document_struct"><a/> [↩](#p_type_Document_struct) | [#](#f_type_Document_struct)
> document struct info<br/>
> <br/>


<pre><code class='go custom'>type Document struct {
	SortTag int    // sort tag
	Title   string // title plain text
	Content string // markdown text or plain text
}</code></pre>


### [func ParseDocument](../../../../../src/github.com/slowfei/gosfdoc/parse.go#L290-L363) <a name="f_func_ParseDocument_fileBuf_+FileBuf____Document"><a/> [↩](#p_func_ParseDocument_fileBuf_+FileBuf____Document) | [#](#f_func_ParseDocument_fileBuf_+FileBuf____Document)
> parse public document content<br/>
> @param `fileBuf`<br/>
> @return document array<br/>
> <br/>


<pre><code class='go custom'>func ParseDocument(fileBuf *FileBuf) []Document { ...... }</code></pre>


### [type FileBuf struct](../../../../../src/github.com/slowfei/gosfdoc/struct.go#L57-L63) <a name="f_type_FileBuf_struct"><a/> [↩](#p_type_FileBuf_struct) | [#](#f_type_FileBuf_struct)
> file content buffer<br/>
> <br/>


<pre><code class='go custom'>type FileBuf struct {
	path       string
	fileInfo   os.FileInfo
	buf        []byte
	lineLenSum []int       // 记录每行长度的总和
	UserData   interface{} // 自定义存储数据
}</code></pre>


### [func NewFileBuf](../../../../../src/github.com/slowfei/gosfdoc/struct.go#L73-L97) <a name="f_func_NewFileBuf_fileContent___byte_path_string_info_os-FileInfo_filter_+regexp-Regexp__+FileBuf"><a/> [↩](#p_func_NewFileBuf_fileContent___byte_path_string_info_os-FileInfo_filter_+regexp-Regexp__+FileBuf) | [#](#f_func_NewFileBuf_fileContent___byte_path_string_info_os-FileInfo_filter_+regexp-Regexp__+FileBuf)
> new file buffer<br/>
> @param `fileContent`<br/>
> @param `path` file path<br/>
> @param `info` file info<br/>
> @param replace regexp, replace text to empty(''), call regexp.ReplaceAll func<br/>
> <br/>


<pre><code class='go custom'>func NewFileBuf(fileContent []byte, path string, info os.FileInfo, filter *regexp.Regexp) *FileBuf { ...... }</code></pre>


### [func (\*FileBuf) Byte](../../../../../src/github.com/slowfei/gosfdoc/struct.go#L348-L360) <a name="f_func__+FileBuf__Byte_index_int___byte_bool_"><a/> [↩](#p_func__+FileBuf__Byte_index_int___byte_bool_) | [#](#f_func__+FileBuf__Byte_index_int___byte_bool_)
> by index get file buffer byte<br/>
> @param `index` buffer index<br/>
> @return `byte`<br/>
> @return `bool` success return true<br/>
> <br/>


<pre><code class='go custom'>func (*FileBuf) Byte(index int) (byte, bool) { ...... }</code></pre>


### [func (\*FileBuf) FileInfo](../../../../../src/github.com/slowfei/gosfdoc/struct.go#L385-L387) <a name="f_func__+FileBuf__FileInfo___os-FileInfo"><a/> [↩](#p_func__+FileBuf__FileInfo___os-FileInfo) | [#](#f_func__+FileBuf__FileInfo___os-FileInfo)
> get file info<br/>
> @return<br/>
> <br/>


<pre><code class='go custom'>func (*FileBuf) FileInfo() os.FileInfo { ...... }</code></pre>


### [func (\*FileBuf) Find](../../../../../src/github.com/slowfei/gosfdoc/struct.go#L114-L116) <a name="f_func__+FileBuf__Find_rex_+regexp-Regexp____byte"><a/> [↩](#p_func__+FileBuf__Find_rex_+regexp-Regexp____byte) | [#](#f_func__+FileBuf__Find_rex_+regexp-Regexp____byte)
> regexp find bytes<br/>
> @param `rex`<br/>
> @return<br/>
> <br/>


<pre><code class='go custom'>func (*FileBuf) Find(rex *regexp.Regexp) []byte { ...... }</code></pre>


### [func (\*FileBuf) FindAll](../../../../../src/github.com/slowfei/gosfdoc/struct.go#L134-L136) <a name="f_func__+FileBuf__FindAll_rex_+regexp-Regexp______byte"><a/> [↩](#p_func__+FileBuf__FindAll_rex_+regexp-Regexp______byte) | [#](#f_func__+FileBuf__FindAll_rex_+regexp-Regexp______byte)
> regexp find all bytes<br/>
> @param `rex`<br/>
> @return<br/>
> <br/>


<pre><code class='go custom'>func (*FileBuf) FindAll(rex *regexp.Regexp) [][]byte { ...... }</code></pre>


### [func (\*FileBuf) FindAllSubmatch](../../../../../src/github.com/slowfei/gosfdoc/struct.go#L144-L146) <a name="f_func__+FileBuf__FindAllSubmatch_rex_+regexp-Regexp________byte"><a/> [↩](#p_func__+FileBuf__FindAllSubmatch_rex_+regexp-Regexp________byte) | [#](#f_func__+FileBuf__FindAllSubmatch_rex_+regexp-Regexp________byte)
> Regexp.FindAllSubmatch<br/>
>  @param `rex`<br/>
>  @return<br/>
> <br/>


<pre><code class='go custom'>func (*FileBuf) FindAllSubmatch(rex *regexp.Regexp) [][][]byte { ...... }</code></pre>


### [func (\*FileBuf) FindAllSubmatchIndex](../../../../../src/github.com/slowfei/gosfdoc/struct.go#L164-L166) <a name="f_func__+FileBuf__FindAllSubmatchIndex_rex_+regexp-Regexp______int"><a/> [↩](#p_func__+FileBuf__FindAllSubmatchIndex_rex_+regexp-Regexp______int) | [#](#f_func__+FileBuf__FindAllSubmatchIndex_rex_+regexp-Regexp______int)
> Regexp.FindAllSubmatchIndex<br/>
> @param `rex`<br/>
> @return<br/>
> <br/>


<pre><code class='go custom'>func (*FileBuf) FindAllSubmatchIndex(rex *regexp.Regexp) [][]int { ...... }</code></pre>


### [func (\*FileBuf) FindSubmatch](../../../../../src/github.com/slowfei/gosfdoc/struct.go#L124-L126) <a name="f_func__+FileBuf__FindSubmatch_rex_+regexp-Regexp______byte"><a/> [↩](#p_func__+FileBuf__FindSubmatch_rex_+regexp-Regexp______byte) | [#](#f_func__+FileBuf__FindSubmatch_rex_+regexp-Regexp______byte)
> regexp find submatch bytes<br/>
> @param `rex`<br/>
> @return<br/>
> <br/>


<pre><code class='go custom'>func (*FileBuf) FindSubmatch(rex *regexp.Regexp) [][]byte { ...... }</code></pre>


### [func (\*FileBuf) FindSubmatchIndex](../../../../../src/github.com/slowfei/gosfdoc/struct.go#L154-L156) <a name="f_func__+FileBuf__FindSubmatchIndex_rex_+regexp-Regexp____int"><a/> [↩](#p_func__+FileBuf__FindSubmatchIndex_rex_+regexp-Regexp____int) | [#](#f_func__+FileBuf__FindSubmatchIndex_rex_+regexp-Regexp____int)
> Regexp.FindSubmatchIndex<br/>
> @param `rex`<br/>
> @return<br/>
> <br/>


<pre><code class='go custom'>func (*FileBuf) FindSubmatchIndex(rex *regexp.Regexp) []int { ...... }</code></pre>


### [func (\*FileBuf) LineLen](../../../../../src/github.com/slowfei/gosfdoc/struct.go#L367-L369) <a name="f_func__+FileBuf__LineLen___int"><a/> [↩](#p_func__+FileBuf__LineLen___int) | [#](#f_func__+FileBuf__LineLen___int)
> get line length<br/>
> @return int<br/>
> <br/>


<pre><code class='go custom'>func (*FileBuf) LineLen() int { ...... }</code></pre>


### [func (\*FileBuf) LineNumberByIndex](../../../../../src/github.com/slowfei/gosfdoc/struct.go#L233-L275) <a name="f_func__+FileBuf__LineNumberByIndex_beginIndex_endIndex_int____int"><a/> [↩](#p_func__+FileBuf__LineNumberByIndex_beginIndex_endIndex_int____int) | [#](#f_func__+FileBuf__LineNumberByIndex_beginIndex_endIndex_int____int)
> line number by begin and end index<br/>
> @param `beginIndex` buffer byte begin index<br/>
> @param `endIndex`	end index<br/>
> @return []int [start line,end line], line number 1 start.<br/>
> <br/>


<pre><code class='go custom'>func (*FileBuf) LineNumberByIndex(beginIndex, endIndex int) []int { ...... }</code></pre>


### [func (\*FileBuf) Path](../../../../../src/github.com/slowfei/gosfdoc/struct.go#L376-L378) <a name="f_func__+FileBuf__Path___string"><a/> [↩](#p_func__+FileBuf__Path___string) | [#](#f_func__+FileBuf__Path___string)
> get file path<br/>
> @return<br/>
> <br/>


<pre><code class='go custom'>func (*FileBuf) Path() string { ...... }</code></pre>


### [func (\*FileBuf) RowByIndex](../../../../../src/github.com/slowfei/gosfdoc/struct.go#L283-L318) <a name="f_func__+FileBuf__RowByIndex_lineNumber_int____byte"><a/> [↩](#p_func__+FileBuf__RowByIndex_lineNumber_int____byte) | [#](#f_func__+FileBuf__RowByIndex_lineNumber_int____byte)
> get row content by line number 1 start.<br/>
> @param `lineNumber` line number<br/>
> @param	content of the specified line number<br/>
> <br/>


<pre><code class='go custom'>func (*FileBuf) RowByIndex(lineNumber int) []byte { ...... }</code></pre>


### [func (\*FileBuf) String](../../../../../src/github.com/slowfei/gosfdoc/struct.go#L394-L396) <a name="f_func__+FileBuf__String___string"><a/> [↩](#p_func__+FileBuf__String___string) | [#](#f_func__+FileBuf__String___string)
> buffer to string<br/>
> @return<br/>
> <br/>


<pre><code class='go custom'>func (*FileBuf) String() string { ...... }</code></pre>


### [func (\*FileBuf) SubBytes](../../../../../src/github.com/slowfei/gosfdoc/struct.go#L327-L339) <a name="f_func__+FileBuf__SubBytes_beginIndex_endIndex_int____byte"><a/> [↩](#p_func__+FileBuf__SubBytes_beginIndex_endIndex_int____byte) | [#](#f_func__+FileBuf__SubBytes_beginIndex_endIndex_int____byte)
> extracts the file buffer from a bytes<br/>
> @param `beginIndex`<br/>
> @param `endIndex`<br/>
> @return bytes<br/>
> <br/>


<pre><code class='go custom'>func (*FileBuf) SubBytes(beginIndex, endIndex int) []byte { ...... }</code></pre>


### [func (\*FileBuf) SubNestAllIndex](../../../../../src/github.com/slowfei/gosfdoc/struct.go#L193-L195) <a name="f_func__+FileBuf__SubNestAllIndex_subNest_+SFSubUtil-SubNest_outBetweens_____int______int"><a/> [↩](#p_func__+FileBuf__SubNestAllIndex_subNest_+SFSubUtil-SubNest_outBetweens_____int______int) | [#](#f_func__+FileBuf__SubNestAllIndex_subNest_+SFSubUtil-SubNest_outBetweens_____int______int)
> all blocks subset<br/>
> @param `subNest`<br/>
> @param `outBetweens` rule out between index<br/>
> @return buffer start and end index list<br/>
> <br/>


<pre><code class='go custom'>func (*FileBuf) SubNestAllIndex(subNest *SFSubUtil.SubNest, outBetweens [][]int) [][]int { ...... }</code></pre>


### [func (\*FileBuf) SubNestAllIndexByBetween](../../../../../src/github.com/slowfei/gosfdoc/struct.go#L206-L214) <a name="f_func__+FileBuf__SubNestAllIndexByBetween_startIndex_endIndex_int_subNest_+SFSubUtil-SubNest_outBetweens_____int______int"><a/> [↩](#p_func__+FileBuf__SubNestAllIndexByBetween_startIndex_endIndex_int_subNest_+SFSubUtil-SubNest_outBetweens_____int______int) | [#](#f_func__+FileBuf__SubNestAllIndexByBetween_startIndex_endIndex_int_subNest_+SFSubUtil-SubNest_outBetweens_____int______int)
> all blocks subset by buffer between index<br/>
> @param `startIndex`<br/>
> @param `endIndex`<br/>
> @param `subNest`<br/>
> @param `outBetweens` rule out between index<br/>
> @return buffer start and end index list<br/>
> <br/>


<pre><code class='go custom'>func (*FileBuf) SubNestAllIndexByBetween(startIndex, endIndex int, subNest *SFSubUtil.SubNest, outBetweens [][]int) [][]int { ...... }</code></pre>


### [func (\*FileBuf) SubNestGetOutBetweens](../../../../../src/github.com/slowfei/gosfdoc/struct.go#L222-L224) <a name="f_func__+FileBuf__SubNestGetOutBetweens_nests_---+SFSubUtil-SubNest______int"><a/> [↩](#p_func__+FileBuf__SubNestGetOutBetweens_nests_---+SFSubUtil-SubNest______int) | [#](#f_func__+FileBuf__SubNestGetOutBetweens_nests_---+SFSubUtil-SubNest______int)
> get between rule out points<br/>
> @param `nests` SubNest objects<br/>
> @return data source points [0] is start point [1] is end point<br/>
> <br/>


<pre><code class='go custom'>func (*FileBuf) SubNestGetOutBetweens(nests ...*SFSubUtil.SubNest) [][]int { ...... }</code></pre>


### [func (\*FileBuf) SubNestIndex](../../../../../src/github.com/slowfei/gosfdoc/struct.go#L176-L184) <a name="f_func__+FileBuf__SubNestIndex_startIndex_int_subNest_+SFSubUtil-SubNest_outBetweens_____int____int"><a/> [↩](#p_func__+FileBuf__SubNestIndex_startIndex_int_subNest_+SFSubUtil-SubNest_outBetweens_____int____int) | [#](#f_func__+FileBuf__SubNestIndex_startIndex_int_subNest_+SFSubUtil-SubNest_outBetweens_____int____int)
> block subset return range index<br/>
> param `startIndex` buffer start index<br/>
> param `subNest`<br/>
> param `outBetweens` rule out between index<br/>
> return buffer start and end index<br/>
> <br/>


<pre><code class='go custom'>func (*FileBuf) SubNestIndex(startIndex int, subNest *SFSubUtil.SubNest, outBetweens [][]int) []int { ...... }</code></pre>


### [func (\*FileBuf) WriteFilepath](../../../../../src/github.com/slowfei/gosfdoc/struct.go#L104-L106) <a name="f_func__+FileBuf__WriteFilepath_path_string__error"><a/> [↩](#p_func__+FileBuf__WriteFilepath_path_string__error) | [#](#f_func__+FileBuf__WriteFilepath_path_string__error)
> out file<br/>
> @param `path` out path<br/>
> <br/>


<pre><code class='go custom'>func (*FileBuf) WriteFilepath(path string) error { ...... }</code></pre>


### [type FileLink struct](../../../../../src/github.com/slowfei/gosfdoc/struct.go#L573-L578) <a name="f_type_FileLink_struct"><a/> [↩](#p_type_FileLink_struct) | [#](#f_type_FileLink_struct)

<pre><code class='go custom'>type FileLink struct {
	menuName string `json:"-"` // type belongs
	Filename string // a tag show text
	Link     string // a tag link

}</code></pre>


### [type FileResultFunc func](../../../../../src/github.com/slowfei/gosfdoc/gosfdoc.go#L127-L127) <a name="f_type_FileResultFunc_func"><a/> [↩](#p_type_FileResultFunc_func) | [#](#f_type_FileResultFunc_func)
> file scan result func<br/>
> @param `path`<br/>
> @param `result`<br/>
> <br/>


<pre><code class='go custom'>type FileResultFunc func</code></pre>


### [type Intro struct](../../../../../src/github.com/slowfei/gosfdoc/struct.go#L537-L539) <a name="f_type_Intro_struct"><a/> [↩](#p_type_Intro_struct) | [#](#f_type_Intro_struct)
> markdown intro<br/>
> <br/>


<pre><code class='go custom'>type Intro struct {
	Content []byte
}</code></pre>


### [func NewDefaultIntro](../../../../../src/github.com/slowfei/gosfdoc/struct.go#L546-L548) <a name="f_func_NewDefaultIntro___+Intro"><a/> [↩](#p_func_NewDefaultIntro___+Intro) | [#](#f_func_NewDefaultIntro___+Intro)
> new default intro<br/>
> @return pointer type<br/>
> <br/>


<pre><code class='go custom'>func NewDefaultIntro() *Intro { ...... }</code></pre>


### [func ParseIntro](../../../../../src/github.com/slowfei/gosfdoc/parse.go#L388-L397) <a name="f_func_ParseIntro_fileBuf_+FileBuf__+Intro"><a/> [↩](#p_func_ParseIntro_fileBuf_+FileBuf__+Intro) | [#](#f_func_ParseIntro_fileBuf_+FileBuf__+Intro)
> commons parse file introduction content<br/>
> @param `fileBuf`<br/>
> @return introduction content<br/>
> <br/>


<pre><code class='go custom'>func ParseIntro(fileBuf *FileBuf) *Intro { ...... }</code></pre>


### [func (\*Intro) WriteFilepath](../../../../../src/github.com/slowfei/gosfdoc/struct.go#L556-L561) <a name="f_func__+Intro__WriteFilepath_path_string__error"><a/> [↩](#p_func__+Intro__WriteFilepath_path_string__error) | [#](#f_func__+Intro__WriteFilepath_path_string__error)
> output file<br/>
> @param `path` output full path<br/>
> @return<br/>
> <br/>


<pre><code class='go custom'>func (*Intro) WriteFilepath(path string) error { ...... }</code></pre>


### [type MainConfig struct](../../../../../src/github.com/slowfei/gosfdoc/config.go#L51-L66) <a name="f_type_MainConfig_struct"><a/> [↩](#p_type_MainConfig_struct) | [#](#f_type_MainConfig_struct)
> main config info<br/>
> output `gosfdoc.json` use<br/>
> <br/>


<pre><code class='go custom'>type MainConfig struct {
	path           string              `json:"-"` // private handle path, save console command path.
	currentVersion string              `json:"-"` // current output version, private record.
	DocUrl         string              // custom link url to document http. e.g.: http://slowfei.github.io/gosfdoc/index.html
	ScanPath       string              // scan document info file path, relative or absolute path, is "/" scan current config file directory path.
	CodeLang       []string            // code languages
	Outpath        string              // output document path, relative or absolute path.
	OutAppendPath  string              // append output source code and markdown relative path(scan path join). defalut ""
	CopyCode       bool                // copy source code to document directory. default false
	CodeLinkRoot   bool                // source code link to root directory, 'CopyCode' is true was invalid, default true
	HtmlTitle      string              // document html show title
	DocTitle       string              // html top tabbar show title
	MenuTitle      string              // html left menu show title
	Languages      []map[string]string // document support the language. key is directory name, value is show text.
	FilterPaths    []string            // filter path, relative or absolute path
}</code></pre>


### [func (\*MainConfig) Check](../../../../../src/github.com/slowfei/gosfdoc/config.go#L104-L181) <a name="f_func__+MainConfig__Check____error_bool_"><a/> [↩](#p_func__+MainConfig__Check____error_bool_) | [#](#f_func__+MainConfig__Check____error_bool_)
> check config param value<br/>
> error value will update default.<br/>
> @return error<br/>
> @return bool    fatal error is false, pass is true. (pass does not mean that there are no errors)<br/>
> <br/>


<pre><code class='go custom'>func (*MainConfig) Check() (error, bool) { ...... }</code></pre>


### [func (MainConfig) GithubLink](../../../../../src/github.com/slowfei/gosfdoc/config.go#L197-L241) <a name="f_func__MainConfig__GithubLink_relMDPath_string_isToMarkdown_bool__string"><a/> [↩](#p_func__MainConfig__GithubLink_relMDPath_string_isToMarkdown_bool__string) | [#](#f_func__MainConfig__GithubLink_relMDPath_string_isToMarkdown_bool__string)
> to github.com link path<br/>
> use on a tag href<br/>
>                                                      append path       relative path<br/>
> e.g.: https://.../project/doc/v1_0_0/md/default/(github.com/slowfei)/(temp/gosfdoc.md)<br/>
>   to: https://.../project/doc/v1_0_0/src/github.com/slowfei/gosfdoc.go (to source code path)<br/>
>   to: https://.../project/doc/v1_0_0/md/default/github.com/test/test.md  (to markdown path)<br/>
> @param `relMDPath` relative markdown out project path.<br/>
>                    relative path: $GOPATH/[github.com/slowfei]/projectname/( .../markdown.md )<br/>
> @param `isToMarkdown` to markdown link? false is source code access path<br/>
> @return use github.com to relative link. "../../../" or "../../src/[projectname]"<br/>
> <br/>


<pre><code class='go custom'>func (MainConfig) GithubLink(relMDPath string, isToMarkdown bool) string { ...... }</code></pre>


### [type MenuFile struct](../../../../../src/github.com/slowfei/gosfdoc/config.go#L257-L261) <a name="f_type_MenuFile_struct"><a/> [↩](#p_type_MenuFile_struct) | [#](#f_type_MenuFile_struct)
> html menu show helper struct<br/>
> src.html File list struct<br/>
> <br/>


<pre><code class='go custom'>type MenuFile struct {
	MenuName string
	Version  string
	List     []FileLink
}</code></pre>


### [type MenuMarkdown struct](../../../../../src/github.com/slowfei/gosfdoc/config.go#L247-L251) <a name="f_type_MenuMarkdown_struct"><a/> [↩](#p_type_MenuMarkdown_struct) | [#](#f_type_MenuMarkdown_struct)
> html menu show helper struct<br/>
> index.html Markdown struct<br/>
> <br/>


<pre><code class='go custom'>type MenuMarkdown struct {
	MenuName string
	Version  string
	List     []PackageInfo
}</code></pre>


### [type OperateResult int](../../../../../src/github.com/slowfei/gosfdoc/gosfdoc.go#L109-L109) <a name="f_type_OperateResult_int"><a/> [↩](#p_type_OperateResult_int) | [#](#f_type_OperateResult_int)
> operate result<br/>
> <br/>


<pre><code class='go custom'>type OperateResult int</code></pre>


### [type PackageInfo struct](../../../../../src/github.com/slowfei/gosfdoc/struct.go#L566-L571) <a name="f_type_PackageInfo_struct"><a/> [↩](#p_type_PackageInfo_struct) | [#](#f_type_PackageInfo_struct)
> package info<br/>
> <br/>


<pre><code class='go custom'>type PackageInfo struct {
	menuName string `json:"-"` // type belongs
	Name     string // package name plain text
	Desc     string // description plain text
	Link     string // markdown link
}</code></pre>


### [type Preview struct](../../../../../src/github.com/slowfei/gosfdoc/struct.go#L592-L598) <a name="f_type_Preview_struct"><a/> [↩](#p_type_Preview_struct) | [#](#f_type_Preview_struct)
> preview struct info<br/>
> <br/>


<pre><code class='go custom'>type Preview struct {
	SortTag  string // sort tag
	Level    int    // hierarchy level show. 0 is >, 1 is >>, 3 is >>> ...(markdown syntax)
	ShowText string // show plain text
	Anchor   string // preferably unique, with the func link
	DescText string // markdown brief description or implement objects, can empty.
}</code></pre>


### [type SortSet struct](../../../../../src/github.com/slowfei/gosfdoc/struct.go#L618-L622) <a name="f_type_SortSet_struct"><a/> [↩](#p_type_SortSet_struct) | [#](#f_type_SortSet_struct)
> Preview,CodeBlock,Document sort implement<br/>
> <br/>


<pre><code class='go custom'>type SortSet struct {
	previews   []Preview
	documents  []Document
	codeBlocks []CodeBlock
}</code></pre>


### [func (SortSet) Len](../../../../../src/github.com/slowfei/gosfdoc/struct.go#L627-L639) <a name="f_func__SortSet__Len___int"><a/> [↩](#p_func__SortSet__Len___int) | [#](#f_func__SortSet__Len___int)
> sort Len() implement<br/>
> <br/>


<pre><code class='go custom'>func (SortSet) Len() int { ...... }</code></pre>


### [func (SortSet) Less](../../../../../src/github.com/slowfei/gosfdoc/struct.go#L644-L656) <a name="f_func__SortSet__Less_i_j_int__bool"><a/> [↩](#p_func__SortSet__Less_i_j_int__bool) | [#](#f_func__SortSet__Less_i_j_int__bool)
> sort Less(...) implement<br/>
> <br/>


<pre><code class='go custom'>func (SortSet) Less(i, j int) bool { ...... }</code></pre>


### [func (SortSet) Swap](../../../../../src/github.com/slowfei/gosfdoc/struct.go#L661-L671) <a name="f_func__SortSet__Swap_i_j_int__"><a/> [↩](#p_func__SortSet__Swap_i_j_int__) | [#](#f_func__SortSet__Swap_i_j_int__)
> sort Swap(...) implement<br/>
> <br/>


<pre><code class='go custom'>func (SortSet) Swap(i, j int)  { ...... }</code></pre>


