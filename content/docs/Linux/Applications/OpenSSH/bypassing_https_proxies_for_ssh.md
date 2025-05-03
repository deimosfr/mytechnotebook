---
weight: 999
url: "/Outrepasser_les_proxy_HTTPS_pour_SSH/"
title: "Bypassing HTTPS Proxies for SSH"
description: "Guide to bypass corporate proxies and allow SSH connections through port 443 when standard ports are blocked"
categories: ["Linux", "Network", "Debian", "Ubuntu"]
date: "2012-07-30T11:19:00+02:00"
lastmod: "2012-07-30T11:19:00+02:00"
tags:
  ["SSH", "Proxy", "HTTPS", "Network", "Security", "Firewall", "Connect-proxy"]
toc: true
---

## Introduction

Those workplace proxies can be really annoying! But there are always solutions!

So here's the situation: I want to access a remote machine via SSH, but only ports 80 and 443 are allowed. Even if you configure the SSH server on port 443, you'll notice it doesn't work.

A solution? Yes: connect-proxy.

## Installation

### Seveur

On the server, simply modify the sshd_config file to make SSH listen on port 443:

```bash
Port 443
```

And restart the SSH service.

PS: If you don't want to run SSH on port 443, you can use [SSLH method]({{< ref "docs/Servers/FileSharing/SFTPandFTP/sslh_multiplexing_ssl_and_ssh_connections_on_the_same_port.md" >}}) to multiplex SSL and SSH on the same port.

### Client

#### Debian / Ubuntu

Install connect-proxy:

```bash
aptitude install connect-proxy
```

#### Mac

Let's compile it:

```bash
cd /private/tmp
wget http://www.meadowy.org/~gotoh/ssh/connect.c
gcc connect.c -o connect -lresolv
sudo cp connect /usr/bin
sudo chmod 555 /usr/bin/connect
sudo chown root:wheel /usr/bin/connect
```

#### Others

You'll need to compile it:

```bash
cd /tmp/
wget http://www.meadowy.org/~gotoh/ssh/connect.c
gcc connect.c -o connect
sudo cp connect /usr/local/bin/
sudo chmod +x /usr/local/bin/connect
```

## Configuration

Create or edit your SSH config file (`~/.ssh/config`):

```bash
## Outside of the firewall, with HTTPS proxy
Host my_ssh_server_i_want_to_reach
ProxyCommand connect -H annoying_proxy:3128 %h 443
## Inside the firewall (do not use proxy)
Host *
  ProxyCommand connect %h %p
```

The configuration is complete, now you just need to connect:

```bash
ssh my_ssh_server
```

## Resources

- http://www.zeitoun.net/articles/ssh-through-http-proxy/start
