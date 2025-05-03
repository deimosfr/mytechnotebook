---
weight: 999
url: "/Gérer_ses_updates_avec_Update_Manager_et_smpatch/"
title: "Managing Updates with Update Manager and smpatch"
description: "How to manage Solaris updates using Update Manager GUI and smpatch CLI tools"
categories: ["Linux", "Storage", "Solaris"]
date: "2009-11-21T07:11:00+02:00"
lastmod: "2009-11-21T07:11:00+02:00"
tags: ["Solaris", "Updates", "Patching", "CLI", "System Administration"]
toc: true
---

## Introduction

As I love Debian and some of my servers are running on Solaris, I had to get in touch with Solaris update solutions. They have a GUI called Solaris Update Manager and the CLI version called smpatch.

As it is on servers, I don't have a graphical interface and I need to run updates with command lines. That's what I'll specifically talk about in this documentation. List of interesting binaries:

- pprosetup: Used to set the rules for downloading and applying patches.
- pprosvc: The automation service program for Patch Manager.
- smpatch: Used to actually download, apply, and remove the patches specified on the command line.

## Update Manager

Just to give you a quick idea, when you have a graphical interface, it looks like this:

![Sum-installing.jpg](/images/sum-installing.avif)

I won't explain how it works as it is very simple.

## smpatch

smpatch is the command line version of Update Manager.

### Configuration

#### Set proxy settings

```bash
smpatch set patchpro.proxy.host=fqdn_of_proxy_host
smpatch set patchpro.proxy.port=proxy_port
smpatch set patchpro.proxy.user=proxy_username_if_needed
```

Now set password by prompt:

```bash
smpatch set patchpro.proxy.passwd
```

#### Configuration files

Here are the configuration files:

- System defaults:
  - `/var/sadm/install/admin/default`
- PatchPro defaults:
  - `/etc/opt/SUNWppro/etc/patchpro.conf`
  - `/opt/SUNWppro/lib/.proxypw`
- SunSolve Account Options
  - `/etc/patch/patch.conf`
  - `/etc/patch/secret.conf`

### Analyze

The analyze will check what needs to be updated on your server:

```bash
$ smpatch analyze
122213-34 GNOME 2.6.0_x86: GNOME Desktop Patch
119901-09 GNOME 2.6.0_x86: Gnome libtiff - library for reading and writing TIFF Patch
142293-01 SunOS 5.10_x86: Place Holder patch
141445-09 SunOS 5.10_x86: kernel patch
141031-05 SunOS 5.10_x86: passwd patch
142335-01 SunOS 5.10_x86: mixer patch
119784-13 SunOS 5.10_x86: bind patch
126869-04 SunOS 5.10_x86: SunFreeware bzip2 patch
120273-27 SunOS 5.10_x86: SMA patch
123896-15 SunOS 5.9_x86 5.10_x86: Common Agent Container (cacao) runtime 2.2.3.1 upgrade patch 15
121082-08 SunOS 5.10_x86: Disable Transport Agentry for Sun Update Connection Hosted EOL
118778-12 SunOS 5.10_x86: Sun GigaSwift Ethernet 1.0 driver patch
141503-02 SunOS 5.10_x86: auditconfig patch
141525-05 SunOS 5.10_x86: ssh and openssl patch
141511-04 SunOS 5.10_x86: ehci, ohci, uhci patch
...
```

### Upgrade

Now you want to upgrade. You have 2 choices:

- Automatic update: will check what is needed and will upgrade your system as well
- Manual update: you'll need to download them and install them afterward

#### Automatic

The automatic **update** will do everything on its own (download + install):

```bash
$ smpatch update
122213-34 has been validated.
119901-09 has been validated.
142293-01 has been validated.
141445-09 has been validated.
141031-05 has been validated.
...
```

That's all. It will tell you if it's necessary to reboot or not.

#### Manual

For the manual upgrade, you'll need first to download updates:

```bash
$ smpatch download
125534-15 has been validated.
126364-08 has been validated.
...
```

And then install them:

```bash
$ smpatch add
```

#### Default storage location

By default, updates are stored in this folder `/var/sadm/spool/`. Updates are in .jar format:

```bash
$ ls /var/sadm/spool/
118668-23.jar                              121431-43.jar                              125953-19.jar                              139100-02.jar                              141589-03.jar
118778-12.jar                              122213-34.jar                              125993-04.jar                              139621-01.jar                              141591-01.jar
119116-35.jar                              122260-02.jar                              126018-05.jar                              140018-03.jar                              141879-08.jar
119214-20.jar                              122471-03.jar                              126036-07.jar                              140130-10.jar                              142241-01.jar
```

## References

http://www.cuddletech.com/blog/pivot/entry.php?id=579  
http://docs.huihoo.com/opensolaris/system-administration-guide-basic-administration/html/ch22s06.html  
http://bob.vancleef.org/index.php?name=News&file=article&sid=56  
http://www.sun.com/service/sunupdate/flash_content.html
