---
title: "WSUS: Save Disk Space (Purge Unused Updates)"
slug: wsus-save-disk-space-purge-unused-updates/
description: "How to free up disk space on a WSUS server by purging unused Windows updates using WSUSDebug tool."
categories: ["Windows", "Servers"]
date: "2007-08-13T09:02:00+02:00"
lastmod: "2007-08-13T09:02:00+02:00"
tags: ["Windows", "WSUS", "Disk Management", "Windows Server"]
---

To save space on your WSUS server, you can purge unused updates. There's no need to keep old updates that are no longer useful.

For this, you need to use the "WSUSDebug" tool ([available here for direct download](../../static/others/wsus_server_debug_tool.zip) or [here](https://www.laboratoire-microsoft.org/d/?id=16846)).

Download and install it in "C:\Windows".
Then open the command prompt and type:

```bash
WsusDebugTool.exe /Tool:PurgeUnneededFiles
wsusutil.exe deleteunneededrevisions
WSUSUTIL.exe Reset
WSUSUTIL.exe Removeinactiveapprovals
```
