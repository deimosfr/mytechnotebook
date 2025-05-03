---
weight: 999
url: "/Debian_\\:_Erreur_GPG_lors_d'apt-get_update/"
title: "Debian: GPG Error During apt-get update"
description: "How to solve GPG key errors when running apt-get update on Debian systems"
categories: ["Security", "Linux", "Debian"]
date: "2006-08-05T13:05:00+02:00"
lastmod: "2006-08-05T13:05:00+02:00"
tags: ["Debian", "Security", "Linux", "APT", "GPG"]
toc: true
---

```bash
W: GPG error: http://security.debian.org testing/updates Release:
The following signatures couldn't be verified because the public key is not available: NO_PUBKEY key-number
```

Is this type of message appearing? It's quite annoying, especially since it happens every year. To solve this problem, follow these steps:

```bash
gpg --keyserver pgpkeys.mit.edu --recv-key key-number
gpg -a --export key-number | apt-key add -
```
