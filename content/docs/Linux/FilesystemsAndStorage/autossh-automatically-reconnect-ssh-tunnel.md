---
weight: 999
url: "/Autossh_\\:_reconnecter_automatiquement_un_tunnel_SSH/"
title: "AutoSSH: Automatically Reconnect SSH Tunnels"
description: "Learn how to use AutoSSH to maintain persistent SSH tunnels that automatically reconnect if the connection drops."
categories: ["Linux", "SSH", "Network"]
date: "2013-03-20T10:44:00+02:00"
lastmod: "2013-03-20T10:44:00+02:00"
tags: ["SSH", "AutoSSH", "Tunnel", "Security"]
toc: true
---

![AutoSSH](/images/openssh_logo.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 1.4 |
| **Operating System** | Debian 7 |
| **Website** | [AutoSSH Website](https://www.harding.motd.ca/autossh/) |
| **Last Update** | 20/03/2013 |
| **Others** | |
{{< /table >}}

## Introduction

If like me you need to maintain a permanent tunnel, and you want it to automatically reconnect in case of disconnection, you need to use a tool like AutoSSH. I'll skip the details of its complete operation, but you should know that it can work in 2 ways:

1. By establishing a loop with the remote server and regularly checking that data is flowing through it
2. By querying a service on the remote machine at regular intervals

## Installation

To install it, it's simple:

```bash
aptitude install autossh
```

## Configuration

Personally, I don't use the method that regularly listens to a port; the other one is quite sufficient for me. So I'll cover the basic method. Let's imagine that I usually use an SSH connection by creating a socks and a local forward:

```bash
ssh -D12345 -N -f -L2222:10.0.0.1:22 deimos@server.deimos.fr
```

To use it with autossh, it's simple, just use all the SSH options and paste them after the autossh command:

```bash
autossh -M 0 -D12345 -N -f -L2222:10.0.0.1:22 deimos@server.deimos.fr
```

The -M 0 option allows you to not use the monitoring option (solution #1). And that's it, autossh manages this connection for you.

### At Boot

If you want to enable it at boot, here's an example to add to rc.local:

```
# /etc.rc.local
autossh -M 0 -q -f -N -oServerAliveInterval=60 -oServerAliveCountMax=3 -L2222:10.0.0.1:22 deimos@server.deimos.fr
```

The SSH options ServerAliveInterval and ServerAliveCountMax allow the connection to be cut if there's a problem to force autossh to restart it.

## References

[https://www.harding.motd.ca/autossh/README](https://www.harding.motd.ca/autossh/README)
