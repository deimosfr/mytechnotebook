---
weight: 999
url: "/Rsyslog_\\:\_Installation_et_configuration_d'Rsyslog/"
title: "Rsyslog: Installation and Configuration"
description: "Guide to installing and configuring Rsyslog for log management, including server setup, templates, and client configuration"
categories: ["Debian", "Security", "Database"]
date: "2011-10-06T11:18:00+02:00"
lastmod: "2011-10-06T11:18:00+02:00"
tags:
  [
    "rsyslog",
    "syslog",
    "logging",
    "monitoring",
    "servers",
    "network",
    "security",
  ]
toc: true
---

## Introduction

Syslog is a standard for forwarding log messages in an IP network. The term "syslog" is often used for both the actual syslog protocol, as well as the application or library sending syslog messages.

Syslog is a client/server protocol: the syslog sender sends a small (less than 1KB) textual message to the syslog receiver. The receiver is commonly called "syslogd", "syslog daemon" or "syslog server". Syslog messages can be sent via UDP and/or TCP.

Syslog is typically used for computer system management and security auditing. While it has a number of shortcomings, syslog is supported by a wide variety of devices and receivers across multiple platforms. Because of this, syslog can be used to integrate log data from many different types of systems into a central repository.

## Installation

There is nothing especially to do to install rsyslog as it is embedded in Debian.

## Configuration

All configuration files are located into `/etc/rsyslog.d/`.

By default, all remote incoming logs are stored with the local logs of the rsyslog server (`/var/log/syslog`). If you want to redirect the incoming traffic in a MySQL database, you should update the configuration file with those lines:

(`/etc/rsyslog.conf`):

```bash
#  /etc/rsyslog.conf    Configuration file for rsyslog v3.
#
#           For more information see
#           /usr/share/doc/rsyslog-doc/html/rsyslog_conf.html


#################
#### MODULES ####
#################

$ModLoad imuxsock # provides support for local system logging
$ModLoad imklog   # provides kernel logging support (previously done by rklogd)
#$ModLoad immark  # provides --MARK-- message capability

# provides UDP syslog reception
$ModLoad imudp
$UDPServerRun 514

# provides TCP syslog reception
$ModLoad imtcp
$InputTCPServerRun 514

###########################
#### GLOBAL DIRECTIVES ####
###########################

#
# Use traditional timestamp format.
# To enable high precision timestamps, comment out the following line.
#
$ActionFileDefaultTemplate RSYSLOG_TraditionalFileFormat

#
# Set the default permissions for all log files.
#
$FileOwner root
$FileGroup adm
$FileCreateMode 0640
$DirCreateMode 0755

# Load MySQL module
$ModLoad ommysql
#
# Include all config files in /etc/rsyslog.d/
#
$IncludeConfig /etc/rsyslog.d/*.conf


###############
#### RULES ####
###############

:hostname, !isequal, "syslog1"   ~
& ~

#
# First some standard log files.  Log by facility.
#
auth,authpriv.*         /var/log/auth.log
# Additional things like Cisco syslog
*.*;auth,authpriv.none,local4.none,local7.none,local6.none,local5.none  -/var/log/syslog
#cron.*             /var/log/cron.log
daemon.*            -/var/log/daemon.log
kern.*              -/var/log/kern.log
lpr.*               -/var/log/lpr.log
mail.*              -/var/log/mail.log
user.*              -/var/log/user.log

############################################
#
#     Data Stored in the Mysql Database
#
# please refer to /etc/rsyslog.d/mysql.conf
#
############################################

#
# Logging for the mail system.  Split it up so that
# it is easy to write scripts to parse these files.
#
mail.info           -/var/log/mail.info
mail.warn           -/var/log/mail.warn
mail.err            /var/log/mail.err

#
# Logging for INN news system.
#
news.crit           /var/log/news/news.crit
news.err            /var/log/news/news.err
news.notice         -/var/log/news/news.notice

#
# Some "catch-all" log files.
#
*.=debug;\
    auth,authpriv.none;\
#    local4.none,local7.none,local6.none,local5.none;\
    news.none;mail.none -/var/log/debug
*.=info;*.=notice;*.=warn;\
    auth,authpriv.none;\
    cron,daemon.none;\
 #   local4.none,local7.none,local6.none,local5.none;\
    mail,news.none      -/var/log/messages

#
# Emergencies are sent to everybody logged in.
#
*.emerg             *

#
# I like to have messages displayed on the console, but only on a virtual
# console I usually leave idle.
#
#daemon,mail.*;\
#   news.=crit;news.=err;news.=notice;\
#   *.=debug;*.=info;\
#   *.=notice;*.=warn   /dev/tty8

# The named pipe /dev/xconsole is for the `xconsole' utility.  To use it,
# you must invoke `xconsole' with the `-file' option:
#
#    $ xconsole -file /dev/xconsole [...]
#
# NOTE: adjust the list below, or you'll go crazy if you have a reasonably
#      busy site..
#
daemon.*;mail.*;\
    news.err;\
    *.=debug;*.=info;\
    *.=notice;*.=warn
