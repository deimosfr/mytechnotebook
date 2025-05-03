---
weight: 999
url: "/Monitorer_les_accès_au_superuser/"
title: "Monitoring Superuser Access"
description: "How to monitor and track superuser access on Unix systems through logging mechanisms."
categories: ["Linux", "Security", "Solaris"]
date: "2009-02-04T18:34:00+02:00"
lastmod: "2009-02-04T18:34:00+02:00"
tags: ["su", "security", "admin", "monitoring", "Solaris", "logging"]
toc: true
---

## Introduction

When the operating system is installed, a superuser is created, with an UID of 0. The usage of the su command is recorded in `/var/adm/sulog`.

## Configuration

To record in the first place you need to do the following.
In the file `/etc/default/su`, uncomment the entry:

```bash
SULOG=/var/adm/sulog.
```

Save it.

The entries look like this (`/var/adm/sulog`):

```bash
MO 02/18 14:21 + pts/0 nrocha-root
TU 02/19 14:45 - pts/0 root-nrocha
WE 02/20 19:47 + pts/0 amaria-nrocha
```

* The first three columns show the time the event occurred.
* The fourth column shows a - for failed access and a + for successful access.
* The fifth column shows which port the access was made from.
* The last column shows the name of the user who tried to switch users and the switched user.

Note: This procedure was tested on the Solaris 10 OS.

## References

[https://wikis.sun.com/display/BigAdmin/Security+Administration+Tech+Tips](https://wikis.sun.com/display/BigAdmin/Security+Administration+Tech+Tips)
