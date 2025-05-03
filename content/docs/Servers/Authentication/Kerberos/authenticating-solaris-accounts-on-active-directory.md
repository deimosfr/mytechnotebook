---
weight: 999
url: "/Authentification_de_comptes_Solaris_sur_un_Active_Directory/"
title: "Authenticating Solaris Accounts on Active Directory"
description: "Learn how to set up authentication of Solaris accounts using Windows Active Directory and the Kerberos protocol."
categories: ["Solaris", "Windows", "Authentication"]
date: "2007-09-08T11:13:00+02:00"
lastmod: "2007-09-08T11:13:00+02:00"
tags: ["Solaris", "Active Directory", "Kerberos", "Authentication", "Windows"]
toc: true
---

## Introduction

Implementation of authentication on Solaris from an Active Directory (AD).

What this implementation allows for user management on a machine:

- Solaris accounts need to be created with an identifier identical to the AD one, with disk space.
- Password verification is done via AD.

This document is based on a scenario of implementing this type of authentication on Solaris 9.

The concepts described here apply to all UNIX operating systems that support Kerberos version 5 protocol.

**Environment:**

- server_ad.domain.com is the Active Directory server,
- domain.com is the domain managed by server_ad.

## Prerequisites

Prerequisites:

- Kerberos version 5 (in Sun Enterprise Authentication Mechanism (SEAM) 1.0.1 product),
- ensure that DNS is properly configured on the domain that is managed by Active Directory,
- ensure that the date is properly synchronized with the AD server (ntpdate).

## Configuration

Files to configure to allow authentication on the Solaris station via AD are:

- `/etc/pam.conf` to indicate that Kerberos should be used for authentication,
- `/etc/krb5/krb5.conf` for using the KDC (Key Distribution Center) of the AD domain.

### krb5.conf

New configuration file `/etc/krb5/krb5.conf`:

```
# PAM Configuration
# 13/03/2007 - Yann Le Thieis
#
# Authentication
#
other   auth sufficient         pam_krb5.so.1
other   auth sufficient         pam_unix.so.1 try_first_pass
#
# Password
#
other   password sufficient     pam_krb5.so.1
other   password sufficient     pam_unix.so.1
#
# Account
#
other   account optional        pam_krb5.so.1
other   account optional        pam_unix.so.1
#
# Session
#
other   session optional        pam_krb5.so.1
other   session optional        pam_unix.so.1
```

### pam.conf

First, back up the original version of krb5.conf:

```bash
$ cp -p /etc/krb5/krb5.conf /etc/krb5/krb5.conf.old
```

The configuration:

```
# krb5.conf configuration for domain domain.com
# 13/03/2007 by Yann Le Thieis
#
[libdefaults]
        default_realm = DOMAIN.COM
        verify_ap_req_nofail = false

[realms]
        domain.com = {
                kdc = server_ad.domain.com:88
                admin_server = server_ad.domain.com:749
                default_domain = domain.com
        }

[domain_realm]
        .domain.com = DOMAIN.COM
        domain.com = DOMAIN.COM

[logging]
        default = FILE:/var/krb5/kdc.log
        kdc = FILE:/var/krb5/kdc.log
        kdc_rotate = {

# How often to rotate kdc.log. Logs will get rotated no more
# often than the period, and less often if the KDC is not used
# frequently.

                period = 1d

# how many versions of kdc.log to keep around (kdc.log.0, kdc.log.1, ...)

                versions = 10
        }

[appdefaults]
        kinit = {
                renewable = true
                forwardable= true
        }
        gkadmin = {
                help_url = http://docs.sun.com:80/ab2/coll.384.1/SEAM/@AB2PageView/1195
        }
```

The line "verify_ap_req_nofail = false" is extremely important if the file /etc/krb5/krb5.keytab is not filled with a line for your domain (i.e., a key that validates the KDC, see the man krb5.conf manual).

## Testing this configuration

The AD account used for the test is ylethieis, which does not exist locally on the Solaris machine. But first, let's try with a dummy account that doesn't exist anywhere.

```bash
$ kinit bidon
Password for bidon@domain.com:
kinit: Client not found in Kerberos database while getting initial credentials
```

Note: kinit – obtain and cache Kerberos ticket-granting ticket.

Try with the ylethieis account but entering a wrong password:

```bash
$ kinit ylethieis
Password for ylethieis@domain.com:
kinit: Pre-authentication failed while getting initial credentials
```

Try with the ylethieis account and the correct password for AD:

```bash
$ kinit ylethieis
Password for ylethieis@domain.com:
```

The Kerberos client service on the Solaris machine correctly queries the AD.

Cached tickets:

```bash
$ klist
Ticket cache: /tmp/krb5cc_0
Default principal: ylethieis@domain.com

Valid starting                   Expires                  Service principal
Tuesday, March 13, 2007, 11:03:14 GMT  Tuesday, March 13, 2007, 21:03:14 GMT  krbtgt/domain.com@domain.com
       renewable until Tuesday, March 13, 2007, 21:03:14 GMT
```

## Creating an AD account environment on the Solaris machine

Create the space for accounts authenticated via AD, and an ad user group to distinguish them from others (not mandatory!):

```bash
$ mkdir /export/home/ad
$ groupadd ad
```

Add an ylethieis account:

```bash
$ useradd -g ad -m -d /export/home/ylethieis ylethieis
UX: useradd: ylethieis name too long.
64 blocks
```

This account has the ad group as its primary group.

The login name is indicated as too long but the account was successfully created!

At this stage, on the Solaris system, the ylethieis user:

- has no password,
- has ad as its primary group,
- is just listed in /etc/passwd.

Login to the Solaris system with the ylethieis account:

```bash
$ telnet server_solaris
Trying 192.168.0.120...
Connected to 192.168.0.120.
Escape character is '^]'.

SunOS 5.9

login: ylethieis
Enter Kerberos password for ylethieis:
Last login: Tue Mar 13 14:27:52 from yuluth
Sun Microsystems Inc.   SunOS 5.9       Generic January 2003
$
```

We can see that authentication via Active Directory has succeeded for the ylethieis account.
