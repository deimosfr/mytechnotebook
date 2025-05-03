---
weight: 999
url: "/Encfs_\\:_Mise_en_place_d'Encfs_avec_FUSE/"
title: "EncFS: Setting up EncFS with FUSE"
description: "A guide on how to implement EncFS with FUSE for encrypted filesystems."
categories: ["Linux", "Security", "Debian"]
date: "2011-02-05T16:49:00+02:00"
lastmod: "2011-02-05T16:49:00+02:00"
tags: ["Linux", "Security", "Encryption", "FUSE", "Filesystem"]
toc: true
---

## 1. Introduction

[EncFS](https://encfs.sourceforge.net/) is an encrypted file system which allows you to store files on your hard drive that others cannot access, even if the machine is physically taken.

EncFS seems particularly interesting for the following reasons:

- No need to be root to use it
- Super simple implementation and trivial usage (see below)
- Separate encryption for each file. This may seem less secure (file sizes, names, and modification dates are known) but it has the enormous advantage of being more efficient (no need to re-encrypt the entire volume if only one file changes) and adaptable (no need to predict in advance the size the file system will take)
- (subjective opinion) the developer seems to have a good understanding of cryptographic tools. You should verify this for yourself ;-)

## 2. Installation

### 2.1 Debian

Let's get the necessary packages:

```bash
aptitude install encfs fuse-utils
```

Now, let's check that you have what you need in your kernel:

```bash
grep FUSE < /boot/config-2.6.19
```

You should have one of these 2 lines:

```
CONFIG_FUSE_FS=y
CONFIG_FUSE_FS=m
```

If this is not the case, you need to recompile your kernel and set FUSE as a module or integrate it directly.

### 2.2 FreeBSD

On FreeBSD, it's preferable to take the packaged versions as the compilation part is long and requires dependencies:

```bash
pkg_add -r fusefs-encfs
```

Then add this line to your rc.conf:

```bash
# /etc/rc.conf
...
fusefs_enable="YES"
```

Next, if you want any user to be able to use fuse, run this command:

```bash
sysctl vfs.usermount=1
```

## 3. Configuration

We will create the encrypted file system:

```bash
encfs ~/.crypt ~/crypt
```

Answer the questions:

```
The directory "/home/deimos/.crypt" does not exist. Should it be created? (y,n) y
The directory "/home/deimos/crypt" does not exist. Should it be created? (y,n) y
Creating a new encrypted volume.
Please choose one of the following options:
  enter "x" for expert configuration mode,
  enter "p" for pre-configured paranoia mode,
  anything else or an empty line will select standard mode.
?> p
```

We're not going to complicate things here and will configure in Paranoia mode to have the best encryption:-)

```
Paranoid configuration selected.

Configuration completed. The filesystem to be created has the following properties:
Filesystem cipher: "ssl/aes", version 2:1:1
Filename encoding: "nameio/block", version 3:0:1
Key size: 256 bits
Block size: 512 bytes, including 8 byte MAC header
Each file contains 8 bytes of header with unique IV data.
Filenames encrypted using IV chaining mode.
File data IV is chained to filename IV.

Now you will need to enter a password for your filesystem.
You will need to remember this password, as there is absolutely no
recovery mechanism. However, the password can be changed
later using encfsctl.
```

Choose your password:

```
New EncFS Password:  
Verify EncFS Password:
```

## 4. Usage

### 4.1 Single User

To mount our encrypted folder:

```bash
encfs ~/.crypt ~/crypt
```

Enter your password now. (Read the Multi-User version for the rest of the explanations). If you want everyone to be able to access your encrypted share, you need to add *--public*:

```bash
encfs --public ~/.crypt ~/crypt
```

### 4.2 Multi Users

Create a file to be encrypted on-the-fly:

```bash
touch ~/crypt/toto
```

Let's check our file:

