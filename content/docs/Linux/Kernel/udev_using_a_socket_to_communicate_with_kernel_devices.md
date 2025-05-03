---
weight: 999
url: "/Udev_\\:_Utilisation_d'un_socket_pour_parler_avec_les_devices_kernel/"
title: "Udev: Using a Socket to Communicate with Kernel Devices"
description: "A comprehensive guide on how to use udev for device management on Linux systems, including monitoring, creating custom rules, and handling device detection events"
categories: ["Debian", "Storage", "Linux"]
date: "2012-07-05T09:33:00+02:00"
lastmod: "2012-07-05T09:33:00+02:00"
tags:
  ["udev", "kernel", "devices", "linux", "system administration", "hardware"]
toc: true
---

## Introduction

udev is a device manager that replaces devfs on Linux kernels series 2.6. Its main function is to manage devices in the `/dev` directory.

udev runs in user mode and communicates with hotplug which runs in kernel mode. It uses and stores information it has discovered in `/sys`. When hardware is detected, udev can assign a device name, create symbolic links, or execute a program when an action occurs on one or more devices.

udev has an open socket on the machine, and the kernel informs udev of new devices through this socket.

Here's how udev works:

1. The kernel discovers a device and sends its status to sysfs
2. udev is informed of this new event via a netlink socket
3. udev creates the device (`/dev/device`) and/or launches a program (defined in udev rules)
4. udev informs hald (Hardware Abstraction Layer Daemon) of this event via a socket
5. The HAL (Hardware Abstraction) retrieves information about this device
6. The HAL builds object structures related to the device with information retrieved previously and via other resources
7. The HAL broadcasts events through D-Bus
8. A userland application watches for this type of event to process the information afterward

## Usage

There isn't much to do with udev for normal operation since it's generally correctly configured by your Linux distribution. However, you may want to play with it and customize your system a bit.

### Monitoring

You can monitor udev by using the 'udevmonitor' command. For example, if I want to monitor when I connect my iPhone:

```bash {linenos=table,hl_lines=[1]}
> udevmonitor --env
udevmonitor prints the received event from the kernel [UEVENT]
and the event which udev sends out after rule processing [UDEV]

UEVENT[1330349745.773910] add@/devices/pci0000:00/0000:00:1a.7/usb1/1-6
ACTION=add
DEVPATH=/devices/pci0000:00/0000:00:1a.7/usb1/1-6
SUBSYSTEM=usb
SEQNUM=1077
PHYSDEVBUS=usb
PHYSDEVDRIVER=usb
...
DEVNAME=/dev/usbdev1.7_ep83

UDEV  [1330349746.567644] add@/class/usb_device/usbdev1.7
UDEV_LOG=3
ACTION=add
DEVPATH=/class/usb_device/usbdev1.7
SUBSYSTEM=usb_device
SEQNUM=1083
PHYSDEVPATH=/devices/pci0000:00/0000:00:1a.7/usb1/1-6
PHYSDEVBUS=usb
PHYSDEVDRIVER=usb
MAJOR=189
MINOR=6
UDEVD_EVENT=1
DEVNAME=/dev/bus/usb/001/007
```

## The Rules

The udev rules are defined in `/etc/udev/rules.d`, which means you can customize them or create your own rules.

```bash {linenos=table,hl_lines=[4]}
# The initial syslog(3) priority: "err", "info", "debug" or its
# numerical equivalent. For runtime debugging, the daemons internal
# state can be changed with: "udevadm control --log-priority=<value>".
udev_log="err"
```

Here are the udev.conf configuration file options:

- udev_root: where devices should be created (`/dev`)
- udev_rules: where the rules are (`/etc/udev/rules.d/*.rules`)
- udev_log: defines the verbosity level

### Creating Your Rules

We'll create a file with the ".rules" extension to be taken into account by udev (`/etc/udev/rules.d/*.rules`), in the form:

```
<match-key><operator>value, <assignment-key><operator>value, <action><operator>value
```

which can correspond to something like:

