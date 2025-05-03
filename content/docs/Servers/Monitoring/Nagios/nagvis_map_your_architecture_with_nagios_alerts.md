---
weight: 999
url: "/NagVis\\:_Cartographier_son_architecture_avec_ses_alertes_nagios/"
title: "NagVis: Map Your Architecture with Nagios Alerts"
description: "Learn how to install and configure NagVis to display your Nagios alerts on custom infrastructure maps."
categories: ["Monitoring", "Linux", "MySQL"]
date: "2008-03-11T07:33:00+02:00"
lastmod: "2008-03-11T07:33:00+02:00"
tags: ["Nagios", "NagVis", "NDOUtils", "MySQL", "Monitoring", "Network"]
toc: true
---

## Introduction

NagVis is a really great product, easy to use but not to install. That's why I did this documentation. You also need to install NDOUtils. We'll see the installation too.

## NDOUtils

### Prepair System

Make sure your System meets the requirements:

- Nagios v2.X (stable) is running 
- mysql (Version > 4)

### Download NDOUtils

You can get the latest version of NDOUtils on the Nagios website [https://www.nagios.org/download/](https://www.nagios.org/download/)

### Unpack NDOUtils

Unpack NDOUtils in a path of your choice

```bash
tar xvzf tar ndoutils-1.4.tar.gz
```

You should have at least the following files and directories:

- Config: This directory contains all the configuration file of NDOUtils.
- db: This directory contains all the sql information to install or upgrade the mysql database.
- docs: This directory contains the technical documentation of NDOUtils.
- Include: This directory contains all the libraries.
- src: This directory contains the source code.
- configure: This file will be use to compile the program

### Compiling NDOUtils

Use the following commands to compile the NDO broker module, NDO2DB daemon, and additional utilities:

```bash
./configure
make 
```

After the compilation open the config.log file and take care that there's no failed status inside. If there is some failed status it's probably because you don't have all the libraries installed in your C compiler. Install it and recompile it.  
Note:

```bash
tail config.log
```

### Initializing the SQL Database

You have to create a database for storing the data.

```bash
mysql -u root -p
Enter password: 
Mysql> CREATE DATABASE nagios; 
Mysql>exit
```

You also have to create a specific user for nagios with the rights: SELECT, INSERT, UPDATE, and DELETE.

```bash
mysql -u root -p
Enter password: 
Mysql> GRANT SELECT,INSERT,UPDATE,DELETE,CREATE,DROP ON nagios.* TO nagios@localhost IDENTIFIED BY 'nagios';
Mysql>exit
```

Then now you can run the DB installation script in the db/ subdirectory.

```bash
cd db
./installdb
```

Make sure the script create the tables in the database (use phpmyadmin for exemple)

### Configuration

Sample config files are included in the config/ subdirectory. You have 2 files to modify:

- Ndomod.cfg, you have to change and uncomment the following lines:

```bash
…
# OUTPUT TYPE
…
#output_type=file
output_type=tcpsocket
#output_type=unixsocket
…
# OUTPUT
…
#output=/usr/local/nagios/var/ndo.dat
output=127.0.0.1
#output=/usr/local/nagios/var/ndo.sock
…
# BUFFER FILE
…
#buffer_file=/usr/local/nagios/var/ndomod.tmp
buffer_file=/var/cache/nagios2/ndomod.tmp
…
```

Feel free to change the other parameters according to your configuration.

- Ndo2db.cfg, you have to change and uncomment the following lines:

```bash
…
# USER/GROUP PRIVILIGES
…
ndo2db_user=nagios
ndo2db_group=nagios
…
# SOCKET TYPE
…
#socket_type=unix
socket_type=tcp
…
# SOCKET NAME
…
socket_name=/var/cache/nagios2/ndo.sock
…
# DATABASE SERVER TYPE
…
db_servertype=mysql
…
# DATABASE NAME
…
db_name=nagios
…
# DATABASE USERNAME/PASSWORD
…
db_user=nagios
db_pass=nagios
…
```

Feel free to change the other parameters according to your configuration.

- There is another file to modify: Nagios.cfg. This file is located in the `/etc/nagios/` directory. You have to enable the event_broker_option like this:

```bash
…
# EVENT BROKER OPTIONS
…
event_broker_options=-1
…
# EVENT BROKER MODULE(S)
…
#broker_module=/somewhere/module1.o
#broker_module=/somewhere/module2.o arg1 arg2=3 debug=0
broker_module=/usr/lib/nagios/ndoutils/ndomod.o config_file=/etc/nagios2/ndomod.cfg
…
```

### Installation

To install the script you have to copy files compiled and the configuration files in the right directories.

```bash
mkdir /usr/lib/nagios/ndoutils/
cp ./src/ndomod-2x.o /usr/lib/nagios/ndoutils/ndomod.o
cp ./src/ndo2db-2x /usr/lib/nagios/ndoutils/ndo2db
cp. /src/file2sock /usr/lib/nagios/ndoutils/file2sock
cp. /src/log2ndo /usr/lib/nagios/ndoutils/log2ndo
```

### Getting Things Running

Start the NDO2DB daemon:

```bash
/usr/lib/nagios/ndoutils/ndo2db -c /etc/nagios2/ndo2db.cfg
```

Start or restart nagios if is already running:

```bash
/etc/init.d/nagios2 restart
```

Check the Nagios logs to make sure it started okay:

```bash
vi /var/log/nagios2
```

You should see some entries in the database.
If you want that the NDO2DB start at the system, you have to create a script.

## NagVis

### Prepair System

Make sure your System meets the requirements:

- Nagios v2.X is running and authentication is configured (.htaccess file placed and configured properly!)
- NDOUtils v1.X is installed and worked ([https://nagios.sourceforge.net/docs/ndoutils/NDOUtils.pdf](https://nagios.sourceforge.net/docs/ndoutils/NDOUtils.pdf))
- Working PHP installation (Version > 4.2.0) with "php-gd" support
- php-mysql support is needed for the strongly recommended "ndomy" backend

### Download NagVis

You can get the latest version of Nagvis on the Nagvis' website. [https://www.nagvis.org/downloads](https://www.nagvis.org/downloads)

### Unpack NagVis

Unpack NagVis in a path of your choice

```bash
tar xvzf nagvis-1.2.x.tar.gz
```

You should have the following files:

- etc: This directory contain configuration file and demo Map.
- nagvis: This directory contain all the templates, images and functions of Nagvis.
- wui: This directory contain all Ajax functionalities (to configure map)
- Config.php: This file is used to configure maps.
- Index.php: This will be the index of Nagvis.

### Move NagVis

Place the NagVis directory tree (with etc, nagvis, wui, config.php and index.php inside) into your Nagios share folder.

```bash
mv nagvis /usr/share/nagios2/htdocs/
```

### Configure

Move to new NagVis directory

```bash
cd /usr/share/nagios2/htdocs/nagvis
```

An example for the configuration file can be found in etc/nagvis.ini.php-sample. Copy this example to nagvis/etc/nagvis.ini.php.

```bash
cp etc/nagvis.ini.php-sample etc/nagvis.ini.php
```

Now you can edit this file with your favorite text editor:

```bash
vi etc/nagvis.ini.php
```

The most lines in the fresh copied config.ini.php are commented out. If you want to set different settings, than there are set, uncomment the line and change the value of it.
Edit the following lines:

```bash
[paths]
; absolute physical NagVis path
base="/usr/share/nagios2/htdocs/nagvis/"
; absolute html NagVis path
htmlbase="/nagios2/nagvis"
; absolute html NagVis cgi path
htmlcgi="/nagios2/cgi-bin"
…
[backend_ndomy_1]
…
; hostname for NDO-db
dbhost="localhost"	<- The hostname of your database host
; portname for NDO-db
dbport=3306	<- The port of listened by your database
; database-name for NDO-db
dbname="nagios"	<- The name of your database
; username for NDO-db
dbuser="nagios"	<- The login
; password for NDO-db
dbpass="nagios"	<- The password
; prefix for tables in NDO-db
dbprefix="nagios_"	<- The prefix of your tables
…
; path to the cgi-bin of this backend
htmlcgi="/nagios2/cgi-bin"
```

### Permissions

First check which user the webserver is running with (In my case it is www-data). If you don't know which user the webserver is running on have a look at the webservers configuration. In case of apache you can do this by the following command:

```bash
grep -e '^User' /etc/apache2/httpd.conf
```

If your configuration file is located in another path you should correct this in the command above. The set the permissions to your NagVis directory (In my case the path are like this):

```bash
chown www-data. /usr/share/nagios2/htdocs/nagvis -R
chmod 664 /usr/share/nagios2/htdocs/nagvis/etc/nagvis.ini.php
chmod 775 /usr/share/nagios2/htdocs/nagvis/nagvis/images/maps
chmod 664 /usr/share/nagios2/htdocs/nagvis/nagvis/images/maps/*
chmod 775 /usr/share/nagios2/htdocs/nagvis/etc/maps
chmod 664 /usr/share/nagios2/htdocs/nagvis/etc/maps/*
```

It's possible to set lower permissions on the files, in my setup it's ok like this. Only change them if you know, what you are doing.

### The WUI

NagVis has an included web based config tool called WUI. If you want to use it use your browser to open the page:

```
http://<nagiosserver>/nagvis/config.php
```

Hint: If you have some script or popup blockers, disable them for the WUI.  
When you see the NagVis image, right click on it, a context menu should open, now you can configure NagVis and create maps with the WUI.

The Config Tools DOES NOT display the current Nagios States of Objects configured. Its only for configuring! To "use" your configured Maps afterwards, see STEP 7!

If this does't work for you, or if you don't want to use the WUI, you can simply edit the map config files in the nagvis/etc/maps/ directory with your favorite text editor. For valid format and values have a look at Map Config Format Description on NagVis.org ([https://www.nagvis.org/docs/1.2/map_config_format_description](https://www.nagvis.org/docs/1.2/map_config_format_description)).

### Watch the Maps

You should now be able to watch your defined maps in your browser:

```
http://<nagiosserver>/nagvis/index.php?map=<mapname>
```

## Resources
- [https://www.nagvis.org/](https://www.nagvis.org/)
