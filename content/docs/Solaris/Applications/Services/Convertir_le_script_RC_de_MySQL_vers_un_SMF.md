---
weight: 999
url: "/Convertir_le_script_RC_de_MySQL_vers_un_SMF/"
title: "Convert MySQL RC Script to SMF"
description: "How to convert MySQL RC startup scripts to Solaris Management Facility (SMF) configuration"
categories: ["Database", "MySQL", "Solaris"]
date: "2008-04-20T08:10:00+02:00"
lastmod: "2008-04-20T08:10:00+02:00"
tags: ["MySQL", "Solaris", "SMF", "Database", "Configuration"]
toc: true
---

## Introduction

Solaris 10 and later now supplies MySQL as part of the OS, provided you've installed the "SUNWmysql[rtu]" pkgs, but it's started via a legacy RC script still. This document details how to create an SMF manifest to start MySQL instead.

Note: This process is only needed if you are running Solaris 10, or you wish to use the MySQL 4.x installation that is supplied with Nevada. Later builds of Nevada (I believe snv_79 and later) now come with MySQL 5.0 and includes a service manifest for this version.

## Find scripts

Find all of the current RC scripts:

```bash
$ find /etc/rc* /etc/init.d | grep -i mysql
/etc/rc0.d/K00mysql
/etc/rc1.d/K00mysql
/etc/rc2.d/K00mysql
/etc/rc3.d/S99mysql
/etc/rcS.d/K00mysql
```

If you've never configured MySQL, you may find these don't exist yet.

Remove all the old legacy scripts, if they exist:

```bash
for x in `find /etc/rc* /etc/init.d
```

## Create SMF Script

Create the MySQL manifest:

```xml
<?xml version="1.0"?>
<!DOCTYPE service_bundle SYSTEM "/usr/share/lib/xml/dtd/service_bundle.dtd.1">
<!--
    Copyright 2005 Sun Microsystems, Inc.  All rights reserved.
    Use is subject to license terms.
    MySQL.xml : MySQL manifest, Scott Fehrman, Systems Engineer
    updated: 2005-09-16
-->

<service_bundle type='manifest' name='MySQL'>
<service name='application/database/mysql' type='service' version='1'>

   <single_instance />

   <dependency
      name='filesystem'
      grouping='require_all'
      restart_on='none'
      type='service'>
      <service_fmri value='svc:/system/filesystem/local' />
   </dependency>

   <exec_method
      type='method'
      name='start'
      exec='/etc/sfw/mysql/mysql.server start'
      timeout_seconds='120' />

   <exec_method
      type='method'
      name='stop'
      exec='/etc/sfw/mysql/mysql.server stop'
      timeout_seconds='120' />

   <instance name='default' enabled='false' />

   <stability value='Unstable' />

   <template>
      <common_name>
         <loctext xml:lang='C'>MySQL RDBMS 4.0.15</loctext>
      </common_name>
      <documentation>
         <manpage title='mysql' section='1' manpath='/usr/sfw/share/man' />
      </documentation>
   </template>

</service>
</service_bundle>
```

And copy it:

```bash
# mkdir /var/svc/manifest/application/database
# cp mysql.txt /var/svc/manifest/application/database/mysql.xml
```

If you've got Postgresql installed already, you'll already have a `/var/svc/manifest/application/database` directory.

Import the manifest:

```bash
$ svccfg validate /var/svc/manifest/application/database/mysql.xml
$ svccfg import /var/svc/manifest/application/database/mysql.xml
```

Check the service:

```bash
$ svcs mysql
STATE STIME FMRI
disabled 13:23:17 svc:/application/database/mysql:default
```

There you have it. Now to enable it, just ensure you've configured MySQL as per the README at `/etc/sfw/mysql/README.solaris.mysql` and then enable the service:

```bash
$ svcadm enable mysql
$ svcs mysql
STATE STIME FMRI
online 13:29:14 svc:/application/database/mysql:default
```
