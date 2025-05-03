---
weight: 999
url: "/ajouter-le-dvd-red-hat-comme-repository/"
title: "Adding Red Hat DVD as a Repository"
description: "How to add a Red Hat DVD as a local or remote YUM repository for installing packages without internet access or subscription."
categories: ["Linux", "Red Hat", "System Administration"]
date: "2012-02-21T16:21:00+02:00"
lastmod: "2012-02-21T16:21:00+02:00"
tags: ["Red Hat", "YUM", "Repository", "DVD", "Package Management"]
toc: true
---

## Introduction

You don't have a Red Hat subscription or simply no internet connection to connect to the NRH, yet you have the installation DVD that contains all the packages you need!

Here's how to add the DVD to the repository.

## Configuration

### DVD

Create a file in /etc/yum.repos.d/redhatcd.repo

```bash
cp /media/RHEL/media.repo /etc/yum.repos.d/redhatcd.repo
```

And adapt it with this content:

```bash
[InstallMedia]
name=Red Hat Enterprise Linux 6.1
mediaid=1305068199.328169
metadata_expire=-1
gpgcheck=0
cost=500
baseurl=file:///media/RHEL/Server
```

Now we'll clear the cache:

```bash
yum clean all
```

That's it, now you can use yum with your DVD.

### Local Repository

If you want to use a local repository, copy the "Server" folder from the DVD to a local folder (for example /home/repo/Server).

Now you need to install the createrepo-0.9.8-4.el6.noarch.rpm package (or another version), as well as all necessary dependencies:

```bash
yum install createrepo
```

Then run this command if you're on Red Hat 6:

```bash
createrepo -v /home/repo/Server
```

Or this command if you're on Red Hat 5:

```bash
createrepo -d -s sha1 /home/repo/Server
```

Then create a file /etc/yum.repos.d/local.repo

```bash
[InstallMedia]
name=Red Hat Enterprise Linux 6.1
mediaid=1305068199.328169
metadata_expire=-1
gpgcheck=0
cost=500
baseurl=file:///home/repo/Server/
enabled=1
```

You can check the list of repositories with the command:

```bash
yum repolist
```

Then before installing a package, run:

```bash
yum clean all
```

### Remote Repository

You can install an httpd server, then put the "Server" folder from the DVD to use as a remote yum DVD repository. Here's an example of client-side configuration:

```bash
[dvd-base]
name=Red Hat Enterprise Linux $releasever Beta - Base - $basearch - DVD
baseurl=http://server/repositories/rhel/$releasever/$basearch/
enabled=1
gpgcheck=1
gpgkey=file:///etc/pki/rpm-gpg/RPM-GPG-KEY-redhat-beta,file:///etc/pki/rpm-gpg/RPM-GPG-KEY-redhat-release
```

For the server part, it's identical to [the local configuration](#local-repository).

## FAQ

### [Errno -3] Error performing checksum

If you get this kind of message:

```
http://server/repositories/rhws/5Client/i386/repodata/primary.xml.gz: [Errno -3] Error performing checksum
Trying other mirror.
Error: failure: repodata/primary.xml.gz from dvd-base: [Errno 256] No more mirrors to try.
```

It's probably because you're running the createrepo command without the necessary arguments for Red Hat 5. Brief explanation:

```
Because RPM packages for Red Hat Enterprise Linux 6 are compressed using the XZ lossless data compression format, and may also be signed using alternative (and stronger) hash algorithms such as
SHA-256, it is not possible to run createrepo on Red Hat Enterprise Linux 5 to create the package metadata for Red Hat Enterprise Linux 6 packages.
The createrepo command relies on rpm to open and inspect the packages, and rpm on Red Hat Enterprise Linux 5 is not able to open the improved Red Hat Enterprise Linux 6 RPM package format.
```

## Resources
- https://access.redhat.com/kb/docs/DOC-9744
- http://samixblog.blogspot.com/2011/11/yum-errno-3-error-performing-checksum.html
- http://nareshov.wordpress.com/2011/12/22/rpmbuild-behaviour-centos5-vs-centos6/
