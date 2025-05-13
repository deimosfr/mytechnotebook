---
title: "Kill and lsof: Killing the process listening on a specific port"
slug: kill-and-lsof-killing-the-process-listening-on-a-specific-port/
description: "How to kill a process listening on a specific port using lsof and kill commands in Linux"
categories: ["Linux"]
date: "2007-03-08T08:58:00+02:00"
lastmod: "2007-03-08T08:58:00+02:00"
tags: ["Network", "Servers", "Linux", "Command line"]
---

You have certainly already looked for ways to kill a process that's listening on a specific port. Here's an example for port 1390:

```bash
kill $( lsof -i:1390 -t )
```
