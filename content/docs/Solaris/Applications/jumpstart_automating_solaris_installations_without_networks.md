---
weight: 999
url: "/Jumpstart\\:_automatiser_les_installations_Solaris_sans_réseaux/"
title: "Jumpstart: Automating Solaris Installations Without Networks"
description: "A guide on how to automate Solaris installations using Jumpstart without requiring network infrastructure."
categories:
  - "Solaris"
date: "2009-08-21T15:17:00+02:00"
lastmod: "2009-08-21T15:17:00+02:00"
tags:
  - "jumpstart"
  - "solaris"
  - "installation"
  - "automation"
  - "iso"
toc: true
---

## Introduction

For my job, I had to automate Solaris installations. Jumpstart exists for this purpose. The problem is that in a new datacenter, you don't always have what you need. That's the real issue. That's why I researched documentation on the internet, which often doesn't work properly, and I'll try to create a comprehensible guide that works (hopefully).

## Prerequisites

You'll need a fresh Solaris installation on which you'll do minimal configuration. Also install anything else you're interested in.

Note: For now, forget about the ZFS root version and use UFS. There is no clean method for doing a flash install with ZFS.

You'll also need the **SUNWmkcd** package to have the mkisofs command.

## Creating the Flar

We'll create a flar image which will make an archive of the current system:

```bash
mkdir -p /export/home/sol10jumpstart
flarcreate -n sol10_jumpstart -c /export/home/sol10jumpstart/sol10_auto.flar
```

## Copying the DVD

We'll copy the DVD content to make necessary modifications:

```bash
cp -Rf /cdrom/cdrom0 /export/home/dvd
rm -Rf /export/home/dvd/Solaris_10/Product
```

## x86.miniroot

### Unpack

Here we'll unpack x86.miniroot to modify its content:

```bash
/boot/solaris/bin/root_archive unpack /export/home/dvd/boot/x86.miniroot /var/tmp/miniroot
```

### Jumpstart CDROM

Now, to address a small issue, we'll edit this file and you should comment out these 2 lines:

(`/var/tmp/miniroot/usr/sbin/install.d/profind`)

```bash
# Factory JumpStart (default) profile search
# Arguments:    none
#
cdrom()
{
   # Factory JumpStart is only allowed with factory
   # stub images, indicated by the file /tmp/.preinstall
   #
   #if [ -f /tmp/.preinstall ]; then
       mount -o ro -F lofs ${CD_CONFIG_DIR} ${SI_CONFIG_DIR} >/dev/null 2>&1

       if [ $? -eq 0 ]; then
           verify_config "CDROM"
       fi
   #fi
}
```

### sysidcfg

Here's probably the most interesting file for setting up your jumpstart. But first, we need to remove the existing symbolic link that would prevent jumpstart from working properly:

```bash
rm /var/tmp/miniroot/etc/sysidcfg
```

Now, let's create a new file with this content:

(`/var/tmp/miniroot/etc/sysidcfg`)

```bash
timezone=UTC
timeserver=localhost
name_service=NONE
network_interface=primary {
  netmask=255.255.255.0
  protocol_ipv6=no
  default_route=NONE }
nfs4_domain=dynamic
security_policy=NONE
#keyboard=US-English
system_locale=en_US
terminal=vt100
root_password=REPLACE_WITH_YOU_OWN
```

For the password, you need to get the encrypted version from /etc/shadow for example. Here's another example:

(`/var/tmp/miniroot/etc/sysidcfg`)

```bash
name_service=none
root_password=TITJXNq6L24dw
network_interface=none
security_policy=none
system_locale=C
terminal=vt100
timeserver=localhost
```

And here's another example:

(`/var/tmp/miniroot/etc/sysidcfg`)

```bash
name_service=none
timezone=UTC
timeserver=localhost
root_password=OngWELbxoVfUU
network_interface=nge0 {hostname=installtemp default_route=1.1.1.2 ip_address
=1.1.1.3 netmask=255.255.0.0 protocol_ipv6=no}
nfs4_domain=dynamic
security_policy=none
system_locale=C
terminal=vt100
timeserver=localhost
```

### Pack

Now we'll repackage everything:

```bash
/var/tmp/miniroot/boot/solaris/bin/root_archive pack /export/home/dvd/boot/x86.miniroot /var/tmp/miniroot
```

## Moving the flar

Let's move the flar to the folder containing the Solaris DVD with the latest modifications we've made:

```bash
mv /export/home/sol10jumpstart/sol10_auto.flar /export/home/dvd
```

## Customizing the Jumpstart

We'll now choose the automations we want to implement:

```bash
cd /export/home/dvd
rm -Rf .install_config
mkdir .install_config
cd .install_config
```

### any_profile

(`any_profile`)

```bash
install_type flash_install
archive_location local_file /cdrom/sol10_auto.flar
fdisk all solaris all
partitioning explicit
filesys rootdisk.s0 20480 /
filesys rootdisk.s1 4096 swap
filesys rootdisk.s3 10240 /var
filesys rootdisk.s4 1024 /globaldevices
filesys rootdisk.s7 free /export/home
```

### begin

(`begin`)

```bash
#!/bin/sh
echo "Begining ISO FLAR based jumpstart."
```

### finish

(`finish`)

```bash
#!/bin/sh
ROOTDIR=${ROOTDIR:-/a}

#echo "Finish script for Jumpstart FINISH"
#echo "Get rid of the nfs prompt"
touch ${ROOTDIR}/etc/.NFS4inst_state.domain

# TODO: keep exit status, return it, use the first error encountered.
```

### rules

(`rules`)

```bash
probe rootdisk
probe disks
probe karch
probe memsize
probe model
probe hostname
any -  begin any_profile finish
```

Then, we need to verify the entire configuration. Fortunately, a small tool exists (this command is required):

```bash
/export/home/dvd/Solaris_10/Misc/jumpstart_sample/check
```

### Grub

We'll edit the boot menu. Add these lines (they should be placed at the beginning of the title section):

(`/export/home/dvd/boot/grub/menu.lst`)

```bash
title Solaris 10 Jumpstart
      kernel /boot/multiboot kernel/unix - install w -B install_media=cdrom
      module /boot/x86.miniroot
```

## Creating the ISO File

All that's left is to create the ISO file with everything we've done:

```bash
cd /export/home/dvd/
mkisofs -b boot/grub/stage2_eltorito -c .catalog -no-emul-boot -boot-load-size 4 -boot-info-table -relaxed-filenames -l -ldots -r -N -d -D -V SOL10JUMPSTART -o /export/home/mysol10u6x86.iso .
```

Now all you need to do is burn it and boot from it :-)

## FAQ

### /export/home/dvd/boot/x86.miniroot: override protection 444

If you get this error, copy the DVD content to a folder and try your command again (generally):

```bash
/boot/solaris/bin/root_archive pack /export/home/sol_10_1008_x86/boot/x86.miniroot /var/tmp/miniroot
```

## References

- http://wikis.sun.com/display/BigAdmin/Creating+a+bootable+ISO+image
- http://run.tournament.org.il/tag/flar/
- http://docs.sun.com/app/docs/doc/819-5776/6n7r9js2j?a=view
- http://www.sun.com/bigadmin/features/articles/jumpstart_x86_x64.jsp
- http://forums.sun.com/thread.jspa?threadID=5372582&tstart=0
- http://amorin.org/professional/jumpstart.php
- http://docs.sun.com/app/docs/doc/820-2315/ggsez?l=fr&a=view
- http://www.eng.auburn.edu/~doug/howtos/multipathing.html
