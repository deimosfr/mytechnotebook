---
weight: 999
url: "/SMF_\\:_Service_Management_Facility/"
title: "SMF: Service Management Facility"
description: "A comprehensive guide to Solaris Service Management Facility (SMF), covering service management, runlevels, milestones, and configuration of services."
categories: ["Linux", "Servers", "System Administration"]
date: "2010-02-03T21:43:00+02:00"
lastmod: "2010-02-03T21:43:00+02:00"
tags: ["Solaris", "SMF", "System Administration", "Services", "Boot Process"]
toc: true
---

## Introduction

SMF services are listed by categories:

- application
- device
- legacy
- milestone
- network
- platform
- site
- system

Example:

```
svc:/system/filesystem/root:default
```

- The prefix svc indicates that it's a service managed by SMF
- The category of the service is "system"
- The service itself is a filesystem
- The instance of the service is the root of the file system
- The word "default" identifies the first, in this case only, instance of the service

Another example:

```
lrc:/etc/rc3_d/S90samba
```

- The name "lrc" indicates that the running service is not managed by SMF
- The pathname "/etc/rc3_d" refers to the folder "/etc/rc3.d" where the script is used to be managed
- The script name is S90samba

To list the names and states of services:

```bash
$ svcs
STATE          STIME    FMRI
legacy_run     Feb_10   lrc:/etc/rc2_d/S10lu
legacy_run     Feb_10   lrc:/etc/rc2_d/S20sysetup
legacy_run     Feb_10   lrc:/etc/rc2_d/S90wbem
legacy_run     Feb_10   lrc:/etc/rc2_d/S99dtlogin
legacy_run     Feb_10   lrc:/etc/rc3_d/S81volmgt
(output removed)
online         Feb_10   svc:/system/system-log:default
online         Feb_10   svc:/system/fmd:default
online         Feb_10   svc:/system/console-login:default
online         Feb_10   svc:/network/smtp:sendmail
online         Feb_10   svc:/milestone/multi-user:default
online         Feb_10   svc:/milestone/multi-user-server:default
online         Feb_10   svc:/system/zones:default
offline        Feb_10   svc:/application/print/ipp-listener:default
offline        Feb_10   svc:/application/print/rfc1179:default
maintenance   10:24:15 svc:/network/rpc/spray:default
```

Here's the list of possible states:

{{< table "table-hover table-striped" >}}
| State | Description |
|-------|-------------|
| online | The service instance is enabled and has successfully started. |
| offline | The service instance is enabled, but the service is not yet running or available to run. |
| disabled | The service instance is not enabled and is not running. |
| legacy_run | The legacy service is not managed by SMF, but the service can be observed. This state is only used by legacy services. |
| uninitialized | This state is the initial state for all services before their configuration has been read. |
| maintenance | The service instance has encountered an error that must be resolved by the administrator. |
| degraded | The service instance is enabled, but is running at a limited capacity. |
{{< /table >}}

## Runlevels and Milestones Management

Here are the different service types:

- single-user
- multi-user
- multi-user-server
- network
- name-services
- sysconfig
- devices

Here's the relationship between milestones and services:

![Sun-milstone-1](/images/sun-milstone-1.avif)

Here's an example of the relationship between dependencies:

![Sun-milstone-2](/images/sun-milstone-2.avif)

To determine the current milestones:

```bash
$ svcs
```

Here are the states that a milestone can take:

- none
- single-user
- multi-user
- multi-user-server
- all

To choose which milestone you want to boot into:

```
ok> boot -m milestone=single-user
```

Note: svc.startd is a daemon.

The database listing all these services is located here:

```
/etc/svc/repository.db
```

This database is managed by the svc.configd service.

If there's an error at this level, your machine won't boot. To repair this, boot in single user mode and run this command:

```bash
/lib/svc/bin/restore_repository
```

Here's a list of milestones and runlevels:

