---
weight: 999
url: "/Installation_et_configuration_d\\'Apache_2/"
title: "Apache 2 Installation and Configuration"
description: "This guide covers the installation and configuration of Apache 2 web server, including authentication, URL rewrites, virtual hosts, and performance optimizations."
categories: ["Monitoring", "Security", "Ubuntu"]
date: "2011-04-06T08:57:00+02:00"
lastmod: "2011-04-06T08:57:00+02:00"
tags:
  [
    "Apache",
    "Authentication",
    "LDAP",
    "VirtualHost",
    "URL rewrite",
    "Web Server",
    "Security",
  ]
toc: true
---

## Introduction

Apache is one of the most widely used web servers in the world, if not THE most used web server in the world.

## Installation

To install it:

```bash
apt-get install apache2
```

## Configuration

### Choose a Default Charset

In the configuration file `/etc/apache2/apache2.conf` or `/etc/apache2/conf.d/charset`, insert this:

```bash
AddDefaultCharset          .latin9
```

Then reload Apache:

```bash
/etc/init.d/apache2 reload
```

### Authentication

#### LDAP

For LDAP authentication, you need to install this:

```bash
apt-get install libapache-authznetldap-perl
```

Then enable the module and restart the server:

```bash
a2enmod authnz_ldap
/etc/init.d/apache2 restart
```

Now, for the configuration part, I'll take the example of Nagios3 where we need to modify the "DirectoryMatch" section as follows:

```apache
<DirectoryMatch (/usr/share/nagios3/htdocs|/usr/lib/cgi-bin/nagios3)>
    Options FollowSymLinks

    DirectoryIndex index.html

    AllowOverride AuthConfig
    Order Allow,Deny
    Allow From All

    AuthName "Nagios Access"
    AuthType Basic
    #AuthUserFile /etc/nagios3/htpasswd.users
    # nagios 1.x:
    #AuthUserFile /etc/nagios/htpasswd.users
    #require valid-user

    # auth from ldap
    AuthzLDAPAuthoritative on
    AuthBasicProvider ldap
    AuthLDAPURL ldap://ldap/dc=openldap,dc=mycompany,dc=lan?uid?sub?(objectClass=posixAccount)
    AuthLDAPRemoteUserIsDN off
    AuthLDAPGroupAttribute memberUid
    AuthLDAPGroupAttributeIsDN off
    Require ldap-group cn=prod,ou=Groupes,dc=openldap,dc=mycompany,dc=local
    Require ldap-group cn=sysnet,ou=Groupes,dc=openldap,dc=mycompany,dc=local
    Require ldap-user nagiosadmin
</DirectoryMatch>
```

Here, I have 2 groups (sysnet and prod) that are authorized to connect.

### Skip Authentication for Specific IP Addresses

I need monitoring screens to access Nagios without authentication while keeping LDAP authentication for other users. Building on the example above, here are the lines to modify:

```apache
    AllowOverride AuthConfig
    Require valid-user
    Order Deny,Allow
    Allow From 10.100.10.0/24
    Satisfy Any
```

This way, IPs from the 10.100.10.0/24 subnet don't need to authenticate while others do. **To decide whether to validate one solution or the other, I use the Satisfy Any directive**. We can put 'Satisfy All' if we want all conditions to be validated.

### Creating Redirects

If you want to protect a specific folder, you have 2 methods:

- Complete prohibition
- Redirection

PS: I won't discuss special cases such as htaccess here, [see this documentation]({{< ref "docs/Servers/Web/Apache/mixing_apache_authentication.md">}}).

To prohibit access:

```apache
<Directory /my/folder/to/block>
        Order allow,deny
        Deny from all
</Directory>
```

If you want to create a redirection, insert this in your "Directory" section:

```apache
RedirectMatch ^/$ http://mysecureshell.sourceforge.net
```

This will redirect to the MySecureShell site :-)
Or if you want to redirect to a local folder:

```apache
RedirectMatch ^/$ /my_folder/
```

There's also the ultimate solution:

```apache
Redirect /myfolder http://mysecureshell.sourceforge.net
```

or even:

```apache
RedirectMatch ^(.*)$ https://www.deimos.fr$1
```

Which gives me:

```apache
<VirtualHost *:80>
        ServerName www.deimos.fr
        ServerAlias deimos.fr www.deimos.fr
        RedirectMatch ^(.*)$ https://www.deimos.fr$1
</VirtualHost>
```

#### HTML Redirector

Here's a very simple solution for creating a redirector. Just place an index.html file in the desired folder with this content:

