---
weight: 999
url: "/Mise_en_place_du_client_OpenLDAP/"
title: "Setting up OpenLDAP Client"
description: "Step-by-step guide for configuring an OpenLDAP client on Solaris systems."
categories: ["Database", "Linux", "Security"]
date: "2009-09-23T09:00:00+02:00"
lastmod: "2009-09-23T09:00:00+02:00"
tags: ["LDAP", "OpenLDAP", "Solaris", "Authentication", "PAM", "FAQ"]
toc: true
---

## Introduction

OpenLDAP is unfortunately not available as standard on Solaris, however it is possible to install it via pkg-get.

## Configuration

* Configure the `/etc/pam.conf` file:

For each line:

```bash
service  auth required           pam_unix_auth.so.1
```

replace "required" with "sufficient" and add behind the line:

```bash
service auth sufficient pam_ldap.so.1 try_first_pass
```

This should result in something like this (`/etc/pam.conf`):

```bash
#
#ident  "@(#)pam.conf   1.28    04/04/21 SMI"
#
# Copyright 2004 Sun Microsystems, Inc.  All rights reserved.
# Use is subject to license terms.
#
# PAM configuration
#
# Unless explicitly defined, all services use the modules
# defined in the "other" section.
#
# Modules are defined with relative pathnames, i.e., they are
# relative to /usr/lib/security/$ISA. Absolute path names, as
# present in this file in previous releases are still acceptable.
#
# Authentication management
#
# login service (explicit because of pam_dial_auth)
#
login   auth requisite          pam_authtok_get.so.1
login   auth required           pam_dhkeys.so.1
login   auth required           pam_unix_cred.so.1
login   auth sufficient         pam_ldap.so.1 try_first_pass
login   auth sufficient         pam_unix_auth.so.1
login   auth required           pam_dial_auth.so.1
#
# rlogin service (explicit because of pam_rhost_auth)
#
rlogin  auth sufficient         pam_rhosts_auth.so.1
rlogin  auth requisite          pam_authtok_get.so.1
rlogin  auth required           pam_dhkeys.so.1
rlogin  auth required           pam_unix_cred.so.1
rlogin  auth sufficient         pam_ldap.so.1 try_first_pass
rlogin  auth sufficient         pam_unix_auth.so.1
#
# Kerberized rlogin service
#
krlogin auth required           pam_unix_cred.so.1
krlogin auth binding            pam_krb5.so.1
krlogin auth sufficient         pam_ldap.so.1
krlogin auth sufficient         pam_unix_auth.so.1
#
# rsh service (explicit because of pam_rhost_auth,
# and pam_unix_auth for meaningful pam_setcred)
#
rsh     auth sufficient         pam_rhosts_auth.so.1
rsh     auth required           pam_unix_cred.so.1
#
# Kerberized rsh service
#
krsh    auth required           pam_unix_cred.so.1
krsh    auth binding            pam_krb5.so.1
krsh    auth sufficient         pam_ldap.so.1
krsh    auth sufficient         pam_unix_auth.so.1
#
# Kerberized telnet service
#
ktelnet auth required           pam_unix_cred.so.1
ktelnet auth binding            pam_krb5.so.1
ktelnet auth sufficient         pam_ldap.so.1
ktelnet auth sufficient         pam_unix_auth.so.1
#
# PPP service (explicit because of pam_dial_auth)
#
ppp     auth requisite          pam_authtok_get.so.1
ppp     auth required           pam_dhkeys.so.1
ppp     auth required           pam_unix_cred.so.1
ppp     auth sufficient         pam_ldap.so.1
ppp     auth sufficient         pam_unix_auth.so.1
ppp     auth required           pam_dial_auth.so.1
#
# Default definitions for Authentication management
# Used when service name is not explicitly mentioned for authentication
#
other   auth requisite          pam_authtok_get.so.1
other   auth required           pam_dhkeys.so.1
other   auth required           pam_unix_cred.so.1
other   auth sufficient         pam_ldap.so.1
other   auth sufficient         pam_unix_auth.so.1
#
# passwd command (explicit because of a different authentication module)
#
passwd  auth required           pam_passwd_auth.so.1
#
# cron service (explicit because of non-usage of pam_roles.so.1)
#
cron    account required        pam_unix_account.so.1
#
# Default definition for Account management
# Used when service name is not explicitly mentioned for account management
#
other   account requisite       pam_roles.so.1
other   account required        pam_unix_account.so.1
#
# Default definition for Session management
# Used when service name is not explicitly mentioned for session management
#
other   session required        pam_unix_session.so.1
#
# Default definition for  Password management
# Used when service name is not explicitly mentioned for password management
#
other   password required       pam_dhkeys.so.1
other   password requisite      pam_authtok_get.so.1
other   password requisite      pam_authtok_check.so.1
other   password required       pam_authtok_store.so.1
#
# Support for Kerberos V5 authentication and example configurations can
# be found in the pam_krb5(5) man page under the "EXAMPLES" section.
#
```

* Configure the `/etc/nsswitch.ldap` file

Keep "ldap" only where it's useful: for now on the passwd: and group: lines.
For the rest, use the content of the `/etc/nsswitch.dns` file.
This gives:

