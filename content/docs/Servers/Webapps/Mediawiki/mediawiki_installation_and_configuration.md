---
weight: 999
url: "/MediaWiki\\:Installation_et_configuration/"
title: "MediaWiki: Installation and Configuration"
description: "A comprehensive guide for installing, configuring and managing MediaWiki including manual and Debian installation methods, extensions, web server configuration and many customization options."
categories: ["Servers", "Debian", "Security", "Nginx"]
date: "2014-11-27T07:38:00+02:00"
lastmod: "2014-11-27T07:38:00+02:00"
tags:
  [
    "MediaWiki",
    "PostgreSQL",
    "Debian Method",
    "Nginx",
    "Web Server",
    "Security",
  ]
toc: true
---

## Introduction

A wiki is a content management system for websites that allows pages to be freely and equally editable by all authorized visitors. Wikis are used to facilitate collaborative document writing with minimal constraints. The wiki was invented by Ward Cunningham in 1995, for a section of a website on computer programming that he called WikiWikiWeb. The word "wiki" comes from the Hawaiian term wiki wiki, which means "fast" or "informal". By the mid-2000s, wikis had reached a good level of maturity and are associated with Web 2.0. Created in 2001, Wikipedia has become the world's most visited wiki.

MediaWiki is a Wiki engine written in PHP and created by Magnus Manske. Initially developed for Wikipedia (which has been using it since 2002), it also serves as the foundation for other WikiMedia Foundation projects (Wiktionary, Wikisource, Wikibooks or Wikiquote). Other associations have adopted it (e.g., Wikitravel, Mozilla or Ekopedia).

## Installation

I'm describing here 2 ways to install Mediawiki. Choose the best for your needs, but I recommend the manual method.

### Manual Method

First of all, download the current version with git:

```bash
git clone https://gerrit.wikimedia.org/r/p/mediawiki/core.git
```

Then list all versions:

```bash
git tag -l | sort -V
```

And choose the wished version you want to use. Example, I want to use 1.21:

```bash
git checkout 1.21.0
```

#### Extensions install

I strongly recommend as well to use git to get all your extensions. You don't need to download all available as it takes some disk space, but you can choose them one by one like that:

```bash
cd extensions
git clone https://gerrit.wikimedia.org/r/p/mediawiki/extensions/<extension_name>.git
```

You can see the names here: https://gerrit.wikimedia.org/r/#/admin/projects/

### Debian Method

To install from the Debian repository:

```bash
aptitude install mediawiki
```

Alternatively, you can install it manually from the main site.

#### PostgreSQL

If after installation you don't want to use MySQL but Postgres as a database, you'll certainly encounter some issues. [Refer to this link]({{< ref "docs/Servers/Databases/PostgreSQL/postgresql_installation_and_configuration.md">}}) and look at the FAQ. If it still doesn't work, if you encounter other types of errors, then try this:

If you are using Postgres, you will need to either have a database and user created for you, or simply supply the name of a Postgres user with "superuser" privileges to the configuration form. Often, this is the database user named **postgres**.

The database that MediaWiki will use needs to have both plpgsql and tsearch2 installed. The installer script will try and install plpgsql, but you may need to install tsearch2 yourself. (tsearch2 is used for searching the text of your wiki). Here's one way to do most of the setup. This is for a Unix-like system, and assumes that you have already installed the plpgsql and tsearch2 modules. In this example, we'll create a database named **wikidb**, owned by a user named **wikiuser**. From the command-line, as the postgres user, perform the following steps.

```
createuser -S -D -R -P -E wikiuser (then enter the password)
createdb -O wikiuser wikidb
createlang plpgsql wikidb
```

Adding tsearch2 to the database is not a simple step, but hopefully it will already be done for you by whatever packaging process installed the tsearch2 module. In any case, the installer will let you know right away if it cannot find tsearch2.

The above steps are not all necessary, as the installer will try and do some of them for you if supplied with a superuser name and password.

