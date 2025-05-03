---
weight: 999
url: "/Yubikey_\\:_Configure_your_yubikey_with_pam/"
title: "Yubikey: Configure Your Yubikey with PAM"
description: "Guide on setting up Yubikey with PAM for two-factor authentication on Linux systems, including configuration for challenge-response mode and automatic screen locking when the key is removed."
categories: ["Security", "Linux", "Debian"]
date: "2015-04-22T16:04:00+02:00"
lastmod: "2015-04-22T16:04:00+02:00"
tags:
  ["Authentication", "Security", "Yubikey", "PAM", "2FA", "Challenge-Response"]
toc: true
---

![Yubikey](/images/yubikey.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Operating System** | Debian 7/8 |
| **Website** | [Yubikey Website](https://www.yubico.com/) |
| **Last Update** | 22/04/2015 |
{{< /table >}}

## Introduction

I've bought Yubikeys to manage several things. They permit 2 different kinds of authentication per key. The authentication methods are:

- Yubico OTP
- OATH-HOTP
- Static Password
- Challenge-Response

The goal was to authenticate through my Yubikey without a password, but still have the possibility to connect with my user password if I lose my key. Another requirement is to lock my computer if I remove the key.

## Installation

To install, we'll use packages. One for PAM and the GUI for configuration:

```bash
aptitude install yubikey-personalization-gui libpam-yubico
```

## Configuration

### Challenge-Response

#### Gui

We're now going to configure the key. Insert it and launch the GUI:

```bash
yubikey-personalization-gui
```

Select the challenge-response menu and SHA-1 challenge:

![Yubi gui1.png](/images/yubi_gui1.avif)

Then select all options as shown in the screen below:

![Yubi gui2.png](/images/yubi_gui2.avif)

The 3rd step, if checked, requires a 2.5 second key press to unlock.

Then save the challenge-response to your user settings:

```bash
ykpamcfg -v -2
```

#### pam

The configuration of PAM is quick and easy, simply add this line (`/etc/pam.d/common-auth`):

```bash {linenos=table,hl_lines=[17]}
#
# /etc/pam.d/common-auth - authentication settings common to all services
#
# This file is included from other service-specific PAM config files,
# and should contain a list of the authentication modules that define
# the central authentication scheme for use on the system
# (e.g., /etc/shadow, LDAP, Kerberos, etc.).  The default is to use the
# traditional Unix authentication mechanisms.
#
# As of pam 1.0.1-6, this file is managed by pam-auth-update by default.
# To take advantage of this, it is recommended that you configure any
# local modules either before or after the default block, and use
# pam-auth-update to manage selection of other modules.  See
# pam-auth-update(8) for details.

# here are the per-package modules (the "Primary" block)
auth    sufficient           pam_yubico.so  mode=challenge-response
auth    [success=2 default=ignore]      pam_unix.so nullok_secure
auth    [success=1 default=ignore]      pam_winbind.so krb5_auth krb5_ccache_type=FILE cached_login try_first_pass
# here's the fallback if no module succeeds
auth    requisite                       pam_deny.so
# prime the stack with a positive return value if there isn't one already;
# this avoids us returning an error just because nothing sets a success code
# since the modules above will each just jump around
auth    required                        pam_permit.so
# and here are more per-package modules (the "Additional" block)
auth    optional                        pam_cap.so
# end of pam-auth-update config
```

#### udev

We'll install the udev rule:

```bash
cp /lib/udev/rules.d/69-yubikey.rules /etc/udev/rules.d/
```

and override it to add a custom script (screensaver lock) (`/lib/udev/rules.d/69-yubikey.rules`):

```bash {linenos=table,hl_lines=["9-11"]}
ACTION!="add|change", GOTO="yubico_end"

# Udev rules for letting the console user access the Yubikey USB
# device node, needed for challenge/response to work correctly.

# Yubico Yubikey II
ATTRS{idVendor}=="1050", ATTRS{idProduct}=="0010|0110|0111|114|116", \
    ENV{ID_SECURITY_TOKEN}="1"

LABEL="yubico_end"
# Launch on remove
ACTION=="remove", SUBSYSTEM=="usb", ENV{ID_VENDOR_ID}=="1050", ENV{ID_MODEL_ID}=="0010", RUN+="/path/yubi_remove_script.sh"
# Launch on insert
# ACTION=="add", SUBSYSTEM=="usb", ATTRS{idVendor}=="1050", ATTRS{idProduct}=="0010", RUN+="/path/yubi_add_script.sh"
```

You can test if udev sees your key correctly with this command (try to insert and remove it):

```bash
udevadm monitor --property
```

Then reload udev rules:

```bash
udevadm control --reload-rules
udevadm trigger
```

And create the script where you've declared it (`yubi_script.sh`):

```bash {linenos=table,hl_lines=[3]}
#! /bin/bash
export DISPLAY=":0"
su <username> -c "/usr/bin/xscreensaver-command -lock"
```

And change the username to your desired one. Don't forget to add execution rights:

```bash
chmod 755 yubi_script.sh
```

## FAQ

### How do I enable debug?

It's easy to add debug mode. Simply add "debug" to the PAM line (`/etc/pam.d/common-auth`):

```bash {linenos=table,hl_lines=[17]}
#
# /etc/pam.d/common-auth - authentication settings common to all services
#
# This file is included from other service-specific PAM config files,
# and should contain a list of the authentication modules that define
# the central authentication scheme for use on the system
# (e.g., /etc/shadow, LDAP, Kerberos, etc.).  The default is to use the
# traditional Unix authentication mechanisms.
#
# As of pam 1.0.1-6, this file is managed by pam-auth-update by default.
# To take advantage of this, it is recommended that you configure any
# local modules either before or after the default block, and use
# pam-auth-update to manage selection of other modules.  See
# pam-auth-update(8) for details.

# here are the per-package modules (the "Primary" block)
auth    sufficient           pam_yubico.so  debug mode=challenge-response
auth    [success=2 default=ignore]      pam_unix.so nullok_secure
auth    [success=1 default=ignore]      pam_winbind.so krb5_auth krb5_ccache_type=FILE cached_login try_first_pass
# here's the fallback if no module succeeds
auth    requisite                       pam_deny.so
# prime the stack with a positive return value if there isn't one already;
# this avoids us returning an error just because nothing sets a success code
# since the modules above will each just jump around
auth    required                        pam_permit.so
# and here are more per-package modules (the "Additional" block)
auth    optional                        pam_cap.so
# end of pam-auth-update config
```

and create debug file information:

```bash
touch /var/run/pam-debug.log
chmod 666 /var/run/pam-debug.log
```

You can now check the `/var/run/pam-debug.log` file.

### USB error: Access denied (insufficient permissions)

A possible solution is to add a group to udev and make your user belong to that group. Example (`/etc/udev/rules.d/70-yubikey.rules`):

```bash
ACTION=="add|change", SUBSYSTEM=="usb", ATTRS{idVendor}=="1050", ATTRS{idProduct}=="0010", MODE="0664", GROUP="yubikey"
```

Here you need to create a "yubikey" group and add your current user to that group.

Now reload the rules:

```bash
udevadm control --reload-rules
udevadm trigger
```

It should work now.

## References

http://craoc.fr/doku.php?id=yubikey#configuration_de_la_yubikey