{{< table "table-hover table-striped" >}}
| Run Level | Milestone | Description |
|-----------|-----------|-------------|
| 0 | | System is running the PROM monitor. |
| s or S | single-user | Solaris OS single-user mode with critical file systems mounted and accessible. |
| 1 | | The system is running in a single-user administrative state with access to all available file systems. |
| 2 | multi-user | The system is supporting multiuser operations. Multiple users can access the system. All system daemons are running except for the Network File System (NFS) server and some other network resource server related daemons. |
| 3 | multi-user-server | The system is supporting multiuser operations and has NFS resource sharing and other network resource servers available. |
| 4 | | This level is currently not implemented. |
| 5 | | A transitional run level in which the Solaris OS is shut down and the system is powered off. |
| 6 | | A transitional run level in which the Solaris OS is shut down and the system reboots to the default run level. |
{{< /table >}}

To know which runlevel you're in, use this command:

```bash
who -r
```

For runlevels, you can find in `/etc/` or `/sbin` the different runlevels:

- rc0
- rc1
- rc2
- rc3
- rc4
- rc5
- rc6
- rcS

When you look in one of these folders, you can see boot and stop processes. To distinguish them:

- K for Kill
- S for Start

To know the current process with respect to the daemon:

```bash
$ ls -i S90samba
4715 samba
```

## The Boot Process

If you've followed along correctly, you should understand that the boot order looks like this:

- PROM boot phase
- Boot program phase
- Kernel initialization phase
- Init phase
- svc.startd phase

![Sun-initd](/images/sun-initd.avif)

During the boot phase, the kernel reads its configuration file `/etc/system`, then loads modules. It uses the "ufsboot" command to load the files.

Then it loads the `/etc/init` daemon.

Here's what you can find in the `/etc/system` file:

- moddir

Searches for any modules to load

- root device and root file system configuration

By default: rootfs:ufs  
This is for the root file system. Example:

```
rootdev:/sbus@1,f8000000/esp@0,800000/sd@3,0:a
```

- exclude

Won't load the listed modules. Example:

```
exclude: sys/shmsys
```

- forceload

Forces loading certain modules. Example:

```
forceload: drv/vx
```

- set

Changes kernel parameters to modify system operations. Example:

```
set maxusers=40
```

Make a copy of the `/etc/system` file before saving changes. If the file is incorrect, you won't be able to boot. If you have a problem with the modified file, here's the solution to repair:

```bash
$ ok boot -a
Enter filename [kernel/sparcv9/unix]:
Enter default directory for modules [/platform...]:
Name of system file [etc/system]: etc/system.orig - or - /dev/null
root filesystem type [ufs]:
Enter physical name of root device [/...]:
(further boot messages omitted)
```

## Inittab

Each line in this file looks like this:

```
id:rstate:action:process
```

Here are the fields:

{{< table "table-hover table-striped" >}}
| Field | Description |
|-------|-------------|
| id | Two character identifier for the entry |
| rstate | Run levels to which this entry applies |
| action | Defines how the process listed should be run<br>For a description of the action keywords see man inittab |
| process | Defines the command to execute |
{{< /table >}}

By default, here's what you find in the inittab file:

```bash
ap::sysinit:/sbin/autopush -f /etc/iu.ap
sp::sysinit:/sbin/soconfig -f /etc/sock2path
smf::sysinit:/lib/svc/bin/svc.startd	>/dev/msglog 2<>/dev/msglog </dev/console
p3:s1234:powerfail:/usr/sbin/shutdown -y -i5 -g0 >/dev/msglog 2<>/dev/msglog
```

For possible actions, we have:

- sysinit

Executes the process before the init process tries to access the console (for example, the console login prompt). The init process waits for completion of the process before it continues to read the inittab file.

- powerfail

Executes the process only if the init process receives a power fail signal.

The svc.startd daemon is the replacement for init. To view the current configuration:

```bash
/var/svc/manifest
```

If you want to see milestone files to edit them:

```
single-user.xml
multi-user.xml
multi-user-server.xml
network.xml
name-services.xml
sysconfig.xml
/sbin/rc2
/lib/svc/method/fs-local
```

## SVCS

Here's the command to monitor SMF services:

