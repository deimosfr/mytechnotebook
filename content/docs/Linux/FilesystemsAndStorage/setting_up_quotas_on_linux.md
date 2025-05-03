---
weight: 999
url: "/Mise_en_place_des_quotas_sous_Linux/"
title: "Setting up quotas on Linux"
description: "A guide on how to implement and manage disk quotas on Linux systems to control user disk space usage."
categories: ["Linux", "Debian"]
date: "2008-05-23T14:43:00+02:00"
lastmod: "2008-05-23T14:43:00+02:00"
tags: ["storage", "disk management", "administration", "quotas", "system configuration"]
toc: true
---

## Introduction

The assignment of quotas in a file system is a tool that allows control of disk space usage. Quotas consist of setting a space limit for a user or a group of users.

For the creation of these quotas, **2 types of limits** are defined:

* **The soft limit**: indicates the maximum amount of space that a user can occupy on the file system. If this limit is reached, the user receives warning messages about exceeding their assigned quota. When used in combination with grace periods, if the user continues to exceed the soft limit after the grace period has elapsed, they will face the same restriction as when reaching a hard limit.

* **The hard limit**: defines an absolute limit for space usage. The user cannot exceed this limit. Beyond this limit, writing to the file system is forbidden.

Additionally, these limits are expressed in blocks and inodes. The block is a unit of space. Quotas expressed in number of blocks therefore represent a space limit not to be exceeded. As for quotas expressed in number of inodes, they represent the maximum number of files and directories that the user can create.

For reference, grace periods set a time period before the soft limit is transformed into a hard limit. It is set in the following units: second, minute, hour, day, week.

## Prerequisites

We need to check for the presence of quotas at the kernel level (use the configuration file of your current kernel):

```bash
$ grep QUOTA < /boot/config-2.6.18-6-amd64
CONFIG_QUOTA=y
```

## Installation

To install quotas on Debian:

```bash
apt-get install quota
```

## Configuration

### /etc/fstab

Quotas are activated at startup with the quotaon command. Quotas are deactivated when the system shuts down with the quotaoff command.

To set quotas on a file system, you need to update the /etc/fstab file. You'll need to add mount options for the file system(s) concerned. Two options can be used (and combined of course):

* usrquota: activates user quotas
* grpquota: activates group quotas

Example:

```bash
/dev/hdc1 /home ext3 defaults,usrquota 1 1
/dev/hdc2 /tmp ext3 defaults,grpquota 1 1
```

### Creating the necessary structures for quotas to function

**One or two files must be created for quota usage: aquota.user and aquota.group**. These files will contain the quota configuration assigned to users and/or groups. These files must be created at the root of the file systems that include these quotas. Examples:

```bash
touch /home/aquota.user
touch /tmp/aquota.group
```

**Attention: don't forget to modify the permissions on these files! They must have read and write permissions for root only.**
Examples:

```bash
chmod 600 /home/aquota.user
chmod 600 /tmp/aquota.group
```

Remount the file system(s) concerned to take into account the use of quotas for this file system:

```bash
mount -o remount /home
mount -o remount /tmp
```

After creating these files, you need to initialize the quota database by executing the following command:

```bash
$ quotacheck -auvg
edquota: Quota file not found or has wrong format.
No filesystems with quota detected.
```

This is what happens if it doesn't work. Enable quotas:

```bash
quotaon -a
```

## Assignment and verification of quotas

### Setting quotas

Quota assignment is done using the edquota command, which can be used for any type of quota (user or group). The command opens an editor (vi or emacs depending on the content of your EDITOR variable), which allows you to directly modify the aquota.user or aquota.group files.

```bash
Syntax: edquota [-u user] [-g group] [-t]

    * -u user defines quotas for one or more users
    * -g group defines quotas for one or more groups
    * -t defines deadlines
```

Example:

```bash
$ edquota -u citrouille
Disk quotas for user anne (uid 500):
  Filesystem         blocks       soft       hard     inodes     soft     hard
  /dev/hdc1           0       9000       10000         0     90000      10000
```

The file consists of 6 columns:

* Filesystem: file system affected by quotas
* blocks: number of blocks occupied by the user in the file system. Here no file has been created yet.
* soft: soft limit in number of blocks. Here it is set at 9,000 blocks or about 9 MB
* hard: hard limit in number of blocks (about 10 MB)
* inodes: number of inodes occupied by the user in the file system
* soft: soft limit in number of inodes
* hard: hard limit in number of inodes

