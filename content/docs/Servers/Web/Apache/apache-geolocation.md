---
weight: 999
url: "/Geolocalisation_avec_Apache_2/"
title: "Geolocation with Apache 2"
description: "A guide to set up and use mod_geoip with Apache 2 for geographical targeting and IP-based filtering."
categories: ["Servers", "Apache", "Web"]
date: "2008-03-16T08:41:00+02:00"
lastmod: "2008-03-16T08:41:00+02:00"
tags: ["Apache", "mod_geoip", "Geolocation", "Web Server"]
toc: true
---

## Introduction

This guide explains how to set up [mod_geoip](https://www.maxmind.com/app/mod_geoip) with Apache2 on a Debian Etch system. mod_geoip looks up the IP address of the client end user. This allows you to redirect or block users based on their country. You can also use this technology for your [OpenX](https://www.openx.org/) (formerly known as OpenAds or phpAdsNew) ad server to allow [geo targeting](https://en.wikipedia.org/wiki/Geo_targeting).

## Installing mod_geoip

To install mod_geoip, we simply run:

```bash
apt-get install libapache2-mod-geoip
```

Then we open `/etc/apache2/mods-available/geoip.conf` and uncomment the GeoIPDBFile line so that the file looks as follows:

```bash
vi /etc/apache2/mods-available/geoip.conf
```

```
<IfModule mod_geoip.c>
  GeoIPEnable On
  GeoIPDBFile /usr/share/GeoIP/GeoIP.dat
</IfModule>
```

Next we restart Apache:

```bash
/etc/init.d/apache2 restart
```

That's it already!

## A Short Test

To see if mod_geoip is working correctly, we can create a small PHP file in one of our web spaces (e.g. `/var/www`):

```bash
vi /var/www/geoiptest.php
```

```php
<html>
<body>
<?php
$country_name = apache_note("GEOIP_COUNTRY_NAME");
print "Country: " . $country_name;
?>
</body>
</html>
```

Call that file in a browser, and it should display your country (**make sure that you're calling the file from a public IP address, not a local one**).

## Use Cases

You can use mod_geoip to redirect or block/allow users based on their country.
You can also use mod_geoip with OpenX/OpenAds/phpAdsNew.

## References

http://www.howtoforge.com/mod-geoip-apache2-debian-etch  
http://www.maxmind.com/app/mod_geoip  
http://www.maxmind.com/openads_geoip.pdf