```
PROGRAM=="script_to_run" RESULT=="what_should_be_returned_by_PROGRAM_for_validation"
```

Here's an example:

```bash
BUS=="usb", SYSFS{product}="iPhone", SYMLINK+="iphone"
```

I ask udev that when it detects an iPhone on a USB port, it automatically creates a symbolic link in `/dev` called `/dev/iphone`.

To find information on which we can build our udev rules, we must first retrieve information about the hardware. First, let's get the information (DEVNAME) on the path of the device we know [using monitoring](#monitoring):

```bash
> udevinfo -q path -n /dev/bus/usb/001/007
/class/usb_device/usbdev1.7
```

This command gave us the path in `/sys` and we'll use it now to get all available information:

```bash {linenos=table,hl_lines=[20]}
> udevinfo -a -p /class/usb_device/usbdev1.7

Udevinfo starts with the device specified by the devpath and then
walks up the chain of parent devices. It prints for every device
found, all possible attributes in the udev rules key format.
A rule to match, can be composed by the attributes of the device
and the attributes from one single parent device.

  looking at device '/class/usb_device/usbdev1.7':
    KERNEL=="usbdev1.7"
    SUBSYSTEM=="usb_device"
    SYSFS{dev}=="189:6"

  looking at parent device '/devices/pci0000:00/0000:00:1a.7/usb1/1-6':
    ID=="1-6"
    BUS=="usb"
    DRIVER=="usb"
    SYSFS{configuration}=="PTP"
    SYSFS{serial}=="8603c60e01995bb89ca2cb39570f9bb1039454df"
    SYSFS{product}=="iPhone"
    SYSFS{manufacturer}=="Apple Inc."
    SYSFS{maxchild}=="0"
    SYSFS{version}==" 2.00"
    SYSFS{devnum}=="7"
    SYSFS{speed}=="480"
    SYSFS{bMaxPacketSize0}=="64"
    SYSFS{bNumConfigurations}=="4"
...
```

I've just obtained the information that interested me about my device.

We'll tell udev to reload its rules so it takes our new rule into account (without having to reboot the machine):

```bash
udevcontrol reload_rules
```

Or depending on the OS version, this also works:

```bash
udevadm control --reload-rules
```

### With iSCSI

When using iSCSI, you can retrieve information like this:

```bash
> scsi_id -g -x -s /block/sda
ID_VENDOR=ATA
ID_MODEL=WDC_WD1600AAJS-6
ID_REVISION=58.0
ID_SERIAL=SATA_WDC_WD1600AAJS-_WD-WMAS20873789
ID_TYPE=disk
ID_BUS=scsi
```

The equivalents for USB or ATA are:

```bash
/lib/udev/ata_id /dev/hdx
/lib/udev/usb_id /dev/sdx
```

### Explanations on Creating Rules

Now, you might say that's fine, but you haven't explained much about how to construct these rules! Let's detail this (largely thanks to the man page).

#### Matching Keys

There are operators for defining udev matching rules, and these are designed to match device properties (some of these properties match the parent device in sysfs):

{{< table "table-hover table-striped" >}}
| Description | Operator |
|-------------|----------|
| Compare for equality | == |
| Compare for inequality | != |
| Assign a value to a key. Keys that represent a list, are reset and only this single value is assigned | = |
| Add the value to a key that holds a list of entries | += |
| Assign a value to a key finally; disallow any later changes, which may be used to prevent changes by any later rules | := |
{{< /table >}}

And here are the elements that can be used for matching:

