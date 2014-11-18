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
//  Create on 2014-10-28
//  Update on 2014-11-14
//  Email  slowfei#foxmail.com
//  Home   http://www.slowfei.com

//  html file constant content
package assets

const HTML_INDEX = `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">

    <title></title>
    
    <link rel="stylesheet" type="text/css" href="assets/gosfdoc.min.css">

    <script src="assets/assets.min.js"></script>
    <script src="assets/gosfdoc.min.js"></script>
    
</head>
<body class="side pushed" ontouchstart>
    <div id="main_sidebar" class="ui large floating vertical inverted labeled sidebar menu active">
        <div class="item" id="menu_title">menu title</div>
        
    </div>
    <div id="btn_show_menu" class="ui launch black right attached button">
            <span class="text"><</span>
    </div>
    
    <div class="ui fixed transparent inverted menu topfixed">
        <div class="container">
            <a id="home" class="left item" href="javascript:;">&nbsp;</a>
            <div class="title item" id="doc_title">doc title</div>
        </div>
    </div>

    <div class="ui menu transparent right_menu">

        <div class="section ui dropdown link item version">
            <div class="textcolor">
                <b>Version:</b> 
                <span id="version_text"></span>
            </div>
            <div id="version_value" class="menu ui transition">
            </div>
        </div>
        
        <div class="section ui dropdown link item language">
            <div class="textcolor">
                <b>Language:</b> 
                <span id="language_text"></span>
            </div>
            <div id="language_value" class="menu ui transition">
            </div>
        </div>

        <div class="item title">
            <a href="javascript:;" id="item_about"><b>About</b></a>
            <div class="ui modal about" id="about_modal">
                <div class="content md common_md" id="about_content">

                </div>
                <div class="actions">
                    <div class="ui button">
                      Close
                    </div>
                </div>
            </div>
        </div>
         <div class="item">&nbsp;</div>
     </div>
    
    
    <div class="segment intro">
        <div class="container">
            <div id="segment_intro">

            </div>
        </div>  
    </div>
    
    <div class="main container">
        <div class="sticky-wrapper" id="sticky-wrapper">
            <div id="sticky" class="ui vertical pointing secondary menu">

            </div>
        </div>

        <div class="md common_md ui sticky" id="main_content">
            
        </div>
    </div>
  
    <br/><br/><br/><br/><br/><br/><br/><br/><br/><br/><br/>
    <br/><br/><br/><br/><br/><br/><br/><br/><br/><br/><br/>
    <br/><br/><br/><br/><br/><br/><br/><br/><br/><br/><br/>
    <br/><br/><br/><br/><br/><br/><br/><br/><br/>
</body>
</html>`

const HTML_SRC = `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">

    <title>source code</title>

    <link rel="stylesheet" type="text/css" href="assets/gosfdoc.min.css">

    <script src="assets/assets.min.js"></script>
    <script src="assets/gosfdoc.src.min.js"></script>

</head>
<body class="side pushed src" ontouchstart>
    <div id="main_sidebar" class="ui large floating vertical inverted labeled sidebar menu active">
        <div class="item" id="menu_title">
            <center>Files</center>
        </div>
        
    </div>
    <div id="btn_show_menu" class="ui launch black right attached button src">
            <span class="text"><</span>
    </div>
    
    <div class="main container src">

        <div class="ui menu transparent right_menu">
            <div class="section ui dropdown link item version">
                <div class="textcolor">
                    <b>Version:</b> 
                    <span id="version_text"></span>
                </div>
                <div id="version_value" class="menu ui transition">
                </div>
            </div>
        </div>

        <div id="src_info">
        </div>
        <div class="ui sticky src" id="main_content">
            <pre><code>source code browse, please select left menu item.</code></pre>
        </div>
    </div>

</body>
</html>`
