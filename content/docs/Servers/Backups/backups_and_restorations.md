---
weight: 999
url: "/Sauvegardes_et_Restaurations/"
title: "Backups and Restorations with tapes"
description: "A comprehensive guide to backup and restoration techniques, covering tape backups, incremental backups, and UFS snapshots"
categories: ["Linux", "Backup", "Database"]
date: "2008-04-05T17:33:00+02:00"
lastmod: "2008-04-05T17:33:00+02:00"
tags: ["cd ~", "tar", "mt", "ufs", "backup", "restore", "snapshots"]
toc: true
---

## Introduction

A crucial function of system administration is to backup file systems. Backups safeguard against data loss, damage, or corruption. Backup tapes are often referred to as dump tapes.

## rmt device

All tape drives have logical device names that you use to reference the device on the command line. The image shows the format that all logical device names use.

```bash
/dev/rmt/#hn
```

- \# : Logical Tape number
- h : Tape density (l,m,h,c,u)
- No rewind

The logical tape numbers in the tape drive names always start with 0. For example:

- The first instance of a tape drive:

```bash
/dev/rmt/0
```

- The second instance of a tape drive:

```bash
/dev/rmt/1
```

- The third instance of a tape drive:

```bash
/dev/rmt/2
```

Two optional parameters further define the logical device name:

- Tape density - Five values can be given in the tape device name: l (low), m (medium), h (high), c (compressed), or u (ultra compressed).
- No rewind - The letter n at the end of a tape device name indicates that the tape should not be rewound when the current operation completes.

Tape densities depend on the tape drive. Check the manufacturer's documentation to determine the correct densities for the tape media.

Tape drives that support data compression contain internal hardware that performs the compression. If you back up a software-compressed file to a tape drive with hardware compression, the resulting file may be larger in size.

## mt

You use the mt command (magnetic tape control) to send instructions to the tape drive. Not all tape drives support all mt commands.

The format for the mt command is:

```bash
mt -f tape-device-name command count
```

You use the -f option to specify the tape device name, typically a no-rewind device name. If no -f option is used, the default tape device file `/dev/rmt/0` is used.

### Using the mt Command

The table lists some of the mt commands that you can use to control a magnetic tape drive.

{{< table "table-hover table-striped" >}}
| Command | Definition |
|---------|------------|
| mt status | Displays status information about the tape drive |
| mt rewind | Rewinds the tape |
| mt offline | Rewinds the tape and, if appropriate, takes the drive unit offline and if the hardware supports it, unloads |
| mt fsf count | Moves the tape forward count records |
{{< /table >}}

Assuming the tape was rewound to the start of tape, the following command positions the tape at the beginning of the third tape record.

```bash
mt -f /dev/rmt/0n fsf 2
```

The most common method to schedule backups is to perform cumulative incremental backups daily. This schedule is recommended for most situations.

To set up a backup schedule, determine:

- The file systems to back up
- A backup device (for example, tape drive)
- The number of tapes to use for the backup
- The type of backup (for example, full or incremental)
- The procedures for marking and storing tapes
- The time it takes to perform a backup

Determining File System Names to Back Up

Display the contents of the `/etc/vfstab` file. Then view the mount point column to find the name of the file system that you want to back up.
Determining the Number of Tapes

You determine the number of tapes for a backup according to the size of the file system you are backing up.

To determine the size of the file system, use the ufsdump command with the S option. The following are the command formats:

```bash
# ufsdump 0S filesystem_name
<number reported>
```

or

```bash
# ufsdump 3S filesystem_name
<number reported>
```

The numeric option determines the appropriate dump level. The output is the estimated number of bytes that the system requires for a complete backup.

Divide the reported bytes by the capacity of the tape to determine how many tapes you need to backup the file system.
Determining Back Up Frequency and Levels

You determine how often and at what level to backup each file system. The level of a backup refers to the amount of information that is backed up.

## Identifying Incremental and Full Backups

You can perform a full backup or an incremental backup of a file system. A full backup is a complete file system backup. An incremental backup copies only files in the file system that have been added or modified since a previous lower-level backup.

You use dump level 0 to perform a full backup. You use dump levels 1 through 9 to schedule incremental backups. The level numbers have no meaning other than their relationship to each other as a higher or lower number.

The Explore shows an example of a file system backup performed in incremental levels.

The table defines the elements of the sample incremental backup strategy shown in The image.

