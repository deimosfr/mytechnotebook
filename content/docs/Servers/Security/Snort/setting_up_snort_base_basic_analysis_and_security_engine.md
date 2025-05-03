---
weight: 999
url: "/Mise_en_place_de_Snort_\\&_BASE_(Basic_Analysis_and_Security_Engine)/"
title: "Setting up Snort & BASE (Basic Analysis and Security Engine)"
description: "How to install and configure Snort IDS with BASE web interface for network intrusion detection on Debian systems"
categories: ["Debian", "Security", "Database"]
date: "2008-02-10T10:43:00+02:00"
lastmod: "2008-02-10T10:43:00+02:00"
tags: ["IDS", "Snort", "BASE", "Network Security", "MySQL"]
toc: true
---

## Introduction

Snort is what we call an [IDS (Intrusion Detection System)](https://fr.wikipedia.org/wiki/Syst%C3%A8me_de_d%C3%A9tection_d%27intrusion) and more specifically a passive NIDS (Network Intrusion Detection System). It can therefore detect who is trying to compromise your system.

Currently it's not perfect, but it's still less expensive than some [IPS (Intrusion Prevention System)](https://fr.wikipedia.org/wiki/Syst%C3%A8me_de_pr%C3%A9vention_d%27intrusion). Snort coupled with BASE provides real convenience for intrusion detection.

## Installation and configuration

[Documentation on installing and configuring Base and Snort](/pdf/ids_snort_base.pdf)

### Using Debian packages

For my part, I only followed a small portion since I used Debian packages directly (advantage of automatic updates). For those who want to take the same route:

```bash
apt-get install snort-mysql php5-gd libpcre3 acidbase python-adodb
```

Explanations for beginners:

- Snort is the tool that will listen in promiscuous mode on one or more of your network cards and thus detect potential intrusion attempts
- BASE/AcidBase is the one that will read the results of Snort recorded in the SQL database (or other)

It's mentioned in the documentation, but for people like me who prefer to read between the lines, here's the command to test your snort configuration:

```bash
snort -c /etc/snort/snort.conf
```

## FAQ

### BASE: Database ERROR: Table 'snort.iphdr' doesn't exist

If you encounter this problem after an update or reinstallation, you just need to reimport an SQL file. Better practice than long explanations:

```bash
cd /usr/share/doc/snort-mysql/
gzip -d create_mysql.gz
mysql -uroot -pPASSWORD -D snort < create_mysql
```

That's it! So it wasn't really a big deal, and BASE is running again :-)

## Resources

- [Another Documentation on Base and Snort](/pdf/snort_base.pdf)
- [Setting up Snort (IDS), OSSEC (HbIDS) and Prelude (HIDS)](/pdf/av04mihf.pdf)
- [https://oinkmaster.sourceforge.net/](https://oinkmaster.sourceforge.net/)
