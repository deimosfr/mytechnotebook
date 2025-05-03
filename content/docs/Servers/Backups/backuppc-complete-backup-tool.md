---
weight: 999
url: "/BackupPC_\\:_Un_outil_complet_de_backup/"
title: "BackupPC: A Complete Backup Tool"
description: "How to install, configure and use BackupPC to backup all your data including Windows hosts and SQL databases."
categories: ["Backup", "Server", "Windows", "SQL"]
date: "2012-07-09T07:24:00+02:00"
lastmod: "2012-07-09T07:24:00+02:00"
tags: ["backuppc", "backup", "windows", "rsync", "smb", "mysql", "postgresql"]
toc: true
---

## Introduction

If you're looking for a beautiful tool for making backups and restorations, look no further, BackupPC is for you.

## Configuration

### Windows Host via SMB

You can backup via network shares. For this on Windows, set up your share, and for the machine here's an example:

```
#============================================================= -*-perl-*-
#
# Configuration file for Windows hosts.
# Note the slashes instead of backslashes
#

###########################################################################
# What to backup and when to do it
###########################################################################
#Array of directories to backup
$Conf{BackupFilesOnly} = ['/Documents and Settings', 'Travail'];
#Array of directories excluded from backup
$Conf{BackupFilesExclude} = '/Documents and Setings/user1/Local Settings/Temp';

###########################################################################
# General per-PC configuration settings
###########################################################################
#NetBios name of the machine
$Conf{ClientNameAlias} = 'netbiosname';
#Backup method used
$Conf{XferMethod} = 'smb';
#Verbosity level of log files
$Conf{XferLogLevel} = 1;
#Name of shares to backup
$Conf{SmbShareName} = ['C$'];
#Network user name
$Conf{SmbShareUserName} = 'Administrateur';
#Network user password
$Conf{SmbSharePasswd} = 'secret';
#Backup compression method
$Conf{ArchiveComp} = 'bzip2';
```

### Windows Host via rsync

We'll use the method where we'll do without Cygwin. You can use the method with Cygwin, but don't install both simultaneously.

