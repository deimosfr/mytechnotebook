---
weight: 999
url: "/Theme_pour_le_Directory_Index_Listing_d\\'Apache/"
title: "Theme for Apache Directory Index Listing"
description: "How to create a custom theme for Apache's Directory Index Listing to improve its appearance."
categories: ["Linux", "Servers", "Debian"]
date: "2012-04-09T13:15:00+02:00"
lastmod: "2012-04-09T13:15:00+02:00"
tags: ["Apache", "Web Server", "Customization", "CSS", "HTML", "Directory Listing"]
toc: true
---

![File.Index](/images/apachefileindex_logo.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 1.0 |
| **Operating System** | Debian 6 |
| **Website** | [Website](https://blog.is-a-geek.org/apache-file-listing-mit-eigenem-theme-verschonern) |
| **Last Update** | 09/04/2012 |
{{< /table >}}

## Introduction

The default Apache directory index listing is... (we can say it)... very ugly. Here is a theme I found [on this site](https://blog.is-a-geek.org/apache-file-listing-mit-eigenem-theme-verschonern) and which [I have slightly modified](/others/directory_index_theme.tgz). Here's a preview:

![Apachefileindex](/images/apachefileindex.avif)

## Installation

It's easy to install. Let's say we want to place it in `/var/www/include` to make things easier:

```bash
wget /var/www/ http://www.deimos.fr/blocnotesinfo/images/f/f8/Directory_index_theme.tgz
tar -xzf Directory_index_theme.tgz
rm -f Directory_index_theme.tgz
chown -Rf www-data. include
```

## Configuration

Now we just need to configure Apache. Again, we'll keep things simple and modify Apache's base file:

```apache {linenos=table,hl_lines=["13-61"]}
<VirtualHost *:80>
	ServerAdmin webmaster@localhost

	DocumentRoot /var/www
	<Directory />
		Options FollowSymLinks
		AllowOverride None
	</Directory>
	<Directory /var/www/>
		Options Indexes FollowSymLinks MultiViews
		AllowOverride None
		Order allow,deny
		allow from all
	</Directory>

    Alias /icons/ /var/www/include/icons/
    Alias /include/ /var/www/include/
    <Location />
        <IfModule mod_autoindex.c>
          Options Indexes FollowSymLinks
          IndexOptions +FancyIndexing
          IndexOptions +VersionSort
          IndexOptions +HTMLTable
          IndexOptions +FoldersFirst
          IndexOptions +IconsAreLinks
          IndexOptions +IgnoreCase
          IndexOptions +SuppressDescription
          IndexOptions +SuppressHTMLPreamble
          IndexOptions +XHTML
          IndexOptions +IconWidth=16
          IndexOptions +IconHeight=16
          IndexOptions +NameWidth=*
          IndexOptions +DescriptionWidth=200
          IndexOrderDefault Descending Date
          HeaderName /include/header.html
          ReadmeName /include/footer.html
          AddIcon /icons/type_application.png .exe .app .EXE .APP
          AddIcon /icons/type_binary.png .bin .hqx .uu .BIN .HQX .UU
          AddIcon /icons/type_box.png .tar .tgz .tbz .tbz2 bundle .rar .TAR .TGZ .TBZ .TBZ2
          AddIcon /icons/type_rar.png .rar .RAR
          AddIcon /icons/type_html.png .htm .html .HTM .HTML
          AddIcon /icons/type_code.png .htx .htmls .dhtml .phtml .shtml .inc .ssi .c .cc .css .h .rb .js .rb .pl .py .sh .shar .csh .ksh .tcl .as
          AddIcon /icons/type_database.png .db .sqlite
          AddIcon /icons/type_disc.png .iso .image
          AddIcon /icons/type_document.png .ttf
          AddIcon /icons/type_excel.png .xlsx .xls .xlm .xlt .xla .xlb .xld .xlk .xll .xlv .xlw
          AddIcon /icons/type_flash.png .flv
          AddIcon /icons/type_illustrator.png .ai .eps .epsf .epsi
          AddIcon /icons/type_pdf.png .pdf
          AddIcon /icons/type_php.png .php .phps .php5 .php3 .php4 .phtm
          AddIcon /icons/type_photoshop.png .psd
          AddIcon /icons/type_monitor.png .ps
          AddIcon /icons/type_powerpoint.png .ppt .pptx .ppz .pot .pwz .ppa .pps .pow
          AddIcon /icons/type_swf.png .swf
          AddIcon /icons/type_text.png .tex .dvi
          AddIcon /icons/type_vcf.png .vcf .vcard
          AddIcon /icons/type_word.png .doc .docx
          AddIcon /icons/type_zip.png .Z .z .tgz .gz .zip
          AddIcon /icons/type_globe.png .wrl .wrl.gz .vrm .vrml .iv
          AddIcon /icons/type_android.gif .apk .APK

          AddIconByType (TXT,/icons/type_text.png) text/*
          AddIconByType (IMG,/icons/type_image.png) image/*
          AddIconByType (SND,/icons/type_audio.png) audio/*
          AddIconByType (VID,/icons/type_video.png) video/*
          AddIconByEncoding (CMP,/icons/type_box.png) x-compress x-gzip

          AddIcon /icons/back.png ..
          AddIcon /icons/information.png README INSTALL
          AddIcon /icons/type_folder.png ^^DIRECTORY^^
          AddIcon /icons/blank.png ^^BLANKICON^^

          DefaultIcon /icons/type_document.png
        </ifModule>
    </Location>

	ScriptAlias /cgi-bin/ /usr/lib/cgi-bin/
	<Directory "/usr/lib/cgi-bin">
		AllowOverride None
		Options +ExecCGI -MultiViews +SymLinksIfOwnerMatch
		Order allow,deny
		Allow from all
	</Directory>

	ErrorLog ${APACHE_LOG_DIR}/error.log

	# Possible values include: debug, info, notice, warn, error, crit,
	# alert, emerg.
	LogLevel warn

	CustomLog ${APACHE_LOG_DIR}/access.log combined

    Alias /doc/ "/usr/share/doc/"
    <Directory "/usr/share/doc/">
        Options Indexes MultiViews FollowSymLinks
        AllowOverride None
        Order deny,allow
        Deny from all
        Allow from 127.0.0.0/255.0.0.0 ::1/128
    </Directory>
</VirtualHost>
```

Adapt this configuration to your needs and then reload Apache:

```bash
/etc/init.d/apache2 reload
```

## Resources
- http://blog.is-a-geek.org/apache-file-listing-mit-eigenem-theme-verschonern
