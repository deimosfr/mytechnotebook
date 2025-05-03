---
weight: 999
url: "/LemonLDAP\\:\\:NG_\\:_Plus_qu'un_simple_SSO/"
title: "LemonLDAP::NG: More than just SSO"
description: "A guide to LemonLDAP::NG, a comprehensive single sign-on solution with security features, authentication methods, and configuration options."
categories: ["Apache", "Nginx", "Debian"]
date: "2013-03-16T08:33:00+02:00"
lastmod: "2013-03-16T08:33:00+02:00"
tags: ["SSO", "SAML", "CAS", "Cookies", "HTTP", "Network", "Authentication"]
toc: true
---

![LemonDAP::NG](/images/lemonldap_logo.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 1.2.3 |
| **Operating System** | Debian 7 |
| **Website** | [LemonLDAP::NG Website](https://lemonldap-ng.org) |
| **Last Update** | 16/03/2013 |
| **Others** | Apache 2.2 |
{{< /table >}}

## Introduction

LemonLDAP::NG was created by Éric German for the French Ministry of Finance. Initially named LevonLDAP, in tribute to Novell, it was designed to be compatible with Novell's single sign-on (SSO) authentication system. It is based on the book "Writing Apache Modules with Perl and C The Apache API and mod_perl" by Doug MacEachern and Lincoln Stein (O'Reilly).

It was later renamed LemonLDAP::NG. From 2004, the project was gradually taken over by the French National Gendarmerie to become LemonLDAP::NG in 2005. Both projects coexisted for some time before LemonLDAP support was definitively abandoned.

It's an SSO system with a unique identifier/password pair. SSO does not handle access control.

## SSO Modes

- SSO by agent: installed on the client machine. No notion of security

![Simple Architecture](/images/archi_simple.avif)[^1]

- SSO by delegation: The user only needs their web browser. The server application points to the authentication portal

![Portal Architecture](/images/archi_portail.avif)[^2]

- Reverse proxy: Authentication is handled through a reverse proxy

![Reverse Proxy Architecture](/images/archi_reverse_proxy.avif)[^3]

## HTTP Requests

Before diving into LemonLDAP, we need to understand how the HTTP protocol works. Let's try to retrieve a website:

```bash {linenos=table,hl_lines=[1,5]}
> telnet www.deimos.fr 80
Trying 88.190.51.112...
Connected to shenzi.deimos.fr.
Escape character is '^]'.
GET /

HTTP/1.1 301 Moved Permanently
Server: nginx
Content-Type: text/html; charset=UTF-8
X-Pingback: http://blog.deimos.fr/xmlrpc.php
X-Powered-By: W3 Total Cache/0.9.2.8
Location: http://localhost/
Content-Length: 0
Accept-Ranges: bytes
Date: Fri, 22 Feb 2013 18:40:10 GMT
X-Varnish: 1562646426
Age: 0
Via: 1.1 varnish
Connection: close

Connection closed by foreign host.
```

```bash {linenos=table,hl_lines=[1,5,6]}
> telnet www.deimos.fr 80
Trying 88.190.51.112...
Connected to shenzi.deimos.fr.
Escape character is '^]'.
GET / HTTP/1.0
Host: www.deimos.fr 80

HTTP/1.1 200 OK
Server: nginx
Content-Type: text/html; charset=UTF-8
Vary: Accept-Encoding
X-Pingback: http://blog.deimos.fr/xmlrpc.php
X-Powered-By: W3 Total Cache/0.9.2.8
Link: <http://wp.me/2Q0VB>; rel=shortlink
Date: Fri, 22 Feb 2013 18:40:02 GMT
X-Varnish: 1562646424 1562646423
Age: 8
Via: 1.1 varnish
Connection: close

<!DOCTYPE html>
<html lang="fr-FR">
<head>
<meta charset="UTF-8" />
<meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=0">
<title>Deimosfr Blog | Parce que la mémoire humaine ne fait pas des Go</title>
<link rel="profile" href="http://gmpg.org/xfn/11" />
<link rel="stylesheet" type="text/css" media="screen" href="http://blog.deimos.fr/wp-content/themes/yoko/style.css" />
[...]
		</script>
</body>
</html>
<!-- Performance optimized by W3 Total Cache. Learn more: http://www.w3-edge.com/wordpress-plugins/

Page Caching using apc
Object Caching 8437/8588 objects using apc

 Served from: blog.deimos.fr @ 2013-02-22 19:39:54 by W3 Total Cache -->Connection closed by foreign host.
```

The return code is 200, everything is fine. You can check the list of HTTP return codes here: http://en.wikipedia.org/wiki/List_of_HTTP_status_codes

## Cookies

The HTTP protocol is stateless. To maintain a persistent connection, we need to use HTTP 1.1 protocol. We also use cookies to store sessions. On the server side, the user session ID is stored. A cookie can be up to 4096 bytes maximum.

There are several types of cookies:

- Session cookies: The cookie will remain active for a defined time or until the browser is closed
- Persistent cookies: They remain active permanently

A cookie exchange works like this:

1. A simple request
2. Server response with a "Set-cookie" header field
3. Client request with a "Cookie" header field

As a reminder, cookies are only valid on a single DNS domain.

## LemonLDAP

### The Components

LemonLDAP::NG uses 3 components:

- The portal: authentication interface, application menu, password change
- The Handler: Apache module that controls access to web applications
- The Manager: the graphical part for configuring LemonLDAP

### Communication

The advantage compared to CAS is that the client does not need to go back through LemonLDAP::NG with each web service change.

![LemonLDAP::NG Architecture](/images/lemonldap-ng-architecture.avif)

## Authentication Phases

![LemonLDAP::NG SSO](/images/lemonldapng-sso.avif)[^4]

1. When a user tries to access a protected application, their request is intercepted by the Handler
2. If the SSO cookies are not detected, the Handler redirects the user to the portal
3. The user authenticates on the portal
4. The portal verifies their authentication
5. If validated, the portal retrieves the user information
6. The portal creates a session where it stores the user information
7. The portal retrieves a session key
8. The portal creates SSO cookies with session key/value
9. The user is redirected to a protected application with their new cookie
10. The Handler retrieves the cookie and session
11. The Handler records user data in its cache
12. The Handler checks access rights and sends headers to protected applications
13. Protected applications send a response to the Handler
14. The Handler sends a response to the user

### The Different Databases

Several databases are used:

- Authentication: how to validate authentication data
- Users: user data
- Password: where to change the user password
- Provider: how to provide identity to an external service

For example, you can use Kerberos with LDAP.

Internal databases:

- Sessions: server-side session storage
- Configuration: configuration storage (versioned)
- Notifications

### Authentication Methods

- LDAP
- Database
- SSL X509
- Apache modules (Kerberos, OTP...)
- SAML 2.0
- OpenID
- Twitter
- CAS
- Yubikey
- Radius

## Session Storage

LemonLDAP::NG uses 3 levels of cache:

- Apache::Session::\*: final storage of sessions
- Cache:Cache\*: allows the Handler to share data between Apache threads and processes
- Internal variables to LemonLDAP::NG::Handler: if the same user uses the same thread or process again.

## Installation

Now for the practical part! You can either use the version available through the official repositories, or add the LemonLDAP::NG repository, which is what we'll do here:

```bash
# LemonLDAP::NG repository
deb     http://lemonldap-ng.org/deb squeeze main
deb-src http://lemonldap-ng.org/deb squeeze main
```

Then update:

```bash
aptitude install lemonldap-ng
```

Let's change the default configuration which is example.com:

```bash
sed -i 's/example\.com/deimos.fr/g' /etc/lemonldap-ng/* /var/lib/lemonldap-ng/conf/lmConf-1 /var/lib/lemonldap-ng/test/index.pl
```

After this, you should not modify these configuration files manually!

Let's enable the sites and the Perl module for Apache:

```bash
a2ensite handler-apache2.conf
a2ensite portal-apache2.conf
a2ensite manager-apache2.conf
a2ensite test-apache2.conf
a2enmod perl
```

Then restart Apache:

```bash
apache2ctl restart
```

## Configuration

### Manager

For the manager to work:

```bash
echo "127.0.0.1 reload.example.com" >> /etc/hosts
```

### DNS Configuration

For DNS resolution to work correctly:

```bash
cat /etc/lemonldap-ng/for_etc_hosts >> /etc/hosts
```

### Test Application

There are test1.example.com and test2.example.com to test your sites with username and password 'dwho': http://test1.example.com

Other admin accounts that exist are:

- rtyler/rtyler
- msmith/msmith
- dwho/dwho

This page will allow you to see important information that will be exchanged with LemonLDAP::NG. It's useful for debugging in addition to Apache logs.

### Protection

In the file `/etc/lemonldap-ng/lemonldap-ng.ini` we can configure who has access to the manager!

```ini {linenos=table,hl_lines=[12]}
[...]
# Manager protection: by default, the manager is protected by a demo account.
# You can protect it :
# * by Apache itself,
# * by the parameter 'protection' which can take one of the following
# values :
#   * authenticate : all authenticated users can access
#   * manager      : manager is protected like other virtual hosts: you
#                    have to set rules in the corresponding virtual host
#   * rule: <rule> : you can set here directly the rule to apply
#   * none         : no protection
protection = manager
[...]
```

By default, the protection is set to Manager, which is good. Only authorized people can connect to it (VirtualHost). It's possible to use rules to specify something specific (a uid for example).  
Authenticate allows any user who connects to have access to the manager. And the last option, none, is strongly discouraged as it allows anyone to access it.

### Macros

Macros allow you to create variables in LemonLDAP. For example:

```
$fullname => $givenname . '' . $surname
```

You can then reuse $fullname.

### Sessions

By default, every 10 minutes (via cron) there is a check for sessions that need to be purged: `/etc/cron.d/liblemonldap-ng-portal-perl`.

### HTTP Headers

### Script Parameters

You can modify GET requests to POST, for example.

### Notifications

They allow information to be validated by users and are stored in persistent sessions.

For example, you can create a disclaimer.

## CAS

LemonLDAP::NG can handle CAS as either a client or server. CAS only performs URL redirections, with the possibility of proxy tickets.

### Server

When acting as a CAS server, you need to enable rewrite rules:

```bash
a2enmod rewrite
/etc/init.d/apache2 restart
```

### Client

```bash
aptitude install libauthcas-perl
```

## SAML

Here are some terms to understand:

- IDP: Identity Provider
- CoT: Circle of Trust
- InterCoT: Circle of Trust between IDPs
- AA: Attribute Authority
- Proxy IDP: proxy for IDP to transfer identity requests
- SP: Service Provider

To install LASSO (SAML), you need to install this:

```bash
aptitude install liblasso-perl
```

## References

[^1]: http://www-igm.univ-mlv.fr/~dr/XPOSE2006/CLERET/techniques.html
[^2]: http://www-igm.univ-mlv.fr/~dr/XPOSE2006/CLERET/techniques.html
[^3]: http://www-igm.univ-mlv.fr/~dr/XPOSE2006/CLERET/techniques.html
[^4]: http://lemonldap-ng.org/documentation/presentation
