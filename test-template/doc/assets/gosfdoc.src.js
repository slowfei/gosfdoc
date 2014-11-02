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
    var $_preCode = null;
    var _codeArrayKey = "CodeArray";
    var _windowWhich = -1;
    var _rexLine = /^(#L\d+|#L\d+[-]L\d+)$/;
    var _sourceDir = "src/";
    
    $(function(){

        init();

        hljs.configure({tabReplace: '    ',useBR:false});
        $_preCode = $("#main_content pre code");

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

        var filepath = getURIQuery("f");
        if(!filepath){
            $_preCode.html("source code browse, please select left menu item.");
        }else{
            //  init load file
            loadSourceCode(filepath);
            //  set left sidebar item
            $("#main_sidebar .item a.item[href='"+filepath+"']").addClass('active');
        }

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
            $("#menu_title").html("<center>Files</center>");

            var $sidebarItem = $('<div class=item></div>');
            //  each markdowns
            $.each(dataJson.Files, function(index, val) {

                $.each(val, function(projectName, items) {

                    var $menu = $('<div class="menu"></div>');
                    var $itemTitle = $('<b>'+projectName+'</b>');

                    $.each(items, function(pkeKey, fileInfo) {
                        var fileName = fileInfo.Filename;
                        var linkHref = fileName;
                        
                        var $item = $('<a class="item" href='+linkHref+'>'+fileName+'</a>');

                        $item.click(function() {
                            if ($(this).hasClass('active')) {return false;}
                            window.location.hash = '';
                            window.location.search = 'f='+$(this).attr('href');
                            return false;
                        });

                        $menu.append($item);
                    });

                    $sidebarItem.append($itemTitle);
                    $sidebarItem.append($menu);
                   
                    return false;

                });// end  $.each(val, function(projectName, items)

            });// end  $.each(dataJson.Markdowns, function(index, val)

            $("#main_sidebar").append($sidebarItem);
       })
       .fail(function(jqXHR, textStatus, errorThrown) {
             $("#menu_title").html('<center>Load "config.json" Error.</center>');
       });
    }
    
    /**
     *  load source code file
     *
     *  @param filepath file url path
     */
    function loadSourceCode(filepath){
        $_preCode.empty();
        // $_preCode.removeData(_codeArrayKey);
        $_preCode.html("Loading...");
        $.ajax({
           url: _sourceDir+filepath,
           async:true,
           cache:false,
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

            //  location line
            var currentHash = window.location.hash;
            if (_rexLine.test(currentHash)) {
                currentHash = currentHash.replace(/L|#/g,'');
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
                appendstr = "\nuse gosfdoc command run document.";
            };

            $_preCode.html("sorry! source code load failed." + appendstr);
        });
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
            var hash = ''

            if (!cancel) {
                var currentHash = window.location.hash;

                if ( -1 == _windowWhich || !_rexLine.test(currentHash) ) {
                    hash = 'L'+currentLine;
                    activeLines.push(currentLine);
                }else{

                    cancelTextlUserSelect();

                    currentHash = currentHash.replace(/L|#/g,'');
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
                    hash = 'L'+ startLine + '-L' + endLine;
                }   
            }
            selectLineNumber(activeLines);
            window.location.hash = hash;
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