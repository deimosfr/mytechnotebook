---
weight: 999
url: "/Piwik_\\:_Des_statistiques_pour_votre_site_web/"
title: "Piwik: Statistics for Your Website"
description: "A guide to installing and configuring Piwik (now Matomo), an open-source alternative to Google Analytics for website statistics tracking and analysis."
categories: ["Nginx", "Linux", "MySQL"]
date: "2014-07-12T09:47:00+02:00"
lastmod: "2014-07-12T09:47:00+02:00"
tags: ["Piwik", "Analytics", "Web", "Statistics", "Nginx", "MySQL", "PHP", "WordPress", "MediaWiki", "Gitweb"]
toc: true
---

![Piwik Logo](/images/piwik-logo.avif)

## Introduction

I've been using [Piwik](https://piwik.org/) for over a year and hadn't written an article about it yet. This is the opportunity to show you this equivalent to Google Analytics.

## Prerequisites

Set up a database on a MySQL instance that you'll use during the installation.

## Installation

Navigate to your web directory:

```bash
cd /var/www
wget http://piwik.org/latest.zip
unzip latest.zip
rm How\ to\ install\ Piwik.html latest.zip
chown -Rf www-data. piwik
```

Then launch the installer at `http://server/piwik`

## Configuration

### Web Server

Here are some web server configurations you might need.

#### Nginx

For Piwik, here's the configuration (`/etc/nginx/sites-available/piwik.deimos.fr`):

```bash
server {
    include listen_port.conf;
    listen 443 ssl;
 
    ssl_certificate /etc/nginx/ssl/deimos.fr/server-unified.crt;
    ssl_certificate_key /etc/nginx/ssl/deimos.fr/server.key;
    ssl_session_timeout 5m;
 
    server_name piwik.deimos.fr;
    root /usr/share/nginx/www/deimos.fr/piwik;
    index index.php;
 
    access_log /var/log/nginx/piwik.deimos.fr_access.log;
    error_log /var/log/nginx/piwik.deimos.fr_error.log;
 
    # Drop config
    include drop.conf; 
 
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
 
    location ~* \.(js|css|png|jpg|jpeg|gif|ico)$ {
        expires max;
        log_not_found off;
    }
}
```

### Piwik

We'll create a cron job to optimize Piwik's performance:

```bash
5 * * * * php /var/www/deimos.fr/piwik/console core:archive --url=https://piwik.deimos.fr >> /var/log/piwik.log
```

Once the crontab is added, we'll configure the web interface to stop automatically processing data when accessing a report. Go to:

* Settings > General Settings tab
* Let Piwik archiving trigger when reports are viewed from a browser: **No**
* Reports for today (or any date range including today) will be processed at most every: **3600 seconds**

The Piwik configuration is now complete. You need to add sites and configure them to send stats to Piwik.

### MediaWiki

Piwik is an equivalent of Google Analytics, but free. It enables you to have statistics for your website with nice graphs etc. There is a [plugin](https://www.mediawiki.org/wiki/Extension:Piwik_Integration) that allows you to insert code in each page (necessary), but it's obsolete and has security vulnerabilities. We'll use another module that will simply allow us to insert this type of code.

You need to install the [PCR GUI Inserts](https://www.mediawiki.org/wiki/Extension:PCR_GUI_Inserts) module which allows you to insert information at various locations on your pages. First, activate the extension by adding these lines:

```php
# PCR Extension for Piwik / Google Ads
require_once("$IP/extensions/pcr/pcr_guii.php");
```

Insert this into the MediaWiki configuration file:

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

### Gitweb

If you want to integrate with Piwik, it's quite simple. I created a patch - you'll need to modify the JavaScript code to display in your page:

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

### WordPress

For WordPress, install the [Piwik Analytics](https://forwardslash.nl/piwik-analytics/) extension which will allow you to easily configure Piwik for your blog.
