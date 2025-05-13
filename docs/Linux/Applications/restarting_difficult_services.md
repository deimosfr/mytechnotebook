---
title: "Restarting difficult services"
slug: restarting-difficult-services/
description: "How to restart difficult services like SSH when there's no simple solution to stop them."
categories: ["Linux"]
date: "2007-11-14T09:49:00+02:00"
lastmod: "2007-11-14T09:49:00+02:00"
tags: ["Servers", "Services", "SSH"]
---

As you've probably noticed, there are some services like SSH that don't have a simple solution to stop the service.

That's why I'm giving you this solution. For example with SSH, if I want to restart it:

```bash
kill -HUP `cat /var/run/sshd.pid`
```

It's as simple as that. Now you can manually restart it:

```bash
/usr/sbin/sshd
```
