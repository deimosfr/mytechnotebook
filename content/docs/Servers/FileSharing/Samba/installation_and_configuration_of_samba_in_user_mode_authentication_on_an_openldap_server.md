---
weight: 999
url: '/Installation_et_configuration_de_Samba_en_mode_"User"_(Authentification_sur_un_serveur_OpenLDAP)/'
title: "Installation and Configuration of Samba in User Mode (Authentication on an OpenLDAP Server)"
description: "This guide explains how to install and configure Samba in user mode with authentication on an OpenLDAP server for secure file sharing between Linux and Windows environments."
categories: ["Linux", "Security"]
date: "2008-12-02T07:31:00+02:00"
lastmod: "2008-12-02T07:31:00+02:00"
tags:
  [
    "Samba",
    "OpenLDAP",
    "Security",
    "Authentication",
    "Network Sharing",
    "Linux",
  ]
toc: true
---

## Introduction

Samba is free software under the GPL license that supports the SMB/CIFS protocol. This protocol is used by Microsoft for sharing various resources (files, printers, etc.) between computers running Windows. Samba allows Unix systems to access resources on these systems and vice versa.

Previously, PCs running DOS and early Windows versions sometimes had to install a TCP/IP stack and a set of Unix-origin software: an NFS client, FTP, telnet, lpr, etc. This was cumbersome and penalizing for the PCs of that time, and it also forced their users to develop a double set of habits, adding Unix habits to those of Windows. Samba therefore adopts the opposite approach.

Its name comes from the file and print sharing protocol from IBM, reused by Microsoft, called SMB (Server message block), to which the two vowels 'a' were added: "SaMBa".

Samba was originally developed by Andrew Tridgell starting in 1991, and today receives contributions from about twenty developers from around the world under his coordination. He gave it this name by choosing a name similar to SMB by querying a Unix dictionary with the grep command: `grep "^s.*m.*b" /usr/dict/words`

When both file sharing systems (NFS, Samba) are installed for comparison, Samba proves less efficient than NFS in terms of transfer rates.

Nevertheless, a study has shown that Samba 3 was up to 2.5 times faster than the SMB implementation of Windows Server 2003. See information on LinuxFr.

However, Samba is not compatible with IPv6.

**The "ADS" mode allows you to use an OpenLDAP server to authenticate users for accessing Samba shares. This solution is quite complex to implement but ensures increased security (authentication via Kerberos) for your server.**

**Notes: It is imperative to have implemented [NT-type ACLs]({{< ref "docs/Linux/FilesystemsAndStorage/acl-implementing-nt-type-permissions-on-linux.md" >}}) before continuing.**

To help you with this documentation, you can also consult: [Installation and Configuration of Samba in ADS Mode (Authentication on an AD Server)]({{< ref "docs/Servers/FileSharing/Samba/installation_and_configuration_of_samba_in_ads_mode_authentication_on_an_ad_server.md" >}})

## Information

To access shared resources (SMB/CIFS) of the Windows domain, a Samba installation is required on the Linux server that will primarily play the role of an SMB/CIFS client.  
Samba will also allow, using MSRPC commands, to communicate with the LDAP server to perform various operations: adding information about the Linux server in the directory, listing user accounts/groups, transmitting authentication requests...

The server will run in standalone mode, with user or share security mode. With this mode, users are required to authenticate when they access a share. Samba can authenticate and manage file permissions using accounts and passwords stored in the LDAP directory.

## Samba

### Installation

To install Samba:

```bash
apt-get install samba smbclient
```

### Configuration

To configure Samba, edit the `/etc/samba/smb.conf` file:

#### Samba

