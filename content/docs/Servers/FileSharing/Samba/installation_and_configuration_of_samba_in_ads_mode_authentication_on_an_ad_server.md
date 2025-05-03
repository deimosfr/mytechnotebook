---
weight: 999
url: '/Installation_et_configuration_de_Samba_en_mode_"ADS"_(Authentification_sur_un_serveur_AD)/'
title: "Installation and Configuration of Samba in ADS Mode (Authentication on an AD Server)"
description: "Guide on how to integrate a Linux server with Samba into an Active Directory domain for unified authentication"
categories: ["Linux", "Samba", "Active Directory", "Authentication"]
date: "2008-12-02T07:24:00+02:00"
lastmod: "2008-12-02T07:24:00+02:00"
tags:
  [
    "Samba",
    "Active Directory",
    "Kerberos",
    "Authentication",
    "Windows",
    "Linux",
  ]
toc: true
---

## Introduction

Samba is free software under GPL license supporting the SMB/CIFS protocol. This protocol is used by Microsoft for sharing various resources (files, printers, etc.) between computers running Windows. Samba allows Unix systems to access the resources of these systems and vice-versa.

Previously, PCs equipped with DOS and early Windows versions sometimes had to install a TCP/IP stack and a set of Unix-originated software: NFS client, FTP, telnet, lpr, etc. This was heavy and penalizing for the PCs of that time, and it also forced their users to adopt dual habits, adding those of UNIX to those of Windows. Samba therefore adopts the opposite approach.

Its name comes from the file and print sharing protocol from IBM and reused by Microsoft called SMB (Server message block), to which were added the two vowels a: "SaMBa".

Samba was originally developed by Andrew Tridgell in 1991 and today receives contributions from about twenty developers from around the world under his coordination. He gave it this name by choosing a name close to SMB by querying a Unix dictionary with the grep command:

```
grep "^s.*m.*b" /usr/dict/words
```

When the two file sharing systems (NFS, Samba) are installed for comparison, Samba proves less efficient than NFS in terms of transfer rates.

Nevertheless, a study has shown that Samba 3 was up to 2.5 times faster than the SMB implementation of Windows Server 2003. See the information on LinuxFr.

However, Samba is not compatible with IPv6.

**The "ADS" mode allows using an LDAP server on a MS Windows AD (Active Directory) to authenticate users for accessing Samba shares. This solution is quite complex to implement but ensures increased security (authentication via Kerberos) for your server.**

**Note: It is imperative to have set up [ACL: Implementation of NT-type rights]({{< ref "docs/Linux/FilesystemsAndStorage/acl-implementing-nt-type-permissions-on-linux.md" >}}) before continuing.**

Active Directory is at the heart of Microsoft Windows systems; it is in charge of managing user accounts, authentications, but also a large number of information about machines.

Active Directory relies on different network protocols:

- DNS (name resolution)
- LDAP (directory querying)
- Kerberos V (authentication, ticket distribution)
- NTP (date and time synchronization of machines)
- SMB/CIFS (resource sharing)

## Introduction to UNIX User Management

The management of LINUX user accounts is carried out using different components in accordance with UNIX philosophy (a program does one thing and does it well). Among these different actors, we find:

PAM (Pluggable Authentication Modules) allows, among other things, to select different authentication procedures and sources (e.g.: Authentication by smart cards, Databases, Directories...).

NSS (Name Services Switch) allows Unix to provide correspondence services between names of all kinds (machine names and user names) and the identifiers of these same objects for the machine (IP addresses and uid/gid) using various sources (Files, Directories...).

## Integration of a Samba Server into an Active Directory Domain

Integrating a Samba server into an Active Directory domain requires configuring a Kerberos client on the Samba machine. Kerberos is an authentication system that allows servers to authenticate users and communicate securely. In order to achieve the integration of the LINUX server to the AD domain, additional components are required.

### Kerberos

It is necessary to configure a Kerberos client to validate the identity of the LINUX server in the Microsoft network. This will communicate with the AD server to make "ticket" requests to the KDC which will be used to ensure the authenticity and security of communications.

#### Time Synchronization

First, we need to ensure that our AD and our Samba server are at the same time. For this, we just need to synchronize the time with an NTP server:

```bash
apt-get install ntpdate
```

```bash
ntpdate ntp.ciril.fr
```

