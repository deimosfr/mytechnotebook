---
weight: 999
url: "/OpenVPN_\\:_Mise_en_place_d'OpenVPN_sur_différentes_plateformes/"
title: "OpenVPN: Setting up OpenVPN on different platforms"
icon: "openvpn"
icontype: "simple"
description: "A comprehensive guide on installing, configuring and using OpenVPN across different operating systems including Debian, FreeBSD, OpenBSD, Windows, macOS and Linux."
categories: ["Debian", "Security", "Linux"]
date: "2013-05-30T15:26:00+02:00"
lastmod: "2013-05-30T15:26:00+02:00"
tags: ["OpenVPN", "VPN", "Tunnel Blick", "Security", "Network"]
toc: true
---

![OpenVPN](/images/openvpn_logo.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | OpenVPN 2.x |
| **Operating System** | Debian 6<br />OpenBSD 5.9<br />FreeBSD 9 |
| **Website** | [OpenVPN Website](https://openvpn.net/) |
| **Last Update** | 30/05/2013 |
{{< /table >}}

## Introduction

OpenVPN is currently the best solution available for software-based VPN.

OpenVPN allows peers to authenticate each other using a pre-shared private key, certificates, or username/password combinations. It heavily utilizes the OpenSSL authentication library and the SSLv3/TLSv1 protocol. Available on Solaris, OpenBSD, FreeBSD, NetBSD, Linux (Debian, Redhat, Ubuntu, etc...), Mac OS X, Windows 2000, XP, Vista and 7, it also offers numerous security and control functions.

OpenVPN is not compatible with IPsec or other VPN software. The software consists of an executable for both client and server connections, an optional configuration file, and one or more keys depending on the chosen authentication method.

## Installation

### Debian

For Debian or Debian-like distributions (such as Ubuntu), it's very simple:

```bash
aptitude install openvpn
```

In addition to this documentation, take a look at `/etc/default/openvpn`. There are some very interesting options there.

### FreeBSD

On FreeBSD, it's also simple:

```bash
pkg_add -vr openvpn
```

### OpenBSD

On OpenBSD:

```bash
pkg_add -iv openvpn
```

## Configuration

### Authentication via credentials

Dual authentication using PAM adds a security level to certificate authentication. This solution can be used for a deployment using a single certificate shared among all users while still having an authentication method.

Finally, PAM authentication can be used to manage users in an LDAP database.

All my configuration was done on OpenBSD. I'll still try to document for Linux as well based on what I've found.

#### Server

Let's install the appropriate package (as it's a plugin not integrated into OpenVPN today):

```bash
pkgadd -iv openvpn_bsdauth
```

- Add these lines to the server configuration for credential authentication on Linux:

```bash
# OpenVPN PAM Auth
plugin /usr/lib/openvpn/openvpn-auth-pam.so common-auth
```

or

```bash
plugin openvpn-auth-passwd.so _openvpnusers
```

\_openvpnusers: corresponds to the name of the group with connection rights

- If you are on BSD:

```bash
# OpenVPN BSD Auth
auth-user-pass-verify /usr/local/libexec/openvpn_bsdauth via-file
```

Then add the people you want to be able to connect in the '\_openvpnusers' group.

Restart the OpenVPN server.

#### Client

- Login with prompt

At the client level, here's what you need to add:

```bash
# OpenVPN PAM/BSD authentication
auth-user-pass
```

Now, when you try to launch your connection, it will ask you for a login and password.

- Automatic login

If you want to have an automatic login and password, you'll need to put the login and password in a file like this:

```bash
login
password
```

Make sure that you are the only one with rights to this file:

```bash
chmod 600 auth.conf
```

And add this to client.conf:

```bash
# OpenVPN PAM/BSD authentication
auth-user-pass auth.conf
```

Now, launch the connection and nothing will be asked of you.

### Authentication with keys

#### Server

Here we'll create certificates for authentication. We'll need to create a root certificate, then certificates for the clients. Edit the following file and adapt it to your configuration:

```bash
 ...
 export KEY_SIZE=1024
 export CA_EXPIRE=36500
 export KEY_EXPIRE=36500
 export KEY_COUNTRY="FR"
 export KEY_PROVINCE="PA"
 export KEY_CITY="Paris"
 export KEY_ORG="Deimos-Corp"
 export KEY_EMAIL="xxx@mycompany.com"
 ...
```

Here, I set the key expiration to 10 years so I don't have to regenerate keys too often. Next, navigate to the OpenVPN documentation folder to find all the scripts that will allow you to generate certificates:

```bash
# FreeBSD: /usr/local/share/easy-rsa
cd /usr/share/doc/openvpn/examples/easy-rsa/2.0/
. ./vars
./clean-all
./build-ca
./build-key-server server
./build-dh
cd keys
openvpn --genkey --secret ta.key
cp *.key *.c* *.pem /etc/openvpn
```

Replace server with the name of your server where openvpn is installed.

#### Generate client certificates

Still on the server side, for clients, proceed like this for each of them:

```bash
./build-key deimos
mv keys/deimos* /etc/openvpn
```

### IP reservation for clients

Add this to the openvpn config if you want to make IP reservations:

```bash
#Directory containing the configuration for each client (e.g., fixed IP address)
client-config-dir /etc/openvpn/clients
```

Don't forget to create the /etc/openvpn/clients directory.

Then, you need a configuration file per client (/etc/openvpn/clients/<CN indicated on the client certificate>)
Let's take the example of machine srv1 (file /etc/openvpn/clients/srv1):

```bash
ifconfig-push 10.8.0.50 10.8.0.51
```

(The address 10.8.0.51 is used as a "Peer Point" for the OpenVPN server)

## Client configuration

Here are the different types of possible configurations.

### Windows

First, you need to download [OpenVPN GUI](https://openvpn.se/). Then, place the keys in **C:\Program Files\OpenVPN\Config** and the configuration file as well. But it must be **renamed to config.ovpn** (or xxxx.ovpn)

```bash
client
dev tun
proto udp

remote @IP 5000

resolv-retry infinite
nobind

tls-client

persist-key
persist-tun

ca ca.crt
cert deimos.crt
key deimos.key

comp-lzo

verb 1
```

Now check the OpenVPN service if you want a permanent connection, or use the GUI.

### Mac OS X

On Mac, the GUI is [Tunnel Blick](https://www.tunnelblick.net/). You need to apply this type of configuration and place the keys in the right location:

```bash
client
dev tun
proto udp

remote @IP 5000

resolv-retry infinite
nobind

tls-client

persist-key
persist-tun

ca /Users/deimos/Library/openvpn/ca.crt
cert /Users/deimos/Library/openvpn/deimos.crt
key /Users/deimos/Library/openvpn/deimos.key

status /Users/deimos/Library/openvpn/openvpn-status.log

comp-lzo

verb 3
```

All that remains is to "Connect openvpn" on the tunnel placed in the top right.

### Linux

On Linux, you just need to install openvpn and apply the client configuration identical to that of Windows.

Then, to launch the client:

```bash
openvpn --config /home/deimos/.openvpn/client.conf
```

For the VPN connection to auto-establish, the config file must be named with the same name as the certificates.

If you want a graphical client on Linux, Ubuntu has its own Network Manager and with a small plugin installed, it handles it very well:

```bash
apt-get install network-manager-openvpn
```

## My configuration

Because it's not always easy to see what a working configuration looks like, here are mine.
Warning for people who want to try quickly without reading the documentation: I use dual authentication. Remove the lines that don't interest you after reading the documentation.

### Server

I launch my server with these lines:

```bash
up
!/usr/local/sbin/openvpn --config /etc/openvpn/server.conf --tmp-dir /tmp --daemon --script-security 2 > /dev/null 2>&1
```

```bash
local 192.168.10.254
port 1194
proto udp
dev tun0
ca /etc/openvpn/ca/ca.crt
cert /etc/openvpn/srv/mufasa.crt
key /etc/openvpn/srv/mufasa.key
dh /etc/openvpn/dh1024.pem
server 192.168.20.0 255.255.255.0
ifconfig-pool-persist /etc/openvpn/ipp.txt
push "route 192.168.0.0 255.255.255.0"
push "route 192.168.100.0 255.255.255.0"
push "route 192.168.200.0 255.255.255.0"
push dhcp-option "DNS 192.168.100.3"
push dhcp-option "DOMAIN deimos.fr"
client-to-client
keepalive 10 120
tls-auth /etc/openvpn/srv/ta.key 0
auth-user-pass-verify /usr/local/libexec/openvpn_bsdauth via-file
comp-lzo
persist-key
persist-tun
status openvpn-status.log
verb 3
```

### Client

This configuration is used to have a fully automatic connection:

```bash
client
dev tun
proto tcp
remote mufasa.deimos.fr 1194
resolv-retry infinite
nobind
tls-client
persist-key
persist-tun
keepalive 10 120
ca /etc/openvpn/ca.crt
cert /etc/openvpn/shenzi.crt
key /etc/openvpn/shenzi.key
tls-auth /etc/openvpn/ta.key 1
status /etc/openvpn/openvpn-status.log
comp-lzo
verb 3
auth-user-pass /etc/openvpn/auth.cfg
auth-retry nointeract
```

For your information, here are the options that need to be configured:

- keepalive: enables automatic reconnection in case of loss
- auth-user-pass: allows you to store your credentials in a file
- auth-retry: allows for no interaction
- auth-nocache: this directive is deliberately not included. If you include it, the credentials will be dropped from memory after the first connection and at the first disconnection, no automatic reconnection will work. This usually results in a message like: "ERROR: could not read Auth username from stdin".

## FAQ

### WARNING: No server certificate verification method has been enabled

This line is simply missing from your client configuration:

```bash
tls-client
```

### Revoking a certificate

In case of compromise of one of the clients, it is important to know how to revoke its certificate to block access to the OpenVPN server.
It is possible to block access to a client using easy-rsa (still positioning yourself in the easy-rsa directory):

```bash
./revoke-full client2
```

Then you just need to copy the revocation list (keys/crl.pem) to the /etc/openvpn/server/ directory and specify to the OpenVPN server to check the revocation list by adding the line:

```bash
crl-verify /etc/openvpn/server1/crl.pem
```

to the server configuration file (/etc/openvpn/server.conf).

You then need to restart the OpenVPN server.

### Bypassing proxies

You may be at work or school where only ports 80 and 443 are open (bad). Additionally, some sites are blocked. To get around this, the server needs to use port 443 to establish the tunnel. Port 80 might not work due to certain restrictions (port 443 uses the CONNECT method due to SSL, while port 80 works in GET and POST mode). I specify that this is how it works if the proxy is correctly configured (and not in a "Swiss cheese" mode, otherwise it goes through on port 80 if the CONNECT mode is enabled).

#### Server

You simply need to modify the connection type (replace UDP with TCP) and automatically change its default gateway with these 2 options:

```bash
...
proto tcp
port 443
push "redirect-gateway"
...
```

Then restart the server :-)

#### Client

Here we will define proxy rules and go through port 443:

```bash
...
remote fire.deimos.fr 443
# Enter here the proxy name or IP with its port
http-proxy srv-proxy 3128
# Reinitialize connections through the proxy in case of loss or failure
http-proxy-retry
# Pretend to be a web client
http-proxy-option AGENT Mozilla/5.0+(Windows;+U;+Windows+NT+5.0;+en-GB;+rv:1.7.6)+Gecko/20050226+Firefox/1.0.1
...
```

With all this, we're good to go :-)

## FAQ

### Advanced routing impossible

If like me you want to do a bit of complex routing on OpenVPN, you absolutely must change your TUN interfaces to TAP. Why? Simply because you're on layer 3 with TUN and layer 2 with TAP.

- On OpenBSD, you need to do it like this:

```bash
...
dev tun0
dev-type tap
...
```

- On Linux:

```bash
...
dev tap0
...
```

### Making OpenVPN work in an OpenVZ VE

If you want to run an OpenVPN server in a VE, add these types of rights and create the necessary devices:

```bash
vzctl set $my_veid --devices c:10:200:rw --save
vzctl set $my_veid --capability net_admin:on --save
vzctl exec $my_veid mkdir -p /dev/net
vzctl exec $my_veid mknod /dev/net/tun c 10 200
vzctl exec $my_veid chmod 600 /dev/net/tun
vzctl set $my_veid --devnodes net/tun:rw --save
```

## Resources

- [OpenVPN Installation](/pdf/installation_openvpn.pdf)
- [Documentation on a complex OpenVPN setup](/pdf/openvpn_2_howto_fr.pdf)
- [Hardware Authentication for OpenVPN](/pdf/authentification_materielle_pour_openvpn.pdf)
- http://blog.innerewut.de/2005/7/4/openvpn-2-0-on-openbsd
- http://www.openbsd-france.org/documentations/OpenBSD-openvpn.html
- http://www.procyonlabs.com/guides/openbsd/openvpn/index.php
- http://purple.monk.free.fr/phiva/?p=90
- http://www.imped.net/oss/misc/openvpn-2.0-howto-edit.html
- http://auth-passwd.sourceforge.net/