For installing tsearch2 to the wikidb database under Windows, do the following steps:

- Find tsearch2.sql (probably under .\PostgreSQL\8.x\share\contrib) and copy it to the postgresql\8.x\bin directory;
- From a command prompt at the postgresql\8.x\bin directory, type:

```
psql wikidb < tsearch2.sql -U wikiuser
```

- It will prompt you for the password for wikiuser;

That's it!

Point (2) seems only to work on windows, because on debian linux 4.0 (etch) only user postgres is allowed to use language c. So there it must be called by:

```
su - postgres -c psql wikidb < tsearch2.sql
```

afterwards you must grant select rights to wikiuser to the tsearch tables and insert the correct locale.

```
su - postgres
psql -d wikidb -c "grant select on pg_ts_cfg to wikiuser;"
psql -d wikidb -c "grant select on pg_ts_cfgmap to wikiuser;"
psql -d wikidb -c "grant select on pg_ts_dict to wikiuser;"
psql -d wikidb -c "grant select on pg_ts_parser to wikiuser;"
psql -d wikidb -c "update pg_ts_cfg set locale = current_setting('lc_collate') where ts_name = 'default' and prs_name='default';"
```

If you receive an error similar to "ERROR: relation "pg_ts_cfg" does not exist" when executing the above statements, try installing tsearch2 to the wikidb database again, but instead use these two separate steps (and then try the grant statements again):

```
su - postgres
psql wikidb -f tsearch2.sql
```

## Upgrade

To upgrade Mediawiki with git, go into your Mediawiki folder instance and check the version you're running on:

```bash
> cd /var/www/mediawiki
> git describe --tags
1.21.0
```

Then get the latest version of Mediawiki:

```bash
git fetch
```

Change your version to the desired version:

```bash
git checkout 1.21.1
```

and upgrade your instance with the database:

```bash
php maintenance/update.php
```

That's all, your MediaWiki core is now up to date. You now need look at extensions.

### Extensions upgrade

To upgrade you extensions, you've used git, that's why it will be easy. Anyway, some of them may not be managed with git, but the old repository version (SVN). I've wrote a little script to handle that:

```bash
#!/bin/sh
extensions=`pwd`
changes=0
for i in * ; do
    if [ -d $extensions/$i ] ; then
        echo ""
        echo "[+] $i"
        cd $extensions/$i
        # Git
        if [ -d .git ] ; then
            git pull
            changes=1
        # SVN
        elif [ -d .svn ] ; then
            svn up
            changes=1
        fi
     fi
done

# Reset rights
if [ $changes -eq 1 ] ; then
    chown -Rf www-data. .
fi
```

Simply copy it in your extension directory and launch it from there:

```bash
> ./update_extensions.sh

[+] CharInsert
Already up-to-date.

[+] Cite
remote: Counting objects: 41, done
remote: Finding sources: 100% (6/6)
remote: Total 6 (delta 4), reused 6 (delta 4)
Unpacking objects: 100% (6/6), done.
From https://gerrit.wikimedia.org/r/p/mediawiki/extensions/Cite
   22f4d9e..1e542ef  master     -> origin/master
Updating 22f4d9e..1e542ef
Fast-forward
 Cite.i18n.php |   35 +++++++++++++++++++++++++++++++++++
 1 file changed, 35 insertions(+)

[+] Gadgets
Already up-to-date.

[+] LdapAuthentication
Already up-to-date.

[+] MsUpload
Already up-to-date.

[+] MultiBoilerplate
At revision 115794.

[+] ParserFunctions
Already up-to-date.

[+] SyntaxHighlight_GeSHi
Already up-to-date.

[+] Vector
Already up-to-date.

[+] WikiEditor
Already up-to-date.
```

Everything is up to date now :-)

## Configuration

### Web Server

#### Nginx

