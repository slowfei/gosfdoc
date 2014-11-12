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
    var GOSFDOC_JSON = "config.json";
    var QUERY_KEY_VERSION = "v";
    var QUERY_KEY_FILE = "f";
    var QUERY_KEY_LINES = "L";
    var VERSION_DEFAULT = "1.0.0";

    var dataGosfdocJson = null;
    var $_preCode = null;
    var _version = null;
    var _windowWhich = -1;
    var _linkRoot = false;
    var _appendPath = "";
    var _rexLine = /^(\d+|\d+[-]\d+)$/;
    
    $(function(){

        hljs.configure({tabReplace: '    ',useBR:false});
        $_preCode = $("#main_content pre code");

        //
        init();

        //
        initUI();

        //
        initEvent();

    });

    /**
     *  init config
     */
    function init(){
        
        $("#menu_title").html("<center>Loading...</center>");

        $.ajax({
           url: GOSFDOC_JSON,
           async:false,
           cache:false,
           type: 'GET',
           dataType: 'json'
       })
       .done(function(dataJson) {
            dataGosfdocJson = dataJson;
             $("#menu_title").html("<center>Files</center>");

            //  handle other params
            if (dataJson.LinkRoot) {
                _linkRoot = dataJson.LinkRoot;
            }
            if (dataJson.AppendPath){
                _appendPath = dataJson.AppendPath;
            }

            parseSelectVersion(dataJson.Versions);
            parseMenuList(dataJson.Files);

            var filepath = getURIQuery(QUERY_KEY_FILE);
            var lines = getURIQuery(QUERY_KEY_LINES);

            //  _version 已经在parseSelectVersion中进行设置
            setURIQuery(filepath,lines,_version);

            //  load code
            loadSourceCode(filepath);
       })
       .fail(function(jqXHR, textStatus, errorThrown) {
             $("#menu_title").html('<center>Load "config.json" Error.</center>');
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
            $('#main_sidebar').sidebar({
                duration:1,
            }).sidebar('toggle');
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

        //  check browser set css
        if ( - 1 == navigator.userAgent.indexOf('Chrome') && -1 != navigator.userAgent.indexOf('Safari')) {
            $("#number").addClass('safari');
        } else if ( - 1 != navigator.userAgent.indexOf('Opera')) {
            $("#number").addClass('opera');
        } else if ( - 1 != navigator.userAgent.indexOf('IE')) {
            $("#number").addClass('ie');
        }
    }

    /**
     *  init event
     */
    function initEvent(){
        //  keydown
        $(window).keydown(function(event) {
            if ( 1 <= $(".line_number.active",$_preCode).length ) {
                _windowWhich = event.which;
            }else{
                _windowWhich = -1;
            }
        }).keyup(function(event) {
            _windowWhich = -1;
        });


    }
    
    /**
     *  load source code file
     *
     *  @param filepath file url path
     */
    function loadSourceCode(filepath){
        if (!filepath) {return}

        //  set left sidebar item active
        var $tempItems =  $("#main_sidebar .item a.item");
        $.each($tempItems, function(index, val) {
            var $item = $(val);
            var href = $item.attr('href');
            if (href == filepath) {
                $item.addClass('active');
            }else{
                $item.removeClass('active');
            }
        });
        // $("#main_sidebar .item a.item[href='"+filepath+"']").addClass('active');

        $_preCode.empty();
        $_preCode.html("Loading...");
        
        //  handle url path
        
        var urlPath = "";
        if (_linkRoot) {
            //  e.g. : ../temp.go

            if ( 0 != _appendPath.length && 0 == filepath.indexOf(_appendPath+"/")) {
                filepath = filepath.substring(_appendPath.length+1,filepath.length);
            }
            urlPath = "../" + filepath;

        }else{
            //  e.g. : v1_0_0/src/temp.go

            var verPath = converToVersionPath(_version);

            if ( 0 != _appendPath.length && 0 != filepath.indexOf(_appendPath)) {
                filepath = _appendPath + "/" + filepath;
            }

            urlPath =  verPath +"/src/"+filepath;
        }
        
        $.ajax({
           url: urlPath,
           async:true,
           cache:true,
           dataType: 'text'
        }).done(function(text){
            $_preCode.attr('src', filepath);
            $_preCode.empty();
            var lang = filepath.substring(filepath.lastIndexOf('.')+1,filepath.length);

            var codehtml = '';
            if ( lang && hljs.getLanguage(lang) ) {
                codehtml = hljs.fixMarkup(hljs.highlightAuto(text,[lang]).value);
            }else{
                codehtml = hljs.fixMarkup(hljs.highlightAuto(text).value);
            }

            var codeArray = new Array();
            var codeLines = codehtml.split(/\n/g);
            var linesCount = codeLines.length;

            for (var i = 0; i < linesCount; i++) {
                var codeLine = codeLines[i];
                if ( !codeLine || 0 == codeLine.length ) { codeLine = "&nbsp;"};

                if ( -1 != codeLine.indexOf('>/*') &&  -1 == codeLine.indexOf('*/<') ) {
                    // handle <span class="hljs-comment">/** 
                    var space = '';
                    var spaceCount = codeLine.indexOf('<');
                    for (var j = 0; j < spaceCount; j++) {
                        space += ' ';
                    }
                    if ( -1 != spaceCount ) {
                         codeLine = codeLine.substring(spaceCount,codeLine.length);
                    };
                    codeLine = codeLine.replace('>/*','><span class="line_code"><a class="line_number" L="'+(i+1)+'"></a>'+space+'/*') + "\n</span>";
                    
                }else if( -1 != codeLine.indexOf('*/<') && -1 == codeLine.indexOf('>/*') ){
                    //  handle */</span> 
                    codeLine = '<span class="line_code"><a class="line_number" L="'+(i+1)+'"></a>' + codeLine.replace('*/<','*/\n</span><');

                }else{
                    codeLine = '<span class="line_code"><a class="line_number" L="'+(i+1)+'"></a>'+codeLine+'\n</span>';
                }

                codeArray.push(codeLine);
            }

            if ( 1500 < linesCount ) {
                $("body").addClass('transition_initial');
                $("#main_sidebar").addClass('transition_initial');
            }else{
                 $("body").removeClass('transition_initial');
                $("#main_sidebar").removeClass('transition_initial');
            }

            $("#src_info").html(linesCount + ' lines '+ text.getBytesLength()  + ' byte');

            $_preCode.html(codeArray.join(''));
            handleLineNumber();

            //  set active lines
            var currentHash = getURIQuery(QUERY_KEY_LINES);
            if ( currentHash && _rexLine.test(currentHash)) {
                currentHash = currentHash.split('-');
                var L1 =  parseInt(currentHash[0]);
                var L2 = currentHash[1];
                var startLine = L1;
                var endLine = L1;
                if (L2) {
                    endLine = parseInt(L2);
                };

                var activeLines = new Array();
                for (var i = startLine; i <= endLine; i++) {
                    activeLines.push(i+'');
                }
                selectLineNumber(activeLines,L1+'');
            }

        }).fail(function() {
            var appendstr = "";
            if ( -1 != window.location.protocol.indexOf('file') ) {
                appendstr = "\n\nTry to use gosfdoc command run document.";
            };
            $_preCode.html("File: "+ filepath + "\nSorry! Unable to load the source code. version: "+ _version + appendstr);
        });
    }

    /**
     *  parse version select info
     *
     *  @param verJson
     */
    function parseSelectVersion(verJson){

        var ver = getURIQuery(QUERY_KEY_VERSION);

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
                if ( _version == value) {
                    return;
                }

                $("#version_text").text(text);
                _version = value;

                var filepath = getURIQuery(QUERY_KEY_FILE);
                var lines = getURIQuery(QUERY_KEY_LINES);
                setURIQuery(filepath,lines,_version);

                if (dataGosfdocJson) {
                    parseMenuList(dataGosfdocJson.Files);
                }

                loadSourceCode(filepath);
            }
        });
    }

    /**
     *  parse left menu list item
     *
     *  @param filesJson
     */
    function parseMenuList(filesJson){

        var $sidebarItem = $('<div class=item></div>');
        var addCount = 0;

        //  each Files
        $.each(filesJson, function(index, val) {
            var proVer = val.Version;

            if ( _version != proVer ) {
                return true;
            }

            var projectName = val.MenuName;
            var $menu = $('<div class="menu"></div>');
            var $itemTitle = $('<b>'+projectName+'</b>');

            $.each(val.List, function(listIndex, fileInfo) {
                var fileName = fileInfo.Filename;
                var linkHref = fileName;
                
                var $item = $('<a class="item" href='+linkHref+'>'+fileName+'</a>');

                $item.click(function() {
                    if ($(this).hasClass('active')) {return false;}

                    var filepath = $(this).attr('href');
                    setURIQuery(filepath,false,_version);

                    loadSourceCode(filepath);

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
    }

    /**
     *  handle code line number element
     */
    function handleLineNumber(singleLine){
        var $lineNumber = null;
        if (singleLine) {
            $lineNumber = $(singleLine);
        }else{
            $lineNumber = $(".line_number",$_preCode);
        }
        $lineNumber.unbind();

        $lineNumber.hover(function() {
            var $parentSapn = $(this).parent();
            $parentSapn.css('background-color', 'rgba(243, 104, 11, 0.5)');
        }, function() {
            var $parentSapn = $(this).parent();
            $parentSapn.css('background-color', 'transparent');
        });

        $lineNumber.click(function(event) {

            var $_thisNumber = $(this);
            var cancel = false;
            if ($_thisNumber.hasClass('active')) {
                cancel = true;
            }
            var currentLine = $(this).attr('L');
            var activeLines = new Array();
            var lines = '';

            if (!cancel) {
                var currentHash = getURIQuery(QUERY_KEY_LINES);

                if ( -1 == _windowWhich || !_rexLine.test(currentHash) ) {
                    lines = currentLine;
                    activeLines.push(currentLine);
                }else{

                    cancelTextlUserSelect();

                    currentHash = currentHash.split('-');
                    var L1 =  parseInt(currentHash[0]);
                    currentLine = parseInt(currentLine);

                    var startLine = 0;
                    var endLine = 0;
                    if ( L1 < currentLine ) {
                        startLine = L1;
                        endLine = currentLine;
                    }else{
                        startLine = currentLine;
                        endLine = L1;
                    }

                    for (var i = startLine; i <= endLine; i++) {
                             activeLines.push(i+'');
                    }
                    lines =  startLine + '-' + endLine;
                }   
            }
            selectLineNumber(activeLines);

            setURIQuery(getURIQuery(QUERY_KEY_FILE),lines,_version);
        }); // end $lineNumber.click(function(event)
    }

    /**
     *  cancel current user select text
     */
    function cancelTextlUserSelect(){
        var sel = window.getSelection ? window.getSelection() : document.selection;
        if (sel) {
            if (sel.removeAllRanges) {
                sel.removeAllRanges();
            } else if (sel.empty) {
                sel.empty();
            }
        }
    }

    /**
     *  select line active
     *
     *  @param activeLines active line number
     *  @param scrollNumber scroll line top
     */
    function selectLineNumber(activeLines,scrollNumber){
         $.each($(".line_number",$_preCode), function(index, val) {
            var $valLineNumber = $(val);
            var line = $valLineNumber.attr('L');

            if ( scrollNumber && scrollNumber == line) {
                var top = $valLineNumber.offset().top - $(window).height() * 0.2;
                $(window).scrollTop(top);
            };
            if ( -1 != $.inArray(line, activeLines) ) {
                if ( !$valLineNumber.hasClass('active') ) {
                    $valLineNumber.addClass('active');
                }
            }else{
                $valLineNumber.removeClass('active');
            }
        });
    }

    /**
     *  set fixed uri query
     *
     *  @param filepath
     *  @param lines
     *  @param version
     */
    function setURIQuery(filepath,lines,version){
        var hash = "";

        if ( version ){
            hash += QUERY_KEY_VERSION + '=' + escape(version);
        }

        if ( filepath ) {
            if ( 0 != hash.length ){
                hash += '&';   
            }
            hash += QUERY_KEY_FILE + '=' + escape(filepath);
        }

        if ( lines ) {
            if ( 0 != hash.length ){
                hash += '&';   
            }
            hash += QUERY_KEY_LINES + '=' + escape(lines);
        }

        window.location.hash = hash;
        window.location.search ='';
        window.location.query = '';
    }

    /**
     *  string byte length  
     *
     *  @param byte length
     */
    String.prototype.getBytesLength = function() { 
        var totalLength = 0;     
        var charCode;  
        for (var i = 0; i < this.length; i++) {  
            charCode = this.charCodeAt(i);  
            if (charCode < 0x007f)  {     
                totalLength++;     
            } else if ((0x0080 <= charCode) && (charCode <= 0x07ff))  {     
                totalLength += 2;     
            } else if ((0x0800 <= charCode) && (charCode <= 0xffff))  {     
                totalLength += 3;   
            } else{  
                totalLength += 4;   
            }          
        }  
        return totalLength;   
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

})(jQuery);