{{< table "table-hover table-striped" >}}
| Description | Operator |
|-------------|----------|
| Match the name of the event action | ACTION |
| Match the name of the device | KERNEL |
| Match the devpath of the device | DEVPATH |
| Match the subsystem of the device | SUBSYSTEM |
| Match the name of the event action | ACTION |
| Search the devpath upwards for a matching device subsystem name | BUS |
| Search the devpath upwards for a matching device driver name | DRIVER |
| Search the devpath upwards for a matching device name | ID |
| Search the devpath upwards for a device with matching sysfs attribute values. Up to five SYSFS keys can be specified per rule. All attributes must match on the same device. Trailing whitespace in the attribute values is ignored, if the specified match value does not contain trailing whitespace itself. | SYSFS{filename} |
| Match against the value of an environment variable. Up to five ENV keys can be specified per rule. This key can also be used to export a variable to the environment. | ENV{key} |
| Execute external program. The key is true, if the program returns without exit code zero. The whole event environment is available to the executed program. The program's output printed to stdout is available for the RESULT key | PROGRAM |
| Match the returned string of the last PROGRAM call. This key can be used in the same or in any later rule after a PROGRAM call | RESULT |
{{< /table >}}

You can use matching patterns for your correspondences:

{{< table "table-hover table-striped" >}}
| Description | Operator |
|-------------|----------|
| Matches zero, or any number of characters | \* |
| Matches any single character | ? |
| Matches any single character specified within the brackets. For example, the pattern string 'tty[SR]' would match either 'ttyS' or 'ttyR'. Ranges are also supported within this match with the '-' character. For example, to match on the range of all digits, the pattern [0-9] would be used. If the first character following the '[' is a '!', any characters not enclosed are matched | [] |
{{< /table >}}

#### Rules on Key Assignment

Here are the solutions for values:

{{< table "table-hover table-striped" >}}
| Description | Operator |
|-------------|----------|
| The name, a network interface should be renamed to. Or as a temporary workaround, the name a device node should be named. Usually the kernel provides the defined node name, or even creates and removes the node before udev even receives any event. Changing the node name from the kernel's default creates inconsistencies and is not supported. If the kernel and NAME specify different names, an error will be logged. Udev is only expected to handle device node permissions and to create additional symlinks, not to change kernel-provided device node names. Instead of renaming a device node, SYMLINK should be used. Symlink names must never conflict with device node names, it will result in unpredictable behavior | NAME |
| The name of a symlink targeting the node. Every matching rule will add this value to the list of symlinks to be created. Multiple symlinks may be specified by separating the names by the space character. In case multiple devices claim the same name, the link will always point to the device with the highest link_priority. If the current device goes away, the links will be re-evaluated and the device with the next highest link_priority will own the link. If no link_priority is specified, the order of the devices, and which one of them will own the link, is undefined. Claiming the same name for a symlink, which is or might be used for a device node, may result in unexpected behavior and is not supported | SYMLINK |
| The permissions for the device node. Every specified value overwrites the compiled-in default value | OWNER, GROUP, MODE |
| The value that should be written to a sysfs attribute of the event device | ATTR{key} |
| Set a device property value. Property names with a leading '.' are not stored in the database or exported to external tool or events | ENV{key} |
| Attach a tag to a device. This is used to filter events for users of libudev's monitor functionality, or to enumerate a group of tagged devices. The implementation can only work efficiently if only a few tags are attached to a device. It is only meant to be used in contexts with specific device filter requirements, and not as a general-purpose flag. Excessive use might result in inefficient event handling | TAG |
| Add a program to the list of programs to be executed for a specific device. This can only be used for very short running tasks. Running an event process for a long period of time may block all further events for this or a dependent device. Long running tasks need to be immediately detached from the event process itself. If the option RUN{fail_event_on_error} is specified, and the executed program returns non-zero, the event will be marked as failed for a possible later handling. If no absolute path is given, the program is expected to live in /lib/udev, otherwise the absolute path must be specified. Program name and arguments are separated by spaces. Single quotes can be used to specify arguments with spaces | RUN |
| Named label where a GOTO can jump to | LABEL |
| Jumps to the next LABEL with a matching name | GOTO |
| Import a set of variables as device properties, depending on type: program: Execute an external program specified as the assigned value and import its output, which must be in environment key format. Path specification, command/argument separation, and quoting work like in RUN. file: Import a text file specified as the assigned value, which must be in environment key format. db: Import a single property specified as the assigned value from the current device database. This works only if the database is already populated by an earlier event. cmdline: Import a single property from the kernel commandline. For simple flags the value of the property will be set to '1'. parent: Import the stored keys from the parent device by reading the database entry of the parent device. The value assigned to IMPORT{parent} is used as a filter of key names to import (with the same shell-style pattern matching used for comparisons). If no option is given, udev will choose between program and file based on the executable bit of the file permissions | IMPORT{type} |
| Wait for a file to become available or until a 10 seconds timeout expires. The path is relative to the sysfs device, i.e. if no path is specified this waits for an attribute to appear | WAIT_FOR |
| Rule and device options: link_priority=value: Specify the priority of the created symlinks. Devices with higher priorities overwrite existing symlinks of other devices. The default is 0. event_timeout=Number of seconds an event will wait for operations to finish, before it will terminate itself. string_escape=none|replace: Usually control and other possibly unsafe characters are replaced in strings used for device naming. The mode of replacement can be specified with this option. static_node=: Apply the permissions specified in this rule to a static device node with the specified name. Static device nodes might be provided by kernel modules, or copied from /lib/udev/devices. These nodes might not have a corresponding kernel device at the time udevd is started, and allow to trigger automatic kernel module on-demand loading. watch: Watch the device node with inotify, when closed after being opened for writing, a change uevent will be synthesised. nowatch: Disable the watching of a device node with inotify. | OPTIONS |
{{< /table >}}

