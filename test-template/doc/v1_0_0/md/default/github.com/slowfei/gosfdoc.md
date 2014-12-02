
## Welcome to gosfdoc v1.0.0
------

&emsp;&emsp;gosfdoc 目的是规范代码文档生成，在写代码同时也可以进行文档的编写，不必等写完代码后再来编写文档，可以避免写文档的懒惰。
将编写代码的点点滴滴都可以记录与文档中，让大脑有了大致的思路然后快速的编写代码，也可以让以后快速的了解代码。

目前支持生成的语言：

> 1.golang

生成的内容：


## 文档生成规范
------

* 文档介绍注释（在整改项目只能是唯一的）

```
/**Intro

 	统一的格式，注意 /**Intro （这里为截取的内容） */
 	在首页上显示，指引用户下载或查看文档。

	支持markdown格式编写。
 */

```

* 关于页面注释(注意文档右上角的About)

```
/**About

 	统一的格式，注意 /**About （这里为截取的内容） */

## 关于
------

Hello world gosfdoc

author:
> slowfei#foxmail.com

links:[http://www.slowfei.com][0]

[0]:http://www.slowfei.com

 */
```

* 关于页面注释(注意文档右上角的About)




## Preview
------
> [func Main()][#]<br/>
> [type TestStruct struct][#]<br/>
> implements: [Test][#]<br/>
>>[func (* TestStruct) hello(str string) string](#func_TestStruct.hello)<a name="preview_TestStruct.hello"><a/><br/>
>>[func (* TestStruct) hello2() string][#]<br/>
>>

###Package files
[a.go](../../../../src/gosfdoc.go#L10-L16) [b.go](../../../../../../../gosfdoc.go#L10) [c.go](../../../../../../../github.com/slowfei/gosfdoc.go#L10)

## Constants
------

<!-- <a href="#HTTP 状态编码参数" id="HTTP 状态编码参数">HTTP 状态编码参数</a> -->

###[HTTP 状态编码参数](#HTTP 状态编码参 数) 
> 状态编码参数的描述内容

```go
const (
	StatusOK                    = 200
	StatusCreated               = 201
	StatusAccepted              = 202
	StatusNonAuthoritativeInfo  = 203
	StatusNonContent			= 204
	StatusResetContent  		= 205
	StatusPartialContent		= 206
)
```

> [src][#] HTTP status codes, defined in RFC 2616.

```go
const DefaultMaxHeaderBytes = 1 << 20 // 1 MB
```

## Variables
------

###[错误参数定义](../../../src) [#](#错误参数定义)
> 错误参数的其他描述(Errors introduced by the HTTP server.)

```go
var (
    ErrWriteAfterFlush = errors.New("Conn.Write called after Flush")
    ErrBodyNotAllowed  = errors.New("http: onquest method or response status code does not allow body")
    ErrHijackedon      = errors.New("Conn has been hijacked")
    ErrContentLength   = errors.New("Conn.Write wrote more than the declared Content-Length")
)
```

> DefaultClient is the default Client and is used by Get, Head, and Post.

```go
var DefaultClient = &Client{}
```

###[可定义标题，也可以不定义标题](#可定义标题，也可以不定义标题) [#](#可定义标题，也可以不定义标题)
> DefaultServeMux is the default ServeMux used by Serve.

```go
var DefaultServeMux = NewServeMux()
```

## Func Details 
------

### [type TestStruct struct][#]
>	结构介绍
>

```go
type TestStruct struct{
	host        string
	routerKeyst []string
	routers     map[string]IRouter
}
```

###[func (*TestStruct) hello](src.html?f=github.com/slowfei/gosfdoc.go) <a name="func_TestStruct.hello"><a/> [↩](#preview_TestStruct.hello)|[#](#func_TestStruct.hello)
> 函数介绍描述<br/>
> <br/>
> @param `str` 字符串传递<br/>
> @return `v1` 返回参数v1<br/>
> @return v2 返回参数v2<br/>

```go
func (* TestStruct) hello(str string) (v1,v2 string)
```


### [type Config struct](../../../../github.com/slowfei/gosfdoc.go)
> leafveingo config<br/>
> default see _defaultConfigJson

```go
type Config struct {
    AppVersion          string   // app version.
    FileUploadSize      int64    // file size upload
    Charset             string   // html encode type
    StaticFileSuffixes  []string // supported static file suffixes
    ServerTimeout       int64    // server time out, default 0
    SessionMaxlifeTime  int32    // http session maxlife time, unit second. use session set
    IPHeaderKey         string   // proxy to http headers set ip key, default ""
    IsReqPathIgnoreCase bool     // request url path ignore case
    MultiProjectHosts   []string // setting integrated multi-project hosts,default nil

    TemplateSuffix string // template suffix
    IsCompactHTML  bool   // is Compact HTML, default true

    LogConfigPath string // log config path, relative or absolute path, relative path from execute file root directory. log config path, relative or absolute path
    LogGroup      string // log group name

    TLSCertPath string // tls cert.pem, relative or absolute path, relative path from execute file root directory
    TLSKeyPath  string // tls key.pem
    TLSPort     int    // tls run prot, default server port+1
    TLSAloneRun bool   // are leafveingo alone run tls server, default false

    // is ResponseWriter writer compress gizp...
    // According Accept-Encoding select compress type
    // default true
    IsRespWriteCompress bool

    UserData map[string]string // user custom config other info
}
```


[#]:javascript:;




