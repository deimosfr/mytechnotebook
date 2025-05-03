---
weight: 999
url: "/Les_extentions_pratiques_de_Mediawiki/"
title: "Practical Extensions for MediaWiki"
description: "A comprehensive guide to useful MediaWiki extensions including Google Analytics, Search, Ads, Syntax highlighting, and other useful tools to enhance functionality."
categories: ["Debian", "Linux"]
date: "2013-08-02T15:04:00+02:00"
lastmod: "2013-08-02T15:04:00+02:00"
tags: ["Google Search", "CharInsert", "Cite", "SphinxSearch", "View source", "search", "Google Adsense", "What links here", "Servers"]
toc: true
---

## Introduction

MediaWiki has extensions that add interesting functionality. The only downside of these extensions is that they can become obsolete with version changes. It's up to you to decide...however, some are only used occasionally.

## Extensions

### SpecialRenameuser.php

This extension is provided in Debian repositories, so I recommend installing it with apt. Otherwise, it's [available here](https://www.mediawiki.org/wiki/Extension:Renameuser). Then you just need to link it:

```bash
cd /etc/mediawiki-extensions/extensions-enabled
ln -s /usr/share/mediawiki-extensions/SpecialRenameuser.php ./
```

Now, if you go to the Special pages, you'll have the ability to rename a user.

### DeleteHistory

