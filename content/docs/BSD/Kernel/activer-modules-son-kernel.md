---
weight: 999
url: "/Activer_modules_son_kernel/"
title: "Activating Kernel Modules"
description: "How to enable, disable, and manage kernel modules in BSD and Linux systems"
categories: ["BSD", "System Administration"]
date: "2008-05-04T06:48:00+02:00"
lastmod: "2008-05-04T06:48:00+02:00"
tags: ["kernel", "modules", "bsd", "linux", "acpi", "config"]
---

## BSD Systems

Here's the type of command you can use to disable ACPI (among other features) in your BSD kernel:

```bash
config -ef /bsd
disable acpi
quit
```

To re-enable it, simply replace `disable` with `enable`.
