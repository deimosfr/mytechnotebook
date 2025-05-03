---
weight: 999
url: "/Encrypter_sa_swap/"
title: "Encrypting Swap Partition"
description: "A guide on how to encrypt your swap partition in OpenBSD systems both with and without rebooting."
categories: ["Linux", "BSD"]
date: "2007-08-01T22:54:00+02:00"
lastmod: "2007-08-01T22:54:00+02:00"
tags: ["OpenBSD", "Security", "Encryption", "Swap", "System Administration"]
toc: true
---

Since OpenBSD 3.7, swap is automatically encrypted. If you're using an earlier version and wish to enable encryption, there are two solutions:

## Without Rebooting

To enable swap encryption without rebooting, use the following command:

```bash
sysctl -w vm.swapencrypt.enable=1
```

## With Rebooting

To enable swap encryption permanently (requires a reboot), edit the `/etc/sysctl.conf` file and uncomment this line:

```bash
vm.swapencrypt.enable=1
```
