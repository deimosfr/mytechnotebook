---
weight: 999
url: "/Automatiser_une_installation_de_Debian/"
title: "Automate Debian Installation"
description: "Learn how to automate Debian installation using preseed files to create identical server setups efficiently."
categories: ["Linux", "Debian", "Installation", "Automation"]
date: "2013-05-07T11:09:00+02:00"
lastmod: "2013-05-07T11:09:00+02:00"
tags: ["Debian", "preseed", "Installation", "Automation", "DHCP"]
toc: true
---

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Operating System** | 7.0 |
| **Website** | [Debian Website](https://www.debian.org) |
| **Last Update** | 07/05/2013 |
{{< /table >}}

## Introduction

It's not always easy to set up 10 identical servers. That's why this section will help you have a clean and controlled installation.

## preseed.cfg

You must first create the preconfiguration file and place it where you want. Here's the preseed I use for Debian Wheezy (check [my Git](https://git.deimos.fr) for the most recent version):

```
# Preseed file for Debian
# Made by Pierre Mavro / Deimosfr

# To create a temporary web server to quickly serve this preseed file,
# simply type one of this command in the same folder than preseed:
# while true; do nc -l -p 8000 -q 1 < preseed.cfg ; done
# python -m SimpleHTTPServer

# For more informations:
# http://wiki.deimos.fr/Automatiser_une_installation_de_Debian

### Contents of the preconfiguration file (for wheezy)
d-i debian-installer/language string en
d-i debian-installer/country string FR
d-i debian-installer/locale string en_US.UTF-8

### Keyboard
d-i console-keymaps-at/keymap select fr
d-i keyboard-configuration/xkb-keymap select fr
d-i console-keymaps-at/keymap select fr
d-i keymap select fr(latin9)

### Network configuration
d-i netcfg/choose_interface select auto
d-i netcfg/get_hostname string unassigned-hostname
d-i netcfg/get_domain string unassigned-domain
d-i netcfg/wireless_wep string

### Apt mirror
d-i mirror/protocol string http
d-i mirror/country string manual
d-i mirror/http/hostname string ftp.fr.debian.org
d-i mirror/http/directory string /debian
d-i mirror/http/proxy string
d-i mirror/suite string wheezy

### Account setup
d-i passwd/root-login boolean false
d-i passwd/make-user boolean true
d-i passwd/root-password password soleil
d-i passwd/root-password-again password soleil
d-i passwd/user-fullname string Deimos
d-i passwd/username string deimos
d-i passwd/user-password password soleil
d-i passwd/user-password-again password soleil

### Clock and time zone setup
d-i clock-setup/utc boolean true
d-i time/zone string Europe/Paris
d-i clock-setup/ntp boolean true

### Partitioning
d-i partman-auto/method string lvm
d-i partman-lvm/device_remove_lvm boolean true
d-i partman-md/device_remove_md boolean true
d-i partman-lvm/confirm boolean true
d-i partman-lvm/confirm_nooverwrite boolean true
d-i partman-auto-lvm/new_vg_name string vgos
# Partition will be:
# /boot: ~128M ext4
# /: [1-∞]G LVM ext4
# /var: [768-2048]M LVM ext4
# swap: [RAM*150%-2048]M LVM
d-i partman-auto/expert_recipe string                         \
      boot-root::                                            \
              128 3000 128 ext4                               \
                      $primary{ }                             \
                      $bootable{ }                            \
                      method{ format } format{ }              \
                      use_filesystem{ } filesystem{ ext4 }    \
                      mountpoint{ /boot }                     \
                      options/noatime{ noatime }              \
              .                                               \
              1024 4000 -1 ext4                               \
                      $lvmok{ }                               \
                      method{ format } format{ }              \
                      use_filesystem{ } filesystem{ ext4 }    \
                      mountpoint{ / }                         \
                      options/noatime{ noatime }              \
                      lv_name{ root }                         \
              .                                               \
              768 1000 2048 ext4                              \
                      $lvmok{ }                               \
                      method{ format } format{ }              \
                      use_filesystem{ } filesystem{ ext4 }    \
                      mountpoint{ /var }                      \
                      options/noatime{ noatime }              \
                      lv_name{ var }                          \
              .                                               \
              100% 1000 150% linux-swap                       \
                      $lvmok{ }                               \
                      method{ swap } format{ }                \
                      lv_name{ swap }                         \
              .
d-i partman-partitioning/confirm_write_new_label boolean true
d-i partman/choose_partition select finish
d-i partman/confirm boolean true
d-i partman/confirm_nooverwrite boolean true
d-i partman-md/confirm boolean true
d-i partman/mount_style select uuid

### Base system installation
d-i base-installer/install-recommends boolean false

### Apt setup
apt-cdrom-setup apt-setup/cdrom/set-first boolean false
d-i apt-setup/non-free boolean true
d-i apt-setup/contrib boolean true
d-i apt-setup/use_mirror boolean true
d-i apt-setup/services-select multiselect security, volatile
d-i apt-setup/security_host string security.debian.org
d-i apt-setup/volatile_host string volatile.debian.org

### Package selection
tasksel tasksel/first multiselect standard
d-i pkgsel/upgrade select safe-upgrade
popularity-contest popularity-contest/participate boolean true
d-i pkgsel/include string openssh-server

### Grub
d-i grub-installer/only_debian boolean true
d-i grub-installer/with_other_os boolean true

# Finish install
d-i finish-install/reboot_in_progress note
d-i cdrom-detect/eject boolean true
```

I've highlighted a little trick to quickly set up a temporary web server to deliver the preseed without having to remake an ISO image. I've also highlighted the login and password part (soleil) which should be changed :-)