```bash
> ls -l ~/crypt/toto
-rw-r--r--  1 deimos deimos 0 2005-04-20 21:27 /home/deimos/crypt/toto
```

And now in the encrypted folder:

```bash
> ls -l ~/.crypt/
 total 0
-rw-r--r--  1 deimos deimos 0 2005-04-20 21:27 FmIxHB3JurWr9jUCCgsUI8Ei
```

### 4.3 Unmounting a volume

To unmount the volume:

```bash
fusermount -u ~/crypt/
```

That's how simple it is :-)

### 4.4 Changing the password

If you want to change the password of an encfs volume, do this:

```bash
encfsctl passwd ~/.crypt
```

## 5. Pam and Encfs

Put pam_encfs.conf in /etc/security and modify your pam to load (for example):

```
auth required        pam_encfs.so
```

and if you want to auto umount on logout:

```
session        required        pam_encfs.so
```

(note that setting "encfs_default --idle=1", means it'll auto umount after 1 minute idletime, so you can ignore this if you want to)

If you want gdm working you'll have to do this: (to allow use of --public / allow_root / allow_other)

```
echo "user_allow_other" >> /etc/fuse.conf
```

```
adduser testuser # (put him in the fuse group if you have one)
mkdir -p /mnt/storage/enc/testuser 
```

Setup your /etc/pam_encfs.conf (default should work)

```
chown testuser:testuser /mnt/storage/enc/testuser
su testuser
encfs /mnt/storage/enc/testuser /home/testuser
```

Use same password as your login atm:

```
fusermount -u /home/testuser
```

when you login, the directory should be mounted.

Example to enable encryption for existing user (logout of any important things, turn off your apps, preferably do this in terminal login/as root):

```
sudo mkdir -p /mnt/storage/enc/anders /mnt/storage/enc/tmp
```

Use your main password on next part:

```
encfs /mnt/storage/enc/anders /mnt/storage/enc/tmp -- -o allow_root
cd /home/anders
find . -print -xdev | cpio -pamd /mnt/storage/enc/tmp
fusermount -u /mnt/storage/enc/tmp
cd /
sudo mv /home/anders /home/anders.BAK
sudo mkdir /home/anders
sudo chown anders:anders /home/anders
sudo rmdir /mnt/storage/enc/tmp
exit
```

On next login (in theory) your homedir should be mounted ;)

## 6. FAQ

### 6.1 Is there an example configuration file?

