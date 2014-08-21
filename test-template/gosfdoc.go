//  Copyright 2013 slowfei And The Contributors All rights reserved.
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.
//
//  Create on 2013-9-14
//  Update on 2014-07-17
//  Email  slowfei#foxmail.com
//  Home   http://www.slowfei.com

//
//	router manager
//
package leafveingo

import (
	// "github.com/slowfei/gosfcore/debug"
	"fmt"
	"github.com/slowfei/gosfcore/utils/strings"
	"net/http"
	"reflect"
)

var (
	//	globalRouterList
	_globalRouterList []globalRouter = nil
)

//#pragma mark global handle ----------------------------------------------------------------------------------------------------

//
//	global router
//
type globalRouter struct {
	appName string
	router  IRouter
}

/**
 *	global add router
 *
 *	@param appName
 *	@param routerKey
 *	@param router
 */
func AddRouter(appName string, router IRouter) {
	_globalRouterList = append(_globalRouterList, globalRouter{appName, router})
}

//#pragma mark interface option ----------------------------------------------------------------------------------------------------

//
//	router element
//
type RouterElement struct {
	host       string
	routerKeys []string
	routers    map[string]IRouter
}

//
//	router option
//
type RouterOption struct {
	/*
		GET http://localhost:8080/home/index.go
		routerKey 		= "/home/"
		routerPath 		= "index"
		requestMethod	= "get"
		urlSuffix       = "go"

		GET http://localhost:8080/home/manager/help.go
		routerKey 		= "/home/"
		routerPath 		= "manager/help"
		requestMethod	= "get"
		urlSuffix       = "go"

		GET http://localhost:8080/home/manager/other
		routerKey 		= "/home/"
		routerPath 		= "manager/other"
		requestMethod	= "get"
		urlSuffix       = ""

		GET http://localhost:8080/home#!other
		routerKey 		= "/home"
		routerPath 		= "#!other"
		requestMethod	= "get"
		urlSuffix       = ""
	*/

	RouterData       interface{}   // custom storage data
	RouterDataRefVal reflect.Value // reflect value data

	RouterKey     string //
	RouterPath    string // have been converted to lowercase
	RequestMethod string // lowercase get post ...
	UrlSuffix     string // lowercase

	AppName string // application name
}

//
//	router interface
//
type IRouter interface {

	/**
	 *	1. after router parse
	 *
	 *	@param context
	 *	@param option
	 *	@return HttpStatus  200 pass
	 */
	AfterRouterParse(context *HttpContext, option *RouterOption) HttpStatus

	/**
	 *	2. parse func name
	 *
	 *	@param context			http context
	 *	@param option			router option
	 *	@return funcName 		function name specifies call
	 *	@return statusCode		http status code, 200 pass, other to status page
	 */
	ParseFuncName(context *HttpContext, option *RouterOption) (funcName string, statusCode HttpStatus, err error)

	/**
	 *	3. before call func
	 *
	 *	@param context
	 *	@param option
	 *	@return HttpStatus
	 */
	CallFuncBefore(context *HttpContext, option *RouterOption) HttpStatus

	/**
	 *	4. request func
	 *
	 *	@param context			http content
	 *	@param funcName			call controller func name
	 *	@param option			router option
	 *	@return returnValue		controller func return value
	 *	@return statusCode		http status code, 200 pass, other to status page
	 */
	CallFunc(context *HttpContext, funcName string, option *RouterOption) (returnValue interface{}, statusCode HttpStatus, err error)

	/**
	 *	5. parse template path
	 *	no need to add the suffix
	 *
	 *	@param context
	 *	@param funcName	 controller call func name
	 *	@param option	 router option
	 *	@return template path, suggest "[routerKey]/[funcName]"
	 */
	ParseTemplatePath(context *HttpContext, funcName string, option *RouterOption) string

	/**
	 *	6. after call func
	 *
	 *	@param context
	 *	@param option
	 */
	CallFuncAfter(context *HttpContext, option *RouterOption)

	/**
	 *	@return router key
	 */
	RouterKey() string

	/**
	 *	@return controller option
	 */
	ControllerOption() ControllerOption

	/**
	 *	@return controller info
	 */
	Info() string
}

//#pragma mark Leafveingo method ----------------------------------------------------------------------------------------------------

/**
 *	router parse
 *
 *	@param context
 *	@param reqPathNoSuffix 		the no suffix request path
 *	@param reqSuffix			url request suffix
 */
func routerParse(context *HttpContext, reqPathNoSuffix, reqSuffix string) (router IRouter, option *RouterOption, statusCode HttpStatus) {
    /*  routerParse*/
	statusCode = Status404

	var lowerReqPath string
	if context.lvServer.IsReqPathIgnoreCase() {
		lowerReqPath = SFStringsUtil.ToLower(reqPathNoSuffix)
	} else {
		lowerReqPath = reqPathNoSuffix
	}

	reqPathLen := len(lowerReqPath)

	// FOR routerList -> IF host -> IF scheme -> FOR keys -> IF reqPath

	listCount := len(context.lvServer.routerList)
	for i := 0; i < listCount; i++ {
		element := context.lvServer.routerList[i]

		if context.reqHost == element.host {

			keyCount := len(element.routerKeys)
			for j := 0; j < keyCount; j++ {

				key := element.routerKeys[j]
				keyLen := len(key)

				if keyLen <= reqPathLen && key == lowerReqPath[:keyLen] {

					//	controller router find
					if iR, ok := element.routers[key]; ok {

						//	uri scheme handle
						statusCode = routerSchemeHandle(context, iR.ControllerOption().Scheme())
						if statusCode != Status200 {
							router = nil
							option = nil
							return
						}

						context.routerElement = element
						router = iR
						option = new(RouterOption)
						option.AppName = context.lvServer.AppName()

						if 0 != len(reqSuffix) {
							if '.' == reqSuffix[0] {
								reqSuffix = reqSuffix[1:]
							}
							option.UrlSuffix = SFStringsUtil.ToLower(reqSuffix)
						}

						option.RequestMethod = SFStringsUtil.ToLower(context.Request.Method)
						if 0 == len(option.RequestMethod) {
							option.RequestMethod = "get"
						}

						option.RouterKey = key
						option.RouterPath = lowerReqPath[keyLen:]

						statusCode = router.AfterRouterParse(context, option)
					} else {
						//	基本上不会进来此处
						context.lvServer.log.Error("lv.routerKeys contains %#v and lv.routers not contains %#v", key, key)
						statusCode = Status404
					}
					return

				} // end  keyLen <= reqPathLen && key == lowerReqPath[:keyLen]

			} // end j := 0; j < keyCount; j++

		} // end context.reqHost == element.host

	} // end i := 0; i < listCount; i++

	return
}

