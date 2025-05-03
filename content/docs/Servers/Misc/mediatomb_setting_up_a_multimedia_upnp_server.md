---
weight: 999
url: "/MediaTomb_\\:_Mise_en_place_d'un_serveur_multimédia_(UPnP)/"
title: "MediaTomb: Setting up a multimedia (UPnP) server"
description: "A guide on how to set up MediaTomb, a multimedia server for streaming movies, music, and photos using UPnP technology."
categories: ["MySQL", "Debian", "Database", "Multimedia"]
date: "2009-04-23T19:19:00+02:00"
lastmod: "2009-04-23T19:19:00+02:00"
tags: ["UPnP", "MediaTomb", "MySQL", "Debian", "PS3", "Streaming", "Multimedia"]
toc: true
---

## Introduction

MediaTomb is a multimedia server capable of streaming movies, music, and photos. It's very practical for game consoles like Xbox360 or PlayStation 3, but has many other uses as well.

For example, if you have a machine that contains movies but isn't connected to your TV screen, you need something to relay the content. This is where a PlayStation 3, for instance, comes in. However, it needs to know where to find the movies, so you need to install a UPnP server on your machine.

One of the best UPnP servers today is [MediaTomb](https://mediatomb.cc/).

## Installation

As usual, for Debian users:

```bash
apt-get install mediatomb
```

## Configuration

### MySQL

This step is optional, but personally I prefer centralizing information on my machine. If you skip this step, a MySQLite database will be used. I decided to set up MediaTomb with a MySQL database. First, connect to your MySQL server:

```bash
mysql [-u <username>] [-p]
```

Then let's create a database dedicated to MediaTomb and an associated user:

```bash
mysql> CREATE DATABASE db_mediatomb;
mysql> GRANT ALL ON db_mediatomb.* TO 'mediatomb_user'@'localhost' IDENTIFIED BY 'mon_pasowrd';
```

### Mediatomb

#### MySQL

Now we're ready to configure this in the `/etc/mediatomb/config.xml` configuration file. First, modify this to disable SQLite:

```xml
<sqlite3 enabled="no">
```

Now, activate MySQL:

```xml
<mysql enabled="yes">
   <host>localhost</host>
   <username>mediatomb_user</username>
   <password>mon_password</password>
   <database>db_mediatomb</database>
</mysql>
```

#### Thumbnails

First, let's install this small package:

```bash
apt-get install ffmpegthumbnailer
```

Then add this code in the '<server>' section:

```xml
<extended-runtime-options>
<ffmpegthumbnailer enabled="yes">
  <thumbnail-size>128</thumbnail-size>
  <seek-percentage>5</seek-percentage>
  <filmstrip-overlay>yes</filmstrip-overlay>
  <workaround-bugs>no</workaround-bugs>
</ffmpegthumbnailer>
</extended-runtime-options>
```

Now restart the server:

```bash
/etc/init.d/mediatomb restart
```

Great! We're done with MySQL.

### Graphical Interface

Now you can use the graphical interface. Connect to the server with your MediaTomb server address on port 49152:

```
http://ip_address:49152
```

All you have to do is configure your folders to scan and their refresh times. It takes some time to index them (depending on the number of items in your folders).

### Security

To limit access to network interfaces on your multimedia server, edit the `/etc/default/mediatomb` file:

```
INTERFACE="bond0"
```

Then restart the MediaTomb service.

## Working with PlayStation 3

For the server to be recognized by PlayStation 3, you need to modify the file:

```xml
<protocolInfo extend="yes"/><!-- For PS3 support change to "yes" -->
<map from="avi" to="video/divx"/><!-- Uncomment the line below for PS3 divx support -->
```

**Important: You must have a PS3 firmware version 1.80 or higher**

Then restart the MediaTomb server. Now, on the PlayStation, you can scan for Multimedia servers and yours will appear magically.

## Resources
- [Fuppes: Set Up A Linux PlayStation 3 Media Server](/pdf/set_up_a_linux_playstation_3_media_server.pdf)
