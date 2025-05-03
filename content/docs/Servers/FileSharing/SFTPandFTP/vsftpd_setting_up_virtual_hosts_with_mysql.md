---
weight: 999
url: "/Vsftpd_\\:_Mise_en_place_d'hôtes_virtuels_avec_MySQL/"
title: "Vsftpd: Setting up virtual hosts with MySQL"
description: "This guide explains how to configure vsftpd with virtual hosts using MySQL, including SSL configuration and security settings"
categories: ["Security", "Linux", "MySQL"]
date: "2010-01-21T13:51:00+02:00"
lastmod: "2010-01-21T13:51:00+02:00"
tags: ["Servers", "FTP", "Security", "SSL", "vsftpd", "MySQL", "OpenBSD"]
toc: true
---

## Introduction

[vsFTPd](https://fr.wikipedia.org/wiki/Vsftpd), short for Very Secure FTP Daemon, is a free, simple, and secure FTP server.

It was developed with the best possible security in mind to address vulnerabilities in classic FTP servers.

It includes all the standard options of classic FTP servers (ProFTPd, Pure-FTPd, etc.). It supports IPv6 and SSL.

VsFTPd is an FTP server designed with maximum security in mind. Unlike other FTP servers (ProFTPd, PureFTPd, etc.), no security vulnerability has ever been found in VsFTPd. This server is widely used by large enterprises.

The default configuration of VsFTPd is very restrictive:

1. Only anonymous accounts are allowed to connect to the server, and in read-only mode
2. Users can only access their own account

VsFTPd features:

- Accessible configuration
- Virtual users
- Virtual IP addresses
- Bandwidth limitation
- IPv6
- Built-in SSL encryption support

It is distributed under the terms of the GNU GPL license.

## Installation

I'm performing this installation on OpenBSD, so here's how to proceed:

```bash
pkg_add -iv vsftpd
```

Yes, that's all :-)

## Configuration

### ssl

We'll need a certificate for SSL, so let's generate one:

```bash
mkdir -p /etc/ssl/vsftpd
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout /etc/ssl/vsftpd/vsftpd.pem -out /etc/ssl/vsftpd/vsftpd.pem
```

And then we set the proper permissions:

```bash
chown root:root /etc/ssl/vsftpd/vsftpd.pem
chmod 600 /etc/ssl/vsftpd/vsftpd.pem
```

### vsftpd

For the configuration, I want all users to be chrooted and SSL to be used at all levels (connections + data). Here's my config (`/etc/vsftpd.conf`):

```bash
# Example config file /etc/vsftpd.conf
#
# The default compiled in settings are fairly paranoid. This sample file
# loosens things up a bit, to make the ftp daemon more usable.
# Please see vsftpd.conf.5 for all compiled in defaults.
#
# READ THIS: This example file is NOT an exhaustive list of vsftpd options.
# Please read the vsftpd.conf.5 manual page to get a full idea of vsftpd's
# capabilities.
#
# Standalone mode
listen=YES
#
# SSL
ssl_enable=YES
#allow_anon_ssl=NO
force_local_data_ssl=NO
force_local_logins_ssl=YES
 
ssl_tlsv1=YES
ssl_sslv2=YES
ssl_sslv3=YES
 
rsa_cert_file=/etc/ssl/vsftpd/vsftpd.pem
rsa_private_key_file=/etc/ssl/vsftpd/vsftpd.pem
 
#
# TCP Wrappers
#tcp_wrappers=YES
#
# Allow anonymous FTP? (Beware - allowed by default if you comment this out).
anonymous_enable=NO
#
# Uncomment this to allow local users to log in.
local_enable=YES
#
# Uncomment this to enable any form of FTP write command.
write_enable=YES
#
# Default umask for local users is 077. You may wish to change this to 022,
# if your users expect that (022 is used by most other ftpd's)
local_umask=022
#
# Uncomment this to allow the anonymous FTP user to upload files. This only
# has an effect if the above global write enable is activated. Also, you will
# obviously need to create a directory writable by the FTP user.
#anon_upload_enable=YES
#
# Uncomment this if you want the anonymous FTP user to be able to create
# new directories.
#anon_mkdir_write_enable=YES
#
# Activate directory messages - messages given to remote users when they
# go into a certain directory.
dirmessage_enable=YES
#
# Activate logging of uploads/downloads.
xferlog_enable=YES
#
# Make sure PORT transfer connections originate from port 20 (ftp-data).
connect_from_port_20=YES
#
# If you want, you can arrange for uploaded anonymous files to be owned by
# a different user. Note! Using "root" for uploaded files is not
# recommended!
#chown_uploads=YES
#chown_username=whoever
#
# You may override where the log file goes if you like. The default is shown
# below.
xferlog_file=/var/log/vsftpd.log
#
# If you want, you can have your log file in standard ftpd xferlog format
#xferlog_std_format=YES
#
# You may change the default value for timing out an idle session.
idle_session_timeout=600
#
# You may change the default value for timing out a data connection.
#data_connection_timeout=120
#
# It is recommended that you define on your system a unique user which the
# ftp server can use as a totally isolated and unprivileged user.
nopriv_user=_vsftpd
#
# Enable this and the server will recognise asynchronous ABOR requests. Not
# recommended for security (the code is non-trivial). Not enabling it,
# however, may confuse older FTP clients.
#async_abor_enable=YES
#
# By default the server will pretend to allow ASCII mode but in fact ignore
# the request. Turn on the below options to have the server actually do ASCII
# mangling on files when in ASCII mode.
# Beware that on some FTP servers, ASCII support allows a denial of service
# attack (DoS) via the command "SIZE /big/file" in ASCII mode. vsftpd
# predicted this attack and has always been safe, reporting the size of the
# raw file.
# ASCII mangling is a horrible feature of the protocol.
#ascii_upload_enable=YES
#ascii_download_enable=YES
#
# You may fully customise the login banner string:
ftpd_banner=Deimos FTP Server
#
# You may specify a file of disallowed anonymous e-mail addresses. Apparently
# useful for combatting certain DoS attacks.
#deny_email_enable=YES
# (default follows)
#banned_email_file=/etc/vsftpd.banned_emails
#
# You may specify an explicit list of local users to chroot() to their home
# directory. If chroot_local_user is YES, then this list becomes a list of
# users to NOT chroot().
chroot_list_enable=YES
chroot_local_user=YES
# (default follows)
chroot_list_file=/etc/ftpchroot
#
# You may activate the "-R" option to the builtin ls. This is disabled by
# default to avoid remote users being able to cause excessive I/O on large
# sites. However, some broken FTP clients such as "ncftp" and "mirror" assume
# the presence of the "-R" option, so there is a strong case for enabling it.
#ls_recurse_enable=YES
 
#
# If enabled, vsftpd will load a list of usernames from the filename
# given by userlist_file. If a user tries to log in using a name in this
# file, they will be denied before they are asked for a password.
# This may be useful in preventing clear text passwords being transmitted.
userlist_enable=YES
#
# This option is the name of the file loaded when the userlist_enable
# option is active.
userlist_file=/etc/ftpusers
#
# This option should be the name of a directory which is empty. Also,
# the directory should not be writable by the ftp user. This directory
# is used as a secure chroot() jail at times vsftpd does not require
# filesystem access.
secure_chroot_dir=/var/empty
 
pasv_enable=YES
#
# The minimum port to allocate for PASV style data connections.
# Can be used to specify a narrow port range to assist firewalling.
pasv_min_port=49152
#
# The maximum port to allocate for PASV style data connections.
# Can be used to specify a narrow port range to assist firewalling.
pasv_max_port=65535
#
# By default, numeric IDs are shown in the user and group fields of
# directory listings. You can get textual names by enabling this parameter.
# It is off by default for performance reasons.
text_userdb_names=YES
```

## Launching

There are several ways to launch it. To test, you can simply run the vsftpd command:

```bash
vsftpd &
```

Once everything is working as you want, it's best to couple it with inetd or xinetd.

## Resources
Here are some documents covering VSFTP in various forms:
- [Virtual Hosting with vsftpd and MySQL](/pdf/virtual_hosting_vsftpd_mysql.pdf)
- [Four highly secure FTP servers with vsftpd](/pdf/quatre_serveurs_ftp_hyper_sécurisés_avec_vsftpd.pdf)