#### Installing the Kerberos Client

Now we need to install Kerberos authentication:

```bash
apt-get install krb5-clients krb5-user
```

#### Configuring the Kerberos Client

Please edit the `/etc/krb5.conf` file:

```bash
[libdefaults]
	    default_realm = EXAMPLE.COM

[realms]
EXAMPLE.COM = {
	    kdc = ad.example.com
	    admin_server = ad.example.com
}
[domain_realms]
	    .example.com = EXAMPLE.COM
```

Here are the correspondences:  
**IMPORTANT: It is imperative to respect the case for all names.**

- EXAMPLE.COM: DNS name of the AD
- ad.example.com: FQDN

#### Verification of the Kerberos Connection

```bash
kinit Administrateur@EXAMPLE.COM
klist
kdestroy
```

## Samba

To access shared resources (SMB/CIFS) of the Windows domain, a Samba installation is required on the Linux server which will mainly play the role of an SMB/CIFS client.  
Samba will also use MSRPC commands to communicate with the AD server to perform various operations: adding information about the LINUX server in the directory, listing user/group accounts, transmitting authentication requests...

### Installation

To install Samba:

```bash
apt-get install samba smbclient
```

### Configuration

To configure Samba, edit the `/etc/samba/smb.conf` file:

```bash
#======================= Global Settings =====================================
[global]
        server string = Samba # Samba server name
        socket options = TCP_NODELAY SO_RCVBUF=8192 SO_SNDBUF=8192 # Socket optimization
        realm = EXAMPLE.COM # Kerberos REALM
        workgroup = workgroup # Domain name
        os level = 80 # Samba server level

        ## Restrictions ##
        hosts deny = ALL # Deny everyone
        hosts allow = 192.168.0.0/255.255.255.0 127.0.0.1 10.8.0.0/255.255.255.0 # Allow only requests from these IPs
        bind interfaces only = yes
        interfaces = eth0 # Allow only requests from this network interface

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
        security = ADS # Active Directory Server manages the security of shared resources
        encrypt passwords = true # Active Directory doesn't accept clear passwords
        unix password sync = no
        passwd program = /usr/bin/passwd %u
        passwd chat = *Enter\snew\sUNIX\spassword:* %n\n *Retype\snew\sUNIX\spassword:* %n\n .
        invalid users = root # Don't authorize these users.

        ## Restrictions ##
        hide special files = no # Hide special files
        hide unreadable = no # Hide unreadable files
        hide dot files = no # Hide hidden files (starting with a ".")

        ## Resolve office save problems ##
        oplocks = no # Resolves compatibility issues with versions > MS Office 2002

        ## ACL SUPPORT ##
        nt acl support = yes
        acl compatibility = auto
        acl check permissions = yes
        acl group control = yes

#======================= Share Definitions =======================

[Homes]
        comment = Home Directories
        browseable = yes # Allows browsing a directory tree
        read only = no # Not read-only
        writable = yes # Allows writing
        create mask = 0777 # Permissions for file creation
        directory mask = 0777 # Permissions for directory creation
        veto files = /.DS_Store/.fuse_*/ # Don't display objects: ".DS_Store" and ".fuse_*"
        dos filemode = yes
        inherit acls = yes
        inherit permissions = yes

[netlogon]
        comment = Network Logon Service
        path = /mnt/test
        read only = no
        dos filemode = yes
        inherit acls = yes
        inherit permissions = yes
        browseable = yes
        writable = yes
        create mask = 0777
        directory mask = 0777
        valid users = @"EXAMPLE.TEST+users"
        admin users = @"EXAMPLE.TEST+administrators"

[Sauvegardes]
        comment = Sauvegardes
        path = /saves # Share folder
        browseable = yes
        writable = yes
        veto files = /.DS_Store/.fuse_*/
```

Adapt all this to your configuration. Then restart Samba:

```bash
/etc/init.d/samba restart
```

### Creating the Machine Account for the Samba Server in Active Directory

This will add the Samba machine to AD:

```bash
net ads join -U Administrateur (to exit 'net ads leave')
```

#### Verification of Access via Kerberos to the Shared Resources of the "AD" Server

```bash
kinit Administrateur@EXAMPLE.COM
smbclient -L //AD -k
kdestroy
```

#### Verification of "Mounting" a Shared Resource

