---
weight: 999
url: "/Securiser_son_architecture_avec_SELinux/"
title: "Secure Your Architecture with SELinux"
description: "A comprehensive guide to understanding and implementing SELinux security policies in Linux systems"
categories: ["Linux", "Security", "System Administration"]
date: "2011-10-06T12:51:00+02:00"
lastmod: "2011-10-06T12:51:00+02:00"
tags: ["SELinux", "Security", "Linux", "Policies", "Access Control"]
toc: true
---

## Introduction

[Security-Enhanced Linux, abbreviated as SELinux](https://en.wikipedia.org/wiki/Selinux), is a LSM (Linux security module) that allows defining a MAC (mandatory access control) access policy to elements of a Linux-based system. Initiated by the NSA based on work conducted with SCC and the University of Utah in the USA (DTMach prototypes, DTOS, FLASK project), its architecture separates policy enforcement from its definition. It notably allows classifying applications of a system into different groups with finer access levels. It also allows assigning a confidentiality level for accessing system objects, such as file descriptors, according to an MLS (Multi Level Security) security model. SELinux uses the Bell LaPadula model with SCC Type enforcement for integrity. It is free software, with some parts licensed under GNU GPL or BSD.

In practice, the innovation is based on defining extended attributes in the UNIX filesystem. Beyond the concept of "read, write, execute rights" for a given user, SELinux defines for each file or process:

- A virtual user (or collection of roles);
- A role;
- A security context.

## Usage

### Defining the Mode

This is where you define the SELinux mode and its type. We'll use the targeted mode which is single-user mode (no one above root). The mls mode is equivalent to RBAC which allows many more security groups and users (for very large companies):

```bash
# This file controls the state of SELinux on the system.
# SELINUX= can take one of these three values:
#     enforcing - SELinux security policy is enforced.
#     permissive - SELinux prints warnings instead of enforcing.
#     disabled - No SELinux policy is loaded.
SELINUX=enforcing
# SELINUXTYPE= can take one of these two values:
#     targeted - Targeted processes are protected,
#     mls - Multi Level Security protection.
SELINUXTYPE=targeted
```

### Getting the Current Mode and Changing It

First, we need to know in which mode we are:

```bash
$ getenforce
Enforcing
```

Here we see I'm in enforcing mode. If I want to switch to permissive mode, I run this command with argument 0:

```bash
setenforce 0
```

And I set it to 1 to go back to enforcing mode.

We can verify that SELinux is properly enabled by listing processes. The SELinux attributes will display with the 'Z' argument:

```bash
$ ps axZ | grep cron
system_u:system_r:crond_t:s0-s0:c0.c1023 1417 ? Ss    0:00 crond
system_u:system_r:crond_t:s0-s0:c0.c1023 1428 ? Ss    0:00 /usr/sbin/atd
unconfined_u:unconfined_r:unconfined_t:s0-s0:c0.c1023 2082 pts/0 S+   0:00 grep cron
```

We can also do this on folders or files:

```bash
$ ls -Z /root
-rw-------. root root system_u:object_r:admin_home_t:s0 anaconda-ks.cfg
-rw-r--r--. root root system_u:object_r:admin_home_t:s0 install.log
-rw-r--r--. root root system_u:object_r:admin_home_t:s0 install.log.syslog
```

#### Disable Only One Domain

You also have the option to disable only one domain if desired. Let's take Apache as an example:

```bash
semanage permissive -a httpd_t
```

This service is now in permissive state.

### Analyzing Blocks

When SELinux decides to block certain access, there are several ways to analyze and accept certain false positives. First, there's the 'audit2allow' command:

```bash
 audit2allow -la


#============= httpd_t ==============
allow httpd_t admin_home_t:file { read getattr open };
```

You also have logs that provide information about blocks in `/var/log/audit/audit.log`:

```bash
type=AVC msg=audit(1317855069.699:15772): avc:  denied  { getattr } for  pid=2273 comm="httpd" path="/var/www/html/mon_fichier.txt" dev=dm-0 ino=30073 scontext=unconfined_u:system_r:httpd_t:s0 tcontext=unconfined_u:object_r:admin_home_t:s0 tclass=file
type=SYSCALL msg=audit(1317855069.699:15772): arch=c000003e syscall=4 success=yes exit=0 a0=7f300213e3c8 a1=7ffffd59c590 a2=7ffffd59c590 a3=0 items=0 ppid=2268 pid=2273 auid=0 uid=48 gid=48 euid=48 suid=48 fsuid=48 egid=48 sgid=48 fsgid=48 tty=(none) ses=2 comm="httpd" exe="/usr/sbin/httpd" subj=unconfined_u:system_r:httpd_t:s0 key=(null)
type=AVC msg=audit(1317855069.700:15773): avc:  denied  { read } for  pid=2273 comm="httpd" name="mon_fichier.txt" dev=dm-0 ino=30073 scontext=unconfined_u:system_r:httpd_t:s0 tcontext=unconfined_u:object_r:admin_home_t:s0 tclass=file
type=AVC msg=audit(1317855069.700:15773): avc:  denied  { open } for  pid=2273 comm="httpd" name="mon_fichier.txt" dev=dm-0 ino=30073 scontext=unconfined_u:system_r:httpd_t:s0 tcontext=unconfined_u:object_r:admin_home_t:s0 tclass=file
type=SYSCALL msg=audit(1317855069.700:15773): arch=c000003e syscall=2 success=yes exit=11 a0=7f300213e480 a1=80000 a2=0 a3=2 items=0 ppid=2268 pid=2273 auid=0 uid=48 gid=48 euid=48 suid=48 fsuid=48 egid=48 sgid=48 fsgid=48 tty=(none) ses=2 comm="httpd" exe="/usr/sbin/httpd" subj=unconfined_u:system_r:httpd_t:s0 key=(null)
```

I can see here that a file named mon_fichier.txt was blocked for httpd_t because the object is not correct.

### The Contexts

As you've seen, there are some special attributes called contexts. We can display this list of contexts via the 'semanage' command which allows us to manage contexts:

```bash
$ semanage fcontext -l | more
SELinux fcontext                                   type               Context

/                                                  directory          system_u:object_r:root_t:s0
/.*                                                all files          system_u:object_r:default_t:s0
/HOME_DIR/\.Xdefaults                              all files          system_u:object_r:config_home_t:s0
/HOME_DIR/\.xine(/.*)?                             all files          system_u:object_r:config_home_t:s0
/[^/]+                                             regular file       system_u:object_r:etc_runtime_t:s0
/\.autofsck                                        regular file       system_u:object_r:etc_runtime_t:s0
/\.autorelabel                                     regular file       system_u:object_r:etc_runtime_t:s0
/\.journal                                         all files          <<None>>
/\.suspended                                       regular file       system_u:object_r:etc_runtime_t:s0
/a?quota\.(user|group)                             regular file       system_u:object_r:quota_db_t:s0
/afs                                               directory          system_u:object_r:mnt_t:s0
/bin                                               directory          system_u:object_r:bin_t:s0
/bin/.*                                            all files          system_u:object_r:bin_t:s0
/bin/alsaunmute                                    regular file       system_u:object_r:alsa_exec_t:s0
/bin/bash                                          regular file       system_u:object_r:shell_exec_t:s0
/bin/bash2                                         regular file       system_u:object_r:shell_exec_t:s0
/bin/d?ash                                         regular file       system_u:object_r:shell_exec_t:s0
/bin/dbus-daemon                                   regular file       system_u:object_r:dbusd_exec_t:s0
/bin/dmesg                                         regular file       system_u:object_r:dmesg_exec_t:s0
/bin/fish                                          regular file       system_u:object_r:shell_exec_t:s0
/bin/fusermount                                    regular file       system_u:object_r:fusermount_exec_t:s0
...
```

Contexts match regexes and are authorized this way.

#### Context Modification

Let's say for my website, I want to create an index file. I'm in `/tmp` and create my index. At that moment, when the file is created on disk, SELinux will tag the index file and specify that it belongs to the `/tmp` directory.

So when I move it to `/var/www`, it will still keep those attributes and the Apache server won't be able to use this file. To fix the issue, I have 2 choices:

- Restore the rights defined in the context database for the parent directory.
- Reassign the correct rights to the specific file

##### Context Restoration

To restore contexts, we'll use the restcon command:

```bash
restcon -Rv /var/www
```

And now we have reset all SELinux rights in `/var/www`.

##### Context Reassignment

To reassign the correct rights, I need to apply the security policy to this file:

```bash
chcon -v -t httpd_sys_content_t /var/www/index.html
```

httpd_sys_content_t: this is the desired context type for the `/var/www` directory

To find the right context, use the semanage command as seen above:

```bash
$ semanage fcontext -l | grep "/var/www"
/var/www(/.*)?                                     all files          system_u:object_r:httpd_sys_content_t:s0
/var/www(/.*)?/logs(/.*)?                          all files          system_u:object_r:httpd_log_t:s0
/var/www/[^/]*/cgi-bin(/.*)?                       all files          system_u:object_r:httpd_sys_script_exec_t:s0
...
```

You can see in the last column that the context we're looking for is "httpd_sys_content_t".

##### Adding a Context

This is something you should avoid doing to solve problems, but rather use to improve security or customize it for your needs. Here we'll add a context to fix the problem with the file listed above that belongs to admin and is located in `/var/www/html`. We'll add a context with the rights of the root directory:

```bash
semanage fcontext -a -t admin_home_t "/var/www(/.*)?"
```

Then you can verify your changes:

```bash
$ semanage fcontext -l | grep www
/srv/([^/]*/)?www(/.*)?                            all files          system_u:object_r:httpd_sys_content_t:s0
/usr/share/awstats/wwwroot(/.*)?                   all files          system_u:object_r:httpd_awstats_content_t:s0
/usr/share/awstats/wwwroot/cgi-bin(/.*)?           all files          system_u:object_r:httpd_awstats_script_exec_t:s0
/var/www(/.*)?                                     all files          system_u:object_r:admin_home_t:s0
```

Now you just need to run "restcon" to fix the permissions.

#### Port Blocking

SELinux also allows only certain services to run on specific ports. Proof:

```bash
$ semanage port -l | grep http
http_cache_port_t              tcp      3128, 8080, 8118, 10001-10010
http_cache_port_t              udp      3130
http_port_t                    tcp      80, 443, 488, 8008, 8009, 8443
pegasus_http_port_t            tcp      5988
pegasus_https_port_t           tcp      5989
```

If you want to run Apache on another port, for example, you'll need to add it to the contexts list:

```bash
semanage port -a -t http_port_t -p tcp 81
```

I chose port 81 in this example.

### Booleans

Booleans are another type of blocking that SELinux implements, typically found on well-known services. To get this list:

```bash
$ getsebool -a
abrt_anon_write --> off
allow_console_login --> on
allow_corosync_rw_tmpfs --> off
allow_cvs_read_shadow --> off
allow_daemons_dump_core --> on
allow_daemons_use_tty --> on
allow_domain_fd_use --> on
allow_execheap --> off
allow_execmem --> on
allow_execmod --> on
allow_execstack --> on
allow_ftpd_anon_write --> off
allow_ftpd_full_access --> off
allow_ftpd_use_cifs --> off
...
```

To change a value, simply do:

```bash
setsebool allow_ftpd_full_access on
```

You can verify afterwards in two ways:

```bash
$ getsebool allow_ftpd_full_access
```

or

```bash
$ semanage boolean -l
```

## FAQ

### My System Refuses to Boot Because of SELinux

To fix this problem, at the grub boot, edit the kernel line and add this at the end:

```
enforcing=0
```

This will set permissive mode at machine boot so you can fix your problem.

### How to Reapply All Security Policies to My System?

If you want to reset all your SELinux security policy on your machine, there are 2 solutions. The first is messier; it consists of checking and reapplying all changes on the fly:

```bash
restorecon -R /
```

I told you... it's ugly! However, another cleaner solution that will correctly reapply all permissions at the next reboot is to create a file at the root:

```bash
touch /autorelabel
```

### I Have a SELinux Problem and Nothing in My Logs

I encountered this with Samba which was causing problems, but neither audit2allow nor logs showed anything. To solve this problem and see the log messages, tell it to log everything:

```bash
semanage dontaudit off
```

Then you just need to check the logs (`/var/log/audit/audit.log` and `/var/log/messages`).

## Resources
- [SELinux, the Kernel Security Agency]({{< ref "docs/Linux/Security/secure_your_architecture_with_selinux.md" >}})
