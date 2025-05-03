---
weight: 999
url: "/amazone-s3-sauvegarde-propres-et-automatisees-avec-amazon-s3/"
title: "Amazon S3: Clean and Automated Backups with Amazon S3"
description: "How to set up clean and automated backups using Amazon S3 cloud storage service with various mounting options like FUSE and AutoFS"
categories: ["Backup", "Cloud", "Server"]
date: "2011-02-15T08:19:00+01:00"
lastmod: "2011-02-15T08:19:00+01:00"
tags: ["Amazon S3", "Backup", "Cloud Storage", "FUSE", "AutoFS"]
toc: true
---

## Introduction

For those who don't know, [Amazon S3](https://aws.amazon.com/s3) is a remote backup service. The advantages are:

- the price
- the bandwidth
- the storage

I'll let you look at the price list to see for yourself.

There isn't really any provided software, but rather an API. And with its current popularity, some guys have worked on it and created some pretty nice tools :-)

I'll describe here how I've set up something quite elegant and cost-effective.

## Installation

First, let's install the basic commands (just in case):

```bash
apt-get install s3cmd
```

Then we'll install what's needed for Fuse, which will allow us to navigate on Amazon S3 as if it were a mounted filesystem. Here we install the development libraries for compilation:

```bash
apt-get install libcurl4-openssl-dev libxml2-dev libfuse-dev make g++ fuse-utils
```

_Note: If you already have the s3fs binary (so no compilation needed), just install these packages:_

```bash
apt-get install libcurl3 libfuse2 fuse-utils
```

### Compilation

Now we need to download the sources for Fuse Over Amazon, then compile them:

```bash
wget http://s3fs.googlecode.com/files/s3fs-1.40.tar.gz
tar -xzf s3fs-1.40.tar.gz
cd s3fs-1.40
./configure --bindir /usr/bin
make
make install
```

## Configuration

We'll first use s3cmd to create our _bucket_ (personal folder for amazon s3), then use it with Fuse.

### s3cmd

Let's first configure our account by running the following command:

```bash
s3cmd --configure
```

Fill in all the fields using the information from your account. The most important are the _access key_ and the _secret key_. Once done, we'll create our bucket:

```bash
s3cmd mb monbucket_with_a_unique_name
```

It's important to create a bucket with a specific name because everyone is basically on the same filesystem. So if the name already exists, you won't be able to create this bucket.

### FuseOverAmazon

Let's create the /etc/passwd-s3fs file and fill in the info as follows (`/etc/passwd-s3fs`):

```
access_key:secret_key
```

Enter your access key, followed by ':', then your secret key.

Then we'll apply the appropriate permissions:

```bash
chmod 600 /etc/passwd-s3fs
```

Finally, we can mount all this in /mnt for example (adapt according to your needs):

```bash
s3fs monbucket_with_a_unique_name /mnt -ouse_cache=/tmp
```

And there you go :-). You can do ls, mkdir, cp, rm etc... in /mnt and it will hit your Amazon S3 backup :-).

### Fstab

And if we want to put all this in fstab? This is of course optional, but very practical:

```
s3fs#monbucket_with_a_unique_name /mnt   fuse    ouse_cache=/tmp,noatime,allow_other 0 0
```

Here's the line to add.

### AutoFS

The best setup is with autofs! I invite you to [read this documentation first]({{< ref "docs/Servers/FileSharing/autofs-mounting-and-unmounting-shares.md" >}}) to not be too confused, then add this line in a file auto.as3 for example:

```
amazons3        -fstype=fuse,ouse_cache=/tmp,noatime,allow_other:s3fs\#monbucket_with_a_unique_name
```

Then modify the auto.master file:

```
/mnt /etc/auto.as3 --timeout=60
```

Restart the autofs service for this to work.

## Backups with specific software

With some software, you need to be clever to bypass the 5GB file size limitation imposed by Amazon S3 or multiple files in too large a number which will cost us dearly.

### Sbackup

For Sbackup, it's quite simple, you need to check if the file size is larger than 5 GB and split the files.tgz file.

### BackupPc

For BackupPc, it is recommended to compress existing data before uploading it because there are often too many small files.

## Resources

- http://code.google.com/p/s3fs/wiki/FuseOverAmazon
- http://s3tools.org/s3cmd
