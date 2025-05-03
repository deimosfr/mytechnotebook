---
weight: 999
url: "/Squid_\\:_Installation_et_configuration_de_Squid/"
title: "Squid: Installation and Configuration of Squid"
description: "Guide for installing and configuring the Squid proxy server on Debian and FreeBSD systems, with example configurations and security best practices."
categories: ["Debian", "Security", "Linux", "FreeBSD", "Network"]
date: "2012-06-06T08:31:00+02:00"
lastmod: "2012-06-06T08:31:00+02:00"
tags: ["Squid", "Proxy", "Cache", "Network", "Security", "FreeBSD", "Debian"]
toc: true
---

![Squid](/images/squid_logo.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 2.7 |
| **Operating System** | FreeBSD 9<br>Debian 6 |
| **Website** | [Squid Website](https://www.squid-cache.org/) |
| **Last Update** | 06/06/2012 |
{{< /table >}}

## Introduction

A Squid server is a proxy server capable of using FTP, HTTP, Gopher, and HTTPS protocols. Unlike conventional proxy servers, a Squid server handles all requests in a single, non-blocking input/output process.

It's free software distributed under the GNU GPL license.

Squid keeps metadata and especially the most frequently used data in memory. It also stores DNS requests in memory, as well as failed requests. DNS requests are non-blocking.

Cached data can be arranged in hierarchies or meshes to use less bandwidth.

Squid is inspired by the Harvest project. It is compatible with IPv6 from version 3 onwards.

## Installation

### Debian

The installation is simple:

```bash
aptitude install squid
```

### FreeBSD

Installation is easy:

```bash
pkg_add -vr squid
```

Then initialize the cache:

```bash
> squid -z
2012/05/29 05:20:40| Creating Swap Directories
```

Then set Squid to start at boot:

```bash {linenos=table}
# Squid
squid_enable="YES"
```

And finally we'll create the most basic configuration file possible:

```bash
cp /usr/local/etc/squid/squid.conf /usr/local/etc/squid/squid.conf.default
grep -v "^#" < /usr/local/etc/squid/squid.conf.default | sed '/^$/d' > /usr/local/etc/squid/squid.conf
```

## Configuration

### Example 1

```bash {linenos=table}
#-------------------------------------------------------------------------------
# Minimum configuration
#-------------------------------------------------------------------------------
acl all src all
acl manager proto cache_object
acl localhost src 127.0.0.1/32
acl to_localhost dst 127.0.0.0/8 0.0.0.0/32
# Squid listening port
http_port 3128

#-------------------------------------------------------------------------------
# Security
#-------------------------------------------------------------------------------
#chroot on                            # Chroot Squid deamon
forwarded_for off                    # Hide source IP
visible_hostname proxy.deimos.fr     # Mask proxy name
httpd_suppress_version_string on     # Hide squid version

#-------------------------------------------------------------------------------
# ACL network definition
#-------------------------------------------------------------------------------
acl wifi_net src x.x.x.x/24  	# Wifi network
acl wan_net src x.x.x.x/24	# Wan local network

#-------------------------------------------------------------------------------
# ACL Ports definition
#-------------------------------------------------------------------------------
acl SSL_ports port 443
acl Safe_ports port 80		    # http
acl Safe_ports port 21		    # ftp
acl Safe_ports port 443		    # https
acl Safe_ports port 70		    # gopher
acl Safe_ports port 210		    # wais
acl Safe_ports port 1025-65535	# unregistered ports
acl CONNECT method CONNECT

#-------------------------------------------------------------------------------
# Specific ACL
#-------------------------------------------------------------------------------
# Apache mod_gzip and mod_deflate known to be broken so don't trust
# Apache to signal ETag correctly on such responses
acl apache rep_header Server ^Apache
broken_vary_encoding allow apache
#We recommend you to use at least the following line
hierarchy_stoplist cgi-bin ? 

#-------------------------------------------------------------------------------
# Allow/Deny access
#-------------------------------------------------------------------------------
# Minimal access
http_access allow manager localhost
http_access deny manager
# Deny requests to unknown ports
http_access deny !Safe_ports
# Deny CONNECT to other than SSL ports
http_access deny CONNECT !SSL_ports
http_access deny to_localhost
# Custom access
http_access allow wifi_net
# And finally deny all other access to this proxy
http_access deny all

#-------------------------------------------------------------------------------
# Internet Cache Protocol
#-------------------------------------------------------------------------------
icp_access allow wifi_net
icp_access deny all

#-------------------------------------------------------------------------------
# Cache properties
#-------------------------------------------------------------------------------
cache_mgr root                     # Email contact in cache die case
# cache_dir ufs Directory-Name Mbytes L1 L2 [options]
cache_dir ufs /var/squid/cache 1024 16 256
maximum_object_size 10240000 KB    # Set maximum file size to be cached
# Cache expiration patterns
refresh_pattern ^ftp:		1440	20%	10080
refresh_pattern ^gopher:	1440	0%	1440
refresh_pattern -i (/cgi-bin/|\?) 0	0%	0
refresh_pattern .		0	20%	4320

#-------------------------------------------------------------------------------
# Performances Tuning
#-------------------------------------------------------------------------------
pipeline_prefetch on    # To boost the performance of pipelined requests

#-------------------------------------------------------------------------------
# Logs
#-------------------------------------------------------------------------------
access_log /var/squid/logs/access.log squid
cache_log /var/squid/logs/cache.log
cache_store_log /var/squid/logs/store.log
coredump_dir /var/squid/cache
buffered_logs on            # Will speed up if there is not a lot of logs
debug_options ALL,1         # Set log level 1 -> 9
```

### Example 2

For configuration, I won't go into details, but here's an overview of a working configuration that is quite restrictive:

```bash {linenos=table}
http_port 3129 
http_port 3128 
icp_port 3131

# Be more anonymous
# That's three pieces of information you may not want to give away:
# - The host name of your proxy server
# - The version of Squid it's running
# - The IP address of the system that's making the request via the proxy
forwarded_for off
visible_hostname proxy.local
httpd_suppress_version_string on

####auth_param basic program /usr/lib/squid/squid_ldap_auth -b dc=openldap,dc=mycompany,dc=lan -f 'uid=%s' -s sub ldap
#auth_param basic program /usr/lib/squid/squid_ldap_auth -v 3 -b dc=openldap,dc=mycompany,dc=lan -f "(&(objectClass=mycompanyUser)(uid=%s))" -s sub -H ldap://ldap
auth_param basic credentialsttl 2 hours
#auth_param basic realm Web-Proxy

#acl Authentified proxy_auth REQUIRED
acl all	src 0.0.0.0/0
acl mycompany dstdomain mycompany.net
acl mycompany dstdomain mycompany.com
acl mycompany dstdomain mycompany.lan

# Concurent access
#url_rewrite_concurrency 20

# White and Black lists
#acl url_blacklist dstdom_regex -i "/etc/squid/bidon.txt" 
#acl good_domains dstdom_regex -i "/etc/squid/good_domains"
acl url_whitelist dstdom_regex -i "/etc/squid/url_whitelist.txt" 
acl url_blacklist dstdom_regex -i "/etc/squid/url_blacklist.txt" 
acl dst_whitelist dst "/etc/squid/dst_whitelist.txt"
acl dst_blacklist dst "/etc/squid/dst_blacklist.txt" 

# Facebook
acl facebook dstdom_regex facebook.com
acl srcfacebook src 10.101.4.253/32
acl srcfacebook src 10.101.0.253/32 
acl srcfacebook src 10.101.0.107/32
acl srcfacebook src 10.101.4.23/32 

# SuperUser
acl SuperUser src 10.101.4.253/32

# Monster
acl monsternet	dst 208.71.196.0/24
acl monsternet	dst 63.112.169.0/24
acl monsterdom  dstdom_regex newjobs.com monster.com

# Bypass proxy
acl binaries_ext  url_regex -i \.iso$ \.zip$ \.deb$ \.rpm$ \.gz$ \.bz2$ \.exe$ \.cab$ \.bin$ \.tgz$ \.msi$ \.sh$
acl binaries_mime req_mime_type -i ^application/x-debian-package$ ^application/x-bzip2$

acl tunnelurl     url_regex ^http://.*/IDLE/[0-9]+$
acl tunnelurl     url_regex ^http://.*/SEND/[0-9]+$
acl tunnelmethod  method POST
acl videoreq req_mime_type -i ^video/x-ms-wmv$
acl audioreq req_mime_type -i ^audio/mpeg$
acl tunnelreq req_mime_type -i ^application/x-fcs$
acl videorep rep_mime_type -i ^video/x-ms-wmv$
acl audiorep rep_mime_type -i ^audio/mpeg$
acl tunnelrep rep_mime_type -i ^application/x-fcs$
acl manager proto cache_object
acl localhost src 127.0.0.1/32
acl to_localhost dst 127.0.0.0/8
acl localnet src 10.0.0.0/8    # RFC1918 possible internal network
acl localnet src 172.16.0.0/12 # RFC1918 possible internal network
acl localnet src 192.168.0.0/16        # RFC1918 possible internal network

# SSL Ports
acl SSL_ports port 443
acl SSL_ports port 5050	# yahoo

# Allow outbound ports
acl Safe_ports port 80		# http
acl Safe_ports port 8080	# http
acl Safe_ports port 11371	# gpg key
acl Safe_ports port 21		# ftp
acl Safe_ports port 21		# ftp

acl CVS_port port 2401		# CVS
acl CONNECT method CONNECT

# Bypass hours
acl timeok time 18:00-23:59
acl timeok time 00:00-09:00
acl timeok time 12:00-14:00
acl timeok time AS

acl visio src 10.101.3.17/32
acl visio src 10.101.3.9/32

# Google talk ACL
acl	port_gtalk	port 5222
acl gtalk_dst dstdom_regex talk.*.google.com

# ulbridge and co
acl	port_63007	port 63007
acl	port_63005	port 63005
acl	port_8090	port 8090
acl	port_2000	port 2000
acl	snutulbrlbnuat01 dstdom_regex 	snutulbrlbnuat01
acl	cnutulomlnprd-om01 dstdom_regex cnutulomlnprd-om01
acl cnutulbrlnprd-br01 dstdom_regex cnutulbrlnprd-br01
acl	public_ip	dst 62.23.35.231/32

http_access allow Superuser
http_access allow timeok
http_access allow url_whitelist
http_access allow dst_whitelist
http_access allow srcfacebook facebook
http_access deny url_blacklist
http_access deny dst_blacklist
http_access deny tunnelurl tunnelmethod
http_access deny audioreq
http_access deny videoreq
http_access deny tunnelreq

http_reply_access allow timeok
http_reply_access allow url_whitelist
http_reply_access deny audiorep
http_reply_access deny videorep
http_reply_access deny tunnelrep

http_access allow snutulbrlbnuat01 port_63005
http_access allow snutulbrlbnuat01 port_63007
http_access allow cnutulomlnprd-om01 port_63007
http_access allow cnutulbrlnprd-br01 port_63005
http_access allow stmartin port_8003
http_access allow public_ip port_2000

# Google talk
http_access allow gtalk_dst port_gtalk

http_access allow manager localhost
http_access deny manager
http_access allow localhost
#http_access allow localnet Safe_ports mycompany Authentified
#http_access deny localnet Safe_ports mycompany !Authentified
#http_access allow localnet Authentified
#http_access allow localnet CONNECT Authentified
http_access allow localnet Safe_ports 
http_access allow localnet CVS_port
#Authentified
http_access allow localnet SSL_ports CONNECT 
http_access allow localnet CVS_port CONNECT 
#Authentified
#http_access deny !Safe_ports
#http_access deny CONNECT !SSL_ports
http_access deny all

icp_access deny all
htcp_access deny all

cache deny monsternet
cache deny monsterdom

deny_info http://proxy.mycompany.lan/not-allowed-new.html all
deny_info http://proxy.mycompany.lan/not-allowed-new.html url_blacklist.txt
deny_info http://proxy.mycompany.lan/not-allowed-new.html dst_blacklist.txt

#hierarchy_stoplist cgi-bin ?
acl QUERY urlpath_regex cgi-bin \?
no_cache deny QUERY

# comment it out to desactivate squidGuard
url_rewrite_program /usr/bin/squidGuard -c /etc/squid/squidGuard.conf

refresh_pattern ^ftp:		1440	20%	10080
refresh_pattern ^gopher:	1440	0%	1440
refresh_pattern (cgi-bin|\?)	0	0%	0
#refresh_pattern -i \.(gif|jpg|avi|iso|txt)$ 60 20% 120
refresh_pattern .		0	20%	4320
#refresh_pattern -i \.(gif|jpg|avi|iso|txt)$ 30 20% 60

coredump_dir /var/spool/squid
access_log /var/log/squid/access.log squid

cache_dir ufs /var/spool/squid 25000 16 256
maximum_object_size 1024000 KB

delay_pools 3

delay_class 1 2
delay_parameters 1 -1/-1 -1/-1

delay_class 2 3
##delay_parameters 2 256000/256000 64000/64000 16000/48000
delay_parameters 2 256000/256000 64000/64000 -1/-1
#delay_parameters 2 -1/-1 256000/1280000 -1/-1
#delay_parameters 2 256000/256000 -1/-1 -1/-1

delay_class 3 1
delay_parameters 3 256000/256000
#delay_parameters 3 1024000/1024000


delay_access 1 deny localnet
delay_access 3 allow binaries_ext
delay_access 3 allow visio
delay_access 2 deny  binaries_ext
delay_access 2 allow localnet


#debug_options ALL,1 29,6 28,6
ignore_expect_100 on
```

Now all that's left is to adapt and restart the server.

### Verify your rules

#### Get

For Get methods, it's possible to check what's available this way:

```bash
printf "GET http://<destination_ip>:<destination_port>\r\n" |  nc -w 1 <proxy> <proxy_port>
```

#### Connect

You should limit CONNECT methods as much as possible and only allow the GET method. CONNECT is generally used by port 443 and is potentially dangerous because it allows tunneling. To exploit a tunnel, here's how:

```bash
nc -w 1 -v -X connect -x <proxy>:<proxy_port> <destination> <destination_port>
```

We can see if it works (and if a port is open on the other side) with a result like:

```
Connection to 88.191.144.199 873 port [tcp/rsync] succeeded!
```

or

```
nc: Proxy error: "HTTP/1.0 504 Gateway Time-out"
```

At this point, we know that the CONNECT mode works on this port. Otherwise, we'll have:

```
nc: Proxy error: "HTTP/1.0 302 Moved Temporarily"
```

## FAQ

### WARNING! Your cache is running out of filedescriptors

This happens when squid hits the max ulimit. This manifests as major slowdowns on the Internet. To solve this problem, simply increase the size of the file descriptors (default is 1024):

```bash
SQUID_MAXFD=4096
```

Then restart the service.

## Resources
- http://www.squid-cache.org/
- http://www.visolve.com/squid/squid27/contents.php
- http://www.visolve.com/squid/squid27/accesscontrols.php
- http://www.davidandrzejewski.com/2011/11/06/squid-proxy-make-outgoing-headers-anonymous/?utm_source=feedburner&utm_medium=feed&utm_campaign=Feed%3A+davidandrzejewski+%28David+Andrzejewski%29
