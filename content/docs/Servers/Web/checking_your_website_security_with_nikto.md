---
weight: 999
url: "/Vérifier_la_sécurité_de_son_site_web_avec_Nikto/"
title: "Checking Your Website Security with Nikto"
description: "Guide on how to use Nikto to test web server security and detect potential vulnerabilities."
categories: ["Security", "Linux", "Apache"]
date: "2011-04-05T21:05:00+02:00"
lastmod: "2011-04-05T21:05:00+02:00"
tags: ["Security", "Web Security", "Pentesting", "Nikto", "Audit", "Wikto"]
toc: true
---

## Introduction

To verify your configuration file and test potential security vulnerabilities, here's a practical Perl script called Nikto. These are the most well-known security audit applications similar to Nikto:

- [WebScarab/ProxMon](https://www.secuobs.com/news/07042007-proxmon.shtml)
- [WebInspect](https://www.secuobs.com/news/12022007-webinspect.shtml) by SPI Dynamics
- [Wikto](https://www.secuobs.com/news/20042005-google-hacking.shtml)
- [Exploit-me](https://www.secuobs.com/news/04122007-firecat.shtml)
- [Paros Proxy](https://www.parosproxy.org/index.shtml)

## Installation and Configuration

Here's the documentation I found:

[Apache Security Testing](/pdf/security_testing_your_apache_configuration_with_nikto.pdf)

For those who don't want to recompile packages:

```bash
aptitude install nikto
```

Then it's simple, as described in the documentation:

```bash
nikto -h localhost
```

## Resources
- [ProxyStrike a new transparent proxy for web application auditing](https://www.secuobs.com/news/15042008-proxystrike.shtml)
