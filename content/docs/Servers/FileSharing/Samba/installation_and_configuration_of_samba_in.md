---
weight: 999
url: "/Installation_et_configuration_de_Samba_en_mode_\"Share\"/"
title: "Installation and Configuration of Samba in \"Share\" Mode"
description: "Learn how to install and configure Samba in Share mode, a simple way to share folders without authentication requirements."
categories: ["Linux", "Security"]
date: "2007-08-27T15:11:00+02:00"
lastmod: "2007-08-27T15:11:00+02:00"
tags: ["Samba", "File Sharing", "SMB", "CIFS", "Networking"]
toc: true
---

## Introduction

Samba is an open source software licensed under the GPL that supports the SMB/CIFS protocol. This protocol is used by Microsoft for sharing various resources (files, printers, etc.) between computers running Windows. Samba allows Unix systems to access resources from these systems and vice versa.

Previously, PCs equipped with DOS and early versions of Windows sometimes had to install a TCP/IP stack and a set of Unix-originated software: NFS client, FTP, telnet, lpr, etc. This was heavy and penalizing for PCs of that time, and it also forced users to adopt a double set of habits, adding those of UNIX to those of Windows. Samba takes the opposite approach.

Its name comes from the file and print sharing protocol from IBM reused by Microsoft called SMB (Server message block), to which the two vowels "a" were added: "SaMBa".

Samba was originally developed by Andrew Tridgell in 1991, and today receives contributions from about twenty developers from around the world under his coordination. He gave it this name by choosing a name close to SMB by querying a Unix dictionary with the command grep: `grep "^s.*m.*b" /usr/dict/words`

When both file sharing systems (NFS, Samba) are installed for comparison, Samba proves less efficient than NFS in terms of transfer rates.

Nevertheless, a study has shown that Samba 3 was up to 2.5 times faster than the SMB implementation of Windows Server 2003. See the information on LinuxFr.

However, Samba is not compatible with IPv6.

**The "Share" mode allows simple folder sharing. No login or password needed, everyone has access to everything, which is not a secure solution, but it's a simple one.**

## Installation

To install Samba:

```bash
apt-get install samba
```

## Configuration

Before you start, define a directory that you want to share (example: `/home/share`):

```bash
mkdir /home/share
chmod 777 /home/share
```

We give it full permissions.

To configure Samba, edit the file `/etc/samba/smb.conf`:

```bash
#======================= Global Settings =====================================
[global]
        server string = Samba # Samba server name
        socket options = TCP_NODELAY SO_RCVBUF=8192 SO_SNDBUF=8192 # Socket optimization
        workgroup = workgroup # Workgroup name
        os level = 20 # Samba server level
 
        ## Encoding ## European display with accents
        dos charset = 850
        display charset = UTF8
 
        ## Name resolution ## Name resolutions
        dns proxy = no
        wins support = no
        name resolve order = lmhosts host wins bcast
 
        ## Logs ##
        max log size = 50
        log file = /var/log/samba/%m.log
        syslog only = no
        syslog = 0
        panic action = /usr/share/samba/panic-action %d
 
        ## Passwords ##
        password server = None # No password server in share mode
        security = SHARE # Chosen mode
        invalid users = root # Do not authorize these users.
 
        ## Restrictions ##
        hide special files = no # Hide special files
        hide unreadable = no # Hide unreadable files
        hide dot files = no # Hide hidden files (starting with a ".")
 
        ## Resolve office save problems ##
        oplocks = no # Resolves compatibility issues with versions > MS Office 2002
 
#======================= Shares ==============================================
 
# tmp share
[tmp]
   comment = Temporary file space
   path = /tmp
   read only = no
   public = yes
 
# share share
[share]
   comment = Share file space
   path = /home/share
   read only = no
   public = yes
```

Some explanations:

- First configure the data in Global
- Set the OS level < 20 unless it acts as a domain controller, then > 50

Adapt all this to your configuration. Then restart Samba:

```bash
/etc/init.d/samba restart
```

## Connection

### Windows

To connect from Windows, in a link window, type this:

```
\\IP_of_samba_server\Share_name
```

You will access the share directly.

### Unix (Linux/Mac...)

You must have smbfs installed before continuing:

```bash
apt-get install smbfs
```

Then, just create a folder and mount the share inside:

```bash
mkdir test
mount -t cifs -o username=nobody,password=nobody //192.168.0.1/tmp ./test
```

## Resources
- [Documentation on a Complete auto discovery and mounting solution with SMB shares](/pdf/complete_auto_discovery_mounting_smb_shares.pdf)
