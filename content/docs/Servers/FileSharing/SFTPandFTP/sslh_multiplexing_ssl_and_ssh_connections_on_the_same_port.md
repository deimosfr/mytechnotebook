---
weight: 999
url: "/SSLH\\:_Multiplexer_les_connections_SSL_et_SSH_sur_le_même_port/"
title: "SSLH: Multiplexing SSL and SSH connections on the same port"
description: "How to configure SSLH to multiplex SSL and SSH connections on the same port to allow both HTTPS and SSH traffic through a single port."
categories: ["Linux", "FreeBSD", "Network"]
date: "2012-06-10T09:31:00+02:00"
lastmod: "2012-06-10T09:31:00+02:00"
tags:
  [
    "SSLH",
    "SSH",
    "SSL",
    "Security",
    "Network",
    "Debian",
    "FreeBSD",
    "OpenBSD",
    "pfSense",
  ]
toc: true
---

## Introduction

[SSLH](https://www.rutschle.net/tech/sslh.shtml) is like a magic trick. It allows you, for example, to have both HTTPS and SSH on a WAN address on port 443. How is this possible? Can we have a single listening port for multiple services?

Indeed! You just need to see SSLH as a layer 7 proxy capable of filtering and differentiating between SSL frames and SSH frames.

## Installation

### Debian

On Debian, it's an easy move:

```bash
aptitude install sslh
```

### FreeBSD

For installation on FreeBSD:

```bash
pkg_add -vr sslh
```

### OpenBSD

On OpenBSD for version 1.7:

```bash
wget http://www.rutschle.net/tech/sslh-1.7a.tar.gz
tar -xzvf sslh-1.7a.tar.gz
cd sslh-1.7a
cc -o sslh -DLIBWRAP sslh.c -lwrap
cp sslh /usr/local/sbin
```

For the latest versions:

```bash
wget http://www.rutschle.net/tech/sslh-1.10.tar.gz
tar -xzvf sslh-1.10.tar.gz
cd sslh-1.10
make
make install
```

_Note: With the latest versions, you'll end up with 2 binaries: sslh-fork and sslh-select. sslh-select is for a single thread while sslh-fork is multi-threaded._

## Configuration

### Debian

If you're on Debian, it's simple. Create a file in `/etc/default/sslh` with the following content:

```bash
LISTEN=ifname:443
SSH=localhost:22
SSL=localhost:443
```

- LISTEN: This is the IP of the interface and the listening port of SSLH. It must be the input listening port (WAN for example), the one through which your clients will pass.
- SSH: The address and port corresponding to your SSH port
- SSL: The address and port corresponding to your HTTPS port

### FreeBSD

To configure SSLH, just add these lines to rc.conf:

```bash
# SSLH
sslh_enable="YES"
sslh_mode="select"
# sslh_fib="NONE"
sslh_pidfile="/var/run/sslh/sslh.pid"
sslh_ssltarget="websrv:443"
sslh_sshtarget="localhost:22"
sslh_sshtimeout="2"
sslh_listening="192.168.0.254:443"
sslh_uid="nobody"
```

### OpenBSD

On OpenBSD, I chose to add these lines to `/etc/rc.local` with my configuration directly in the command line:

```bash {linenos=table,hl_lines=[3,5]}
if [ -x /usr/local/sbin/sslh ] ; then
    # Versions < 1.7a
    /usr/local/sbin/sslh -p <interface:443> -s <ssh_srv:22> -l <web_https:443>
    # Versions >= 1.10
    /usr/local/sbin/sslh -P /tmp/sslh.pid -p <interface:443> --ssh <ssh_srv:22> --ssl <web_https:443>
    echo 'SSLH'
fi
```

Again, adapt these lines according to your needs.

### pfSense

On pfSense, we'll create an init-like file:

```bash {linenos=table,hl_lines=["2-7"]}
#!/bin/sh
sslh_listen_ip=<interface_ip>
sslh_listen_port=<port>
ssh_redirect_ip=<interface_ip>
ssh_redirect_port=<port>
ssl_redirect_ip=<interface_ip>
ssl_redirect_port=<port>

rc_start() {
	/sbin/sslh-fork -P /tmp/sslh.pid -p $sslh_listen_ip:$sslh_listen_port --ssh $ssh_redirect_ip:$ssh_redirect_port --ssl $ssl_redirect_ip:$ssl_redirect_port &
}

rc_stop() {
	/usr/bin/killall sslh-fork
}

case $1 in
	start)
		rc_start
		;;
	stop)
		rc_stop
		;;
	restart)
		rc_stop
		rc_start
		;;
esac
```

Adapt the beginning of the script with your desired information.

Then set the appropriate permissions:

```bash
chmod +x /usr/local/etc/rc.d/sslh.sh
```

## FAQ

### I have Zombie sslh processes on OpenBSD

I've experienced zombie processes with each connection attempt to the SSH server. To fix this issue, here's a patch for version 0.7:

```c
462a463,464
>    if (fork() > 0) exit(0); /* Detach */
>
467,468d468
<    if (fork() > 0) exit(0); /* Detach */
<
```
