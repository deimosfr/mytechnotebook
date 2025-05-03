---
weight: 999
url: "/Pam_cracklib_\\:_Choisir_la_complexité_des_mots_de_passe/"
title: "PAM Cracklib: Configure Password Complexity"
description: "How to configure password complexity requirements using PAM Cracklib to enforce strong password policies on Linux systems."
categories: ["Linux", "Security"]
date: "2011-02-06T08:49:00+02:00"
lastmod: "2011-02-06T08:49:00+02:00"
tags: ["PAM", "Security", "Password", "Authentication"]
toc: true
---

## Introduction

If you're tired of users choosing passwords that are too simple on your systems, compromising security in the process, there's a solution. PAM Cracklib allows you to specify the minimum size of passwords, the number of lowercase letters, uppercase letters, digits, and more.

It's almost essential, especially if you rely on a backend like LDAP.

## Installation

Installation is straightforward:

```bash
aptitude install libpam-cracklib
```

## Configuration

We only have one file to edit, which greatly simplifies things. Since we're using Debian, they've made our lives easier - we just need to uncomment some lines that already come with clear explanations!

```bash {linenos=table,hl_lines=[23]}
#
# /etc/pam.d/common-password - password-related modules common to all services
#
# This file is included from other service-specific PAM config files,
# and should contain a list of modules that define the services to be
# used to change user passwords.  The default is pam_unix.

# Explanation of pam_unix options:
#
# The "sha512" option enables salted SHA512 passwords.  Without this option,
# the default is Unix crypt.  Prior releases used the option "md5".
#
# The "obscure" option replaces the old `OBSCURE_CHECKS_ENAB' option in
# login.defs.
#
# See the pam_unix manpage for other options.

# As of pam 1.0.1-6, this file is managed by pam-auth-update by default.
# To take advantage of this, it is recommended that you configure any
# local modules either before or after the default block, and use
# pam-auth-update to manage selection of other modules.  See
# pam-auth-update(8) for details.

# here are the per-package modules (the "Primary" block)
password    requisite           pam_cracklib.so retry=3 minlen=10 difok=3 dcredit=-1 ucredit=-1 lcredit=-1
password    [success=1 default=ignore]  pam_unix.so obscure use_authtok try_first_pass sha512
# here's the fallback if no module succeeds
password    requisite           pam_deny.so
# prime the stack with a positive return value if there isn't one already;
# this avoids us returning an error just because nothing sets a success code
# since the modules above will each just jump around
password    required            pam_permit.so
# and here are more per-package modules (the "Additional" block)
# end of pam-auth-update config
```

I've commented the first line in bold and then uncommented the last two.
On the last line, I also removed 'nullok' which allows empty passwords. Any account with an empty password will be rejected (be careful though if you have a system user that needs this type of account for maintenance operations).

Now, let's explain the pam_cracklib.so line parameters:

- retry: the number of times the user can retry if they enter the wrong password
- minlen: the minimum length of the password
- difok: this is a clever one, and a bit tricky - it remembers previous passwords that users have set. Here a user can't reuse a previously used password until the 5th time.
- dcredit: if the number is negative, it means the password must contain at least x decimal digits to be validated (here at least one digit is required)
- ucredit: if the number is negative, it means the password must contain at least x uppercase letters to be validated (here at least one uppercase letter is required)
- lcredit: if the number is negative, it means the password must contain at least x lowercase letters to be validated (here at least one lowercase letter is required)

I could have also added ocredit which allows you to specify special characters. For dcredit, ucredit, lcredit and ocredit, if they equal a positive number, they subtract from minlen when used.

Check the references below if you want more information :-)

## References

http://www.deer-run.com/~hal/sysadmin/pam_cracklib.html  
http://linux.die.net/man/8/pam_cracklib