## Loading the Preconfiguration File

For loading this file, you can choose what you want (file, http...) during the installation boot (grub):

```bash
# Web server version
preseed/url=http://host/path/to/preseed.cfg
# CD version
preseed/file=/cdrom/preseed.cfg
# USB key version
preseed/file=/hd-media/preseed.cfg
```

You can also edit the txt.cfg file on a CD-ROM to tell it where the file is located:

```
# isolinux/txt.cfg
default install
label install
    menu label ^Install
    menu default
    kernel /install.amd/vmlinuz
    append preseed/file=/cdrom/preseed.cfg auto=true priority=critical lang=fr locale=en_US.UTF-8 console-keymaps-at/keymap=fr-latin9 vga=788 initrd=/install.amd/initrd.gz -- quiet
```

## Using a DHCP Server to Specify Preconfiguration Files

It's also possible to use DHCP to specify a file to download from the network. DHCP allows you to specify a filename. Normally this file is used for network booting. If it's a URL, the installation system that allows network-type preconfiguration will download the file and use it as a preconfiguration file. Here's an example showing how to configure the dhcpd.conf file belonging to version 3 of the ISC DHCP server (debian package dhcp3-server).

```
# /etc/dhcp/dhcpd3.cfg
if substring (option vendor-class-identifier, 0, 3) = "d-i" {
    filename "http://host/preseed.cfg";
}
```

Note that the above example only allows the file for DHCP clients that identify themselves as "d-i". Other DHCP clients are not affected. You can also put the text in a paragraph targeted at a single host to avoid preconfiguring all installations done in your network.

A good way to use this technique is to only preconfigure values related to your network, for example the name of your Debian mirror. This way installations automatically use the right mirror and the rest of the installation can be done interactively. You need to be very careful if you want to automate the entire installation with DHCP-type preconfiguration.

## Resources
- [https://www.debian.org/releases/stable/s390/apbs02.html.fr](https://www.debian.org/releases/stable/s390/apbs02.html.fr)
- [https://www.unixgarden.com/index.php/gnu-linux-magazine-hs/une-installation-de-debian-automatique-2](https://www.unixgarden.com/index.php/gnu-linux-magazine-hs/une-installation-de-debian-automatique-2)
