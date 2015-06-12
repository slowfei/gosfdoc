
## Preview
------

> [const ( DBTypeMemory... ...DBTypeFile )](#f_const___DBTypeMemory---_---DBTypeFile__)<a name="p_const___DBTypeMemory---_---DBTypeFile__"><a/>

> [var ( ErrInvalidIndex... ...ErrInvalidIndex )](#f_var___ErrInvalidIndex---_---ErrInvalidIndex__)<a name="p_var___ErrInvalidIndex---_---ErrInvalidIndex__"><a/>

> [type DBType int](#f_type_DBType_int)<a name="p_type_DBType_int"><a/>

> [type IndexDB interface](#f_type_IndexDB_interface)<a name="p_type_IndexDB_interface"><a/>

>> [func CreateIndexDB(langName string, dbT DBType) IndexDB](#f_func_CreateIndexDB_langName_string_dbT_DBType__IndexDB)<a name="p_func_CreateIndexDB_langName_string_dbT_DBType__IndexDB"><a/>

> [type TypeInfo struct](#f_type_TypeInfo_struct)<a name="p_type_TypeInfo_struct"><a/>

>> [func (\*mendb) Close() ](#f_func__+mendb__Close___)<a name="p_func__+mendb__Close___"><a/>

>> [func (\*mendb) Open() error](#f_func__+mendb__Open___error)<a name="p_func__+mendb__Open___error"><a/>

>> [func (\*mendb) SetType(t TypeInfo) error](#f_func__+mendb__SetType_t_TypeInfo__error)<a name="p_func__+mendb__SetType_t_TypeInfo__error"><a/>

>> [func (\*mendb) Type(packageName, packagePath, typeName string) (TypeInfo, bool)](#f_func__+mendb__Type_packageName_packagePath_typeName_string___TypeInfo_bool_)<a name="p_func__+mendb__Type_packageName_packagePath_typeName_string___TypeInfo_bool_"><a/>

<br/>
### Directory files
[index.go ](../../../../../src/github.com/slowfei/gosfdoc/index/index.go)[memory.go ](../../../../../src/github.com/slowfei/gosfdoc/index/memory.go)

## Constants
------
### [source code](../../../../../src/github.com/slowfei/gosfdoc/index/index.go#L23-L26) <a name="f_const___DBTypeMemory---_---DBTypeFile__"><a/> [↩](#p_const___DBTypeMemory---_---DBTypeFile__) | [#](#f_const___DBTypeMemory---_---DBTypeFile__)
> database type definit<br/>
> <br/>


<pre><code class='go custom'>const (
	DBTypeMemory DBType = iota // disposable memory cache
	DBTypeFile                 // file type storage
)</code></pre>


## Variables
------
### [source code](../../../../../src/github.com/slowfei/gosfdoc/index/index.go#L29-L31) <a name="f_var___ErrInvalidIndex---_---ErrInvalidIndex__"><a/> [↩](#p_var___ErrInvalidIndex---_---ErrInvalidIndex__) | [#](#f_var___ErrInvalidIndex---_---ErrInvalidIndex__)
> error definit<br/>
> <br/>


<pre><code class='go custom'>var (
	ErrInvalidIndex = errors.New("gosfdoc/index: Invalid unique index, package or type name nil.")
)</code></pre>


## Func Details
------
### [type DBType int](../../../../../src/github.com/slowfei/gosfdoc/index/index.go#L20-L20) <a name="f_type_DBType_int"><a/> [↩](#p_type_DBType_int) | [#](#f_type_DBType_int)
> database type<br/>
> <br/>


<pre><code class='go custom'>type DBType int</code></pre>


### [type IndexDB interface](../../../../../src/github.com/slowfei/gosfdoc/index/index.go#L36-L72) <a name="f_type_IndexDB_interface"><a/> [↩](#p_type_IndexDB_interface) | [#](#f_type_IndexDB_interface)
> data storage interface<br/>
> <br/>


<pre><code class='go custom'>type IndexDB interface {

	/**
	 *  Open (operating data)-> Close -> Open (operating data)-> Close...
	 *
	 *  @return `error`
	 */
	Open() error

	/**
	 *  all finished operating data can close
	 */
	Close()

	/**
	 *  save as type info, the same data is overwritten
	 *  package and name identifies a unique index
	 *
	 *  @param `t`
	 *  @return `error`
	 */
	SetType(t TypeInfo) error

	/**
	 *	package path or package name or type name query
	 *
	 *	包名和包路径可二选一，两个参数一起可以提高精确率。
	 *	package path and package name can choose one, can also together.
	 *
	 *	@param `packageName`
	 *	@param `packagePath`
	 *	@param `typeName`
	 *	@return `TypeInfo`
	 *	@return `bool`
	 */
	Type(packageName, packagePath, typeName string) (TypeInfo, bool)
}</code></pre>


### [func CreateIndexDB](../../../../../src/github.com/slowfei/gosfdoc/index/index.go#L80-L83) <a name="f_func_CreateIndexDB_langName_string_dbT_DBType__IndexDB"><a/> [↩](#p_func_CreateIndexDB_langName_string_dbT_DBType__IndexDB) | [#](#f_func_CreateIndexDB_langName_string_dbT_DBType__IndexDB)
> open or create IndexDB<br/>
> @param `langName` language name string<br/>
> @param `dbT`      storage type, default DBTypeMemory<br/>
> <br/>


<pre><code class='go custom'>func CreateIndexDB(langName string, dbT DBType) IndexDB { ...... }</code></pre>


### [type TypeInfo struct](../../../../../src/github.com/slowfei/gosfdoc/index/index.go#L89-L96) <a name="f_type_TypeInfo_struct"><a/> [↩](#p_type_TypeInfo_struct) | [#](#f_type_TypeInfo_struct)
> index type struct<br/>
> type info in various languages<br/>
> <br/>


<pre><code class='go custom'>type TypeInfo struct {
	DocHttpUrl  string // document http url e.g.: http://slowfei.github.io/gosfdoc/index.html
	PackageName string // package name unique index
	PackagePath string // helper package name the set path
	TypeName    string // only within the scope of the package name unique index
	LineStart   int    // line number start
	LineEnd     int    // line number end
}</code></pre>


### [func (\*mendb) Close](../../../../../src/github.com/slowfei/gosfdoc/index/memory.go#L33-L35) <a name="f_func__+mendb__Close___"><a/> [↩](#p_func__+mendb__Close___) | [#](#f_func__+mendb__Close___)

<pre><code class='go custom'>func (*mendb) Close()  { ...... }</code></pre>


### [func (\*mendb) Open](../../../../../src/github.com/slowfei/gosfdoc/index/memory.go#L29-L31) <a name="f_func__+mendb__Open___error"><a/> [↩](#p_func__+mendb__Open___error) | [#](#f_func__+mendb__Open___error)

<pre><code class='go custom'>func (*mendb) Open() error { ...... }</code></pre>


### [func (\*mendb) SetType](../../../../../src/github.com/slowfei/gosfdoc/index/memory.go#L37-L46) <a name="f_func__+mendb__SetType_t_TypeInfo__error"><a/> [↩](#p_func__+mendb__SetType_t_TypeInfo__error) | [#](#f_func__+mendb__SetType_t_TypeInfo__error)

<pre><code class='go custom'>func (*mendb) SetType(t TypeInfo) error { ...... }</code></pre>


### [func (\*mendb) Type](../../../../../src/github.com/slowfei/gosfdoc/index/memory.go#L48-L82) <a name="f_func__+mendb__Type_packageName_packagePath_typeName_string___TypeInfo_bool_"><a/> [↩](#p_func__+mendb__Type_packageName_packagePath_typeName_string___TypeInfo_bool_) | [#](#f_func__+mendb__Type_packageName_packagePath_typeName_string___TypeInfo_bool_)

<pre><code class='go custom'>func (*mendb) Type(packageName, packagePath, typeName string) (TypeInfo, bool) { ...... }</code></pre>


