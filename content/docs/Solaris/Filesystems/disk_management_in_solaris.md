---
weight: 999
url: "/Management_des_disques_sous_Solaris/"
title: "Disk Management in Solaris"
description: "A comprehensive guide on managing disks in Solaris, covering physical disk structure, partitioning, filesystems, and troubleshooting."
categories: ["Linux", "Unix", "System Administration"]
date: "2012-01-30T11:07:00+02:00"
lastmod: "2012-01-30T11:07:00+02:00"
tags: ["Solaris", "Disks", "Filesystem", "Partitioning", "VTOC", "UFS", "Format", "Mount"]
toc: true
---

## Introduction

Compared to Linux, Solaris is quite similar except for certain aspects which I will clarify here. This guide will not only cover Solaris but also include general information about disk architecture and filesystems.

## Physical Operation

### Files, Inodes and Blocks

On your hard drive, where you store your data, there is a hierarchical structure:

* Directories and files
* Inodes
* Blocks

Inodes are what know exactly where each directory/file is located. When you create or call a file, it points to an inode. This inode is then able to identify which data block it belongs to (binary slice). Here's a small explanation with an image:

![Hdd function](/images/hdd_fonction.avif)

### Hardware Recognition

To recognize where a specific device is located on our system, here's a brief explanation:
Assuming we have a file in /dev like: c0t0d0s0. This means:

* c0: On **C**ontroller **0**
* t0: On SCSI **T**arget **0**
* d0: I have **D**isk **0**
* s0: And I am positioned on **S**lice **0**

Here are some examples:

#### SCSI

![Sun hdd scsi](/images/sun_hdd_scsi.avif)

#### IDE

![Sun hdd ide](/images/sun_hdd_ide.avif)

### Slices

Then, the slices (also called partitions) are defined as follows for better performance optimization:

![Sun slice](/images/sun_slice.avif)

## What About My Machine?

To start, you should know what's on your machine. For this, there are 3 solutions:

### path_to_inst

```bash
cat /etc/path_to_inst
```

```
#
#       Caution! This file contains critical kernel state
#
"/pseudo" 0 "pseudo"
"/options" 0 "options"
"/xsvc" 0 "xsvc"
"/objmgr" 0 "objmgr"
"/scsi_vhci" 0 "scsi_vhci"
"/isa" 0 "isa"
"/isa/i8042@1,60" 0 "i8042"
"/isa/i8042@1,60/keyboard@0" 0 "kb8042"
"/isa/i8042@1,60/mouse@1" 0 "mouse8042"
"/isa/lp@1,378" 0 "ecpp"
"/isa/asy@1,3f8" 0 "asy"
"/isa/asy@1,2f8" 1 "asy"
"/isa/fdc@1,3f0" 0 "fdc"
"/isa/fdc@1,3f0/fd@0,0" 0 "fd"
"/ramdisk" 0 "ramdisk"
"/pci@0,0" 0 "pci"
"/pci@0,0/display@f" 0 "vgatext"
"/pci@0,0/pci8086,7191@1" 0 "pci_pci"
"/pci@0,0/pci-ide@7,1" 0 "pci-ide"
"/pci@0,0/pci-ide@7,1/ide@0" 0 "ata"
"/pci@0,0/pci-ide@7,1/ide@0/cmdk@0,0" 0 "cmdk"
"/pci@0,0/pci-ide@7,1/ide@1" 1 "ata"
"/pci@0,0/pci-ide@7,1/ide@1/sd@0,0" 16 "sd"
"/pci@0,0/pci1000,30@10" 0 "mpt"
"/pci@0,0/pci1022,2000@11" 0 "pcn"
"/iscsi" 0 "iscsi"
```

### Prtconf

```bash
prtconf | grep -v not
```

