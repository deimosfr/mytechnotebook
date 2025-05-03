---
weight: 999
url: "/problemes_d\\'enregistrements_de_fichiers_de_type_office,_adobe.../" 
title: "Problems with saving Office, Adobe... file types"
description: "How to solve issues with saving Office and Adobe files on Samba shares by disabling oplocks or veto files."
categories: ["Linux"]
date: "2010-09-21T21:47:00+02:00"
lastmod: "2010-09-21T21:47:00+02:00"
tags: ["Samba", "Network", "Servers"]
toc: true
---

## Solving Office files saving issues

If you encounter problems saving office files on Samba shares, it's likely due to **Oplock**. Oplock is a feature increasingly used by software developers as it allows much faster saving than normal.

It seems that this problem occurs from version **3.0.6** of Samba.

To avoid these issues, you need to disable this feature. Here's what you need to insert in the Samba configuration file in the "**Global**" section:

```bash
#       Resolve office save problems
        oplocks = no
```

Restart Samba afterwards.

If this solution still doesn't work for you, then it might be because of "**veto files**". Disable temporary files at the share level:

```bash
#       veto files = /*.tmp*/*.TMP*/
```

**Be careful with .DSStore file blocking on Mac which can crash network shares since version 10.6 of Mac OS.**
