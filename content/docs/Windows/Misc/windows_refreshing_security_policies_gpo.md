---
weight: 999
url: "/Windows_\\:_Rafraîchir_les_politiques_de_sécurité_(GPO)/"
title: "Windows: Refreshing Security Policies (GPO)"
description: "How to refresh Windows security policies (GPO) without disconnecting from the session."
categories: ["Windows"]
date: "2006-10-23T10:59:00+02:00"
lastmod: "2006-10-23T10:59:00+02:00"
tags: ["Windows", "Security", "GPO", "Administration"]
toc: true
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
