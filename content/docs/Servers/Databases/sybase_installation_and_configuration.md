---
weight: 999
url: "/Sybase_\\:_installation_et_configuration/"
title: "Sybase: Installation and Configuration"
description: "A guide to installing and configuring Sybase Adaptive Server Enterprise database on Linux systems."
categories: ["Database", "Linux", "Servers"]
date: "2008-04-06T08:41:00+02:00"
lastmod: "2008-04-06T08:41:00+02:00"
tags: ["Sybase", "Database", "SQL", "Linux", "Installation"]
toc: true
---

## 1 Introduction

Adaptive Server Enterprise is the data management system designed to manage the explosion of data volume in critical contexts.

Adaptive Server Enterprise (ASE) has long been recognized for its reliability, low total cost of ownership and optimal performance. The latest version ASE 15 focuses on key features that form the foundation for long-term strategic agility and continuous innovation in critical environments. ASE 15 offers unique security options and numerous new features aimed at improving performance while reducing costs and operational risks. Discover how to take advantage of new technologies such as grids and clusters, service-oriented architectures, and real-time messaging.

## 2 Installation

First, create a "sybase" user on our system and give it all rights on the partition:

```bash
useradd sybase
chown -R sybase. /
```

Next, we need to know that Linux only shares 32mb by default with Sybase, which needs at least 64mb, so we need to use sysctl to increase the shared memory (to be added to rc.local):

```bash
# rc.local
/sbin/sysctl -w kernel.shmmax=67108864
/sbin/sysctl kernel.shmmax (to check)
```

We may need some additional libraries during installation, which can be found using aptitude:

```bash
apt-get install libaio-dev libstdc++5-3.3-dev
```

First, download the latest sources: ASE-x.x.x.gz
Extract the archive, enter the "sybase" directory and type:

```bash
./setup -console
```

After answering all the installation questions, the installation completes correctly (we hope). However, you need to accept configuring all services (if you did a complete installation) to display all the information.
During installation, all server information is displayed, save it somewhere:

```
 The installer will now configure the new servers with
 the following values. Click Next to continue configuring the
 servers.

   Adaptive Server

      Adaptive Server Name                                      DEBSYBASE5
      Port number                                               5012
      Page size                                                 2k
      Error log                                                 /opt/sybase/ASE-15_0/install/DEBSYBASE5.log
      Master device                                             /opt/sybase/data/master.dat
      Master device size (MB)                                   30
      Master database size (MB)                                 13
      System procedure device                                   /opt/sybase/data/sysprocs.dat
      System procedure device size (MB)                         132
      System procedure database size (MB)                       132
      System Device                                             /opt/sybase/data/sybsysdb.dat
      System Device Size (MB)                                   1
      System Database Size (MB)                                 1

   Backup Server

      Backup Server Name       DEBSYBASE_BS3
      Port number              5013
      Error log                /opt/sybase/ASE-15_0/install/DEBSYBASE_BS3.log

   Monitor Server

      Monitor Server Name      DEBSYBASE_MS2
      Port number              5014
      Error log                /opt/sybase/ASE-15_0/install/DEBSYBASE_MS2.log

   XP Server

      XP Server Name        DEBSYBASE_XP2
      Port number           5015
      Error log             /opt/sybase/ASE-15_0/install/DEBSYBASE_XP2.log

   Job Scheduler

      Job Scheduler Agent Name        DEBSYBASE5_JSAGENT
      Port number                     4902
      Management Device               /opt/sybase/data/sybmgmtdb.dat
      Management Device Size (MB)     75
      Management Database Size (MB)   75

   Self Management

      Self Management User Name       sa
      Self Management User Password   ******

   Web Services

      HTTP port number for production task     8181
      HTTPS port number for production task    8182
      Hostname for production task             deb-sybase
      Certificate password                     ******
      Keystore password                        ******
      Log file for production task             /opt/sybase/WS-15_0/logs/producer.log
      Port number for consumption task         8183
      Log file for consumption task            /opt/sybase/WS-15_0/logs/consumer.log

   Unified Agent - Self Discovery Service Adaptor

      Adaptor   UDP

   Unified Agent - Security Login Modules

      CSI.loginModule.1.provider             com.sybase.ua.services.security.simple.SimpleLoginModule
      CSI.loginModule.1.controlFlag          sufficient
      CSI.loginModule.1.options.moduleName   Simple Login Module
      CSI.loginModule.1.options.username     uafadmin
      CSI.loginModule.1.options.password     ******
      CSI.loginModule.1.options.roles        uaAgentAdmin,uaPluginAdmin
      CSI.loginModule.1.options.encrypted    false

      CSI.loginModule.2.provider             com.sybase.ua.services.security.ase.ASELoginModule
      CSI.loginModule.2.controlFlag          sufficient
      CSI.loginModule.2.options.moduleName   ASE Login Module
```