{{< table "table-hover table-striped" >}}
| Level | Example |
|-------|---------|
| 0 (Full) | Performed once each month. |
| 3 | Performed every Monday. The backup copies new or modified files since the last lower-level backup (for example, 0). |
| 4 | Performed every Tuesday. The backup copies new or modified files since the last lower-level backup (for example, 3). |
| 5 | Performed every Wednesday. The backup copies new or modified files since the last lower-level backup (for example, 4). |
| 6 | Performed every Thursday. The backup copies new or modified files since the last lower-level backup (for example, 5). |
| 2 | Performed every Friday. The backup copies new or modified files since the last lower-level backup, which is the Level 0 backup at the beginning of the month. |
{{< /table >}}

_Note:_ Many system administrators use the crontab utility to start a script that runs the ufsdump command.

The `/etc/dumpdates` file records backups if the -u option is used with the ufsdump command. Each line in the `/etc/dumpdates` file shows the file system that was backed up and the level of the last backup. It also shows the day, the date, and the time of the backup.

The following is an example `/etc/dumpdates` file:

```bash
# cat /etc/dumpdates
/dev/rdsk/c0t2d0s6  0 Fri Nov 5  19:12:27  2004
/dev/rdsk/c0t2d0s0  0 Fri Nov 5  20:44:02  2004
/dev/rdsk/c0t0d0s7  0 Tue Nov 9  09:58:26  2004
/dev/rdsk/c0t0d0s7  1 Tue Nov 9  16:25:28  2004
```

When an incremental backup is performed, the ufsdump command consults the `/etc/dumpdates` file. It looks for the date of the next lower-level backup. Then, the ufsdump command copies to the backup media all of the files that were modified or added since the date of that lower-level backup.

When the backup is complete, the `/etc/dumpdates` file records a new entry that describes this backup. The new entry replaces the entry for the previous backup at that level.

You can view the `/etc/dumpdates` file to determine if the system is completing backups. If a backup does not complete because of equipment failure, the `/etc/dumpdates` file does not record the backup.

Note: When you are restoring an entire file system, check the `/etc/dumpdates` file for a list of the most recent dates and levels of backups. Use this list to determine which tapes are needed to restore the entire file system. The tapes should be physically marked with the dump level and date of the backup.

Check that the file system is inactive, or unmounted, before you back the system up. If the file system is active, the output of the backup can be inconsistent, and you could find it impossible to restore some of the files correctly.

## Backing Up an Unmounted File System

The standard Solaris OS command for ufs file system backups is `/usr/sbin/ufsdump`.

The format for the ufsdump command is:

```bash
ufsdump option(s) argument(s) filesystem_name
```

You can use this command to back up a complete or a partial file system. Backups are often referred to as dumps.

The table defines several common options for the ufsdump command.

{{< table "table-hover table-striped" >}}
| Option | Description |
|--------|-------------|
| 0-9 | Back up level. Level 0 is a full backup of the file system. Levels 1 through 9 are incremental backups of files that have changed since the last lower-level backup. When no backup level is given, the default is level 9. |
| v | Verify. After each tape is written, the system verifies the contents of the media against the source file system. If any discrepancies occur, the system prompts the operator to insert new media and repeat the process. Use this option only on an unmounted file system. Any activity in the file system causes the system to report discrepancies. |
| S | Size estimate. This option allows you to estimate the amount of space that will be needed on the tape to perform the level of backup you want. |
| l | Autoload. You use this option with an autoloading (stackloader) tape drive. |
| o | Offline. When the backup is complete, the system takes the drive offline, rewinds the tape (if you use a tape), and, if possible, ejects the media. |
| u | Update. The system creates an entry in the `/etc/dumpdates` file with the device name for the file system disk slice, the backup level (0-9), and the date. If an entry already exists for a backup at the same level, the system replaces the entry. |
| n | Notify. The system sends messages to the terminals of all logged-in users who are members of the sys group to indicate that the ufsdump command requires attention. |
| f device | Specify. The system specifies the device name of the file system backup. When you use the default tape device, `/dev/rmt/0`, you do not need the -f option. The system assumes the default. |
{{< /table >}}

You use the ufsdump command to create file system backups to tape. The dump level (0-9) specified in the ufsdump command determines which files to back up.

### Using the ufsdump Command

Perform the following steps to use the ufsdump command to start a tape backup:

- Become the root user to change the system to single-user mode, and unmount the file systems.

```bash
      # /usr/sbin/shutdown -y -g300 "System is being shutdown for backup"

      Shutdown started.    Mon Oct 11 12:22:33 BST 2004

      Broadcast Message from root (pts/1) on host1 Mon Oct 11 12:22:33...
      The system host1 will be shut down in 5 minutes
      System is being shutdown for backup
      (further output omitted)
```