```
System Configuration:  Sun Microsystems  i86pc
Memory size: 512 Megabytes
System Peripherals (Software Nodes):
 
i86pc
    scsi_vhci, instance #0
    isa, instance #0
        i8042, instance #0
            keyboard, instance #0
            mouse, instance #0
        fdc, instance #0
    pci, instance #0
        pci8086,7191, instance #0
        pci-ide, instance #0
            ide, instance #0
                cmdk, instance #0
            ide, instance #1
                sd, instance #16
        display, instance #0
        pci1000,30, instance #0
        pci1022,2000, instance #0
    iscsi, instance #0
    pseudo, instance #0
    options, instance #0
    xsvc, instance #0
    objmgr, instance #0
```

### Format

```bash
format
```

```
Searching for disks...done


AVAILABLE DISK SELECTIONS:
       0. c0d0 <DEFAULT cyl 4174 alt 2 hd 255 sec 63>
          /pci@0,0/pci-ide@7,1/ide@0/cmdk@0,0
```

## Adding a Device

### With Reboot

To have our disk detected at boot, we need to create a `/reconfigure` file:

```bash
touch /reconfigure
```

Then, simply connect your device and restart your machine. Once done, configure the slices.

### Without Reboot

If it's a "critical" machine, you need to run the "devfsadm" command. This command will try to match the loaded kernel drivers with the devices in `/devices`.

Here are some usage examples:

* Defining devices such as disk, tape, port, audio or pseudo:

```bash
devfsadm -c disk -c tape -c audio
```

* Configure only one device based on the driver:

```bash
devfsadm -i driver_name
```

* Configure disks only supported by certain controllers (dad, st or sd)

```bash
devfsadm -i dad
```

* For verbose mode

```bash
devfsadm -v
```

* To flush (clear) symbolic links that point to non-existent devices:

```bash
devfsadm -c
```

## Partitioning

The disk partitioning is made of slices that are delimited by cylinders. Indeed, a slice occupies a strip of cylinders (ex: 1 to 2850). Then the next slice will go from 2850 to 5000.

The partitions are therefore determined from the first cylinder of each slice:

* Slice 1: Cylinder 0 to 2850
* Slice 2: Cylinder 2850 to 5000

...

### Waste

Wasting cylinders creates a potentially empty slice. You can use it later. However, in case of defective sectors, it is possible that the partition will shrink due to lost cylinders.

### Overlapping

Overlapping occurs when multiple slices access the same cylinder (usually one slice encroaching on another). To fix this problem, when you edit your partition, use the "modify" command:

```bash
modify
```

```
        Select partitioning base:
                0. Current partition table (unnamed)
                1. All Free Hog
        Choose base (enter number) [0]? 0
        Warning: Overlapping partition (1) in table.
        Warning: Fix, or select a different partition table.
```

### Defining Partitions

The format command automatically handles partitioning according to data in `/etc/format.dat`. The advantage is that it's super fast and easy when adding a disk. Now, manually, here's how to do it:

```bash
format
```

```
Searching for disks...done

AVAILABLE DISK SELECTIONS:
       0. c0t0d0 <ST38410A cyl 16706 alt 2 hd 16 sec 63>
          /pci@1f,0/pci@1,1/ide@3/dad@0,0
       1. c1t3d0 <SUN9.0G cyl 4924 alt 2 hd 27 sec 133>
          /pci@1f,0/pci@1/scsi@1/sd@3,0
Specify disk (enter its number):
```

Now we have the list of detected disks. We'll choose the second disk and continue:

```
Specify disk (enter its number): 1
selecting c1t3d0
[disk formatted]

FORMAT MENU:
        disk       - select a disk
        type       - select (define) a disk type
        partition  - select (define) a partition table
        current    - describe the current disk
        format     - format and analyze the disk
        repair     - repair a defective sector
        label      - write label to the disk
        analyze    - surface analysis
        defect     - defect list management
        backup     - search for backup labels
        verify     - read and display labels
        save       - save new disk/partition definitions
        inquiry    - show vendor, product and revision
        scsi       - independent SCSI mode selects
        cache      - enable, disable or query SCSI disk cache
        volname    - set 8-character volume name
        !<cmd>     - execute <cmd>, then return
        quit
format>
```