```bash
#
# Copyright 2006 Sun Microsystems, Inc.  All rights reserved.
# Use is subject to license terms.
#
 
#
# /etc/nsswitch.dns:
#
# An example file that could be copied over to /etc/nsswitch.conf; it uses
# DNS for hosts lookups, otherwise it does not use any other naming service.
#
# "hosts:" and "services:" in this file are used only if the
# /etc/netconfig file has a "-" for nametoaddr_libs of "inet" transports.
 
# DNS service expects that an instance of svc:/network/dns/client be
# enabled and online.
 
passwd:     files ldap
group:      files ldap
 
# You must also set up the /etc/resolv.conf file for DNS name
# server lookup.  See resolv.conf(4).
hosts:      files dns 
 
# Note that IPv4 addresses are searched for in all of the ipnodes databases
# before searching the hosts databases.
ipnodes:   files dns 
 
networks:   files
protocols:  files
rpc:        files
ethers:     files
netmasks:   files
bootparams: files
publickey:  files
# At present there isn't a 'files' backend for netgroup;  the system will 
#   figure it out pretty quickly, and won't use netgroups at all.
netgroup:   files
automount:  files
aliases:    files
services:   files
printers:       user files
 
auth_attr:  files
prof_attr:  files
project:    files
 
tnrhtp:     files
tnrhdb:     files
```

Once that is done, we can proceed with the configuration. **Note: if you are in a cluster environment, adapt to your initial configuration**:

```bash
cp /etc/nsswitch.ldap /etc/nsswitch.conf
```

* Launch the LDAP client configuration

Simply type the command:

```bash
ldapclient manual -v -a authenticationMethod=simple -a proxyDN=cn=admin,dc=openldap,dc=mydomain,dc=local -aproxyPassword=bidon -a defaultSearchBase=dc=openldap,dc=mydomain,dc=local -a defaultServerList=ldap.mydomain.local -a serviceSearchDescriptor=passwd:dc=openldap,dc=mydomain,dc=local?sub -a serviceSearchDescriptor=shadow:dc=openldap,dc=mydomain,dc=local?sub -a serviceSearchDescriptor=group:dc=openldap,dc=mydomain,dc=local?sub -a serviceAuthenticationMethod=pam_ldap:simple=
```

Note: it seems that the ldapclient command is bugged and requires the proxyDN and proxyPassword parameters even if they are unused (and even if they contain anything)!

* Be careful with the home directory, you need to configure `/etc/auto_home` (http://www.solaris-fr.org/home/docs/base/utilisateurs). For me, this gives:

```perl
#
# Copyright 2003 Sun Microsystems, Inc.  All rights reserved.
# Use is subject to license terms.
#
# ident "@(#)auto_home  1.6     03/04/28 SMI"
#
# Home directory map for automounter
#
+auto_home
* localhost:/export/home/&
```

In case you want to automatically create the home directory, you need to port the pam_mkhomedir module from Linux:

* http://mega.ist.utl.pt/~filipe/pam_mkhomedir-sol/?C=D;O=A
* http://www.keutel.de/pam_mkhomedir/index.html

A good idea would also be to automatically mount the home directory from an NFS server.

User accounts in the LDAP directory need to have the "shadowAccount" objectClass in their objectClass list to be recognized by Solaris.

## FAQ

### WARNING: /var/ldap/ldap_client_file is missing or not readable

I got this message when I tried to initialize the service with the ldapclient command and the LDAP servers were down. Check that:

* The client service (svc:/network/ldap/client:default) is offline
* The LDAP servers are available

If it still doesn't work, you'll need to create 2 files manually:

```bash
#
# Do not edit this file manually; your changes will be lost.Please use ldapclient (1M) instead.
#
NS_LDAP_FILE_VERSION= 2.0
NS_LDAP_SERVERS= server2, server1
NS_LDAP_SEARCH_BASEDN= dc=openldap,dc=mycompany,dc=com
NS_LDAP_AUTH= simple
NS_LDAP_CACHETTL= 0
NS_LDAP_SERVICE_SEARCH_DESC= passwd:dc=openldap,dc=mycompany,dc=com?sub?&(&(objectClass=posixAccount)(!(objectClass=computer)))
NS_LDAP_SERVICE_SEARCH_DESC= shadow:dc=openldap,dc=mycompany,dc=com?sub
NS_LDAP_SERVICE_SEARCH_DESC= group:dc=openldap,dc=mycompany,dc=com?sub
NS_LDAP_SERVICE_AUTH_METHOD= pam_ldap:simple
```

And the file containing the credentials:

```bash
#
# Do not edit this file manually; your changes will be lost.Please use ldapclient (1M) instead.
#
NS_LDAP_BINDDN= cn=admin,dc=openldap,dc=mycompany,dc=lan
NS_LDAP_BINDPASSWD= {NS1}4a3788e8c053424f
```

To generate a password, use the ldap_gen_profil command:

```bash
ldap_gen_profile -P profile -b dc=mycompany,dc=local -w test 127.0.0.1 | grep SolarisBindPassword
```

## Resources
- http://docs.sun.com/app/docs/doc/806-4077/6jd6blbev?a=view#ldapsecure-74
- http://www.ypass.net/solaris/openldap/ldapcachemgr.php
