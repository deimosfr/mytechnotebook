---
weight: 999
url: "/Tiny_Tiny_RSS_\\:_Une_alternative_à_Google_Reader/" 
title: "Tiny Tiny RSS: An Alternative to Google Reader"
description: "How to set up and configure Tiny Tiny RSS as an alternative to Google Reader, including installation, MySQL configuration, and importing from Google Reader"
categories: ["MySQL", "Database", "Linux"]
date: "2013-06-29T16:26:00+02:00"
lastmod: "2013-06-29T16:26:00+02:00"
tags: ["RSS", "Tiny Tiny RSS", "Google Reader", "MySQL", "Database", "Web Services"]
toc: true
---

![Tiny Tiny RSS Logo](/images/tinytinyrss_logo.avif)

## Introduction

For years I've been happily using Google Reader. But increasingly, I'm reading from RSS feeds about people who have lost their data through Google. Even in enterprise mode, it's not possible today to easily and automatically back up your data.

I've therefore chosen not to leave my data with Google and to take care of it myself. I've selected [Tiny Tiny RSS](https://tt-rss.org) which aims to be a very good alternative to Google Reader.

## Installation

For prerequisites, we need a LAMP-type server with some specific options:

```bash
aptitude install php5-curl php5-cli
```

Then we'll install Tiny Tiny RSS:

```bash
cd /var/www
wget http://tt-rss.org/download/tt-rss-1.6.2.tar.gz
tar -xzf tt-rss-1.6.2.tar.gz
rm -f tt-rss-1.6.2.tar.gz
mv tt-rss-1.6.2 tinyrss
chown -Rf www-data. tinyrss
```

### MySQL

I've chosen MySQL as a backend for Tiny Tiny RSS. We'll initialize the database and its users (adapt as needed):

```bash
CREATE DATABASE `tinyrss` ;
CREATE USER 'tinyrss_user'@'localhost' IDENTIFIED BY '***';
GRANT USAGE ON * . * TO 'tinyrss_user'@'localhost' IDENTIFIED BY '***';
GRANT ALL PRIVILEGES ON `tinyrss` . * TO 'tinyrss_user'@'localhost';
```

Then we'll import the SQL schema:

```bash
mysql -utinyrss_user -p tinyrss < schema/ttrss_schema_mysql.sql
```

## Configuration

For configuration, we'll use a config file:

```bash
mv config.php{-dist,}
```

Then we'll edit this file to insert our database values:

```php
    // Database server configuration
    define('DB_TYPE', "mysql"); // or pgsql
    define('DB_HOST', "localhost");
    define('DB_USER', "tinyrss_user");
    define('DB_NAME', "tinyrss");
    define('DB_PASS', "***");
    define('SELF_URL_PATH', 'http://www.deimos.fr/tinyrss');
    define('SINGLE_USER_MODE', false);
```

You can then connect to your server with admin/password.

### Updating Feeds

To update your feeds and manage them from the web interface, we'll use the provided daemon:

```bash
su - www-data -c "php /var/www/tinytinyrss/update.php --daemon --quiet &"
```

### Exporting Google Reader Feeds

If like me you're migrating from Google Reader to Tiny Tiny RSS, you can import your data.

To export your Google Reader data in OPML format via the command line, create a small script:

```bash
#!/bin/sh
curl -sH "Authorization: GoogleLogin auth=$(curl -sd "Email=$1&Passwd=$2&service=reader" https://www.google.com/accounts/ClientLogin | grep Auth | sed 's/Auth=\(.*\)/\1/')" http://www.google.com/reader/subscriptions/export;
```

Then to use it, run the script with your login and password as arguments:

```bash
/backup_google_reader.sh user@gmail.com password > google-reader-subscriptions.xml
```