In the menus that we can see, change or confirm partition choices, we have:

{{< table "table-hover table-striped" >}}
| Elements | Functions |
|---------|----------|
| partition | Displays the partition menu |
| label | Writes the current partition name list to the disk |
| verify | Reads and displays disk names |
| quit | Exit the format utility |
{{< /table >}}

Then type **partition** to see the menu:

```bash
format> partition
```

```
PARTITION MENU:
        0      - change '0' partition
        1      - change '1' partition
        2      - change '2' partition
        3      - change '3' partition
        4      - change '4' partition
        5      - change '5' partition
        6      - change '6' partition
        7      - change '7' partition
        select - select a predefined table
        modify - modify a predefined partition table
        name   - name the current table
        print  - display the current table
        label  - write partition map and label to the disk
        !<cmd> - execute <cmd>, then return
        quit
```

Here are the options offered:

{{< table "table-hover table-striped" >}}
| Elements | Functions |
|---------|----------|
| 0-7 | Specify partition size and offset |
| select | Choose a predefined slice in /etc/format.dat |
| modify | Change current partition in the table |
| quit | Used to identify the partition table in /etc/format.dat |
| print | Display the current partition table |
| label | Write the current partition table |
| !<cmd> | Execute an external command at the shell level |
{{< /table >}}

To display the new partition table, type **print**:

```bash
partition> print
```

```
Current partition table (original):
Total disk cylinders available: 4924 + 2 (reserved cylinders)

Part      Tag    Flag   Cylinders     Size            Blocks
  0 unassigned    wm     0            0         (0/0/0)           0
  1 unassigned    wm     0            0         (0/0/0)           0
  2     backup    ru     0 - 4923     8.43GB    (4924/0/0) 17682084
  3 unassigned    wu     0            0         (0/0/0)           0
  4 unassigned    wm     0            0         (0/0/0)           0
  5 unassigned    wm     0            0         (0/0/0)           0
  6 unassigned    wu     0            0         (0/0/0)           0
  7 unassigned    wm     0            0         (0/0/0)           0
```

Here is the meaning of the columns:

{{< table "table-hover table-striped" >}}
| Column Name | Description |
|---------|----------|
| Part | Slice number of the disk |
| Tag | Predefined tag (optional) |
| Flag | Predefined flag (optional) |
| Cylinders | Start and end cylinder of the slice |
| Size | Size of the slice in blocks (b), cylinders (c), Mbytes (MB), or Gbytes (GB) |
| Blocks | Total number of cylinders and sectors per slice |
{{< /table >}}

To start configuring the disk, type **0**:

```bash
partition> 0
```

```
Part      Tag    Flag  Cylinders    Size        Blocks
  0 unassigned    wm    0           0     (0/0/0)           0
```

Type **?** to get the list of possible choices:

```bash
Enter partition id tag[unassigned]: ?
```

```
Expecting one of the following: (abbreviations ok):
        unassigned    boot          root          swap
        usr           backup        stand         var
        home          alternates    reserved
Enter partition id tag[unassigned]:
```

Type **alternates**:

```bash
Enter partition id tag[unassigned]: alternates
```

Type **?** to get the list of possible choices:

```bash
Enter partition permission flags[wm]: ?
```

```
Expecting one of the following: (abbreviations ok):
        wm    - read-write, mountable
        wu    - read-write, unmountable
        rm    - read-only, mountable
        ru    - read-only, unmountable

Enter partition permission flags[wm]:
```

Press the "Enter" key:

```bash
Enter partition permission flags[wm]:
```

Press "Enter" again to accept cylinder 0 as starting point:

```bash
Enter new starting cyl[0]:
```

