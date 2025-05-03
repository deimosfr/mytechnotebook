---
weight: 999
url: "/Mise_en_place_des_ACLs_pour_CVS/"
title: "Setting up ACLs for CVS"
description: "A guide on how to install and configure ACLs for CVS repositories."
categories: ["Linux", "Development"]
date: "2006-11-24T09:42:00+02:00"
lastmod: "2006-11-24T09:42:00+02:00"
tags: ["CVS", "Linux", "Development", "Servers", "Security", "Administration"]
toc: true
---

## Installation

Download the latest pre-patched version of CVSACL from the internet (http://cvsacl.sourceforge.net), then extract it.

```bash
wget http://switch.dl.sourceforge.net/sourceforge/cvsacl/cvs-1.11.22-cvsacl-1.2.5-patched.tar.gz
tar -xzvf cvs-1.11.22-cvsacl-1.2.5-patched.tar.gz
```

Configuration, compilation and installation:

```bash
cd cvs-1.11.22-cvsacl-1.2.5-patched.tar.gz
./configure
make
make install
```

If CVS is properly installed with its CVSACL Patch, the command `cvs --version` should output:

```
Concurrent Versions System (CVS) 1.11.22 (client/server)
with CVSACL Patch 1.2.5 (cvsacl.sourceforge.net)
```

## Configuration

### Preparation of the repository

#### If the repository doesn't exist yet

* Create the repository directory (the directory can be created anywhere, preferably in a location with sufficient space)

```bash
mkdir -p /home/cvsadmin/cvsroot
```

* Create a symbolic link to the repository

```bash
ln -s /home/cvsadmin/cvsroot /usr/local/cvsroot
```

* Define the $CVSROOT variable and initialize the repository

```bash
export CVSROOT=/usr/local/cvsroot
cvs -d $CVSROOT init
```

After this operation, a CVSROOT directory will be created at the root of the repository.
This directory, seen as the first module of the repository, contains all the CVS configuration files.

#### If the repository already exists

* Copy the `aclconfig.default` file from the sources to the CVSROOT folder at the root of the repository.

```bash
sudo cp /root/cvs-1.11.22-cvsacl-1.2.5-patched.tar.gz/aclconfig.default $CVSROOT/CVSROOT
```

* Rename the file to `aclconfig`

```bash
sudo mv $CVSROOT/CVSROOT/aclconfig.default $CVSROOT/CVSROOT/aclconfig
```

### Configuration of Lockfiles

Create a `.lock` directory at the root of the repository (`/usr/local/cvsroot`) and allow full access to all users:

```bash
mkdir /usr/local/cvsroot/.lock
chmod 777 /usr/local/cvsroot/.lock
```

Edit the `config` file in `/usr/local/cvsroot/CVSROOT/` and modify it as follows:

```bash
# Put CVS lock files in this directory rather than directly in the repository.
LockDir=/home/cvsadmin/cvsroot/.lock
```

### Configuration of users

CVSACL offers the possibility to manage access rights by group either by using system groups (`/etc/group`) or by using its own group file.
It's preferable to use the latter method as it allows for CVS rights management completely independent from the system.
To do this, edit the `aclconfig` file in `/usr/local/cvsroot/CVSROOT` and modify it as follows:

```bash
# Set `UseSystemGroups' to yes to use system group definitions (/etc/group).
#UseSystemGroups=yes
# Set `UseCVSGroups' to yes to use another group file.
UseCVSGroups=yes
```

Then create the `group` file in `/usr/local/cvsroot/CVSROOT`:

```bash
touch /usr/local/cvsroot/CVSROOT/group
```

The group file must be in the form "group:user1, user2, user3 ..."  
Groups are totally independent of system groups and can bear any name.  
Users are those from the system. You will need to create them with `useradd` **as well as their home directory**.

### Configuration of ACLs

First, define "root" as the CVS owner.  
This will allow them to perform all CVS tasks and administer access to the repository.

```bash
cvs -d /usr/local/cvsroot racl root:p -r ALL ALL
```

#### Definition of rights

Syntax of the rights assignment command:

```bash
cvs -d </Path/to/repository> racl <user or group>:<right> <TAG> <BRANCH>
```

* no access

```
Command line character: n
```

No possible action on the repository

* read

```
Command line character: r
```

Read-only. With these rights, only the following actions are possible: annotate, checkout, diff, export, log, rannotate, rdiff, rlog, status.

* write

```
Command line character: w
```

This permission only allows cvs commit/checkin actions. It does not allow adding/removing a file from/to the repository; other permissions are defined for this.

* tag

```
Command line character: t
```

This permission authorizes the cvs tag and rtag sub-commands, so it is possible to control Tag and Untag operations. The "t" permission includes the "r" permission because reading is mandatory for tagging. However, "t" does not include writing; it is not possible to commit with just this permission.

* create

```
Command line character: c
```

The "c" permission authorizes the creation/deletion of files from/to the repository but once again this permission does not include "w"; we can only import or export files. After adding a file, it is necessary to perform a commit which will be accepted because we are adding a file, not modifying it.

* delete

```
Command line character: d
```

"d" authorizes deletion and does not include "w"

* full access except admin rights

```
Command line character: a
```

"a" includes all permissions listed above except ACL management rights.

* acl admin

```
Command line character: p
```

"p" indicates that the user is an owner. They have full control of the repository and can manage ACLs.

### Example of ACL

```bash
cvs -d /usr/local/cvsroot group1:r -r ALL ALL
cvs -d /usr/local/cvsroot group2:n -r ALL module1
cvs -d /usr/local/cvsroot user1:w -Rr ALL module2
```

The first line authorizes reading on all directories of the repository for the group "group1"  
The second line prohibits access to "module1" for the group "group2"  
The third line allows "user1" to modify files in ALL subdirectories and files of "module2" (-R = recursively)
