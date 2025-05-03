---
weight: 999
url: "/Awstats_\\:_Mise_en_place_d'Awstats,_interpréteur_de_logs_web/"
title: "Awstats: Setting up Awstats, a Web Logs Interpreter"
description: "Learn how to install and configure Awstats, a powerful web log analyzer that provides graphical reports on website traffic statistics."
categories: ["Server", "Web", "Analytics"]
date: "2010-04-11T15:51:00+02:00"
lastmod: "2010-04-11T15:51:00+02:00"
tags: ["awstats", "statistics", "analytics", "log", "apache", "lighttpd"]
toc: true
---

## Introduction

AWStats is a web log analyzer (but also FTP, Streaming, and mail) offering static but also dynamic graphical views of access statistics for your web servers.  
It displays the number of visits, unique visitors, pages, hits, transfers by domain/country, host, time, browser, OS, etc. It can be run using CGI scripts or command line.  
AWStats is free software under the GPL license.

The installation of Awstats is simple: this software interprets Apache logs (all commands to which the server responded) to provide you with understandable graphs and tables in the form of a web page.

## Installation

Install Awstats using the command:

```bash
aptitude install awstats
```

## Configuration

### Apache

In the configuration, I'm modifying these lines if I'm using Apache:

```bash
LogFile="/var/log/apache2/access.log"
SiteDomain="deimos.fr"
```

Displaying the web page interpreting logs will require images: you have two solutions so that the address http://your_address/awstats-icon directs to awstats images: use an Apache alias or create an awstats-icon folder in your Apache space and put the awstats images there. For this last solution, type:

```bash
mkdir /var/www/awstats-icon
cp -r /usr/share/awstats/icon/* /var/www/awstats-icon
```

### Lighttpd

Use this configuration if you're using Lighttpd:

```bash
LogFile="/var/log/lighttpd/access.log"
SiteDomain="deimos.fr"
LogFormat=1
```

We'll also need to configure lighttpd:

```perl
alias.url = (
                "/awstats-icon" => "/usr/share/awstats/icon/",
                "/awstats/" => "/usr/lib/cgi-bin/",
                "/icon/" => "/usr/share/awstats/icon/"
              )
# provide awstats cgi-bin access
$HTTP["url"] =~ "/awstats/" {
      cgi.assign = ( ".pl" => "/usr/bin/perl" )
}
```

Then we'll activate this configuration:

```bash
cd /etc/lighttpd/conf-enabled && ln -s /etc/lighttpd/conf-available/50-awstats.conf .
```

### Crontab and Multi-domains

By default, there is a line in the crontab that does its job well. But **if you have multiple domains** and therefore several configuration files in /etc/awstats/, you'll need to be a bit tricky. I recommend commenting out the current line in this cron file:

```bash
#0,10,20,30,40,50 * * * * www-data [ -x /usr/lib/cgi-bin/awstats.pl -a -f /etc/awstats/awstats.conf -a -r /var/log/apache/access.log ] && /usr/lib/cgi-bin/awstats.pl -config=awstats -update >/dev/null
```

And either add lines like:

```bash
0,10,20,30,40,50 * * * * www-data /usr/lib/cgi-bin/awstats.pl -config=deimos.fr
0,10,20,30,40,50 * * * * www-data /usr/lib/cgi-bin/awstats.pl -config=mavro.fr
```

To differentiate your domains (here deimos.fr and mavro.fr), or create a small script like this and put it in the cron instead:

```bash
#!/bin/bash
# path to cgi-bin
AWS=/usr/lib/cgi-bin/awstats.pl

# append your domain
DOMAINS="deimos.fr mavro.fr"

# loop through all domains
for d in ${DOMAINS}
do
   ${AWS} -update -config=${d}
done
```

A small chmod for execution rights and you're good to go.

### Protection of Awstats

You probably don't want your statistics to be accessible from anywhere, so you can protect yourself with htaccess or something else.

#### htaccess under Apache

This documentation explains how to protect a directory with htaccess (login + password).

Insert these lines and adapt to your configuration (/etc/apache2/sites-enabled/000-default):

```apache
<Directory /var/www/myhtaccess>
    AllowOverride AuthConfig
    Order allow,deny
    allow from all
</Directory>
```

Then create a .htaccess file in **/var/www/myhtaccess** and add this:

```apache
AuthType Basic
AuthName "Acces Prive"
AuthGroupFile /dev/null
AuthUserFile /etc/apache2/htaccesspassword

<Limit GET POST>
    Require valid-user
</Limit>

php_value magic_quotes_runtime 1
php_value magic_quotes_gpc 1
```

Then create your access file with the user (/etc/apache2/htaccesspassword):

```bash
htpasswd -c /etc/apache2/htaccesspassword username
```

For the next time, to add users, just remove "-c" like this:

```bash
htpasswd /etc/apache2/htaccesspassword username
```

Don't forget to restart apache :-)

For good documentation, follow this:
[Documentation on Htaccess](/pdf/htaccess.pdf)

#### htdigest under Lighttpd

Htaccess files don't exist in Lighttpd, but there is an equivalent. Check before starting that the mod_auth module is properly loaded.
We'll first generate (with -c for the first time, like htaccess) a file containing the credentials to be authorized to view a specific site:

```bash
htdigest -c /etc/lighttpd/.passwd 'Authorized users only' deimos
```

Here I'm creating the user deimos. The realm (here 'Authorized users only') will allow us to differentiate between different login/password files that we can have since we can only specify one for the entire server.

Then add these lines to the global lighttpd configuration:

```perl
auth.backend = "htdigest"
auth.backend.htdigest.userfile = "/etc/lighttpd/.passwd"
auth.debug = 2
```

Then I add the protection where I need it:

```perl
auth.require = ( "/docs/" =>
   (
      "method" => "digest",
      "realm" => "Authorized users only",
      "require" => "valid-user"
   )
)
```

Restart lighty and you're good.
The example above shows how to add the restriction where we need it, so we'll do it by modifying our awstats configuration:

```perl
alias.url = (
                "/awstats-icon" => "/usr/share/awstats/icon/",
                "/awstats/" => "/usr/lib/cgi-bin/",
                "/icon/" => "/usr/share/awstats/icon/"
              )
# provide awstats cgi-bin access
$HTTP["url"] =~ "/awstats/" {
    cgi.assign = ( ".pl" => "/usr/bin/perl" )
    auth.require = ( "/awstats/" =>
        (
        "method" => "digest",
        "realm" => "Trusted users only",
        "require" => "valid-user"
        )
    )   
}
```