Enter the size of the partition (here 980mb):

```bash
Enter partition size[0b, 0c, 0e, 0.00mb, 0.00gb]: 980mb
```

Let's check:

```bash
partition> print
```

```
Current partition table (unnamed):
Total disk cylinders available: 1965 + 2 (reserved cylinders)

Part      Tag    Flag   Cylinders     Size            Blocks
  0 alternates    wm     0 -  558   980.16MB    (559/0/0)   200736
  1 unassigned    wm     0            0         (0/0/0)          0
  2     backup    ru     0 - 4923     8.43GB    (4924/0/0) 17682084
  3 unassigned    wm     0            0         (0/0/0)          0
  4 unassigned    wm     0            0         (0/0/0)          0
  5 unassigned    wm     0            0         (0/0/0)          0
  6 unassigned    wu     0            0         (0/0/0)          0
  7 unassigned    wm     0            0         (0/0/0)          0
```

We can see the changes. Let's adjust the start cylinder of slice 1:

```bash
partition> 1
```

```
Part      Tag    Flag   Cylinders    Size        Blocks
  1 unassigned    wm     0           0     (0/0/0)          0
```

Enter "**swap**":

```bash
Enter partition id tag[unassigned]: swap
```

Type "**wu**":

```bash
Enter partition permission flags[wm]: wu
```

Enter the start cylinder of slice 1:

```bash
Enter new starting cyl[0]: 559
```

Enter the new size of the partition:

```bash
Enter partition size[0b, 0c, 603e, 0.00mb, 0.00gb]: 512mb
```

Let's check:

```bash
partition> print
```

```
Current partition table (unnamed):
Total disk cylinders available: 1965 + 2 (reserved cylinders)

Part      Tag    Flag   Cylinders     Size          Blocks
  0 alternates    wm     0 -  558   980.16MB  (559/0/0)   2007369
  1       swap    wu   559 -  851   513.75MB  (293/0/0)   1052163
  2     backup    ru     0 - 4923     8.43GB  (4924/0/0) 17682084
  3 unassigned    wm     0            0       (0/0/0)          0
  4 unassigned    wm     0            0       (0/0/0)          0
  5 unassigned    wm     0            0       (0/0/0)          0
  6 unassigned    wu     0            0       (0/0/0)          0
  7 unassigned    wm     0            0       (0/0/0)          0
```

Let's do the same for slice 7:

```bash
partition> 7
```

```
Part      Tag    Flag     Cylinders        Size            Blocks
  7 unassigned    wm       0               0         (0/0/0) 0
```

Type "**home**":

```bash
Enter partition id tag[unassigned]: home
```

Press the "Enter" key:

```bash
Enter partition permission flags[wm]:
```

Enter the starting cylinder:

```bash
Enter new starting cyl[0]: 852
```

Enter the value "**$**" to occupy all available free space in this partition:

```bash
Enter partition size[0b, 0c, 694e, 0.00mb, 0.00gb]: $
```

Let's check:

```bash
partition> print
```

```
Current partition table (unnamed):
Total disk cylinders available: 1965 + 2 (reserved cylinders)

Part      Tag    Flag   Cylinders      Size            Blocks
  0 alternates    wm     0 -  558    980.16MB    (559/0/0)   2007369
  1       swap    wu   559 -  851    513.75MB    (293/0/0)   1052163
  2     backup    ru     0 - 4923      8.43GB    (4924/0/0) 17682084
  3 unassigned    wm     0             0         (0/0/0)          0
  4 unassigned    wm     0             0         (0/0/0)          0
  5 unassigned    wm     0             0         (0/0/0)          0
  6 unassigned    wu     0             0         (0/0/0)          0
  7       home    wm   852 - 4923      6.97GB    (4072/0/0) 14622552
```

After checking that there are no errors, type **label**:

```bash
partition> label
Ready to label disk, continue? y
```

### Checking Labels