```bash
$ svcs
STATE          STIME    FMRI
legacy_run     13:45:11 lrc:/etc/rcS_d/S29wrsmcfg
legacy_run     13:45:37 lrc:/etc/rc2_d/S10lu
legacy_run     13:45:38 lrc:/etc/rc2_d/S20sysetup
legacy_run     13:45:38 lrc:/etc/rc2_d/S40llc2
legacy_run     13:45:38 lrc:/etc/rc2_d/S42ncakmod
legacy_run     13:45:39 lrc:/etc/rc2_d/S47pppd
(output omitted)
online         13:45:36 svc:/network/smtp:sendmail
online         13:45:38 svc:/network/ssh:default
online         13:45:38 svc:/system/fmd:default
online         13:45:38 svc:/application/print/server:default
online         13:45:39 svc:/application/print/rfc1179:default
online         13:45:41 svc:/application/print/ipp-listener:default
online         13:45:45 svc:/milestone/multi-user:default
online         13:45:53 svc:/milestone/multi-user-server:default
online         13:45:54 svc:/system/zones:default
online          8:46:25 svc:/system/filesystem/local:default
online          8:46:26 svc:/network/inetd:default
online          8:46:32 svc:/network/rpc/meta:tcp
online          8:46:32 svc:/system/mdmonitor:default
online          8:46:38 svc:/milestone/multi-user:default
online         13:14:35 svc:/network/telnet:default
maintenance     8:46:21 svc:/network/rpc/keyserv:default
```

To verify the status of a service:

```bash
$ svcs svc:/system/console-login:default

STATE          STIME    FMRI
online         14:38:27 svc:/system/console-login:default
```

To view the dependencies of a service:

```bash
svcs -d svc:/system/filesystem/local:default

STATE          STIME    FMRI
online         14:38:15 svc:/system/filesystem/minimal:default
online         14:38:26 svc:/milestone/single-user:default
```

This shows what a service needs (dependencies):

```bash
$ svcs -d milestone/multi-user:default
STATE          STIME    FMRI
online         13:44:53 svc:/milestone/name-services:default
online         13:45:12 svc:/milestone/single-user:default
online         13:45:13 svc:/system/filesystem/local:default
online         13:45:15 svc:/network/rpc/bind:default
online         13:45:16 svc:/milestone/sysconfig:default
online         13:45:17 svc:/system/utmp:default
online         13:45:19 svc:/network/inetd:default
online         13:45:31 svc:/network/nfs/client:default
online         13:45:34 svc:/system/system-log:default
online         13:45:36 svc:/network/smtp:sendmail
```

Here we can see other services that depend on /system/filesystem/local:

```bash
$ svcs -D svc:/system/filesystem/local
STATE          STIME    FMRI
disabled       14:38:00 svc:/network/inetd-upgrade:default
disabled       14:38:07 svc:/network/nfs/server:default
online         14:38:30 svc:/network/inetd:default
online         14:38:30 svc:/network/smtp:sendmail
online         14:38:30 svc:/system/cron:default
online         14:38:30 svc:/system/sac:default
online         14:38:45 svc:/system/filesystem/autofs:default
online         14:38:47 svc:/system/dumpadm:default
online         14:38:51 svc:/milestone/multi-user:default
```

### svcadm

This command is used to change the state of a service:

```bash
$ ps -ef
```

```bash
$ svcs cron
```

```
STATE          STIME    FMRI
online         14:38:30 svc:/system/cron:default
```

```bash
$ svcadm -v disable system/cron:default
system/cron:default disabled.
```

```bash
svcs cron
STATE          STIME    FMRI
disabled       20:35:25 svc:/system/cron:default
```

```bash
ps -ef
```

```bash
$ svcadm -v enable system/cron:default
system/cron:default enabled.
```

```bash
$ svcs cron
STATE          STIME    FMRI
online         20:35:59 svc:/system/cron:default
```

```bash
$ ps -ef
```

To temporarily disable the cron service:

```bash
svcadm -v disable -t system/cron:default
svc:/system/cron:default temporarily disabled.
```

## Managing a Non-SMF Service

init.d is here:

```bash
$ svcs
```

```bash
$ ps -ef
```

```bash
ls /etc/init.d/volmgt
/etc/init.d/volmgt
/etc/init.d/volmgt stop
```

