
## golang语言标签
------
　　这里将介绍golang语言标签的基本使用，golang提供多种方式进行文档的描述，可以使用&#47;&#42;&#42; &#42;&#47;，也可以使用&#47;&#47;。


## package标签
------
　　在每个.go文件中都可以进行描述，多个描述将会以";"进行分割。

> &#47;&#47; 请用简短的一句话描述当前包的作用，最好是一行的内容，否则也会被去除换行。<br/>
> package main

　

> &#47;&#42;&#42;<br/>
> &nbsp;&#42;&nbsp;&nbsp;为使用的灵活性提供两种格式使用。<br/>
> &nbsp;&#42;&#47;<br/>
> package main


## func标签
------

> &#47;&#42;&#42;<br/>
> &nbsp;&#42;&nbsp;&nbsp;函数的描述，可以使用markdown格式，这里描述内容的格式没有特殊要求，可根据要求定义。<br/>
> &nbsp;&#42;&nbsp;&nbsp;（注：需要紧跟着函数体）<br/>
> &nbsp;&#42;&nbsp;&nbsp;@param &#96;src&#96; 描述参数的作用，<br/>
> &nbsp;&#42;&nbsp;&nbsp;@return 返回值的意思描述<br/>
> &nbsp;&#42;&#47;<br/>
> func(t &#42;Test) funcName(src string) string{<br/>
> &nbsp;&nbsp;...<br/>
> }<br/>

　

> &#47;&#47;&nbsp;&nbsp;函数的描述，可以使用markdown格式，这里描述内容的格式没有特殊要求，可根据要求定义。<br/>
> &#47;&#47;&nbsp;&nbsp;（注：需要紧跟着函数体）<br/>
> &#47;&#47;&nbsp;&nbsp;@param &#96;src&#96; 描述参数的作用，<br/>
> &#47;&#47;&nbsp;&nbsp;@return 返回值的意思描述<br/>
> func(t &#42;Test) funcName(param string) string{<br/>
> &nbsp;&nbsp;...<br/>
> }<br/>


## type标签
------

> &#47;&#42;&#42;<br/>
> &nbsp;&#42;&nbsp;&nbsp;类型的描述，可以使用markdown格式，这里描述内容的格式没有特殊要求，可根据要求定义。<br/>
> &nbsp;&#42;&nbsp;&nbsp;（注：需要紧跟类型定义）<br/>
> &nbsp;&#42;&#47;<br/>
> type Test interface{<br/>
> &nbsp;&nbsp;...<br/>
> }<br/>

　

> &#47;&#47;&nbsp;&nbsp;函数的描述，可以使用markdown格式，这里描述内容的格式没有特殊要求，可根据要求定义。<br/>
> &#47;&#47;&nbsp;&nbsp;（注：需要紧跟类型定义）<br/>
> type Test interface{<br/>
> &nbsp;&nbsp;...<br/>
> }<br/>


## var & const 标签
------

> &#47;&#42;&#42;<br/>
> &nbsp;&#42;&nbsp;&nbsp;var描述，可以使用markdown格式，这里描述内容的格式没有特殊要求，可根据要求定义。<br/>
> &nbsp;&#42;&nbsp;&nbsp;（注：需要紧跟var定义）<br/>
> &nbsp;&#42;&#47;<br/>
> var (<br/>
> &nbsp;&nbsp;&#47;&#47;&nbsp;&nbsp;注：全部定义的参数为首字母大写才会被输出到文档中<br/>
> &nbsp;&nbsp;TestVar = "temp"<br/>
> &nbsp;&nbsp;...<br/>
> )<br/>
> <br/>
> &#47;&#42;&#42;<br/>
> &nbsp;&#42;&nbsp;&nbsp;可以两种形式进行定义与描述<br/>
> &nbsp;&#42;&nbsp;&nbsp;（注：需要紧跟var定义）<br/>
> &nbsp;&#42;&#47;<br/>
> var TestVar = "temp"<br/>

　

> &#47;&#47;&nbsp;&nbsp;const描述，可以使用markdown格式，这里描述内容的格式没有特殊要求，可根据要求定义。<br/>
> &#47;&#47;&nbsp;&nbsp;（注：需要紧跟const定义）<br/>
> const (<br/>
> &nbsp;&nbsp;&#47;&#47;&nbsp;&nbsp;注：全部定义的参数为首字母大写才会被输出到文档中<br/>
> &nbsp;&nbsp;TestConst = "temp"<br/>
> &nbsp;&nbsp;...<br/>
> )<br/>
> <br/>
> &#47;&#47;&nbsp;&nbsp;可以两种形式进行定义与描述<br/>
> &#47;&#47;&nbsp;&nbsp;（注：需要紧跟const定义）<br/>
> const TestConst = "temp"<br/>