#### Substitution Rules

The NAME, SYMLINK, PROGRAM, OWNER, GROUP, MODE and RUN fields support substitution rules. These are built-in to help you set certain variables:

{{< table "table-hover table-striped" >}}
| Description | Operator |
|-------------|----------|
| The kernel name for this device | $kernel, %k |
| The kernel number for this device. For example, 'sda3' has kernel number of '3' | $number, %n |
| The devpath of the device | $devpath, %p |
| The name of the device matched while searching the devpath upwards for SUBSYSTEMS, KERNELS, DRIVERS and ATTRS. | $id, %b |
| The driver name of the device matched while searching the devpath upwards for SUBSYSTEMS, KERNELS, DRIVERS and ATTRS | $driver |
| The value of a sysfs attribute found at the device, where all keys of the rule have matched. If the matching device does not have such an attribute, and a previous KERNELS, SUBSYSTEMS, DRIVERS, or ATTRS test selected a parent device, use the attribute from that parent device. If the attribute is a symlink, the last element of the symlink target is returned as the value | $attr{file}, %s{file} |
| A device property value | $env{key}, %E{key} |
| The kernel major number for the device | $major, %M |
| The kernel minor number for the device | $minor, %m |
| The string returned by the external program requested with PROGRAM. A single part of the string, separated by a space character may be selected by specifying the part number as an attribute: %c{N}. If the number is followed by the '+' char this part plus all remaining parts of the result string are substituted: %c{N+} | $result, %c |
| The node name of the parent device | $parent, %P |
| The current name of the device node. If not changed by a rule, it is the name of the kernel device | $name |
| The current list of symlinks, separated by a space character. The value is only set if an earlier rule assigned a value, or during a remove events | $links |
| The udev_root value | $root, %r |
| The sysfs mount point | $sys, %S |
| The name of a created temporary device node to provide access to the device from a external program before the real node is created | $tempnode, %N |
| The '%' character itself | %% |
| The '$' character itself | $$ |
{{< /table >}}

## Examples

Here's a small list of examples to help you build rules:

This always maps a specific USB device (in this case, a pendrive) to `/dev/usbpen`, which is then set in fstab to mount on `/mnt/usbpen`:

```bash
# Symlink USB pen
SUBSYSTEMS=="usb", ATTRS{serial}=="1730C13B18000B84", KERNEL=="sd?", NAME="%k", SYMLINK+="usbpen", GROUP="storage"
SUBSYSTEMS=="usb", ATTRS{serial}=="1730C13B18000B84", KERNEL=="sd?1", NAME="%k", SYMLINK+="usbpen", GROUP="storage"
```

