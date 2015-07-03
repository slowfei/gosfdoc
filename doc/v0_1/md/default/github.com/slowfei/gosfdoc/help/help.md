
## gosfdoc安装
------
使用终端go命令

##### 1.下载gosfdoc和web server使用的模版包

    go get github.com/slowfei/gosfdoc
    go get github.com/slowfei/leafveingo/template

##### 2.安装可执行文件

windows:

    go build -o="gosfdoc.exe" github.com/slowfei/gosfdoc/bin

linux or unix:

    go build -o="gosfdoc" github.com/slowfei/gosfdoc/bin

go build -o="指定的文件名或绝对路径文件"<br/>
执行文件可以放入自定义的环境变量中，方便直接使用。

##### 涉及包

    github.com/slowfei/gosfcore
    github.com/slowfei/leafveingo/template


## 基本使用
------
使用终端执行命令<br/>
（注：gosfdoc默认放置在环境变量中，演示直接使用命令，如果未设置环境变量则需要加上文件全路径）

    $ gosfdoc -h

    gosfdoc v0.1.000

    usage help:
    'create' command init create default gosfdoc.json file, can be custom to modify file content.
    'output' command by gosfdoc.json output document 
    'web [8080]' run web server can specify port, default 8080, command by gosfdoc.json 

    other param:
      -config="gosfdoc.json": custom config file path.
      -lang="": ["go"] specify code language type ',' separated, default all language.
      -v="": output the document version string.


##### `$ gosfdoc create`
　　创建配置文件命令，在命令行输入`cd`转到需要需要生成项目的目录中，然后再执行 `$ gosfdoc create` 此命令可选参数`-lang`（根据指定语言生成文档，目前只支持go语言），执行后会在当前目录中生成默认的配置文件(gosfdoc.json)，查看配置文件详细说明。

##### `$ gosfdoc output -v="1.0"`
　　生成文档命令，必须指定生成的版本号，版本号出现重复会提示是否覆盖。此命令可选参数`-config="gosfdoc.json"`（指定配置文件），默认会寻找当前目录的配置文件，也可以指定绝对路径的配置文件。

##### `$ gosfdoc web 8080`
　　运行web服务，可在浏览器中浏览项目文档，`8080`为端口号，可根据需求进行更改。此命令可选参数`-config="gosfdoc.json"`，默认会寻找当前目录的配置文件，也可以指定绝对路径的配置文件。


## gosfdoc.json配置说明
------

    {
        // 项目扫描路径，根据路径下的文件生成文档。该路径可以是绝对路径或则相对路径，相对路径是根据配置文件的目录连接计算。
        // "/" 表示从配置文件所在目录开始扫描，建议使用绝对路径或"/"。 使用create 命令生成出来的配置文件是绝对路径。
        "ScanPath"         : "/",

        // 指定生成文档的语言，例如多语言混合编程的时候，就可以指定需要生成语言的文档。（注：目前只支持golang）
        "CodeLang"         : ["go"],

        // 指定文档生成的目录路径，路径可以是绝对路径或相对路径，相对路径是根据ScanPath路径连接。默认"doc"
        "Outpath"          : "doc",

        // 追加的前缀路径，在使用go的同学可能都知道使用github时在项目工作目录下还需要建立"github.com/slowfei/[projectname]"
        // 然而在生成文档时操作路径是在项目的根目录操作的，OutAppendPath的设置就是为此准备。如果不了解的同学可以先试试使用空""。
        "OutAppendPath"    : "github.com/slowfei/gosfdoc",

        // 在生成文档时是否将源代码复制一份，默认true。如果false将不复制源代码，此时在浏览文档中就无法浏览源代码。
        // 需要设置CodeLinkRoot为true，浏览文档时直接浏览项目目录文件的源代码。
        "CopyCode"         : true,

        // CodeLinkRoot与CopyCode是辅助作用的，如果设置true浏览源代码的路径是项目的根目录，则CopyCode可以设置false。
        // 此时在浏览项目文档时，查看源代码的路径就会以项目根目录为主。如果设置false，CopyCode最好设置为true，否则在浏览
        // 项目文档时无法浏览源代码。
        "CodeLinkRoot"     : false,

        // 生成的文档是以html形式展现的，这里可以设置html的标题显示。
        "HtmlTitle"        : "gosfdoc document",

        // 在文档页面中有显示文档标题，具体效果可以运行文档查看。
        "DocTitle"         : "<b>Github: <a target='_blank' href='https://github.com/slowfei/gosfdoc'>https://github.com/slowfei/gosfdoc</a></b>",

        // 浏览文档页面中左侧导航的标题显示
        "MenuTitle"        : "<center><b>markdown</b></center>",

        // 文档支持的本地化语言选项，gosfdoc可以支持本地化多语言浏览文档，但需要的是自行翻译。本地话语言使用可以查看"本地化语言支持"
        // 设置例子： 
        //  [
        //      {"default" : "Default"}
        //      {"zh-cn" : "简体中文"}
        //      {"en-us" : "English"}
        //  ],
        //
        "Languages"        : [
                                {"default" : "Default"}
                             ],

        // 需要过滤不扫描的子目录文件路径，使用的是相对路径。例如："test-template"是一个测试静态页面的目录，所以不需要生成文档。
        "FilterPaths"      : [
                                "assets",
                                "test-template"
                             ] 
    }


## 公用文档标签
------

