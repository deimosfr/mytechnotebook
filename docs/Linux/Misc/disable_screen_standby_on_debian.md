---
title: "Disable Screen Standby on Debian"
slug: disable-screen-standby-on-debian/
description: "How to disable screen standby and screen locking on Debian systems."
categories: ["Debian", "Linux"]
date: "2010-06-27T07:48:00+02:00"
lastmod: "2010-06-27T07:48:00+02:00"
tags: ["Servers", "Linux", "Debian", "Configuration"]
---

## Introduction

This is a simple but sometimes useful feature. By default, Debian turns the screen black after 30 minutes and locks the screen after 60 minutes.

It's possible to disable all of these features.

## Instructions

Edit the `/etc/console-tools/config` file and change these lines:

```bash
BLANK_TIME=0
POWERDOWN_TIME=0
```

Reboot and you're all set! :-)
