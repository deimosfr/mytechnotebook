---
weight: 999
url: "/Gitweb_\\:_Installation_et_configuration_d'une_interface_web_pour_git/"
title: "Gitweb: Installation and configuration of a web interface for git"
description: "Learn how to install and configure Gitweb, the web interface that allows you to view all commits and comments for git repositories."
categories: ["Server", "Git", "Web"]
date: "2014-02-17T21:59:00+02:00"
lastmod: "2014-02-17T21:59:00+02:00"
tags: ["git", "gitweb", "web", "apache", "nginx", "lighttpd"]
toc: true
---

## Introduction

Gitweb is a web interface that allows you to view all commits, comments, etc. for [git]({{< ref "docs/Servers/Versionning/Git">}}).

## Installation

On Debian, it's easy:

```bash
apt-get install gitweb
```

## Configuration

Personally, I use lighttpd, but I'll try to provide correct configuration for Apache as well.

### Web Server

#### Lighttpd

If you're using Lighttpd, make sure the following server modules are loaded:

- mod_cgi
- mod_redirect

Then, create a file for your lighttpd configuration:

```perl
# Gitweb

url.redirect += (
    "^/gitweb$" => "http://www.deimos.fr/gitweb/",
)
alias.url += (
    "/gitweb/" => "/usr/lib/cgi-bin/gitweb.cgi",
    "/gitweb.css" => "/usr/share/gitweb/gitweb.css",
    "/git-favicon.png" => "/usr/share/gitweb/git-favicon.png",
    "/git-logo.avif" => "/usr/share/gitweb/git-logo.avif",
)

$HTTP["url"] =~ "^/gitweb/" {
    setenv.add-environment = (
        "GITWEB_CONFIG" => "/etc/gitweb.conf",
    )
    cgi.assign = ( "" => "" )
}

$HTTP["host"] =~ "^(git|gitweb)\.(.*)" { url.redirect = ( "^/(.*)" => "http://www.%2/gitweb/" )}
```

Then enable it:

```bash
cd /etc/lighttpd/conf-enabled/
ln -s /etc/lighttpd/conf-available/50-gitweb.conf .
```

Restart or reload your lighttpd server afterward.

#### Apache

And here's the configuration for Apache, for those who use it:

```apache
<VirtualHost *:80>
    ServerName git.example.org
    DocumentRoot /pub/git
    SetEnv  GITWEB_CONFIG   /etc/gitweb.conf
    RewriteEngine on
    # make the front page an internal rewrite to the gitweb script
    RewriteRule ^/$  /cgi-bin/gitweb.cgi
    # make access for "dumb clients" work
    RewriteRule ^/(.*\.git/(?!/?(HEAD|info|objects|refs)).*)?$ /cgi-bin/gitweb.cgi%{REQUEST_URI}  [L,PT]
</VirtualHost>
```

Or, in the Debian Squeeze version, you can find this:

```apache
Alias /gitweb /usr/share/gitweb

<Directory /usr/share/gitweb>
  Options FollowSymLinks +ExecCGI
  AddHandler cgi-script .cgi
</Directory>
```

Restart or reload your Apache server afterward.

#### Nginx

With Nginx, you'll need to authorize certain extensions in php-fpm:

```bash
# In /etc/php5/fpm/pool.d/www.conf
[...]
security.limit_extensions = .php .php3 .php4 .php5 .cgi
[...]
```

And here's a configuration example to adapt to your needs:

```bash
server {
    listen 80;
    listen 443 ssl;

    ssl_certificate /etc/nginx/ssl/deimos.fr/server-unified.crt;
    ssl_certificate_key /etc/nginx/ssl/deimos.fr/server.key;
    ssl_session_timeout 5m;

    server_name git.deimos.fr;
    root /usr/share/gitweb/;

    access_log /var/log/nginx/git.deimos.fr_access.log;
    error_log /var/log/nginx/git.deimos.fr_error.log;

    index gitweb.cgi;

    location /gitweb.cgi {
        fastcgi_cache mycache;
        fastcgi_cache_key $request_method$host$request_uri;
        fastcgi_cache_valid any 1h;
        include fastcgi_params;
        fastcgi_pass  unix:/run/fcgiwrap.socket;
    }
}
```

Then we'll need a cgi wrapper:

```bash
aptitude install fcgiwrap
```

Create the link, then reload the configuration:

```bash
cd /etc/nginx/sites-enabled
ln -s /etc/nginx/sites-available/git.deimos.fr .
/etc/init.d/nginx reload
/etc/init.d/fcgiwrap restart
/etc/init.d/php5-fpm reload
```

### Gitweb

Adapt this file according to your needs:

