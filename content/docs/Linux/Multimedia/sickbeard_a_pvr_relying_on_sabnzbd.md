---
weight: 999
url: "/Sick-Beard_\\:_Un_PVR_s'appuyant_sur_SABnzbd/"
title: "Sick-Beard: A PVR Relying on SABnzbd"
description: "This guide explains how to install and configure Sick-Beard, a PVR tool that works with SABnzbd to easily manage TV series episodes on Debian systems."
categories: ["Debian", "Linux", "Ubuntu"]
date: "2012-09-26T19:26:00+02:00"
lastmod: "2012-09-26T19:26:00+02:00"
tags: ["Sick-Beard", "PVR", "SABnzbd", "Servers", "Network"]
toc: true
---

![Sick-Beard](/images/sick-beard-logo.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 830b3b1 |
| **Operating System** | Debian 7 |
| **Website** | [Sick-Beard Website](https://sickbeard.com) |
| **Last Update** | 26/09/2012 |
{{< /table >}}

## Introduction

This tool relies on [SABnzbd]({{< ref "docs/Linux/Multimedia/sabnzbd_a_web_interface_for_managing_newsgroups.md" >}}) and allows for easy management of TV series episodes.

## Installation

To install it, we need the following package:

```bash
aptitude install python-cheetah
```

Then we install Sick-Beard:

```bash
cd /tmp
wget -O sickbeard.tgz "https://nodeload.github.com/midgetspy/Sick-Beard/tarball/master"
tar -xzf sickbeard.tgz -C /usr/share/
mv /usr/share/midgetspy-Sick-Beard-* /usr/share/sick-beard
chown -Rf www-data. /usr/share/sick-beard/
```

And set up the init file:

```bash
cp /usr/share/sick-beard/init.ubuntu /etc/init.d/sickbeard
update-rc.d sickbeard defaults
```

## Configuration

We'll use a configuration file for default values:

```bash
# SickBeard configuration
SB_USER=www-data
SB_HOME=/usr/share/sick-beard
SB_DATA=/usr/share/sick-beard
SB_OPTS=/usr/share/sick-beard/config.ini
```

You can then start the service:

```bash
> /etc/init.d/sickbeard start
Removing stale /var/run/sickbeard/sickbeard.pid
Starting SickBeard
```

### Apache Redirection

You might want to have a simple URL for using this service, so let's create our configuration (`/etc/apache2/sites-enabled/000-default`):

```apache {linenos=table,hl_lines=["50-57"]}
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

    # Sabnzb
    <Location /sabnzbd>
    	order deny,allow
   	deny from all
    	allow from all
    	ProxyPass http://server:8080/sabnzbd
    	ProxyPassReverse http://server:8080/sabnzbd
    </Location>

    # Sickbeard
    <Location /sickbeard/>
    	order deny,allow
   	deny from all
    	allow from all
    	ProxyPass http://server:8081/sickbeard/
    	ProxyPassReverse http://server:8081/sickbeard/
    </Location>
</VirtualHost>
```

Next, stop Sickbeard if it's still running:

```bash
/etc/init.d/sickbeard stop
```

Then edit the configuration file to modify the web-root like this (`/usr/share/sick-beard/config.ini`):

```bash
[...]
web_root = "/sickbeard"
[...]
```

Finally, start Sickbeard and reload Apache:

```bash
/etc/init.d/sickbeard start
/etc/init.d/apache2 reload
```

## Usage

All you have to do now is type the URL: http://server/sickbeard/ and configure the software through the interface. It's quite simple, so I won't go into details.