/**
 *	router uri scheme handle
 *
 *	@param context
 *	@param scheme
 *	@return statusCode Status200 pass
 */
func routerSchemeHandle(context *HttpContext, scheme URIScheme) (statusCode HttpStatus) {
	statusCode = Status400

	switch context.RequestScheme() {
	case URI_SCHEME_HTTP:

		if scheme&URI_SCHEME_HTTP == URI_SCHEME_HTTP {
			statusCode = Status200
			break
		}

		//	TODO 下一步考虑委托函数处理不支持的协议，可让用户来处理。

		if scheme&URI_SCHEME_HTTPS == URI_SCHEME_HTTPS {
			if "GET" == context.Request.Method {
				//	重定向操作
				statusCode = StatusNil
				routerSchemeRedirect(context, URI_SCHEME_HTTPS)
			}
		}

	case URI_SCHEME_HTTPS:
		if scheme&URI_SCHEME_HTTPS == URI_SCHEME_HTTPS {
			statusCode = Status200
			break
		}

		if scheme&URI_SCHEME_HTTP == URI_SCHEME_HTTP {
			if "GET" == context.Request.Method {
				//	重定向操作
				statusCode = StatusNil
				routerSchemeRedirect(context, URI_SCHEME_HTTP)
			}
		}

	default:

	}

	return
}

/**
 *  does not support the scheme redirect handle
 *
 *  @param context
 *  @param scheme
 */
func routerSchemeRedirect(context *HttpContext, scheme URIScheme) {
	schemeStr := ""
	port := 0

	switch scheme {
	case URI_SCHEME_HTTP:
		schemeStr = "http://"
		port = context.lvServer.port
	case URI_SCHEME_HTTPS:
		if nil == context.lvServer.tlsListener {
			return
		}
		schemeStr = "https://"
		port = context.lvServer.tlsPort
	default:
		return
	}

	//	host
	host := context.Request.Host
	hostLen := len(host)

	//	remove port
	retScope := hostLen - 7
	if 0 > retScope {
		retScope = 0
	}
	index := -1
	for i := hostLen - 1; i >= retScope; i-- {
		if ':' == host[i] {
			index = i
			break
		}
	}
	if -1 != index {
		host = fmt.Sprintf("%s:%d", host[:index], port)
	}

	//	path
	path := context.Request.URL.Path
	if 0 == len(path) {
		path = "/"
	} else if '/' != path[0] {
		path = "/" + path
	}

	query := context.Request.URL.Query().Encode()
	if 0 != len(query) {
		query = "?" + query
	}

	urlstr := schemeStr + host + path + query

	http.Redirect(context.RespWrite, context.Request, urlstr, int(Status302))

	context.LVServer().Log().Debug("scheme redirect to: " + urlstr)
}

//  Copyright 2013 slowfei And The Contributors All rights reserved.
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.
//
//  Create on 2013-9-14
//  Update on 2014-07-17
//  Email  slowfei#foxmail.com
//  Home   http://www.slowfei.com

//
//	router manager
//
package leafveingo

import (
	// "github.com/slowfei/gosfcore/debug"
	"fmt"
	"github.com/slowfei/gosfcore/utils/strings"
	"net/http"
	"reflect"
)

var (
	//	globalRouterList
	_globalRouterList []globalRouter = nil
)

//#pragma mark global handle ----------------------------------------------------------------------------------------------------

//
//	global router
//
type globalRouter struct {
	appName string
	router  IRouter
}

/**
 *	global add router
 *
 *	@param appName
 *	@param routerKey
 *	@param router
 */
func AddRouter(appName string, router IRouter) {
	_globalRouterList = append(_globalRouterList, globalRouter{appName, router})
}

//#pragma mark interface option ----------------------------------------------------------------------------------------------------

//
//	router element
//
type RouterElement struct {
	host       string
	routerKeys []string
	routers    map[string]IRouter
}

//
//	router option
//
type RouterOption struct {
	/*
		GET http://localhost:8080/home/index.go
		routerKey 		= "/home/"
		routerPath 		= "index"
		requestMethod	= "get"
		urlSuffix       = "go"

		GET http://localhost:8080/home/manager/help.go
		routerKey 		= "/home/"
		routerPath 		= "manager/help"
		requestMethod	= "get"
		urlSuffix       = "go"

		GET http://localhost:8080/home/manager/other
		routerKey 		= "/home/"
		routerPath 		= "manager/other"
		requestMethod	= "get"
		urlSuffix       = ""

		GET http://localhost:8080/home#!other
		routerKey 		= "/home"
		routerPath 		= "#!other"
		requestMethod	= "get"
		urlSuffix       = ""
	*/

	RouterData       interface{}   // custom storage data
	RouterDataRefVal reflect.Value // reflect value data

	RouterKey     string //
	RouterPath    string // have been converted to lowercase
	RequestMethod string // lowercase get post ...
	UrlSuffix     string // lowercase

	AppName string // application name
}