```bash
$ ps -ef
```

```bash
$ /etc/init.d/volmgt start
volume management starting.
```

```bash
$ ps -ef
```

```bash
svcs
```

## Creating a Service Managed by SMF

This procedure can be a bit complex for the uninitiated. Here's the chronological order to follow:

- Determine which milestones and run levels this service should be available at, and the appropriate command to start and stop the service.
- Establish relationships between dependencies, the service, and other services.
- Create a script in `/lib/svc/method` to start the process if necessary.
- Create an .xml file in the appropriate subdirectory.
- Make a copy of the "Service Repository Database".
- Integrate this script into SMF using the svccfg utility.

Create the file `/lib/svc/method/newservice`:

```bash
#!/sbin/sh
#
# Copyright 2004 Sun Microsystems, Inc.  All rights reserved.
# Use is subject to license terms.
#
# ident "@(#)newservice 1.14    04/08/30 SMI"

case "$1" in
'start')
        /usr/bin/newservice &
;;

'stop')
	/usr/bin/pkill -x -u 0 newservice
        ;;
*)
        echo "Usage: PAGECONTENT { start | stop }"
        ;;
esac
exit 0
```

Set permissions:

```bash
chmod 744 /lib/svc/method/newservice
```

Then create the file `/var/svc/manifest/site/newservice.xml`:

```xml
<?xml version="1.0"?>
<!DOCTYPE service_bundle SYSTEM "/usr/share/lib/xml/dtd/service_bundle.dtd.1">
<!--
        Copyright 2004 Sun Microsystems, Inc.  All rights reserved.
        Use is subject to license terms.

        ident   "@(#)newservice.xml     1.2     04/09/13 SMI"
-->

<service_bundle type='manifest' name='OPTnew:newservice'>

<service
        name='site/newservice'
        type='service'
        version='1'>

        <single_instance/>

        <exec_method
                type='method'
                name='start'
                exec='/lib/svc/method/newservice start'
                timeout_seconds='30' />

        <exec_method
                type='method'
                name='stop'
                exec='/lib/svc/method/newservice stop'
                timeout_seconds='30' />

        <property_group name='startd' type='framework'>
                <propval name='duration' type='astring' value='transient' />
        </property_group>

        <instance name='default' enabled='true' />

        <stability value='Unstable' />

        <template>
                <common_name>
                        <loctext xml:lang='C'>
                                New service
                        </loctext>
                </common_name>
        </template>
</service>

</service_bundle>
```

```bash
cd /var/svc/manifest/milestone
cp multi-user.xml /var/tmp
vi multi-user.xml
```

Here's an example of the content:

```xml
	<dependency
                name='fs'
                grouping='require_all'
                restart_on='none'
                type='service'>
                <service_fmri value='svc:/system/filesystem/local' />
        </dependency>

        <dependency
                name='newservice'
                grouping='require_all'
                restart_on='none'
                type='service'>
                <service_fmri value='svc:/site/newservice' />
        </dependency>
```

The new service must be imported into SMF:

```bash
svccfg import /var/svc/manifest/site/newservice.xml
```

Now it should be visible:

```bash
$ svcs newservice
STATE          STIME    FMRI
online          8:43:45 svc:/site/newservice:default
```

It should also be possible to manipulate the service using:

```bash
$ svcadm -v disable site/newservice

site/newservice disabled.
```

```bash
$ svcs newservice

STATE          STIME    FMRI
disabled        9:11:38 svc:/site/newservice:default
```

```bash
svcadm -v enable site/newservice

site/newservice enabled.
```

```bash
svcs newservice

STATE          STIME    FMRI
online          9:11:54 svc:/site/newservice:default
```

We can see that the multiuser milestone for our new service is necessary to finish:

```bash
$ svcs -d milestone/multi-user:default
STATE          STIME    FMRI
disabled        8:43:16 svc:/platform/sun4u/sf880drd:default
online          8:43:16 svc:/milestone/name-services:default
online          8:43:33 svc:/system/rmtmpfiles:default
online          8:43:42 svc:/network/rpc/bind:default
online          8:43:46 svc:/milestone/single-user:default
online          8:43:46 svc:/system/utmp:default
online          8:43:47 svc:/system/system-log:default
online          8:43:47 svc:/system/system-log:default
online          8:43:49 svc:/system/filesystem/local:default
online          8:44:01 svc:/system/mdmonitor:default
online          9:11:54 svc:/site/newservice:default
```