For devices with multiple partitions, the following example maps the device to `/dev/usbdisk`, and partitions 1, 2, 3 etc., to `/dev/usbdisk1`, `/dev/usbdisk2`, `/dev/usbdisk3`, etc.

```bash
# Symlink multi-part device
SUSSYSTEMS=="usb", ATTRS{serial}=="1730C13B18000B84", KERNEL=="sd?", NAME="%k", SYMLINK+="usbdisk", GROUP="storage"
SUBSYSTEMS=="usb", ATTRS{serial}=="1730C13B18000B84", KERNEL=="sd?[1-9]", NAME="%k", SYMLINK+="usbdisk%n", GROUP="storage"
```

The above rules are equivalent to the following one:

```bash
# Symlink multi-part device
SUBSYSTEMS=="usb", ATTRS{serial}=="1730C13B18000B84", KERNEL=="sd*", NAME="%k", SYMLINK+="usbdisk%n", GROUP="storage"
```

It's also possible to omit the NAME and GROUP statements, so that the defaults from udev.rules are used. The shortest and simplest solution would be adding this rule:

```bash
# Symlink multi-part device
SUBSYSTEMS=="usb", ATTRS{serial}=="1730C13B18000B84", KERNEL=="sd*", SYMLINK+="usbdisk%n"
```

This always maps an Olympus digicam to `/dev/usbcam`, which can be stated in fstab to mount on `/mnt/usbcam`:

```bash
# Symlink USB camera
SUBSYSTEMS=="usb", ATTRS{serial}=="000207532049", KERNEL=="sd?", NAME="%k", SYMLINK+="usbcam", GROUP="storage"
SUBSYSTEMS=="usb", ATTRS{serial}=="000207532049", KERNEL=="sd?1", NAME="%k", SYMLINK+="usbcam", GROUP="storage"
```

And this maps a Packard Bell MP3 player to `/dev/mp3player`:

```bash
# Symlink MP3 player
SUBSYSTEMS=="usb", ATTRS{serial}=="0002F5CF72C9C691", KERNEL=="sd?", NAME="%k", SYMLINK+="mp3player", GROUP="storage"
SUBSYSTEMS=="usb", ATTRS{serial}=="0002F5CF72C9C691", KERNEL=="sd?1", NAME="%k", SYMLINK+="mp3player", GROUP="storage"
```

To map a selected USB key to `/dev/mykey` and all other keys to `/dev/otherkey`:

```bash
# Symlink USB keys
SUBSYSTEMS=="usb", ATTRS{serial}=="insert serial key", KERNEL=="sd?1", NAME="%k", SYMLINK+="mykey"
SUBSYSTEMS=="usb", KERNEL=="sd?1", NAME="%k", SYMLINK+="otherkey"
```

Note the order of the lines. Since all USB keys should create the `/dev/sd<a||b>` node, udev will first check if it is a rules-stated USB key, defined by serial number. But if an unknown USB key is plugged, it will also create a node, using the previously stated generic name, "otherkey". That rule should be the last one in the rules file so that it does not override the others.

This is an example on how to distinguish USB HDD drives and USB sticks:

```bash
BUS=="usb", ATTRS{product}=="USB2.0 Storage Device", KERNEL=="sd?", NAME="%k", SYMLINK+="usbdisk", GROUP="storage"
BUS=="usb", ATTRS{product}=="USB2.0 Storage Device", KERNEL=="sd?[1-9]", NAME="%k", SYMLINK+="usbdisk%n", GROUP="storage"
BUS=="usb", ATTRS{product}=="USB Mass Storage Device", KERNEL=="sd?1", NAME="%k", SYMLINK+="usbflash", GROUP="storage"
```

## Resources
- http://reactivated.net/writing_udev_rules.html
- http://www.redhat.com/magazine/002dec04/features/udev/
- http://wiki.debian.org/udev
- https://wiki.archlinux.org/index.php/Map_Custom_Device_Entries_with_udev
