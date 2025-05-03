---
weight: 999
url: "/Gérer_des_certificats_SSL_signés_par_une_autorité_de_certification/"
title: "Managing SSL Certificates Signed by a Certificate Authority"
description: "Guide for generating and configuring SSL certificates signed by a certificate authority like StartSSL for Nginx and Lighttpd web servers."
categories: ["Linux", "Debian", "Nginx"]
date: "2015-02-17T03:58:00+02:00"
lastmod: "2015-02-17T03:58:00+02:00"
tags: ["SSL", "Security", "Lighttpd", "Nginx", "StartSSL", "Certificates", "HTTPS"]
toc: true
---

![StartSSL Logo](/images/startssl_logo.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Operating System** | Debian 7 |
| **Website** | [StartSSL Website](https://www.startssl.com) |
| **Last Update** | 17/02/2015 |
| **Others** | Nginx 1.2.1<br />Lighttpd 1.4.31 |
{{< /table >}}

## Introduction

You may need SSL certificates for your company or for personal needs on your website. The drawback of self-generated and self-signed certificates is that the first time you visit your site, you'll get a warning message.

To avoid this warning and to have a nice little padlock on your browser indicating that you're protected, you typically need to pay a certification authority a lot of money to get a valid SSL certificate.

However, there are kind companies that offer free or inexpensive certificates for your domain name that are properly signed :-). We'll look at how to do this with [StartCom](https://www.startcom.org/).

## Certificate Generation

Before you start, create your account and go to the private key step on [StartCom](https://www.startcom.org/).

Here, I'll present several servers:

1. Lighttpd: Generating a class 1 key[^1]
2. Nginx: Generating a wildcard key in class 2[^2]

Feel free to adapt to your configuration if you need to switch these.

### Lighttpd

We'll see here how to set up certificates with Lighttpd. First, let's create the essentials:

```bash
cd /etc/lighttpd
mkdir ssl
cd ssl
```

Next, we'll generate the RSA private key and secure it:

```bash
openssl genrsa -out server.key 4096
chmod 400 server.key
```

Copy the content of this key to the website so it can generate the rest. Then we'll create the CSR:

```bash
openssl req -new -nodes -key server.key -out server.csr
```

For the common name part, enter your default site (ex: www.deimos.fr).

Then download the StartCom certificates:

```bash
wget http://www.startssl.com/certs/ca.pem
wget http://www.startssl.com/certs/sub.class1.server.ca.crt
wget http://www.startssl.com/certs/sub.class1.server.ca.pem
```

### Nginx

We'll see here how to set up certificates with Nginx. First, let's create the essentials:

```bash
cd /etc/nginx
mkdir ssl
cd ssl
```

Next, we'll generate the RSA private key and secure it:

```bash
openssl genrsa -out server.key 4096
chmod 400 server.key
```

Copy the content of this key to the website so it can generate the rest. Then we'll create the CSR:

```bash
openssl req -new -nodes -key server.key -out server.csr
```

For the common name part, enter your default site (ex: www.deimos.fr).

Then download the StartCom certificates:

```bash
wget http://www.startssl.com/certs/ca.pem
wget http://www.startssl.com/certs/sub.class2.server.ca.crt
wget http://www.startssl.com/certs/sub.class2.server.ca.pem
```

### Certificate Signing

Now we'll generate a certificate on the StartSSL website. To begin, create your domain with the Validation Wizard:

![Startssl3](/images/startssl3.avif)

Choose Domain Name Validation:

![Startssl1](/images/startssl1.avif)

Then create the domain you want:

![Startssl2](/images/startssl2.avif)

Finish creating the domain and click on Certificates Wizard:

![Startssl3](/images/startssl3.avif)

Then select "Web Server SSL/TLS Certificate" as that's what we need:

![Startssl4](/images/startssl4.avif)

Skip this part since we've generated our own certificate:

![Startssl5](/images/startssl5.avif)

And paste the contents of the server.csr file into the text area:

![Startssl6](/images/startssl6.avif)

Complete the process, then create a **server.crt** file with the SSL certificate content that will be provided.

## Configuration

### Lighttpd

Next, we'll create a PEM certificate from those we've generated along with a CRT file:

```bash
cat server.key server.crt > server.pem
cat ca.pem sub.class1.server.ca.pem > ca-certs.crt
```

Then we'll configure our Lighttpd server to use our new keys (`/etc/lighttpd/conf-enabled/10-ssl.conf`):

```perl
## lighttpd support for SSLv2 and SSLv3
## 
## Documentation: /usr/share/doc/lighttpd-doc/ssl.txt
##  http://www.lighttpd.net/documentation/ssl.html 
 
#### SSL engine
$SERVER["socket"] == "0.0.0.0:443" {
                  ssl.engine                  = "enable"
                  ssl.pemfile                 = "/etc/lighttpd/ssl/server.pem"
                  ssl.ca-file                 = "/etc/lighttpd/ssl/ca-certs.crt"
}
```

Don't forget to restart your Lighttpd server for the parameters to take effect :-)

### Nginx

For Nginx, it's a bit different from Lighttpd. We'll create the unified certificate like this:

```bash
cat ssl.crt sub.class2.server.ca.pem ca.pem > /etc/nginx/ssl/server-unified.crt
```

Then configure Nginx (`/etc/nginx/sites-enabled/www.deimos.fr`):

```nginx
[...]
ssl on;
ssl_certificate /etc/nginx/ssl/server-unified.crt;
ssl_certificate_key /etc/nginx/ssl/server.key;
[...]
```

Then restart Nginx for the certificates to work.

## Resources

[^1]: http://forum.startcom.org/viewtopic.php?t=719
[^2]: http://www.startssl.com/?app=42
