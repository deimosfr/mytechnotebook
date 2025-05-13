---
title: "Windows: Refreshing Security Policies (GPO)"
slug: windows-refreshing-security-policies-gpo/
description: "How to refresh Windows security policies (GPO) without disconnecting from the session."
categories: ["Windows"]
date: "2006-10-23T10:59:00+02:00"
lastmod: "2006-10-23T10:59:00+02:00"
tags: ["Windows", "Security", "GPO", "Administration"]
---

Instead of closing and reopening your session to verify that security policies are working properly, there are commands to avoid disconnecting.

* Windows 2000:

```bash
secedit /refreshpolicy machine_policy /enforce
```

* Windows XP:

```bash
gpupdate /force
```
