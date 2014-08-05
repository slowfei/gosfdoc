
# Welcome to gosfdoc
------

测试锚点<a name="pookie"></a>

&emsp;&emsp;gosfdoc 目的是规范代码文档生成，在写代码同时也可以进行文档的编写，不必等写完代码后再来编写文档，可以避免写文档的懒惰。
将编写代码的点点滴滴都可以记录与文档中，让大脑有了大致的思路然后快速的编写代码，也可以让以后快速的了解代码。

目前支持生成的语言：

> 1.golang

生成的内容：


## 文档生成规范
------

* 文档介绍注释（在整改项目只能是唯一的）

```gosfdoc

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
> [func Main()](#funcMain)<br/>
> [type TestStruct struct][#]<br/>
>> implements: [Test][#]<br/>
>>[func (* TestStruct) hello(str string) string](#func_TestStruct.hello)<a name="preview_TestStruct.hello"></a><br/>
>>[func (* TestStruct) hello2() string][#]<br/>
>>

###Package files
[a.go][#] [b.go][#] [c.go][#]

## Constants
------

<!-- <a href="#HTTP 状态编码参数" id="HTTP 状态编码参数">HTTP 状态编码参数</a> -->

###[HTTP 状态编码参数](#HTTP 状态编码参 数) 
> 状态编码参数的描述内容

```
const (
    StatusOK                   = 200
    StatusCreated              = 201
    StatusAccepted             = 202
    StatusNonAuthoritativeInfo = 203
    StatusNoContent            = 204
    StatusResetContent         = 205
    StatusPartialContent       = 206
)
```

> [src][#] HTTP status codes, defined in RFC 2616.

```
const DefaultMaxHeaderBytes = 1 << 20 // 1 MB
```

## Variables
------

###[错误参数定义](#错误参数定义)
> 错误参数的其他描述(Errors introduced by the HTTP server.)

```
var (
    ErrWriteAfterFlush = errors.New("Conn.Write called after Flush")
    ErrBodyNotAllowed  = errors.New("http: request method or response status code does not allow body")
    ErrHijacked        = errors.New("Conn has been hijacked")
    ErrContentLength   = errors.New("Conn.Write wrote more than the declared Content-Length")
)
```

> DefaultClient is the default Client and is used by Get, Head, and Post.

```
var DefaultClient = &Client{}
```

###[可定义标题，也可以不定义标题](#可定义标题，也可以不定义标题)
> DefaultServeMux is the default ServeMux used by Serve.

```
var DefaultServeMux = NewServeMux()
```

## Func Details 
------

### [type TestStruct struct][#]
>	结构介绍
>

```
type TestStruct struct{

}
```


###[func (*TestStruct) hello](../../../src/github.com/slowfei/gosfdoc.go) [↩](#preview_TestStruct.hello)|[#](#func_TestStruct.hello)<a name="func_TestStruct.hello"></a>
> 函数介绍描述<br/>
> <br/>
> @param `str` 字符串传递<br/>
> @return `v1` 返回参数v1<br/>
> @return v2 返回参数v2<br/>

```
func (* TestStruct) hello(str string) (v1,v2 string)
```



[#]:javascript:;