To check labels (also called VTOC), there are 2 solutions:

* Use **verify** in the **format** utility
* Use the **prtvtoc** command

#### Reading the VTOC with Format

Open the format utility, then type verify:

```bash
format> verify
```

```
Primary label contents:

Volume name = <        >
ascii name  = <SUN9.0G cyl 4924 alt 2 hd 27 sec 133>
pcyl        = 4926
ncyl        = 4924
acyl        =    2
nhead       =   27
nsect       =  133
Part      Tag    Flag     Cylinders        Size            Blocks
  0 alternates    wm       0 -  558      980.16MB    (559/0/0)   2007369
  1       swap    wu     559 -  851      513.75MB    (293/0/0)   1052163
  2     backup    ru       0 - 4923        8.43GB    (4924/0/0) 17682084
  3 unassigned    wu       0               0         (0/0/0)           0
  4 unassigned    wm       0               0         (0/0/0)           0
  5 unassigned    wm       0               0         (0/0/0)           0
  6 unassigned    wu       0               0         (0/0/0)           0
  7       home    wm     852 - 4923        6.97GB    (4072/0/0) 14622552
```

To quit, type **q**.

#### Reading the VTOC with Prtvtoc

Run the command on a disk:

```bash
prtvtoc /dev/dsk/c1t3d0s0
```

```
* /dev/dsk/c1t3d0s0 partition map
*
* Dimensions:
*     512 bytes/sector
*     133 sectors/track
*      27 tracks/cylinder
*    3591 sectors/cylinder
*    4926 cylinders
*    4924 accessible cylinders
*
* Flags:
*   1: unmountable
*  10: read-only
*
*                          First     Sector    Last
* Partition  Tag  Flags    Sector     Count    Sector  Mount Directory
       0      9    00          0   2007369   2007368
       1      3    01    2007369   1052163   3059531
       2      5    11          0  17682084  17682083
       7      8    00    3059532  14622552  17682083
```

Here are some explanations:

{{< table "table-hover table-striped" >}}
| Field | Description |
|---------|----------|
| Dimensions | Describes the logical dimensions of the disk |
| Flags | Describes the flags listed in the partition table |
| Partition | The slice number described later in the partition table |
| Tag | Value used to specify how the disk will be used, described later in the partition table |
| Flags | Flag 00 is for "read/write, mountable"; 01 is "read/write, unmountable"; and 10 is "read only" |
| First Sector | Defines the first sector for the slice |
| Sector Count | Defines the total number of sectors in the slice |
| Last Sector | Defines the last sector in the slice |
| Mount Directory | If this field is empty, no entry will be defined in "/etc/vfstab" and the slice will not be mounted at startup |
{{< /table >}}

## In Case of Problems

### Relabeling Disks

The fmthard command allows you to relabel disks. First, let's save the current VTOC in a file:

```bash
prtvtoc /dev/dsk/c1t3d0s0 > /var/tmp/c1t3d0.vtoc
```

We can save the VTOC of another disk in a file to relabel it on a new disk:

```bash
fmthard -s datafile /dev/rdsk/c#t#d#s2
```

Open format, select the disk and give it the desired name. Then reinject the saved VTOC:

```bash
mthard -s /var/tmp/c1t3d0.vtoc /dev/rdsk/c1t3d0s2
```

Finally, to initialize the disk:

```bash
fmthard -s /dev/null /dev/rdsk/c1t3d0s2
```

## FileSystems

There are 4 types of FileSystems in Solaris:

* ufs: The most used FileSystem in Solaris. It can easily go up to Terabits, is based on the Berkeley system
* hsfs: A somewhat special Sierra system
* pcfs: For FAT32 and DOS
* udfs: Universal Disk File System, this is for CD/DVD media...

Here is a description of UFS:

![Sun ufs](/images/sun_ufs.avif)

As well as how inodes work:

![Sun ufs inodes](/images/sun_ufs_inodes.avif)