#### 目录路径简介信息
在首页显示目录路径的简介信息

> &nbsp;&#35;package-info 简短的介绍当前目录路径下的概要介绍，这标签基本在全是".dc"文档格式的目录中使用。

<br/>
#### 项目内容介绍
在首页主窗口显示项目介绍的内容

> &#47;&#42;&#42;Intro<br/>
> &nbsp;&#42;&nbsp;&nbsp;以markdown的格式编写项目介绍内容<br/>
> &nbsp;&#42;&nbsp;&nbsp;<br/>
> &nbsp;&#42;&nbsp;&nbsp;&#35;&#35;&#35;&#35;&#35;标题<br/>
> &nbsp;&#42;&nbsp;&nbsp;其他相关的内容<br/>
> &nbsp;&#42;&#47;<br/>
> <br/>
> &#47;&#47;Intro<br/>
> &#47;&#47;&nbsp;&nbsp;为使用的灵活性提供两种格式，此格式也可以在源代码中进行编写。<br/>
> &#47;&#47;&nbsp;&nbsp;此格式需要注意的是头和结尾的约束。<br/>
> &#47;&#47;End<br/>

<br/>
#### “关于”信息的介绍
注意到右上角的关于吗？下面介绍的标签格式就是“关于”信息内容介绍的编写

> &#47;&#42;&#42;About<br/>
> &nbsp;&#42;&nbsp;&nbsp;以markdown的格式编写介绍内容<br/>
> &nbsp;&#42;&nbsp;&nbsp;<br/>
> &nbsp;&#42;&nbsp;&nbsp;&#35;&#35;&#35;&#35;&#35;标题<br/>
> &nbsp;&#42;&nbsp;&nbsp;其他相关的内容<br/>
> &nbsp;&#42;&#47;<br/>
> <br/>
> &#47;&#47;About<br/>
> &#47;&#47;&nbsp;&nbsp;为使用的灵活性提供两种格式，此格式也可以在源代码中进行编写。<br/>
> &#47;&#47;&nbsp;&nbsp;此格式需要注意的是头和结尾的约束。<br/>
> &#47;&#47;End<br/>

<br/>
#### 文档内容标签

> &#47;&#42;&#42;&#42;1-标题(注意三个星号)<br/>
> &nbsp;&#42;&nbsp;&nbsp;以markdown的格式编写文档内容<br/>
> &nbsp;&#42;&nbsp;&nbsp;<br/>
> &nbsp;&#42;&nbsp;&nbsp;&#35;&#35;&#35;&#35;&#35;[index]-[标题内容]<br/>
> &nbsp;&#42;&nbsp;&nbsp;index使用整数表示排序，然后紧跟着"-"，最后自定义标题。<br/>
> &nbsp;&#42;&nbsp;&nbsp;这个标题会在左上角的标题导航中显示。<br/>
> &nbsp;&#42;&#47;<br/>
> <br/>
> &#47;&#47;&#47;2-标题(注意三个斜杠)<br/>
> &#47;&#47;&nbsp;&nbsp;为使用的灵活性提供两种格式，此格式也可以在源代码中进行编写。<br/>
> &#47;&#47;&nbsp;&nbsp;此格式需要注意的是头和结尾的约束。<br/>
> &#47;&#47;End<br/>

<br/>
#### 私有化标签
私有化标签是为了提供在生成文档时，将需要隐藏的隐私内容不输出文档中。私有化标签基本在源代码中使用，可以指定源代码文件只输出文档或只输出源代码。

> 该标签放置在源代码的第首行<br/>
> &#47;&#47;&#35;private-doc-code &#47;&#47;代码和文档都不输出<br/>
> &#47;&#47;&#35;private-code &#47;&#47;只输出源代码，gosfdoc.json > CopyCode = tue的时候<br/>
> &#47;&#47;&#35;private-doc &#47;&#47;只输出文档，不会输出源代码<br/>
> <br/>
> 该标签可防止在源代码中的任意位置<br/>
> &#47;&#47;&#35;private<br/>
> &nbsp;&nbsp;func main(){...} &#47;&#47;需要标记的代码块，此代码块不会进行输出。<br/>
> &nbsp;&nbsp;(注：CopyCode = false && CodeLinkRoot = true 浏览源代码时行数会出现错误，主要是因为此标签隐藏的代码块中的行数不计算在内。)<br/>
> &#47;&#47;&#35;private-end<br/>



## 本地化语言支持
------
　　gosfdoc以为多语言浏览提供方便的操作方式，需要注意的一点，gosfdoc不会自定进行本地化语言的翻译，需要根据项目管理者的要求自行翻译。

#### 那该如何操作呢？

> 第一步

修改`gosfdoc.json`配置文件中的`Languages`参数，根据需要添加需要的相应语言选项。default也可根据需要进行修饰，default是默认的输出语言。<br/>
接下来以建立English为例子:
```
"Languages" : [
                 {"default" : "Default"}
                 {"en-us" : "English"}
              ],
``` 
> 第二步

建立"en-us"文件夹，观察输出的文档目录，比如现在生成的是0.1版本目录名称会是"v0_1"。接下来进入的路径为：`../project/doc/v0_1/md`，然后复制default目录一份，包括里面的所有文件，将复制的一份改名为"en-us"。

> 第三步

在en-us目录下的所有文件，只要关于本地化的语言进行翻译就可以了。此时重新加载浏览器就会出现相应的本地化语言选项。


