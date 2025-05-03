---
weight: 999
url: "/LVM_\\:_Utilisation_des_LVM/"
title: "LVM: Working with Logical Volume Management"
icon: "description"
description: "How to create, manage, and optimize logical volumes on Linux systems using LVM for flexible storage management"
categories: ["Linux", "Storage", "System Administration"]
date: "2014-09-15T11:24:00+02:00"
lastmod: "2014-09-15T11:24:00+02:00"
tags:
  [
    "LVM",
    "Storage",
    "Partitioning",
    "Linux",
    "Volume Management",
    "Snapshot",
    "Filesystem",
  ]
toc: true
---

## Introduction

Logical Volume Management (LVM) is a method and software for partitioning, concatenating, and utilizing storage spaces on a server. It allows flexible management, security, and online optimization of storage spaces in UNIX/Linux-type operating systems.

We also refer to it as Volume Manager.

Since LVM is not very simple to use, and since I don't handle it every day either, I thought a small documentation was essential. I'll fill it in as needed.

Here are the main files and folders used by LVM:

{{< table "table-hover table-striped" >}}
| File/Folder | Description |
|-------------|-------------|
| `/etc/lvm/lvm.conf` | LVM configuration file |
| `/etc/lvm/cache.cache` | Device name cache file |
| `/etc/lvm/backup` | Folder containing automatic backups of VG metadata |
| `/etc/lvm/archive` | Folder containing automatic archives of VG metadata |
| `/var/lock/lvm` | Lock file to prevent simultaneous execution of multiple LVM tools, avoiding metadata corruption |
{{< /table >}}

## Creating a Partition

When you use fdisk to create a partition and assign it as LVM, it may not appear in the devices. To avoid rebooting to see it, simply run this command on the disk in question (sda for example):

```bash
partx -a /dev/sda
```

If you then need to add it to the fstab, it's preferable to use UUIDs. To locate them, there's a command:

```bash
$ blkid
blkid/dev/sda1: UUID="fd292b5c-091f-4a2f-b694-7d881e2eaa54" TYPE="ext4"
/dev/sda2: UUID="a4ljAA-KRHZ-sSHH-E4yy-tKmc-ZsHF-Scrc2N" TYPE="LVM2_member"
/dev/mapper/vg_redhatsrv1-lv_root: UUID="74706659-f76e-4dd3-8f93-b3629923d356" TYPE="ext4"
/dev/mapper/vg_redhatsrv1-lv_swap: UUID="c02b0c7b-88b3-4521-820f-a42c27871e35" TYPE="swap"
/dev/sdb1: UUID="c7b14f0b-1ccd-40f8-8183-7e85cbcb9d64" TYPE="crypto_LUKS"
```

All you have to do is add a line like this to the fstab (`/etc/fstab`):

```bash
UUID=fd292b5c-091f-4a2f-b694-7d881e2eaa54 /mnt                   ext4    defaults        0 0
```

## Usage

### Activating an LVM Partition

To activate an LVM partition, we'll use pvcreate:

```bash
pvcreate /dev/sda1
```

You can then see the status with the pvdisplay command.

### Creating a Volume Group

We must then create a VG if we want to be able to create volumes:

```bash
vgcreate my_vg /dev/sda1
```

You should put the name of the VG you want here instead of 'my_vg'.

### Creating Volumes

To create volumes:

```bash
lvcreate -n my_lv1 -L 5G my_vg
```

So you need:

- my_lv1: the desired LV name
- -L 5G: the volume group size
- my_vg: the name of the VG on which the LV should be stored.

### Resizing

#### Change block storage size (optional)

You may need to update the partition block size after updating the partition table. First get `growpart`:

```bash
apt-get install cloud-utils
```

Now get the partition table:

```bash
$ lsblk
NAME            MAJ:MIN RM  SIZE RO TYPE MOUNTPOINTS
sda               8:0    0  300G  0 disk
├─sda1            8:1    0  512M  0 part /boot/efi
├─sda2            8:2    0  488M  0 part /boot
└─sda3            8:3    0  299G  0 part
  ├─vg-root     254:0    0 74.5G  0 lvm  /
  └─vg-longhorn 254:1    0  200G  0 lvm  /mnt/longhorn
sr0              11:0    1  633M  0 rom
sr1              11:1    1 1024M  0 rom
```

Then run the following command to resize the partition (3 in this case):

```bash
sudo growpart /dev/sda 3
```

This command will resize the third partition on the `/dev/sda` disk. You can change the partition number and disk name as needed.

#### Increasing Size

To increase the size of a partition:

```bash
lvextend -L +10G /dev/mapper/my_partition
```

Then you just need to resize the filesystem (xfs_grow for example for XFS)

#### Reducing Size

Decreasing the size of a filesystem is a bit more delicate. Indeed, if you make the mistake of decreasing the size of the logical volume before reducing the content size (the filesystem itself), then you destroy your filesystem... same if you reduce the logical volume size too much.

To avoid any risk, I recommend using the following method (a bit longer than normal, but much more reliable):

- Reduce the filesystem size more than necessary
- Reduce the logical volume size to give it exactly the new desired size.
- Enlarge the filesystem so that it occupies all available space.

This way, the risk of error is much lower.

##### Reducing the FS

If you don't want to do this operation manually, use a live CD with gparted included.

Caution, not all filesystems can be "reduced". For ext3 and reiserfs, it works very well. Here's an example with reiserfs...

```bash
df -h
```

In this example, the "ca" volume is in the volume group svg. On this logical volume exists a reiserfs filesystem of 512 MB size. However, I only use 230 MB. Moreover, I know I'll never add anything to this volume. So I want to decrease its size to 256 MB (to leave a safety margin, and because it makes a round number ;) I start by unmounting the filesystem:

```bash
umount /home/ca
```

Then I will reduce the size of the filesystem, more than necessary. Rather than remove 256 MB, I'll remove 258. I can do this because there's 283 MB free... Obviously, removing more space than remains would be suicidal...

```bash
resize_reiserfs -s -258M /dev/svg/ca
```

CAUTION: If you're using ext3, you can't indicate the amount of space to remove, you must give the final desired size (512-258). The right command would have been:

```bash
resize2fs -p /dev/svg/ca 254M
```

It's possible that it asks you to run this command before. If that's the case, do it:

```bash
e2fsck -f /dev/svg/ca
```

##### Reducing the LV

In case you don't see the LV in /dev/mapper, you'll need to activate it:

```bash
lvchange -a y /dev/svg/ca
```

Note: to deactivate an LV, change the "y" to "n".

Now that the filesystem has decreased, we need to give the logical volume its new size, 256 MB instead of 512:

```bash
lvresize -L -256M /dev/svg/ca
  WARNING: Reducing active logical volume to 256.00 MB
  THIS MAY DESTROY YOUR DATA (filesystem etc.)
Do you really want to reduce ca? [y/n]: y
  Reducing logical volume ca to 256.00 MB
```

Just one last step, we tell the filesystem that it can automatically expand to take all available space. It should therefore be able to grow by 2 MB. It will find the exact size in blocks etc. on its own... We didn't take the risk of making an error by reducing it "exactly" to the same size as the logical volume, because the slightest error could have corrupted the filesystem by a few blocks.

```bash
resize_reiserfs /dev/svg/ca
```

or, if you're using ext3:

```bash
resize2fs /dev/svg/ca
```

All that's left is to remount the filesystem:

```bash
$ mount /dev/svg/ca /home/ca
$ df -h | grep ca
/dev/mapper/svg-ca    256M  230M   27M  90% /home/ca
```

It's done... The filesystem is now 256 MB, and we still have our 230 MB of data inside. Conclusion: Playing with the size of logical volumes works very well, you just have to take your time and not do anything silly :)

### Volume Recovery

If I'm in recovery mode, and I want to mount my filesystems, how do I do it?

Make sure the module is properly loaded:

```bash
modprobe dm-mod
```

Scan all LVMs:

