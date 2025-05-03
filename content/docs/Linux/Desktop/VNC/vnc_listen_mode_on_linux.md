---
weight: 999
url: "/VNC_\\:_Mode_Listen_sous_Linux/"
title: "VNC: Listen Mode on Linux"
description: "How to use VNC in Listen mode on Linux systems"
categories: ["Linux", "Network"]
date: "2006-08-05T19:40:00+02:00"
lastmod: "2006-08-05T19:40:00+02:00"
tags: ["VNC", "Remote Access", "Network", "Linux"]
toc: true
---

## Introduction

Yes, in Windows, listen mode is easy! But on Linux, where is the "Run listening mode" icon? ;-)

It's actually quite simple:

```bash
xvnc4viewer -listen -Shared
```

The "Shared" option allows **multiple users to connect simultaneously**.

On the client side, you need to open incoming port **5500** and you're good to go.