For setting up Mediawiki with Nginx and short URLs, here is the configuration to adopt. I also added SSL and forced redirects from the login page to SSL:

```bash
server {
    include listen_port.conf;
    listen 443 ssl;

    ssl_certificate /etc/nginx/ssl/deimos.fr/server-unified.crt;
    ssl_certificate_key /etc/nginx/ssl/deimos.fr/server.key;
    ssl_session_timeout 5m;

    server_name wiki.deimos.fr wiki.m.deimos.fr;
    root /usr/share/nginx/www/deimos.fr/blocnotesinfo;

    client_max_body_size 5m;
    client_body_timeout 60;

    access_log /var/log/nginx/wiki.deimos.fr_access.log;
    error_log /var/log/nginx/wiki.deimos.fr_error.log;

    location / {
        rewrite ^/$ $scheme://$host/index.php permanent;
        # Short URL redirect
        try_files $uri $uri/ @rewrite;
    }

    location @rewrite {
        if (!-f $request_filename){
            rewrite ^/(.*)$ /index.php?title=$1&$args;
        }
    }

    # Force SSL Login
    set $ssl_requested 0;
    if ($arg_title ~ Sp%C3%A9cial:Connexion) {
        set $ssl_requested 1;
    }
    if ($scheme = https) {
        set $ssl_requested 0;
    }
    if ($ssl_requested = 1) {
        return 301 https://$host$request_uri;
    }

    # Drop config
    include drop.conf;

    # Deny direct access to specific folders
    location ^~ /(maintenance|images)/ {
        return 403;
    }

    location ~ \.php$ {
        fastcgi_cache mycache;
        fastcgi_cache_key $request_method$host$request_uri;
        fastcgi_cache_valid any 1h;
        include fastcgi_params;
        fastcgi_pass unix:/var/run/php5-fpm.sock;
    }

    location = /_.gif {
        expires max;
        empty_gif;
    }

    location ^~ /cache/ {
        deny all;
    }

    location /dumps {
        root /usr/share/nginx/www/deimos.fr/blocnotesinfo/local;
        autoindex on;
    }

    # BEGIN W3TC Browser Cache
    gzip on;
    gzip_types text/css application/x-javascript text/x-component text/richtext image/svg+xml text/plain text/xsd text/xsl text/xml image/x-icon;
    location ~ \.(css|js|htc)$ {
        expires 31536000s;
        add_header Pragma "public";
        add_header Cache-Control "public, must-revalidate, proxy-revalidate";
        add_header X-Powered-By "W3 Total Cache/0.9.2.4";
    }

    location ~ \.(html|htm|rtf|rtx|svg|svgz|txt|xsd|xsl|xml)$ {
        expires 3600s;
        add_header Pragma "public";
        add_header Cache-Control "public, must-revalidate, proxy-revalidate";
        add_header X-Powered-By "W3 Total Cache/0.9.2.4";
    }

    location ~ \.(asf|asx|wax|wmv|wmx|avi|bmp|class|divx|doc|docx|eot|exe|gif|gz|gzip|ico|jpg|jpeg|jpe|mdb|mid|midi|mov|qt|mp3|m4a|mp4|m4v|mpeg|mpg|mpe|mpp|otf|odb|odc|odf|odg|odp|ods|odt|ogg|pdf|png|pot|pps|ppt|pptx|ra|ram|svg|svgz|swf|tar|tif|tiff|ttf|ttc|wav|wma|wri|xla|xls|xlsx|xlt|xlw|zip)$ {
        expires 31536000s;
        add_header Pragma "public";
        add_header Cache-Control "public, must-revalidate, proxy-revalidate";
        add_header X-Powered-By "W3 Total Cache/0.9.2.4";
        try_files $uri $uri/ @rewrite;
    }
    # END W3TC Browser Cache
}
```

### Left Menu (Sidebar)

- To edit your menu on the left:

Access the modifications via this link:

