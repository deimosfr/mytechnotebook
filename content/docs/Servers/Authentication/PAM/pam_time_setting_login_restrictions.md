
---
weight: 999
url: "/Pam_time\\:_Mettre_des_restrictions_sur_les_logins/"
title: "PAM Time: Setting Login Restrictions"
description: "How to use PAM Time module to set various login restrictions like time-based and access-based restrictions on Linux systems."
categories: ["Linux", "Security"]
date: "2011-05-06T20:36:00+02:00"
lastmod: "2011-05-06T20:36:00+02:00"
tags: ["PAM", "security", "authentication", "login", "access control"]
toc: true
---

## Introduction

pam_time is able to make several kinds of restrictions like:

- Time Based Restrictions
- Access Based Restrictions

I'll explain here how to use those options.

## Usage

### Time Based Restrictions

These examples will limit the login times of certain users. See `/etc/security/time.conf` for more information/examples. In order to place time restrictions on user logins, the following must be placed in `/etc/pam.d/login`:

```bash
account    required    /lib/security/pam_time.so
```

The remaining lines should be placed in `/etc/security/time.conf`.

- Only allow user nikesh to login during on weekdays between 7 am and 5 pm:

```bash
login;*;nikesh;Wd0700-1700
```

- Allow users A & B to login on all days between 8 am and 5 pm except for Sunday.

```bash
login;*;A
```

If a day is specified more than once, it is unset. So in the above example, Sunday is specified twice (Al = All days, Su = Sunday). This causes it to be unset, so this rule applies to all days except Sunday.

### Access Based Restrictions

`/etc/security/access.conf` can be used to restrict access by terminal or host. The following must be placed in `/etc/pam.d/login` in order for these examples to work:

```bash
account    required   /lib/security/pam_access.so
```

- Deny nikesh login access on all terminals except for tty1:

```bash
-:nikesh:ALL EXCEPT tty1
```

- Users in the group operator are only allowed to login from a local terminal:

```bash
-:operator:ALL EXCEPT LOCAL
```

- Allow user A to only login from a trusted server:

```bash
-:A:ALL EXCEPT trusted.somedomain.com
```

## Resources
- [https://linuxpoison.blogspot.com/2009/05/how-to-set-login-time-based.html?utm_source=feedburner&utm_medium=feed&utm_campaign=Feed%3A+blogspot%2FfrEh+%28Linux+Poison%29&utm_content=Google+Reader](https://linuxpoison.blogspot.com/2009/05/how-to-set-login-time-based.html?utm_source=feedburner&utm_medium=feed&utm_campaign=Feed%3A+blogspot%2FfrEh+%28Linux+Poison%29&utm_content=Google+Reader)
