---
weight: 999
url: "/_Echange_de_clefs_SSH/"
title: "OpenSSH: SSH Key Exchange"
description: "How to set up and use SSH key exchange for passwordless authentication, including basic and advanced configurations, ssh-add usage, and troubleshooting common issues."
categories: ["Networking", "Linux", "Security"]
date: "2013-10-25T09:02:00+02:00"
lastmod: "2013-10-25T09:02:00+02:00"
tags: ["SSH", "Authentication", "Security", "Linux", "Debian"]
toc: true
---

## Introduction

SSH key exchange is great for logging in without having to type your password. It's also very simple to set up.

## Basic

### Server

On the server, the user account to which the client will connect (for example **root**) must have the .ssh directory present:

```bash
mkdir .ssh
```

### Client

On the client, you need to generate a key pair, **unless you already have one (~/.ssh/id_dsa.pub)**:

```bash
ssh-keygen -t rsa
```

{{< alert context="info" text="For better security, encrypt your private key during creation by adding -p. If your keys are stolen, it won't be as severe:" />}}

```bash
ssh-keygen -t rsa -p
```

Then, you need to send the key to the server:

```bash
cat .ssh/id_rsa.pub | ssh **root**@**remote_host** "cat >> .ssh/authorized_keys"
```

or

```bash
ssh-copy-id -i ~/.ssh/id_rsa.pub root@remote-server
```

Now, if we connect to the server, we won't be prompted for a password:

```bash
ssh **root**@**remote_host**
```

### Change ssh key passphrase

You can change your ssh passphrase:

```bash
ssh-keygen -f ~/.ssh/id_rsa -p
```

## Complex with restrictions

If, for example, you don't want root to be accessible from anywhere, you need to perform a basic key exchange between the client machine and the server, then edit the following on the server:

* The OpenSSH server configuration file `/etc/ssh/sshd_config`:

```bash
PermitRootLogin without-password
```

* The authorized key file `/root/.ssh/authorized_keys`:

```bash
from="10.0.0.1" ssh-dss..... (the key in question)
```

Finally, restart OpenSSH :). Now, only the machine at 10.0.0.1 will be authorized to connect directly as root and only via the key.

If you have multiple machines or hosts to add, separate them with commas.

## ssh-add

ssh-add[^1] is a tool that allows you to have an SSH private key with a passphrase and not have to type it each time, but simply once during the first use. There is also an X counterpart called ssh-askpass. It's also possible to define a timeout:

```bash
ssh-agent
ssh-add -t 3600
```

So here, after an hour, you'll need to enter the passphrase again.

## FAQ

### Authentication refused: bad ownership or modes for directory

If you encounter this type of error:

```bash
Authentication refused: bad ownership or modes for directory /home/client
Failed publickey for client from x.x.x.x port 57113 ssh2
Connection closed by x.x.x.x
```

You have permission problems in your user's home directory. Check that it has permissions like 755. If it's not possible to change the permissions, then you need to tell SSH to be less restrictive about permissions. You need to modify the file `/etc/ssh/sshd_config` and add this option:

```bash
StrictModes no
```

### Disabling protocol version 1. Could not load host key

I had this small issue, particularly with the Xen Enterprise live CD for performing P2Vs. I wanted to connect remotely to check the progress of the migration. I needed to generate SSH keys to start the server. Here's the procedure:

```bash
ssh-keygen -t rsa1 -f /etc/ssh_host_rsa_key -N ""
ssh-keygen -t dsa -f /etc/ssh_host_dsa_key -N ""
/usr/sbin/sshd
```

And there you go, the problem of the server with missing keys is resolved :)

### I can't change the root password and I absolutely want to connect to the machine

Be careful with this technique because anyone will be able to connect. But for the more adventurous among you, modify these parameters:

```bash
PermitRootLogin yes
StrictModes no
PermitEmptyPasswords yes
```

Restart your SSH service, and there you go, your server is now completely insecure :)

## Resources
- [How To Set Up SSH With Public-Key Authentication](/pdf/how_to_set_up_ssh_with_public-key_authentication_on_debian.pdf)

[^1]: http://docstore.mik.ua/orelly/networking_2ndEd/ssh/ch06_03.htm