//
//	router interface
//
type IRouter interface {

	/**
	 *	1. after router parse
	 *
	 *	@param context
	 *	@param option
	 *	@return HttpStatus  200 pass
	 */
	AfterRouterParse(context *HttpContext, option *RouterOption) HttpStatus

	/**
	 *	2. parse func name
	 *
	 *	@param context			http context
	 *	@param option			router option
	 *	@return funcName 		function name specifies call
	 *	@return statusCode		http status code, 200 pass, other to status page
	 */
	ParseFuncName(context *HttpContext, option *RouterOption) (funcName string, statusCode HttpStatus, err error)

	/**
	 *	3. before call func
	 *
	 *	@param context
	 *	@param option
	 *	@return HttpStatus
	 */
	CallFuncBefore(context *HttpContext, option *RouterOption) HttpStatus

	/**
	 *	4. request func
	 *
	 *	@param context			http content
	 *	@param funcName			call controller func name
	 *	@param option			router option
	 *	@return returnValue		controller func return value
	 *	@return statusCode		http status code, 200 pass, other to status page
	 */
	CallFunc(context *HttpContext, funcName string, option *RouterOption) (returnValue interface{}, statusCode HttpStatus, err error)

	/**
	 *	5. parse template path
	 *	no need to add the suffix
	 *
	 *	@param context
	 *	@param funcName	 controller call func name
	 *	@param option	 router option
	 *	@return template path, suggest "[routerKey]/[funcName]"
	 */
	ParseTemplatePath(context *HttpContext, funcName string, option *RouterOption) string

	/**
	 *	6. after call func
	 *
	 *	@param context
	 *	@param option
	 */
	CallFuncAfter(context *HttpContext, option *RouterOption)

	/**
	 *	@return router key
	 */
	RouterKey() string

	/**
	 *	@return controller option
	 */
	ControllerOption() ControllerOption

	/**
	 *	@return controller info
	 */
	Info() string
}

//#pragma mark Leafveingo method ----------------------------------------------------------------------------------------------------

/**
 *	router parse
 *
 *	@param context
 *	@param reqPathNoSuffix 		the no suffix request path
 *	@param reqSuffix			url request suffix
 */
func routerParse(context *HttpContext, reqPathNoSuffix, reqSuffix string) (router IRouter, option *RouterOption, statusCode HttpStatus) {
    /*  routerParse*/
	statusCode = Status404

	var lowerReqPath string
	if context.lvServer.IsReqPathIgnoreCase() {
		lowerReqPath = SFStringsUtil.ToLower(reqPathNoSuffix)
	} else {
		lowerReqPath = reqPathNoSuffix
	}

	reqPathLen := len(lowerReqPath)

	// FOR routerList -> IF host -> IF scheme -> FOR keys -> IF reqPath

	listCount := len(context.lvServer.routerList)
	for i := 0; i < listCount; i++ {
		element := context.lvServer.routerList[i]

		if context.reqHost == element.host {

			keyCount := len(element.routerKeys)
			for j := 0; j < keyCount; j++ {

				key := element.routerKeys[j]
				keyLen := len(key)

				if keyLen <= reqPathLen && key == lowerReqPath[:keyLen] {

					//	controller router find
					if iR, ok := element.routers[key]; ok {

						//	uri scheme handle
						statusCode = routerSchemeHandle(context, iR.ControllerOption().Scheme())
						if statusCode != Status200 {
							router = nil
							option = nil
							return
						}

						context.routerElement = element
						router = iR
						option = new(RouterOption)
						option.AppName = context.lvServer.AppName()

						if 0 != len(reqSuffix) {
							if '.' == reqSuffix[0] {
								reqSuffix = reqSuffix[1:]
							}
							option.UrlSuffix = SFStringsUtil.ToLower(reqSuffix)
						}

						option.RequestMethod = SFStringsUtil.ToLower(context.Request.Method)
						if 0 == len(option.RequestMethod) {
							option.RequestMethod = "get"
						}

						option.RouterKey = key
						option.RouterPath = lowerReqPath[keyLen:]

						statusCode = router.AfterRouterParse(context, option)
					} else {
						//	基本上不会进来此处
						context.lvServer.log.Error("lv.routerKeys contains %#v and lv.routers not contains %#v", key, key)
						statusCode = Status404
					}
					return

				} // end  keyLen <= reqPathLen && key == lowerReqPath[:keyLen]

			} // end j := 0; j < keyCount; j++

		} // end context.reqHost == element.host

	} // end i := 0; i < listCount; i++

	return
}

/**
 *	router uri scheme handle
 *
 *	@param context
 *	@param scheme
 *	@return statusCode Status200 pass
 */
func routerSchemeHandle(context *HttpContext, scheme URIScheme) (statusCode HttpStatus) {
	statusCode = Status400

	switch context.RequestScheme() {
	case URI_SCHEME_HTTP:

		if scheme&URI_SCHEME_HTTP == URI_SCHEME_HTTP {
			statusCode = Status200
			break
		}

		//	TODO 下一步考虑委托函数处理不支持的协议，可让用户来处理。

		if scheme&URI_SCHEME_HTTPS == URI_SCHEME_HTTPS {
			if "GET" == context.Request.Method {
				//	重定向操作
				statusCode = StatusNil
				routerSchemeRedirect(context, URI_SCHEME_HTTPS)
			}
		}

	case URI_SCHEME_HTTPS:
		if scheme&URI_SCHEME_HTTPS == URI_SCHEME_HTTPS {
			statusCode = Status200
			break
		}

		if scheme&URI_SCHEME_HTTP == URI_SCHEME_HTTP {
			if "GET" == context.Request.Method {
				//	重定向操作
				statusCode = StatusNil
				routerSchemeRedirect(context, URI_SCHEME_HTTP)
			}
		}

	default:

	}

	return
}

/**
 *  does not support the scheme redirect handle
 *
 *  @param context
 *  @param scheme
 */
