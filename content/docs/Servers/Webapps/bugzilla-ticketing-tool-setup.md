---
weight: 999
url: "/Bugzilla_\\:_mise_en_place_d'un_outil_de_ticketing/"
title: "Bugzilla: Setting Up a Ticketing Tool"
description: "How to install and configure Bugzilla as a ticketing and bug tracking system"
categories: ["Ticketing", "Web Applications", "Bug Tracking"]
date: "2013-07-03T14:37:00+02:00"
lastmod: "2013-07-03T14:37:00+02:00"
tags: ["Bugzilla", "Ticketing", "Bug Tracking", "Apache", "MySQL", "Perl"]
toc: true
---

![Bugzilla](/images/bugzilla_logo.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 4.2.5 |
| **Operating System** | Debian 6 |
| **Website** | [Bugzilla Website](https://www.bugzilla.org/) |
| **Last Update** | 03/07/2013 |
{{< /table >}}

## Introduction

Bugzilla is a free web-based bug tracking system with a web interface, developed and used by the Mozilla Organization. It allows tracking of bugs or "Request for Enhancement" (RFE) in the form of "tickets". It is a server-side software with a three-tier architecture, written in Perl. It is available on UNIX (Linux, BSD, etc.) and is distributed under the tri-license MPL/LGPL/GPL.

It is used by many organizations to track the development of numerous software applications, on the Internet or in private networks. The most well-known users include the Mozilla Foundation, Facebook, NASA, YAHOO, GNOME, KDE, Red Hat, Novell, and Mandriva.[^1]

## Installation

We will need a web server to set up Bugzilla. We'll keep it simple and use Apache, along with a MySQL database:

```bash
aptitude install apache2 mysql-server libapache2-mod-perl2
```

Then we'll install all the Bugzilla dependencies:

```bash
aptitude install libtimedate-perl libdatetime-perl libtemplate-perl libemail-send-perl libemail-mime-perl libdbi-perl liburi-perl libmath-random-isaac-perl libdbd-mysql-perl libgd-gd2-perl libchart-perl libtemplate-plugin-gd-perl libmime-tools-perl libwww-perl libxml-twig-perl libnet-ldap-perl libauthen-sasl-perl libauthen-radius-perl libsoap-lite-perl libjson-rpc-perl libtest-taint-perl libhtml-scrubber-perl libencode-detect-perl
```

Then we'll download the latest version of Bugzilla and extract it:

```bash
cd /var/www
wget http://ftp.mozilla.org/pub/mozilla.org/webtools/bugzilla-4.2.5.tar.gz
tar -xzf bugzilla-4.2.5.tar.gz
```

If you're still missing modules, you can use the non-Debian-integrated solution that will install everything needed:

```bash
cd bugzilla-4.2.5
/usr/bin/perl install-module.pl --all
```

## Configuration

### MySQL

For MySQL configuration, we'll tune some parameters by adding these lines to the configuration:

```bash
[mysqld]
# Allow packets up to 4MB
max_allowed_packet=4M
# Allow small words in full-text indexes
ft_min_word_len=2
```

Restart MySQL for these parameters to take effect. Then create the database, a user, and its MySQL permissions (replace the password with what you want):

```sql
CREATE DATABASE bugs;
CREATE USER 'bugs'@'localhost' IDENTIFIED BY 'bugs';
GRANT SELECT, INSERT, UPDATE, DELETE, INDEX, ALTER, CREATE, LOCK TABLES, CREATE TEMPORARY TABLES, DROP, REFERENCES ON bugs.* TO bugs@localhost IDENTIFIED BY 'bugs';
FLUSH PRIVILEGES;
```

### Apache

For Apache to handle perl/CGI scripts, we've installed mod_perl. Now we need to configure the Apache directories:

```apache {linenos=table,hl_lines=[3,4,5,6,7,8]}
<VirtualHost *:80>
[...]
        <Directory /var/www/bugzilla-4.2.5>
                AddHandler cgi-script .cgi
                Options +Indexes +ExecCGI
                DirectoryIndex index.cgi
                AllowOverride Limit FileInfo Indexes
        </Directory>
[...]
</VirtualHost>
```

Reload your Apache configuration afterward.

### Bugzilla

We need to run the check tool for the first time to create our configuration file:

```bash
./checksetup.pl
```

Now, set your variables in the configuration file to match your database and web server information:

```perl
[...]
$webservergroup = 'www-data';
$db_name = 'bugs';
$db_user = 'bugs';
$db_pass = 'bugs';
[...]
```

Feel free to adapt with your own information.

Then, rerun the configuration tool, which will create everything you need for your database:

```bash
> ./checksetup.pl
Adding new table bz_schema...
Initializing bz_schema...
Creating tables...
Converting attach_data maximum size to 100G...
Setting up choices for standard drop-down fields:
   priority bug_status rep_platform resolution bug_severity op_sys
Creating ./data directory...
Creating ./data/attachments directory...
Creating ./data/db directory...
Creating ./data/extensions directory...
Creating ./data/mining directory...
[...]
Enter the e-mail address of the administrator: xxx@mycompany.com
Enter the real name of the administrator: Deimos
Enter a password for the administrator account:
Please retype the password to verify:
xxx@mycompany.com is now set up as an administrator.
```

It should have now asked for and validated the Administration information (login + password).

#### Crontab

We'll also create a crontab that will execute a set of scripts for graphs and whining:

```bash
#!/bin/sh
5 0 * * * cd /var/www/bugzilla && ./collectstats.pl
55 0 * * * cd /var/www/bugzilla && ./whineatnews.pl
*/15 * * * * cd /var/www/bugzilla && ./whine.pl
```

Add the proper permissions and restart the cron service:

```bash
chmod 755 /etc/cron.daily/bugzilla
service cron restart
```

### Adding an Administrator

If you want to add an Admin, it's simple. The account must already exist on Bugzilla, then:

```bash
> ./checksetup.pl --make-admin=xxx@mycompany.com
[...]
Removing existing compiled templates...
Precompiling templates...done.
Fixing file permissions...

xxx@mycompany.com is now set up as an administrator.

Now that you have installed Bugzilla, you should visit the 'Parameters'
page (linked in the footer of the Administrator account) to ensure it
is set up as you wish - this includes setting the 'urlbase' option to
the correct URL.
```

Adjust this command with the email address of the user you want to make an Admin. After that, all you need to do is connect to it: http://myserver/bugzilla :-)

## References

[^1]: http://fr.wikipedia.org/wiki/Bugzilla