```
https://wiki.deimos.fr/index.php?title=MediaWiki:Sidebar
```

Then, you can edit it by putting **"link|name"**

Here's an example:

```
* navigation
** mainpage|mainpage
** portal-url|portal
** currentevents-url|currentevents
** recentchanges-url|recentchanges
** randompage-url|randompage
** helppage|help
** sitesupport-url|sitesupport

* quick navigation
** linux|Linux
** unix|Unix
** Mac OS X|Mac OS X
** windows|Windows
** serveurs|Servers
** réseaux|Networks
** divers|Miscellaneous
```

- To create an interwiki link:

[[Mylink#Interwiki|The name I want to give]]

### One instance, multiple wikis (multi-tenant / wiki-family)

#### Method 1

If for example you want to do multi-languages or simply have multiple databases and not make n updates for each of your bases, then opt for this solution.

The goal is to modify the code of the LocalSettings.php file in order to be able to detect in the http header (URL), a succession of characters that can refer to one base or another.

Proceed like this:

- Run a wiki setup. Once the LocalSettings.php file is created, rename it to **fr.php** for example.
- Run another wiki setup. Once the LocalSettings.php file is created, rename it to **en.php** for example.
- Create the LocalSettings.php file, then insert and adapt the following:

```php
<?php

$callingurl = strtolower($_SERVER['SERVER_NAME']); // identify the asking url

if ( $callingurl == "fr.deimos.fr" ) {
        require_once( 'fr.php' );
}

if ( $callingurl == "en.deimos.fr" ) {
        require_once( 'en.php' );
}

?>
```

In this configuration, if the url corresponds to http://fr.deimos.fr/mediawiki/index.php, for example, then the fr.php file will be taken into account.

#### Method 2

Here is another method that allows you to choose a different configuration file depending on the first parameter following the url. Example:

- http://www.deimos.fr/**wiki1**
- http://www.deimos.fr/**wiki2**
- http://www.deimos.fr/**wiki3**
- ...

Modify the LocalSettings.php file:

```php
<?php
list($null, $toplevel, $dontcare) = explode('/', strtolower($_SERVER['REQUEST_URI']));

if (is_file("LocalSettings-$toplevel.php"))
{
    require_once("LocalSettings-$toplevel.php");
}
else
{
    print <<<EOF
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.1//EN" "http://www.w3.org/TR/xhtml11/DTD/xhtml11-flat.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="en">
<head>
<title>wiki not found </title>
</head>
<body>
<p>The wiki you were looking for was not found.</p>
</body>
EOF;
}
?>
```

All you have to do is create your configurations like this:

- LocalSettings-wiki1.php
- LocalSettings-wiki2.php
- LocalSettings-wiki3.php

### URL Change

You may have wanted to start a wiki with a certain name and then want to change the folder name later. This means that the url has changed. Let me explain:

I installed media wiki and it points to:

```
http://www.deimos.fr/mediawiki
```

I've already created articles. I would now like to rename this to:

```
http://www.deimos.fr/blocnotesinfo
```

The problem is that the old articles have kept the first url in memory. To update all this, do:

```
cd /usr/share/mediawiki/maintenance
php clear_interwiki_cache.php
php alltrans.php
php refreshLinks.php
```

Try now and it should be much better :-)

### Refresh Your Cache

To refresh your cache, just add **?action=purge** at the end of the wiki url. Ex:

```
https://wiki.deimos.fr/index.php?action=purge
```

### RSS Feed

To have a link with these RSS feeds, it's quite simple, judge for yourself:

- Recent changes:

```
https://wiki.deimos.fr/index.php?title=Special:Recentchanges&feed=rss
https://wiki.deimos.fr/index.php?title=Special:Recentchanges&feed=atom
```

- New Pages:

```
https://wiki.deimos.fr/index.php?title=Special:Newpages&feed=rss
https://wiki.deimos.fr/index.php?title=Special:Newpages&feed=atom
```

- Page history (MW1.7+):

```
https://wiki.deimos.fr/index.php?title=MediaWiki:Installation_et_configuration&action=history&feed=atom
```

### Allow Searches from 2 Characters

By default, searching for terms in mediawiki doesn't take into account words of 3 letters or less. However, in computing, words like XEN, VNC, etc. are very common.

Mediawiki uses mysql's full-text search functionality for its searches. This option allows indexing certain text fields to perform advanced searches that respect certain language constraints. So that the indexes don't take up too much space, the default configuration ignores words with fewer than 4 letters. This is the ft_min_word_len parameter

- Modify the /etc/mysql/my.cnf file to modify the ft_min_word_len parameter

```
[mysqld]
ft_min_word_len=3
```

- restart mysql
- connect to the database and launch an index repair:

```
mysql -u root -p wikidb
mysql> REPAIR TABLE searchindex;
```

There you go, now you can search for 2-letter words in your wiki.

### Remove Title on Main Page

After a long search, I finally found a small javascript that removes the "Home" title from the start page :-).

