//        ___                          ___               
//       /\_ \                       / ___\       __     
//   ____\//\ \     ___   __  __  __/\ \__/   __ /\_\    
//  /  __\ \ \ \   / __ \/\ \/\ \/\ \ \  __\/ __ \/\ \   
// /\__   \ \_\ \_/\ \s\ \ \ \_/ \_/ \ \ \_/\  __/\ \ \  
// \/\____/ /\____\ \____/\ \___f___/ \ \_\\ \____\\ \_\ 
//  \/___/  \/____/\/___/  \/__//__/   \/_/ \/____/ \/_/ 
//                                                       
//              http://www.slowfei.com                   
//               slowfei#foxmail.com                     
//-------------------------------------------------------

;(function($){
    var MD_BASE_URL = "md";
    var MD_FILE_SUFFIX = ".md";
    var GOSFDOC_JSON = "config.json";
    var QUERY_KEY_PACKAGE = "p";
    var QUERY_KEY_TITLE = "t";
    var QUERY_KEY_ANCHOR = "a";
    var QUERY_KEY_VERSION = "v";
    var DATA_KEY_STICKYITEMS = "stickyItems";
    var COOKIE_LANG_KEY = "language";
    var COOKIE_VERSION_KEY = "version"
    var LANG_DEFAULT = "default";
    var VERSION_DEFAULT = "1.0.0";
    var _language = LANG_DEFAULT;
    var _version = null;
    var dataGosfdocJson = null;
    var $mainContent = null;
    var navigationHeight = 46;

    /**
     *  jquery document load
     */
    $(function(){
        // Synchronous highlighting with highlight.js
        hljs.configure({
          tabReplace: '    ', // 4 spaces
          useBR: true
        });

        // markdown set
        marked.setOptions({
            langPrefix:''
        });

        marked.setOptions({
            highlight: function (code,lang) {
                if (lang) {
                    if ( hljs.getLanguage(lang) ) {
                         return hljs.fixMarkup(hljs.highlightAuto(code,[lang]).value);
                    }else{
                        return hljs.fixMarkup(hljs.highlightAuto(code).value);
                    }
                   
                }
            }
        });

        $mainContent = $("#main_content");

        // menu show and hide button
        $('#btn_show_menu').click(function(){
            var t = $(this).text();
            $(this).text(t == ">" ? "<":">");
            $('#main_sidebar').sidebar('toggle');
        });

        //  init doc
        initDoc();
        
        //  menu item content text out of range handle
        $.each($("#main_sidebar .item a.item"), function(index, val) {
            var $_this = $(val);

            var centHtml = $_this.html();
            centHtml = centHtml.replace(/\//g,"\/<wbr>").replace(/\./g,".<wbr>");
            $_this.html(centHtml);

            var frameWidth = $_this.outerWidth() - 10;
            var realWidth = getInnerWidth(val);
            if (realWidth >= frameWidth) {
                $_this.popup({
                    on:'hover',
                    content:"<a href=\""+$_this.attr('href')+"\">"+$_this.text()+"</a>",
                    transition:'none',
                    position:'left center',
                    className: {
                         popup: 'ui popup inverted item_popup'
                    },
                    onShow:function(element){
                        var $_item_popup = $(".ui.popup.item_popup");
                        var popupTop = parseFloat($_item_popup.css("top"));
                        var popupPadd = parseFloat($_item_popup.css("padding-top"));
                        var popupHeight = $_item_popup.outerHeight();
                        var itemHeight = $_this.outerHeight();

                        var padd = (itemHeight - popupHeight) / 2.0;

                        if ( 0 >= padd ) {
                            padd = '0.8em';
                        }else{
                            popupTop = popupTop - padd + 0.5;
                            padd += popupPadd;
                        }

                        $_item_popup.css({
                            left            : '-11px',
                            paddingTop      : padd,
                            paddingBottom   : padd,
                            top             : popupTop
                        }).mouseout(function(){
                            $(this).hide();
                            $(this).remove();
                        });
                    }
                });
            };
        });

        //  about model show
        $("#item_about").click(function(){
            $("#about_modal").modal('setting',{
                duration:500
            }).modal('show');
            return false;
        });

        //  to home page
        $("#home").click(function(event) {
            reHome();
        }); 


        //  monitor sticky
        $(window).scroll(stickyScroll);

        $(window).resize(function(event) {
            var $stickyWrapper = $("#sticky-wrapper");
            $stickyWrapper.css('height', 'auto');
            var swTop = parseFloat($stickyWrapper.css('top'));
            var $window = $(window);
            if ($window.height() - swTop < $stickyWrapper.outerHeight() ) {
                 $stickyWrapper.height($window.height() - swTop - 10);
            };
        });

    });


    /**
     *  init doc
     */
    function initDoc(){
    
       $.ajax({
           url: GOSFDOC_JSON,
           async:false,
           cache:false,
           type: 'GET',
           dataType: 'json'
       })
       .done(function(dataJson) {
            dataGosfdocJson = dataJson;
            var queryPackage = getURIQuery(QUERY_KEY_PACKAGE);

            parseSelectVersion(dataJson.Versions);
            parseContentJson(dataJson.ContentJson);
            parseAbout(dataJson.AboutMd);
            parseLanguages(dataJson.Languages);
            queryPackage = parseMenuList(dataJson.Markdowns,queryPackage);

            if (queryPackage) {
                //  handle query package
                $(".segment.intro").css("paddingTop","20px");
                parsePackageMarkdown(queryPackage);
            }else{
                // home show
                reHome();
            }
       })
       .fail(function(jqXHR, textStatus, errorThrown) {
            $('#btn_show_menu').click();
            $("#segment_intro").empty().html("Sorry! can not read \""+ GOSFDOC_JSON + "\" file.<br/><br/> try to use <code>gosfdoc</code> command run document: ");
            //  TODO 读取 gosfdoc.json 的错误处理
       });
    }

    /**
     *  reset Home page
     */
    function reHome(){
        if (null == dataGosfdocJson) {
            return;
        }

        $(window).scrollTop(0);

        $(".segment.intro").css("paddingTop","");
        $("#main_sidebar a").removeClass('active');

        parseIntro(dataGosfdocJson.IntroMd);

        if ( $("#main_sidebar").hasClass('active')) {
            $('#btn_show_menu').click();
        };
       
        var $contentPackage = $('<div class="ui list"></div>');
        var addCount = 0;

        //  each markdowns
       $.each(dataGosfdocJson.Markdowns, function(index, val) {
            var proVer = val.Version;

            if (_version != proVer) {
                return true;
            }

            var projectName = val.MenuName;
            var $listItem = $('<div class="item"></div>');
            var $listContent = $('<div class="content"><div class="header">'+projectName+'</div></div>');
            var $list = $('<div class="list"></div>');
          
            $.each(val.List, function(listIndex, packageInfo) {
                var packageName = packageInfo.Name;
                var packageDesc = packageInfo.Desc;
                var mdpath = packageName + MD_FILE_SUFFIX;
          
                var $item = $('<div class="item"><div class="content"><a class="header" href=?p='+mdpath+'>'+packageName+'</a><div class="description">'+packageDesc+'</div></div></div>')
                
                $("a",$item).click(function() {
                    $(window).scrollTop(0);
                    if(!$(this).hasClass('active')){
                        setURIQuery(mdpath,false,false,_version);
                    }
                    parsePackageMarkdown(mdpath);
                    return false;
                });

                $list.append($item);
            });

            addCount++;
            $listItem.append($listContent);
            $listItem.append($list);
            $contentPackage.append($listItem);

        });// end  $.each(dataJson.Markdowns, function(index, val)

        if ( 0 == addCount) {
            $contentPackage.html("No info data. Version: "+_version);
        };

        $mainContent.attr('mdpath', null);
        $mainContent.empty().html($contentPackage);
        $("#sticky").empty().hide();
        $mainContent.css("marginLeft","20px");

        setURIQuery(false,false,false,_version);
    }

    /**
     *  monitor sticky item
     */
    function stickyScroll(){
        var windowTop = $(window).scrollTop();

        var stickyItems =  $mainContent.data(DATA_KEY_STICKYITEMS);
        if (!stickyItems) {
            return;
        };
        var item = null;
        $.each(stickyItems, function(key, val) {
            if ( windowTop >= key ) {
                 item = val;
            };
        });

        if (item) {
            setStickyItemActive(item);
        };
    }

    /**
     *  set sticky item active
     *
     *  @param item
     */
    function setStickyItemActive(item){
        var windowTop = $(window).scrollTop();
        var $_item = $(item);

        if (!$_item.hasClass('active')) {
            $("#sticky a").removeClass('active');
            $(item).addClass('active');

            //  set uri
            var packageName = getURIQuery(QUERY_KEY_PACKAGE);
            var anchor = getURIQuery(QUERY_KEY_ANCHOR);

            //
            setURIQuery(packageName,$_item.text(),anchor,_version);

            // overflow out scroll
            var $stickyWrapper = $("#sticky-wrapper");
            var swTop = parseFloat($stickyWrapper.css('top'));
            var itemOffsetTop = $_item.offset().top - windowTop - swTop;
            var itemouterHeight = $_item.outerHeight();

            var swHeight = $stickyWrapper.height();
            var swScrollTop = $stickyWrapper.scrollTop();

            if (itemOffsetTop + itemouterHeight > swHeight - swScrollTop ) {
                $stickyWrapper.scrollTop(swScrollTop + itemouterHeight * 2);
            }else if(itemOffsetTop < swScrollTop ){
                $stickyWrapper.scrollTop(swScrollTop - itemouterHeight * 2);
            }

        };
    }

    /**
     *  parse package markdown
     *  
     *  @param mdpath
     */
    function parsePackageMarkdown(mdpath){

        $("#segment_intro").empty();
        $(".segment.intro").css("paddingTop","20px");

        var $menuItems = $("#main_sidebar a");
        $menuItems.removeClass('active');
        $.each($menuItems, function(index, val) {
            var $_this = $(val);
            var href = $_this.attr('href');
            if (href == '?'+QUERY_KEY_PACKAGE +'=' + mdpath) {
                $_this.addClass('active');
                return false;
            };
        });

        var queryTitle = getURIQuery(QUERY_KEY_TITLE);
        var queryAnchor = getURIQuery(QUERY_KEY_ANCHOR);
        setURIQuery(mdpath,queryTitle,queryAnchor,_version);

        $mainContent.removeData(DATA_KEY_STICKYITEMS);
        $("#sticky").empty().hide();
        $mainContent.css("marginLeft","20px");

        //  load content package markdown
        ajaxGet({
                path : mdpath,
                async: true,
                dataType : 'text',
                doneFunc:function(text){
                    $mainContent.empty().html(marked(text));
                    $mainContent.attr('mdpath', mdpath);

                    var queryTitle = getURIQuery(QUERY_KEY_TITLE);
                    var queryAnchor = getURIQuery(QUERY_KEY_ANCHOR);

                    var titleScrollTop = false;

                    //  each #main_content element child node h2 tag
                    var h2Items = {}
                    $.each($("h2",  $mainContent), function(index, val) {
                        var $stickyItem = $(val);
                        var offsetTop = $stickyItem.offset().top - navigationHeight;  //  top navigation height
                        h2Items[offsetTop] = val;
                    });

                    //  save sticky items
                    var stickyItems = {};
                    $.each(h2Items, function(key, val) {
                        var text = $(val).text();
                        var $item = $('<a class="item" href="javascript:;">'+text+'</a>');

                        $item.click(function() {
                            var queryPackage = getURIQuery(QUERY_KEY_PACKAGE);
                            var queryTitle = getURIQuery(QUERY_KEY_TITLE); 
                            setURIQuery(queryPackage,queryTitle,false,_version);

                            $(window).scrollTop(key);
                            return false;
                        });

                        $("#sticky").append($item);
                        stickyItems[parseInt(key)] = $item;

                        if (queryTitle == text && !titleScrollTop ) {
                            titleScrollTop = key;
                        };
                    });

                    //  handle a tag anchor point scroll
                    $.each($("a", $mainContent), function(index, val) {
                        var $atag = $(val);
                        var href = $atag.attr('href');

                        if ( href && 0 < href.length ) {

                            if ('#' == href[0]) {

                                //  set anchor tag click
                                $atag.click(function() {
                                    var $_this = $(this);
                                    var href = $_this.attr('href');
                                    var text = $_this.text();

                                    if (!href && 0 >= href.length) {
                                        return false;
                                    };

                                    var anchor = href.substring(1,href.length);

                                    var queryPackage = getURIQuery(QUERY_KEY_PACKAGE);
                                    var queryTitle = getURIQuery(QUERY_KEY_TITLE);

                                    setURIQuery(queryPackage,queryTitle,unescape(anchor),_version);

                                    if ( '#' != text ) {
                                        var $anchorTag = $("*[id='"+anchor+"']",$mainContent);
                                        if ( 0 == $anchorTag.length ) {
                                            $anchorTag = $("*[name='"+anchor+"']",$mainContent);
                                        }
                                        if ( 0 != $anchorTag.length ) {
                                            $anchorTag = $anchorTag.eq(0);
                                            var offsetTop = $anchorTag.offset().top - navigationHeight;
                                            $(window).scrollTop(offsetTop);
                                        }
                                    };
                                    return false;
                                });

                            }else if( 0 == href.indexOf('../') ){
                                href = href.replace(/\.\.\//g,"");
                                $atag.attr('href', href);
                            }
                        }

                    });
                    
                    //  scroll anchor
                    if ( queryAnchor ) {

                        queryAnchor = escape(queryAnchor);
                        
                        var $anchorTag = $("*[id='"+queryAnchor+"']",$mainContent);
                        if ( 0 == $anchorTag.length ) {
                            $anchorTag = $("*[name='"+queryAnchor+"']",$mainContent);
                        }
                        if ( 0 != $anchorTag.length ) {
                            $anchorTag = $anchorTag.eq(0);
                            titleScrollTop = $anchorTag.offset().top - navigationHeight;
                        }

                        var queryPackage = getURIQuery(QUERY_KEY_PACKAGE);
                        var queryTitle = getURIQuery(QUERY_KEY_TITLE);
                        setURIQuery(queryPackage,queryTitle,queryAnchor,_version);
                    };

                    //  scroll to func position
                    if (titleScrollTop) {
                        $(window).scrollTop(titleScrollTop);
                    };

                    if ( 0 != $("#sticky a").length ) {
                        $("#sticky").show();
                        $mainContent.css("marginLeft","170px");

                        var $stickyWrapper = $("#sticky-wrapper");
                        var swTop = parseFloat($stickyWrapper.css('top'));
                        var $window = $(window);

                        if ($window.height() - swTop < $stickyWrapper.outerHeight() ) {
                             $stickyWrapper.height($window.height() - swTop - 10);
                        };
                    }

                    $mainContent.data(DATA_KEY_STICKYITEMS, stickyItems);
                },
                failFunc:function(){
                    $mainContent.empty().html("Sorry! Can not load "+ mdpath);
                    $mainContent.attr('mdpath', '');
                }
            }
        );
    }

    /**
     *  parse left menu list data show info
     *
     *  @param markdownJson
     *  @param findPackage Need to check whether the existing package, null not check
     *  @return confirmed the package, other return null
     */
    function parseMenuList(markdownJson,findPackage){
        var resultPackage = null;

        var $sidebarItem = $('<div class=item></div>');
        var addCount = 0;

        //  each markdowns
        $.each(markdownJson, function(index, val) {
            var proVer = val.Version;

            if ( _version != proVer ) {
                return true;
            }

            var projectName = val.MenuName;
            var $menu = $('<div class="menu"></div>');
            var $itemTitle = $('<b>'+projectName+'</b>');

            $.each(val.List, function(listIndex, packageInfo) {
                var packageName = packageInfo.Name;

                if ( findPackage ) {
                    var mdpath = packageName + MD_FILE_SUFFIX;
                    if (mdpath == findPackage) {
                        resultPackage = findPackage;
                    }
                }
                
                var $item = $('<a class="item" href=?p='+mdpath+'>'+packageName+'</a>');

                $item.click(function() {
                    if(!$(this).hasClass('active')){
                        $(window).scrollTop(0);
                        setURIQuery(mdpath,false,false,_version);
                        parsePackageMarkdown(mdpath);
                    }
                    return false;
                });

                $menu.append($item);
            });

            addCount++;
            $sidebarItem.append($itemTitle);
            $sidebarItem.append($menu);

        });// end  $.each(dataJson.Markdowns, function(index, val)

        if ( 0 == addCount) {
            $sidebarItem.html("No info data. Version: "+_version);
        };

        var $tempMenuTitle = $("#menu_title");
        $("#main_sidebar").empty();
        $("#main_sidebar").append($tempMenuTitle);

        $("#main_sidebar").append($sidebarItem);

        return resultPackage;
    }

    /**
     *  parse version select info
     *
     *  @param verJson
     */
    function parseSelectVersion(verJson){

        var ver = getURIQuery(QUERY_KEY_VERSION);

        if (!ver) {
            ver = getCookie(COOKIE_VERSION_KEY);
        }

        if ( ver ) {
            _version = ver;
        };

        var tempOne = "";
        var checkVer = false;
        var $versionElements = $("#version_value");
        $versionElements.empty();

        $.each(verJson, function(index, val) {

            var verstr = val.toString();

            //  默认选择排序在第一位的版本
            if ( 0 == tempOne.length) {
                tempOne = verstr;
            }
            
            //  效验设置的版本信息是否与获取的版本信息相符
            if ( _version == verstr ) {
                checkVer = true;
            }   

            var html = '<div class="item" data-value="'+verstr+'">'+verstr+'</div>';
            $versionElements.append(html);

        });

        if (!checkVer) {
            if ( 0 == tempOne.length ) {
                _version = VERSION_DEFAULT;
            }else{
                _version = tempOne;
            }
        }

        $("#version_text").text(_version);

       //  language handle
        $('.ui.dropdown.version').unbind().dropdown({
            on:"hover",
            onChange:function(value, text){
                if (_version == value) {
                    return;
                }
                
                $("#version_text").text(text);
                setCookie(COOKIE_VERSION_KEY,value);
                _version = value;

                var packageName = getURIQuery(QUERY_KEY_PACKAGE);
                var queryTitle = getURIQuery(QUERY_KEY_TITLE);
                var anchor = getURIQuery(QUERY_KEY_ANCHOR);
                setURIQuery(packageName,queryTitle,anchor,_version);

                var jumpHome = true;
                var mdpath = $mainContent.attr('mdpath');
                if (mdpath && null != dataGosfdocJson) {
                    var checkPackage = parseMenuList(dataGosfdocJson.Markdowns,mdpath);
                    if (checkPackage) {
                        parsePackageMarkdown(checkPackage);
                        jumpHome = false;
                    }
                }

                if (jumpHome) {
                     reHome();
                }
            }
        });
    }

    /**
     *  parse language elements
     *
     *  @param langJson
     */
    function parseLanguages(langJson){
        var lang = getCookie(COOKIE_LANG_KEY);
        if ( null != lang ) {
            _language = lang;
        };

        var showText = "";
        var $langElements = $("#language_value");
        $langElements.empty();

        $.each(langJson, function(index, val) {

            var key = "";
            var mapVal = "";

            $.each(val, function(mk, mv) {
                key = mk;
                mapVal = mv;
                return false;
            });

            if ( 0 == showText.length ) {
                showText = mapVal;
            }

            if ( key == _language ) {
                showText = mapVal;
            }

            var html = '<div class="item" data-value="'+key+'">'+mapVal+'</div>';
            $langElements.append(html);


        });

        $("#language_text").text(showText);

       //  language handle
        $('.ui.dropdown.language').unbind().dropdown({
            on:"hover",
            onChange:function(value, text){
                $("#language_text").text(text);
                setCookie(COOKIE_LANG_KEY,value);
                _language = value;
                var mdpath = $mainContent.attr('mdpath');
                if (mdpath) {
                    parsePackageMarkdown(mdpath);
                }else{
                    reHome();
                }
            }
        });
    }

    /**
     *  parse about content
     *
     *  @param aboutpath
     */
    function parseAbout(aboutpath){
        ajaxGet({
                path : aboutpath,
                async: true,
                dataType : 'text',
                doneFunc:function(text){
                    $("#about_content").empty().html(marked(text));
                },
                failFunc:function(){
                    $("#about_content").empty().html("Can not read "+ aboutpath + " version:"+_version);
                }
            }
        );
    }

    /**
     *  parse content json
     */
    function parseContentJson(contpath){
        ajaxGet({
                path : contpath,
                async: true,
                dataType : 'json',
                doneFunc:function(dataJson){
                    document.title = dataJson.HtmlTitle;
                    $("#doc_title").empty().html(dataJson.DocTitle);
                    $("#menu_title").empty().html(dataJson.MenuTitle);
                },
                failFunc:function(){
                    alert("Can not read "+ contpath + " version:"+_version);
                }
            }
        );
    }

    /**
     *  parse intro
     */
    function parseIntro(mdpath){
        ajaxGet({
                path : mdpath,
                async: true,
                dataType : 'text',
                doneFunc:function(text){
                    $("#segment_intro").empty().html(marked(text));
                },
                failFunc:function(){
                    $("#segment_intro").empty().html("Can not read "+ mdpath + " version:"+_version);
                }
            }
        );
    }

    //#pragma mark other method ----------------------------------------------------------------------------

    /**
     *  sync get request 
     *      
     *  @param path      does not include "md/default/" prefix
     *  @param dataType  default 'text'
     *  @param doneFunc     
     *  @param failFunc
     */
    function ajaxGet(options){

        var option = $.extend ({ 
            path:'', 
            dataType:'text',
            async   :false,
            cache   :true,
            doneFunc:null,
            failFunc:null
        },options );

        //  先根据选择的语言选择发送请求，如果出现错误在选择默认路径。
        $.ajax({
           url: getMDUrl(option.path),
           async:option.async,
           cache:option.cache,
           type: 'GET',
           dataType: option.dataType
       })
       .done(option.doneFunc)
       .fail(function() {
            $.ajax({
               url: getMDDefaultUrl(option.path),
               async:option.async,
               cache:option.cache,
               type: 'GET',
               dataType: option.dataType
            })
            .done(option.doneFunc)
            .fail(option.failFunc);
       });
    }

    /**
     *  conversion version use path
     * 
     *  @param version 
     *  @return version path
     */
    function converToVersionPath(version){
        var result = "";
        if ( version && 0 != version.length) {
            result = "v"+version.replace(/\./g,"_") + "/";
        }
        return result;
    }

    /**
     *  get markdown directory url
     *  by select language
     *
     *  @param path 
     *  @return
     */
    function getMDUrl(path){
        var verPath = converToVersionPath(_version);
        return verPath + MD_BASE_URL + '/' + _language +'/' + path;
    }

    /**
     *  get markdown directory url
     *  by default language
     *  
     *  @param path 
     *  @return
     */
    function getMDDefaultUrl(path){
        var verPath = converToVersionPath(_version);
        return verPath + MD_BASE_URL + '/' + LANG_DEFAULT + '/' + path;
    }

    /**
     *  determine whether the default directory
     *
     *  @param url
     *  @return true is default directory
     */
    function isMDDefaultUrl(url){
        var result = true;
        var defaultDir = MD_BASE_URL + '/' + LANG_DEFAULT;
        var indexOf = url.indexOf(defaultDir);

        if (0 != indexOf) {
            result = false;
        };

        return result;
    }

    /**
     *  set fixed uri query
     *
     *  @param packageval
     *  @param title
     *  @param anchor
     *  @param version
     */
    function setURIQuery(packageval,title,anchor,version){
        var hash = "";

        if ( version ){
            hash += QUERY_KEY_VERSION + '=' + escape(version);
        }

        if ( packageval ) {
            if ( 0 != hash.length ){
                hash += '&';   
            }
            hash += QUERY_KEY_PACKAGE + '=' + escape(packageval);
        }

        if ( title ) {
            if ( 0 != hash.length ){
                hash += '&';   
            }
            hash += QUERY_KEY_TITLE + '=' + escape(title);
        }

        if ( anchor ) {
            if ( 0 != hash.length ){
                hash += '&';   
            }
            hash += QUERY_KEY_ANCHOR + '=' + escape(anchor);
        }

        window.location.hash = hash;
        window.location.search ='';
        window.location.query = '';
    }

    /**
     *  get uri query param
     *
     *  support: http://localhost:8080/?name=value || http://localhost:8080/#name=value
     *
     *  @param name query name
     *  @return exist to query string, else to false
     */
    function getURIQuery(name){

        var queryStr = unescape(window.location.search) + '&';

        if (queryStr == '&') {
            queryStr = unescape(window.location.hash) + '&';
        }

        var regex = new RegExp('.*?[&\\?#]' + name + '=(.*?)&.*');

        var val = queryStr.replace(regex, "$1");

        var result = val == queryStr ? false : val;

        return result;
    }

    /**
     *  get real width content 
     *
     *  @param element javascript element
     */
    function getInnerWidth(element) {
        var result;
        var wrapper = document.createElement('span');
        while (element.firstChild) {
            wrapper.appendChild(element.firstChild);
        }
        element.appendChild(wrapper);
        result = wrapper.offsetWidth;
        element.removeChild(wrapper);
        while (wrapper.firstChild) {
            element.appendChild(wrapper.firstChild);
        }
        return result;
    }

    /**
     *  set cookie
     *  expires default 30 day
     *
     *  @param name
     *  @param value
     */
    function setCookie(name,value){
        var exp = new Date();
        exp.setTime(exp.getTime() + 30*24*60*60*1000);
        document.cookie = name + "="+ escape (value) + ";expires=" + exp.toGMTString();
    }

    /**
     *  get cookie
     *
     *  @param name
     */
    function getCookie(name){
        var result = null
        var arr,reg = new RegExp("(^| )"+ name +"=([^;]*)(;|$)");
        if(arr = document.cookie.match(reg)){
            result = (arr[2]);
        }
        return result;
    }

    /**
     *  delete cookie
     *
     *  @param name
     */
    function delCookie(name){

        var exp = new Date();
        exp.setTime(exp.getTime() - 1);

        var cval = getCookie(name);

        if(cval != null){
            document.cookie = name + "="+cval+";expires="+exp.toGMTString();
        }
    }
})(jQuery);