```html
<?xml version="1.0" encoding="ISO-8859-1"?>
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
  <head>
    <meta http-equiv="content-type" content="text/html; charset=ISO-8859-1" />
    <meta http-equiv="refresh" content="0; url=http://www.google.com" />
    <title>www.deimos.fr</title>
    <meta name="robots" content="noindex,follow" />
  </head>

  <body>
    <p><a href="http://www.google.com">Please wait while redirecting...</a></p>
  </body>
</html>
```

#### VirtualHost

When we have an Apache server at the front end and want to redirect traffic to other Apache servers at the back end, we need to activate mod_proxy. Here's an example:

```bash
ProxyRequests Off
NameVirtualHost 1.2.3.4 # IP of your box

<VirtualHost 1.2.3.4> # Website managed by Apache
   ServerName www.domain.tld
   DocumentRoot /var/www/htdocs/ # etc...
</VirtualHost>

<VirtualHost 1.2.3.4>
   ServerName www.domain2.tld
   ErrorLog blabla
   CustomLog blabla
   ProxyPassReverse /
   http://127.0.0.1:8002/
</VirtualHost>

<VirtualHost 1.2.3.4>
   ServerName www.domain3.tld
   ErrorLog blabla
   CustomLog blabla
   ProxyPassReverse /
   http://127.0.0.1:8003/
</VirtualHost>
```

Here, depending on the URL the client entered, there will be automatic redirects to other servers.

### URL Rewriting Redirects

Here's an example of URL rewriting. This allows redirecting cvsweb.mydomain.com automatically to the correct URL and cleaning up the URL as well. I changed from:

- http://machine.mydomain.com/cgi-bin/cvsweb

to

- http://cvsweb.mydomain.com

Here's the solution. First, let's load the module:

```bash
a2enmod rewrite
```

Then we'll write this in our configuration file (`/etc/apache2/sites-enabled/000-default`):

```bash
<VirtualHost cvsweb.mydomain.com:80>
        ServerName http://cvsweb.mydomain.com
        ServerAlias cvsweb
        ServerAdmin it-system@mydomain.com
        DocumentRoot /var/www/
        ScriptAlias /cgi-bin/ /usr/lib/cgi-bin/
        LogLevel warn
        ServerSignature On
        RedirectMatch ^/$ /cgi-bin/cvsweb/

        ### Rewrite http://cvsweb as http://cvsweb.mydomain.com
        RewriteEngine On
        RewriteCond %{HTTP_HOST} ^cvsweb$
        RewriteRule ^(.*)$ http://cvsweb.mydomain.com/$1 [R=301,L]
</VirtualHost>
```

Then force Apache to reload everything:

```bash
/etc/init.d/apache2 force-reload
```

#### Block Internet Explorer Access to Your Site

It's possible to block access to all kinds of browsers. If like me, you're not friends with IE which breaks your PNGs in version 6, doesn't respect standards, breaks CSS, etc., it might be convenient to block it and politely direct the user to download Firefox as soon as possible.

For this, we'll use the rewrite mode. It must be enabled as described above. Then add these lines in the desired folder (Directory for the entire site for example) in sites-enabled/000-default:

```apache
<Directory />
   ...
   AllowOverride FileInfo
   <IfModule mod_rewrite.c>
      RewriteEngine on
      RewriteCond %{HTTP_USER_AGENT} .*MSIE.*
      # opera sometimes pretends to be IE
      RewriteCond %{HTTP_USER_AGENT} !.*Opera.*
      # avoid infinite loop in conditions
      RewriteCond %{REQUEST_FILENAME} !.*ie.html
      # redirect to a page explaining the reasons for rejection
      RewriteRule .* /ie.html  [L]
   </IfModule>
</Directory>
```

All that's left is to create the ie.html file and put your nice text in it (you can also make a simple text file). Here's what I use (`/var/www/ie.html`):

```html
<!DOCTYPE html PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN" "http://www.w3.org/TR/html4/loose.dtd">
<html>
  <head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
    <title>Access Forbidden with Internet Explorer</title>
  </head>
  <body>
    Dear Internet User,<br />
    <br />
    This site cannot be accessed using Internet Explorer.<br />
    You should now understand that times are changing.<br />
    <br />
    You are currently using a browser (Internet Explorer) that doesn't<br />
    respect <a href="http://www.w3.org">standards</a> and holds a monopoly due
    to its mandatory omnipresence in<br />
    your dear OS (Windows). That said, perhaps you're at work and don't have a
    choice of OS.<br />
    <br />
    However, Internet Explorer should no longer be used when there are many<br />
    other free, open-source browsers that respect standards!<br />
    But since you don't seem to be aware of this, it's okay, let me help you.<br />
    <br />
    To begin with, you can download a clean browser like
    <a href="http://www.mozilla.org">Firefox</a>.<br />
    This would already help you get on the right track and will allow you to
    access<br />
    my site.<br />
    <br />
    Still on your new quest to the light side of the force, you should switch
    to<br />
    a free, open-source OS (<a href="http://www.ubuntu.com">Ubuntu</a> for
    example) that will surely make you happy.<br />
    <br />
    I encourage you to take control as soon as possible.<br />
    <br />
    Regards,<br />
    Pierre (aka Deimos)
  </body>
</html>
```