```bash
# path to git projects (<project>.git)
$projectroot = "/var/lib/git";

# directory to use for temp files
$git_temp = "/tmp";

# target of the home link on top of all pages
#$home_link = $my_uri || "/";

# html text to include at home page
$home_text = "indextext.html";

# file with project list; by default, simply scan the projectroot dir.
$projects_list = $projectroot;

# stylesheet to use
$stylesheet = "/gitweb.css";

# logo to use
$logo = "/git-logo.avif";

# the 'favicon'
$favicon = "/git-favicon.png";

# change default git logo url
$logo_url = "http://www.deimos.fr/gitweb";
$logo_label = "Deimos.fr Git Repository";

# This prevents gitweb to show hidden repositories
$export_ok = "git-daemon-export-ok";
$strict_export = 1;

# This lets it make the URLs you see in the header
@git_base_url_list = ( 'git://www.deimos.fr/git' );
```

#### Restrict access

It's possible to disable access to certain repositories. To do this, enable this in your gitweb configuration:

```bash
$export_ok = "gitweb-export-ok";
```

Then create this file in each repository that you want to display:

```bash
touch /var/cache/git/git_deimosfr.git/gitweb-export-ok
```

Only repositories with this file will be displayed.

#### Change the header

If you want to change the header of your gitweb, create an indextext.html file at the location of the cgi and insert HTML code:

```html
<h2>Deimos Git</h2>

Welcome on my Gitweb. I store here my configurations that I usually use every
days.<br /><br />

Other links:<br />
- My blog: <a href="http://www.deimos.fr/blog">http://www.deimos.fr/blog</a
><br />
- My wiki:
<a href="http://www.deimos.fr/blocnotesinfo"
  >http://www.deimos.fr/blocnotesinfo</a
><br />
```

## Web interface

Your server is now accessible via the following web interface: http://server/gitweb

### More attractive theme

There is a more attractive theme than the default gitweb theme. [Download the archive here](https://github.com/kogakure/gitweb-theme/) and run these commands:

```bash
tar -xzvf kogakure-gitweb-theme.tar.gz
cd kogakure-gitweb-theme
mv /usr/share/gitweb/gitweb.css /usr/share/gitweb/gitweb-old.css
cp gitweb.css /usr/share/gitweb/
cp -Rf img /usr/share/gitweb/
chown -Rf www-data. /usr/share/gitweb/
```

And there you go, your new theme is in place.

### Piwik integration

If you want integration with Piwik, it's quite simple. I've made a patch - you'll need to modify the JavaScript code to display in your page in this patch:

```diff
*** gitweb.old	2011-04-05 14:05:06.120951481 +0200
--- gitweb.cgi	2011-04-05 14:04:41.913944817 +0200
***************
*** 3612,3617 ****
--- 3612,3633 ----
  		      qq!</script>\n!;
  	}

+
+ 	print <<PIWIK;
+ <!-- Piwik -->
+ <script type="text/javascript">
+ var pkBaseURL = (("https:" == document.location.protocol) ? "https://www.deimos.fr/piwik/" : "http://www.deimos.fr/piwik/");
+ document.write(unescape("%3Cscript src='" + pkBaseURL + "piwik.js' type='text/javascript'%3E%3C/script%3E"));
+ </script><script type="text/javascript">
+ try {
+ var piwikTracker = Piwik.getTracker(pkBaseURL + "piwik.php", 10);
+ piwikTracker.trackPageView();
+ piwikTracker.enableLinkTracking();
+ } catch( err ) {}
+ </script><noscript><p><img src="http://www.deimos.fr/piwik/piwik.php?idsite=10" style="border:0" alt="" /></p></noscript>
+ <!-- End Piwik Tracking Code -->
+ PIWIK
+
  	print "\n</body>\n" .
  	      "</html>";
  }
***************
*** 7033,7038 ****
--- 7049,7066 ----
  	}
  	print <<XML;
  </outline>
+ <!-- Piwik -->
+ <script type="text/javascript">
+ var pkBaseURL = (("https:" == document.location.protocol) ? "https://www.deimos.fr/piwik/" : "http://www.deimos.fr/piwik/");
+ document.write(unescape("%3Cscript src='" + pkBaseURL + "piwik.js' type='text/javascript'%3E%3C/script%3E"));
+ </script><script type="text/javascript">
+ try {
+ var piwikTracker = Piwik.getTracker(pkBaseURL + "piwik.php", 10);
+ piwikTracker.trackPageView();
+ piwikTracker.enableLinkTracking();
+ } catch( err ) {}
+ </script><noscript><p><img src="http://www.deimos.fr/piwik/piwik.php?idsite=10" style="border:0" alt="" /></p></noscript>
+ <!-- End Piwik Tracking Code -->
  </body>
  </opml>
  XML
```

## FAQ

### Unnamed repository; edit this file to name it for gitweb.

I have this displayed at the beginning of my web interface, how do I change it?

It's quite simple. In your git repository, you have a 'description' file. Edit it and put whatever you want in it.