func routerSchemeRedirect(context *HttpContext, scheme URIScheme) {
	schemeStr := ""
	port := 0

	switch scheme {
	case URI_SCHEME_HTTP:
		schemeStr = "http://"
		port = context.lvServer.port
	case URI_SCHEME_HTTPS:
		if nil == context.lvServer.tlsListener {
			return
		}
		schemeStr = "https://"
		port = context.lvServer.tlsPort
	default:
		return
	}

	//	host
	host := context.Request.Host
	hostLen := len(host)

	//	remove port
	retScope := hostLen - 7
	if 0 > retScope {
		retScope = 0
	}
	index := -1
	for i := hostLen - 1; i >= retScope; i-- {
		if ':' == host[i] {
			index = i
			break
		}
	}
	if -1 != index {
		host = fmt.Sprintf("%s:%d", host[:index], port)
	}

	//	path
	path := context.Request.URL.Path
	if 0 == len(path) {
		path = "/"
	} else if '/' != path[0] {
		path = "/" + path
	}

	query := context.Request.URL.Query().Encode()
	if 0 != len(query) {
		query = "?" + query
	}

	urlstr := schemeStr + host + path + query

	http.Redirect(context.RespWrite, context.Request, urlstr, int(Status302))

	context.LVServer().Log().Debug("scheme redirect to: " + urlstr)
}

//  Copyright 2013 slowfei And The Contributors All rights reserved.
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.
//
//  Create on 2013-9-14
//  Update on 2014-07-17
//  Email  slowfei#foxmail.com
//  Home   http://www.slowfei.com

//
//	router manager
//
package leafveingo

import (
	// "github.com/slowfei/gosfcore/debug"
	"fmt"
	"github.com/slowfei/gosfcore/utils/strings"
	"net/http"
	"reflect"
)

var (
	//	globalRouterList
	_globalRouterList []globalRouter = nil
)

//#pragma mark global handle ----------------------------------------------------------------------------------------------------

//
//	global router
//
type globalRouter struct {
	appName string
	router  IRouter
}

/**
 *	global add router
 *
 *	@param appName
 *	@param routerKey
 *	@param router
 */
func AddRouter(appName string, router IRouter) {
	_globalRouterList = append(_globalRouterList, globalRouter{appName, router})
}

//#pragma mark interface option ----------------------------------------------------------------------------------------------------

//
//	router element
//
type RouterElement struct {
	host       string
	routerKeys []string
	routers    map[string]IRouter
}

//
//	router option
//
type RouterOption struct {
	/*
		GET http://localhost:8080/home/index.go
		routerKey 		= "/home/"
		routerPath 		= "index"
		requestMethod	= "get"
		urlSuffix       = "go"

		GET http://localhost:8080/home/manager/help.go
		routerKey 		= "/home/"
		routerPath 		= "manager/help"
		requestMethod	= "get"
		urlSuffix       = "go"

		GET http://localhost:8080/home/manager/other
		routerKey 		= "/home/"
		routerPath 		= "manager/other"
		requestMethod	= "get"
		urlSuffix       = ""

		GET http://localhost:8080/home#!other
		routerKey 		= "/home"
		routerPath 		= "#!other"
		requestMethod	= "get"
		urlSuffix       = ""
	*/

	RouterData       interface{}   // custom storage data
	RouterDataRefVal reflect.Value // reflect value data

	RouterKey     string //
	RouterPath    string // have been converted to lowercase
	RequestMethod string // lowercase get post ...
	UrlSuffix     string // lowercase

	AppName string // application name
}

//
//	router interface
//
type IRouter interface {

	/**
	 *	1. after router parse
	 *
	 *	@param context
	 *	@param option
	 *	@return HttpStatus  200 pass
	 */
	AfterRouterParse(context *HttpContext, option *RouterOption) HttpStatus

	/**
	 *	2. parse func name
	 *
	 *	@param context			http context
	 *	@param option			router option
	 *	@return funcName 		function name specifies call
	 *	@return statusCode		http status code, 200 pass, other to status page
	 */
	ParseFuncName(context *HttpContext, option *RouterOption) (funcName string, statusCode HttpStatus, err error)

	/**
	 *	3. before call func
	 *
	 *	@param context
	 *	@param option
	 *	@return HttpStatus
	 */
	CallFuncBefore(context *HttpContext, option *RouterOption) HttpStatus

	/**
	 *	4. request func
	 *
	 *	@param context			http content
	 *	@param funcName			call controller func name
	 *	@param option			router option
	 *	@return returnValue		controller func return value
	 *	@return statusCode		http status code, 200 pass, other to status page
	 */
	CallFunc(context *HttpContext, funcName string, option *RouterOption) (returnValue interface{}, statusCode HttpStatus, err error)

	/**
	 *	5. parse template path
	 *	no need to add the suffix
	 *
	 *	@param context
	 *	@param funcName	 controller call func name
	 *	@param option	 router option
	 *	@return template path, suggest "[routerKey]/[funcName]"
	 */
	ParseTemplatePath(context *HttpContext, funcName string, option *RouterOption) string

	/**
	 *	6. after call func
	 *
	 *	@param context
	 *	@param option
	 */
	CallFuncAfter(context *HttpContext, option *RouterOption)

	/**
	 *	@return router key
	 */
	RouterKey() string

	/**
	 *	@return controller option
	 */
	ControllerOption() ControllerOption

	/**
	 *	@return controller info
	 */
	Info() string
}

//#pragma mark Leafveingo method ----------------------------------------------------------------------------------------------------

/**
 *	router parse
 *
 *	@param context
 *	@param reqPathNoSuffix 		the no suffix request path
 *	@param reqSuffix			url request suffix
 */
