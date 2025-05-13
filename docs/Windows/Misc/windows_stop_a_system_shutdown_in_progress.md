---
title: "Windows: Stop a system shutdown in progress"
slug: windows-stop-a-system-shutdown-in-progress/
description: "How to stop a Windows system shutdown that is already in progress using a simple command"
categories:
  - Windows
date: "2007-03-06T13:21:00+02:00"
lastmod: "2007-03-06T13:21:00+02:00"
tags:
  - Windows
  - System
  - Administration
---

Sometimes it's necessary to create scripts that reboot your machine after modifications. Or you may have caught a virus that tells you your PC will reboot in x seconds.

To cancel this reboot, simply use:

```bash
shutdown -a
```
