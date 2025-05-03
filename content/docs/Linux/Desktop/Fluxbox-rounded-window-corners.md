---
weight: 999
url: "/Fluxbox_\\:_Arrondir_les_bords_de_toutes_les_fenêtres/"
title: "Fluxbox: Rounded Corners for All Windows"
description: "How to add rounded corners to all windows in Fluxbox window manager"
categories: ["Linux", "Desktop"]
date: "2007-11-22T07:20:00+02:00"
lastmod: "2007-11-22T07:20:00+02:00"
tags: ["fluxbox", "desktop", "customization", "linux"]
toc: true
---

Edit `/usr/share/fluxbox/styles/`**current_theme_name**, and insert these lines:

```bash
menu.roundCorners:                      topleft topright bottomleft bottomright
window.roundCorners:                    topleft topright bottomleft bottomrondie
toolbar.roundCorners:                   topleft topright bottomleft bottomrondie
```

## Resources
- [Get familiar with alternative Linux desktops](/pdf/get_familiar_with_alternative_linux_desktops.pdf)
