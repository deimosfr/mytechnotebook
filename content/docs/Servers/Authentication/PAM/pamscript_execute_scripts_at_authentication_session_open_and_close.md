---
weight: 999
url: "/PAM-script_\\:_Executer_des_scripts_à_l'authentification,_l'ouverture_et_la_fermeture_de_session/"
title: "PAM-script: Execute Scripts at Authentication, Session Open and Close"
description: "Learn how to use PAM-script module to execute scripts during authentication, session opening and closing on Linux systems."
categories: ["Linux", "Security"]
date: "2010-03-14T19:40:00+02:00"
lastmod: "2010-03-14T19:40:00+02:00"
tags: ["PAM", "Security", "Authentication", "Session", "Linux"]
toc: true
---

## Introduction

You may need to run some operations at authentication, session opening or closing. Here is a PAM module I've found that allows this functionality.

## Installation

Download the module from [the Freshmeat project](https://freshmeat.net/projects/pam_script/) and extract it:

```bash
wget http://freshmeat.net/redir/pam_script/22413/url_tgz/libpam-script_0.1.12.tar.gz
tar -xzvf libpam-script_0.1.12.tar.gz
```

Now install the dependencies:

```bash
aptitude install libpam-dev gcc make
```

Now compile it:

```bash
$ make
gcc -Wall -pedantic -fPIC -shared  -o pam_script.so pam_script.c
```

Now you just need to copy it:

```bash
cp pam_script.so /lib/security
```

## Configuration

### PAM

#### Session

If you want to launch something with root permissions at session startup, edit the `/etc/pam.d/common-session` and add this line:

```bash
session required        pam_mkhomedir.so skel=/etc/skel/ umask=0022
session required        pam_script.so runas=root onsessionopen=/etc/security/onsessionopen
session sufficient      pam_ldap.so
session required        pam_unix
```

After pam_script, you can configure:

* runas: choose the user you want to run script (runas=root)
* onsessionopen: this script will be launched on started session (onsessionopen=/etc/security/onsessionopen)
* onsessionclose: this script will be launched on closed session (onsessionclose=/etc/security/onsessionclose)

#### Auth

You may also want to launch something at authentication:

```bash
auth    required        pam_unix.so nullok_secure
auth     required pam_script.so onauth=/etc/security/onauth
```

### Scripts

Just create the default scripts and add the necessary permissions:

```bash
touch /etc/security/onsessionopen /etc/security/onsessionclose /etc/security/onauth
chmod 755 /etc/security/onsessionopen /etc/security/onsessionclose /etc/security/onauth
```

And add this minimum content:

```bash
#!/bin/sh
```

## Test & Debug

You can now test by adding for example "touch /tmp/test_ok" on the "onsessionopen" script. To have more details, please look at the logs:

```bash
$ tail /var/log/auth.log
Jul 15 13:03:35 moonlight sshd[3777]: PAM-script: Real User is: pmavro
Jul 15 13:03:35 moonlight sshd[3777]: PAM-script: Command is:   /etc/security/onsessionopen
Jul 15 13:03:35 moonlight sshd[3777]: PAM-script: Executing uid:gid is: 0:0
```

All looks good :-)