You will proceed in the same way for assigning quotas to a group. (Don't try to edit these files directly; they are not in text format.)

### Setting a grace period

We've also seen that we can adjust the grace period between when a user reaches the soft limit and when they are banned from any further occupation in the file system. So we're going to set the duration of this grace period. It will be the same for any user and/or group.

Example:

```bash
$ edquota -t
Grace period before enforcing soft limits for users:
Time units may be: days, hours, minutes, or seconds
  Filesystem             Block grace period     Inode grace period
  /dev/hdc1                    7days                  7days
```

So you just need to replace the values with your values in the unit that suits you: second, minute, hour, day, week.

### Exceeding quotas: what happens

For once, we'll put ourselves in the user's position. We'll describe the main cases of exceeding quotas and the messages sent to the user.

Let's take the following example: User Anne has 9MB as a soft limit and 10MB as a hard limit. Her grace period is 7 minutes. Below is the content of the file system subject to these quotas:

```bash
$ ls -l /home/anne
total 1842
-rw-------    1 root     root         7168 fév 28 23:50 aquota.user
-rw-r--r--    1 anne     anne      1857516 mar  1 12:19 fic1
drwx------    2 root     root        12288 nov 28 12:59 lost+found
```

We are well below the quotas. We will now copy file fic1 4 times. The first 3 copies go well and we have fic2, fic3 and fic4. Below is the last copy with user anne:

```bash
$ cp fic1 fic5
ide1(22,10): warning, user block quota exceeded.

$ ls -l
total 9134
-rw-------    1 root     root         7168 fév 28 23:50 aquota.user
-rw-r--r--    1 anne     anne      1857516 mar  1 12:19 fic1
-rw-r--r--    1 anne     anne      1857516 mar  1 13:18 fic2
-rw-r--r--    1 anne     anne      1857516 mar  1 13:18 fic3
-rw-r--r--    1 anne     anne      1857516 mar  1 13:18 fic4
-rw-r--r--    1 anne     anne      1857516 mar  1 13:18 fic5
drwx------    2 root     root        12288 nov 28 12:59 lost+found
```

The soft limit is exceeded. The user receives a message but the write is performed because we haven't exceeded the hard limit.

2 scenarios can then arise if the user does not contact the administrator or if they do not free up space to go back below the soft limit:

* 1st case: the user tries to write to the file system which leads them to exceed the hard limit.

```bash
$ cp fic1 fic6
ide1(22,10): write failed, user block limit reached.
cp: écriture de `fic6': Débordement du quota d'espace disqueL'opération échoue. Une partie du fichier seulement a été copiée. l'utilisateur de pourra plus écrire dans le système de fichiers.
```

* 2nd case: the user lets the 7-minute grace period set by the administrator elapse. They then try to copy the contents of the /etc/passwd file for example. The total space occupied still remains less than the hard limit.

The penalty will be identical to the 1st case. The operation fails:

```bash
$ cp /etc/passwd .
ide1(22,10): write failed, user block quota exceeded too long.
cp: écriture de `./passwd': Débordement du quota d'espace disque L'opération a échoué comme en témoigne le listage ci-dessous:

anne@pingu$ ls -l passwd
        -rw-r--r--    1 anne     anne            0 mar  1 14:48 passwd
```

Similarly if you try to write to the passwd file, you will get the following message in your editor when saving:

```bash
"passwd" erreur d'écriture (système de fichiers plein?) Appuyez sur ENTRÉE ou tapez une commande pour continuer
Il vous est impossible d'écrire.
```

## Verification and display of quotas

The following commands will allow you to verify the quotas assigned to each group and/or user and possibly synchronize the information needed by the system to track these quotas.

### Editing quota-related information

The repquota command displays a summary of quota usage and grace periods.

```bash
Syntax: repquota [ -vug ] -a | filesystem

    * -v: verbose mode, displays additional info
    * -u: displays information about user quotas
    * -g: displays information about group quotas
    * -a: displays information about all file systems with quotas
    * filesystem: displays information about quotas for the specified file system
```

For the example, I added a user Bob.

```bash
$ repquota -avug
 *** Report for user quotas on device /dev/hdc10
 Block grace time: 00:07; Inode grace time: 00:07
                         Block limits                File limits
 User            used    soft    hard  grace    used  soft  hard  grace
 ----------------------------------------------------------------------
 root      --      19       0       0              2     0     0
 anne      --    7293    9000   10000              5  9000 10000
 bob       +-    9000    8000    9000  00:07       5  8000  9000
 +         --      19       0       0              2     0     0

 Statistics:
 Total blocks: 7
 Data blocks: 1
 Entries: 3
 Used average: 3,000000
```

Here we find information related to the quota imposed on users. There will be as many lines as there are users, groups and file systems concerned.

The quotas set in number of blocks and inodes are recalled. We also find the number of blocks and the number of inodes used. When a timestamp appears in the grace column, as for example for Bob, this means that the user (or group) has exceeded the soft limit. The grace period is therefore counting down.

You can also use the quota command followed by the name of a user or a group. Again you will get all the information relating to quotas and the use of the allocated space.

Example: to obtain information related to quotas concerning Anne:

```bash
$ quota anne
Disk quotas for user anne (uid 500):
     Filesystem  blocks   quota   limit   grace   files   quota   limit   grace
     /dev/hdc10    7293    9000   10000               5    9000   10000
```

### Verification and synchronization of quota files

Quota files can sometimes become inconsistent. Management of these then becomes impossible. On the other hand, when you add a new user or a new group using the edquota command, you also need to synchronize the files to take into account this new information.

```bash
Syntax: quotacheck [ -vug ] -a | filesystem

    * -v: verbose mode, displays additional info
    * -u: check only user quota files
    * -g: check only group quota files
    * -a: check quota files for all file systems that have them
    * filesystem: check quota files for the specified file system
```

Example: check all quota files, regardless of the file system concerned:

```bash
$ quotaoff -a
$ quotacheck -auvg
quotacheck: Scanning /dev/hdc10 [/home/anne/quota] done
quotacheck: Checked 2 directories and 10 files
```

That's it for this tutorial on quotas. For more information, consult the man pages of the commands: repquota, quotaon, quotaoff, quotacheck, edquota.

## Resources
- http://www.lea-linux.org/cached/index/Admin-admin_fs-quotas.html