func routerParse(context *HttpContext, reqPathNoSuffix, reqSuffix string) (router IRouter, option *RouterOption, statusCode HttpStatus) {
    /*  routerParse*/
	statusCode = Status404

	var lowerReqPath string
	if context.lvServer.IsReqPathIgnoreCase() {
		lowerReqPath = SFStringsUtil.ToLower(reqPathNoSuffix)
	} else {
		lowerReqPath = reqPathNoSuffix
	}

	reqPathLen := len(lowerReqPath)

	// FOR routerList -> IF host -> IF scheme -> FOR keys -> IF reqPath

	listCount := len(context.lvServer.routerList)
	for i := 0; i < listCount; i++ {
		element := context.lvServer.routerList[i]

		if context.reqHost == element.host {

			keyCount := len(element.routerKeys)
			for j := 0; j < keyCount; j++ {

				key := element.routerKeys[j]
				keyLen := len(key)

				if keyLen <= reqPathLen && key == lowerReqPath[:keyLen] {

					//	controller router find
					if iR, ok := element.routers[key]; ok {

						//	uri scheme handle
						statusCode = routerSchemeHandle(context, iR.ControllerOption().Scheme())
						if statusCode != Status200 {
							router = nil
							option = nil
							return
						}

						context.routerElement = element
						router = iR
						option = new(RouterOption)
						option.AppName = context.lvServer.AppName()

						if 0 != len(reqSuffix) {
							if '.' == reqSuffix[0] {
								reqSuffix = reqSuffix[1:]
							}
							option.UrlSuffix = SFStringsUtil.ToLower(reqSuffix)
						}

						option.RequestMethod = SFStringsUtil.ToLower(context.Request.Method)
						if 0 == len(option.RequestMethod) {
							option.RequestMethod = "get"
						}

						option.RouterKey = key
						option.RouterPath = lowerReqPath[keyLen:]

						statusCode = router.AfterRouterParse(context, option)
					} else {
						//	基本上不会进来此处
						context.lvServer.log.Error("lv.routerKeys contains %#v and lv.routers not contains %#v", key, key)
						statusCode = Status404
					}
					return

				} // end  keyLen <= reqPathLen && key == lowerReqPath[:keyLen]

			} // end j := 0; j < keyCount; j++

		} // end context.reqHost == element.host

	} // end i := 0; i < listCount; i++

	return
}

/**
 *	router uri scheme handle
 *
 *	@param context
 *	@param scheme
 *	@return statusCode Status200 pass
 */
func routerSchemeHandle(context *HttpContext, scheme URIScheme) (statusCode HttpStatus) {
	statusCode = Status400

	switch context.RequestScheme() {
	case URI_SCHEME_HTTP:

		if scheme&URI_SCHEME_HTTP == URI_SCHEME_HTTP {
			statusCode = Status200
			break
		}

		//	TODO 下一步考虑委托函数处理不支持的协议，可让用户来处理。

		if scheme&URI_SCHEME_HTTPS == URI_SCHEME_HTTPS {
			if "GET" == context.Request.Method {
				//	重定向操作
				statusCode = StatusNil
				routerSchemeRedirect(context, URI_SCHEME_HTTPS)
			}
		}

	case URI_SCHEME_HTTPS:
		if scheme&URI_SCHEME_HTTPS == URI_SCHEME_HTTPS {
			statusCode = Status200
			break
		}

		if scheme&URI_SCHEME_HTTP == URI_SCHEME_HTTP {
			if "GET" == context.Request.Method {
				//	重定向操作
				statusCode = StatusNil
				routerSchemeRedirect(context, URI_SCHEME_HTTP)
			}
		}

	default:

	}

	return
}

/**
 *  does not support the scheme redirect handle
 *
 *  @param context
 *  @param scheme
 */
func routerSchemeRedirect(context *HttpContext, scheme URIScheme) {
	schemeStr := ""
	port := 0

	switch scheme {
	case URI_SCHEME_HTTP:
		schemeStr = "http://"
		port = context.lvServer.port
	case URI_SCHEME_HTTPS:
		if nil == context.lvServer.tlsListener {
			return
		}
		schemeStr = "https://"
		port = context.lvServer.tlsPort
	default:
		return
	}

	//	host
	host := context.Request.Host
	hostLen := len(host)

	//	remove port
	retScope := hostLen - 7
	if 0 > retScope {
		retScope = 0
	}
	index := -1
	for i := hostLen - 1; i >= retScope; i-- {
		if ':' == host[i] {
			index = i
			break
		}
	}
	if -1 != index {
		host = fmt.Sprintf("%s:%d", host[:index], port)
	}

	//	path
	path := context.Request.URL.Path
	if 0 == len(path) {
		path = "/"
	} else if '/' != path[0] {
		path = "/" + path
	}

	query := context.Request.URL.Query().Encode()
	if 0 != len(query) {
		query = "?" + query
	}

	urlstr := schemeStr + host + path + query

	http.Redirect(context.RespWrite, context.Request, urlstr, int(Status302))

	context.LVServer().Log().Debug("scheme redirect to: " + urlstr)
}

//  Copyright 2013 slowfei And The Contributors All rights reserved.
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.
//
//  Create on 2013-9-14
//  Update on 2014-07-17
//  Email  slowfei#foxmail.com
//  Home   http://www.slowfei.com

//
//	router manager
//
package leafveingo

import (
	// "github.com/slowfei/gosfcore/debug"
	"fmt"
	"github.com/slowfei/gosfcore/utils/strings"
	"net/http"
	"reflect"
)

var (
	//	globalRouterList
	_globalRouterList []globalRouter = nil
)

//#pragma mark global handle ----------------------------------------------------------------------------------------------------

//
//	global router
//
type globalRouter struct {
	appName string
	router  IRouter
}

/**
 *	global add router
 *
 *	@param appName
 *	@param routerKey
 *	@param router
 */
func AddRouter(appName string, router IRouter) {
	_globalRouterList = append(_globalRouterList, globalRouter{appName, router})
}

//#pragma mark interface option ----------------------------------------------------------------------------------------------------

//
//	router element
//
type RouterElement struct {
	host       string
	routerKeys []string
	routers    map[string]IRouter
}

