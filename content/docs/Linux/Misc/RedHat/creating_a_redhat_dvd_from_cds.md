---
weight: 999
url: "/Créer_un_DVD_RedHat_à_partir_des_CD/"
title: "Creating a RedHat DVD from CDs"
description: "Guide for creating a bootable RedHat DVD from CD images"
categories: ["Linux", "RHEL", "Red Hat"]
date: "2007-08-17T15:30:00+02:00"
lastmod: "2007-08-17T15:30:00+02:00"
tags: ["Linux", "Red Hat", "ISO", "DVD", "Installation"]
toc: true
---

## Introduction

As of this writing, the current RHEL release (5) is only available on CDs, not DVD. It is possible to create a bootable DVD ISO from these CDs using Chris Kloiber's mkdvdiso.sh script.

## Method 1

### Script mkdvdiso.sh

Insert this in a mkdvdiso.sh file:

```bash
#!/bin/bash
 
# by Chris Kloiber <ckloiber@redhat.com>
 
# A quick hack that will create a bootable DVD iso of a Red Hat Linux
# Distribution. Feed it either a directory containing the downloaded
# iso files of a distribution, or point it at a directory containing
# the "RedHat", "isolinux", and "images" directories.
 
# This version only works with "isolinux" based Red Hat Linux versions.
 
# Lots of disk space required to work, 3X the distribution size at least.
 
# GPL version 2 applies. No warranties, yadda, yadda. Have fun.
 
 
if [ $# -lt 2 ]; then
	echo "Usage: `basename $0` source /destination/DVD.iso"
	echo ""
	echo "        The 'source' can be either a directory containing a single"
	echo "        set of isos, or an exploded tree like an ftp site."
	exit 1
fi
 
cleanup() {
[ ${LOOP:=/tmp/loop} = "/" ] && echo "LOOP mount point = /, dying!" && exit
[ -d $LOOP ] && rm -rf $LOOP 
[ ${DVD:=~/mkrhdvd} = "/" ] && echo "DVD data location is /, dying!" && exit
[ -d $DVD ] && rm -rf $DVD 
}
 
cleanup
mkdir -p $LOOP
mkdir -p $DVD
 
if [ !`ls $1/*.iso 2>&1>/dev/null ; echo $?` ]; then
	echo "Found ISO CD images..."
	CDS=`expr 0`
	DISKS="1"
 
	for f in `ls $1/*.iso`; do
		mount -o loop $f $LOOP
		cp -av $LOOP/* $DVD
		if [ -f $LOOP/.discinfo ]; then
			cp -av $LOOP/.discinfo $DVD
			CDS=`expr $CDS + 1`
			if [ $CDS != 1 ] ; then
                        	DISKS=`echo ${DISKS},${CDS}`
                	fi
		fi
		umount $LOOP
	done
	if [ -e $DVD/.discinfo ]; then
		awk '{ if ( NR == 4 ) { print disks } else { print ; } }' disks="$DISKS" $DVD/.discinfo > $DVD/.discinfo.new
		mv $DVD/.discinfo.new $DVD/.discinfo
	fi
else
	echo "Found FTP-like tree..."
	cp -av $1/* $DVD
	[ -e $1/.discinfo ] && cp -av $1/.discinfo $DVD
fi
 
rm -rf $DVD/isolinux/boot.cat
find $DVD -name TRANS.TBL | xargs rm -f
 
cd $DVD
mkisofs -J -R -v -T -o $2 -b isolinux/isolinux.bin -c isolinux/boot.cat -no-emul-boot -boot-load-size 4 -boot-info-table .
/usr/lib/anaconda-runtime/implantisomd5 --force $2
 
cleanup
echo ""
echo "Process Complete!"
echo ""
```

### Install the anaconda-runtime

Prerequisite package and its dependencies:

```bash
yum -y install anaconda-run || timerpm -q anaconda-runtime
```

### Make DVD

Place the CD ISO files and mkdvdiso.sh in a directory, and run the mkdvdiso.sh script.

```bash
chmod 755 ./mkdvdiso.sh && ./mkdvdiso.sh
```

The following command creates a RHEL bootable DVD ISO:

```bash
./mkdvdiso.sh . $(pwd)/RedHat-dvd.iso
```

The file .discinfo should look like:

```
1170972069.396645
Red Hat Enterprise Linux Server 5
i386
1,2,3,4,5
Server/base
Server/RPMS
Server/pixmaps
```

## Method 2

Create folders:

```bash
mkdir /mnt/cd{1,2,3,4,5}
```

Mount every iso image with:

```bash
mount -o loop cd1.iso /mnt/cd1
```

Copy the minimum:

```bash
cp -a cd1/isolinux cd1/.discinfo .
```

Edit the .discinfo and replace 1 with:

```
1,2,3,4,5
```

Do this for all your CDs. Then:

```bash
mkisofs -o redhat.iso -b isolinux/isolinux.bin -c isolinux/boot.cat -no-emul-boot -boot-load-size 4 -boot-info-table -R -m TRANS.TBL -x \
/mnt/cd1/.discinfo -x /mnt/cd1/isolinux -graft-points /mnt/cd1 .discinfo=.discinfo isolinux/=isolinux RedHat/=/mnt/cd2/RedHat RedHat/=/mnt/cd3/RedHat \
RedHat/=/mnt/cd4/RedHat RedHat/=/mnt/cd5/RedHat
```

Adapt this for your specific needs.