#### Version < 1.12.0

Insert these lines into your **MediaWiki:Common.js** file (https://wiki.deimos.fr/MediaWiki:Common.js) and insert these lines:

```javascript
/* hide heading on [[Main_Page]] */
   var mpTitle = "Home";
   var isMainPage = (document.title.substr(0, document.title.lastIndexOf(" - ")) == mpTitle);
   var isDiff = (document.location.search && (document.location.search.indexOf("diff=") != -1 || document.location.search.indexOf("oldid=") != -1));

   if (isMainPage && !isDiff) {
      document.write('<style type="text/css">/*<![CDATA[*/ #lastmod, #siteSub, #contentSub, h1.firstHeading { display: none !important; } /*]]>*/</style>');
   }
/*
```

#### Version >= 1.12.0

Edit your file https://wiki.deimos.fr/MediaWiki:Common.css and insert this:

```css
/** CSS placed here will be applied to all appearances. */

.page-Accueil * .firstHeading,
.page-Accueil * h3#siteSub,
.page-Accueil * #contentSub {
  display: none;
}
```

### Limit User Creation

This will prohibit user creation, and editing from anonymous users. To be put in your LocalSettings.php file:

```
# This snippet prevents new registrations from anonymous users
# (Sysops can still create user accounts)
$wgGroupPermissions['*']['createaccount'] = false;

# This snippet prevents editing from anonymous users
$wgGroupPermissions['*']['edit'] = false;
require_once( "includes/DefaultSettings.php" );
```

### File Import

To authorize certain extensions, edit these few lines in your LocalSettings.php file:

```
## To enable image uploads, make sure the 'images' directory
## is writable, then set this to true:
$wgEnableUploads       = true;
$wgCheckFileExtensions = false;
$wgStrictFileExtensions = true;
$wgFileExtensions = array( 'png', 'gif', 'jpg', 'jpeg', 'ogg', 'zip', 'pdf', 'gz', 'tgz', 'sxw', 'ipcc');
$wgVerifyMimeType = false;
$wgUseImageResize = false;
$wgUseImageMagick = true;
$wgImageMagickConvertCommand = "/usr/bin/convert";
```

### Improve Search Ease

Add these lines to your LocalSettings file for advanced search options. To have Google-like search suggestions:

```
# Drop-down AJAX search suggestions
$wgEnableMWSuggest  = true;
```

And for more relevant searches:

```
# More relevant search snippets
$wgAdvancedSearchHighlighting = true;
```

### Add Text to the Login Part

I would like, for example, to add text to offer an SSL connection. Edit the Loginend Page:

```
https://wiki.deimos.fr/blocnotesinfo/index.php?title=MediaWiki:Loginend
```

And we can put text like this, for example:

```html
<div style="clear:both; font-size:.85em; line-height:1.4em; margin-left:2em">
  <u>Note</u>: Information sent during this connection will not be encrypted.<br />
  For more security, it is possible to connect via the
  [https://{{SERVERNAME}}/blocnotesinfo/index.php?title=Sp%C3%A9cial:Connexion&returnto=Accueil
  secure server]. ---

  <u>Note</u>: The information sent during this connection will not be
  encrypted.<br />
  For more security, it is possible to connect through the
  [https://{{SERVERNAME}}/blocnotesinfo/index.php?title=Sp%C3%A9cial:Connexion&returnto=Accueil
  secure server].
</div>
```

### Stay on SSL

The $wgServer value allows you to define the address of the wiki, which will be regularly used to access certain pages. If like me, you have forced authentication in SSL, you will notice that you return to non-encrypted just after logging in. And this is not necessarily the desired case. If you want to stay on SSL when you start using it, edit your configuration on the $wgServer part and replace it with this:

```php
if ( !empty( $_SERVER['HTTPS'] ) ) {
    $wgServer           = "https://wiki.deimos.fr";
}else{
    $wgServer           = "http://wiki.deimos.fr";
}
```

### Modify the Header for All Articles

If you need to modify all the pages of your wiki so that the header is identical everywhere, you should use hooks. For this hook (ArticlePageDataBefore), just add this at the end of the LocalSettings.php file:

```php
# Article Header
function MyArticleHeader( $article, $fields )
{
    global $wgOut;
    $wgOut->addWikiText('{{Template:ArticleHeader}}');
    return true;
}
$wgHooks['ArticlePageDataBefore'][] = 'MyArticleHeader';
```

Here I have an ArticleHeader template that will be displayed in all articles.

### Hide Version Information

#### For Non-Connected Users

You don't necessarily want everyone to be able to get information about the platform on which MediaWiki runs, or information about the MediaWiki version itself. Here's how to remove the Version page for non-connected people:

```php
# Hide Version Special Page
function DisableSpecialPages(&$SpecialPageslist)
{
    global $wgUser;
    if ( $wgUser->isAnon() )
    {
        unset( $SpecialPageslist['Version'] );
        return true;
    }
    return true;
}
$wgHooks['SpecialPage_initList'][]='DisableSpecialPages';
```

#### For Non-Admin Users

If you want only administrators to have access to the Version page:

```php
# Hide Version Special Page
function DisableSpecialPages(&$SpecialPageslist)
{
    global $wgUser;
    if ( !$wgUser->isAllowed('protect') )
    {
        unset( $SpecialPageslist['Version'] );
        return true;
    }
    return true;
}
$wgHooks['SpecialPage_initList'][]='DisableSpecialPages';
```

### Hide Sidebar During Edits

It is possible to hide the sidebar when a connected user makes modifications. I found [this solution](https://en.wikipedia.org/wiki/User:PleaseStand/Hide_Vector_sidebar) using JavaScript that I liked. The problem is that it's for each user, you have to edit a file that is specific to them. So I modified this JavaScript somewhat to make it work for all users once installed. I also added a tab for file uploads. To set it up, it's quite simple and I remind you that it's only compatible with the Vector theme. Edit this page https://wiki.deimos.fr/MediaWiki:Vector.js and insert this:

```javascript
/* hide-vector-sidebar.js: Adds a button to toggle visibility of the Vector sidebar.
   Written by PleaseStand. Public domain; all copyright claims waived as described
   in http://en.wikipedia.org/wiki/Template:PD-self
   Modified by Deimosfr <xxx@mycompany.com>
*/

/*global document, window, addOnloadHook, addPortletLink, skin*/

var sidebarSwitch;

function sidebarHide() {
  document.getElementById("mw-panel").style.visibility = "hidden";
  document.getElementById("mw-head-base").style.marginLeft = "0";
  document.getElementById("content").style.marginLeft = "0";
  document.getElementById("left-navigation").style.left = "0";
  document.getElementById("footer").style.marginLeft = "0";
  if (typeof sidebarSwitch == "object") {
    sidebarSwitch.parentNode.removeChild(sidebarSwitch);
  }
  sidebarSwitch = addPortletLink(
    "p-cactions",
    "javascript:sidebarShow()",
    "Show sidebar",
    "ca-sidebar",
    "Show the navigation links",
    "a"
  );
}

function sidebarShow() {
  document.getElementById("mw-panel").style.visibility = "";
  document.getElementById("mw-head-base").style.marginLeft = "";
  document.getElementById("content").style.marginLeft = "";
  document.getElementById("left-navigation").style.left = "";
  document.getElementById("footer").style.marginLeft = "";
  if (typeof sidebarSwitch == "object") {
    sidebarSwitch.parentNode.removeChild(sidebarSwitch);
  }
  sidebarSwitch = addPortletLink(
    "p-cactions",
    "javascript:sidebarHide()",
    "Hide sidebar",
    "ca-sidebar",
    "Hide the navigation links",
    "a"
  );
}

function createTab() {
  addPortletLink(
    "p-cactions",
    wgArticlePath.replace("$1", "Special:Upload"),
    "Import a file"
  );
}

// Only activate on Vector skin
if (skin == "vector") {
  addOnloadHook(function () {
    if (document.getElementById("editform")) {
      // Change this if you want to show the sidebar by default
      sidebarHide();
      // Add custom tab
      addOnloadHook(createTab);
    }
  });
}
```

And there you go :-). Edit a page to see the result! You can simply show and hide the sidebar from the tabs.

### Shorter URL

It is possible to make shorter URLs to go from http://www.domainname.com/wiki/index.php?title=Accueil to http://wiki.domainname.com/Accueil. That's better, isn't it? You'll see it's not very complicated. First, in your Mediawiki configuration file, make sure you have these lines:

```php
[...]
$wgMetaNamespace = "Name of my wiki";
$wgScriptPath = "";
$wgArticlePath = "/$1";
$wgScriptExtension = ".php";
[...]
```

And then, for the Apache part, add this to your VirtualHost for your wiki:

```apache
[...]
RewriteEngine On
RewriteCond %{DOCUMENT_ROOT}%{REQUEST_URI} !-f
RewriteCond %{DOCUMENT_ROOT}%{REQUEST_URI} !-d
RewriteRule ^(.*)$ %{DOCUMENT_ROOT}/index.php [L,QSA]
[...]
```

Restart your Apache and you're good to go :-)

### Remove Discussion Pages

It is possible to remove discussion pages (Talk pages) by adding this Hook to your configuration:

```php
# Remove talkpage tab
function RemoveVectorTabs( SkinTemplate &$sktemplate, array &$links ) {
        global $wgUser, $wgHVTFUUviewsToRemove;
        if ( isset( $links['namespaces']['talk'] ))
                unset( $links['namespaces']['talk'] );
        return true;
}
$wgHooks['SkinTemplateNavigation'][] = 'RemoveVectorTabs';
```

### Enable Debug Mode

It is possible to struggle somewhat with setting up certain extensions. To address this, you need to activate the superb debug mode by adding these lines to your LocalSettings.php:

```php
# Debug mode
error_reporting( -1 );
ini_set( 'display_errors', 1 );
$wbDebugLogFile = '/tmp/debug.log';
$wgDebugComments = true;
$wgDebugToolbar = true;
```

### Open a Link in a New Window

If you want to permanently open external links in a new window/tab, put this in your configuration:

```php
[...]
$wgExternalLinkTarget = '_blank';
[...]
```

Otherwise, if you only want to open a particular link in a new window/tab, you'll need to add a JavaScript to the MediaWiki:Common.js page:

```javascript
addOnloadHook(function () {
  var pops = function (elems) {
    for (var i = 0; i < elems.length; i++) {
      if (!(" " + elems[i].className + " ").match(/ pops /)) continue;
      var anchs = elems[i].getElementsByTagName("a");
      for (var j = 0; j < anchs.length; j++) anchs[j].target = "_blank";
    }
  };
  var bc = document.getElementById("bodyContent");
  var tags = ["span", "div", "table", "td", "th"];
  for (var i = 0; i < tags.length; i++) pops(bc.getElementsByTagName(tags[i]));
});
```

Then add your new link in a page like this:

```html
<span class="pops">http://www.deimos.fr</span>
```

## Use Information Windows

You've certainly seen on Wikipedia or other mediawiki windows at the top left of articles presenting certain information. This is called an infobox. While browsing the web, I could see infoboxes more complicated than others requiring plugins when in fact it is possible to create very beautiful ones without plugins.

For those interested, here's mine: http://www.deimos.fr/blocnotesinfo/index.php?title=Template:Infobox

And here's how to use it. Insert this into a page:

```
{{Infobox
| infobox_image       = [[Image:logo.png|220px|Name]]
| infobox_softversion = x.x
| infobox_os          = OS
| infobox_website     = [http://www.deimos.fr/ Deimos Website]
| infobox_others      =
}}
```

## Monobook Display

To change the Monobook CSS for all users, add this to your url:

```
https://wiki.deimos.fr/blocnotesinfo/index.php/MediaWiki:Monobook.css
```

Then, for CSS, here's what I put:

```css
/* CSS placed here will be applied to all appearances. */

/* Hide Acceuil page on the main page */
.page-Accueil * .firstHeading,
.page-Accueil * h3#siteSub,
.page-Accueil * #contentSub {
  display: none;
}

/* Set default colors for geshi addon */
div.mw-geshi {
  background-color: #f9f9f9;
  padding: 1em;
  margin: 1em 0;
  border: 1px dashed #2f6fab;
}

/* Change default background image */
body {
  background: white url("images/a/a6/Headbg.jpg") 0px 0px no-repeat;
}

/* Rounding corners for Firefox and Mozilla browser. Warning: This could make your website slower */

.pBody {
  padding: 0.3em 0.3em;
  -moz-border-radius-topright: 0.5em;
  -moz-border-radius-bottomright: 0.5em;
}

.portlet h5 {
  -moz-border-radius-topright: 0.5em;
}

#p-cactions ul li,
#p-cactions ul li a {
  -moz-border-radius-topright: 0.5em;
  -moz-border-radius-topleft: 0.5em;
}

#content {
  -moz-border-radius-topleft: 0.5em;
  -moz-border-radius-bottomleft: 0.5em;
}

/* Changing background opacity */
.ns-0 * #content,
.ns-0 * #p-cactions li,
.ns-0 * #p-cactions li a {
  filter: alpha(opacity=90);
}
```

## Disable Skins

To disable skins, it's very simple, just add this to your LocalSettings:

```php
# To remove various skins from the User Preferences choices
$wgSkipSkins = array("standard", "cologneblue", "modern", "monobook", "myskin", "nostalgia", "simple");
```

## FAQ

Here's the site info:
View the [User Guide](https://meta.wikipedia.org/wiki/Aide:Contenu) for more information on using this software.

Please see [documentation on customizing the interface](https://meta.wikipedia.org/wiki/MediaWiki_i18n) and the [User's Guide](https://meta.wikipedia.org/wiki/MediaWiki_User%27s_Guide) for usage and configuration help.

### Know Your Version

To know your version, just add **Special:version** after "title=" in your URL:

```
https://wiki.deimos.fr/index.php?title=Special:version
```

### Various Problems After an Update

Are you experiencing issues following a MediaWiki update? Try this, it fixes a lot of problems:

```bash
php maintenance/update.php
```

## Resources
- [Your Knowledge Base with Mediawiki](/pdf/votre_base_de_connaissance_avec_mediawiki.pdf)