//
//	router option
//
type RouterOption struct {
	/*
		GET http://localhost:8080/home/index.go
		routerKey 		= "/home/"
		routerPath 		= "index"
		requestMethod	= "get"
		urlSuffix       = "go"

		GET http://localhost:8080/home/manager/help.go
		routerKey 		= "/home/"
		routerPath 		= "manager/help"
		requestMethod	= "get"
		urlSuffix       = "go"

		GET http://localhost:8080/home/manager/other
		routerKey 		= "/home/"
		routerPath 		= "manager/other"
		requestMethod	= "get"
		urlSuffix       = ""

		GET http://localhost:8080/home#!other
		routerKey 		= "/home"
		routerPath 		= "#!other"
		requestMethod	= "get"
		urlSuffix       = ""
	*/

	RouterData       interface{}   // custom storage data
	RouterDataRefVal reflect.Value // reflect value data

	RouterKey     string //
	RouterPath    string // have been converted to lowercase
	RequestMethod string // lowercase get post ...
	UrlSuffix     string // lowercase

	AppName string // application name
}

//
//	router interface
//
type IRouter interface {

	/**
	 *	1. after router parse
	 *
	 *	@param context
	 *	@param option
	 *	@return HttpStatus  200 pass
	 */
	AfterRouterParse(context *HttpContext, option *RouterOption) HttpStatus

	/**
	 *	2. parse func name
	 *
	 *	@param context			http context
	 *	@param option			router option
	 *	@return funcName 		function name specifies call
	 *	@return statusCode		http status code, 200 pass, other to status page
	 */
	ParseFuncName(context *HttpContext, option *RouterOption) (funcName string, statusCode HttpStatus, err error)

	/**
	 *	3. before call func
	 *
	 *	@param context
	 *	@param option
	 *	@return HttpStatus
	 */
	CallFuncBefore(context *HttpContext, option *RouterOption) HttpStatus

	/**
	 *	4. request func
	 *
	 *	@param context			http content
	 *	@param funcName			call controller func name
	 *	@param option			router option
	 *	@return returnValue		controller func return value
	 *	@return statusCode		http status code, 200 pass, other to status page
	 */
	CallFunc(context *HttpContext, funcName string, option *RouterOption) (returnValue interface{}, statusCode HttpStatus, err error)

	/**
	 *	5. parse template path
	 *	no need to add the suffix
	 *
	 *	@param context
	 *	@param funcName	 controller call func name
	 *	@param option	 router option
	 *	@return template path, suggest "[routerKey]/[funcName]"
	 */
	ParseTemplatePath(context *HttpContext, funcName string, option *RouterOption) string

	/**
	 *	6. after call func
	 *
	 *	@param context
	 *	@param option
	 */
	CallFuncAfter(context *HttpContext, option *RouterOption)

	/**
	 *	@return router key
	 */
	RouterKey() string

	/**
	 *	@return controller option
	 */
	ControllerOption() ControllerOption

	/**
	 *	@return controller info
	 */
	Info() string
}

//#pragma mark Leafveingo method ----------------------------------------------------------------------------------------------------

/**
 *	router parse
 *
 *	@param context
 *	@param reqPathNoSuffix 		the no suffix request path
 *	@param reqSuffix			url request suffix
 */
func routerParse(context *HttpContext, reqPathNoSuffix, reqSuffix string) (router IRouter, option *RouterOption, statusCode HttpStatus) {
    /*  routerParse*/
	statusCode = Status404

	var lowerReqPath string
	if context.lvServer.IsReqPathIgnoreCase() {
		lowerReqPath = SFStringsUtil.ToLower(reqPathNoSuffix)
	} else {
		lowerReqPath = reqPathNoSuffix
	}

	reqPathLen := len(lowerReqPath)

	// FOR routerList -> IF host -> IF scheme -> FOR keys -> IF reqPath

	listCount := len(context.lvServer.routerList)
	for i := 0; i < listCount; i++ {
		element := context.lvServer.routerList[i]

		if context.reqHost == element.host {

			keyCount := len(element.routerKeys)
			for j := 0; j < keyCount; j++ {

				key := element.routerKeys[j]
				keyLen := len(key)

				if keyLen <= reqPathLen && key == lowerReqPath[:keyLen] {

					//	controller router find
					if iR, ok := element.routers[key]; ok {

						//	uri scheme handle
						statusCode = routerSchemeHandle(context, iR.ControllerOption().Scheme())
						if statusCode != Status200 {
							router = nil
							option = nil
							return
						}

						context.routerElement = element
						router = iR
						option = new(RouterOption)
						option.AppName = context.lvServer.AppName()

						if 0 != len(reqSuffix) {
							if '.' == reqSuffix[0] {
								reqSuffix = reqSuffix[1:]
							}
							option.UrlSuffix = SFStringsUtil.ToLower(reqSuffix)
						}

						option.RequestMethod = SFStringsUtil.ToLower(context.Request.Method)
						if 0 == len(option.RequestMethod) {
							option.RequestMethod = "get"
						}

						option.RouterKey = key
						option.RouterPath = lowerReqPath[keyLen:]

						statusCode = router.AfterRouterParse(context, option)
					} else {
						//	基本上不会进来此处
						context.lvServer.log.Error("lv.routerKeys contains %#v and lv.routers not contains %#v", key, key)
						statusCode = Status404
					}
					return

				} // end  keyLen <= reqPathLen && key == lowerReqPath[:keyLen]

			} // end j := 0; j < keyCount; j++

		} // end context.reqHost == element.host

	} // end i := 0; i < listCount; i++

	return
}

/**
 *	router uri scheme handle
 *
 *	@param context
 *	@param scheme
 *	@return statusCode Status200 pass
 */
func routerSchemeHandle(context *HttpContext, scheme URIScheme) (statusCode HttpStatus) {
	statusCode = Status400

	switch context.RequestScheme() {
	case URI_SCHEME_HTTP:

		if scheme&URI_SCHEME_HTTP == URI_SCHEME_HTTP {
			statusCode = Status200
			break
		}

		//	TODO 下一步考虑委托函数处理不支持的协议，可让用户来处理。

		if scheme&URI_SCHEME_HTTPS == URI_SCHEME_HTTPS {
			if "GET" == context.Request.Method {
				//	重定向操作
				statusCode = StatusNil
				routerSchemeRedirect(context, URI_SCHEME_HTTPS)
			}
		}

	case URI_SCHEME_HTTPS:
		if scheme&URI_SCHEME_HTTPS == URI_SCHEME_HTTPS {
			statusCode = Status200
			break
		}

		if scheme&URI_SCHEME_HTTP == URI_SCHEME_HTTP {
			if "GET" == context.Request.Method {
				//	重定向操作
				statusCode = StatusNil
				routerSchemeRedirect(context, URI_SCHEME_HTTP)
			}
		}

	default:

	}

	return
}

