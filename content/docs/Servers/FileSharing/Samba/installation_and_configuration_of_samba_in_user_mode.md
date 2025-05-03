---
weight: 999
url: '/Installation_et_configuration_de_Samba_en_mode_"User"/'
title: "Installation and Configuration of Samba in 'User' Mode"
description: "A guide for setting up Samba in user mode to share folders with Windows and Unix systems securely"
categories: ["FreeBSD", "Security", "Ubuntu"]
date: "2011-10-06T12:53:00+02:00"
lastmod: "2011-10-06T12:53:00+02:00"
tags: ["Samba", "File Sharing", "Security", "Networking"]
toc: true
---

## Introduction

Samba is a free software licensed under the GPL that supports the SMB/CIFS protocol. This protocol is used by Microsoft for sharing various resources (files, printers, etc.) between computers running Windows. Samba allows Unix systems to access resources from these systems and vice versa.

Previously, PCs with DOS and early versions of Windows sometimes needed to install a TCP/IP stack and a set of Unix applications: an NFS client, FTP, telnet, lpr, etc. This was cumbersome and penalizing for PCs of that time, and it also forced users to adopt two sets of habits, adding those of UNIX to those of Windows. Samba therefore adopts the opposite approach.

Its name comes from the file and print sharing protocol from IBM and reused by Microsoft called SMB (Server message block), to which two vowels "a" were added: "SaMBa".

Samba was originally developed by Andrew Tridgell in 1991 and now receives contributions from about twenty developers from around the world under his coordination. He gave it this name by choosing a name close to SMB by querying a Unix dictionary with the command grep: `grep "^s.*m.*b" /usr/dict/words`

When both file sharing systems (NFS, Samba) are installed for comparison, Samba proves less efficient than NFS in terms of transfer rates.

Nevertheless, a study showed that Samba 3 was up to 2.5 times faster than Windows Server 2003's SMB implementation. See the information on LinuxFr.

However, Samba is not compatible with IPv6.

**The "User" mode allows you to share folders simply by user. You then need a login and a password. This is a sufficiently secure solution for small businesses.**

## Installation

### Debian

To install Samba:

```bash
aptitude install samba
```

### FreeBSD

For FreeBSD:

```bash
pkg_add -r samba34
```

## Configuration

To configure Samba, edit the file `/etc/samba/smb.conf` (`/usr/local/etc/smb.conf` under FreeBSD):

```bash
#======================= Global Settings =====================================
[global]
        # Samba server name
        server string = Samba
        # Socket optimization
        socket options = TCP_NODELAY SO_RCVBUF=8192 SO_SNDBUF=8192
        # Workgroup name
        workgroup = workgroup
        # Samba server level
        os level = 20

        ## Restrictions ##
        # Deny everyone
        hosts deny = ALL
        # Allow only requests from these IPs
        hosts allow = 192.168.0.0/255.255.255.0 127.0.0.1 10.8.0.0/255.255.255.0
        bind interfaces only = yes
        # Allow only requests from this network interface
        interfaces = eth0

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
        # User mode
        security = user
        encrypt passwords = true
        unix password sync = no
        passwd program = /usr/bin/passwd %u
        passwd chat = *Enter\snew\sUNIX\spassword:* %n\n *Retype\snew\sUNIX\spassword:* %n\n .
        # Do not allow these users
        invalid users = root

        ## Restrictions ##
        # Hide special files
        hide special files = no
        # Hide unreadable files
        hide unreadable = no
        # Hide hidden files (starting with a ".")
        hide dot files = no

        ## Resolve office save problems ##
        # Solves compatibility issues with versions > MS Office 2002
        oplocks = no

#======================= Shares ==============================================

[Homes]
        comment = Home Directories
        browseable = yes # Allows browsing a directory tree
        read only = no # No read-only
        writable = yes # Allows writing
        create mask = 0700 # File creation rights
        directory mask = 0700 # Directory creation rights
        veto files = /.DS_Store/.fuse_*/ # Do not display objects: ".DS_Store" and ".fuse_*"

[Sauvegardes]
        comment = Sauvegardes
        path = /saves # Shared folder
        browseable = yes
        writable = yes
        validusers = deimos, @smbusers # Authorized users. To define a group, start with "@"
        veto files = /.DS_Store/.fuse_*/

[Share2]
        comment = Sauvegardes
        path = /share
        # I only allow authenticated access
        guest ok = no
        # Everyone can read
        read only = yes
        # But only my smbusers group can write
        write list = @smbusers
```

Some explanations:

- First configure the data in Global
- Set the OS level < 20 unless it acts as a domain controller, then > 50

Adapt all this to your configuration. Then restart Samba:

```bash
/etc/init.d/samba restart
```

Or like this under FreeBSD:

```bash
/usr/local/etc/rc.d/samba restart
```

Now, you need to add users! It's quite simple, but it's the kind of thing you often forget:

```bash
smbpasswd -a deimos
```

This will add my user deimos. And here is the list of possible options:

```bash
smbpasswd --help

When run by root:
    smbpasswd [options] [username]
otherwise:
    smbpasswd [options]

options:
  -L                   local mode (must be first option)
  -h                   print this usage message
  -s                   use stdin for password prompt
  -c smb.conf file     Use the given path to the smb.conf file
  -D LEVEL             debug level
  -r MACHINE           remote machine
  -U USER              remote username extra options when run by root or in local mode:
  -a                   add user
  -d                   disable user
  -e                   enable user
  -i                   interdomain trust account
  -m                   machine trust account
  -n                   set no password
  -W                   use stdin ldap admin password
  -w PASSWORD          ldap admin password
  -x                   delete user
  -R ORDER             name resolve order
```

Under FreeBSD, you can find the file containing the list of authorized users in `/usr/local/etc/samba34`.

For FreeBSD still, if you want Samba to start automatically at boot:

```bash
echo 'samba_enable="YES"' >> /etc/rc.conf.local
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
aptitude install smbfs
```

Then, simply create a folder and mount the share in it:

```bash
mkdir saves
mount -t cifs -o username=user,password=password //192.168.0.1/saves ./saves
```

## FAQ

### Test your configuration

If you have some problems, here is a way to check your configuration (FreeBSD):

```bash
/usr/local/bin/testparm -s
```

### Unable to connect to CUPS server localhost:631 - Connection refused

If you don't use a CUPS print server, then make the following changes to disable it and allow Samba to start:

```bash
        # Disable printers
        load printers = no
        show add printer wizard = no
        printing = none
        printcap name = /dev/null
        disable spoolss = yes
```

### Migrate your smbpasswd file to tdbsam

It is possible that when updating your Samba, your smbpasswd file no longer works and must be replaced by a tdbsam. All your users will then no longer work. The easiest way is to convert all your old accounts to this new format:

```bash
cd /etc/samba
pdbedit -i smbpasswd -e tdbsam
```

You may need to modify your samba configuration file to add this:

```bash
passdb backend = tdbsam:/var/lib/samba/passdb.tdb
```

## Resources

- If you want to push folder permissions, check the documentation on [ACL: Implementation of NT-type rights]({{< ref "docs/Linux/FilesystemsAndStorage/acl-implementing-nt-type-permissions-on-linux.md" >}}).

- [Other documentation on Samba]({{< ref "docs/Servers/#samba" >}})
- [CIFS Solaris Workgroup](/pdf/cifs_solaris_workgroup.pdf)
- http://www.csua.berkeley.edu/~ranga/notes/freebsd_samba.html+samba+freebsd&cd=6&hl=fr&ct=clnk&client=ubuntu
- http://www.tobanet.de/dokuwiki/samba:upgrade
