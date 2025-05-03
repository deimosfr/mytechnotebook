---
weight: 999
url: "/Mixing_Apache_Authentication/"
title: "Mixing Apache Authentication"
description: "How to mix different authentication methods in Apache including PAM, htaccess, IP restrictions, country-based access, and Radius authentication."
categories: ["Linux", "Ubuntu", "Apache"]
date: "2010-04-11T15:02:00+02:00"
lastmod: "2010-04-11T15:02:00+02:00"
tags: ["Apache", "Authentication", "Security", "PAM", "Radius"]
toc: true
---

## Mixing PAM

### Linux

How to mix PAM authentication (mod_auth_pam) and text file authentication (mod_auth) with Apache. First install this package:

```bash
apt-get install libapache2-mod-auth-pam
```

Then configure your htaccess:

```bash
AuthPAM_Enabled on
AuthPAM_FallThrough on
AuthAuthoritative Off
AuthUserFile /etc/apache2/htpassword
AuthType Basic
AuthName "Restricted Access"
Require valid-user
```

If mod_auth_pam doesn't find a valid user, it falls back to mod_auth authentication automatically.

Here is another example with webdav:

```apache {linenos=table}
    Alias /webdav /var/www/ngs
    <Location /webdav>
        DAV On
        AuthPAM_Enabled on
        AuthBasicAuthoritative Off 
        AuthPAM_FallThrough off 
        AuthUserFile /dev/null
        AuthType Basic
        AuthName "Webdav Authentication"
        Require group ngs 
    </Location>
```

### OpenBSD

On OpenBSD, I had to install [mod_auth_bsd](https://www.25thandclement.com/~william/projects/bsdauth.html):

```bash
pkg_add -iv mod_auth_bsd
```

Then, enable the module for Apache:

```bash
/usr/local/sbin/mod_auth_bsd-enable
```

Then restart Apache this way:

```bash
apachectl stop
apachectl start
```

Then in the Apache configuration `/var/www/conf/http.conf`, add this:

```bash
AuthBSDGroup auth

<Directory /var/www/htdocs/private>
   SSLRequireSSL
   AuthType Basic
   AuthName "ACME Login"
   AuthBSD On
   Require valid-user
</Directory>
```

## Restriction by IP address

Imagine using Jinzora. You don't want all your music to be accessible on the web. Simply add this to your VirtualHost configuration:

```bash
vi /etc/apache2/sites-enabled/000-default@
```

```bash
<Location /jinzora>
        Order deny,allow
        Deny from all
        Allow from 192.168.0.0/24
</Location>
```

This will allow all the 192.168.0.0 subnet to access your website. Then reload Apache:

```bash
/etc/init.d/apache2 reload
```

## Restriction by htaccess

This documentation is on how to protect a directory by htaccess (login + password).

Insert these lines and adapt to your configuration (`/etc/apache2/sites-enabled/000-default`):

```apache {linenos=table}
        <Directory /var/www/myhtaccess>
                AllowOverride AuthConfig
                Order allow,deny
                allow from all
        </Directory>
```

Then create a file `.htaccess` in **/var/www/myhtaccess** and put this:

```apache {linenos=table}
AuthType Basic
AuthName "Acces Prive"
AuthGroupFile /dev/null
AuthUserFile /etc/apache2/htaccesspassword

<Limit GET POST>
        Require valid-user
</Limit>

php_value magic_quotes_runtime 1
php_value magic_quotes_gpc 1
```

Then create your access file with the user (`/etc/apache2/htaccesspassword`):

```bash
htpasswd -c /etc/apache2/htaccesspassword username
```

For the next time, to add users, just remove "-c" like this:

```bash
htpasswd /etc/apache2/htaccesspassword username
```

Don't forget to restart Apache.

For a good documentation, follow this:
[Documentation on Htaccess](/pdf/htaccess.pdf)

## Authentication by Countries

[Deny Or Allow Countries With Apache htaccess](/pdf/deny_or_allow_countries_with_apache_htaccess.pdf)

## Authentication through Radius

Here is how to authenticate through a radius server:

[Radius Authentication](/pdf/apache_radius_authentication.pdf)  
[How To Configure Apache To Use Radius For Two-Factor Authentication On Ubuntu](/pdf/how_to_configure_apache_to_use_radius_for_two-factor_authentication_on_ubuntu.pdf)