## Creating a Non-SMF Managed Service

- First, we'll create our script in init.d:

```bash
vi /etc/init.d/filename
```

See above for content, then:

```bash
chmod 744 /etc/init.d/filename
chgrp sys /etc/init.d/filename
```

- Now let's create the proper links (do this for each desired runlevel):

```bash
cd /etc/init.d
ln filename /etc/rc#.d/S##filename
ln filename /etc/rc#.d/K##filename
```

- Let's check:

```bash
ls -li /etc/init.d/filename
ls -li /etc/rc#.d/S##filename
ls -li /etc/rc#.d/K##filename
```

- Now let's test:

```bash
/etc/init.d/filename start
```

## Setting Boot Time for Milestones

Here's an example:

```bash
svcadm -v milestone -d multi-user-server:default
```

And the available options:

- all
- none
- svc:/milestone/single-user:default
- svc:/milestone/multi-user:default
- svc:/milestone/multi-user-server:default

Also remember to make a copy of the milestones database:

```bash
pstop svc.startd
pkill svc.configd
cp /etc/svc/repository.db /etc/svc/safe_repository.db
cp /lib/svc/seed/global.db /etc/svc/repository.db
init 0
ok boot -m verbose
```

## FAQ

### svc.configd: smf(5) database integrity check of: /etc/svc/repository.db

I encountered a message like this after a reboot, thanks UFS. The full message was:

```
svc.configd: smf(5) database integrity check of:

    /etc/svc/repository.db

  failed.  The database might be damaged or a media error might have
  prevented it from being verified.  Additional information useful to
  your service provider is in:

    /etc/svc/volatile/db_errors

  The system will not be able to boot until you have restored a working
  database.  svc.startd(1M) will provide a sulogin(1M) prompt for recovery
  purposes.  The command:

    /lib/svc/bin/restore_repository

  can be run to restore a backup version of your repository.  See
  http://sun.com/msg/SMF-8000-MY for more information.
```

To resolve this issue:

- Reboot in failsafe mode (grub)
- Fix all filesystem fragmentation issues (wizard)
- Mount your root partition in read/write mode in /a (wizard)
- Chroot the /a partition:

```bash
chroot /a /a/bin/bash
```

- Run the command /lib/svc/bin/restore_repository and tell it to repair /boot:

```
The following backups of /etc/svc/repository.db exists, from
oldest to newest:

... list of backups ...

The backups are named based on their type and the time when they were taken.
Backups beginning with "boot" are made before the first change is made to
the repository after system boot.  Backups beginning with "manifest_import"
are made after svc:/system/manifest-import:default finishes its processing.
The time of backup is given in YYYYMMDD_HHMMSS format.

Please enter one of:
        1) boot, for the most recent post-boot backup
        2) manifest_import, for the most recent manifest_import backup.
        3) a specific backup repository from the above list
        4) -seed-, the initial starting repository.  (All customizations
           will be lost.)
        5) -quit-, to cancel.

Enter response [boot]:
```

Just press Enter here. And confirm by typing yes:

```
After confirmation, the following steps will be taken:

svc.startd(1M) and svc.configd(1M) will be quiesced, if running.
/etc/svc/repository.db
    -- renamed --> /etc/svc/repository.db_old_YYYYMMDD_HHMMSS
/etc/svc/volatile/db_errors
    -- copied --> /etc/svc/repository.db_old_YYYYMMDD_HHMMSS_errors
repository_to_restore
    -- copied --> /etc/svc/repository.db
and the system will be rebooted with reboot(1M).

Proceed [yes/no]? yes
```

## Resources
- [Solaris Features: Service Management Facility](https://www.c0t0d0s0.org/archives/4149-Solaris-Features-Service-Management-Facility.html)
- [Using Service Management Facility (SMF)](/pdf/using_service_management_facility_smf.pdf)
