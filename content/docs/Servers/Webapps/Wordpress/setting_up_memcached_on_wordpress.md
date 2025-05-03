---
weight: 999
url: "/Mise_en_place_de_Memcached_sur_Wordpress/"
title: "Setting up Memcached on WordPress"
description: "How to set up Memcached on WordPress to improve website performance"
categories:
  - "Debian"
  - "Linux"
date: "2009-09-08T21:38:00+02:00"
lastmod: "2009-09-08T21:38:00+02:00"
tags:
  - "WordPress"
  - "Memcached"
  - "Performance"
  - "Web"
  - "PHP"
toc: true
---

## Introduction

Memcached is a cache server that, unlike some PHP accelerators, does not consume additional CPU. It's therefore an ideal compromise. For WordPress, there is currently no simple solution to quickly set up this solution (although it's not very time-consuming anyway).

## Installation

On Debian, it's easy:

```bash
apt-get install memcached
```

Your memcached server is now running on port 11211.

## Configuration

### Server

Nothing special needs to be done, the basic configuration is sufficient.

### WordPress

Go to the wp-content folder of your WordPress, then download these files and set the proper permissions:

```bash
cd ./wp-content
wget http://svn.wp-plugins.org/memcached/branches/1.0/memcached-client.php
wget http://svn.wp-plugins.org/memcached/branches/1.0/object-cache.php
chown www-data. memcached-client.php object-cache.php
cd ..
```

Then modify the WordPress configuration file wp-config.php and add this line (`wp-config.php`):

```php
$memcached_servers = array('127.0.0.1:11211');
```

If you have multiple memcached servers, here's the syntax to use (`wp-config.php`):

```php
$memcached_servers = array('192.168.1.1:11211', '192.168.1.2:11211');
```

## Resources
- http://ryan.wordpress.com/2005/12/23/memcached-backend/