- Verify that the /export/home file system was unmounted with the shutdown command. If not, unmount it manually.
- Check the integrity of the file system data with the fsck command.

```
# fsck /export/home
```

- Perform a full (Level 0) backup of the /export/home file system.

```bash
      # ufsdump 0uf /dev/rmt/0 /export/home
      # ufsdump 0uf /dev/rmt/0 /export/home
        DUMP: Writing 32 Kilobyte records
        DUMP: Date of this level 0 dump: Mon Oct 11 12:30:44 2004
        DUMP: Date of last level 0 dump: the epoch
        DUMP: Dumping /dev/rdsk/c0t0d0s7 (host1:/export/home) to /dev/rmt/0.
        DUMP: Mapping (Pass I) [regular files]
        DUMP: Mapping (Pass II) [directories]
        DUMP: Estimated 1126 blocks (563KB).
        DUMP: Dumping (Pass III) [directories]
        DUMP: Dumping (Pass IV) [regular files]
        DUMP: Tape rewinding
        DUMP: 1086 blocks (543KB) on 1 volume at 1803 KB/sec
        DUMP: DUMP IS DONE
        DUMP: Level 0 dump on Mon Oct 11 12:42:12 2004
```

You can use the ufsdump command to perform a backup on a remote tape device.

The format for the ufsdump command is:

```bash
ufsdump options remotehost:tapedevice filesystem
```

To perform remote backups across the network, the system with the tape drive must have an entry in its /.rhosts file for every system that uses the tape drive.
Using the ufsdump Command

The following example shows how to perform a full (Level 0) backup of the /export/home file system on the host1 system, to the remote tape device on the host2 system.

```bash
 # ufsdump 0uf host2:/dev/rmt/0 /export/home
  DUMP: Writing 32 Kilobyte records
  DUMP: Date of this level 0 dump: Mon Oct 11 13:30:44 2004
  DUMP: Date of last level 0 dump: the epoch
  DUMP: Dumping /dev/rdsk/c0t0d0s7 (host1:/export/home) to   host2:/dev/rmt/0.
  DUMP: Mapping (Pass I) [regular files]
  DUMP: Mapping (Pass II) [directories]
  DUMP: Estimated 320 blocks (160KB).
  DUMP: Dumping (Pass III) [directories]
  DUMP: Dumping (Pass IV) [regular files]
  DUMP: Tape rewinding
  DUMP: 318 blocks (159KB) on 1 volume at 691 KB/sec
  DUMP: DUMP IS DONE
  DUMP: Level 0 dump on Mon Oct 11 13:44:12 2004
```

## Restoring

When you are restoring data to a system, consider the following questions:

- Can the system boot on its own (regular file system restore)?
- Do you need to boot the system from CD-ROM, DVD, or network (critical file system restore)?
- Do you need to boot the system from CD-ROM, DVD, or network and repair the boot drive (special case recovery)?

To restore files or file systems, determine the following:

- The file system backup tapes that are needed
- The device name to which you will restore the file system
- The name of the temporary directory to which you will restore individual files
- The type of backup device to be used (local or remote)
- The backup device name (local or remote)

To restore a regular file system, such as the `/export/home` or `/opt` file system, back up to the disk, you use the ufsrestore command. The ufsrestore command copies files to the disk, relative to the current working directory, from backup tapes that were created by the ufsdump command.

You can use the ufsrestore command to reload an entire file system hierarchy from a Level 0 backup and related incremental backups. You can also restore one or more single files from any backup tape.

The format for the ufsrestore command is:

```bash
ufsrestore option(s) argument(s) filesystem
```

The table describes some options that you can use with the ufsrestore command.

{{< table "table-hover table-striped" >}}
| Option | Description |
|--------|-------------|
| t | Lists the table of contents of the backup media. |
| r | Restores the entire file system from the backup media. |
| x file1 file2 | Restores only the files named on the command line. |
| i | Invokes an interactive restore. |
| v | Specifies verbose mode. This mode displays the path names to the terminal screen as each file is restored. |
| f device | Specifies the tape device name. When not specified, the `/dev/rmt/0` device file is used. |
{{< /table >}}

When you restore an entire file system from a backup tape, the system creates a restoresymtable file. The ufsrestore command uses the restoresymtable file for check-pointing or passing information between incremental restores. You can remove the restoresymtable file when the restore is complete.

