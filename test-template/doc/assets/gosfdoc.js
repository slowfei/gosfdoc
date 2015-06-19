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
    var NAVIGATION_HEIGHT = 46;
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

    var dataGosfdocJson = null;
    var $mainContent = null;
    var _language = LANG_DEFAULT;
    var _version = null;
    var _userRehash = false;
    var _rehashTimeout = null;
    var _appendPath = "";
    var _rexLineGithub = /^(L(\d+)|L(\d+)[-]L(\d+))$/;
    var _markAtag = /<a .+<\/a>/g;

    /**
     *  history struct
     *
     *  @param key
     *  @param value
     *  @param version
     */
    function HistoryStruct(key,value,version){
        this.key = key;
        this.hash = value;
        this.scrollTop = 0;
        this.version = version;
    }

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

        // highlight set
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

        //  init doc
        initDoc();

        //  
        initUI();

        //
        initEvent();

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
             if (dataJson.AppendPath){
                _appendPath = dataJson.AppendPath;
            }
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
     *  init ui
     */
    function initUI(){
        // menu show and hide button
        $('#btn_show_menu').click(function(){
            var t = $(this).text();
            $(this).text(t == ">" ? "<":">");
            $('#main_sidebar').sidebar('toggle');
        });
        
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

    }

    /**
     *  init event
     */
    function initEvent(){

        var histroyIndex = -1;
        var historyList = new Array();
        var maxHistroy = 20;
        var userReplace = false;

        //  First loaded record
        var key = getURIQuery(QUERY_KEY_PACKAGE);
        if (key) {
            var hash = window.location.hash;
            var history = new HistoryStruct(key,hash,_version);
            histroyIndex = historyList.push(history);
        }

        //  monitor sticky
        $(window).scroll(function(event) {
            var windowTop = $(window).scrollTop();

            var stickyItems =  $mainContent.data(DATA_KEY_STICKYITEMS);
            if (!stickyItems) {
                return;
            }

            var item = null;
            $.each(stickyItems, function(key, val) {
                if ( windowTop >= key ) {
                     item = val;
                };
            });

            if (item) {
                setStickyItemActive(item);
            }

            //  scroll height record browsing history
            if( 1 <= histroyIndex && historyList.length >=  histroyIndex ){
                var history = historyList[histroyIndex-1];
                if ( history ) {
                    history.scrollTop = windowTop;
                    historyList[histroyIndex-1] = history;
                }
            }
        });

        //
        $(window).resize(function(event) {
            var $stickyWrapper = $("#sticky-wrapper");
            $stickyWrapper.css('height', 'auto');
            var swTop = parseFloat($stickyWrapper.css('top'));
            var $window = $(window);
            if ($window.height() - swTop < $stickyWrapper.outerHeight() ) {
                 $stickyWrapper.height($window.height() - swTop - 10);
            };
        });

        //
        $(window).on('hashchange', function() {
            var key = getURIQuery(QUERY_KEY_PACKAGE);
            var currnetVer = _version;
            var hash = window.location.hash;
            if (!key || userReplace){
                return;
            }

            if (!_userRehash && 1 <= histroyIndex && historyList.length >=  histroyIndex ) {
                //  TODO 由于无法控制浏览器向后和向前的事件控制，所以历史浏览功能改用上下左右键来代替。
                var history = historyList[histroyIndex-1];
                var scrollTop = 0;
                var historyVersion = false;
                if (history) {
                    scrollTop = history.scrollTop;
                    historyVersion = history.version;
                }

                //  hashchange switch
                userReplace = true;

                //  version handle
                if ( historyVersion ) {
                    if ( historyVersion != currnetVer) {
                        $("#version_text").text(historyVersion);
                        $.each($("div.item",$("#version_value")), function(index, val) {
                            var $item = $(val);
                            if ($item.attr('data-value') == historyVersion) {
                                $item.addClass('active');
                            }else{
                                $item.removeClass('active');
                            }
                        });
                        _version = historyVersion;
                        if (null != dataGosfdocJson) {
                            parseMenuList(dataGosfdocJson.Markdowns,key);  
                        }
                    }
                }

                parsePackageMarkdown(key,scrollTop);

                userReplace = false;
                return;
            }
            
            var isAdd = false;
            if ( 0 == historyList.length) {
                isAdd = true;
            }else if( 1 <= histroyIndex && historyList.length >=  histroyIndex ){
                var history = historyList[histroyIndex-1];
                if ( history.key == key && currnetVer == history.version ) {
                    history.hash = hash;
                    history.scrollTop = $(window).scrollTop();
                    history.version = currnetVer;
                    historyList[histroyIndex-1] = history;
                }else{
                    isAdd = true;
                    var count = historyList.length;
                    for (var i = histroyIndex; i < count; i++) {
                        historyList.pop();
                    }
                }
            }

            if (isAdd) {
                if ( maxHistroy <= historyList.length ) {
                    historyList.shift();
                }
                var history = new HistoryStruct(key,hash,currnetVer);
                histroyIndex = historyList.push(history);
            }

            // console.log(historyList.length);
        });

        //  up down left right key
        //  up key control left menu previous item 
        //  down key control left menu next item 
        //  left key control browsing history back
        //  right key control browsing history forward
        $(document).keyup(function(event) {
            var which = event.which;

            switch (which){
                case 38:{
                    // alert("up");
                    var $items = $("#main_sidebar .menu a.item");
                    var itemCount = $items.length;
                    if ( 0 < itemCount ) {
                        var selectIndex = 0;
                        $.each($items, function(index, val) {
                            var $item = $(val);
                            if ($item.hasClass('active')) {
                                selectIndex = index-1;
                                return false;
                            }
                            return true;
                        });
                        if (0 > selectIndex) {
                            $items.eq(itemCount-1).click();
                        }else{
                            $items.eq(selectIndex).click();
                        }
                    }
                }
                break;
                case 40:
                    // alert("down");
                    var $items = $("#main_sidebar .menu a.item");
                    var itemCount = $items.length;
                    if ( 0 < itemCount ) {
                        var selectIndex = 0;
                        $.each($items, function(index, val) {
                            var $item = $(val);
                            if ($item.hasClass('active')) {
                                selectIndex = index+1;
                                return false;
                            }
                            return true;
                        });
                        if (itemCount <= selectIndex) {
                            $items.eq(0).click();
                        }else{
                            $items.eq(selectIndex).click();
                        }
                    }
                break;
                case 37:{
                    // alert("left");
                    histroyIndex--;

                    if( 1 <= histroyIndex &&  historyList.length >= histroyIndex ){
                        var history = historyList[histroyIndex-1];
                        var urlPackage = getURIQuery(QUERY_KEY_PACKAGE);
                        var urlVersion = getURIQuery(QUERY_KEY_VERSION);
                        if (urlPackage && urlVersion) {
                            if (history.key != urlPackage || history.version != urlVersion) {
                                _userRehash = false;
                                window.location.hash = history.hash;
                            }
                        }
                    }else if( 1 >= histroyIndex ){
                        histroyIndex = 1;
                    }
                }
                break;
                case 39:{
                    // alert("right");
                    histroyIndex++;
                    if( 1 <= histroyIndex && historyList.length >= histroyIndex ){
                        var history = historyList[histroyIndex-1];
                        var urlPackage = getURIQuery(QUERY_KEY_PACKAGE);
                        var urlVersion = getURIQuery(QUERY_KEY_VERSION);
                        if (urlPackage && urlVersion) {
                            if (history.key != urlPackage || history.version != urlVersion) {
                                _userRehash = false;
                                window.location.hash = history.hash;
                            }
                        }
                    }else if( historyList.length <= histroyIndex ){
                        histroyIndex = historyList.length;
                    }
                }
                break;
            }
            // console.log("histroyIndex:"+histroyIndex);
            return false;
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
                var mdpath = packageInfo.Link;
          
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

            //  set uri, 考虑到用户体验，滚动窗口不进行URI的设置，只有点击的时候才进行设置。
            // var packageName = getURIQuery(QUERY_KEY_PACKAGE);
            // var anchor = getURIQuery(QUERY_KEY_ANCHOR);
            // setURIQuery(packageName,$_item.text(),anchor,_version);

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
     *  @param scrollTop Specify window scroll top
     */
    function parsePackageMarkdown(mdpath,scrollTop){

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

                    //  set pre code highlight
                    $.each($("pre code.custom",$mainContent), function(index, val) {
                        var $_code = $(val);
                        var cssstr = $_code.attr('class');
                        var lang = "";
                        var repMap = {}; // 用户存储替换的<a>
                        var tempRep = ";temprep;_"; // 占位替换标签
                        var block = $_code.html();

                        //  使用占位符替换a标签
                        var atags = block.match(_markAtag);
                        if (atags) {
                            $.each(atags, function(index, atag) {
                                var key = tempRep+index;
                                repMap[key] = atag;
                                block = block.replace(atag,key);
                            });
                        }
                       
                        //  选取第一位的css作为语言高亮的选择
                        if (cssstr) {
                            cssstr = cssstr.split(" ");
                            if ( 1 <= cssstr.length ) {
                                lang = cssstr[0];
                            }
                        }

                        var newHtml = "";
                        if ( 0 != lang.length && hljs.getLanguage(lang) ) {
                            newHtml = hljs.fixMarkup(hljs.highlightAuto(block,[lang]).value);
                        }else{
                            newHtml = hljs.fixMarkup(hljs.highlightAuto(block).value);
                        }
                        newHtml = newHtml.replace(/&amp;/g, '&');

                        $.each(repMap, function(key, atag) {
                            newHtml = newHtml.replace(key,atag);
                        });

                        $_code.empty().html(newHtml);
                    });

                    var queryTitle = getURIQuery(QUERY_KEY_TITLE);
                    var queryAnchor = getURIQuery(QUERY_KEY_ANCHOR);

                    var titleScrollTop = false;

                    //  each #main_content element child node h2 tag
                    var h2Items = {}
                    $.each($("h2",  $mainContent), function(index, val) {
                        var $stickyItem = $(val);
                        var offsetTop = $stickyItem.offset().top - NAVIGATION_HEIGHT;  //  top navigation height
                        h2Items[offsetTop] = val;
                    });

                    //  save sticky items
                    var stickyItems = {};
                    $.each(h2Items, function(key, val) {
                        var text = $(val).text();
                        var $item = $('<a class="item" href="javascript:;">'+text+'</a>');

                        $item.click(function() {
                            var queryPackage = getURIQuery(QUERY_KEY_PACKAGE);
                            setURIQuery(queryPackage,text,false,_version);
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

                                    if ( '#' != text ) {
                                        var $anchorTag = $("*[id='"+anchor+"']",$mainContent);
                                        if ( 0 == $anchorTag.length ) {
                                            $anchorTag = $("*[name='"+anchor+"']",$mainContent);
                                        }
                                        if ( 0 != $anchorTag.length ) {
                                            $anchorTag = $anchorTag.eq(0);
                                            var offsetTop = $anchorTag.offset().top - NAVIGATION_HEIGHT;
                                            $(window).scrollTop(offsetTop);
                                        }
                                    }else{
                                        setURIQuery(queryPackage,false,unescape(anchor),_version);
                                        // copyToClipboard(window.location.href);
                                    }
                                    return false;
                                });

                            }else if ( 0 == href.indexOf("src.html") ){
                                // target="_blank"
                                $atag.attr('target', "_blank");
                            }else if ( 0 == href.indexOf("http://") || 0 == href.indexOf("https://") ){
                                $atag.attr('target', "_blank");
                            }else if( 0 == href.indexOf('../') ){
                                var newHref = "javascript:;";
                                var isBlank = true;
                                /*
                                    github:
                                    mdUrl   = https://github.com/../project/doc/v1_0_0/md/default/github.com/slowfei/gosfdoc.md
                                    srcFile = https://github.com/../project/doc/v1_0_0/src/default/github.com/slowfei/gosfdoc.go

                                    to src dir:
                                    href = ../../../../src/gosfdoc.go#L10-L16
                                    guthub result  = https://github.com/.../project/doc/v1_0_0/src/gosfdoc.go#L10-L16
                                    gosfdoc result = http://.../project/doc/src.html?v=1.0.0&f=gosfdoc.go&L=10-16
                                    
                                    to root dir:
                                    href = ../../../../../../../gosfdoc.go#L10-L16
                                    guthub result  = https://github.com/.../project/gosfdoc.go#L10-L16
                                    gosfdoc result = http://.../project/doc/src.html?v=1.0.0&f=gosfdoc.go&L=10-16
                                */
                                var tempHref = href.replace(/\.\.\//g,"");
                                if ( 0 == tempHref.indexOf("src/")) {
                                    tempHref = tempHref.substring(4,href.length);
                                }

                                var tempSplit = tempHref.split("#");
                               
                                if ( 1 == tempSplit.length) {

                                    var fileName = tempSplit[0];
                                    if ( 0 != _appendPath.length && 0 != fileName.indexOf(_appendPath)){
                                        if ('/' == _appendPath[_appendPath.length - 1]) {
                                            fileName = _appendPath + fileName;
                                        }else{
                                            fileName = _appendPath + "/" +fileName;
                                        }
                                    }
                                        
                                    if (MD_FILE_SUFFIX == fileName.substr(fileName.length-3,fileName.length)) {
                                        
                                        $atag.click(function(event) {
                                            setURIQuery(fileName,false,false,_version);
                                            parsePackageMarkdown(fileName);
                                            return false;
                                        });

                                        newHref = "index.html?v=" + _version + "&p=" + fileName;
                                        isBlank = false;
                                    }else{
                                        newHref = "src.html?v=" + _version + "&f=" + fileName;
                                    }

                                }else if ( 2 == tempSplit.length ){
                                    //  ../../gosfdoc.go#L10-L16 存在参数的处理

                                    var fileName = tempSplit[0];
                                    if ( 0 != _appendPath.length && 0 != fileName.indexOf(_appendPath)){
                                        if ('/' == _appendPath[_appendPath.length - 1]) {
                                            fileName = _appendPath + fileName;
                                        }else{
                                            fileName = _appendPath + "/" +fileName;
                                        }
                                    }

                                    if (MD_FILE_SUFFIX == fileName.substr(fileName.length-3,fileName.length)) {
                                        var anchor = tempSplit[1];
                                        
                                        $atag.click(function(event) {
                                            setURIQuery(fileName,false,anchor,_version);
                                            parsePackageMarkdown(fileName);
                                            return false;
                                        });

                                        newHref = "index.html?v=" + _version + "&p=" + fileName + "&a=" + anchor;
                                        isBlank = false;
                                    }else{
                                        var lineParam = "";
                                        var tempLines = tempSplit[1];

                                        //  tempLines = L10-L16 or L10
                                        //  get L10 to 10
                                        var l2 = tempLines.replace(_rexLineGithub,"$2");
                                        if ( l2 != tempLines) {
                                            if ( 0 == l2.length ) {
                                                //  get L10 to 10
                                                //  get L16 to 16
                                                var l3 = tempLines.replace(_rexLineGithub,"$3");
                                                var l4 = tempLines.replace(_rexLineGithub,"$4");

                                                lineParam = "&L=" + l3 + "-" + l4;
                                            }else{
                                                lineParam = "&L="+l2;
                                            }
                                        }
                                        newHref = "src.html?v=" + _version + "&f=" + fileName + lineParam;
                                    }
                                }

                                $atag.attr('href', newHref);
                                if (isBlank) {
                                    $atag.attr('target', "_blank");
                                }
                            }else{
                                $atag.click(function(event) {
                                    return false;
                                });
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
                            titleScrollTop = $anchorTag.offset().top - NAVIGATION_HEIGHT;
                        }
                    };

                    //  scroll to func position
                    if (scrollTop) {
                        $(window).scrollTop(scrollTop);
                    }else if (titleScrollTop) {
                        $(window).scrollTop(titleScrollTop);
                    }

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
                var mdpath = packageInfo.Link;

                if ( findPackage ) {
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
            var classActive = "";

            //  默认选择排序在第一位的版本
            if ( 0 == tempOne.length) {
                tempOne = verstr;
            }
            
            //  效验设置的版本信息是否与获取的版本信息相符
            if ( _version == verstr ) {
                checkVer = true;
                classActive = " active";
            }   

            var html = '<div class="item'+classActive+'" data-value="'+verstr+'">'+verstr+'</div>';
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
                }else{
                     parseMenuList(dataGosfdocJson.Markdowns,false);
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

        _userRehash = true;

        window.location.hash = hash;
        window.location.search ='';
        window.location.query = '';

        window.clearTimeout(_rehashTimeout);
        _rehashTimeout = window.setTimeout(function(){
            _userRehash = false;
        },600);
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

    /**
     *  copy text to clipboard
     *
     *  @param txt
     */
    function copyToClipboard(txt) {
        //  TOOD 由于复制到剪切板有些复杂，暂时这样。
        window.prompt("Copy link to clipboard?", txt);
    }

    /**
     *
     */
    function p_escape(html,encode) {
        return html.replace(!encode ? /&(?!#?\w+;)/g : /&/g, '&amp;')
                   .replace(/</g, '&lt;')
                   .replace(/>/g, '&gt;')
                   .replace(/"/g, '&quot;')
                   .replace(/'/g, '&#39;');
    }

    /**
     *  
     */
    function p_unescape(html) {
        return html.replace(/&amp;/g, '&')
                   .replace(/&lt;/g, '<')
                   .replace(/&gt;/g, '>')
                   .replace(/&quot;/g, '"')
                   .replace(/&#39;/g, "'");
    }

})(jQuery);