```bash
vgscan
vgchange -ay
```

And mount the LV:

```bash
mkdir -p /mnt/VolGroup00/LogVol00
mount /dev/VolGroup00/LogVol00 /mnt/VolGroup00/LogVol00
```

Once finished, you can deactivate your volume:

```bash
vgchange -an VolGroup00
```

### Exporting/Importing a VG

- You can completely unmount a volume group like this:

```bash
vgexport <my_vg>
```

This will ensure that nothing is mounted and active on your machine.

- To remount it, it's quite simple:

```bash
vgimport <my_vg>
```

Then activate the vg to display the volumes:

```bash
vgchange -ay <my_vg>
```

### Scanning New LVM Volumes

You can verify that there's a new volume by scanning the PVs:

```bash
pvscan
```

The new volumes will then appear.

### Disk Replacement

If you have a VG with multiple disks in it. One of them is defective and you've added its replacement to the VG. There is a hot solution to remove the defective disk from the VG without data loss:

```bash
pvmove /dev/<my_defective_disk>
```

All that's left is to remove the disk.

### Controlling Volume Visibility

You can disable visibility of a volume like this:

```bash
lvchange -an my_lv
```

You can enable visibility of a volume this way:

```bash
lvchange -ay my_lv
```

You can force an lv to be read-only accessible:

```bash
lvchange -pr my_lv
```

Or writable:

```bash
lvchange -pw my_lv
```

### Snapshot

#### Creating a Snapshot

To create a snapshot of an existing LV:

```bash
lvcreate -s -n datasnap -L +4G vgdata/lvdata
```

Here I just created a 4G Snapshot.

#### Deleting the Snapshot

If I'm satisfied with the changes made, I want to delete the snapshot:

```bash
lvremove /dev/vgdata/datasnap
```

#### Rolling Back Changes

If I'm not satisfied with the changes and want to roll back from the snapshot, I'll need to unmount the partition and revert:

```bash
umount <mountpoint>
lvconvert --merge -v vgdata/datasnap
lvchange -an /dev/vgdata/data
lvchange -ay /dev/vgdata/data
mount /dev/vgdata/data <mountpoint>
```

The old partition will be remounted at the right place.

### Dumping the Configuration

You may need to dump the LVM configuration. To do this, run:

```bash
$ lvm dumpconfig
  devices {
  	dir="/dev"
  	scan="/dev"
  	obtain_device_list_from_udev=1
  	preferred_names=["^/dev/mpath/", "^/dev/mapper/mpath", "^/dev/[hs]d"]
  	filter="a/.*/"
  	cache_dir="/etc/lvm/cache"
  	cache_file_prefix=""
  	write_cache_state=1
  	sysfs_scan=1
  	md_component_detection=1
  	md_chunk_alignment=1
  	data_alignment_detection=1
  	data_alignment=0
  	data_alignment_offset_detection=1
  	ignore_suspended_devices=0
  	disable_after_error_count=0
  	require_restorefile_with_uuid=1
  	pv_min_size=2048
  	issue_discards=0
  }
  dmeventd {
  	mirror_library="libdevmapper-event-lvm2mirror.so"
...
```

### Viewing Locks

To view current locks on LVMs, use this command:

```bash
$ lvmdiskscan
  /dev/ram0  [      16,00 MiB]
  /dev/root  [       2,81 GiB]
  /dev/ram1  [      16,00 MiB]
  /dev/sda1  [     200,00 MiB]
  /dev/ram2  [      16,00 MiB]
  /dev/sda2  [    1000,00 MiB]
  /dev/ram3  [      16,00 MiB]
  /dev/sda3  [       2,83 GiB] LVM physical volume
  ...
  /dev/ram13 [      16,00 MiB]
  /dev/ram14 [      16,00 MiB]
  /dev/ram15 [      16,00 MiB]
  /dev/sdb   [       4,00 GiB]
  2 disks
  18 partitions
  0 LVM physical volume whole disks
  1 LVM physical volume
```

### Hot FSCK on ext3/4