### ufsrestore

The following procedure demonstrates how to use the ufsrestore command to restore the /opt file system on the c0t0d0s5 slice.

- Create the new file system structure.

```bash
newfs /dev/rdsk/c0t0d0s5
```

- Mount the file system to the /opt directory, and change to that directory.

```bash
mount /dev/dsk/c0t0d0s5 /opt
cd /opt
```

- Restore the entire /opt file system from the backup tape.

```bash
ufsrestore rf /dev/rmt/0
```

Note: Always restore a file system by starting with the Level 0 backup tape, continuing with the next-lower-level tape, and continuing through the highest-level tape.

- Remove the restoresymtable file.

```bash
rm restoresymtable
```

- Unmount the new file system.

```bash
cd /
umount /opt
```

- Use the fsck command to check the restored file system.

```bash
fsck /dev/rdsk/c0t0d0s5
```

- Perform a full backup of the file system.

```bash
ufsdump 0uf /dev/rmt/0 /dev/rdsk/c0t0d0s5
```

Note: The system administrator should always back up the newly created file system because the ufsrestore command repositions the files and changes the inode allocation.

```bash
init 6
```

The ufsrestore i command invokes an interactive interface. Through the interface, you can browse the directory hierarchy of the backup tape and select individual files to extract. The term volume is used by ufsrestore and should be considered a single tape.
Using the ufsrestore i Command

**The following procedure demonstrates how to use the ufsrestore i command to extract individual files from a backup tape.**

- Become the root user, and change to the temporary directory that you want to receive the extracted files.

```bash
cd /export/home/tmp
```

- Perform the ufsrestore i command.

```bash
 $ ufsrestore ivf /dev/rmt/0
 Verify volume and initialize maps
 Media block size is 64
 Dump   date: Mon Oct 11 12:30:44 2004
 Dumped from: the epoch
 Level 0 dump of /export/home on sys43:/dev/dsk/c0t0d0s7
 Label: none
 Extract directories from tape
 Initialize symbol table.
```

- Display the contents of the directory structure on the backup tape.

```bash
 $ ufsrestore > ls
 .:
 2 *./            13  directory1    15  directory3     11  file2
 2 *../           14  directory2    10  file1          12  file3
```

- Change to the target directory on the backup tape.

```bash
 $ ufsrestore > cd directory1
 $ ufsrestore > ls
 ./directory1:
 3904  ./          2 *../      3905  file1    3906  file2    3907  file3
```

- Add the files you want to restore to the extraction list.

```bash
 $ ufsrestore > add file1 file2
 $ Make node ./directory1
```

Files you want to restore are marked with an asterisk (\*) for extraction. If you extract a directory, all of the directory contents are marked for extraction.

In this example, two files are marked for extraction. The ls command displays an asterisk in front of the selected file names, file1 and file2.

```bash
 $ ufsrestore > ls
 ./directory1:
 3904 *./          2 *../      3905 *file1    3906 *file2    3907  file3
```

- To delete a file from the extraction list, use the delete command.

```bash
ufsrestore > delete file1
```

The ls command displays the file1 file without an asterisk.

```bash
 $ ufsrestore > ls
 ./directory1:
 3904 *./          2 *../      3905  file1    3906 *file2    3907  file3
```

- To view the files and directories marked for extraction, use the marked command.

```bash
 $ ufsrestore > marked
 ./directory1:
 3904 *./          2 *../      3906 *file2
```

- To restore the selected files from the backup tape, perform the command:

```bash
 $ ufsrestore >> extract
 Extract requested files
 You have not read any volumes yet.
 Unless you know which volume your file(s) are on you should start
 with the last volume and work towards the first.
 Specify next volume #: 1
```

_Note:_ The ufsrestore command has to find the selected files. If you used more than one tape for the backup, first insert the tape with the highest volume number and type the appropriate number at this point. Then repeat, working towards Volume #1 until all files have been restored.

```bash
 extract file ./directory1/file2
 Add links
 Set directory mode, owner, and times.
 set owner/mode for ."? [yn] n
```

_Note:_ Answering y sets ownership and permissions of the temporary directory to those of the mount point on the tape.

- To exit the interactive restore after the files are extracted, perform the command:

```bash
ufsrestore> quit
```

- Move the restored files to their original or permanent directory location, and delete the files from the temporary directory.

```bash
 mv /export/home/tmp/directory1/file2 /export/home
 rm -r /export/home/tmp/directory1
```

_Note:_ You can use the help command in an interactive restore to display a list of available commands.

