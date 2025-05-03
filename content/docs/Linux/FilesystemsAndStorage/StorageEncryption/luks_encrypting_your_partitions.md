---
weight: 999
url: "/Luks_\\:_Chiffrer_ses_partitions/"
title: "LUKS: Encrypting Your Partitions"
description: "Learn how to use LUKS to encrypt partitions on Linux, including creating encrypted partitions, unlocking them, and managing passphrases."
categories: ["Linux"]
date: "2013-12-23T21:16:00+02:00"
lastmod: "2013-12-23T21:16:00+02:00"
tags: ["encryption", "security", "partitions", "luks"]
toc: true
---

## Introduction

[LUKS](https://en.wikipedia.org/wiki/Linux_Unified_Key_Setup) is one of the best disk encryption tools for Linux. We'll see here how to use it.

## Usage

### Creating an Encrypted Partition

Be aware that if you use an existing partition, all its data will be erased when initializing the encrypted partition. To initialize it (sdb1 for example):

```bash
$ cryptsetup luksFormat /dev/sdb1 

WARNING!
========
This will overwrite data on /dev/sdb1 irrevocably.

Are you sure? (Type uppercase yes): YES
Enter LUKS passphrase: 
Verify passphrase:
```

Enter the password you want to use to decrypt the partition.

### Unlocking

Next, we'll unlock the encrypted partition to use it:

```bash
cryptsetup luksOpen /dev/sdb1 secret
Enter passphrase for /dev/sdb1:
```

'secret' corresponds here to the device mapper name. We can then verify its existence:

```bash
$ ls /dev/mapper/
control  secret
```

### Preparing the Partition

Now we just need to format this partition:

```bash
mkfs.ext4 /dev/mapper/secret
```

And mount it in a directory.

### Unmounting the Encrypted Disk

Once you've finished, you need to properly close the disk by unmounting and locking it:

```bash
umount /dev/mapper/secret
cryptsetup luksClose secret
```

### Mounting the Encrypted Partition Permanently

If you want to mount the partition permanently, you'll need to use fstab and crypttab. In crypttab:

```bash
secret /dev/sdb1 /root/password luks
```

* secret: name of the device mapper
* /dev/sdb1: the physical device
* /root/password: the file containing your password (you can alternatively put the password directly in the crypttab file)

If you chose to use a file containing the key, create it like this:

```bash
cryptsetup luksAddKey /dev/sdb1 /root/password
chmod 600 /root/password
```

Then add the following line to fstab:

```bash
...
/dev/mapper/secret   /mnt   ext4   defaults   1 2
...
```

Your encrypted partition will now mount automatically at startup (there's less benefit, but it might interest some users).

### Adding a Passphrase

To add a passphrase (maximum of 8 total), here's how to proceed. First locate your encrypted partition:

```bash
cryptsetup luksDump /dev/sdb1
```

Once you're sure it's the right one, add an additional passphrase:

```bash
cryptsetup luksAddKey /dev/sdb1
```

If you want to change a passphrase, you'll need to delete the old one using the method mentioned below.

### Removing a Passphrase

If you want to remove one of your passphrases:

```bash
cryptsetup luksRemoveKey /dev/sdb1
```