### Creating the FileSystem

The **newfs** command allows you to do this:

```bash
newfs /dev/rdsk/c1t3d0s7
```

Answer **y** to this confirmation:

```bash
newfs: construct a new file system /dev/rdsk/c1t3d0s7: (y/n)?
```

It now displays information about the filesystem creation:

```
/dev/rdsk/c1t3d0s7: 6295022 sectors in 1753 cylinders of 27 tracks, 
133 sectors 3073.7MB in 110 cyl groups (16 c/g, 28.05MB/g, 3392 i/g)
super-block backups (for fsck -F ufs -o b=#) at:
 32, 57632, 115232, 172832, 230432, 288032, 345632, 403232, 460832,
 518432, 5746208, 5803808, 5861408, 5919008, 5976608, 6034208, 6091808,
 6149408, 6207008, 6264608,
```

We display the free space

```bash
fstyp -v /dev/dsk/c0t1d0s6 |head
```

```
(output omitted for brevity)
minfree 10%     maxbpg  2048    optim   time
```

The -m option defines the percentage of disk space we want to use:

```bash
newfs -m 2 /dev/dsk/c0t1d0s6
```

```
newfs: construct a new file system /dev/rdsk/c0t1d0s6: (y/n)? y
(output omitted for brevity)
```

Here's the result:

```bash
fstyp -v /dev/dsk/c0t1d0s6 |head
```

```
(output omitted for brevity
minfree 2%      maxbpg  2048    optim   time
```

Verification:

```
fstyp -v /dev/rdsk/c0t0d0s0 | head
ufs
magic   11954   format  dynamic time    Fri Oct 22 10:09:11 2004
sblkno  16      cblkno  24      iblkno  32      dblkno  456
sbsize  5120    cgsize  5120    cgoffset 72     cgmask  0xffffffe0
ncg     110     size    3147511 blocks  3099093
bsize   8192    shift   13      mask    0xffffe000
fsize   1024    shift   10      mask    0xfffffc00
frag    8       shift   3       fsbtodb 1
minfree 2%      maxbpg  2048    optim   time
maxcontig 128   rotdelay 0ms    rps     120
```

To change the available space:

```bash
tunefs -m 1 /dev/rdsk/c0t0d0s0
```

```
minimum percentage of free space changes from 10% to 1%
```

### Checking Disk Status

