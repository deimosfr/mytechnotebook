---
weight: 999
url: "/RPM_\\:_Build_a_binary_RPM_package_from_sources_on_Red_Hat/"
title: "RPM: Build a Binary RPM Package from Sources on Red Hat"
description: "Guide on how to build a binary RPM package from source files on Red Hat systems, including compilation and packaging steps."
categories: ["Linux", "Red Hat"]
date: "2011-06-08T14:23:00+02:00"
lastmod: "2011-06-08T14:23:00+02:00"
tags: ["RPM", "Red Hat", "packaging", "compilation", "Development"]
toc: true
---

## Introduction

Building an RPM is not a complicated task, even from source files. In my job I had to compile [Open On Load](https://www.openonload.org/) kernel module and package it.

The problem is when you want to deploy this kernel module on several production servers and you don't want to have all development tools installed on all production servers. This is the documentation I've made to create RPMs with binaries from sources.

## Compilation

It is better to test compilation before trying to package it (just to be sure it works).  
So you need to have a Red Hat with development tools installed to compile Open On Load module:

```bash
yum groupinstall "Development Tools"
```

Now we can compile:

```bash
wget http://www.openonload.org/download/openonload-201104.tgz
tar -xzf openonload-201104.tgz
rm -f openonload-201104.tgz
cd openonload-201104
./scripts/onload_install
```

## Create package

Now let's build a binary RPMS from the source RPM:

```bash
cd ..
tar -cf openonload-201104.tgz openonload-201104
rpmbuild -ts openonload-201104.tgz
cd /root/rpmbuild/SRPMS/
rpmbuild --rebuild openonload-201104-1.src.rpm
```

You now have got packages in `/root/rpmbuild/RPMS/x86_64/` and can deploy them without any development packages on production servers.
