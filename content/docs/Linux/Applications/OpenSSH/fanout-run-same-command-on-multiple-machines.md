---
weight: 999
url: "/fanout-run-same-command-on-multiple-machines/"
title: "Fanout: Run the Same Command on Multiple Machines Simultaneously"
description: "How to use Fanout to execute commands on multiple servers at the same time"
categories: ["Linux", "System Administration", "Tools"]
date: "2006-10-03T15:52:00+02:00"
lastmod: "2006-10-03T15:52:00+02:00"
tags: ["Fanout", "SSH", "Automation", "Administration"]
toc: true
---

These two tools by William Stearns can quickly become essential when managing multiple machines...

As indicated in the title, these tools allow you to run commands simultaneously on multiple machines, and even interactively (with fanterm).

I didn't find anything in the repositories or as a .deb package. So I "alienated" the rpm:

```bash
wget http://www.stearns.org/fanout/fanout-0.6.1-0.noarch.rpm
sudo alien fanout-0.6.1-0.noarch.rpm
sudo dpkg -i fanout_0.6.1-1_all.deb
```

For SSH on non-standard ports, I used the config file (see man ssh_config)

```bash
cat .ssh/config
```

```
Host bipbip
Hostname 172.16.2.200
User yannick
Port 2222
```

Authentication is done by key (see man ssh)

I replaced xterm with Eterm in `/usr/bin/fanterm`
