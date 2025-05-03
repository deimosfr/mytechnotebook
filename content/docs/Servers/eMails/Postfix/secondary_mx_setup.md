---
weight: 999
url: "/MX_Secondaire\\:_mise_en_place/"
title: "Secondary MX: Setup"
description: "How to set up a secondary MX server with Postfix to handle emails during primary server downtime."
categories: ["Linux", "Servers", "Network"]
date: "2008-05-23T07:43:00+02:00"
lastmod: "2008-05-23T07:43:00+02:00"
tags: ["Postfix", "Email", "MX", "DNS", "Backup"]
toc: true
---

## Introduction

When there's only one mail server handling emails for a domain, if it becomes unavailable, emails sent by third parties to the primary server will be stored in the spool (outgoing queue) of the remote server for a few days in the best case, or immediately returned with an error message (depending on the problem or the configuration of the remote server).

Since problems and unavailability are inherent in modern computing, it's necessary to set up a system that can, at least in a degraded mode, transparently recover emails for the sender.

### What does the Secondary MX do?

Its job isn't very exciting (well, I think that's the case for many web services)... It spends its time waiting for emails that arrive when the sending server couldn't deliver to or contact the primary MX. When it receives them, it keeps them in its spool and tries at regular intervals to contact the primary server to transmit these emails.

**Don't forget to set up your mail server as an MX in DNS!**

## Configuration of Postfix

### On the secondary MX

#### Main.cf

Edit the `/etc/postfix/main.cf` file and adapt **(use only one of these)**:

```bash
mydestination = burnin.deimos.fr, deimos.fr, burnin, localhost
relay_domains = $mydestination, mavro.fr
```

The form enclosed in [] eliminates DNS MX lookups.

By default, the SMTP client performs DNS queries even if you specify a relay machine. If your machine doesn't have access to the DNS server, disable DNS lookups of the SMTP client as follows:

```bash
disable_dns_lookups = yes
```

For smtpd_recipient_restrictions, check that you have these two lines:

```bash
smtpd_recipient_restrictions = permit_mx_backup, permit_mynetworks, reject_unauth_destination
```

Add this line to indicate the primary server:

```bash
transport_maps = hash:/etc/postfix/transport
```

SMTP Banner: If you've set one with the machine name first, don't forget to change it to your secondary mx name:

```bash
smtpd_banner = burnin.deimos.fr - Microsoft Exchange (5.5)
```

Then insert this line to indicate the people to relay:

```bash
relay_recipient_maps = hash:/etc/postfix/relay_recipients
```

#### transport

Create a file `/etc/postfix/transport` and insert this:

```bash
domain_to_relay      smtp:[primary_server_FQDN]
```

Example:

```bash
deimos.fr     smtp:[fire.deimos.fr]
mavro.fr        smtp:[fire.deimos.fr]
```

#### relay_recipients

And now create the relay_recipients file:

```bash
xxx@mycompany.com     x
xxx@mycompany.com         x
```

This must contain the names of the people to relay.

#### mailname

Edit the `/etc/mailname` file and put your DNS:

```bash
deimos.fr
```

#### Validation

Let's validate everything now:

```bash
postmap /etc/postfix/transport
postmap /etc/postfix/relay_recipients
```

And we reload the Postfix configuration:

```bash
/etc/init.d/postfix reload
```

### On the primary MX

Nothing to do :-)

## Verifications

To see what we have in the spool:

```bash
mailq
```

If you're in a hurry to retrieve your emails after bringing your primary server back up:

```bash
mailq -q
```

## References

[Other documentation](/pdf/postfix_as_a_backup_mx.pdf)