```

The '-' char in front of the log files name permit to bufferise before writing into file. Use it when a lot of informations has to be written in the logs!

### Define a template

A template is use to determine some action to take by the server syslog if this template is use.
For example, the following template sendthe info into a mysql database name "my_server". In the table "Message" will be send the syslog value propertie %msg% with a certain regex (we want to customize the message) and into the table "Facility" the syslog value propertie %syslogfacility%, etc...

(`/etc/rsyslog.d/servers_template`):

```bash
$template my_server,"insert into `my_server` (Message, Facility, FromHost, Priority, DeviceReportedTime, ReceivedAt, InfoUnitID, SysLogTag) values ('%msg:R,ERE,1,FIELD:%\[.*\](.*)--end%', %syslogfacility%, '%fromhost%', '%syslogpriority%', '%timereported:::date-mysql%', '%timegenerated:::date-mysql%', %iut%, '%programname%')",stdsql
```

Explanations:

- $template my_server --> we create a template named "my_server"
- insert into `my_server` --> this template send to the "my_server" table
- stdsql --> for "standard sql"

Rsyslog is very flexible, we can make a lot a thing with the logs.This template can be applied to some remote incoming logs, wtih a lot of condition (if..then). For example:

```bash
if $fromhost == '192.168.59.17'  then :ommysql:127.0.0.1,syslog-sys,rsyslog,password;my_server
local7.*  :ommysql:127.0.0.1,syslog-sys,rsyslog,password;my_server
if $syslogfacility-text == 'local4' then :ommysql:127.0.0.1,syslog-sys,rsyslog,password;my_server
```

It's possible to use "and" or "and not"

```bash
if $syslogfacility-text == 'local4' and $fromhost == '192.168.79.17' and not ($syslogtag contains '%ASA-4-106023:' and $msg contains 'src inside:10.101') then :ommysql:127.0.0.1,syslog-sys,rsyslog,password;my_server
```

Explanations:

- :ommysql --> send to a mysql server
- 127.0.0.1 --> ip of the mysql server
- syslog-sys --> name of the database
- rsyslog --> username
- password --> password
- my_server --> the template to use for these conditions.

"$syslogfacility-text" or "$fromhost" are called "property replacer", here is the list of the possible replacer: http://www.rsyslog.com/doc/property_replacer.html

It's possible to use some filters, with this nomeclature:

```bash
:property, [!]compare-operation, "value"  destination
```

For example, this line send all the logs which are not send from the "syslog" server to /var/logs/remotelogs

```bash
:hostname, !isequal, "syslog" /var/logs/remotelogs
&~
```

the tilde character is mandatory for applying the condition. This line will drop all remote logs (so they won't be write on a local file), but if the logs match a condition with a mysql database, they will be send to the database:

```bash
:hostname, !isequal, "syslog"   ~
& ~
```

This line send all the logs from 192.0.2.\* to a file:

```bash
if $fromhost-ip startswith '192.0.2.' then /var/log/network2.log
```

The possibles "compare-operation" are:

{{< table "table-hover table-striped" >}}
| Compare-operation | Effect |
|------------------|--------|
| contains | Checks if the string provided in value is contained in the property. There must be an exact match, wildcards are not supported. |
| isequal | Compares the "value" string provided and the property contents. These two values must be exactly equal to match. The difference to contains is that contains searches for the value anywhere inside the property value, whereas all characters must be identical for isequal. As such, isequal is most useful for fields like syslogtag or FROMHOST, where you probably know the exact contents. |
| startswith | Checks if the value is found exactly at the beginning of the property value. For example, if you search for "val" with:msg, startswith, "val" it will be a match if msg contains "values are in this message" but it won't match if the msg contains "There are values in this message" (in the later case, contains would match). Please note that "startswith" is by far faster than regular expressions. So even once they are implemented, it can make very much sense (performance-wise) to use "startswith". |
| regex | Compares the property against the provided POSIX BRE regular expression. |
| ereregex | Compares the property against the provided POSIX ERE regular expression. |
{{< /table >}}

### Client Configuration

On unix computers, the standard syslog daemon is generally used. To send the logs on a remote server, you have to add the following line for udp connections:

(`/etc/rsyslog.conf`):

```bash
*.*          @remote-syslog-server
```

or for tcp connections:

(`/etc/rsyslog.conf`):

```bash
*.*          @@remote-syslog-server
```

where remote-syslog-server is the name or address of your syslog servers.
To test, you can use logger command to valid everythings is OK:

```bash
logger test
```