```bash
#======================= Global Settings =======================

[global]
   workgroup = mydomain
   server string = %h server
   dns proxy = no

#### Debugging/Accounting ####

   log file = /var/log/samba/log.%m
   max log size = 1000
   syslog = 0
   panic action = /usr/share/samba/panic-action %d

####### Authentication #######

   security = user

   encrypt passwords = true
   obey pam restrictions = yes

   guest account = Invite
   invalid users = root

   unix password sync = yes
   passwd program = /usr/bin/passwd %u
   passwd chat = *Enter\snew\sUNIX\spassword:* %n\n *Retype\snew\sUNIX\spassword:* %n\n *password\supdated\ssuccessfully* .

   pam password change = yes

   passdb backend = ldapsam:ldap://myLDAP_IP_Server/

   ldap suffix = dc=openldap,dc=mydomain,dc=local
   ldap machine suffix = ou=Computers
   ldap idmap suffix =
   ldap user suffix =
   ldap group suffix =
   ldap admin dn = "cn=admin,dc=openldap,dc=mydomain,dc=local"
   ldap passwd sync = yes
   ldap ssl = no
   idmap backend = ldap:ldap://ldap.mydomain.local
   idmap uid = 40000-60000
   idmap gid = 40000-60000

   winbind use default domain = Yes
   winbind trusted domains only = yes

############ Misc ############

   socket options = TCP_NODELAY
   unix extensions = yes
   case sensitive = yes
   delete readonly = yes
   ea support = yes

   ### ACL SUPPORT ###
   nt acl support = yes
   acl compatibility = auto
   acl check permissions = yes
   acl group control = yes

#======================= Share Definitions =======================


[Homes]
        comment = Home Directories
        browseable = yes
        read only = no
        writable = yes
        create mask = 0700
        directory mask = 0700
        dos filemode = yes
        inherit acls = yes
        inherit permissions = yes
```

Adapt all this for your configuration. Then restart Samba:

```bash
/etc/init.d/samba restart
```

## Integration

### Samba

Next, we need to set the LDAP password otherwise we will get an error like:

```
ldap_connect_system: Failed to retrieve password from secrets.tdb
```

Let's use the smbpasswd command to create the `/var/lib/samba/secrets.tdb` file:

```bash
$ smbpasswd -w 'mypassowd'
Setting stored password for "cn=admin,dc=openldap,dc=mydomain,dc=local" in secrets.tdb
```

You can verify using the smbpasswd command. If it asks you for a password, enter it, that means it's working. Otherwise you can try deleting `/var/lib/samba/secrets.tdb` and retry the process.

### OpenLDAP

If you are in LDAP mode, you need to get your SSID on the Samba server:

```bash
net getlocalsid
SID for domain MOONLIGHT is: S-1-5-21-2096052081-3008433157-1548381139-1339
```

If you cannot connect after that, check the logs, you should get something like this:

```
User pmavro with invalid SID S-1-5-21-2096052081-3008433157-1548381139-1319 in passdb
```

You should then check on the OpenLDAP server how it sees the Samba server (which has normally self-registered). You may then have the same bug as me, which is a different sambaSID (notably truncated). Two solutions are available:

- **Solution 1:**

Modify your sambaSID on your machine so that it is the same as on the server:

```bash
net setlocalsid S-1-5-21-2096052081-3008433157-1548381139
```

- **Solution 2:**

Modify on the OpenLDAP server the sambaSID field of the samba server (field 'sambaDomainName') so that it exactly matches the sambaSID via the "net getlocalsid" command.

### Active Directory

We will need Winbind for this to work:

```bash
apt-get install winbind
```

And register your machine in the domain:

```bash
wbinfo --set-auth-user='Administrator%HisPassword'
```

**WARNING: The SID of the machine running samba must be equal to the SID of the domain.** To check, run the commands below:

```bash
net getlocalsid
net getdomainsid
```

They must return the same value!

If this is not the case, copy the value returned by 'net getdomainsid' and do:

```bash
net setlocalsid new_SID_value
```

Check again and if the values are still different, directly modify the SID of the machine in the LDAP directory.

## Resources

- [Samba topics]({{< ref "docs/Servers/FileSharing/Samba/" >}})
- [Documentation on Samba and OpenLDAP installation](/pdf/openldap_samba_domain_controller_ubuntu.pdf)