You've dreamed of it, right? It is indeed possible to do hot fsck thanks to LVM snapshots. Here's a small script that allows scanning all ext3 and ext4 type LVs and defragmenting them (_Note: You need 10% free space compared to the largest LV in your VGs_):

```bash
#!/bin/bash
# Made by deimosfr
# Inspired from https://www.redhat.com/archives/ext3-users/2008-January/msg00032.html

EMAIL=xxx@mycompany.com

# Get line by line
IFS=$'\n'

for line in `lvs --unbuffered -o lv_name,vg_name,lv_size,vg_free --units m --noheadings --nosuffix`; do
	# Get informations on the device
	VOLUME=`echo $line | awk '{ print $1 }'`
	VG=`echo $line | awk '{ print $2 }'`
	LV_SIZE=`echo $line | awk '{ print $3 }' | sed 's/\..*//g'`
	VG_FREE_SIZE=`echo $line | awk '{ print $4 }' | sed 's/\..*//g'`
	SNAPSIZE=`echo $LV_SIZE | awk '{ printf("%s\n"), int($1*0.1) }'`

	# Check if 10% of the LV size is free on the VG
	if [ $SNAPSIZE -lt $VG_FREE_SIZE ]; then
		# Check if it is ext3 or ext4
		LV_DEV=`echo ${VOLUME} | sed 's/\-/--/'`
		if [ `mount | grep /dev/mapper/${VG}-${LV_DEV} | grep -c "ext[3-4]"` -ge 1 ]; then
		        echo "Analyzing ${VG}-${VOLUME}"
		        TMPFILE=`mktemp -t e2fsck_${VG}_${VOLUME}.log.XXXXXXXXXX`
		        START=$(date +'%Y%m%d%H%M%S')

		        # Create snapshot
		        lvcreate -s -L ${SNAPSIZE} -n "${VOLUME}-snap" "${VG}/${VOLUME}"
		        if nice logsave -as $TMPFILE e2fsck -p -C 0 "/dev/${VG}/${VOLUME}-snap" && nice logsave -as $TMPFILE e2fsck -fy -C 0 "/dev/${VG}/${VOLUME}-snap"; then
		                # Set to 0 the fsck counter
		                tune2fs -C 0 -T "${START}" "/dev/${VG}/${VOLUME}"
		        else
		                # Set the fsck counter at max to fsck next time and send an email
		                tune2fs -C 16000 -T "19000101" "/dev/${VG}/${VOLUME}"
		                if test -n "EMAIL"; then
		                        mail -s "E2fsck of /dev/${VG}/${VOLUME} failed on host `hostname`!" $EMAIL < $TMPFILE
		                fi
		        fi

		        # Merge / Remove snapshot
		        lvremove -f "${VG}/${VOLUME}-snap"
		        rm $TMPFILE
		fi
	fi
done
```

## Having a Nice GUI

Well, for LVM and Red Hat connoisseurs, you certainly know the `system-config-lvm` command. For Debian it would be cool if we could have it! And well, bibi got the solution:

```bash
sudo apt-get install python-gtk2 python-gtk-1.2 python-glade-1.2 python-glade2 python-gnome2
wget http://download.fedora.redhat.com/pub/fedora/linux/core/updates/5/i386/system-config-lvm-1.0.18-1.2.FC5.noarch.rpm
sudo alien system-config-lvm-1.0.18-1.2.FC5.noarch.rpm
sudo dpkg -i system-config-lvm_1.0.18-2.2_all.deb
sudo ln -s /usr/bin/python2.4 /usr/bin/python2
```

All that's left is to launch and say "woooow":

```bash
system-config-lvm
```

## FAQ

### I Have a UUID Conflict

If you have the misfortune of running into a UUID conflict, you can regenerate one via the following command:

```bash
uuidgen
```

## Resources
- [Download the LVM documentation](/pdf/lvm.pdf)
- [Documentation for beginners](/pdf/lvm_beginner.pdf)
- [Documentation on LVM Snapshots](/pdf/lvm_snapshots.pdf)
