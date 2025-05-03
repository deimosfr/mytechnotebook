---
weight: 999
url: "/Swaks_-_Swiss_Army_Knife_SMTP/"
title: "Swaks - Swiss Army Knife SMTP"
description: "Learn about Swaks (Swiss Army Knife SMTP), a versatile tool for SMTP testing and management with examples of use cases."
categories:
  - "Linux"
  - "Network"
date: "2006-10-03T13:40:00+02:00"
lastmod: "2006-10-03T13:40:00+02:00"
tags:
  - "smtp"
  - "email"
  - "networking"
  - "tools"
toc: true
---

You may know Netcat, the "TCP/IP Swiss Army Knife", now here's [Swaks](https://jetmore.org/john/code/#swaks), the "SMTP Swiss Army Knife" :)

```bash
apt-get install swaks
```

The description is accurate, I can't say it better... you can do almost anything with it.
Just read the man page to get ideas for how to use it.

An example of usage that allows you to reinject emails without having to go through a local server:

```bash
for i in *.eml; do cat "$i" | swaks -g -n -t recipient@yop.tld -f sender@yop.tld -s 1.2.3.4 ; done
```

If you don't specify a server, it even resolves the MX record...  
In short, it's essential for anyone who works with SMTP.