The fsck command, like in Linux, allows you to check the integrity of the filesystem to repair any orphaned inodes (**don't forget to unmount the partition before performing this operation**):

```bash
fsck /dev/rdsk/c0t0d0s7
```

```
** /dev/rdsk/c0t0d0s7
** Last Mounted on /export/home
** Phase 1 - Check Blocks and Sizes
INCORRECT BLOCK COUNT I=743 (5 should be 2)
CORRECT?
```

If there are files that can be recovered, you will find them in "lost+found". Check if they are correct with the "file" command. It is estimated that if "file" can determine what type of file it is, then it's correct.

## Vfstab

The `/etc/vfstab` file is the equivalent of `/etc/fstab` in Linux. It lists all partitions and their mount points.

```
#device         device          mount           FS      fsck    mount   mount
#to mount       to fsck         point           type    pass    at boot options
#
fd      -       /dev/fd fd      -       no      -
/proc   -       /proc   proc    -       no      -
/dev/dsk/c0t0d0s1       -       -       swap    -       no      -
/dev/dsk/c0t0d0s0   /dev/rdsk/c0t0d0s0   /      ufs     1       no      -
/dev/dsk/c0t0d0s6   /dev/rdsk/c0t0d0s6   /usr   ufs     1       no      -
/dev/dsk/c0t0d0s3   /dev/rdsk/c0t0d0s3   /var   ufs     1       no      -
/dev/dsk/c0t0d0s7   /dev/rdsk/c0t0d0s7   /export/home  ufs   2  yes     -
/devices        -   /devices        devfs   -       no      -
ctfs    -       /system/contract        ctfs    -       no      -
objfs   -       /system/object  objfs   -       no      -
swap    -       /tmp    tmpfs   -       yes -
```

### mtab

The `/etc/mtab` file informs us about mountings in relation to the kernel:

```bash
more /etc/mnttab
```

```
/dev/dsk/c0t0d0s0    /    ufs     rw,intr,largefiles,logging,xattr,onerror=panic,dev=2200008  1098604644
/devices        /devices        devfs   dev=4a80000     1098604620
ctfs    /system/contract        ctfs    dev=4ac0001     1098604620
proc    /proc   proc    dev=4b00000     1098604620
mnttab  /etc/mnttab     mntfs   dev=4b40001     1098604620
swap    /etc/svc/volatile       tmpfs   xattr,dev=4b80001       1098604620
objfs   /system/object  objfs   dev=4bc0001     1098604620
/dev/dsk/c0t0d0s6       /usr    ufs     rw,intr,largefiles,logging,xattr,onerror=panic,dev=220000e  1098604645
fd      /dev/fd fd      rw,dev=4d40001  1098604645
/dev/dsk/c0t0d0s3       /var    ufs     rw,intr,largefiles,logging,xattr,onerror=panic,dev=220000b  1098604647
swap    /var/run        tmpfs   xattr,dev=4b80002       1098604647
swap    /tmp    tmpfs   xattr,dev=4b80003       1098604647
/dev/dsk/c0t0d0s7       /export/home    ufs     rw,intr,largefiles,logging,xattr,onerror=panic,dev=220000f   1098604661
-hosts  /net    autofs  nosuid,indirect,ignore,nobrowse,dev=4dc0001     1098604678
auto_home       /home   autofs  indirect,ignore,nobrowse,dev=4dc0002    1098604678
sys-01:vold(pid491)     /vol    nfs     ignore,noquota,dev=4e00001      1098604701
```

```bash
mount
```

```
/ on /dev/dsk/c0t0d0s0 read/write/setuid/devices/intr/largefiles/logging/xattr/onerror=panic/dev=2200008 on Sun Oct 24 08:57:24 2004
/devices on /devices read/write/setuid/devices/dev=4a80000 on Sun Oct 24 08:57:00 2004
/system/contract on ctfs read/write/setuid/devices/dev=4ac0001 on Sun Oct 24 08:57:00 2004
/proc on proc read/write/setuid/devices/dev=4b00000 on Sun Oct 24 08:57:00 2004
/etc/mnttab on mnttab read/write/setuid/devices/dev=4b40001 on Sun Oct 24 08:57:00 2004
/etc/svc/volatile on swap read/write/setuid/devices/xattr/dev=4b80001 on Sun Oct 24 08:57:00 2004
/system/object on objfs read/write/setuid/devices/dev=4bc0001 on Sun Oct 24 08:57:00 2004
/usr on /dev/dsk/c0t0d0s6 read/write/setuid/devices/intr/largefiles/logging/xattr/onerror=panic/dev=220000e on Sun Oct 24 08:57:25 2004
/dev/fd on fd read/write/setuid/devices/dev=4d40001 on Sun Oct 24 08:57:25 2004
/var on /dev/dsk/c0t0d0s3 read/write/setuid/devices/intr/largefiles/logging/xattr/onerror=panic/dev=220000b on Sun Oct 24 08:57:27 2004
/var/run on swap read/write/setuid/devices/xattr/dev=4b80002 on Sun Oct 24 08:57:27 2004
/tmp on swap read/write/setuid/devices/xattr/dev=4b80003 on Sun Oct 24 08:57:27 2004
/export/home on /dev/dsk/c0t0d0s7 read/write/setuid/devices/intr/largefiles/logging/xattr/onerror=panic/dev=220000f on Sun Oct 24 08:57:41 2004
```

## Mounting Partitions

To manually mount partitions, there is the mount command. Here are some examples:

Mount the filesystem as read-only:

```bash
mount -o ro /dev/dsk/c0t0d0s7 /export/home
```

Set sticky bits across the entire partition:

```bash
mount -o ro,nosuid /dev/dsk/c0t0d0s7 /export/home
```

Remove access dates for each file, which optimizes access times

```bash
mount -o noatime /dev/dsk/c0t0d0s7 /export/home
```

If this partition only contains small files, use this option:

```bash
mount -o nolargefiles /dev/dsk/c0t0d0s7 /export/home
```

To mount all the contents of your `/etc/vfstab` file, use this command:

```bash
mountall
```

To only mount what is local:

```bash
mountall -l
```

### Determining the Mount Type

To know which options to pass, here are some interesting files:

* `/etc/vfstab` for FS
* `/etc/default/fs` for a local filesystem
* `/etc/dfs/fstypes` for remote filesystems

To know the characteristics of a partition:

```bash
fstyp /dev/rdsk/c0t0d0s7
```

```
ufs
```

You can specify during partition mounting if it's hsfs or pcfs:

```bash
mount -F hsfs -o ro /dev/dsk/c0t6d0s0 /cdrom
```

To unmount a partition, do this:

```bash
umount mount_point
```

And to force, use the -f option:

```bash
umount -f mount_point
```

## What's Happening on My Partition?

fuser is what allows in Solaris to know what's happening on the partition. In Linux, it's lsof. To list all processes running on this partition:

```bash
fuser -cu mount_point
```

To kill all processes:

```bash
fuser -ck mount_point
```

Check that no processes are on the partition:

```bash
fuser -c mount_point
```

## Problems with My Root Partition

If you want to fsck the root partition, insert the Sun CD/DVD then type this:

```bash
ok boot cdrom -s
```

```
Boot device: /pci@1f,0/pci@1,1/ide@3/cdrom@2,0:f File and args -s
SunOS Release 5.10 Generic 64 bit 
Copyright 1983-2004 by Sun Microsystems, Inc. All rights reserved.
Booting to milestone "milestone/single-user:default"
Configuring /dev and /devices
Use is subject to license terms
Using RPC Bootparams for network configuration information.
Skipping interface hme0
-
INIT: SINGLE USER MODE
```

Run fsck on your root partition:

```bash
fsck /dev/rdsk/c0t0d0s0
```

If everything worked well, you should be able to mount everything:

```bash
mount /dev/dsk/c0t0d0s0 /a
```

Otherwise, you need to fine-tune `/etc/vfstab`:

```bash
TERM=sun
export TERM
vi /a/etc/vfstab
```

Then we exit and restart:

```bash
cd /
umount /a
```

## Access to Removable Devices

### With Vold

Where to find the peripherals:

{{< table "table-hover table-striped" >}}
| Media | Filesystem Access | Mounted Access |
|---------|----------|----------|
| diskette | /floppy/floppy0 | /vol/dev/aliases/floppy0 |
| CD-ROM | /cdrom/cdrom0 | /vol/dev/aliases/cdrom0 |
| Jaz | /rmdisk/jaz0 | /vol/dev/aliases/jaz0 |
| Zip | /rmdrive/zip0 | /vol/dev/aliases/zip0 |
| PCMCIA | /pcmem0 | /vol/dev/aliases/pcmem0 |
{{< /table >}}

There are 2 files that manage actions during media insertion/ejection:

* `/etc/vold.conf`
* `/etc/rmmount.conf`

Vold is a service (start, stop...):

```bash
/etc/init.d/volmgt restart
```

If it really doesn't want to quit:

```bash
pkill -9 vold
```

### Without Vold

Obviously you have to do everything manually:

```bash
mount -F hsfs -o ro /dev/dsk/c0t6d0s0 /cdrom
mount -F pcfs /dev/diskette /pcfs
```
