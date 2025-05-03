---
weight: 999
url: "/Mise_en_place_d'un_serveur_de_rebond_pour_ses_connections_SSH/"
title: "Setting up an SSH Bouncer Server for Your SSH Connections"
description: "This guide explains how to set up an SSH bouncer server to facilitate connections to machines in a DMZ, using key authentication and automatic mounting with FUSE."
categories: ["Debian", "Linux", "Ubuntu"]
date: "2008-03-08T07:28:00+02:00"
lastmod: "2008-03-08T07:28:00+02:00"
tags: ["SSH", "Security", "Network", "Servers", "SSHFS", "FUSE", "Tunneling"]
toc: true
---

## Installation

For what follows, I'll base this on a standard Kubuntu 7.10 installation. We'll need the following packages:

- sshfs
- tsocks
- afuse

If you choose to use aptitude for your installation, proceed as follows:

```bash
sudo aptitude install fuse tsocks afuse
```

The installation shouldn't pose any insurmountable problems, so I won't elaborate further on this subject.

However, you must ensure that your user (in my case deimos) belongs to the fuse group:

```bash
deimos@deimos-desktop:~$ id -a
uid=1000(deimos) gid=1000(deimos) groups=4(adm),20(dialout),...106(fuse),108(lpadmin),...1000(deimos)
```

If this is not the case, you can add the user with the following command:

```bash
deimos@deimos-desktop:~$ sudo usermod -G fuse deimos
```

or

```bash
deimos@deimos-desktop:~$ adduser deimos fuse
```

Be sure to log out/log in if you're in a graphical session to apply this group change. You'll also need to restart your terminal if you're SSH'ed into your machine. Otherwise, you can use this command if you don't want to log out:

```bash
newgrp fuse
```

### Setting Up

The first step is to set up our means of communication with the bouncer server, particularly setting up authentication keys. To do this, we'll use authentication keys that we'll deposit in the appropriate directory of the bouncer server user:

### Generating the Key

```bash
deimos@deimos-desktop:~$ ssh-keygen -t dsa
Generating public/private dsa key pair.
Enter file in which to save the key (/home/deimos/.ssh/id_dsa):
Enter passphrase (empty for no passphrase):
Enter same passphrase again:
Your identification has been saved in /home/deimos/.ssh/id_dsa.
Your public key has been saved in /home/deimos/.ssh/id_dsa.pub.
The key fingerprint is:
a9:xx:7d:xx:d9:xx:ea:xx:bd:xx:66:xx:98:xx:47:xx deimos@deimos-desktop
```

### Depositing the Key on the Bouncer Server

```bash
deimos@deimos-desktop:~$ cat .ssh/id_dsa.pub | ssh deimos@rebond "cat >> /home/deimos/.ssh/authorized_keys"
The authenticity of host 'rebond (xx.xx.xx.xx)' can't be established.
RSA key fingerprint is 78:xx:ab:xx:d7:xx:26:xx:49:xx:ec:xx:aa:xx:47:xx.
Are you sure you want to continue connecting (yes/no)? yes
Warning: Permanently added 'rebond' (RSA) to the list of known hosts.
deimos@rebond's password:
```

Alternatively, on Debian, you can use this command:

```bash
ssh-copy-id deimos@rebond
```

We're copying the content of the public key we just generated into the list of keys authorized to connect to my account on the bouncer server. Thus, the next time we try to connect to the bouncer machine, I won't have to enter a password:

```bash
deimos@deimos-desktop:~$ ssh deimos@rebond
Last login: Wed Mar 5 09:21:35 2008 from xx.xx.xx.xx
Authorized uses only. All activity may be monitored and reported.
deimos@rebond#
```

## Creating an SSHFS Mount Point

```bash
deimos@deimos-desktop:~$ mkdir test
deimos@deimos-desktop:~$ sshfs deimos@rebond:/tmp test
deimos@deimos-desktop:~$ ls test
croxxLa crxxLa getxxxt psxxta
crouxxiLa get_xxp lxxe
croxxiLa gexxxt PPCxx303
```

Looks like it's working!

## Accessing Our Servers via Bouncer Server

We have now configured a connection to our bouncer server and we can even mount the file system via SSH from the bouncer server locally. The problem is, if we need to access servers behind the bouncer server, we are forced to reconnect to the latter each time to launch the connection:

```bash
deimos@deimos-desktop:~$ ssh -o ConnectTimeout=2 deimos@host-dmz1
ssh: connect to host host-dmz1 port 22: Connection timed out
```

**It's impossible to connect directly to the host-dmz1 server**

```bash
deimos@deimos-desktop:~$ ssh deimos@rebond
Last login: Wed Mar 5 09:45:36 2008 from 10.251.100.134
yperre@rebond#ssh deimos@host-dmz1
deimos@host-dmz1's password:
```

This kind of thing leads to several disadvantages:

- Multiplication of connections on the bouncer machine
- Loss of time and multiplication of operations to connect to your machines. When you have a park of 200 machines to manage, you don't necessarily want to reconnect 50 times a day.

The idea is therefore to reuse the same connection all the time to transit your connections to the DMZ. To do this, we'll use SSH tunnels, particularly the allocation of dynamic connections (option -D).

To do this, let's restart our connection to the bouncer server by adding the '-D 8888' option to create a dynamic port on port 8888 (the dynamic port is actually seen as a SOCKS server):

