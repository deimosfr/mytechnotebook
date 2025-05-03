---
weight: 999
url: "/Dkfilter_\\:_Proxy_SMTP_(signatures_mails)_made_by_Yahoo/"
title: "Dkfilter: Proxy SMTP (Signature Mails) Made by Yahoo"
description: "A guide to set up DomainKeys for Postfix using dkfilter, a SMTP proxy for email signature verification and signing."
categories: ["Debian", "Red Hat", "Linux"]
date: "2008-01-31T07:06:00+02:00"
lastmod: "2008-01-31T07:06:00+02:00"
tags: ["Servers", "Mail", "SMTP", "Security", "Postfix"]
toc: true
---

## Introduction

DomainKeys is an anti-spam software application in development at [Yahoo](https://antispam.yahoo.com/domainkeys) that uses a form of public key cryptography to authenticate the sender's domain. **dkfilter is an SMTP-proxy designed for Postfix**. It implements DomainKeys message signing and verification. It comprises two separate filters, an "outbound" filter for **signing outgoing email on port 587**, and an "inbound" filter for **verifying signatures of incoming email on port 25**. This document is to describe step by step how to install dkfilter for postfix to deploy domainkeys signing and verification.

```

              +------+
              |verify|            (verify)
              +--+---+          SpamAssassin
                 ^                   ^v
incoming:        |              +----++-----+
  MX ---->  25 smtpd ---> 10024 >           >--> 10025 smtpd -->
submission:                     |           |
  SASL -->  25 smtpd \          |  amavisd  |
                      +->       |           |
  mynets->  25 smtpd ---> 10026 >ORIGINATING>--> 10027 smtpd -->
       --> 587 smtpd --->       +-----------+            |
               (convert to 7-bit)                        v
                                                       +----+
                                                       |sign|
                                                       +----+
```

## Installation

### Postfix

Install postfix for your domain to send and receive mails with this doc:

[Installation and Configuration of Postfix and Courier]({{< ref "docs/Servers/eMails/Postfix/installation_and_configuration_of_postfix_and_courier.md">}})

### CPAN Perl modules

Dkfilter is written in Perl. It requires the following Perl Modules from CPAN archive.

- Crypt::OpenSSL::RSA
- Mail::Address
- MIME::Base64
- Net::DNS
- Test::More
- Text::Wrap
- Mail::DomainKeys

Following commands would help:

```bash
perl -MCPAN -e'CPAN::Shell->install("Crypt::OpenSSL::RSA")'
perl -MCPAN -e'CPAN::Shell->install("Mail::Address")'
perl -MCPAN -e'CPAN::Shell->install("MIME::Base64")'
perl -MCPAN -e'CPAN::Shell->install("Net::DNS")'
perl -MCPAN -e'CPAN::Shell->install("Test::More")'
perl -MCPAN -e'CPAN::Shell->install("Text::Wrap")'
perl -MCPAN -e'CPAN::Shell->install("Email::Address")'
perl -MCPAN -e'CPAN::Shell->install("Mail::DomainKeys")'
```

_Note:_ Also resolve any dependent Perl Module required in installing the above Perl modules.

### Dkfilter

The following steps are recommended for installing dkfilter.

#### Download

Download dkfilter from following URL:

```bash
wget http://jason.long.name/dkfilter/dkfilter-0.11.tar.gz
```

#### Installation

Installing dkfilter:

```bash
tar xvf dkfilter-0.11.tar.gz
cd dkfilter-0.11
./configure --prefix=/usr/local/dkfilter
make install
useradd dkfilter
```

The filter scripts will be installed in _/usr/local/dkfilter/bin_ and the Perl module files will be in _/usr/local/dkfilter/lib_.

## Configuration

### Inbound Filter

We need to make relevant changes inside Postfix configuration files to check incoming mails for the signature.

Edit the Postfix master.cf file:

```bash
vi /etc/postfix/master.cf
```

Add this configuration:

```
# Dkfilter Inbound Part1
smtp      inet  n       -       n       -       -       smtpd
        -o smtpd_proxy_filter=127.0.0.1:10029
        -o smtpd_client_connection_count_limit=10

### Amavis filter ###
amavis unix - - n - 2 smtp
        -o smtp_data_done_timeout=1200
        -o disable_dns_lookups=yes

127.0.0.1:10025 inet n - n - - smtpd
        -o content_filter=
        -o local_recipient_maps=
        -o relay_recipient_maps=
        -o smtpd_restriction_classes=
        -o smtpd_client_restrictions=
        -o smtpd_helo_restrictions=
        -o smtpd_sender_restrictions=
        -o smtpd_recipient_restrictions=permit_mynetworks,reject
        -o mynetworks=127.0.0.0/8
        -o strict_rfc821_envelopes=yes
### End of Amavis Filter ###

# Dkfilter Inbound Part2
127.0.0.1:10026 inet n  -       n       -        -      smtpd
        -o smtpd_authorized_xforward_hosts=127.0.0.0/8
        -o smtpd_client_restrictions=
        -o smtpd_helo_restrictions=
        -o smtpd_sender_restrictions=
        -o smtpd_recipient_restrictions=permit_mynetworks,reject
        -o smtpd_data_restrictions=
        -o mynetworks=127.0.0.0/8
        -o receive_override_options=no_unknown_recipient_checks
```

Insert above lines at the end of the file. Here we define that mail will be received after smtp for verification on 127.0.0.1 with port 10026. You can define your own desired IP address on which you want to listen for signature checking.

### Outbound filter

The outbound filter needs access to the private key used for signing messages. In addition, it needs to know the name of the key selector being used, and what domain it should sign messages for. This information is specified with command-line arguments to dkfilter.out.

#### Key pair

- Generate a private/public key pair and publish the public key in DNS:

```bash
cd /usr/local/dkfilter
openssl genrsa -out private.key 1024
openssl rsa -in private.key -pubout -out public.key
```

This creates the files private.key and public.key in the current directory, containing the private key and public key. Make sure private.key is not world-readable, but still readable by the dkfilter user.

- Pick a selector name... e.g. m1

#### Bind

- Put the public-key data in DNS, in your domain, using the selector name you picked. Copy the contents of the public.key file and remove the PEM header and footer, and paste it in dns zone file by creating a TXT entry, like this:

```
_deimos.fr    TXT     "t=y; o=-;"
m1._deimos.fr TXT     "g=; k=rsa; p=MIGxm...MeQIDAQAB;"
```

where m1 is the name of the selector chosen in the last step and the p= parameter contains the public-key as one long string of characters.

#### Postfix

Finally, configure Postfix to filter outgoing, authorized messages only through the dkfilter.out service on port 10027. In the following example, messages sent via SMTP on port 587 (the submission port) will go through an After-Queue content filter that signs messages with DomainKeys.

```bash
vi /etc/postfix/master.cf
```

Add this configuration:

```
submission  inet  n     -       n       -       -       smtpd
        -o smtpd_etrn_restrictions=reject
        -o smtpd_sasl_auth_enable=yes
        -o content_filter=dksign:[127.0.0.1]:10027
        -o receive_override_options=no_address_mappings
        -o smtpd_recipient_restrictions=permit_mynetworks,permit_sasl_authenticated,reject

dksign    unix  -       -       n       -       10      smtp
        -o smtp_send_xforward_command=yes
        -o smtp_discard_ehlo_keywords=8bitmime

127.0.0.1:10028 inet  n  -      n       -       10      smtpd
        -o content_filter=
        -o receive_override_options=no_unknown_recipient_checks,no_header_body_checks
        -o smtpd_helo_restrictions=
        -o smtpd_client_restrictions=
        -o smtpd_sender_restrictions=
        -o smtpd_recipient_restrictions=permit_mynetworks,reject
        -o mynetworks=127.0.0.0/8
        -o smtpd_authorized_xforward_hosts=127.0.0.0/8
```

Execute postfix reload for Postfix to respond to changes in /etc/postfix/master.cf:

```bash
postfix reload
```

## Start / Stop script

This is the startup/shutdown script:

```bash
#!/bin/sh
#
# Copyright (c) 2005 Messiah College.
# Modified by Pierre Mavro

DKFILTERUSER=dkfilter
DKFILTERGROUP=dkfilter
DKFILTERDIR=/usr/local/dkfilter

HOSTNAME=`hostname -f`
DOMAIN=`hostname -d`
DKFILTER_IN_ARGS="--hostname=$HOSTNAME 127.0.0.1:10029 127.0.0.1:10026"
DKFILTER_OUT_ARGS="--keyfile=$DKFILTERDIR/private.key --selector=postfix --domain=$DOMAIN --method=nofws --headers 127.0.0.1:10027 127.0.0.1:10028"

DKFILTER_IN_BIN="$DKFILTERDIR/bin/dkfilter.in"
DKFILTER_OUT_BIN="$DKFILTERDIR/bin/dkfilter.out"
PIDDKFILTER_IN="/var/run/dkfilter.in"
PIDDKFILTER_OUT="/var/run/dkfilter.out"

case "$1" in
       start)
               echo -n "Starting inbound DomainKeys-filter (dkfilter.in)..."
                start-stop-daemon --start -q -p$PIDDKFILTER_IN -u $DKFILTERUSER -g $DKFILTERGROUP -x `$DKFILTER_IN_BIN $DKFILTER_IN_ARGS` &
               RETVAL=$?
               if [ $RETVAL -eq 0 ]; then
                       echo done.
               else
                       echo failed.
                       exit $RETVAL
               fi
               echo -n "Starting outbound DomainKeys-filter (dkfilter.out)..."
                start-stop-daemon --start -q -p $PIDDKFILTER_OUT -u $DKFILTERUSER -g $DKFILTERGROUP -x `$DKFILTER_OUT_BIN $DKFILTER_OUT_ARGS` &
               RETVAL=$?
               if [ $RETVAL -eq 0 ]; then
                       echo done.
               else
                       echo failed.
                       exit $RETVAL
               fi
               ;;
       stop)
               echo -n "Shutting down inbound DomainKeys-filter (dkfilter.in)..."
                start-stop-daemon --stop -p $PIDDKFILTER_IN
               RETVAL=$?
               if [ $RETVAL -eq 0 ]; then
                       echo done.
               else
                       echo failed.
               fi
               echo -n "Shutting down outbound DomainKeys-filter (dkfilter.out)..."
               start-stop-daemon --stop -p $PIDDKFILTER_OUT
               RETVAL=$?
               if [ $RETVAL -eq 0 ]; then
                       echo done.
               else
                       echo failed.
                      exit $RETVAL
               fi
               ;;
       restart)
               $0 stop
               $0 start
               ;;
       *)
               echo "Usage: $0 {start|stop|restart}"
               exit 1
               ;;
esac
```

Copy that script in _/etc/rc.d/init.d (Red Hat System)_ or _/etc/init.d (Debian System)_ and edit it as per your requirement.

Then add the proper permissions:

```bash
chmod 755 /etc/init.d/dkfilter
```

You can also use this command on Debian:

```bash
update-rc.d dkfilter defaults
```

## Resources
- [Set Up Postfix DKIM With dkim-milter](/pdf/set_up_postfix_dkim_with_dkim-milter.pdf)