## Restoring a ufs File System

When performing incremental restores, start with the last volume and work towards the first. The system uses information in the restoresymtable file to restore incremental backups on top of the latest full backup.

The following procedure demonstrates how to restore the /export/home file system from incremental tapes.

Note: This procedure makes use of the interactive restore to assist in showing the concept of incremental restores. You would typically use a command, such as ufsrestore rf, for restoring entire file systems.

- View the contents of the /etc/dumpdates file for information about the /export/home file system.

```bash
 # more /etc/dumpdates |grep c0t0d0s7
 /dev/rdsk/c0t0d0s7       0 Wed Apr 07 09:55:34 2004
 /dev/rdsk/c0t0d0s7       1 Web Apr 07 09:57:30 2004
```

- Create the new file system structure for the /export/home file system.

```bash
newfs /dev/rdsk/c0t0d0s7
```

- Mount the file system and change to that directory.

```bash
 mount /dev/dsk/c0t0d0s7 /export/home
 cd /export/home
```

- Insert the Level 0 backup tape.
- Restore the /export/home file system from the backup tapes.

```bash
 # ufsrestore rvf /dev/rmt/0
 Verify volume and initialize maps
 Media block size is 64
 Dump   date: Wed Apr 07 09:55:34 2004
 Dumped from: the epoch
 Level 0 dump of /export/home on sys41:/dev/dsk/c0t0d0s7
 Label: none
 Begin level 0 restore
 Initialize symbol table.
 Extract directories from tape
 Calculate extraction list.
 Make node ./directory1
 Make node ./directory2
 Make node ./directory3
 Extract new leaves.
 Check pointing the restore
 extract file ./file1
 extract file ./file2
 extract file ./file3
 Add links
 Set directory mode, owner, and times.
 Check the symbol table.
 Check pointing the restore
```

- Load the next lower-level tape into the tape drive:

```bash
 # ufsrestore rvf /dev/rmt/0
 Verify volume and initialize maps
 Media block size is 64
 Dump   date: Wed Apr 07 09:57:30 2004
 Dumped from: Wed Apr 07 09:55:34 2004
 Level 1 dump of /export/home on sys41:/dev/dsk/c0t0d0s7
 Label: none
 Begin incremental restore
 Initialize symbol table.
 Extract directories from tape
 Mark entries to be removed.
 Calculate node updates.
 Make node ./directory4
 Make node ./directory5
 Make node ./directory6
 Find unreferenced names.
 Remove old nodes (directories).
 Extract new leaves.
 Check pointing the restore
 extract file ./file4
 extract file ./file5
 extract file ./file6
 Add links
 Set directory mode, owner, and times.
 Check the symbol table.
 Check pointing the restore
```

### Alternative Steps

The following steps are an alternative to the previous Steps 5 and 6.

- Restore the /export/home file system from the backup tapes. (This example uses an interactive, verbose restore to provide more detailed information.)

```bash
 # ufsrestore ivf /dev/rmt/0
 Verify volume and initialize maps
 Media block size is 64
 Dump   date: Mon Oct 11 13:10:12 2004
 Dumped from: the epoch
 Level 0 dump of /export/home on sys41:/dev/dsk/c0t0d0s7
 Label: none
 Extract directories from tape
 Initialize symbol table.
```

```bash
 $ ufsrestore > ls
 .:
      2 *./            8  directory2     5  file2
      2 *../           9  directory3     6  file3
      7  directory1    4  file1          3  lost+found/

 The system lists files from the last Level 0 backup.
```

```bash
 $ ufsrestore > add *
 Warning: ./lost+found: File exists
```

```bash
 $ ufsrestore > extract
 Extract requested files
 You have not read any volumes yet.
 Unless you know which volume your file(s) are on you should start
 with the last volume and work towards the first.
 Specify next volume #: 1
 extract file ./file1
 extract file ./file2
 extract file ./file3
 extract file ./directory1
 extract file ./directory2
 extract file ./directory3
 Add links
 Set directory mode, owner, and times.
 set owner/mode for '.'? [yn] n
 Directories already exist, set modes anyway? [yn] n
 ufsrestore > q
```

- The information in the /etc/dumpdates file shows an incremental backup that was taken after the Level 0 backup. Load the next tape and perform the incremental restore.

```bash
 $ ufsrestore iv
 Verify volume and initialize maps
 Media block size is 64
 Dump   date: Wed Apr 07 09:57:30 2004
 Dumped from: Wed Apr 07 09:55:34 2004
 Level 1 dump of /export/home on sys41:/dev/dsk/c0t0d0s7
 Label: none
 Extract directories from tape
 Initialize symbol table.
```