/**
 *  does not support the scheme redirect handle
 *
 *  @param context
 *  @param scheme
 */
func routerSchemeRedirect(context *HttpContext, scheme URIScheme) {
	schemeStr := ""
	port := 0

	switch scheme {
	case URI_SCHEME_HTTP:
		schemeStr = "http://"
		port = context.lvServer.port
	case URI_SCHEME_HTTPS:
		if nil == context.lvServer.tlsListener {
			return
		}
		schemeStr = "https://"
		port = context.lvServer.tlsPort
	default:
		return
	}

	//	host
	host := context.Request.Host
	hostLen := len(host)

	//	remove port
	retScope := hostLen - 7
	if 0 > retScope {
		retScope = 0
	}
	index := -1
	for i := hostLen - 1; i >= retScope; i-- {
		if ':' == host[i] {
			index = i
			break
		}
	}
	if -1 != index {
		host = fmt.Sprintf("%s:%d", host[:index], port)
	}

	//	path
	path := context.Request.URL.Path
	if 0 == len(path) {
		path = "/"
	} else if '/' != path[0] {
		path = "/" + path
	}

	query := context.Request.URL.Query().Encode()
	if 0 != len(query) {
		query = "?" + query
	}

	urlstr := schemeStr + host + path + query

	http.Redirect(context.RespWrite, context.Request, urlstr, int(Status302))

	context.LVServer().Log().Debug("scheme redirect to: " + urlstr)
}

//  Copyright 2013 slowfei And The Contributors All rights reserved.
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.
//
//  Create on 2013-9-14
//  Update on 2014-07-17
//  Email  slowfei#foxmail.com
//  Home   http://www.slowfei.com

//
//	router manager
//
package leafveingo

import (
	// "github.com/slowfei/gosfcore/debug"
	"fmt"
	"github.com/slowfei/gosfcore/utils/strings"
	"net/http"
	"reflect"
)

var (
	//	globalRouterList
	_globalRouterList []globalRouter = nil
)

//#pragma mark global handle ----------------------------------------------------------------------------------------------------

//
//	global router
//
type globalRouter struct {
	appName string
	router  IRouter
}

/**
 *	global add router
 *
 *	@param appName
 *	@param routerKey
 *	@param router
 */
func AddRouter(appName string, router IRouter) {
	_globalRouterList = append(_globalRouterList, globalRouter{appName, router})
}

//#pragma mark interface option ----------------------------------------------------------------------------------------------------

//
//	router element
//
type RouterElement struct {
	host       string
	routerKeys []string
	routers    map[string]IRouter
}

//
//	router option
//
type RouterOption struct {
	/*
		GET http://localhost:8080/home/index.go
		routerKey 		= "/home/"
		routerPath 		= "index"
		requestMethod	= "get"
		urlSuffix       = "go"

		GET http://localhost:8080/home/manager/help.go
		routerKey 		= "/home/"
		routerPath 		= "manager/help"
		requestMethod	= "get"
		urlSuffix       = "go"

		GET http://localhost:8080/home/manager/other
		routerKey 		= "/home/"
		routerPath 		= "manager/other"
		requestMethod	= "get"
		urlSuffix       = ""

		GET http://localhost:8080/home#!other
		routerKey 		= "/home"
		routerPath 		= "#!other"
		requestMethod	= "get"
		urlSuffix       = ""
	*/

	RouterData       interface{}   // custom storage data
	RouterDataRefVal reflect.Value // reflect value data

	RouterKey     string //
	RouterPath    string // have been converted to lowercase
	RequestMethod string // lowercase get post ...
	UrlSuffix     string // lowercase

	AppName string // application name
}

//
//	router interface
//
type IRouter interface {

	/**
	 *	1. after router parse
	 *
	 *	@param context
	 *	@param option
	 *	@return HttpStatus  200 pass
	 */
	AfterRouterParse(context *HttpContext, option *RouterOption) HttpStatus

	/**
	 *	2. parse func name
	 *
	 *	@param context			http context
	 *	@param option			router option
	 *	@return funcName 		function name specifies call
	 *	@return statusCode		http status code, 200 pass, other to status page
	 */
	ParseFuncName(context *HttpContext, option *RouterOption) (funcName string, statusCode HttpStatus, err error)

	/**
	 *	3. before call func
	 *
	 *	@param context
	 *	@param option
	 *	@return HttpStatus
	 */
	CallFuncBefore(context *HttpContext, option *RouterOption) HttpStatus

	/**
	 *	4. request func
	 *
	 *	@param context			http content
	 *	@param funcName			call controller func name
	 *	@param option			router option
	 *	@return returnValue		controller func return value
	 *	@return statusCode		http status code, 200 pass, other to status page
	 */
	CallFunc(context *HttpContext, funcName string, option *RouterOption) (returnValue interface{}, statusCode HttpStatus, err error)

	/**
	 *	5. parse template path
	 *	no need to add the suffix
	 *
	 *	@param context
	 *	@param funcName	 controller call func name
	 *	@param option	 router option
	 *	@return template path, suggest "[routerKey]/[funcName]"
	 */
	ParseTemplatePath(context *HttpContext, funcName string, option *RouterOption) string

	/**
	 *	6. after call func
	 *
	 *	@param context
	 *	@param option
	 */
	CallFuncAfter(context *HttpContext, option *RouterOption)

	/**
	 *	@return router key
	 */
	RouterKey() string

	/**
	 *	@return controller option
	 */
	ControllerOption() ControllerOption

	/**
	 *	@return controller info
	 */
	Info() string
}

