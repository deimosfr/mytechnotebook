---
weight: 999
url: "/Crypter_un_mot_de_passe_en_MD5/"
title: "Encrypting a Password with MD5"
description: "Learn how to encrypt passwords using MD5 hash and check file integrity on Linux and Solaris systems."
categories: ["Linux", "Solaris"]
date: "2008-12-17T16:51:00+02:00"
lastmod: "2008-12-17T16:51:00+02:00"
tags: ["Security", "Encryption", "Linux", "Solaris", "MD5"]
toc: true
---

## Introduction

MD5 hashing allows you to verify the integrity of a file. In short, it creates a fingerprint that helps you verify (after a transfer, for example) whether the file arrived correctly and is not corrupted.

Today, MD5 encryption has been broken, and it's possible to simulate a fake fingerprint. Therefore, you should be careful about how you use it. Personally, I use it to verify the integrity of downloaded files (as it's a very common method). For example, when I download Sun Cluster, I verify the MD5 hash to ensure the installation is complete and not corrupted.

## Obtaining an MD5 Hash

### A Word

To encrypt a password in MD5, it's not very complicated, but when you don't do it every day, you might not remember how. So here it is:

```bash
echo -n password | md5sum
```

### A File on Linux

On Linux, the md5sum command is used:

```bash
md5sum /bin/ls
c854f8350b2a2873d4c2635813a797cc  /bin/ls
```

### A File on Solaris

On Solaris, we'll use the digest command to calculate the MD5 hash:

```bash
$ digest -a md5 -v /bin/ls
md5 (/bin/ls) = b57e173220af4b919f1d4bef9db11482
```
