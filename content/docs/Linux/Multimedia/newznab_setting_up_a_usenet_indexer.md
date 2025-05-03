---
weight: 999
url: "/Newznab_\\:_Mise_en_place_d'un_indexeur_de_usenet/"
title: "Newznab: Setting up a Usenet Indexer"
description: "A comprehensive guide on how to install and configure Newznab as a Usenet indexer to work with applications like SABnzbd and Sick-Beard."
categories: ["Debian", "Database", "Linux"]
date: "2012-12-28T16:47:00+02:00"
lastmod: "2012-12-28T16:47:00+02:00"
tags: ["Newznab", "Website", "Usenet", "SABnzbd", "Database", "Apache", "PHP", "MySQL", "MariaDB"]
toc: true
---

![Newznab](/images/newznab_logo.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 0.2.3 |
| **Operating System** | Debian 6 |
| **Website** | [Newznab Website](https://www.newznab.com/) |
| **Last Update** | 28/12/2012 |
{{< /table >}}

## Introduction

Newznab is a Usenet indexer. It can be used with applications like [SABnzbd](../sabnzbd_\:_une_interface_web_pour_gérer_les_newsgroups/) or [Sick-Beard](../Sick-Beard_\:_Un_PVR_s'appuyant_sur_SABnzbd/).

## Installation

Here are the required packages:

```bash
aptitude install apache2 php5 libapache2-mod-php5 php5-curl php5-mysql php5-gd php-pear
```

Since we'll need a MySQL-type database, we'll take the opportunity to set up [MariaDB (follow this link)](../mariadb_\:_migration_depuis_mysql/).

Next, let's set up Newznab:

```bash
cd /var/www
wget http://www.newznab.com/newznab-0.2.3.zip
unzip newznab-0.2.3.zip
mv newznab-0.2.3 newznab
chown -Rf www-data. newznab
rm -f newznab-0.2.3.zip
```

Finally, let's enable Apache's rewrite mode:

```bash
a2enmod rewrite
```

## Configuration

### PHP

We'll need to configure some PHP parameters for Newznab to work properly:

```bash
perl -pi -e 's/;date.timezone =.*/date.timezone = Europe\/Paris/' /etc/php5/cli/php.ini
perl -pi -e 's/;date.timezone =.*/date.timezone = Europe\/Paris/' /etc/php5/apache2/php.ini
perl -pi -e "s/max_execution_time = \d*/max_execution_time = 120/" /etc/php5/cli/php.ini
perl -pi -e "s/max_execution_time = \d*/max_execution_time = 120/" /etc/php5/apache2/php.ini
perl -pi -e "s/memory_limit = .*/memory_limit = 256M/" /etc/php5/apache2/php.ini
```

### Apache

Here's a typical configuration. Feel free to adapt it to your needs:

```apache
<VirtualHost *:80>
        <Directory /var/www/newznab/www>
                Options FollowSymLinks
                AllowOverride All
                Order allow,deny
                allow from all
        </Directory>
        ServerAdmin admin@example.com
        ServerName example.com
        ServerAlias www.example.com
        DocumentRoot /var/www/newznab/www
        LogLevel warn
        ServerSignature Off
</VirtualHost>
```

Then enable it and reload your Apache configuration:

```bash
a2ensite newznab
/etc/init.d/apache2 reload
```

### MySQL

Let's configure MySQL/MariaDB by creating a database and a user:

```bash
create database newznab;
create user 'newznab'@'localhost' identified by 'password';
grant usage ON * . * TO 'newznab'@'localhost' IDENTIFIED BY 'password';
grant ALL on `newznab` .* to 'newznab'@'localhost';
```

Replace password with your desired password.

### Newznab

After all these modifications, restart your Apache server:

```bash
/etc/init.d/apache2 restart
```

Then connect to your server at http://server/newznab/www and follow the instructions. Once you have added a group, you'll need to run a script to retrieve the latest headers. The result should be similar to this:

```bash
> cd /var/www/newznab/misc/update_scripts
> php update_binaries.php
PHP Warning:  mysql_pconnect(): Headers and client library minor version mismatch. Headers:50149 Library:50311 in /var/www/newznab/www/lib/framework/db.php on line 12

newznab 0.2.3 Copyright (C) 2012 newznab.com


This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation.


This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

Updating: 1 groups - Using compression? No
Processing alt.binaries.teevee
Group alt.binaries.teevee has 50,001 new parts.
First: 334617788 Last: 454617805 Local last: 0
New group starting with 50000 messages worth.
Getting 20,001 parts (454567805 to 454587805) - 30,000 in queue
Received 20001 articles of 20001 requested, 0 blacklisted, 0 not binaries 
500 bin adds...1000 bin adds...1500 bin adds...2000 bin adds...
2,456 new, 0 updated, 20,001 parts. 4.50 headers, 6.17 update, 10.67 range.
Getting 20,001 parts (454587806 to 454607806) - 9,999 in queue
Received 20001 articles of 20001 requested, 0 blacklisted, 0 not binaries 
500 bin adds...1000 bin adds...1500 bin adds...
1,989 new, 268 updated, 20,001 parts. 4.05 headers, 6.98 update, 11.03 range.
Getting 9,999 parts (454607807 to 454617805) - 0 in queue
Received 9999 articles of 9999 requested, 0 blacklisted, 1 not binaries 
500 bin updates...
487 new, 754 updated, 9,998 parts. 2.84 headers, 4.27 update, 7.11 range.
Group processed in 29.74 seconds 

Updating completed in 30.09 seconds
```

Next, we'll need to create the releases:

```bash
> php update_releases.php
PHP Warning:  mysql_pconnect(): Headers and client library minor version mismatch. Headers:50149 Library:50311 in /var/www/newznab/www/lib/framework/db.php on line 12

newznab 0.2.3 Copyright (C) 2012 newznab.com


This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation.


This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.


Starting release update process (2012-12-16 00:42:54)
updated regexes to revision 1742
Applying regex 2 for group alt.binaries.*
Applying regex 1 for group alt.binaries.*
Stage 2
Stage 3
Found 0 nfos in 0 releases
0 nfo files processed
Site config (site.checkpasswordedrar) prevented checking releases are passworded
Processing tv for 0 releases
Lookup tv rage from the web - Yes
Tidying away binaries which cant be grouped after 2 days
Deleting parts which are older than 2 days
Deleting binaries which are older than 2 days
Processed 0 releases
```

Once these actions have been executed manually, we'll set them up automatically in crontab:

```bash
0 * * * * /var/www/newznab/misc/update_scripts/cron_scripts/newznab.sh start > /dev/null
```

We'll also need to set permissions and modify default paths:

```bash
chmod 755 /var/www/newznab/misc/update_scripts/cron_scripts/newznab.sh
perl -pi -e 's/\/usr\/local\/www/\/var\/www/' /var/www/newznab/misc/update_scripts/cron_scripts/newznab.sh
```

## Usage

### Import a server list

The default server list is decent, but may not meet your expectations. I've created a small example script that retrieves a list of servers from a website, creates a "clean" list, checks if they're already in the database, then adds them all at once and activates them at the same time. Here's my script:

```bash
#!/bin/sh

# Vars
ngs_list='liste_newsgroups'
mysql_batch='mysql_batch.sql'
database='newznab'
login='newznab'
password='password'
activate=1
url='http://www.binnews.in/_bin/newsgroup.php?country=fr'

# Getting URL
echo 'Getting newsgroups servers from URL'
wget -O $ngs_list.old "$url"
echo "Extracting all newsgroups servers"
grep title $ngs_list.old | awk -F\" '{ print $2 }' > $ngs_list

# Generate sql batch
echo "Generating SLQ batch file in $mysql_batch"
echo "" > $mysql_batch
for i in `cat $ngs_list` ; do
	if [ `mysql -u$login -p$password -e "select name from newznab.groups where name = '$i'" | grep -c $i` -eq 0 ] ; then
		echo "insert into $database.groups (name,backfill_target,first_record,first_record_postdate,last_record,last_record_postdate,last_updated,minfilestoformrelease,minsizetoformrelease,active,description) values ('$i',1,0,NULL,0,NULL,NULL,NULL,NULL,$activate,'$i');" >> $mysql_batch
	fi
done

# Import SQL
echo "Importing batch file in $database database"
mysql -u$login -p$password newznab < $mysql_batch
```

Just fill in the vars section with the correct elements, set the permissions, and run this script for a massive insertion.
