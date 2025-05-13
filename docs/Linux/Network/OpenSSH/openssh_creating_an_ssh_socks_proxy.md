---
title: "OpenSSH: Creating an SSH SOCKS Proxy"
slug: openssh-creating-an-ssh-socks-proxy/
description: "Learn how to create a SOCKS proxy using SSH to securely route your traffic through an encrypted tunnel."
categories: ["Linux", "Network", "Security"]
date: "2012-02-18T12:24:00+02:00"
lastmod: "2012-02-18T12:24:00+02:00"
tags: ["SSH", "Proxy", "SOCKS", "Network", "Security", "OpenSSH", "Tunneling"]
---

![SSH Socks](../../../static/images/ssh_socks.avif)

## Introduction

This tutorial will be brief, but it's highly effective. The utility of creating a SOCKS proxy via SSH is to be able to route any traffic through an external connection once the SSH connection is established. You simply use the proxy that SSH creates and you're ready to go.

With [SSLH](../../../Servers/File Sharing/SFTP and FTP/sslh_multiplexing_ssl_and_ssh_connections_on_the_same_port.md) as a frontend, you have an almost ultimate tool.

For more advanced techniques, I also recommend checking out [the documentation on proxychains](../proxychains_proxy_any_outbound_connection.md).

## Usage

To establish an SSH connection while opening a SOCKS proxy, simply run this command from server A:

```bash
ssh -D <port> <user>@<destination>
```

For example:

```bash
ssh -D 12345 user@serverB
```

Once the connection is established, configure your web browser or other applications to use localhost as a SOCKS proxy on the specified port (in this case 12345).
