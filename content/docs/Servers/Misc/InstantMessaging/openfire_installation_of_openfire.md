---
weight: 999
url: "/OpenFire_\\:_Installation_d'OpenFire/"
title: "OpenFire: Installation of OpenFire"
description: "A guide to install and configure OpenFire, the most functional Jabber server, with integration capabilities for LDAP, Oracle, and MySQL."
categories: ["MySQL", "Database", "Linux", "Servers", "Network"]
date: "2007-11-12T15:03:00+02:00"
lastmod: "2007-11-12T15:03:00+02:00"
tags: ["OpenFire", "Jabber", "Java", "MySQL", "Configuration", "Server", "Communication"]
toc: true
---

## Introduction

OpenFire is currently the most functional Jabber server available. It offers very interesting features, though some are paid.

It's still very integratable with technologies such as LDAP, Oracle, or MySQL.

## Download

To install it, it's really not difficult. We download it from [the official website](https://www.igniterealtime.org/projects/openfire):

```bash
wget 'http://www.igniterealtime.org/downloadServlet?filename=openfire/openfire_3_3_0.tar.gz'
```

Then we need to download Java [Linux (self-extracting file)](https://www.java.com/en/download/manual.jsp):

```bash
wget 'http://javadl.sun.com/webapps/download/AutoDL?BundleId=11187'
```

## Installation

Next we extract Java:

```bash
sh jre-6u1-linux-i586.bin
```

We'll do the same with OpenFire:

```bash
tar -xzvf openfire_3_3_0.tar.gz
```

Put everything in `/usr/share/openfire` for example and install it all:

```bash
mkdir /usr/share/openfire
mv jre1.6.0_01 openfire /usr/share/openfire
```

## Configuration

### MySQL

We need to create a database where data will be stored, then create the necessary components:

```bash
mysql -uroot -p
create database openfire;
quit;
mysql -uroot -p openfire < /usr/share/openfire/openfire/resources/database/openfire_mysql.sql
```

I recommend doing at least minimal security measures. That means at least a dedicated user and restricted rights to the database. Using everything as root is far from secure!

### Environment

All that's left is to add this to the `~/.bashrc`, `~/.zshrc` or other configuration file for the user who will run OpenFire:

```bash
export INSTALL4J_JAVA_HOME=/usr/share/openfire/jre
```

Then, restart the session and launch it:

```bash
cd /usr/share/openfire/openfire/bin && openfire start &
```

## Automatic Startup at Boot

To launch it at boot, it's easy as usual, just add the line above to `/etc/rc.local`.

Alternatively, in the extras directory (`/usr/share/openfire/openfire/bin/extras`), you'll find a ready-made script to start it as a service.

## FAQ

### Managing Memory Usage

If you're using OpenFire for personal needs, there's too much RAM allocated by default. If you need more, this FAQ is also good for you ;-).

For those with limited Java knowledge, I'll explain a bit how it works. We have 2 limits:

* Xms: RAM will be directly allocated for the program, even if it uses less
* Xmx: Maximum memory the program can use before using the [garbage collector](https://en.wikipedia.org/wiki/Garbage_collector) (attempting to recover unused RAM).

Edit the `bin/openfire` file in your OpenFire folder and find the variable **INSTALL4J_ADD_VM_PARAM**. Then adjust to your needs:

```bash
INSTALL4J_ADD_VM_PARAMS="-Xms16m -Xmx32m"
```

## Resources
- [Documentation on OpenFire and Spark installation](/pdf/openfire_3_3_3_spark_2_5_7.pdf)
