---
weight: 999
url: "/Crypter_et_décrypter_un_fichier_avec_OpenSSL/"
title: "Encrypt and Decrypt a File with OpenSSL"
description: "Learn how to encrypt and decrypt files using OpenSSL with simple commands and password protection."
categories: ["Security", "Linux"]
date: "2009-12-11T20:49:00+02:00"
lastmod: "2009-12-11T20:49:00+02:00"
tags: ["OpenSSL", "Encryption", "Security", "Cryptography", "Linux"]
toc: true
---

## Introduction

The OpenSSL Project is a collaborative effort to develop a robust, commercial-grade, full-featured, and Open Source toolkit implementing the Secure Sockets Layer (SSL v2/v3) and Transport Layer Security (TLS v1) protocols as well as a full-strength general purpose cryptography library. The project is managed by a worldwide community of volunteers that use the Internet to communicate, plan, and develop the OpenSSL toolkit and its related documentation.

OpenSSL is based on the excellent SSLeay library developed by Eric A. Young and Tim J. Hudson. The OpenSSL toolkit is licensed under an Apache-style licence, which basically means that you are free to get and use it for commercial and non-commercial purposes subject to some simple license conditions.

## Encrypt

```bash
$ openssl des3 -salt -in file.log -out file.des3
enter des-ede3-cbc encryption password:
Verifying - enter des-ede3-cbc encryption password:
```

The above will prompt for a password, or you can put it in with a -k option, assuming you're on a trusted server.

## Decrypt

```bash
openssl des3 -d -salt -in file.des3 -out file.txt -k mypassword
```
