---
weight: 999
url: "/Gnome-shell_\\:_changement_automatique_de_fond_d'écran/"
title: "Gnome-shell: Automatic Wallpaper Change"
description: "How to set up automatic wallpaper change in Gnome-shell using a simple script and optional tools like Nitrogen"
categories: ["Linux", "Debian"]
date: "2012-10-07T19:14:00+02:00"
lastmod: "2012-10-07T19:14:00+02:00"
tags: ["Gnome", "Desktop", "Customization", "Wallpaper"]
toc: true
---

![Gnome](/images/gnome.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 3.4.2 |
| **Operating System** | Debian 7 |
| **Website** | [Gnome Website](https://www.gnome.org/) |
| **Last Update** | 07/10/2012 |
{{< /table >}}

## Introduction

I couldn't find a native function to regularly change wallpapers in Gnome-shell. That's why I took inspiration from [a small script](https://superuser.com/questions/298050/periodically-changing-wallpaper-under-gnome-3) and improved it.

## Implementation

All my wallpapers are in ~/Images. We'll create a small script, which I put in /usr/bin but you can put it wherever you want:

(`/usr/bin/wallpaper-changer`)

```bash {linenos=table,hl_lines=[9,10]}
#!/bin/bash
IMAGES_FOLDER='/home/pmavro/Images'
RANDOM_EVERY=1800

cd $IMAGES_FOLDER
while [ 1 ] ; do
  random_num=`ls $IMAGES_FOLDER | sort -R | tail -n 1`
  # Gnome shell usage
  # gsettings set org.gnome.desktop.background picture-uri "file:$IMAGES_FOLDER/$random_num"
  # Nitrogen usage
  /usr/bin/nitrogen --set-scaled $IMAGES_FOLDER/$random_num
  # XFCE
  xfconf-query -c xfce4-desktop -p /backdrop/screen0/monitor0/image-path -s $IMAGES_FOLDER/$random_num
  sleep $RANDOM_EVERY
done
```

Uncomment one of the 3 lines between Gnome Shell, XFCE or Nitrogen to choose the desired method.

Let's set the execution rights:

```bash
chmod 755 /usr/bin/wallpaper-changer
```

Then add an autostart to your session:

(`~/.config/autostart/wallpaper.desktop`)

```text
[Desktop Entry]
Name=wallpaper-changer
Exec=/usr/bin/wallpaper-changer
Comment=change wallpaper every so often
Hidden=false
Type=Application
X-GNOME-Autostart-enabled=true
```

Restart your session and you're done :-)

### Nitrogen

If you have multiple screens, you might want to have properly placed dual-screen wallpapers. For this, we'll use the nitrogen utility:

```bash
aptitude install nitrogen
```

Then, you must disable desktop icons if you want to use Nitrogen:

```bash
gsettings set org.gnome.desktop.background show-desktop-icons false
```

## References

http://superuser.com/questions/298050/periodically-changing-wallpaper-under-gnome-3
