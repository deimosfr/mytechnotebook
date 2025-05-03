---
weight: 999
url: "/Symfony_\\:_Installation_et_configuration_du_framework_PHP/"
title: "Symfony: Installation and Configuration of the PHP Framework"
description: "A guide to install and configure the Symfony PHP framework, including prerequisites, installation steps, and basic project configuration."
categories: ["Development", "PHP", "MySQL", "Linux", "Apache"]
date: "2010-05-02T20:44:00+02:00"
lastmod: "2010-05-02T20:44:00+02:00"
tags: ["Symfony", "PHP", "Framework", "MySQL", "Apache", "Development"]
toc: true
---

## Introduction

[Symfony](https://www.symfony-project.org) is a free MVC framework written in PHP 5. As a framework, it facilitates and accelerates the development of Internet and Intranet websites and applications.

## Prerequisites

Before starting, we need a web server of your choice (I'm using Apache2) and a database server (MySQL):

```bash
aptitude install apache2 libapache2-mod-php5 php5 mysql-server
```

Symfony provides a small utility to test your server:

```bash
> wget http://sf-to.org/1.4/check.php
> php check_configuration.php 
********************************
*                              *
*  symfony requirements check  *
*                              *
********************************

php.ini used by PHP: /etc/php5/cli/php.ini

** WARNING **
*  The PHP CLI can use a different php.ini file
*  than the one used with your web server.
*  If this is the case, please launch this
*  utility from your web server.
** WARNING **

** Mandatory requirements **

  OK        PHP version is at least 5.2.4 (5.2.6-1+lenny8)

** Optional checks **

  OK        PDO is installed
[[WARNING]] PDO has some drivers installed: : FAILED
            *** Install PDO drivers (mandatory for Propel and Doctrine) ***
  OK        PHP-XML module is installed
[[WARNING]] XSL module is installed: FAILED
            *** Install and enable the XSL module (recommended for Propel) ***
  OK        The token_get_all() function is available
  OK        The mb_strlen() function is available
  OK        The iconv() function is available
  OK        The utf8_decode() is available
  OK        The posix_isatty() is available
[[WARNING]] A PHP accelerator is installed: FAILED
            *** Install a PHP accelerator like APC (highly recommended) ***
[[WARNING]] php.ini has short_open_tag set to off: FAILED
            *** Set it to off in php.ini ***
[[WARNING]] php.ini has magic_quotes_gpc set to off: FAILED
            *** Set it to off in php.ini ***
  OK        php.ini has register_globals set to off
  OK        php.ini has session.auto_start set to off
  OK        PHP version is not 5.2.9
```

Clearly, there are some small issues. Let's fix them now:

```bash
aptitude install php5-mysql php5-xsl php-apc
perl -pe 's/^(short_open_tag = )on/\1off/i' < /etc/php5/cli/php.ini > /tmp/sym_php_changes_tmp
perl -pe 's/^(magic_quotes_gpc = )on/\1off/i' < /tmp/sym_php_changes_tmp > /etc/php5/cli/php.ini
```

Now, if we check again:

```bash
> php check_configuration.php 
********************************
*                              *
*  symfony requirements check  *
*                              *
********************************

php.ini used by PHP: /etc/php5/cli/php.ini

** WARNING **
*  The PHP CLI can use a different php.ini file
*  than the one used with your web server.
*  If this is the case, please launch this
*  utility from your web server.
** WARNING **

** Mandatory requirements **

  OK        PHP version is at least 5.2.4 (5.2.6-1+lenny8)

** Optional checks **

  OK        PDO is installed
  OK        PDO has some drivers installed: mysql
  OK        PHP-XML module is installed
  OK        XSL module is installed
  OK        The token_get_all() function is available
  OK        The mb_strlen() function is available
  OK        The iconv() function is available
  OK        The utf8_decode() is available
  OK        The posix_isatty() is available
  OK        A PHP accelerator is installed
  OK        php.ini has short_open_tag set to off
  OK        php.ini has magic_quotes_gpc set to off
  OK        php.ini has register_globals set to off
  OK        php.ini has session.auto_start set to off
  OK        PHP version is not 5.2.9
```

Everything is ok :-)

## Installation

There are several solutions to install Symfony. I personally chose SVN, but PEAR would have been just as good, or even the Sandbox (all-in-one). For this, I need to have SVN installed:

```bash
aptitude install subversion
```

Next, we will create a space to place Symfony:

```bash
cd /usr/share
svn checkout http://svn.symfony-project.com/branches/1.4/
mv 1.4 symfony
```

Then we'll verify that everything is installed correctly:

```bash
symfony/data/bin/symfony -V
symfony version 1.4.5-DEV (/usr/share/symfony/lib)
```

Yippee :-)

## Configuration

Let's update our PATH to easily use the new binaries:

```bash
export PATH=$PATH:/usr/share/symfony/data/bin
```

### Initializing a New Project

Let's imagine I have a project called 'phpwol'. I'll create my project folder and initialize it:

```bash
mkdir -p /var/www/phpwol
cd /var/www/phpwol
symfony generate:project phpwol
```

A lot of things have just been created:

{{< table "table-hover table-striped" >}}
| Directory | Description |
|---------|--------|
| `apps/` | Contains all project applications |
| `cache/` | Files cached by the framework |
| `config/` | Project configuration files |
| `data/` | Data files such as initial data sets |
| `lib/` | Project libraries and classes |
| `log/` | Framework log files |
| `plugins/` | Installed plugins |
| `test/` | Unit and functional test files |
| `web/` | The Web root directory (see below) |
{{< /table >}}

## Resources
- http://www.symfony-project.org/getting-started/1_4/fr/
- http://www.lafermeduweb.net/tutorial/symfony-creer-un-site-web-avec-le-framework-php-symfony-14.html
