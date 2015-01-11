#### Test2
1

2

3

4

5

6

7

8

9

10

11

12

13

14
15
### [type Config struct](../../../../github.com/slowfei/gosfdoc.go)<a name="type-config-struct"><a/>
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
