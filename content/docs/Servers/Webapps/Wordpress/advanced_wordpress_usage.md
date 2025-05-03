---
weight: 999
url: "/Utilisation_avancée_de_Wordpress/"
title: "Advanced WordPress Usage"
description: "Tips and techniques for advanced WordPress configuration, including Nginx setup, JavaScript integration, and file handling."
categories: ["Nginx", "Linux"]
date: "2013-05-07T07:28:00+02:00"
lastmod: "2013-05-07T07:28:00+02:00"
tags: ["WordPress", "Nginx", "PHP", "JavaScript", "Web Development"]
toc: true
---

## Introduction

WordPress is great, but it also has limitations that can quickly become annoying. There are ways to overcome these limitations, which I'm going to explain here.

## Nginx Configuration

For WordPress configuration under Nginx, here's an example:

```nginx {linenos=table}
server {
    include listen_port.conf;
    listen 443 ssl;

    ssl_certificate /etc/nginx/ssl/deimos.fr/server-unified.crt;
    ssl_certificate_key /etc/nginx/ssl/deimos.fr/server.key;
    ssl_session_timeout 5m;

    server_name blog.deimos.fr;
    root /usr/share/nginx/www/deimos.fr/blog;
    index index.php;

    access_log /var/log/nginx/blog.deimos.fr_access.log;
    error_log /var/log/nginx/blog.deimos.fr_error.log;

    location / {
        try_files $uri $uri/ /index.php?$args;
    }

    location ~ \.php$ {
        fastcgi_cache mycache;
        fastcgi_cache_key $request_method$host$request_uri;
        fastcgi_cache_valid any 1h;
        include fastcgi_params;
        fastcgi_pass unix:/var/run/php5-fpm.sock;
        fastcgi_intercept_errors on;
    }

    # Drop config
    include drop.conf;

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
    }
    # END W3TC Browser Cache

    # BEGIN W3TC Minify core
    rewrite ^/wp-content/w3tc/min/w3tc_rewrite_test$ /wp-content/w3tc/min/index.php?w3tc_rewrite_test=1 last;
    rewrite ^/wp-content/w3tc/min/(.+\.(css|js))$ /wp-content/w3tc/min/index.php?file=$1 last;
    # END W3TC Minify core
}
}
```

## Using JavaScript in a Post

I want to use JavaScript here to prevent my email address from being captured by spam robots. We first need to insert this in a post:

```javascript
<script type="text/javascript" src="/scripts/updatepage.js"></script><script type="text/javascript">
<!--
```

Then our JavaScript code:

```javascript
emailE=('deimos@' + 'deimos.fr')
document.write('<A href="mailto:' + emailE + '">' + "click here" + '</a>')
```

And finally, we call our function:

```javascript
updatepage();
//--></script>
```

So the full anti-spam system looks like:

```javascript
<script type="text/javascript" src="/scripts/updatepage.js"></script><script type="text/javascript">
<!--
emailE=('deimos' + 'deimos.fr')
document.write('<A href="mailto:' + emailE + '">' + "click here" + '</a>')
updatepage();
//--></script>

<span><NOSCRIPT>
    <em>Email address protected by JavaScript.<BR>
    Please enable JavaScript to contact me.</em>
</NOSCRIPT></span>
```

## Adding Unsupported Extensions

For example, I want to enable the OGV format for my blog. I'll need to modify the file wp-includes/functions.php and add these lines:

```php
...
function wp_ext2type( $ext ) {
    $ext2type = apply_filters('ext2type', array(
        'audio' => array('aac','ac3','aif','aiff','mp1','mp2','mp3','m3a','m4a','m4b','ogg','ram','wav','wma'),
        'video' => array('asf','avi','divx','dv','mov','mpg','mpeg','mp4','mpv','ogm','qt','rm','vob','wmv', 'm4v','ogv'),
        'document' => array('doc','docx','pages','odt','rtf','pdf'),
...
    if ( !$mimes ) {
        // Accepted MIME types are set here as PCRE unless provided.
        $mimes = apply_filters( 'upload_mimes', array(
        'jpg|jpeg|jpe' => 'image/jpeg',
        'gif' => 'image/gif',
        'png' => 'image/png',
        'bmp' => 'image/bmp',
        'tif|tiff' => 'image/tiff',
        'ico' => 'image/x-icon',
        'asf|asx|wax|wmv|wmx' => 'video/asf',
        'avi' => 'video/avi',
        'divx' => 'video/divx',
        'flv' => 'video/x-flv',
        'ogv' => 'video/ogg',
        'mov|qt' => 'video/quicktime',
...
```

There are two places to modify: the wp_ext2type function and the get_allowed_mime_types function.

## Automatically Delete Comments from Trash

It's possible to automatically delete comments in the trash after a specified number of days. Edit your WordPress configuration and add this:

```php
[...]
/* Empty the trash */
define('EMPTY_TRASH_DAYS', 2 );
[...]
```

Here, I've indicated that I want the trash to empty itself every 2 days.

## FAQ

### Failed to write file to disk

I had this small problem when I tried to upload files that were too large. Here are the different solutions to fix the problem:

- Check the owner permissions and write access on the upload folder
- Look in the file `/etc/php5/cgi/php.ini` and adjust the size that suits you:

```
post_max_size = 50M
upload_max_filesize = 50M
```

- Check the size of `/tmp` on your server (for example, my vservers are by default set to 16M, which was problematic since it's the temporary location where WordPress stores its files)

## Resources
- http://codex.wordpress.org/Using_Javascript
