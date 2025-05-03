---
weight: 999
url: "/Modifier_le_firmware_par_le_OpenWrt/"
title: "Modifying the firmware with OpenWrt"
description: "A guide on how to modify the firmware of a Fonera+ router using OpenWrt, including extracting the original firmware, modifying it, and uploading the modified version."
categories: ["Linux", "Network", "Development"]
date: "2009-07-03T01:45:00+02:00"
lastmod: "2009-07-03T01:45:00+02:00"
tags:
  ["OpenWrt", "Firmware", "Router", "Fonera", "Development", "Kernel", "Linux"]
toc: true
---

## Introduction

As you know, some time ago, I tried to install **OpenWrt on a Fonera+** (the larger one with two Ethernet connectors). After persisting, I finally managed to get something partially functional: recompiling the Fon firmware and getting both Ethernet ports working. However, WiFi still remains a problem.

After a period of working on my Frankenstein-modified Foneras, I decided to try again with another approach: using the Fon binary firmware and activating a serial console. The first step was to get a binary update for the Fonera+ in the form of the file foneraplus_1.1.1.1.fon.

This is a gzip compressed tar archive to which Fon added a specific header as explained by [Stefan Tomanek on his site](https://stefans.datenbruch.de/lafonera/upgrades.shtml). His explanations for the classic Fonera firmware seem to be perfectly valid for the Fonera+.

## Modifying the firmware

We use cat and tail to remove the header and get a real archive:

```bash
cat foneraplus_1.1.1.1.fon  | tail -c +520 - > arch.tar.gz
```

The resulting file contains three files that are used by the router to update itself:

```bash
hotfix
image_full
upgrade
```

The image_full file is built with the kernel image, the root file system, and a checksum for the stage2 loader launched by RedBoot. A look at the content of the Fon BuildRoot Makefile (the firmware sources) and we understand the build procedure:

- the lzma compressed kernel image is built classically
- the Perl script fonimage.pl adds a 12-byte header containing a CRC via the Digest::CRC module
- dd is used to get a file multiple of 64KB with zero padding if necessary (bs=64k conv=sync)
- the kernel image file doesn't seem to be used
- fonimage.pl is also used to group the kernel and the root file system in squashfs. The Perl script is launched to build the complete system image with the CRC.

We can conclude that the image_full file is a header/CRC of 12 bytes followed by the lzma kernel then the squashfs image of the root file system. A look at the Perl script also tells us the composition of the 12-byte header:

- 4 bytes for the size
- 4 bytes for the CRC
- 4 bytes for the location of the root file system

Checking is easy, just display the values using fonimage.pl manually:

```bash
./fonimage.pl daimage vmlinux.lzma root.squashfs
size  : 2295397
offset: 715783
```

```bash
ls -l vmlinux.lzma daimage
2295409 2007-11-18 12:20 daimage
715783 2007-09-29 20:31 vmlinux.lzma
```

The offset corresponds to the kernel size and the size corresponds to the final image size minus the header size. The numerical values are stored in the header in a specific format called unsigned long in "network" (big-endian) order. Indeed, in the header of the daimage file, we find:

```bash
hexdump daimage | head -n 1
0000000 2300 6506 3616 ac9f 0a00 07ec 006d 8000
```

23006506 is actually 00230665, which is 2295397, the size of daimage minus 12. The image_full file has this header:

```bash
hexdump image_full | head -n 1
0000000 2100 debf 14a2 9bd3 0a00 3450 006d 8000
```

2100debf is 0021bfde, which is 2211806. 2211806+12 equals 2211818. But the file size is 2293764. It doesn't match! Back to the Makefile.

After using the fonimage.pl script, the Makefile calls prepare_generic_squashfs declared in image.mk, an OpenWrt file. The function uses dd:

```bash
dd if=$(1) of=$(KDIR)/tmpfile.1 bs=64k conv=sync
$(call add_jffs2_mark,$(KDIR)/tmpfile.1)
dd of=$(1) if=$(KDIR)/tmpfile.1 bs=64k conv=sync
$(call add_jffs2_mark,$(1))
```

We find the padding in 64k blocks and a call to add_jffs2_mark which does:

```bash
echo -ne '\xde\xad\xc0\xde' >> $(1)
```

In short:

- we resize the image in 64k blocks
- we add 0xdeadc0de (Dead Code???)
- we resize again
- we add 0xdeadc0de again

Let's verify with a Fon BuildRoot compiled by ourselves:

```bash
cd bin
find .. -name *.lzma
../build_mips/linux-2.6-fonera/vmlinux.lzma
find .. -name root.squashfs
../build_mips/linux-2.6-fonera/root.squashfs
```

```bash
../target/linux/fonera-2.6/image/fonimage.pl \
daimage \
../build_mips/linux-2.6-fonera/vmlinux.lzma \
../build_mips/linux-2.6-fonera/root.squashfs
```

```bash
dd if=daimage of=temp1 bs=64k conv=sync
35+1 enregistrements lus
36+0 enregistrements écrits
2359296 octets (2.4 MB) copiés, 0.00566277 seconde, 417 MB/s
```

```bash
echo -ne '\xde\xad\xc0\xde' >> temp1
```

```bash
dd if=temp1 of=daimage1 bs=64k conv=sync
36+1 enregistrements lus
37+0 enregistrements écrits
2424832 octets (2.4 MB) copiés, 0.00597735 seconde, 406 MB/s
```

```bash
echo -ne '\xde\xad\xc0\xde' >> daimage1
```

```bash
ls -l daimage1 openwrt-fonera-2.6.image
2424836 2007-11-18 13:37 daimage1
2424836 2007-09-29 20:31 openwrt-fonera-2.6.image
```

```bash
md5sum daimage1 openwrt-fonera-2.6.image
b57f8b2279bdb9a6483e094c58fc3381  daimage1
b57f8b2279bdb9a6483e094c58fc3381  openwrt-fonera-2.6.image
```

Bingo! We are able to manually recreate an image. We have perfectly analyzed the build process. Now we need to reverse this process to get a kernel file and a rootfs image that we can modify.

![Fon mem map](/images/fon_mem_map.avif)

The simplest way to be sure is to dismantle our own attempt. So we start by looking at the header of daimage1:

```bash
hexdump daimage1 |  head -n 1
0000000 2300 6506 3616 ac9f 0a00 07ec 006d 8000
```

The values we're interested in here are:

- the image size: 23006506 or 00230665 or 2295397 bytes
- the kernel size (the offset to find the rootfs): 0a0007ec or 000aec07 or 715783

We remove the header:

```bash
cat daimage1 | tail -c +13 > nohead
```

Then the padding using the image size specified in the header:

```bash
cat nohead | head -c +2295397 > nopad
```

Our file is now the concatenation of the kernel and rootfs. We use the offset specified in the header to retrieve the rootfs. image size - offset = rootfs size (2295397-715783=1579614):

```bash
cat nopad | tail -c 1579614 > squash
```

While we're at it, we also extract the kernel to have all the elements:

```bash
cat nopad | head -c +715783 > dakern
```

Let's check:

```bash
md5sum dakern ../build_mips/linux-2.6-fonera/vmlinux.lzma
0c50b77f12c3e18a91db1d027fe0ecc6  dakern
0c50b77f12c3e18a91db1d027fe0ecc6  ../build_mips/linux-2.6-fonera/vmlinux.lzma
```

```bash
md5sum squash ../build_mips/linux-2.6-fonera/root.squashfs
19576d7b96f07ba7694d615e2afe78d1  squash
19576d7b96f07ba7694d615e2afe78d1  ../build_mips/linux-2.6-fonera/root.squashfs
```

It works! We are able to dismantle an official Fon update. Let's apply this to image_full:

```bash
hexdump image_full | head -n 1
0000000 2100 debf 14a2 9bd3 0a00 3450 006d 8000
```

Some calculations:

- image size: 2100debf or 0021bfde or 2211806
- kernel size: 0a003450 or 000a5034 or 675892
- rootfs size: 2211806-675892=1535914

Dismantling:

```bash
cat image_full | tail -c +13 > nohead
cat nohead | head -c +2211806 > nopad
cat nopad | tail -c 1535914 > squash
cat nopad | head -c +675892 > dakern
```

Quick verification:

```bash
hexdump dakern | head -n 1
0000000 006d 8000 ff00 ffff ffff ffff 00ff 0204
```

That's indeed a lzma kernel header.

Now we need to dismantle the squashfs. We're facing several problems including lzma compression and endianness. The image is built for a big endian system (MIPS) but a PC x86 is little endian. Simply forget the idea of mounting the file system in loopback. We need to look at the unsquashfs tool, but again, many problems. After a long period of research, I finally came across a firmware modification kit: [firmware_mod_tools_prebuilt.tar.gz](https://download.berlios.de/firmwaremodkit/firmware_mod_tools_prebuilt.tar.gz).

This archive contains an unsquashfs-lzma that works perfectly with the images for the Fonera (both the new and old ones):

```bash
/tmp/unsquashfs-lzma squash
Reading a different endian SQUASHFS filesystem on squash
created 407 files
created 67 directories
created 165 symlinks
created 0 devices
created 0 fifos
```

And we find a squashfs-root directory containing the useful tree structure. A quick look in /etc/inittab and indeed, there, something is missing that I would have liked:

```bash
ttyS0::askfirst:/bin/ash --login
```

We add it and build a brand new squashfs using the BuildRoot tools:

```bash
/path/to/openwrt/staging_dir_mips/bin/mksquashfs-lzma \
squashfs-root root.squashfs -nopad -noappend -root-owned -be
Creating big endian 3.0 filesystem on root.squashfs, block size 65536.
Big endian filesystem, data block size 65536, compressed data,
compressed metadata, compressed fragments
Filesystem size 1499.77 Kbytes (1.46 Mbytes)
31.47% of uncompressed filesystem size (4766.19 Kbytes)
Inode table size 4801 bytes (4.69 Kbytes)
23.44% of uncompressed inode table size (20479 bytes)
Directory table size 5572 bytes (5.44 Kbytes)
57.00% of uncompressed directory table size (9776 bytes)
Number of duplicate files found 4
Number of inodes 639
Number of files 407
Number of fragments 28
Number of symbolic links  165
Number of device nodes 0
Number of fifo nodes 0
Number of socket nodes 0
Number of directories 67
Number of uids 1
root (0)
Number of gids 0
```

## Uploading the Firmware

We have a kernel and a new rootfs. We reuse the commands to create an image: fonimage.pl, dd, echo, dd, echo. We get a new image that we just need to flash on the Fonera+ via RedBoot:

```bash
fis delete image
load -r -b 0x80041000 myimage
fis create -b 0x80041000 -f 0xA8040000 -l 0x00230004 -e 0x80040400 image
```

After this operation, the Fonera+ is restarted by the reset command. No kernel messages appear, which is normal, but after the complete boot phase we do have a usable shell. The DHCP messages keep cluttering the console but it works!

## Resources
- http://www.lefinnois.net/wp/index.php/2007/11/18/modification-du-firmware-fon-pour-une-fonera/
- http://www.dd-wrt.com/wiki/index.php/LaFonera_Software_Flashing
- http://www.dd-wrt.com/wiki/index.php/LaFonera_(fr)
- http://www.cure.nom.fr/blog/archives/141-Fonera-et-le-firmware-alternatif-DD-WRT.html
- http://www.dd-wrt.com/wiki/index.php/Wireless_Access_Point
