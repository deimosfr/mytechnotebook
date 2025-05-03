---
weight: 999
url: "/OpenSSH_\\:_Tunneling_VPN/"
title: "OpenSSH : Tunneling VPN"
description: "Guide to setting up OpenSSH tunneling VPN, including server and client configuration for creating secure VPN connections."
categories: ["Linux", "Network", "FreeBSD"]
date: "2010-05-14T22:46:00+02:00"
lastmod: "2010-05-14T22:46:00+02:00"
tags: ["SSH", "VPN", "Networking", "OpenBSD", "Security"]
toc: true
---

## Introduction

Since version 4.3 of OpenSSH, the option to create IP tunnels has been added.

First, you need to check the OpenSSH version on both the server and client.

```bash
ssh -v
```

You need root privileges on both machines. There are operations to perform, both at the configuration and network levels.

## Configuration

### Server

#### OpenBSD

The first thing to do is to tell OpenSSH to authorize tunnels by adding this directive:

```bash
# Enable layer-3 tunneling. Change the value to 'ethernet' for layer-2 tunneling
PermitTunnel point-to-point
```

We need to ensure that forwarding is activated:

```bash
sysctl net.inet.ip.forwarding=1
```

And will be activated at reboot:

```bash
net.inet.ip.forwarding=1
```

Let's create a tun interface:

```bash
ifconfig tun0 create
ifconfig tun0 10.0.0.1 10.0.0.2 netmask 0xfffffffc
```

And again, make the configuration permanent:

```bash
10.0.0.1 10.0.0.2 netmask 0xfffffffc
```

Now we can restart SSH:

```bash
pkill -HUP sshd ; /usr/sbin/sshd
```

You also need to disable privilege separation, or adjust permissions on `/dev/tun`. For simplicity, I've added:

```
UsePrivilegeSeparation no
```

Another solution is to grant read-write permissions to a specific group on `/dev/tun`, which is much simpler and safer.

```bash
chmod :mygroup /dev/tun
```

And of course, be in that group.

You then need to load the tun module:

```bash
modprobe tun
```

And add it to `/etc/modprobe.preload` for loading at next boot:

```bash
echo tun >> /etc/modprobe.preload
```

### Client

On the client side, we also need to add this directive but in the `/etc/ssh/ssh_config` file:

```bash
# Enable layer-3 tunneling. Change the value to 'ethernet' for layer-2 tunneling
PermitTunnel point-to-point
```

Edit the `/etc/network/interfaces` file and add this interface:

```bash
      iface tun0 inet static
      pre-up ssh -S /var/run/ssh-myvpn-tunnel-control -M -f -w 0:0 5.6.7.8 true
      pre-up sleep 5
      address 10.254.254.2
      pointopoint 10.254.254.1
      netmask 255.255.255.252
      up route add -net 10.99.99.0 netmask 255.255.255.0 gw 10.254.254.1 tun0
      post-down ssh -S /var/run/ssh-myvpn-tunnel-control -O exit 5.6.7.8
```

You only need permissions on `/dev/tun`, so either run as root or have write permission on `/dev/tun`, as mentioned above, then do (where client is the server):

```bash
ssh -w any:any client
```

You can look at the `-f` and `-N` options to avoid launching a shell on the remote machine. And of course, the usual options still work (key, tunnel, master/slave).

Then, as root, you can change the IP of the new tun0 interface on the server:

```bash
ifconfig tun0 10.0.0.1
```

On FreeBSD:

```bash
ifconfig tun100 inet 10.0.0.1 10.0.0.2 netmask 255.255.255.255
```

And do the same on the client:

```bash
ifconfig tun0 10.0.0.2
```

or

```bash
ifconfig tun100 inet 10.0.0.2 10.0.0.1 netmask 255.255.255.255
```

Finally, you can now test the ping from the client:

```bash
ping 10.0.0.1
```

The rest is normal interface configuration. You can add routes, a firewall, anything.

However, you should know that TCP connections over TCP (as is the case with SSH) are not recommended, due to the nature of TCP.

## FAQ

### Connection closed by ...

This is generally due to the server struggling. Check that it has the correct DNS settings and that in the configuration (`/etc/ssh/sshd_config`) the LoginGraceTime value is high enough.

### Cannot fork into background without a command to execute

You may encounter this error message:

```
Cannot fork into background without a command to execute
Failed to bring up tun1.
```

To resolve this issue, add the `-N` option to the SSH command.

## Resources
- [Documentation on SSH VPN](/pdf/tunnelling_vpn_ssh.pdf)
- http://www.kernel-panic.it/openbsd/vpn/vpn5.html
- http://www.debian-administration.org/article/Setting_up_a_Layer_3_tunneling_VPN_with_using_OpenSSH
