---
weight: 999
url: "/Proxy\\:\_Créer_un_proxy_avec_Apache/"
title: "Proxy: Creating a proxy with Apache"
description: "Guide on how to set up and configure an Apache proxy server for different use cases including tunneling SSH over HTTP and setting up a reverse proxy for applications."
categories: ["Linux", "Network", "Servers"]
date: "2012-03-29T07:16:00+02:00"
lastmod: "2012-03-29T07:16:00+02:00"
tags:
  [
    "Apache",
    "Proxy",
    "Tomcat",
    "SSH Tunneling",
    "Reverse Proxy",
    "Network",
    "Web Server",
    "OpenBSD",
    "Debian",
  ]
toc: true
---

## Introduction

With Apache's mod_proxy, there are several use cases. I will propose two scenarios here.

### Scenario 1

Here's the situation! I'm in a computer school where (like in many schools) only port 80 is open, and the class isn't always interesting.

So what can you do to access your SSH server, play World of Warcraft, or download heavily from eMule?

Well, Uncle Tom has a super pattern for you who wants to break the laws: the APACHE MOD_PROXY PLATINUM EDITION!

Here we're working on Debian, but the configuration is essentially the same on other systems as long as you're using Apache2's mod_proxy.

### Scenario 2

In this scenario, I want to redirect incoming traffic on my standard port (80) to an application (on the same machine or not) using URL rewriting. The advantage is that with mod_proxy, there's no need to use RewriteEngine & Co! The proxy module can handle most of the rewriting, especially hiding the port number (useful for applications running on Tomcat).

## Installation

```bash
aptitude install apache2 apache2-utils apache2.2-common libapache2-mod-proxy-html
```

Then activate modules:

```bash
a2enmod proxy_connect
a2enmod proxy_http
a2enmod proxy_html
```

And restart Apache.

## Configuration

### Scenario 1

#### Debian

First, we'll configure the mod_proxy in question.
Here's my detailed `/etc/apache2/mods-available/proxy.conf` file:

```bash
<IfModule mod_proxy.c>
#On autorise les requêtes de type proxy
        ProxyRequests On
#On autorise le serveur à répondre à ces requêtes
        ProxyVia On
#On autorise les requêtes proxy en destination du port 22, 80 et 443
        AllowCONNECT 22
        AllowCONNECT 80
        AllowCONNECT 443
#On autorise le proxy à destination de n'importe quelle adresse
# (Pour restreindre qu'a une seule adresse il faut mettre quelque chose comme
#"<Proxy google.fr>" ou encore "<Proxy 88.191.31.151>")
        <Proxy *>
# Nous allons restreindre l'accès par mot de passe
                AllowOverride AuthConfig
                AuthName "Proxy Auth"
                AuthType Basic
# Le fichier htpasswd à utiliser
                AuthUserFile /etc/apache2/.htpasswd-proxy
# seuls les utilisateurs authentifiés ont accès
                Require valid-user
                AddDefaultCharset off
                Order deny,allow
                Allow from all
        </Proxy>
</IfModule>
```

Next, we create the "htpasswd" file (e.g., for the user toto)

```bash
htaccess -c /etc/apache2/.htpasswd-proxy toto
```

Now we just need to load the modules

```bash
cd /etc/apache2/mods-enabled/
ln -s ../mods-available/proxy.load .
ln -s ../mods-available/proxy.conf .
ln -s ../mods-available/proxy_connect.load .
ln -s ../mods-available/proxy_http.load .
```

Then restart Apache2

```bash
/etc/init.d/apache2 restart
```

#### OpenBSD

With OpenBSD, no specific installation is needed since Apache is installed by default. Just add this to the configuration:

```bash
<VirtualHost _default_:3128>

#  General setup for the virtual host
DocumentRoot /var/www/htdocs
ServerName mufasa.deimos.fr
ServerAdmin xxx@mycompany.com
ErrorLog logs/error_log
TransferLog logs/access_log

#   SSL Engine Switch:
#   Enable/Disable SSL for this virtual host.
SSLEngine off

<IfModule mod_proxy.c>
      ProxyRequests On
      ProxyVia On
      <Directory proxy:*>
              Order deny,allow
              Allow from all
      </Directory>
</IfModule>

</VirtualHost>
```

Then restart the service:

