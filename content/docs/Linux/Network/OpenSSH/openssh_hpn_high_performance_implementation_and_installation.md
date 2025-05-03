---
weight: 999
url: "/OpenSSH_HPN_(High_Performance_Enabled)_\\:_Implémentation_et_installation/"
title: "OpenSSH HPN (High Performance): Implementation and Installation"
description: "A guide on how to implement and install OpenSSH HPN (High Performance Enabled) which removes performance bottlenecks in standard OpenSSH."
categories: ["Linux", "Networking"]
date: "2006-12-27T11:07:00+02:00"
lastmod: "2006-12-27T11:07:00+02:00"
tags: ["OpenSSH", "Network", "Performance", "SSH", "Security"]
toc: true
---

## Introduction

Here's the introduction provided by the [website](https://www.psc.edu/networking/projects/hpn-ssh/):

SCP and the underlying SSH2 protocol implementation in [OpenSSH](https://www.openssh.com) is network performance limited by statically defined internal flow control buffers. These buffers often end up acting as a bottleneck for network throughput of SCP, especially on long and high bandwidth network links. Modifying the SSH code to allow the buffers to be defined at runtime eliminates this bottleneck. We have created a patch that will remove the bottlenecks in OpenSSH and is fully interoperable with other servers and clients. In addition, HPN clients will be able to download faster from non-HPN servers, and HPN servers will be able to receive uploads faster from non-HPN clients. However, the host receiving the data must have a properly tuned TCP/IP stack. Please refer to this tuning page for more information.

The amount of improvement any specific user will see is dependent on a number of issues. Transfer rates cannot exceed the capacity of the network nor the throughput of the I/O subsystem including the disk and memory speed. The improvement will also be highly influenced by the capacity of the processor to perform the encryption and decryption. Less computationally expensive ciphers will often provide better throughput than more complex ciphers.

## Applying the Patch

First, you need to have the source code of [OpenSSH](https://www.openssh.com) which you can download from: http://www.openssh.com/portable.html

For this tutorial, we'll use the up-to-date version, which is 4.5p1. Download the source:

```bash
wget ftp://ftp.fr.openbsd.org/pub/OpenBSD/OpenSSH/portable/openssh-4.5p1.tar.gz
```

Next, download the OpenSSH HPN patch for our version (4.5p1):

```bash
wget http://www.psc.edu/networking/projects/hpn-ssh/openssh-4.5p1-hpn12v14.diff.gz
```

Now decompress the archives:

```bash
tar -xzvf openssh-4.5p1.tar.gz
gzip -d openssh-4.5p1-hpn12v14.diff.gz
```

Now let's patch the source code:

```bash
cd openssh-4.5p1
patch -p1 < ../openssh-4.5p1-hpn12v14.diff
```

If everything went well, the last lines should look like:

```bash
...
patching file session.c
patching file ssh.c
patching file sshconnect.c
patching file sshconnect2.c
patching file sshd.c
patching file sshd_config
patching file version.h
```

## Compilation

Before starting the configuration, you need to install the dependencies (openssl-dev):

```bash
apt-get install libcurl3-openssl-dev make gcc
```

Next, we can start the configuration:

```bash
./configure --bindir=/usr/bin --sbindir=/usr/sbin --sysconfdir=/etc/ssh --with-md5-passwords
```

Add any arguments you may need if necessary.

Compile:

```bash
make
```

Install:

```bash
make install
```

Your installation is now complete, and you have access to OpenSSH's HPN features! :-)

## FAQ

### What About on Dedibox?

Many people have struggled with recompiling SSH. Here's the solution! Install this:

```bash
apt-get install libpam0g-dev
```

Then, during configuration, add this option:

```bash
./configure --bindir=/usr/bin --sbindir=/usr/sbin --sysconfdir=/etc/ssh --with-md5-passwords --with-pam
```