```bash
 $ ufsrestore > ls
 .:
      2 *./        13  directory4    15  directory6    11  file5
      2 *../       14  directory5    10  file4         12  file6

 $ ufsrestore > add *

 $ ufsrestore > extract
 Extract requested files
 You have not read any volumes yet.
 Unless you know which volume your file(s) are on you should start
 with the last volume and work towards the first.
 Specify next volume #: 1
 extract file ./file4
 extract file ./file5
 extract file ./file6
 extract file ./directory4
 extract file ./directory5
 extract file ./directory6
 Add links
 Set directory mode, owner, and times.
 set owner/mode for '.'? [yn] n
 $ ufsrestore > q
```

## Creating a UFS Snapshot

The UFS Copy on Write Snapshots feature provides administrators an online backup solution for ufs file systems. This utility enables you to use a point-in-time copy of a ufs file system, called a snapshot, to create an online backup. You can create the backup while the file system is mounted and the system is in multiuser mode.

_Note:_ The UFS snapshots are similar to the Sun StorEdge Instant Image product. Instant Image allocates space equal to the size of the entire file system that is being captured. However, the file system data saved by UFS snapshots occupies only as much disk space as needed.

You use the fssnap command to create, query, or delete temporary read-only snapshots of ufs file systems.

The format for the fssnap command is:

```bash
/usr/sbin/fssnap -F FSType -V -o special_option(s) mount-point
```

The table shows some of the options for the fssnap command.

{{< table "table-hover table-striped" >}}
| Option | Description |
|--------|-------------|
| -d | Deletes the snapshot associated with the given file system. If the -o unlink option was used when you built the snapshot, the backing-store file is deleted together with the snapshot. Otherwise, the backing-store file (which contains file system data) occupies disk space until you delete it manually. |
| -F FSType | Specifies the file system type to be used. |
| -i | Displays the state of an FSType snapshot. |
| -V | Echoes the complete command line but does not execute the command. |
| -o | Enables you to use special_options, such as the location and size of the backing-store (bs) file. |
{{< /table >}}

To create a UFS snapshot, specify a backing-store path and the actual file system to be captured. The following is the command format:

```bash
fssnap -F ufs -o bs=backing_store_path /file-system
```

_Note:_ The backing_store_path can be a raw device, the name of an existing directory, or the name of a file that does not already exist.

The following example uses the fssnap command to create a snapshot of the /export/home file system.

```bash
# fssnap -F ufs -o bs=/var/tmp /export/home
/dev/fssnap/0
```

The snapshot subsystem saves file system data in a file called a backing-store file before the data is overwritten. Some important aspects of a backing-store file are:

- A backing-store file is a bit-mapped file that takes up disk space until you delete the UFS snapshot.
- The size of the backing-store file varies with the amount of activity on the file system being captured.
- The destination path that you specify on the fssnap command line must have enough free space to hold the backing-store file.
- The location of the backing-store file must be different from that of the file system you want to capture in a UFS snapshot.
- A backing-store file can reside on different types of file systems, including another ufs file system or a mounted nfs file system.

The fssnap command creates the backing-store file and two read-only virtual devices. The block virtual device, `/dev/fssnap/0`, can be mounted as a read-only file system. The raw virtual device, `/dev/rfssnap/0`, can be used for raw read-only access to a file system.

These virtual devices can be backed up with any of the existing Solaris OS backup commands. The backup created from a virtual device is a backup of the original file system when the UFS snapshot was taken.

Note: When a UFS snapshot is first created, the file system locks temporarily. Users might notice a slight pause when writing to this file system. The length of the pause increases with the size of the file system. There is no performance impact when users are reading from the file system.

Before creating a UFS snapshot, use the df -k command to check that the backing-store file has enough disk space to grow. The size of the backing-store file depends on how much data has changed since the previous snapshot was taken.

You can limit the size of the backing-store file by using the -o maxsize=n option of the fssnap command, where n (k, m, or g) is the maximum size of the backing-store file specified in Kbytes, Mbytes, or Gbytes.

Caution: If the backing-store file runs out of disk space, the system automatically deletes the UFS snapshot, which causes the backup to fail. The active ufs file system is not affected. Check the `/var/adm/messages` file for possible UFS snapshot errors.

Note: You can force an unmount of an active ufs file system, for which a snapshot exists (for example, with the umount -f command). This action deletes the appropriate snapshot automatically.