```bash
mkdir /mnt/test
mount -t cifs -o username=Administrateur //ad/Public /mnt/test
```

Also see the commands:

```bash
net ads info
```

and

```bash
net ads status -U Administrateur
```

## Unified Authentication UNIX / Windows

The Winbind component of Samba helps solve unified authentication problems. It mainly allows, with the help of PAM (Pluggable Authentication Modules) and NSS (Name Service Switch), to make Windows domain users appear as UNIX accounts.

### Installation

```bash
apt-get install winbind
```

### Configuration

The configuration of Winbind is done in the Samba configuration file. So edit `/etc/samba/smb.conf`:

```bash
        ## Integration Winbind in AD ##
        idmap uid = 10000-20000 # Correspondences of uids between the Linux server and Active Directory
        idmap gid = 10000-20000 # Correspondences of gids between the Linux server and Active Directory
        winbind enum users = yes # List users at Winbind startup
        winbind enum groups = yes # List groups at Winbind startup
        winbind separator = + # Domain/username separation character (ex: DOMAIN+user)
        winbind use default domain = yes # If the domain is not specified, use the default one
        template shell = /bin/bash # Default shell
        template homedir = /home/win2k3/%D/%U # Default home directory
```

### Verification of Winbind Functionality

Adapt all this to your configuration. Then restart Samba and Winbind:

```bash
/etc/init.d/samba restart
/etc/init.d/winbind restart
```

- User account query: `wbinfo -u`
- Group query: `wbinfo -g`

### Adding Winbind Support to NSS

Edit the `/etc/nsswitch.conf` file:

```bash
passwd: compat winbind
group:  compat winbind
```

Adapt these lines to your configuration.

### Verification of NSS+Winbind Functionality

This should display a mix between your local user configuration (`/etc/passwd`), groups (`/etc/group`) and accounts in the AD:

```bash
getent passwd
getent group
```

### Adding Winbind Support to PAM

#### Debian

For all of the following, adapt to your configuration. Edit `/etc/pam.d/common-auth`:

```bash
auth    sufficient      pam_winbind.so
auth    required        pam_unix.so nullok_secure
```

Then edit `/etc/pam.d/common-account`:

```bash
account sufficient      pam_winbind.so
account required        pam_unix.so
```

Now edit `/etc/pam.d/common-session`:

```bash
session required        pam_unix.so
session required        pam_mkhomedir.so skel=/etc/skel/ umask=0077
```

If you don't have these files, it is possible that everything is in `/etc/pam.d/system-auth`.

#### Red-Hat

For all of the following, adapt to your configuration. Edit `/etc/pam.d/login`:

```bash
auth       required     pam_securetty.so
auth       sufficient   pam_winbind.so
auth       sufficient   pam_unix.so use_first_pass
auth       required     pam_stack.so service=system-auth
auth       required     pam_nologin.so
account    sufficient   pam_winbind.so
account    required     pam_stack.so service=system-auth
password   required     pam_stack.so service=system-auth
session    required     pam_stack.so service=system-auth
session    optional     pam_console.so
```

## Verification of PAM+Winbind Authentication

Let's create a folder:

```bash
mkdir -p /home/win2k3/EXAMPLE0
```

It is possible to log in on a console with an account declared in AD:

- 1st time try to authenticate on the AD.
- 2nd password for the local system.

You can also do an ssh, but first, you need to restart the service:

```bash
/etc/init.d/ssh restart
```

```bash
ssh Administrateur@localhost
```

No need to create POSIX or Samba accounts for shares anymore:

```bash
/etc/init.d/samba restart
```

```bash
smbclient -L localhost -U utilisateurAD
```

## Connection

### Windows

To connect from Windows, in a link window, type this:

```
\\IP_of_samba_server\Share_name
```

You will directly access the share

### Unix (Linux/Mac...)

You must have smbfs installed before continuing:

```bash
apt-get install smbfs
```

Then, just create a folder and mount the share in it:

```bash
mkdir saves
mount -t cifs -o username=user,password=password //192.168.0.1/saves ./saves
```

## References

Samba: [https://www.samba.org](https://www.samba.org)  
[ADS Documentation](/pdf/ads.pdf)  
[Samba ADS on CentOS Documentation](/pdf/samba_ads_centos.pdf)
