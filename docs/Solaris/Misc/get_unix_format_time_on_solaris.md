---
title: "Get Unix Format Time on Solaris"
slug: get-unix-format-time-on-solaris/
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
