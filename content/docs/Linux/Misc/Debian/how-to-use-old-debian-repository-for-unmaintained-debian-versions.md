---
weight: 999
url: "/How_to_use_old_Debian_repository_for_unmaintained_Debian_versions/"
title: "How to use old Debian repository for unmaintained Debian versions"
description: "A guide on how to configure Debian to use archive repositories for unmaintained Debian versions."
categories: ["Linux", "Debian", "System Administration"]
date: "2013-07-05T08:18:00+02:00"
lastmod: "2013-07-05T08:18:00+02:00"
tags: ["Debian", "Repository", "Archive", "Legacy"]
toc: true
---

![Debian](/images/debian.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Operating System** | Debian X |
| **Website** | [Debian Website](https://www.debian.org) |
| **Last Update** | 05/07/2013 |
{{< /table >}}

## Introduction

Installing old Debian to prepare a migration is quite good when repositories are still available. But for old unmaintained versions of Debian, we need to change some things.

This is what we're gonna look here.

## Usage

For this example, we're going to take lenny version which is totally deprecated. To make repository work, we need to change the repository url by 'archives':

```bash
# /etc/apt/sources.list
deb http://archive.debian.org/debian/ lenny main non-free contrib
deb-src http://archive.debian.org/debian/ lenny main non-free contrib
# Volatile:
deb http://archive.debian.org/debian-volatile lenny/volatile main contrib non-free
deb-src http://archive.debian.org/debian-volatile lenny/volatile main contrib non-free
# Backports:
deb http://archive.debian.org/debian-backports lenny-backports main contrib non-free
# Previously announced security updates:
deb http://archive.debian.org/debian-security lenny/updates main
```

That's all, you can update.

{{< alert context="info" text="It works with all Debian versions" />}}

## References

[https://superuser.com/questions/404806/did-debian-lenny-repositories-vanish](https://superuser.com/questions/404806/did-debian-lenny-repositories-vanish)
