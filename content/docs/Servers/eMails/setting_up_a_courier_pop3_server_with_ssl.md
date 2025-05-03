---
weight: 999
url: "/Mise_en_place_d'un_serveur_Courier_POP3_avec_SSL/"
title: "Setting up a Courier POP3 server with SSL"
description: "Guide to install and configure a Courier POP3 server with SSL encryption on OpenBSD"
categories: ["Linux", "Servers", "Network"]
date: "2008-04-10T13:58:00+02:00"
lastmod: "2008-04-10T13:58:00+02:00"
tags: ["SSL", "Courier", "POP3", "Mail", "OpenBSD", "Server"]
toc: true
---

## Introduction

If you want to install a POP3 server so that your emails can be retrieved from a mail client or another server, then Courier POP3 is a good solution. Additionally, setting up SSL is recommended as it provides security.

For the POP3 server to be fully functional, you'll obviously need it to be accompanied by a mail server ([see documentation here]({{< ref "docs/Servers/eMails">}})).

I have done this installation on OpenBSD, so this entire documentation will be for OpenBSD, but given the simplicity of the process, it's very easily portable to other systems.

## Installation

Installation of courier-pop3:

```bash
$ pkg_add -iv courier-pop3-4.1.1p0

--- courier-pop3-4.1.1p0 -------------------
You now need to edit appropriately the Courier-POP3 configuration files
installed in /etc/courier/courier/.

To use POP3-SSL, be sure to read ssl(8) and run the mkpop3dcert script
if you require a self-signed certificate.

To control the daemon use /usr/local/libexec/pop3.rc and
/usr/local/libexec/pop3-ssl.rc.
```

You also need to install courier-imap because there are commands needed for the POP3 server to start (which is not very well designed):

```bash
$ pkg_add -iv courier-imap-4.1.1p2

parsing courier-imap-4.1.1p2
Dependencies for courier-imap-4.1.1p2 resolve to: courier-authlib-0.58p3, gdbm-1.8.3p0
found libspec c.41.0 in /usr/lib
found libspec courierauth.0.0 in package courier-authlib-0.58p3
found libspec courierauthsasl.0.0 in package courier-authlib-0.58p3
found libspec crypto.13.0 in /usr/lib
found libspec gdbm.3.0 in package gdbm-1.8.3p0
found libspec ssl.11.0 in /usr/lib
adding group _courier
adding user _courier
installed /etc/courier/imapd-ssl from /usr/local/share/examples/courier/imapd-ssl.dist**************************************************************************** | 99%
installed /etc/courier/imapd.cnf from /usr/local/share/examples/courier/imapd.cnf
installed /etc/courier/imapd from /usr/local/share/examples/courier/imapd.dist*************************************************************************************| 100%
installed /etc/courier/quotawarnmsg from /usr/local/share/examples/courier/quotawarnmsg.example
courier-imap-4.1.1p2: complete

--- courier-imap-4.1.1p2 -------------------
You now need to edit appropriately the Courier-IMAP configuration files
installed in /etc/courier/courier/.

Pay particular attention to the details in imapd.cnf, and read ssl(8) if
necessary. You MUST set the CN in imapd.cnf to the hostname by which
your IMAP server is accessed, or else clients will complain. When this
is done, you can use the 'mkimapdcert' script to automatically generate
a server certificate, which is installed into /etc/ssl/private/imapd.pem

To control the daemon use /usr/local/libexec/imapd.rc and
/usr/local/libexec/imapd-ssl.rc, and to run the authdaemon, place the
following in /etc/rc.local:

mkdir -p /var/run/courier{,-auth}/
/usr/local/sbin/authdaemond start
```

## Configuration

Let's go to the /etc/courier directory:

```bash
cd /etc/courier
```

Now edit the pop3d-ssl file to replace these fields:

```bash
POP3DSSLSTART=YES
POP3_TLS_REQUIRED=1
SSLPIDFILE=/var/run/pop3d-ssl.pid
```

Now we'll edit the pop3d.cnf file to generate a proper file for our server:

```bash
RANDFILE = /usr/local/sbin/pop3d.rand

[ req ]
default_bits = 2048
encrypt_key = yes
distinguished_name = req_dn
x509_extensions = cert_type
prompt = no

[ req_dn ]
C=FR
ST=IDF
L=Paris
O=Company Mail Server
OU=Company Mail Server SSL
CN=niceday
emailAddress=xxx@mycompany.com


[ cert_type ]
nsCertType = server
```

Then we'll generate the SSL key:

```bash
$ mkpop3dcert

Generating a 2048 bit RSA private key
...+++
............................+++
writing new private key to '/etc/ssl/private/pop3d.pem'
-----
512 semi-random bytes loaded
Generating DH parameters, 512 bit long safe prime, generator 2
This is going to take a long time
.............+................................+.
```

The key is therefore located in `/etc/ssl/private/pop3d.pem`.

## Service Control

Now let's start our service now that the configuration is complete:

```bash
/usr/local/libexec/pop3d-ssl.rc start
```