The following example creates a snapshot of the /export/home file system, and limits the backing-store file to 500 Mbytes.

```bash
# fssnap -F ufs -o bs=/var/tmp,maxsize=500m /export/home
/dev/fssnap/0
```

You can use either fssnap command to display UFS snapshot information.

The following example displays a list of all the current UFS snapshots on the system. The list also displays the corresponding virtual device for each snapshot.

```bash
# fssnap -i
0    /export/home
1    /usr
2    /database
```

You use the -i option to the `/usr/lib/fs/ufs/fssnap` command to display detailed information for a specific UFS snapshot that was created by the fssnap command.

The following example shows the details for the /export/home snapshot.

```bash
# /usr/lib/fs/ufs/fssnap -i /export/home
Snapshot number              : 0
Block Device                 : /dev/fssnap/0
Raw Device                   : /dev/rfssnap/0
Mount point                  : /export/home
Device state                 : idle
Backing store path           : /var/tmp/snapshot0
Backing store size           : 0 KB
Maximum backing store size   : 512000 KB
Snapshot create time         : Mon Oct 11 08:58:33 2004
Copy-on-write granularity    : 32 KB
```

### tar

You can use the tar command or the ufsdump command to back up a UFS snapshot.
Using the tar Command to Back Up a Snapshot File

If you use the tar command to back up the UFS snapshot, mount the snapshot before backing it up. The following procedure demonstrates how to do this type of mount.

- Create the mount point for the block virtual device.

```bash
mkdir -p /backups/home.bkup
```

- Mount the block virtual device to the mount point.

```bash
mount -F ufs -o ro /dev/fssnap/0 /backups/home.bkup
```

- Change directory to the mount point.

```bash
cd /backups/home.bkup
```

- Use the tar command to write the data to tape.

```bash
tar cvf /dev/rmt/0 .
```

Using the ufsdump Command

If you use the ufsdump command to back up a UFS snapshot, you can specify the raw virtual device during the backup.

```bash
ufsdump 0uf /dev/rmt/0 /dev/rfssnap/0
```

Verify that the UFS snapshot is backed up.

```bash
ufsrestore tf /dev/rmt/0
```

### Incrementals restores

Incremental backups of snapshots contain files that have been modified since the last UFS snapshot. You use the ufsdump command with the N option to create an incremental UFS snapshot, which writes the name of the device being backed up, rather than the name of the snapshot device to the `/etc/dumpdates` file.

The following example shows how to use the ufsdump command to create an incremental backup of a file system.

Note: It is important to note the use of the N argument when backing up a snapshot. This argument ensures proper updates to the `/etc/dumpdates` file.

```bash
ufsdump 1ufN /dev/rmt/0 /dev/rdsk/c1t0d0s0 /dev/rfssnap/0
```

Next you would verify that the UFS snapshot is backed up to tape.

```bash
ufsrestore tf /dev/rmt/0
```

To understand incremental backups of snapshots, consider the following demonstration:

- Create a snapshot of the /extra file system that is going to be backed up while the file system is mounted.

```bash
 # fssnap -o bs=/var/tmp /extra
 /dev/fssnap/0
```

- Verify that the snapshot was successful, and view detailed information about the snapshot.

```bash
 # fssnap -i
   0    /extra
 # /usr/lib/fs/ufs/fssnap -i /extra
 Snapshot number              : 0
 Block Device                 : /dev/fssnap/0
 Raw Device                   : /dev/rfssnap/0
 Mount point                  : /extra
 Device state                 : idle
 Backing store path           : /var/tmp/snapshot0
 Backing store size           : 0 KB
 Maximum backing store size   : Unlimited
 Snapshot create time         : Mon Oct 11 10:34:21 2004
 Copy-on-write granularity    : 32 KB
```

- Make a directory that will be used to mount and view the snapshot data.

```bash
 mkdir /extrasnap
```

- Mount the snapshot to the new mount point, and compare the size of the file system and the snapshot device.

```bash
 # mount -o ro /dev/fssnap/0 /extrasnap
 # df -k |grep extra
 /dev/dsk/c1t0d0s0    1294023       9 1242254     1%    /extra
 /dev/fssnap/0        1294023       9 1242254     1%    /extrasnap
```

- Edit a file under the /extra directory and make it larger, and then compare the size of the file system and the snapshot device.

