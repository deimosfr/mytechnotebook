---
weight: 999
url: "/Changer_le_hostname_de_sa_solaris/"
title: "Changing Hostname on Solaris"
description: "How to change the hostname on Solaris systems both temporarily and permanently without requiring a reboot."
categories: ["Solaris"]
date: "2008-12-05T14:18:00+02:00"
lastmod: "2008-12-05T14:18:00+02:00"
tags: ["solaris", "hostname", "system administration"]
toc: true
---

## Introduction

While there are things about Solaris that can be frustrating, this is one aspect where it's impressive! You don't need to restart your machine to change the hostname - it updates live. I have to admit that's pretty cool.

## Temporary Method

Here's the temporary method (which will be lost after reboot):

```bash
hostname machine_name
```

And to verify:

```bash
hostname
```

## Manual Method

The method provided here doesn't use any utilities.

Edit the following files and modify the HOSTNAME value:

* `/etc/inet/hosts`
* `/etc/net/ticlts/hosts`
* `/etc/net/ticots/hosts`
* `/etc/net/ticotsord/hosts`
* `/etc/nodename`
* `/etc/hostname.[identifier of the interface that associates the IP address with the HOSTNAME]`

Then go to `/var/crash` and rename the directory:

```bash
cd /var/crash
mv oldname newname
```

All that's left is to restart.

In the case of a cluster, modify these files (**except if the cluster is not yet configured on this machine**):

* `/etc/cluster/nodeid`