//#pragma mark Leafveingo method ----------------------------------------------------------------------------------------------------

/**
 *	router parse
 *
 *	@param context
 *	@param reqPathNoSuffix 		the no suffix request path
 *	@param reqSuffix			url request suffix
 */
func routerParse(context *HttpContext, reqPathNoSuffix, reqSuffix string) (router IRouter, option *RouterOption, statusCode HttpStatus) {
    /*  routerParse*/
	statusCode = Status404

	var lowerReqPath string
	if context.lvServer.IsReqPathIgnoreCase() {
		lowerReqPath = SFStringsUtil.ToLower(reqPathNoSuffix)
	} else {
		lowerReqPath = reqPathNoSuffix
	}

	reqPathLen := len(lowerReqPath)

	// FOR routerList -> IF host -> IF scheme -> FOR keys -> IF reqPath

	listCount := len(context.lvServer.routerList)
	for i := 0; i < listCount; i++ {
		element := context.lvServer.routerList[i]

		if context.reqHost == element.host {

			keyCount := len(element.routerKeys)
			for j := 0; j < keyCount; j++ {

				key := element.routerKeys[j]
				keyLen := len(key)

				if keyLen <= reqPathLen && key == lowerReqPath[:keyLen] {

					//	controller router find
					if iR, ok := element.routers[key]; ok {

						//	uri scheme handle
						statusCode = routerSchemeHandle(context, iR.ControllerOption().Scheme())
						if statusCode != Status200 {
							router = nil
							option = nil
							return
						}

						context.routerElement = element
						router = iR
						option = new(RouterOption)
						option.AppName = context.lvServer.AppName()

						if 0 != len(reqSuffix) {
							if '.' == reqSuffix[0] {
								reqSuffix = reqSuffix[1:]
							}
							option.UrlSuffix = SFStringsUtil.ToLower(reqSuffix)
						}

						option.RequestMethod = SFStringsUtil.ToLower(context.Request.Method)
						if 0 == len(option.RequestMethod) {
							option.RequestMethod = "get"
						}

						option.RouterKey = key
						option.RouterPath = lowerReqPath[keyLen:]

						statusCode = router.AfterRouterParse(context, option)
					} else {
						//	基本上不会进来此处
						context.lvServer.log.Error("lv.routerKeys contains %#v and lv.routers not contains %#v", key, key)
						statusCode = Status404
					}
					return

				} // end  keyLen <= reqPathLen && key == lowerReqPath[:keyLen]

			} // end j := 0; j < keyCount; j++

		} // end context.reqHost == element.host

	} // end i := 0; i < listCount; i++

	return
}

/**
 *	router uri scheme handle
 *
 *	@param context
 *	@param scheme
 *	@return statusCode Status200 pass
 */
func routerSchemeHandle(context *HttpContext, scheme URIScheme) (statusCode HttpStatus) {
	statusCode = Status400

	switch context.RequestScheme() {
	case URI_SCHEME_HTTP:

		if scheme&URI_SCHEME_HTTP == URI_SCHEME_HTTP {
			statusCode = Status200
			break
		}

		//	TODO 下一步考虑委托函数处理不支持的协议，可让用户来处理。

		if scheme&URI_SCHEME_HTTPS == URI_SCHEME_HTTPS {
			if "GET" == context.Request.Method {
				//	重定向操作
				statusCode = StatusNil
				routerSchemeRedirect(context, URI_SCHEME_HTTPS)
			}
		}

	case URI_SCHEME_HTTPS:
		if scheme&URI_SCHEME_HTTPS == URI_SCHEME_HTTPS {
			statusCode = Status200
			break
		}

		if scheme&URI_SCHEME_HTTP == URI_SCHEME_HTTP {
			if "GET" == context.Request.Method {
				//	重定向操作
				statusCode = StatusNil
				routerSchemeRedirect(context, URI_SCHEME_HTTP)
			}
		}

	default:

	}

	return
}

/**
 *  does not support the scheme redirect handle
 *
 *  @param context
 *  @param scheme
 */
func routerSchemeRedirect(context *HttpContext, scheme URIScheme) {
	schemeStr := ""
	port := 0

	switch scheme {
	case URI_SCHEME_HTTP:
		schemeStr = "http://"
		port = context.lvServer.port
	case URI_SCHEME_HTTPS:
		if nil == context.lvServer.tlsListener {
			return
		}
		schemeStr = "https://"
		port = context.lvServer.tlsPort
	default:
		return
	}

	//	host
	host := context.Request.Host
	hostLen := len(host)

	//	remove port
	retScope := hostLen - 7
	if 0 > retScope {
		retScope = 0
	}
	index := -1
	for i := hostLen - 1; i >= retScope; i-- {
		if ':' == host[i] {
			index = i
			break
		}
	}
	if -1 != index {
		host = fmt.Sprintf("%s:%d", host[:index], port)
	}

	//	path
	path := context.Request.URL.Path
	if 0 == len(path) {
		path = "/"
	} else if '/' != path[0] {
		path = "/" + path
	}

	query := context.Request.URL.Query().Encode()
	if 0 != len(query) {
		query = "?" + query
	}

	urlstr := schemeStr + host + path + query

	http.Redirect(context.RespWrite, context.Request, urlstr, int(Status302))

	context.LVServer().Log().Debug("scheme redirect to: " + urlstr)
}

//  Copyright 2013 slowfei And The Contributors All rights reserved.
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.
//
//  Create on 2013-9-14
//  Update on 2014-07-17
//  Email  slowfei#foxmail.com
//  Home   http://www.slowfei.com

//
//	router manager
//
package leafveingo

import (
	// "github.com/slowfei/gosfcore/debug"
	"fmt"
	"github.com/slowfei/gosfcore/utils/strings"
	"net/http"
	"reflect"
)

var (
	//	globalRouterList
	_globalRouterList []globalRouter = nil
)

//#pragma mark global handle ----------------------------------------------------------------------------------------------------

//
//	global router
//
type globalRouter struct {
	appName string
	router  IRouter
}