```bash
apachectl stop
apachectl start
```

Obviously, this allows everyone access, so make sure to add some security.

Personally, my Apache is bound to a port that only the local network and people connected via VPN can access.

#### PuTTY: Tunneling SSH

So now we have a nice proxy, but how to make the most of it?

We'll use PuTTY to simplify things, as it's one of the few cross-platform SSH clients that offers all the functions we need: Tunneling + HTTP Proxy.

The principle is as follows:

1. Establish an SSH connection on port 22
2. Go through this proxy server which authorizes connections on port 22
3. Using SSH, we establish encrypted local tunnels that redirect to different services
4. We access the services on localhost through the tunnels

Here's how I configure my PuTTY client to play World of Warcraft:

**_{Session Menu}_**

- Host Name: <SSH server destination>
- Port: 22

**_{Proxy Menu}_**

- Proxy type: HTTP
- Proxy hostname: <proxy server address>
- Port: 80
- Username: <username created in htpasswd>
- Password: <password for this user>

**_{SSH / Tunnels Menu}_**

- Local Ports accept connections from other hosts: ON
- Source port: <local port to open> (ex. 3724)
- Destination: <ip:port of the service you want to forward> (ex. eu.logon.worldofwarcraft.com:3724)

Click on "add" to add others e.g., 5900:vnc; 143:imap; 25:smtp
(for WoW, don't forget this one)

```
Source port: "6112"
Destination: "80.239.185.41:6112"
```

That's good! For World of Warcraft, all that's left is to modify the "realmlist.wtf" file and put:

```
set realmlist localhost
```

As a famous philosopher would say: "And the show begins!"

### Scenario 2

Here I'll use the example of a "myapp" tool running on Tomcat, port 8080. First, I need to tell Tomcat that it will be "proxified," and then I need to set up the proxy part on Apache.

#### Tomcat

On the server side, you'll need to modify the connector for the application in question to add the proxy parameters:

```apache
    <Connector port="8080" protocol="HTTP/1.1"
               connectionTimeout="20000"
               URIEncoding="UTF-8"
               redirectPort="8443"
               proxyName="myapp.mycompany.lan" proxyPort="80"/>
```

We're telling Tomcat that our site will be accessible from myapp.mycompany.lan on port 80.
You can restart your Tomcat now.

#### Apache

We'll activate the proxy module:

```apache
<IfModule mod_proxy.c>

# If you want to use apache2 as a forward proxy, uncomment the
# 'ProxyRequests On' line and the <Proxy *> block below.
# WARNING: Be careful to restrict access inside the <Proxy *> block.
# Open proxy servers are dangerous both to your network and to the
# Internet at large.
#
# If you only want to use apache2 as a reverse proxy/gateway in
# front of some web application server, you DON'T need
# 'ProxyRequests On'.

ProxyRequests Off
ProxyPreserveHost On
<Proxy *>
        AddDefaultCharset off
        Order deny,allow
        Allow from all
        #Allow from .example.com
</Proxy>

# Enable/disable the handling of HTTP/1.1 "Via:" headers.
# ("Full" adds the server version; "Block" removes all outgoing Via: headers)
# Set to one of: Off | On | Full | Block
ProxyVia On

</IfModule>
```

Then configure it for our site. We'll use a VirtualHost for our application:

```apache
<VirtualHost myapp.mycompany.lan:80>
    ServerName http://myapp.mycompany.lan
    ServerAlias myapp.mycompany.lan
    ServerAdmin xxx@mycompany.com
    DocumentRoot /mnt/myapp/datas/www
    ScriptAlias /cgi-bin/ /usr/lib/cgi-bin/
    LogLevel warn
    ServerSignature On

    ProxyPass / http://localhost:8080/
    ProxyPassReverse / http://localhost:8080/

    <Location />
        Order allow,deny
        Allow from all
    </Location>
</VirtualHost>
```

The ProxyPass part tells where to redirect the proxy. Here the Apache proxy and Tomcat are running on the same machine, which is why the URLs point to localhost.

All that's left is to restart Apache, and your service that was originally available at this address:
http://myapp.mycompany.lan:8080/
will be available at:
http://myapp.mycompany.lan/

## Resources
- http://httpd.apache.org/docs/2.0/mod/mod_proxy.html