Yes, both in svn (link at [https://hollowtube.mine.nu/wiki/index.php?n=Projects.PamEncfs](https://hollowtube.mine.nu/wiki/index.php?n=Projects.PamEncfs)), and in the downloaded archive from my release. Some distributions have chosen an extremely simple example configuration file, mine is a bit more explained.

### 6.2 What command will pam_encfs run to mount a directory?

It depends on your options, but something like:

```
encfs -S --idle=1 -v /mnt/storage/enc/test /home/test -- -o allow_other,allow_root,nonempty
```

### 6.3 My KDE doesn't work

Login through KDE sometimes fails because KDE tries to store files to the home directory before mounting, and expect them to be there afterwards. To work around this you'll need to set 3 things in /etc/kde3/kdm/kdmrc, "DmrcDir=/tmp" (in the general section). And "UserAuthDir,ClientLogFile", both can be set to /tmp, these are in the [X-*-Core] section.
There might be security related issues with this solution, I haven't looked into that. If your paranoid about it you could make a temp directory tmp/user that only you have access to.

### 6.4 Can I mount multiple under one login directories with pam_encfs?

No, there is however an unofficial patch here: [https://bugs.gentoo.org/show_bug.cgi?id=102112](https://bugs.gentoo.org/show_bug.cgi?id=102112) ([https://joshua.haninge.kth.se/~sachankara/pam_encfs-0.1.3-multiple-mount-points.patch](https://joshua.haninge.kth.se/~sachankara/pam_encfs-0.1.3-multiple-mount-points.patch)). This has not been applied to the main tree, as it segfaults when I test it with a very basic encfs configuration file (but might work with more advanced ones).

### 6.5 pam_encfs does not find my encfs executable

pam_encfs uses execvp, that means that in some systems it wont find it if it's in /usr/local/bin, make a symlink to /usr/bin.

### 6.6 It works on normal login, but not in gdm

* Problem1, /etc/pam.d/gdm has a different system than /etc/pam.d/login, fix it ;).
* Problem2, You dont have the fuse option user_allow_root(or other) set,
  * Make sure /etc/fuse.conf has user_allow_other (or user_allow_root).
  * Make sure /etc/pam_encfs.conf has fuse_default allow_root, or the fuse option allow_root set.

### 6.7 It asks me for my password twice

Try adding use_first_pass after pam_unix (or any other module that supports it).

### 6.8 I've tried to use pam_encfs as my main authentication scheme, it doesn't work!

I return PAM_IGNORE on errors, this can't work reliably as a main system, because of for example logging in twice (in which case the directory would already be mounted, and we therefore can't check password ok).

### 6.9 I can't login to X because the filesystem doesn't support locks

This could be a problem if your not using drop_permission, use it. And if you REALLY want to mount as root, put:

```
export XAUTHORITY=/tmp/.Xauthority-$USER
export ICEAUTHORITY=/tmp/.ICEauthority-$USER
```

in your ~/.bashrc

My system-auth file on gentoo:

```
auth       required     pam_env.so
auth       sufficient   /lib/security/pam_encfs.so
auth       sufficient   /lib/security/pam_sha512.so pwdfile /etc/security/pam.sha
auth       sufficient   pam_unix.so likeauth nullok
auth       required     pam_deny.so

account    required     pam_unix.so

password   required     pam_cracklib.so retry=3
password   sufficient   pam_unix.so nullok md5 shadow use_authtok
password   required     pam_deny.so

session    required     pam_limits.so
session    required     pam_unix.so
```

Here it'll ask for the password twice, my modules (pam_encfs/pam_sha512) will try to use any previous password if it finds one.
So if you move pam_unix.so in auth to under pam_env.so, it'll ask for the password once.
Note that if pam_unix gets a password it finds ok, pam_encfs/pam_sha512 wont be used at all.

### 6.10 I can't create hard link

This is because of the External IV Chaining:

```
There is a cost associated with this. When External IV Chaining is enabled, hard links will not be allowed within the filesystem, as there would be no way to properly decode two different filenames pointing to the same data.

Also, renaming a file requires modifying the file header. So renames will only be allowed when the user has write access to the file.

Because of these limits, this option is disabled by default for standard mode (and enabled by default for paranoia mode).
```

If you create as the paranoia mode a crypted partition, this mode will be automatically enable. So use the default mode or create as expert mode without this :-)

### 6.11 Encfs in an OpenVZ VE

If you want to use encfs in a VE, you may encounter permission issues:

```
EncFS Password: 
fuse: device not found, try 'modprobe fuse' first
fuse failed. Common problems:
 - fuse kernel module not installed (modprobe fuse)
 - invalid options -- see usage message
```

Be aware that you need to load the fuse module at the VZ level so that VEs inherit it. Add this to your VZ to avoid having to load the module at each boot:

```bash
# /etc/modules
...
# Load Fuse
fuse
...
```

Then load it dynamically to access it afterwards:

```bash
modprobe fuse
```

To work around this, we'll create the fuse device from the host on the VE in question and add admin rights to it (somewhat problematic in terms of security, but no choice):

```bash
vzctl set $my_veid --devices c:10:229:rw --save
vzctl exec $my_veid mknod /dev/fuse c 10 229
vzctl set $my_veid --capability sys_admin:on --save
```

The second line may not work when the VE is turned off. Run it once it's on and then mount your encfs partition.
