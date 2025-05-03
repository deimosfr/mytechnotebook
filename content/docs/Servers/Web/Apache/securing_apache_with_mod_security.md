---
weight: 999
url: "/Sécuriser_Apache_avec_mod_security/"
title: "Securing Apache with mod_security"
description: "Learn how to increase Apache security with mod_security module, a web application firewall to protect against SQL injection, XSS and other common attacks."
categories: ["Debian", "Linux", "Apache"]
date: "2009-01-18T03:41:00+02:00"
lastmod: "2009-01-18T03:41:00+02:00"
tags: ["Apache", "Security", "Web Server", "Firewall", "mod_security"]
toc: true
---

## Introduction

This is what I've been looking for quite some time! A module specifically designed for Apache security.

This module increases the security level of an Apache web server or other servers if used with Apache in proxy mode. Modsecurity acts as an application firewall embedded in Apache. It protects web applications against common attacks (SQL injection, Cross Site Scripting, etc.)

I found this nice documentation, but like most docs, it's missing some things. It's not much but I'm adding it anyway.

## Installation

If your Debian distribution doesn't have the packages, download them from the [Debian website](https://www.debian.org) then:

```bash
dpkg -i libapache2-mod-security* mod-security-common*
```

Then create a symbolic link to activate the module:

```bash
ln -s /etc/apache2/mods-available/mod-security.load /etc/apache2/mods-enabled/
```

Then restart to load everything:

```bash
/etc/init.d/apache2 restart
```

## Configuration

All you need to do is read [Mod security.pdf](/pdf/mod_security.pdf).

And finally, here's my configuration:

```bash
# Security discoverd with Nikto
TraceEnable "off"

# More Security
<IfModule mod_security.c>
   # Turn the filtering engine On or Off
   SecFilterEngine On

   # Server Signature
   SecServerSignature "Microsoft-IIS/5.0"

   # Make sure that URL encoding is valid
   SecFilterCheckURLEncoding On

   # Unicode encoding check
   SecFilterCheckUnicodeEncoding Off

   # Only allow bytes from this range
   SecFilterForceByteRange 0 255

   # Only log suspicious requests
   SecAuditEngine RelevantOnly

   # The name of the audit log file
   SecAuditLog /var/log/apache2/audit_log

   # Debug level set to a minimum
   SecFilterDebugLog /var/log/apache2/modsec_debug_log
   SecFilterDebugLevel 0

   # Should mod_security inspect POST payloads
   SecFilterScanPOST On

   # By default log and deny suspicious requests
   # with HTTP status 500
   SecFilterDefaultAction "deny,log,status:500"

   # Require HTTP_USER_AGENT and HTTP_HOST in all requests
   SecFilterSelective "HTTP_USER_AGENT|HTTP_HOST" "^$"

   # Weaker XSS protection but allows common HTML tags
   SecFilter "<[[:space:]]*script"

   # Prevent XSS atacks (HTML/Javascript injection)
   #SecFilter "<(.|n)+>"

   # Very crude filters to prevent SQL injection attacks
   SecFilter "delete[[:space:]]+from"
   SecFilter "insert[[:space:]]+into"
   # Replace "elect" with "select" in the line below
   SecFilter "elect.+from"
   SecFilter "drop[[:space:]]table"

   # Protecting from XSS attacks through the PHP session cookie
   SecFilterSelective ARG_PHPSESSID "!^[0-9a-z]*$"
   SecFilterSelective COOKIE_PHPSESSID "!^[0-9a-z]*$"

</IfModule>
```

## Resources
- [Mod Security Debian Etch Documentation](/pdf/modsecurity2_debian_etch.pdf)
- [Advanced Apache web server security: mod_security and mod_dosevasive](/pdf/mod_security_mod_dosevasive.pdf)
