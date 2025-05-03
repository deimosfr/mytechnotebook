---

weight: 999
url: "/Installation_FreeBSD_sur_ZFS/"
title: "Installing FreeBSD on ZFS"
description: "A comprehensive guide on how to install FreeBSD with ZFS as the root filesystem, including disk formatting, partitioning, and system configuration."
categories: ["FreeBSD", "Linux"]
date: "2010-09-21T21:23:00+02:00"
lastmod: "2010-09-21T21:23:00+02:00"
tags: ["ZFS", "FreeBSD", "Storage", "Filesystem", "RAID"]
toc: true

---

## Introduction

I love ZFS and for building a large NAS, I need FreeBSD which is capable of implementing ZFS and, most importantly, using it on the root partition as well.

For this purpose, I needed 5 disks of the same capacity and the FreeBSD DVD (I emphasize the DVD because the livefs or CD versions don't contain everything needed to boot from ZFS).

## Disk Formatting

### Creating Partitions

Boot from the FreeBSD DVD and launch the Fixit menu.

Once inside, I recommend checking the names of the devices installed on your system:

```bash
> ls /dev/ad*
/dev/ad10 /dev/ad12 /dev/ad4 /dev/ad6 /dev/ad8
```

Before creating partitions and slices, note that you can check the status of your disks at any time with:

```bash
gpart show ad10
```

Let's define the disks format as GPT:

```bash
gpart create -s gpt ad10
gpart create -s gpt ad12
gpart create -s gpt ad4
gpart create -s gpt ad6
gpart create -s gpt ad8
```

Then we'll create a boot partition:

```bash
gpart add -b 34 -s 128 -t freebsd-boot ad10
gpart add -b 34 -s 128 -t freebsd-boot ad12
gpart add -b 34 -s 128 -t freebsd-boot ad4
gpart add -b 34 -s 128 -t freebsd-boot ad6
gpart add -b 34 -s 128 -t freebsd-boot ad8
```

And a swap partition (4GB for example):

```bash
gpart add -b 162 -s 8388608 -t freebsd-swap -l swap0 ad10
gpart add -b 162 -s 8388608 -t freebsd-swap -l swap1 ad12
gpart add -b 162 -s 8388608 -t freebsd-swap -l swap2 ad4
gpart add -b 162 -s 8388608 -t freebsd-swap -l swap3 ad6
gpart add -b 162 -s 8388608 -t freebsd-swap -l swap4 ad8
```

To calculate the size in cylinders, here's the simple formula in MB with an example for our 4GB:

```
cylinder size = x * 4 * 512
8388608 = 4096 * 4 * 512
```

And finally, the data partitions where the raid-z will reside:

```bash
gpart add -b 8388770 -s 125829120 -t freebsd-zfs -l disk0 ad12
gpart add -b 8388770 -s 125829120 -t freebsd-zfs -l disk1 ad10
gpart add -b 8388770 -s 125829120 -t freebsd-zfs -l disk2 ad4
gpart add -b 8388770 -s 125829120 -t freebsd-zfs -l disk3 ad6
gpart add -b 8388770 -s 125829120 -t freebsd-zfs -l disk4 ad8
```

You should adjust according to the remaining space (see 'gpart show' to know how much space is left).
We install the protected MBR and gptzfsboot on each disk:

```bash
gpart bootcode -b /mnt2/boot/pmbr -p /mnt2/boot/gptzfsboot -i 1 ad12
gpart bootcode -b /mnt2/boot/pmbr -p /mnt2/boot/gptzfsboot -i 1 ad10
gpart bootcode -b /mnt2/boot/pmbr -p /mnt2/boot/gptzfsboot -i 1 ad4
gpart bootcode -b /mnt2/boot/pmbr -p /mnt2/boot/gptzfsboot -i 1 ad6
gpart bootcode -b /mnt2/boot/pmbr -p /mnt2/boot/gptzfsboot -i 1 ad8
```

### Creating the ZFS

We'll need to load the ZFS modules:

```bash
kldload /mnt2/boot/kernel/opensolaris.ko
kldload /mnt2/boot/kernel/zfs.ko
```

And finally we create the raidz:

```bash
mkdir /boot/zfs
zpool create zroot raidz1 /dev/gpt/disk0 /dev/gpt/disk1 /dev/gpt/disk2 /dev/gpt/disk2 /dev/gpt/disk3 /dev/gpt/disk4
zpool set bootfs=zroot zroot
```

Now we'll create the necessary mount points and ZFS options to install the system.
We enable checksum validation on the filesystem:

```bash
zfs set checksum=fletcher4 zroot
```

We create a partition for /tmp with some useful options:

```bash
zfs create -o compression=on -o exec=on -o setuid=off zroot/tmp
chmod 1777 /zroot/tmp
```

And that's it for the ZFS part.

## Installing FreeBSD

We'll decompress part of what's on the DVD into our freshly created zpool:

```bash
cd /dist/8.0-*
export DESTDIR=/zroot
for dir in base catpages dict doc games info lib32 manpages ports; do (cd $dir ; ./install.sh) ; done
cd src ; ./install.sh all
cd ../kernels ; ./install.sh generic
cd /zroot/boot ; cp -Rlp GENERIC/* /zroot/boot/kernel/
```

Let's chroot into our new environment:

```bash
chroot /zroot
```

And configure the rc.conf file:

```bash
zfs_enable="YES"
hostname="freebsd.deimos.fr"
ifconfig_re0="DHCP"
```

And the bootloader file:

```bash
zfs_load="YES"
vfs.root.mountfrom="zfs:zroot"
```

### Configuration

Let's configure the root password:

```bash
passwd
```

The timezone:

```bash
tzsetup
```

The mail aliases:

```bash
cd /etc/mail
make aliases
```

We unmount and exit the chroot:

```bash
umount /dev
exit
```

And we copy the zpool cache:

```bash
cp /boot/zfs/zpool.cache /zroot/boot/zfs/zpool.cache
```

We create the fstab file (`/zroot/etc/fstab`):

```
# Device                       Mountpoint              FStype  Options         Dump    Pass#
/dev/gpt/swap0                 none                    swap    sw              0       0
/dev/gpt/swap1                 none                    swap    sw              0       0
/dev/gpt/swap2                 none                    swap    sw              0       0
/dev/gpt/swap3                 none                    swap    sw              0       0
/dev/gpt/swap4                 none                    swap    sw              0       0
```

We unmount the zpool:

```bash
export LD_LIBRARY_PATH=/mnt2/lib
zfs unmount -a
```

And we configure the mount points of our ZFS:

```bash
zfs set mountpoint=legacy zroot
zfs set mountpoint=/tmp zroot/tmp
```

All that's left is to exit the fixit mode and the sysinstall to reboot.

## FAQ

### RaidZ2 degraded

What happens if you have a disk in degraded mode? First, let's check the status to see what's happening:

```bash
> zpool status -x
  pool: zroot
 state: DEGRADED
status: One or more devices could not be opened.  Sufficient replicas exist for
        the pool to continue functioning in a degraded state.
action: Attach the missing device and online it using 'zpool online'.
   see: http://www.sun.com/msg/ZFS-8000-2Q
 scrub: none requested
...
```

Here we can see the "DEGRADED" status. To put it simply:

* If the machine is running: replace the defective disk with a new one.
* If the machine is off or if you want to shut it down: boot from the FreeBSD DVD, load the ZFS modules.

Then recreate the partitions as you did for the other disks. Next, we'll add the new disk to the raidz2 so it automatically rebuilds what's needed:

```bash
zpool replace zroot /dev/gpt/disk3
```

Here, my disk3 was defective, so I recreated exactly the same partitions on the new disk and added it. When we do a status check, we can see that the other disks are rebuilding the new one:

```bash
> zpool status
  pool: zroot
 state: DEGRADED
status: One or more devices is currently being resilvered.  The pool will
	continue to function, possibly in a degraded state.
action: Wait for the resilver to complete.
 scrub: resilver in progress for 0h4m, 92.79% done, 0h0m to go
config:

	NAME                        STATE     READ WRITE CKSUM
	zroot                       DEGRADED     0     0     0
	  raidz2                    DEGRADED     0     0     0
	    gpt/disk0               ONLINE       0     0     0  4.63M resilvered
	    gpt/disk1               ONLINE       0     0     0  4.63M resilvered
	    gpt/disk2               ONLINE       0     0     0  4.59M resilvered
	    replacing               DEGRADED     0     0     0
	      10045559159691639333  UNAVAIL      0 1.49K     0  was /dev/gpt/disk3/old
	      gpt/disk3             ONLINE       0     0     0  8.22G resilvered
	    gpt/disk4               ONLINE       0     0     0  4.63M resilvered

errors: No known data errors
```

Once completed, you can either reboot if you chose the solution of booting from the FreeBSD DVD, or do nothing if your machine was already running.

## Resources
- http://wiki.freebsd.org/RootOnZFS/GPTZFSBoot/RAIDZ1
- [Documentation on How to install FreeBSD 7.0 under ZFS](/pdf/how_to_install_freebsd_7.0_under_zfs.pdf)
- http://www.sun.com/bigadmin/features/articles/zfs_part2_ease.jsp
