---
weight: 999
url: "/Squirrelmail_\\:_Mise_en_place_d'un_webmail_simple_mais_évolué/"
title: "SquirrelMail: Setting up a Simple yet Advanced Webmail"
description: "Guide for installing and configuring SquirrelMail webmail solution"
categories: ["Linux"]
date: "2008-04-06T20:00:00+02:00"
lastmod: "2008-04-06T20:00:00+02:00"
tags: ["Servers", "Network", "Mail"]
toc: true
---

## Introduction

SquirrelMail is a webmail client that I used for years before switching to RoundCube. While it's not as visually attractive as RoundCube, it has the advantage of being both simple and feature-rich. The obvious prerequisites are a mail server like Postfix and a connection interface such as IMAP or POP.

## Installation

To install it, it's always very simple:

```bash
apt-get install squirrelmail
```

## Configuration

To configure SquirrelMail, simply use the tool provided:

```bash
squirrelmail-config
```

## Resources
- [Documentation on checking password strength for squirrelmail](/pdf/checking_password_strength_for_squirrelmail.pdf)