To begin, download [cwrsync](https://www.itefix.no/i2/node/10650) and install it on the Windows machine. During installation, it will create a user with a randomly generated password, leave it as is. This user is dedicated to starting the rsyncd service.

Then, go to the folder "C:\Program Files (x86)\ICW", then edit the configuration file:

```
use chroot = false
strict modes = false
hosts allow = *
log file = rsyncd.log

# Module definitions
# Remember cygwin naming conventions : c:\work becomes /cygwin/c/work
#
[share]
path = /cygdrive/c/share
comment = share rsync
read only = yes
transfer logging = yes
hosts allow = 192.168.0.14
secrets file = rsyncd.secrets
auth users = backuppc
```

Adapt the following lines:

- path: /cygdrive/ is mandatory. Then use the letter of the drive you're interested in (here 'c'), then the folder in question (here 'share') (which gives '/cygdrive/c/share' for 'C:\share').
- comment: a small comment line
- read only: set it to yes, because only backuppc should access it and it doesn't need specific write permissions
- hosts allow: specify the IP of the backuppc server
- secret file: contains a file with logins and passwords of users authorized to connect
- auth users: specifies which user is authorized to connect.

I have therefore created a specific user called backuppc on the machine and given it specific rights (security tab) to the "C:\share" folder

Now, we'll create a rsyncd.secrets file containing users and passwords. We'll need to authorize the backuppc user to connect to the rsyncd service:

```
user:password
```

The configuration is quite simple, which would give me for example: backuppc:password  
Once that's done, **restart the 'RsyncServer' service** in the services list.

On the backuppc server, the configuration of the server in question looks like this:

```
$Conf{XferMethod} = 'rsyncd';
$Conf{RsyncShareName} = [
  'factory'
];
$Conf{RsyncdPasswd} = 'password';
$Conf{RsyncdUserName} = 'backuppc';
```

Reload backuppc and you're good to go.

## Backing up SQL Databases

### MySQL

To backup SQL databases (MySQL for example), it's preferable to create an SQL account dedicated to backups (backuppc for example) and assign it select and lock rights on all databases:

```bash
$ mysql -uroot -p
CREATE USER 'backuppc'@'localhost' IDENTIFIED BY 'password';
GRANT SELECT , LOCK TABLES ON * . * TO 'backuppc'@'localhost' IDENTIFIED BY 'password' WITH MAX_QUERIES_PER_HOUR 0 MAX_CONNECTIONS_PER_HOUR 0 MAX_UPDATES_PER_HOUR 0 MAX_USER_CONNECTIONS 0;
FLUSH PRIVILEGES;
```

#### Backup All Databases at Once

The advantage of this method is simplicity, but it doesn't allow restoration database by database.

In your host configuration on backuppc, add this line and adapt it to your needs:

```
...
$Conf{DumpPreUserCmd} = '$sshPath -q -x -l root $host /usr/bin/mysqldump -ubackuppc -ppassword -e --single-transaction --opt --all-databases > /tmp/dump.sql';
...
```

Ideally, at the end of the backup, you should delete this dump (for security reasons):

```
...
$Conf{DumpPostUserCmd} = '$sshPath -q -x -l root $host rm -f /tmp/dump.sql';
...
```

#### Backup Database by Database

This more tedious method has the advantage of backing up database by database which allows you to restore only the database you're interested in in case of a problem.

Additionally, it includes on-the-fly compression of your database. However, you'll need to install 7zip first (I chose 7zip for better compression).

We'll create a script that we'll place in /etc/scripts for example:

```bash
#!/bin/bash
user='root'
password='password'
destination='/tmp/backups_sql'
mail='my@mail.com'

mkdir -p $destination
for i in `echo "show databases;" | mysql -u$user -p$password | grep -v Database`; do
       mysqldump -u$user -p$password --opt --add-drop-table --routines --triggers --events --single-transaction --master-data=2 -B $i | 7z a -t7z -mx=9 -si $destination/$i.sql.7z
done

problem_text=''
problem=0
for i in `ls $destination/*`; do
    size=`du -sk $i | awk '{ print $1 }'`
    if [ $size -le 4 ]; then
        problem_text="$problem_text- $i database. Backupped database size is equal or under 4k ($size)\n"
        problem=1
    fi
done

if [ $problem -ne 0 ]; then
    echo -e "Backups problem detected on:\n\n$problem_text" | mail -s "$HOSTNAME - MySQL backup problem" $mail
fi
```

The problem here is the password in clear text. So make sure to restrict to the user who will backup:

```bash
chmod 700 /etc/scripts/backup_mysql_databases.sh
```

In your host configuration on backuppc, add this line and adapt it to your needs:

```
...
$Conf{DumpPreUserCmd} = '$sshPath -q -x -l root $host /etc/scripts/backup_mysql_databases.sh';
...
```

Ideally, at the end of the backup, you should delete this dump (for security reasons):

```
...
$Conf{DumpPostUserCmd} = '$sshPath -q -x -l root $host rm -Rf /tmp/backups_sql';
...
```

### PostgreSQL

To backup Postgres databases, **we need to do as usual, an [SSH key exchange]({{< ref "docs/Linux/Network/OpenSSH/openssh_ssh_key_exchange.md" >}}) but for the postgres user.**

#### Backup All Databases at Once

The advantage of this method is simplicity, but it doesn't allow restoration database by database.

In your host configuration on backuppc, add this line and adapt it to your needs:

```
...
$Conf{DumpPreUserCmd} = '$sshPath -q -x -l postgres $host /usr/bin/pg_dump > /tmp/dump.sql';
...
```

Ideally, at the end of the backup, you should delete this dump (for security reasons):

```
...
$Conf{DumpPostUserCmd} = '$sshPath -q -x -l root $host rm -f /tmp/dump.sql';
...
```

#### Backup Database by Database

This more tedious method has the advantage of backing up database by database which allows you to restore only the database you're interested in in case of a problem.

Additionally, it includes on-the-fly compression of your database. However, you'll need to install 7zip first (I chose 7zip for better compression).

We'll create a script that we'll place in /etc/scripts for example:

```bash
#!/bin/bash
destination='/tmp/backups_pgsql'
mail='my@mail.fr'

mkdir -p $destination || echo -e "Backups problem detected on:\n\n$problem_text" | mail -s "Can't create $destination folder" $mail
for i in `psql -l | grep "^\ [a-zA-Z0-9]" | grep -v 'template[0|1]' | cut -d\| -f1`; do
       /usr/bin/pg_dump $i | 7z a -t7z -mx=9 -si $destination/$i.sql.7z
done

problem_text=''
problem=0
for i in `ls $destination/*`; do
    size=`du -sk $i | awk '{ print $1 }'`
    if [ $size -le 4 ]; then
        problem_text="$problem_text- $i database. Backupped database size is equal or under 4k ($size)\n"
        problem=1
    fi
done

if [ $problem -ne 0 ]; then
    echo -e "Backups problem detected on:\n\n$problem_text" | mail -s "$HOSTNAME - Postgres backup problem" $mail
fi
```

A little security doesn't hurt:

```bash
chmod 744 /etc/scripts/backup_postgres_databases.sh
chown postgres /etc/scripts/backup_postgres_databases.sh
```

In your host configuration on backuppc, add this line and adapt it to your needs:

```
...
$Conf{DumpPreUserCmd} = '$sshPath -q -x -l postgres $host /etc/scripts/backup_postgres_databases.sh';
...
```

Ideally, at the end of the backup, you should delete this dump (for security reasons):

```
...
 $Conf{DumpPostUserCmd} = '$sshPath -q -x -l root $host rm -Rf /tmp/backups_sql';
...
```

## Restoration by Script

Here's a script that allows you to do restoration:

```bash
#!/bin/bash

# Script for restoring hosts (last full backup) from command line.
# The restored backups can be found in $RESTOREDIR (defined below),
# and are to be written on tape.

BACKUPPCDIR=/srv/backuppc-data
HOSTSDIR=$BACKUPPCDIR/pc
RESTOREDIR=$HOSTSDIR/restore/restore

# put the hosts/directories you do not want to restore into egrep...
HOSTS=$(ls $HOSTSDIR | egrep -v '(HOST_CONFIG_FILES|restore)' | tr / " ")
# or use:
# HOSTS="HOST1 HOST2 REMOTE3"


# no need to change anything below...

DATE=$(date +%F)
mkdir -p $RESTOREDIR/$DATE

for HOST in $HOSTS
do
# find the last full backup
NUMBER=$(grep full $HOSTSDIR/$HOST/backups| tail -1 | cut -f1)

if [ "$NUMBER" ]
then

# do the backup for the host
$BACKUPPCDIR/bin/BackupPC_archiveHost $BACKUPPCDIR/bin/BackupPC_tarCreate /usr/bin/split /usr/bin/par2 \
            "$HOST" "$NUMBER" /usr/bin/gzip .gz 0000000 $RESTOREDIR/$DATE 0 \*
fi

done
```

## FAQ

### Problem Creating Link When Starting the Service

If you encounter this type of error message:

```
2008-04-20 17:55:46 Can't create a test hardlink between a file in /var/lib/backuppc/pc and /var/lib/backuppc/cpool.  Either these are different file systems, or this file system doesn't
support hardlinks, or these directories don't exist, or there is a permissions problem, or the file system is out of inodes or full.  Use df, df -i, and ls -ld to check each of these possibilities. Quitting...
```

Check the rights etc... otherwise, if you're using encfs cryptology, then it comes from that and I invite you to follow [this link](Encfs:_Mise_en_place_d'Encfs_avec_FUSE/#I_can.27t_create_hard_link)[this link]({{< ref "docs/Linux/FilesystemsAndStorage/StorageEncryption/encfs_setting_up_encfs_with_fuse.md#i-cant-create-hard-link" >}})

### I Lost the Backup Numbers When I Want to Restore Data

If you still have your data but your backups and backups.old files are corrupted, there's still a way to recreate the backup indexing to be able to recover the data. If you want to reindex everything, run this:

```bash
/usr/share/backuppc/bin/BackupPC_fixupBackupSummary
```

Otherwise, if you just want to do a single machine, add this at the end:

```bash
$ /usr/share/backuppc/bin/BackupPC_fixupBackupSummary -l localhost
Doing host localhost
   Reading /var/lib/backuppc/pc/localhost/0/backupInfo
   Reading /var/lib/backuppc/pc/localhost/1/backupInfo
   Reading /var/lib/backuppc/pc/localhost/2/backupInfo
   Reading /var/lib/backuppc/pc/localhost/204/backupInfo
   Reading /var/lib/backuppc/pc/localhost/206/backupInfo
   Reading /var/lib/backuppc/pc/localhost/207/backupInfo
```

If you have a Perl error, you're probably missing the Perl "Time::ParseDate" package, do this:

```bash
apt-get install libtime-modules-perl
```

## Resources

- [BackupPC Documentation](/pdf/backuppc.pdf)
- [Documentation on Backuppc on Debian](/pdf/backuppc.pdf)