### Public Folders

Public folders are used to have multiple clients on a server where each has their own personal space. The practice is quite simple: we have, for example, the user toto who has a "public_html" folder in their home directory, and their web server is accessible via "http://server/~toto". I did this on OpenBSD with Apache 1.3; normally for version 2, the syntax is the same. So here's the configuration to add:

```apache
UserDir public_html

<Directory /home/clients/*/public_html>
   AllowOverride FileInfo AuthConfig Limit
   Options MultiViews Indexes SymLinksIfOwnerMatch IncludesNoExec

   # Look
   HeaderName /header.htm

   <Limit GET POST OPTIONS PROPFIND>
       Order allow,deny
       Allow from all
   </Limit>
   <Limit PUT DELETE PATCH PROPPATCH MKCOL COPY MOVE LOCK UNLOCK>
       Order deny,allow
       Deny from all
   </Limit>
</Directory>
```

You can also see that I changed the header of my main pages with the "HeaderName" option. This header.htm file must be located in the "DocumentRoot" folder when called by "/".

Here's an example with a mix of BSD authentication + IP restriction:

```bash
UserDir download

AuthBSDGroup auth

<Directory "/home/clients/*/download">
    AllowOverride FileInfo AuthConfig Limit
    Options MultiViews Indexes SymLinksIfOwnerMatch IncludesNoExec

    ### IP security for client connections ###
    Order deny,allow
    Deny from all
    Allow from 192.168.0.0/24
    #################

    # Look
    HeaderName /header.htm

    # Mod Auth BSD
    SSLRequireSSL
    AuthType Basic
    AuthName "Client Authentication"
    AuthBSD On
    AuthBSDKeepPass Off
    AuthBSDStrictRequire On
    <Limit GET POST OPTIONS PROPFIND PUT DELETE PATCH PROPPATCH MKCOL COPY MOVE LOCK UNLOCK>
        Require valid-user
    </Limit>
</Directory>
```

### Modify Header and Footer

#### Standard

If you want to change the top and bottom of your navigation pages (those that allow browsing files and folders), you can use:

```apache
   HeaderName /.header.htm
   ReadmeName /.footer.htm
```

Just create the .header.htm and .footer.htm files and put whatever you want in them :-)

#### Advanced

Sometimes you might want to make things a little more interactive than just simple HTML. But you'll run into a significant problem since it simply won't be able to interpret your code. In my case, I wanted to do it in PHP, so here's the solution. In your _Directory_ section, where you already have your lines containing HeaderName and ReadmeName, you should insert these lines:

```apache
   # In order for the PHP file to execute in a header, need to have a major type of text
   AddType text/html .php
   AddHandler application/x-httpd-php .php
   Options -Indexes

   HeaderName /.header.htm
   ReadmeName /.footer.php
```

And now I have my footer in PHP :-). You can follow the explanations on [Apache's website](https://httpd.apache.org/docs/2.0/mod/mod_autoindex.html#headername). You can also use CGI, etc.

### Enable PHP Compression

PHP5 compression will save us precious seconds on page display. To enable it, edit the following file and set the parameter to "on":

```apache
; Transparent output compression using the zlib library
; Valid values for this option are 'off', 'on', or a specific buffer size
; to be used for compression (default is 4KB)
; Note: Resulting chunk size may vary due to nature of compression. PHP
;   outputs chunks that are few hundreds bytes each as a result of
;   compression. If you prefer a larger chunk size for better
;   performance, enable output_buffering in addition.
; Note: You need to use zlib.output_handler instead of the standard
;   output_handler, or otherwise the output will be corrupted.
; http://php.net/zlib.output-compression
zlib.output_compression = on
```

Then restart your web server for the configuration to be applied.

## Resources
- [Apache documentation on OpenBSD](/pdf/apache_openbsd.pdf)
- [Documentation on reducing Apache's load with lighttpd](/pdf/reduce_apache's_load_with_lighttpd.pdf)
- http://wiki.gcu.info/doku.php?id=unix:apache_mod_rewrite&s=internet%20explorer
- [The Useful Uses Of Mod Rewrite](/pdf/the_useful_uses_of_mod_rewrite.pdf)
- [How To Specify A Custom php ini For A Web Site (Apache2 With mod_php)](/pdf/how_to_specify_a_custom_php_ini_for_a_web_site__apache2_with_mod_php.pdf)
