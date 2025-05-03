---
weight: 999
url: "/Simulate_a_black_hole_for_a_domain_with_Postfix/"
title: "Simulate a black hole for a domain with Postfix"
description: "How to set up a black hole for emails in Postfix to test outgoing mail services."
categories: ["Debian", "Linux"]
date: "2013-05-07T05:54:00+02:00"
lastmod: "2013-05-07T05:54:00+02:00"
tags: ["Postfix", "Email", "SMTP", "Server Configuration"]
toc: true
---

![Postfix](/images/postfix_logo.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 2.10 |
| **Operating System** | Debian 7 |
| **Website** | [Postfix Website](https://www.postfix.org/) |
| **Last Update** | 07/05/2013 |
{{< /table >}}

## Introduction

When you manage outgoing emails through SMTP, you may sometimes need to test if a service is able to send correctly emails and itself check that there were no issue during sending. You can create a black hole for a specific domain and Postfix will answer from the same manner as if it is ok. It also permits to test Postfix sending capacities on a server.

## Installation

Of course you need Postfix:

```bash
aptitude install postfix
```

## Configuration

In the Postfix main configuration, add transport map line (`/etc/postfix/main.cf`):

```bash
[...]
transport_maps = dbm:/etc/postfix/blackhole_map
[...]
```

Now add the fake domain to the transport map file (`/etc/postfix/blackhole_map`):

```bash
blackhole.com     discard:silently
```

Let's generate the map now:

```bash
postmap -c /etc/postfix /etc/postfix/blackhole_map
```

Now restart postfix and you will see in your logs something like this:

```
postfix/discard[1435]: [ID 897546 mail.info] 4D847A5B: to=<john@blackhole.com>, relay=none, delay=13, delays=13/0/0/0, dsn=2.0.0, status=sent (silently)
```

## References

http://www.memoire-partagee.fr/2011/01/smtp-sortant-faire-un-trou-noir-avec-postfix/
