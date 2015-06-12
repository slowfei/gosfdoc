
## Preview
------

> [const ( GO_NAME... ...GO_TEST_SUFFIX )](#f_const___GO_NAME---_---GO_TEST_SUFFIX__)<a name="p_const___GO_NAME---_---GO_TEST_SUFFIX__"><a/>

> [var ( REXFunc... ...SNBetweens )](#f_var___REXFunc---_---SNBetweens__)<a name="p_var___REXFunc---_---SNBetweens__"><a/>

> [type GolangParser struct](#f_type_GolangParser_struct)<a name="p_type_GolangParser_struct"><a/>

>> [func NewParser() \*GolangParser](#f_func_NewParser___+GolangParser)<a name="p_func_NewParser___+GolangParser"><a/>

>> [func (\*GolangParser) CheckFile(filePath string, info os.FileInfo) bool](#f_func__+GolangParser__CheckFile_filePath_string_info_os-FileInfo__bool)<a name="p_func__+GolangParser__CheckFile_filePath_string_info_os-FileInfo__bool"><a/>

>> [func (\*GolangParser) EachIndexFile(filebuf \*gosfdoc.FileBuf) ](#f_func__+GolangParser__EachIndexFile_filebuf_+gosfdoc-FileBuf__)<a name="p_func__+GolangParser__EachIndexFile_filebuf_+gosfdoc-FileBuf__"><a/>

>> [func (\*GolangParser) Name() string](#f_func__+GolangParser__Name___string)<a name="p_func__+GolangParser__Name___string"><a/>

>> [func (\*GolangParser) ParseCodeblock(filebuf \*gosfdoc.FileBuf) []gosfdoc.CodeBlock](#f_func__+GolangParser__ParseCodeblock_filebuf_+gosfdoc-FileBuf____gosfdoc-CodeBlock)<a name="p_func__+GolangParser__ParseCodeblock_filebuf_+gosfdoc-FileBuf____gosfdoc-CodeBlock"><a/>

>> [func (\*GolangParser) ParseEnd() ](#f_func__+GolangParser__ParseEnd___)<a name="p_func__+GolangParser__ParseEnd___"><a/>

>> [func (\*GolangParser) ParsePackageInfo(filebuf \*gosfdoc.FileBuf) string](#f_func__+GolangParser__ParsePackageInfo_filebuf_+gosfdoc-FileBuf__string)<a name="p_func__+GolangParser__ParsePackageInfo_filebuf_+gosfdoc-FileBuf__string"><a/>

>> [func (\*GolangParser) ParsePreview(filebuf \*gosfdoc.FileBuf) []gosfdoc.Preview](#f_func__+GolangParser__ParsePreview_filebuf_+gosfdoc-FileBuf____gosfdoc-Preview)<a name="p_func__+GolangParser__ParsePreview_filebuf_+gosfdoc-FileBuf____gosfdoc-Preview"><a/>

>> [func (\*GolangParser) ParseStart(config gosfdoc.MainConfig) ](#f_func__+GolangParser__ParseStart_config_gosfdoc-MainConfig__)<a name="p_func__+GolangParser__ParseStart_config_gosfdoc-MainConfig__"><a/>

> [type Temp struct](#f_type_Temp_struct)<a name="p_type_Temp_struct"><a/>

<br/>
### Directory files
[golang.go ](../../src/github.com/slowfei/gosfdoc/lang/golang/golang.go)

## Constants
------
### [source code](../../src/github.com/slowfei/gosfdoc/lang/golang/golang.go#L26-L30) <a name="f_const___GO_NAME---_---GO_TEST_SUFFIX__"><a/> [↩](#p_const___GO_NAME---_---GO_TEST_SUFFIX__) | [#](#f_const___GO_NAME---_---GO_TEST_SUFFIX__)

<pre><code class='go custom'>const (
	GO_NAME        = "go"
	GO_SUFFIX      = ".go"
	GO_TEST_SUFFIX = "_test.go"
)</code></pre>


## Variables
------
### [source code](../../src/github.com/slowfei/gosfdoc/lang/golang/golang.go#L32-L62) <a name="f_var___REXFunc---_---SNBetweens__"><a/> [↩](#p_var___REXFunc---_---SNBetweens__) | [#](#f_var___REXFunc---_---SNBetweens__)

<pre><code class='go custom'>var (
	// e.g.: func (t type)funcname(params) return val{;
	// https://www.debuggex.com/r/Su6Ns1LhVxpfD_Di
	// [0-1:prototype] [2-3:comment or null] [4-5:func type or null]
	// [6-7:func name] [8-9:func params] [10-11:single return value or null]
	// [12-13:multi return value or null] [14-15:"{"]
	REXFunc = regexp.MustCompile(`(/\*\*[\s]*(?:[ ]+.*?\n)+[ ]*\*/[ ]*\n|(?:(?:[ ]*//.*?\n)+))?func[ ]*(?:\(([\w \*\n\r\.\[\]]*)\))?[ \n]*([A-Z]\w*)[ \n]*(?:\(([\w ,\*\n\r\.\{\}\[\]]*)\))+?[ \n]*(?:([\w\.\*\{\}\[\]]*)|(?:\(([\w ,\*\n\r\.\{\}\[\]]*)\)))?[ \n]*({)`)
	// e.g.: type Temp struct {; [0-1:prototype][2-3:comment][4-5:type define name][6-7:type name][8-9:"{"]
	REXType = regexp.MustCompile(`(/\*\*[\s]*(?:[ ]+.*?\n)+[ ]*\*/[ ]*\n|(?:(?:[ ]*//.*?\n)+))?type[ ]+([A-Z]\w*)[ ]+(\w+)[ ]*(\{)?`)
	// e.g.: package main
	REXPackage = regexp.MustCompile(`package (\w+)\s*`)
	// e.g.: /** ... */[\n]package main; //...[\n]package main
	REXPackageInfo = regexp.MustCompile(`(/\*\*[\s]*(?:[ ]+.*?\n)+[ ]*\*/[ ]*\n|(?:(?:[ ]*//.*?\n)+))[ ]*package \w+`)
	// e.g.: /** ... */[\n]const|var TConst = 1; //...[\n]const (
	// https://www.debuggex.com/r/2qeKD9vwnBjkgORT
	// [0-1:prototype] [2-3:comment or null] [4-5:const|var] [6-7: define name]
	REXDefine = regexp.MustCompile(`(/\*\*[\s]*(?:[ ]+.*?\n)+[ ]*\*/[ ]*\n|(?:(?:[ ]*//.*?\n)+))?[ ]*(const|var)\s+(?:\(|(?:([A-Z]\w*)\s*=.+))`)
	// e.g: rows data
	REXRows = regexp.MustCompile("\\w+.*")

	SNRoundBrackets = SFSubUtil.NewSubNest([]byte("("), []byte(")"))
	SNBraces        = SFSubUtil.NewSubNest([]byte("{"), []byte("}"))
	SNBetweens      = []*SFSubUtil.SubNest{
		SFSubUtil.NewSubNest([]byte(`"`), []byte(`"`)),
		SFSubUtil.NewSubNest([]byte(`'`), []byte(`'`)),
		SFSubUtil.NewSubNest([]byte("`"), []byte("`")),
		SFSubUtil.NewSubNest([]byte("/*"), []byte("*/")),
		SFSubUtil.NewSubNotNest([]byte("//"), []byte("\n")),
		SNBraces,
	}
)</code></pre>


## Func Details
------
### [type GolangParser struct](../../src/github.com/slowfei/gosfdoc/lang/golang/golang.go#L138-L141) <a name="f_type_GolangParser_struct"><a/> [↩](#p_type_GolangParser_struct) | [#](#f_type_GolangParser_struct)
> golang parser<br/>
> <br/>


<pre><code class='go custom'>type GolangParser struct {
	config  gosfdoc.MainConfig
	indexDB index.IndexDB
}</code></pre>


### [func NewParser](../../src/github.com/slowfei/gosfdoc/lang/golang/golang.go#L146-L150) <a name="f_func_NewParser___+GolangParser"><a/> [↩](#p_func_NewParser___+GolangParser) | [#](#f_func_NewParser___+GolangParser)
> new golang parser<br/>
> <br/>


<pre><code class='go custom'>func NewParser() *GolangParser { ...... }</code></pre>


### [func (\*GolangParser) CheckFile](../../src/github.com/slowfei/gosfdoc/lang/golang/golang.go#L176-L187) <a name="f_func__+GolangParser__CheckFile_filePath_string_info_os-FileInfo__bool"><a/> [↩](#p_func__+GolangParser__CheckFile_filePath_string_info_os-FileInfo__bool) | [#](#f_func__+GolangParser__CheckFile_filePath_string_info_os-FileInfo__bool)
> see DocParser interface<br/>
> <br/>


<pre><code class='go custom'>func (*GolangParser) CheckFile(filePath string, info os.FileInfo) bool { ...... }</code></pre>


### [func (\*GolangParser) EachIndexFile](../../src/github.com/slowfei/gosfdoc/lang/golang/golang.go#L192-L270) <a name="f_func__+GolangParser__EachIndexFile_filebuf_+gosfdoc-FileBuf__"><a/> [↩](#p_func__+GolangParser__EachIndexFile_filebuf_+gosfdoc-FileBuf__) | [#](#f_func__+GolangParser__EachIndexFile_filebuf_+gosfdoc-FileBuf__)
> see DocParser interface<br/>
> <br/>


<pre><code class='go custom'>func (*GolangParser) EachIndexFile(filebuf *gosfdoc.FileBuf)  { ...... }</code></pre>


### [func (\*GolangParser) Name](../../src/github.com/slowfei/gosfdoc/lang/golang/golang.go#L154-L156) <a name="f_func__+GolangParser__Name___string"><a/> [↩](#p_func__+GolangParser__Name___string) | [#](#f_func__+GolangParser__Name___string)

<pre><code class='go custom'>func (*GolangParser) Name() string { ...... }</code></pre>


### [func (\*GolangParser) ParseCodeblock](../../src/github.com/slowfei/gosfdoc/lang/golang/golang.go#L581-L761) <a name="f_func__+GolangParser__ParseCodeblock_filebuf_+gosfdoc-FileBuf____gosfdoc-CodeBlock"><a/> [↩](#p_func__+GolangParser__ParseCodeblock_filebuf_+gosfdoc-FileBuf____gosfdoc-CodeBlock) | [#](#f_func__+GolangParser__ParseCodeblock_filebuf_+gosfdoc-FileBuf____gosfdoc-CodeBlock)
> see DocParser interface<br/>
> <br/>


<pre><code class='go custom'>func (*GolangParser) ParseCodeblock(filebuf *gosfdoc.FileBuf) []gosfdoc.CodeBlock { ...... }</code></pre>


### [func (\*GolangParser) ParseEnd](../../src/github.com/slowfei/gosfdoc/lang/golang/golang.go#L169-L171) <a name="f_func__+GolangParser__ParseEnd___"><a/> [↩](#p_func__+GolangParser__ParseEnd___) | [#](#f_func__+GolangParser__ParseEnd___)
> see DocParser interface<br/>
> <br/>


<pre><code class='go custom'>func (*GolangParser) ParseEnd()  { ...... }</code></pre>


### [func (\*GolangParser) ParsePackageInfo](../../src/github.com/slowfei/gosfdoc/lang/golang/golang.go#L766-L813) <a name="f_func__+GolangParser__ParsePackageInfo_filebuf_+gosfdoc-FileBuf__string"><a/> [↩](#p_func__+GolangParser__ParsePackageInfo_filebuf_+gosfdoc-FileBuf__string) | [#](#f_func__+GolangParser__ParsePackageInfo_filebuf_+gosfdoc-FileBuf__string)
> see DocParser interface<br/>
> <br/>


<pre><code class='go custom'>func (*GolangParser) ParsePackageInfo(filebuf *gosfdoc.FileBuf) string { ...... }</code></pre>


### [func (\*GolangParser) ParsePreview](../../src/github.com/slowfei/gosfdoc/lang/golang/golang.go#L507-L576) <a name="f_func__+GolangParser__ParsePreview_filebuf_+gosfdoc-FileBuf____gosfdoc-Preview"><a/> [↩](#p_func__+GolangParser__ParsePreview_filebuf_+gosfdoc-FileBuf____gosfdoc-Preview) | [#](#f_func__+GolangParser__ParsePreview_filebuf_+gosfdoc-FileBuf____gosfdoc-Preview)
> see DocParser interface<br/>
> <br/>


<pre><code class='go custom'>func (*GolangParser) ParsePreview(filebuf *gosfdoc.FileBuf) []gosfdoc.Preview { ...... }</code></pre>


### [func (\*GolangParser) ParseStart](../../src/github.com/slowfei/gosfdoc/lang/golang/golang.go#L161-L164) <a name="f_func__+GolangParser__ParseStart_config_gosfdoc-MainConfig__"><a/> [↩](#p_func__+GolangParser__ParseStart_config_gosfdoc-MainConfig__) | [#](#f_func__+GolangParser__ParseStart_config_gosfdoc-MainConfig__)
> see DocParser interface<br/>
> <br/>


<pre><code class='go custom'>func (*GolangParser) ParseStart(config gosfdoc.MainConfig)  { ...... }</code></pre>


### [type Temp struct](../../src/github.com/slowfei/gosfdoc/lang/golang/golang.go#L1053-L1053) <a name="f_type_Temp_struct"><a/> [↩](#p_type_Temp_struct) | [#](#f_type_Temp_struct)

<pre><code class='go custom'>type Temp struct</code></pre>