```bash
 # vi file1
 (yank and put text, or read text in from another file)
 # df -k |grep extra
 /dev/dsk/c1t0d0s0    1294023      20 1242243     1%    /extra
 /dev/fssnap/0        1294023       9 1242254     1%    /extrasnap
```

Observe that the file system grew in size while the snapshot file did not.

- Perform a full backup with the N option of the ufsdump command.

```bash
 # ufsdump 0ufN /dev/rmt/0 /dev/rdsk/c1t0d0s0 /dev/rfssnap/0
 DUMP: Writing 32 Kilobyte records
 DUMP: Date of this level 0 dump: Mon Oct 11 10:49:38 2004
 DUMP: Date of last level 0 dump: the epoch
 DUMP: Dumping /dev/rfssnap/0 (sys41:/extrasnap) to /dev/rmt/0.
 DUMP: Mapping (Pass I) [regular files]
 DUMP: Mapping (Pass II) [directories]
 DUMP: Estimated 262 blocks (131KB).
 DUMP: Dumping (Pass III) [directories]
 DUMP: Dumping (Pass IV) [regular files]
 DUMP: Tape rewinding
 DUMP: 254 blocks (127KB) on 1 volume at 1814 KB/sec
 DUMP: DUMP IS DONE
 DUMP: Level 0 dump on Mon Oct 11 11:03:46 2004
```

- Verify the backup.

```bash
 # ufsrestore tf /dev/rmt/0
         2      .
         3      ./file1
         4      ./file2
         5      ./file3
         6      ./file4
```

- Unmount the back up device and remove the snapshot.

```bash
 umount /extrasnap
 fssnap -d /extra
 rm /var/tmp/snapshot0
```

- Make some changes to the /extra file system, such as copying some files, and then re-create the snapshot.

```bash
 # cp file1 file5
 # cp file1 file6
 # fssnap -o bs=/var/tmp /extra
 /dev/fssnap/0
```

- Re-mount the snapshot device, and compare the size of the file system and the snapshot device.

```bash
 # mount -o ro /dev/fssnap/0 /extrasnap
 # df -k |grep extra
 /dev/dsk/c1t0d0s0    1294023      46 1242217     1%    /extra
 /dev/fssnap/0        1294023      46 1242217     1%    /extrasnap
```

- Perform an incremental backup with the N option of the ufsdump command.

```bash
 # ufsdump 1ufN /dev/rmt/0 /dev/rdsk/c1t0d0s0 /dev/rfssnap/0
 DUMP: Writing 32 Kilobyte records
 DUMP: Date of this level 0 dump: Mon Oct 11 13:13:03 2004
 DUMP: Date of last level 0 dump: Mon Oct 11 12:30:44 2004
 DUMP: Dumping /dev/rfssnap/0 (sys41:/extrasnap) to /dev/rmt/0.
 DUMP: Mapping (Pass I) [regular files]
 DUMP: Mapping (Pass II) [directories]
 DUMP: Estimated 294 blocks (147KB).
 DUMP: Dumping (Pass III) [directories]
 DUMP: Dumping (Pass IV) [regular files]
 DUMP: Tape rewinding
 DUMP: 254 blocks (127KB) on 1 volume at 1693 KB/sec
 DUMP: DUMP IS DONE
 DUMP: Level 1 dump on Mon Oct 11 13:22:36 2004
```

- Verify the backup.

```bash
 # ufsrestore tf /dev/rmt/0
         2      .
         7      ./file5
         8      ./file6
```

Notice that the backup of the snapshot contains only the files that were added since the previous Level 0 backup.

The backup created from a virtual device is a backup of the original file system when the UFS snapshot was taken.

You restore a UFS snapshot from a backup tape in the same manner as you would the backup of an original file system. Data written to a tape by ufsdump is simply data, whether it is a snapshot or a file system.

To restore the demo directory from the snapshot backup of the /usr file system, complete the following steps:

- Load the tape that contains the snapshot backup of the /usr file system into the tape drive.
- Change to the /usr file system.

```bash
cd /usr
```

- Perform the a ufsrestore command.

```bash
 # ufsrestore if /dev/rmt/0
 ufsrestore > add demo
 ufsrestore > extract
 Specify next volume #: 1
 set owner/mode for '.'? [yn] n
 ufsrestore > quit
```

- Verify that the demo directory exists, and eject the tape.

Deleting a UFS snapshot from the system is a multistep process and order-dependant. First, unmount the snapshot device, and then delete the snapshot. Finally, remove the backing-store file.

```bash
 umount /dev/fssnap/0
 fssnap -d /export/home
 rm /backing_store_file
```
