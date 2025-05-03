---
weight: 999
url: "/SASL_\\:_Envoie_de_mails_à_distance_sécurisé_avec_son_serveur_Postfix_(SASL+TLS)/"
title: "SASL: Secure Remote Email Sending with Postfix Server (SASL+TLS)"
description: "How to configure SASL and TLS with Postfix for secure remote email sending"
categories: ["Linux", "Servers", "Network"]
date: "2012-06-06T08:27:00+02:00"
lastmod: "2012-06-06T08:27:00+02:00"
tags:
  [
    "Postfix",
    "Email",
    "Security",
    "SASL",
    "TLS",
    "Linux",
    "Server Configuration",
  ]
toc: true
---

![Postfix](/images/postfix_logo.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 2.7.1 |
| **Operating System** | Debian 6 |
| **Website** | [Postfix Website](https://www.postfix.org/) |
| **Last Update** | 06/06/2012 |
{{< /table >}}

## Introduction

For several years now, most mail servers no longer authorize sending emails from outside to prevent becoming an open relay. It's therefore necessary to add an additional layer called [SASL](https://en.wikipedia.org/wiki/SASL).

That's what we'll try to implement here with added security (TLS) :-)

## Installation

To install SASL, you must already have [installed and configured Postfix]({{< ref "docs/Servers/eMails/Postfix/installation_and_configuration_of_postfix_and_courier.md">}}). Once that's done, install SASL:

```bash
aptitude install postfix-tls sasl2-bin libsasl2-2 libsasl2-modules openssl
```

## Configuration

### Preparation of directories

We'll create the necessary directories for our chrooted Postfix:

```bash
mkdir -p /var/spool/postfix/var/run
mv /var/run/saslauthd /var/spool/postfix/var/run/
ln -s /var/spool/postfix/var/run/saslauthd /var/run
chgrp sasl /var/spool/postfix/var/run/saslauthd
adduser postfix sasl
mkdir -p /etc/postfix/ssl
```

### Certificate creation

Create the /etc/postfix/ssl folder:

```bash
cd /etc/postfix/ssl
openssl genrsa -des3 -rand /dev/urandom -out smtpd.key 1024
chmod 600 smtpd.key
openssl req -new -key smtpd.key -out smtpd.csr
openssl x509 -req -days 3650 -in smtpd.csr -signkey smtpd.key -out smtpd.crt
openssl rsa -in smtpd.key -out smtpd.key.unencrypted
mv -f smtpd.key.unencrypted smtpd.key
openssl req -new -x509 -extensions v3_ca -keyout cakey.pem -out cacert.pem -days 3650
```

### main.cf

In your Postfix configuration file, you'll need to add these lines. Edit the file `/etc/postfix/main.cf`:

```bash
# SASL
smtp_sasl_auth_enable = yes
smtpd_sasl_auth_enable = yes
smtp_sasl_security_options =
# To fix the bug with some clients (Outlook...)
broken_sasl_auth_clients = yes
smtpd_sasl_application_name= smtpd
smtpd_sasl_security_options = noanonymous
smtpd_sasl_local_domain =
smtp_sasl_password_maps = hash:/etc/postfix/sasl_passwd

# TLS/SSL
smtpd_tls_security_level = may
smtpd_tls_auth_only = no
smtp_use_tls = yes
smtpd_use_tls = yes
smtp_tls_loglevel = 1
smtpd_tls_received_header = yes
smtpd_tls_loglevel = 1
smtp_tls_note_starttls_offer = yes
smtpd_tls_key_file = /etc/postfix/ssl/smtpd.key
smtpd_tls_cert_file = /etc/postfix/ssl/smtpd.crt
smtpd_tls_CAfile = /etc/postfix/ssl/cacert.pem
```

This should give you something like this:

```bash
# See /usr/share/postfix/main.cf.dist for a commented, more complete version

# Security
smtpd_banner = fire.deimos.fr - Microsoft Exchange (5.5)
biff = no
disable_vrfy_command = yes
smtpd_helo_required = yes

# Reject unknow domain
reject_unknown_recipient_domain = yes

# appending .domain is the MUA's job.
append_dot_mydomain = no

mydomain = deimos.fr
myhostname = fire.deimos.fr
alias_maps = hash:/etc/aliases
alias_database = hash:/etc/aliases
myorigin = /etc/mailname
mydestination = deimos.fr
relayhost =
mynetworks = 127.0.0.0/8, 192.168.0.0/24, 10.8.0.0/24
home_mailbox = Maildir/
mailbox_size_limit = 0
recipient_delimiter = +
inet_interfaces = all
mailbox_command = procmail -a "$EXTENSION"

# Virtual Domains
virtual_alias_domains = deimos.fr
virtual_alias_maps = hash:/etc/postfix/virtual

# Protection against Open Relay
smtpd_client_restrictions = reject_rbl_client bl.spamcop.net

# Protection against Spam
smtpd_recipient_restrictions =  permit_mynetworks,
                                permit_sasl_authenticated,
                                reject_unauth_destination,
                                reject_invalid_hostname,
                                reject_non_fqdn_sender,
                                reject_unknown_sender_domain,
                                reject_non_fqdn_recipient,
                                reject_unknown_recipient_domain,
                                reject_rhsbl_client blackhole.securitysage.com,
                                reject_rhsbl_sender blackhole.securitysage.com,
                                reject_rbl_client list.dsbl.org,
                                reject_rbl_client cbl.abuseat.org,
                                reject_rbl_client dul.dnsbl.sorbs.net,
                                reject_rbl_client multi.uribl.com,
                                reject_rbl_client dsn.rfc-ignorant.org,
                                reject_rbl_client dul.dnsbl.sorbs.net,
                                reject_rbl_client sbl-xbl.spamhaus.org,
                                reject_rbl_client bl.spamcop.net,
                                reject_rbl_client cbl.abuseat.org,
                                reject_rbl_client ix.dnsbl.manitu.net,
                                reject_rbl_client combined.rbl.msrbl.net,
                                reject_rbl_client rabl.nuclearelephant.com,
                                permit
smtpd_data_restrictions = reject_unauth_pipelining
mime_header_checks = regexp:/etc/postfix/mime_header_checks.regexp

# SASL
smtp_sasl_auth_enable = yes
smtpd_sasl_auth_enable = yes
smtp_sasl_security_options =
# Pour corriger le bug de certains client (Outlook...)
broken_sasl_auth_clients = yes
smtpd_sasl_application_name= smtpd
smtpd_sasl_security_options = noanonymous
smtpd_sasl_local_domain =

# TLS/SSL
smtpd_tls_security_level = may
smtpd_tls_auth_only = no
smtp_use_tls = yes
smtpd_use_tls = yes
smtp_tls_loglevel = 1
smtpd_tls_received_header = yes
smtpd_tls_loglevel = 1
smtp_tls_note_starttls_offer = yes
smtpd_tls_key_file = /etc/postfix/ssl/smtpd.key
smtpd_tls_cert_file = /etc/postfix/ssl/smtpd.crt
smtpd_tls_CAfile = /etc/postfix/ssl/cacert.pem

# Divers
smtpd_tls_session_cache_timeout = 3600s
tls_random_source = dev:/dev/urandom

# Use Amavis
content_filter = amavis:[127.0.0.1]:10024
```

### saslauthd

Edit the file `/etc/defaut/saslauthd` and fill in these fields:

```bash {linenos=table,hl_lines=[1,3,4,7]}
START=yes
DESC="SASL Authentication Daemon"
NAME="saslauthd"
MECHANISMS="pam"
MECH_OPTIONS=""
THREADS=5
OPTIONS="-c -m /var/spool/postfix/var/run/saslauthd -r"
```

### smtpd.conf

Here, we'll define the authentication method. Let's create a smtpd.conf file in postfix:

```bash
mkdir /etc/postfix/sasl
```

Then create a file `/etc/postfix/sasl/smtpd.conf` and add this inside:

```bash
pwcheck_method: saslauthd
mech_list: plain login
```

### sasl_passwd

Now, we need to grant users access to SASL. Create the file `/etc/postfix/sasl_passwd` and fill it in like this:

```bash
mail   login:password
fqdn   login
```

Here's a simple example:

```bash
deimos@deimos.fr   deimos:xxxxx
```

Finally, we'll set rather restrictive rights as passwords are in clear text, then run postmap:

```bash
chmod 600 /etc/postfix/sasl_passwd
postmap /etc/postfix/sasl_passwd
```

## Validation

To validate everything, simply restart the services in question:

```bash
/etc/init.d/postfix restart
/etc/init.d/saslauthd restart
```

## Verification

To verify, a simple telnet will do the trick:

```bash {linenos=table,hl_lines=[1,5,11,12]}
> telnet deimos.fr 25
Connected to deimos.fr.
Escape character is '^]'.
220 *******************************************
ehlo test
250-fire.deimos.fr
250-PIPELINING
250-SIZE 10240000
250-ETRN
250-XXXXXXXA
250-AUTH LOGIN PLAIN
250-AUTH=LOGIN PLAIN
250-ENHANCEDSTATUSCODES
250-8BITMIME
250 DSN
```

Some might get this:

```bash
250-STARTTLS
```

instead of:

```bash
250-XXXXXXXA
```

## Resources
- [Postfix SMTP Authentication - On The Secure Port Only](/pdf/postfix_smtp_authentication_-_on_the_secure_port_only.pdf)
