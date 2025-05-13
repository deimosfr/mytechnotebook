---
title: "Fluxbox: Rounded Corners for All Windows"
slug: fluxbox-rounded-corners-for-all-windows/
description: "How to add rounded corners to all windows in Fluxbox window manager"
categories: ["Linux", "Desktop"]
date: "2007-11-22T07:20:00+02:00"
lastmod: "2007-11-22T07:20:00+02:00"
tags: ["fluxbox", "desktop", "customization", "linux"]
---

Edit `/usr/share/fluxbox/styles/`**current_theme_name**, and insert these lines:

```bash
menu.roundCorners:                      topleft topright bottomleft bottomright
window.roundCorners:                    topleft topright bottomleft bottomrondie
toolbar.roundCorners:                   topleft topright bottomleft bottomrondie
```

## Resources
- [Get familiar with alternative Linux desktops](../../static/pdf/get_familiar_with_alternative_linux_desktops.pdf)
