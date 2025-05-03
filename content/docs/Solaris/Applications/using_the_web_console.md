---
weight: 999
url: "/Utiliser_la_webconsole/"
title: "Using the Web Console"
description: "Guide to using Sun's web console for managing applications through a web interface"
categories:
  - "Solaris"
  - "Linux"
date: "2008-11-27T15:55:00+02:00"
lastmod: "2008-11-27T15:55:00+02:00"
tags:
  - "Solaris"
  - "Servers"
  - "Network"
  - "Development"
toc: true
---

## Introduction

The web console is a tool that allows you to access SUN application management via a web interface. For example, it's possible to administer ZFS pools and partitions or manage your cluster entirely through a web interface.

This is very convenient for the average user and even more so when you can save time by delegating recurring tasks to a third party (non-experienced) person. That's why I find the web console very useful. To use it, simply connect to this address: https://127.0.0.1:6789

## Registering an Application

Why register an application? Well, because for example, you've updated your Solaris and as usual, the web console goes haywire. So to re-register your applications, we can first list what's working:

```bash
wcadmin list -a
```

```
Deployed web applications (application name, context name, status):

    console  ROOT            [running]
    console  com_sun_web_ui  [running]
    console  console         [running]
    console  manager         [running]
```

Let's list the existing applications:

```bash
ls /usr/share/webconsole/webapps/
```

```
$ ls /usr/share/webconsole/webapps/
com_sun_web_ui/ console/        zfs/
```

Here I see that I have ZFS, and that's what I decide to reactivate. To do this, it's simple:

```bash
smreg add -a /usr/share/webconsole/webapps/zfs
```

Now I just need to reboot the web console for the change to take effect:

```bash
svcadm restart webconsole
```

Now access the web console and voilà, we've got ZFS back. The wcadmin command now gives us this information:

```bash
wcadmin list -a
```

```
Deployed web applications (application name, context name, status):

    console  ROOT            [running]
    console  com_sun_web_ui  [running]
    console  console         [running]
    console  manager         [running]
    legacy   zfs             [running]
```

## Resources
- http://docs.sun.com/app/docs/doc/817-1985/gcrrb?a=view