```bash
deimos@deimos-desktop:~$ ssh -D 8888 deimos@rebond
Could not request local forwarding.
Last login: Wed Mar 5 09:48:38 2008 from 10.251.100.134
deimos@rebond#
```

Note, if you see the following lines:

```
bind: Address already in use
channel_setup_fwd_listener: cannot listen to port: 8888
```

You have 2 possibilities:

- You already have an open connection with a tunnel
- You have a local program that uses port 8888 => Change it!

Note: From now on, I'll talk about SOCKS server rather than dynamic port.

All that is good, but SSH can't use a SOCKS server to connect to our servers. We'll need to find another solution: a 'socksifying' library (phew!)

You have the choice between dante-client and tsocks. My choice fell on tsocks because of its simplicity, but what follows is perfectly usable under dante!

As we saw above, tsocks (under \*Ubuntu) is simply installed via the packaging system. By default, it will offer you a configuration file /etc/tsocks.conf. I suggest you modify it as follows:

```bash
deimos@deimos-desktop:~$ cat /etc/tsocks.conf
#
server = 127.0.0.1
# Server type defaults to 4 so we need to specify it as 5 for this one
server_type = 5
# The port defaults to 1080 but I've stated it here for clarity
server_port = 8888
```

Now we just need to socksify our SSH calls and voila:

```bash
deimos@deimos-desktop:~$ LD_PRELOAD=/usr/lib/libtsocks.so ssh deimos@host-dmz1
deimos@host-dmz1's password:
```

We are now able to access our DMZ server directly from our workstation. Now let's try to combine this with an SSHFS mount encapsulated in an SSH tunnel:

```bash
deimos@deimos-desktop:~$ LD_PRELOAD=/usr/lib/libtsocks.so sshfs deimos@host-dmz1:/usr/local/lib test
deimos@deimos-desktop:~$ ls test
libcharset.a libgcc_s.so.1 libpopt.la
libcharset.la libiconv.la libpopt.so
libcharset.so libiconv.so libpopt.so.0
libcharset.so.1 libiconv.so.2 libpopt.so.0.0.0
libcharset.so.1.0.0 libiconv.so.2.4.0 libstdc++.a
libg2c.a libintl.a libstdc++.la
libg2c.la libintl.la libstdc++.so
libg2c.so libintl.so libstdc++.so.6
libg2c.so.0 libintl.so.3 libstdc++.so.6.0.3
libg2c.so.0.0.0 libintl.so.3.4.0 preloadable_libiconv.so
libgcc_s.so libpopt.a
```

And there we have it, we have direct access locally, in a transparent way, to our files on the DMZ machine. From there, it's entirely possible to copy our files from one server to another by relying on these mount points.

You might say that's already a pretty good situation, but unfortunately I have to tell you that we can do even better: the use of a fuse automount!

## Automount with FUSE

So far, we've seen the following points:

- Using a bouncer server
- Exchange of private/public keys
- Mounting a file system using the SSH protocol
- Connecting to a server through a SOCKS server/tunnel/dynamic port
- Connecting a file system by using a SOCKS server.

We will now focus on mounting partitions automatically using afuse.

To do this, we'll run a command that will take an SSHFS mount template as a parameter.

Here's the command in question:

```bash
afuse -f -o \
mount_template="LD_PRELOAD=/usr/lib/libtsocks.so sshfs deimos@%r:/ %m" -o \
unmount_template="fusermount -u -z %m" ~/sshfs
```

Note that this command will block your terminal. If you wish to run it as a daemon, you'll need to precede it with the nohup command as well as '&' to run it in the background.

Another important note, if you're using a recent distribution (\*ubuntu, debian, mandriva, etc), your distribution will certainly use UTF8 encoding. If you're using an old Unix/proprietary Unix (Solaris 8, AIX 5.x etc), you'll probably have an ISO8859-1 type encoding. You'll probably need to specify the option '-o from_code=ISO8859-1'.

Let's look at the result:

```bash
deimos@deimos-desktop:~$ cd sshfs/
deimos@deimos-desktop:~/sshfs$ ls
deimos@deimos-desktop:~/sshfs$ cd host-dmz1
deimos@deimos-desktop:~/sshfs/host-dmz1$ ls
bin cdrom etc initrd lib media opt root srv tmp var
boot dev home initrd.img lost+found mnt proc sbin sys usr vmlinuz
deimos@deimos-desktop:~/sshfs/host-dmz1$ cd ..
deimos@deimos-desktop:~/sshfs$ ls
host-dmz1
deimos@deimos-desktop:~/sshfs$ cd recette-dmz1
deimos@deimos-desktop:~/sshfs/recette-dmz1$ ls
1 dead.letter HDS lost+found proc root
11 dev home mnt prod sbin
devices kernel mnt2 doc legal noautoshutdown
bin etc lib opt var usr
LICENSE.txt platform
deimos@deimos-desktop:~/sshfs/recette-dmz1$ cd ..
deimos@deimos-desktop:~/sshfs$ ls
host-dmz1 recette-dmz1
```

You are now able to transparently copy between 2 machines that may be on 2 different DMZs from your workstation (or even editing this with emacs or other kate and vi) and all this in a completely transparent way while facilitating access to your DMZ machines.

The security guys will be pleased!

## Resources
- http://linuxfr.org/2008/03/07/23807.html
- http://ensl.free.fr/softrez/faq/faq-9.html#ss9.1
