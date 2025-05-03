---
weight: 999
url: "/Windows_XP_\\:_Problèmes_lors_d'un_dépassement_mémoire/"
title: "Windows XP: Memory Overflow Problems"
description: "How to solve memory overflow problems in Windows XP by editing registry settings"
categories: ["Windows"]
date: "2007-08-21T07:15:00+02:00"
lastmod: "2007-08-21T07:15:00+02:00"
tags: ["Windows", "Registry", "Memory", "Windows XP"]
toc: true
---

## Memory Overflow Problems

If you encounter a memory overflow problem that prevents you from launching any application, you need to use "regedit" and edit this registry path:

```
HKEY_LOCAL_MACHINE\System\CurrentControlSet\Control\Session Manager\SubSystems
```

The "Windows" value contains a long line. You'll find a part that looks like this:

```
SharedSection=xxxx,yyyy,zzzz
```

Edit the yyyy section and increase it as shown in this example:

```
Windows SharedSection=1024,3072,512 
```

Modified line:

```
Windows SharedSection=1024,8192,512
```

Here, we're increasing the GDI memory (the application that manages windows) from 3MB to 8MB.