### 2.1 Environment Variables

Here are the SYBASE variables that need to be exported:

```bash
export JAVA_JRE=/opt/sybase/_jvm/
export SYBASE_OCS=OCS-15_0
export SYBASE=/opt/sybase
export JAVA_HOME=/opt/sybase/_jvm/
export SYBASE_JRE=/opt/sybase/shared/jre142_013/
export SYBASE_WS=WS-15_0
export LANG=fr
export PATH=$PATH:/etc/init.d/
export SYBASE_ASE=ASE-15_0/
export PATH=$PATH:/opt/sybase/ASE-15_0/install/
export PATH=$PATH:/opt/sybase/OCS-15_0/bin/
```

## 3 Administration

### 3.1 Starting the Server

If you have declared all the variables above:

```bash
startserver -f RUN_DEBSYBASE (according to our server name)
```

Otherwise, navigate to the Sybase directory and run the startup script:

```bash
cd /opt/sybase/ASE-15.0/install
./startserver -f RUN_DEBSYBASE
```

### 3.2 Stopping the Server

To properly shut it down, use isql:

```bash
isql -Usa -P(sa password) -S(sybase server name) (or navigate to the directory where isql is located)
1>shutdown
2>go
```

### 3.3 Server Status

If you have declared all the variables above, you can check if the server is running (similar to ps -aux | grep SYBASE), otherwise find where the showserver file is located:

```bash
showserver
```

### 3.4 Creating a User, Device, and Database

Everything is done in isql:

```bash
1>sp_addlogin 'username', 'password'
2>go
```

To create a device, which is a file that will contain a database:

```bash
-- file in /data/sybase/data01.dat of 100mb
1>disk init name='data01', physname='/data/sybase/data01.dat', size='100m'
2>go
```

To create a database on this file:

```bash
CREATE DATABASE <name> ON <device> = <size in MB>
```

### 3.5 Changing a Password

In isql using the account whose password is to be changed:

```bash
1> exec sp_password NULL, "Secr3t"   (old_password, new_password)(here we set a password in place of an non-existent password)
2> go
Password correctly set.
(return status = 0)
```

## 4 FAQ

ERROR:

```
00:00000:00000:2008/03/06 16:02:18.10 kernel  kbcreate: couldn't create kernel region.
00:00000:00000:2008/03/06 16:02:18.10 kernel  kistartup: could not create shared memory
```

This means the shared memory is not large enough.

SOLUTION:

```bash
/sbin/sysctl -w kernel.shmmax=67108864
```

## 5 Links

Sybase website: http://www.sybase.fr/
All documentation: http://infocenter.sybase.com/help/index.jsp?topic=/com.sybase.help.ase_15.0.sag1/html/sag1/sag11.htm
A tutorial to manage your server with the Windows client: http://www.ianywhere.com/developer/product_manuals/sqlanywhere/0901/fr/html/dbfgfr9/00000165.htm
