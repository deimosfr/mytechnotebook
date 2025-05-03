---
weight: 999
url: "/Décompressions_sous_différents_formats/"
title: "Uncompress for Different Formats"
description: "A guide to uncompress different types of archive formats in Linux including RPM, DEB, ZIP, CAB, and more."
categories: ["Linux"]
date: "2006-10-10T16:27:00+02:00"
lastmod: "2006-10-10T16:27:00+02:00"
tags: ["Servers", "Linux", "Commands", "File Formats"]
toc: true
---

Ah, file uncompression! Always having to remember that "-xjvfirhfidopgnfudjs" command! Not so easy, right? So here's a little reminder:

Before anything else, if you don't know what a file contains: **file \<filename\>**  
This command displays the file type based on its content, not its extension! (This information is based on the `/etc/magic` file)

## RPM Files

- Extract an rpm:

```bash
rpm2cpio <file.rpm> | cpio -mid
```

rpm2cpio belongs to the "rpm" package  
cpio belongs to the "cpio" package

## DEB Files

- Extract a deb:

```bash
ar xv <file.deb>
```

ar belongs to the binutils package

## ZIP Files

- Extract a zip:

```bash
unzip <file.zip>
```

unzip belongs to the infozip package

## Microsoft CAB Files

- Extract a Microsoft cab:

```bash
cabextract <file.cab>
```

cabextract can be obtained from uklinux.net

## InstallShield CAB Files

- Extract an InstallShield cab:

```bash
unshield <file.cab>
```

unshield can be obtained from synce.sourceforge.net  
Note: InstallShield cab files are usually named data1.cab, data1.hdr, data2.cab, etc.

## ARJ Files

- Extract an arj:

```bash
unarj x <file.arj>
```

unarj belongs to the "bin" package, and a complete version of arj can be obtained from arj.sourceforge.net (in which case you would use "arj x" instead of "unarj x")

## RAR Files

- Extract a rar:

```bash
unrar x <file.rar>
```

unrar can be obtained from rarlab.com

## ACE Files

- Extract an ace:

```bash
unace x <file.ace>
```

unace ("LinUnAce") can be obtained from winace.com

## LHA Files

- Extract a lha:

```bash
lha x <file.lha or file.lzh>
```

lha is available from its official site

## JAR Files

- Extract a jar:

```bash
jar xvf <file.jar>
```

jar can be obtained from Sun's JRE or JDK  
Note: xpi files are actually jar files.

## 7z Files

- Extract a 7z:

```bash
7za x <file.7z>
```

7za can be obtained from the p7zip project page on Sourceforge.  
For those who don't know what 7z format is, take a look at the 7zip homepage, which is a free zip/7z archiver for Windows.

## Common Formats

- Those that need no introduction:

Uncompress a Z file:

```bash
uncompress <file.Z>
```

Uncompress a gz file:

```bash
gzip -d <file.gz>
```

Uncompress a bz2 file:

```bash
bzip2 -d <file.bz2>
```

Extract a tar:

```bash
tar xvf <file.tar>
```

And the combinations...

Extract a tgz or tar.gz:

```bash
tar zxvf <file.tgz>
```

Extract a tar.bz2:

```bash
tar jxvf <file.tar.bz2>
```
