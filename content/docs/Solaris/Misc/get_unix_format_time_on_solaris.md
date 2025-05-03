---
weight: 999
url: "/Obtenir_l'heure_au_format_Unix_sous_Solaris/"
title: "Get Unix Format Time on Solaris"
description: "A script to get the epoch time format on Solaris systems."
categories:
  - "Linux"
  - "Solaris"
date: "2007-03-06T16:45:00+02:00"
lastmod: "2007-03-06T16:45:00+02:00"
tags:
  - "Solaris"
  - "Unix"
  - "Time"
  - "Scripts"
toc: true
---

## Getting Unix Format Time on Solaris

Here is a little script to obtain Unix format on Solaris (Epoch Time):

```bash
#!/bin/sh
/usr/bin/truss /usr/bin/date 2>&1 |  nawk -F= '/^time\(\)/ {gsub(/ /,"",$2);print $2}'
exit $?
```

And finally, let's see it in action:

```bash
# ./edate
1149276150
```
