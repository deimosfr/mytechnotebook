---
weight: 999
url: "/Vérifier_l'intégrité_des_fichiers_sur_sa_Debian/"
title: "Checking File Integrity on Debian"
description: "How to check the integrity of files on a Debian system using a script that verifies package files against their original versions."
categories: ["Debian", "Linux"]
date: "2012-08-29T12:02:00+02:00"
lastmod: "2012-08-29T12:02:00+02:00"
tags: ["Debian", "Security", "File Integrity", "System Administration"]
toc: true
---

![Debian](/images/debian_logo.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Operating System** | Debian 6 |
| **Website** | [Debian](https://www.debian.org) |
| **Last Update** | 29/08/2012 |
{{< /table >}}

## Introduction

For some context, I had a former colleague who found himself working in a company with many compromised servers that had been owned for years with no immediate possibility of replacing them. Knowing that binaries had been modified, he had to verify the integrity of all systems. For this purpose, he created a small script to check everything.

## Usage

(`debian_integrity_check.sh`)

```bash
#!/bin/bash

# Get packages list
dpkg -l | awk '($1=="ii") {print $2"="$3}' > /tmp/pkgs

cat > /tmp/apt.conf <<__EOF
// Only needed if arch_of(broken_system) != uname -m
// APT::Architecture "amd64";

APT::Get::Download-Only "true";
APT::Get::Reinstall "true";

Dir "/"
{
  State::status "/var/lib/dpkg/status";
  Cache "/tmp/new-ar";
};

// the filesystem is read-only, hence we need no root permission to
// run apt-get to get file locks
Debug::NoLocking "true";
__EOF

# Install packages in a temp folder
mkdir -p /tmp/new-ar/archives/partial
APT_CONFIG=/tmp/apt.conf apt-get --reinstall install $(cat /tmp/pkgs)

# Diff temp content with root
debsums --all --changed --generate=all --root=/ --deb-path=/tmp/new-ar/archives $(awk -F= '{print $1}' /tmp/pkgs)
```

All you have to do is make the script executable and run it :-)