I developed [this extension (DeleteHistory)](https://www.mediawiki.org/wiki/Extension:DeleteHistory) myself. It allows you to delete the history of your articles. Besides leaving some traces (passwords, confidential information...), you might wonder why not keep the history.

I had the idea to create it because I previously used [SpecialDeleteOldRevisions](https://www.mediawiki.org/wiki/Extension:SpecialDeleteOldRevisions), but the developer wasn't reactive enough for my taste to update his extension and make it work everywhere. Moreover, he developed the entire deletion part himself, which certainly has advantages, but also a major disadvantage: it's not maintained by the MediaWiki team, and if changes to the database happen during an upgrade, the extension must be modified to accommodate these changes.

For my extension, I simply opted for a maintenance script provided by MediaWiki, which removes the problem mentioned above and allows me to be compatible with many different versions.

Depending on your wiki's purpose, this may or may not be important. But for mine, it's not very important!

[This extension is available here](https://www.mediawiki.org/wiki/Extension:DeleteHistory). For those of you who are less interested, let me give you some numbers to consider. In one year of wiki usage, my database was 14 MB. With this extension that removes histories, I was able to reduce it to 4.8 MB. You can see the gain you can achieve by removing histories!

### Google Ads

For [Google Adsense](https://www.google.com/adsense), you must first create an account and then generate ad banners of your desired size.

#### Without skin modification

You need to install the [PCR GUI Inserts](https://www.mediawiki.org/wiki/Extension:PCR_GUI_Inserts) module which will allow you to insert information at various places on your pages. First, activate the extension by adding these lines:

```php
# PCR Extension for Piwik / Google Ads
require_once("$IP/extensions/pcr/pcr_guii.php");
```

##### Bottom section

Here are the lines to add to have a bottom section with Google ads:

```php
# Google Ads (bottom page)
$wgPCRguii_Inserts['SkinAfterBottomScripts']['on'] = true;
$wgPCRguii_Inserts['SkinAfterBottomScripts']['content'] = '
<!-- Graphically aligned -->
<div style="margin-left: 160px;">
 
<!-- Google Ads 1 -->
<script type="text/javascript"><!--
google_ad_client = "*****PLACE YOUR CLIENT ID HERE******";
/* 728x90 */
google_ad_slot = "***CODE HERE***";
google_ad_width = 728;
google_ad_height = 90;
//-->
</script>
<script type="text/javascript"
src="http://pagead2.googlesyndication.com/pagead/show_ads.js">
</script>
<!-- End Google Ads 1 -->
 
</div>
';
```

#### With skin modification

The disadvantage of this method is that when you update MediaWiki, these modifications are likely to be overwritten. You will then have to redo the configuration. Opt for the manual method if you want to be worry-free.

##### Monobook

###### Bottom section

Insert these lines in the file **mediawiki/skins/Monobook.php**:

```php
<script type="text/javascript"><!--
google_ad_client = "*****PLACE YOUR CLIENT ID HERE******";
/* 728x90 */
google_ad_slot = "***CODE HERE***";
google_ad_width = 728;
google_ad_height = 90;
//-->
</script>
<script type="text/javascript"
src="http://pagead2.googlesyndication.com/pagead/show_ads.js">
</script>
```

To be added below:

```php
                </div><!-- end of the left (by default at least) column -->
                        <div class="visualClear"></div>
                        <div id="footer">
```

###### Left section

```php
<div id="custom-advert" class ="portlet">
<h5>Google Ads</h5>
<div class = "pBody">
<script type="text/javascript"><!--
google_ad_client = "*****PLACE YOUR CLIENT ID HERE******";
/* Vertical text */
google_ad_slot = "***CODE HERE***";
google_ad_width = 120;
google_ad_height = 600;
//-->
</script> 
<script type="text/javascript"
src="http://pagead2.googlesyndication.com/pagead/show_ads.js">
</script>
</div></div>
```

And place it above this section:

```php
                </div><!-- end of the left (by default at least) column -->
                        <div class="visualClear"></div>
                        <div id="footer">
<?php
                if($this->data['poweredbyico']) { ?>
                                <div id="f-poweredbyico"><?php $this->html('poweredbyico') ?></div>
<?php   }
                if($this->data['copyrightico']) { ?>
                                <div id="f-copyrightico"><?php $this->html('copyrightico') ?></div>
```

##### Vector

###### Bottom section

Insert these lines in the file **mediawiki/skins/Vector.php**:

```php {linenos=table,hl_lines=["2-3","9-10"]}
<script type="text/javascript"><!--
google_ad_client = "*****PLACE YOUR CLIENT ID HERE******";
google_ad_slot = "***CODE HERE***";
google_ad_width = 728;
google_ad_height = 90;
//-->
</script>
<script type="text/javascript"
src="http://pagead2.googlesyndication.com/pagead/show_ads.js">
</script>
<script type="text/javascript"><!--
google_ad_client = "*****PLACE YOUR CLIENT ID HERE******";
google_ad_slot = "***CODE HERE***";
google_ad_width = 728;
google_ad_height = 90;
//-->
</script>
<script type="text/javascript"
src="http://pagead2.googlesyndication.com/pagead/show_ads.js">
</script>
```

To add in this section:

```php
<?php           endforeach; ?>
                </ul>
            <?php endif; ?>
            <div style="clear:both"></div>
 
/* THE CODE HERE!!! */ 
        </div>
        <!-- /footer -->
        <!-- fixalpha -->
        <script type="<?php $this->text('jsmimetype') ?>"> if ( window.isMSIE55 ) fixalpha(); </script>
        <!-- /fixalpha -->
```

##### Left section

Insert these lines in the file **mediawiki/skins/Vector.php**:

```php
<script type="text/javascript"><!--
google_ad_client = "*****PLACE YOUR CLIENT ID HERE******";
google_ad_slot = "***CODE HERE***";
google_ad_width = 120;
google_ad_height = 600;
//-->
</script>
<script type="text/javascript"
src="http://pagead2.googlesyndication.com/pagead/show_ads.js">
</script>
```

To add in this section:

```php
                <!-- logo -->
                    <div id="p-logo"><a style="background-image: url(<?php $this->text( 'logopath' ) ?>);" href="<?php echo htmlspecialchars( $this->data['nav_urls']['mainpage']['href'] ) ?>" <?php echo Xml::expandAttributes( Linker::tooltipAndAccesskeyAttribs( 'p-logo' ) ) ?>></a></div>
                <!-- /logo -->
                <?php $this->renderPortals( $this->data['sidebar'] ); ?>
 
/* THE CODE HERE!!! */ 
            </div>
        <!-- /panel -->
        <!-- footer -->
        <div id="footer"<?php $this->html( 'userlangattributes' ) ?>>
```

### Google Search

#### Without skin modification

I created an extension that doesn't require modifying the skin: [https://www.mediawiki.org/wiki/Extension:GoogleSearch](https://www.mediawiki.org/wiki/Extension:GoogleSearch)  
All explanations are on the MediaWiki page and it's very simple, so I won't rewrite those lines.

#### With skin modification

Implementing Google search is very practical, as the basic search isn't always the best. To have a small Google search at the bottom left of your menus, edit the file **mediawiki/skins/Monobook.php**. You need to place and adapt (obviously you need a [Google Adsense](https://www.google.com/adsense) account) this text:

```php
 <div id="custom-advert" class ="portlet">
     <h5>Google Search</h5>
     <div class = "pBody">
<!-- SiteSearch Google -->
<form method="get" action="http://www.google.fr/custom" target="_top">
<table border="0" bgcolor="#ffffff">
<tr><td nowrap="nowrap" valign="top" align="left" height="32">
 
<input type="hidden" name="domains" value="www.deimos.fr"></input>
<label for="sbi" style="display: none">Enter the terms you wish to search for.</label>
<input type="text" name="q" size="12" maxlength="255" value="" id="sbi"></input>
</td></tr>
<tr>
<td nowrap="nowrap">
<table>
<tr>
<td>
<input type="radio" name="sitesearch" value="" id="ss1"></input>
<label for="ss1" title="Search the Web"><font size="-2" color="#000000">Web</font></label>
<br />
<input type="radio" name="sitesearch" value="www.deimos.fr" checked id="ss0"></input>
<label for="ss0" title="Search www.deimos.fr"><font size="-2" color="#000000">www.deimos.fr</font></label></td>
</tr>
</table>
<label for="sbb" style="display: none">Submit search form</label>
<input type="submit" name="sa" value="Search" id="sbb"></input>
<input type="hidden" name="client" value="******CLIENT_ID******"></input>
<input type="hidden" name="forid" value="1"></input>
<input type="hidden" name="ie" value="ISO-8859-1"></input>
<input type="hidden" name="oe" value="ISO-8859-1"></input>
<input type="hidden" name="cof" value="GALT:#008000;GL:1;DIV:#336699;VLC:663399;AH:center;BGC:FFFFFF;LBGC:336699;ALC:0000FF;LC:0000FF;T:000000;GFNT:0000FF;GIMP:0000FF;FORID:1"></input>
<input type="hidden" name="hl" value="fr"></input>
</td></tr></table>
</form>
<!-- SiteSearch Google -->
 
      </div></div>
```

To add below:

```php
                </div><!-- end of the left (by default at least) column -->
                        <div class="visualClear"></div>
                        <div id="footer">
```

### Google Reader

First activate the feeds you want, then add these lines to Monobook.php:

```
 <div id="custom-advert" class ="portlet">
        <h5>Watch list</h5>
        <div class = "pBody">
<table border="0"><tr><td height="16" align="left">
<script type="text/javascript" src="http://www.google.fr/reader/ui/publisher-fr.js"></script>
<script type="text/javascript" src="http://www.google.fr/reader/public/javascript/user/00168335903816875978/state/com.google/starred?n=8&amp;callback=GRC_p(%7Bc%3A%22blue%22%2Ct%3A%22%22%2Cs%3A%22false%22%2Cb%3A%22false%22%7D)%3Bnew%20GRC"></script>
</td></tr></table>
</div></div>
```

Adapt these lines to your needs. You need to place this between the two Google searches and you're done :).

### Google Analytics

For this plugin, I recommend first looking at the excellent [documentation on the official MediaWiki site](https://www.mediawiki.org/wiki/Extension:Google_Analytics_Integration#Installation). There's nothing to add :-)

### Group Based Access Control

This extension is used to set up ACLs on MediaWiki. Which is very practical :-). Here's the link:  
[https://www.mediawiki.org/wiki/Extension:Group_Based_Access_Control](https://www.mediawiki.org/wiki/Extension:Group_Based_Access_Control)

In this example (to be placed on one of the pages):

```
<accesscontrol>Administrators,,IT-Department,,Sales(ro)</accesscontrol>
```

* Administrators and IT-Department have all rights
* Sales only has read rights

All of this is fine, but be careful with the search function. When I tried it, it wasn't great.

### Syntax Highlight GeSHi

To colorize your code and number it, use this extension:  
[https://www.mediawiki.org/wiki/Extension:SyntaxHighlight_GeSHi](https://www.mediawiki.org/wiki/Extension:SyntaxHighlight_GeSHi)

If you want to keep the dotted lines around your sources, as well as the gray color (just like a `<syntaxhighlight lang=text></syntaxhighlight>`), here's what you need to add to your MediaWiki:Common.css address (e.g., [https://www.deimos.fr/blocnotesinfo/index.php?title=MediaWiki:Common.css](https://www.deimos.fr/blocnotesinfo/index.php?title=MediaWiki:Common.css)):

```css
div.mw-geshi {
  background-color: #f9f9f9;
  padding: 1em; 
  margin:1em 0; 
  border: 1px dashed #2f6fab;
}
```

I am very satisfied with it. I'll let you see how to do it on the link, as installation and configuration take only 2 minutes.

#### Setting default values

I have the annoying habit of always putting line=1 and start=0. But after a while, it gets annoying! That's why there are default variables, and here are the ones I've modified:

```php
    /**
     * Number at which line numbers should start at
     * @var int
     */
    var $line_numbers_start = 0;
 
    /**  
     * Flag for how line numbers are displayed
     * @var boolean
     */
    var $line_numbers = GESHI_FANCY_LINE_NUMBERS;
 
    /**
     * The size of tab stops
     * @var int
     */
    var $tab_width = 4;
 
    /**  
     * The "nth" value for fancy line highlighting
     * @var int
     */
    var $line_nth_row = 1;
```

If the changes aren't visible right away, run the update.php maintenance script:

```bash
php update.php
```

### Piwik

Piwik is an equivalent to Google Analytics, but free. It allows you to have statistics of your website with nice graphs etc... There is a [plugin](https://www.mediawiki.org/wiki/Extension:Piwik_Integration) to insert code into each page (necessary), but it's obsolete and has security flaws. So we will use another module that will allow us to insert this type of code easily.

You need to install the [PCR GUI Inserts](https://www.mediawiki.org/wiki/Extension:PCR_GUI_Inserts) module which will allow you to insert information at various places on your pages. Activate the extension first by adding these lines:

```
# PCR Extension for Piwik / Google Ads
require_once("$IP/extensions/pcr/pcr_guii.php");
```

Insert into MediaWiki's configuration file:

```php
# PCR Piwik
$wgPCRguii_Inserts['SkinAfterBottomScripts']['on'] = true;
$wgPCRguii_Inserts['SkinAfterBottomScripts']['content'] = ' 
<!-- Piwik -->
<script type="text/javascript">
var pkBaseURL = (("https:" == document.location.protocol) ? "https://www.deimos.fr/piwik/" : "http://www.deimos.fr/piwik/");
document.write(unescape("%3Cscript src=\'" + pkBaseURL + "piwik.js\' type=\'text/javascript\'%3E%3C/script%3E"));
</script><script type="text/javascript">
try {
var piwikTracker = Piwik.getTracker(pkBaseURL + "piwik.php", 2);
piwikTracker.trackPageView();
piwikTracker.enableLinkTracking();
} catch( err ) {}
</script><noscript><p><img src="http://www.deimos.fr/piwik/piwik.php?idsite=x" style="border:0" alt="" /></p></noscript>
<!-- End Piwik Tracking Code -->
';
```

### UsabilityInitiative

This extension ([UsabilityInitiative](https://www.mediawiki.org/wiki/Extension:UsabilityInitiative)) allows you to activate many new features that came with MediaWiki version 1.16 and the Vector theme.

#### MediaWiki 1.16

Download the extension and place it in the extensions folder. Then, here are the lines I added to LocalSettings.php to have additional features:

```php
# UsabilityInitiative
require_once('extensions/UsabilityInitiative/UsabilityInitiative.php');
 
# WikiEditor
require_once("$IP/extensions/UsabilityInitiative/WikiEditor/WikiEditor.php");
$wgWikiEditorModules['toolbar']['global'] = true; // Enable the WikiEditor toolbar for everyone
$wgWikiEditorModules['toolbar']['user'] = false;// Don't allow users to turn the WikiEditor toolbar on/off individually
 
# Collapse Menu
require_once("$IP/extensions/UsabilityInitiative/Vector/Vector.php");
$wgVectorModules['collapsiblenav']['user'] = true;
$wgVectorModules['collapsiblenav']['global'] = true; 
 
# Expandable search
$wgVectorModules['expandablesearch']['user'] = true;
$wgVectorModules['expandablesearch']['global'] = true;
$wgVectorUseSimpleSearch = true;
```

As you can see, I activated the WikiEditor and the collapsible menu (Sidebar).

#### MediaWiki 1.17

In version 1.17, everything has been split up. So I use [WikiEditor](https://www.mediawiki.org/wiki/Extension:WikiEditor) and [Vector](https://www.mediawiki.org/wiki/Extension:Vector) to get back the features I was interested in from version 1.16. Once installed in the extensions, activate everything:

```php
# Vector
require_once( "$IP/extensions/Vector/Vector.php" );
$wgDefaultUserOptions['vector-collapsiblenav'] = 1;
$wgVectorUseSimpleSearch = true;
 
# WikiEditor
require_once( "$IP/extensions/WikiEditor/WikiEditor.php" );
$wgDefaultUserOptions['usebetatoolbar'] = 1;
$wgDefaultUserOptions['usebetatoolbar-cgd'] = 1;
$wgDefaultUserOptions['wikieditor-preview'] = 1;
```

### Change the skins for all your users

Here's a solution to change the skin of all users who are still on Monobook to Vector (the script is in the maintenance folder):

```bash
php userOptions.php skin --old "monobook" --new "vector"
```

### MobileDetect

[MobileDetect](https://www.mediawiki.org/wiki/Extension:MobileDetect) is a simple solution to have a semblance of a version optimized for smartphones. This extension loads the "Chick" theme by default when the detected User Agent is a PDA/Smartphone. This solution is handy when you don't want to set up a complex system to make the pages as well adapted as possible to portable devices.

Put the code in the "extensions" folder, then add this to your LocalSettings.php:

```php
...
# MobileDetect
require_once("$IP/extensions/MobileDetect/MobileDetect.php");
$mobile = mobiledetect();
if ($mobile == true) $wgDefaultSkin = "chick";
```

Your MediaWiki is ready for mobile devices :-)

### DumpHTML

This extension is very practical because it allows you to dump your wiki in HTML format with images and everything else for offline use. [Download the extension](https://www.mediawiki.org/wiki/Extension:DumpHTML), then all you have to do is execute it with some optional parameters to get a good result. For my part, I made a script that dumps the wiki and compresses it. Here's the script:

```bash
#!/bin/sh
# Dump MediaWiki to html
# This permit to get an offline site
# Made by Pierre Mavro / Deimos
 
mediawiki_path="/var/www/deimos.fr/blocnotesinfo"
archive_name="bni_offline"
 
####################################################
# Check if extension is here
if [ -f "$mediawiki_path/extensions/DumpHTML/dumpHTML.php" ] ; then
    dumpHTML_php="$mediawiki_path/extensions/DumpHTML/dumpHTML.php"
elif [ -f "$mediawiki_path/maintenance/dumpHTML.php" ] ; then
    dumpHTML_php="$mediawiki_path/dumpHTML.php"
else
    echo "Can't perform a dump as dumpHTML.php is not found"
    exit 1
fi
skin=`grep '^\$wgDefaultSkin' $mediawiki_path/LocalSettings.php | awk -F\' '{ print \$2 }'`
archive_folder="$mediawiki_path/$archive_name"
 
echo "-- Start to dump the wiki --"
php $dumpHTML_php -d $archive_folder -k $skin --image-snapshot --force-copy
echo "-- Compress the wiki dump archive --"
cd $mediawiki_path
tar -czf $archive_name.tgz $archive_name
echo "-- Remove the uncompressed dump archive"
rm -Rf $archive_folder
echo "Done"
```

### CharInsert

[CharInsert](https://www.mediawiki.org/wiki/Extension:CharInsert) is an extension that allows you to add a small bar in edit mode with the wiki code that we use most often or the one that is most difficult to find. The goal is to simplify page editing as much as possible.

Download and install the extension. Once done, you need to create the following page [https://wiki.deimos.fr/blocnotesinfo/index.php?title=MediaWiki:Edittools#](https://wiki.deimos.fr/blocnotesinfo/index.php?title=MediaWiki:Edittools#)  
Insert the text you want with the `<charinsert>+</charinsert>` tags, and put in the middle the code you want (use + so that your cursor automatically goes to this point after insertion):

```
<!-- Any text entered here will be displayed under the edit boxes or file upload forms. -->
<!-- Text here will be shown below edit and upload forms. -->
<!-- Please don't translate this page with sub pages (it will render support of that menu for your language very likely unmaintainable) -->
<div id="specialchars" class="my-buttons" style="margin-top:10px; border:1px solid #aaaaaa; padding:1px; text-align:center; font-size:110%;" title="Click on the wanted special character.">
<p class="specialbasic" id="Standard">&nbsp;
<charinsert>__TOC__</charinsert> ·
<charinsert>[[:Image:+|thumb|Name]]</charinsert> ·
<charinsert>[[:Media:+]]</charinsert> ·
<charinsert>[[:Image:+]]</charinsert> ·
<charinsert>[[Category:+]]</charinsert>
<br />
<charinsert><pre></charinsert> ·
<charinsert></pre></charinsert> ·
<charinsert><nowiki><syntaxhighlight lang=text></nowiki>+</charinsert> ·
<br />
<charinsert><nowiki>{{command|+|<syntaxhighlight lang=text></nowiki></charinsert> ·
<charinsert><nowiki>{{config|+|<syntaxhighlight lang=text></nowiki></charinsert> ·
</p>
</div>
```

Then, to make it visually clean, add this to the Commoncss ([https://wiki.deimos.fr/blocnotesinfo/index.php?title=MediaWiki:Common.css](https://wiki.deimos.fr/blocnotesinfo/index.php?title=MediaWiki:Common.css))

```css
/* Extra buttons for 'edittools' */
.my-buttons {
        /* padding: 1em; */
        /* margin:5px; */
}
.my-buttons a {
        color: black;
        background-color: #d0e0f0 ;
        font-family:monospace;
        font-size: 115%;
        text-decoration: none;
        border: thin #069 outset;
        line-spacing:5pt;
}
.my-buttons a:hover {
        background-color: #99CCFF;
        border-style:outset;
}
.my-buttons a:active {
        background-color: #0645AD;
        border-style: inset;
}
```

Now when editing an article, you will have your new toolbar.

### Cite

The [Cite](https://www.mediawiki.org/wiki/Extension:Cite/Cite.php) extension allows you to have a references section. Until today, I managed it manually, but now I use this extension which allows me to make automatic links. It's practical and beautiful. However, we can make it even more beautiful by editing the CSS:

```css
/* make the Cite extension list of references look smaller and highlight clicked reference in blue */
ol.references { font-size: 90%; }
ol.references > li:target { background-color: #ddeeff; }
sup.reference:target { background-color: #ddeeff; }
```

Next, we'll update a few fields:

On the MediaWiki:Cite_references_link_one page:

```
 <li id="$1">[[#$2|?]] $3</li>
To be replaced by:
 <li id="$1">[[#$2|^]] $3</li>
```

On the MediaWiki:Cite_references_link_many page:

```
 <li id="$1">? $2 $3</li>
To be replaced by:
 <li id="$1">^ $2 $3</li>
```

On the MediaWiki:Cite_references_link_many_format page:

```
 [[#$1|<sup>$2</sup>]]
To be replaced by:
 [[#$1|<sup>$3</sup>]]
```

### Boilerplate

This is really the essential extension, it allows you to create a "template" page when creating a blank page. Needless to say, it's very useful. Installation is very easy and explained on the [https://www.mediawiki.org/wiki/Extension:Boilerplate](https://www.mediawiki.org/wiki/Extension:Boilerplate) site. Then, you just need to fill in your page that will serve as a template [https://www.deimos.fr/blocnotesinfo/index.php?title=Boilerplate](https://www.deimos.fr/blocnotesinfo/index.php?title=Boilerplate).

### MsUpload

I had been telling myself for a long time that I should find an extension that does multi-file uploads. And I'm very happy to have found one that's well-made, pretty, and as usual simple to set up. It's [MsUpload](https://www.ratin.de/msupload.html). To install it, simply unzip the desired version into the extensions folder and add this to the LocalSettings file:

```
$wgMSU_ShowAutoKat = false;    #autocategorisation
$wgMSU_CheckedAutoKat = false;  #checkbox: checked = true/false
$wgMSU_debug = false;
require_once("$IP/extensions/MsUpload/msupload.php");
```

When editing, a blue bar will appear, and you can just drag and drop.

### SphinxSearch

SphinxSearch replaces the basic search provided by MediaWiki, which is not fulltext. Explaining how to set this up here would be too long, so there is [a dedicated article about it](./sphinx_:_setup_a_full_text_indexer).
