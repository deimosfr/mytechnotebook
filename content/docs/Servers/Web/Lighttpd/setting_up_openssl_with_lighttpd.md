---
weight: 999
url: "/Mise_en_place_d'OpenSSL_avec_Lighttpd/"
title: "Setting up OpenSSL with Lighttpd"
description: "This guide explains how to create and insert SSL certificates in Lighttpd for better website security."
categories: ["Linux", "Servers"]
date: "2009-04-15T07:09:00+02:00"
lastmod: "2009-04-15T07:09:00+02:00"
tags: ["SSL", "Lighttpd", "Security", "Web Server"]
toc: true
---

## Introduction

Adding security to your website is important. In this guide, we'll see how to create and insert SSL certificates in Lighttpd.

## Installation

We only need OpenSSL:

```bash
apt-get install openssl
```

## Configuration

### Generating SSL keys

Let's create an ssl directory in the Lighttpd configuration folder, then generate the certificates:

```bash
mkdir /etc/lighttpd/ssl
openssl req -new -x509 -keyout /etc/lighttpd/ssl/selfcert.pem -out /etc/lighttpd/ssl/selfcert.pem -days 3650 -nodes
```

* selfcert.pem: use the name that interests you (e.g., deimos.fr.pem)
* 3650: number of days the certificate is valid (10 years, we're safe for a good while)

### Lighttpd

Let's enable the SSL module for Lighttpd:

```bash
lighty-enable-mod ssl
```

Then let's modify the SSL configuration file so it takes our new certificate into account (`/etc/lighttpd/conf-available/10-ssl.conf`):

```bash
$SERVER["socket"] == "0.0.0.0:443" {
                 ssl.engine                  = "enable"
                 ssl.pemfile                 = "/etc/lighttpd/ssl/deimos.fr.pem"
}
```

And that's it! All you need to do now is restart Lighttpd, and port 443 will be open with your certificate activated :-)
