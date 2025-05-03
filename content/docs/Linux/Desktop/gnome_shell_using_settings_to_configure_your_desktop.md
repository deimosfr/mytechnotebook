---
weight: 999
url: "/Gnome-shell_\\:_utilisation_de_settings_pour_configurer_votre_desktop/"
title: "GNOME Shell: Using Settings to Configure Your Desktop"
description: "A guide on how to use GNOME Shell settings to configure your desktop environment, including showing date/time, workspace settings, changing backgrounds and more."
categories: ["Linux", "Debian"]
date: "2012-11-06T07:19:00+02:00"
lastmod: "2012-11-06T07:19:00+02:00"
tags: ["GNOME", "Desktop Environment", "Linux Configuration", "gsettings"]
toc: true
---

![Gnome](/images/gnome.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 3.4.2 |
| **Operating System** | Debian 7 |
| **Website** | [Gnome Website](https://www.gnome.org/) |
| **Last Update** | 06/11/2012 |
{{< /table >}}

## Introduction

Since GNOME Shell has been released, there is a tool that allows you to configure all sorts of things for your GNOME. Similar to gconf-editor, but from the command line.

## Usage

### Enable Date in the Top Bar

To enable the date next to the time in the top bar of GNOME Shell:

```bash
gsettings set org.gnome.shell.clock show-date true
```

### Display Seconds

To display seconds in the time in the top bar:

```bash
gsettings set org.gnome.shell.clock show-seconds true
```

### Multi Workspace in Dual Screen

To enable all workspaces when you have multiple screens:

```bash
gsettings set org.gnome.shell.overrides workspaces-only-on-primary false
```

### Change Wallpaper

To change the wallpaper:

```bash
gsettings set org.gnome.desktop.background picture-uri "file:/home/pmavro/Images/wallpaper.png"
```

I have also created [an article about automatic wallpaper changes](/Gnome-shell_\:_changement_automatique_de_fond_d'écran/).

### Change Dock Position

If you have enabled the dock, it's possible to change its position:

```bash
gsettings set org.gnome.shell.extensions.dock position left
gsettings set org.gnome.shell.extensions.dock position right
```

### Change Default Applications

Default applications can be changed graphically. However, if what you need is not in the list of available applications, it's possible to modify this. Here is an example I found[^1]:

```bash
[Default Applications]
application/javascript=gvim.desktop
application/lrf=calibre-lrfviewer.desktop
application/msword=libreoffice-writer.desktop
application/rtf=libreoffice-writer.desktop
application/vnd.oasis.opendocument.spreadsheet=libreoffice-calc.desktop
application/vnd.oasis.opendocument.text=libreoffice-writer.desktop
application/vnd.rn-realmedia=mplayer.desktop
application/x-cbr=comix.desktop
application/x-extension-htm=firefox.desktop
application/x-extension-html=firefox.desktop
application/x-extension-shtml=firefox.desktop
application/x-extension-xhtml=firefox.desktop
application/x-extension-xht=firefox.desktop
application/x-perl=gvim.desktop
application/x-php=gvim.desktop
application/x-rar=comix.desktop
application/x-shellscript=gvim.desktop
application/xhtml+xml=firefox.desktop
application/xml=gvim.desktop
application/x-yaml=gvim.desktop
application/zip=comix.desktop
audio/mp4=audacious.desktop
image/gif=gqview.desktop
image/jpeg=gqview.desktop
image/png=gqview.desktop
inode/directory=pcmanfm.desktop;
text/css=gvim.desktop
text/html=firefox.desktop;chromium-browser.desktop;
text/plain=leafpad.desktop;gvim.desktop;
text/x-chdr=gvim.desktop
text/x-csrc=gvim.desktop
text/x-python=gvim.desktop
video/mp4=mplayer.desktop
video/mpeg=mplayer.desktop
video/quicktime=mplayer.desktop
video/webm=mplayer.desktop
video/x-flv=mplayer.desktop
video/x-matroska=mplayer.desktop
video/x-ms-wmv=mplayer.desktop
video/x-msvideo=mplayer.desktop
video/x-ogm+ogg=mplayer.desktop
x-scheme-handler/http=firefox.desktop;chromium-browser.desktop;
x-scheme-handler/https=firefox.desktop;chromium-browser.desktop;
x-scheme-handler/feed=thunderbird.desktop
x-scheme-handler/ftp=firefox.desktop;chromium-browser.desktop;
x-scheme-handler/mailto=thunderbird.desktop
x-scheme-handler/news=thunderbird.desktop
x-scheme-handler/nntp=thunderbird.desktop
x-scheme-handler/snews=thunderbird.desktop

[Added Associations]
application/javascript=gvim.desktop;
application/msword=libreoffice-writer.desktop;
application/rtf=libreoffice-writer.desktop;
application/vnd.oasis.opendocument.spreadsheet=libreoffice-calc.desktop;
application/vnd.oasis.opendocument.text=libreoffice-writer.desktop;
application/vnd.rn-realmedia=mplayer.desktop;
application/x-cbr=comix.desktop;
application/x-extension-htm=firefox.desktop;chromium-browser.desktop;
application/x-extension-html=firefox.desktop;chromium-browser.desktop;
application/x-extension-shtml=firefox.desktop;chromium-browser.desktop;
application/x-extension-xhtml=firefox.desktop;chromium-browser.desktop;
application/x-extension-xht=firefox.desktop;chromium-browser.desktop;
application/x-perl=gvim.desktop;
application/x-php=gvim.desktop;
application/x-rar=comix.desktop;
application/x-shellscript=gvim.desktop;
application/xhtml+xml=firefox.desktop;
application/xml=gvim.desktop;
application/x-yaml=gvim.desktop;
application/zip=comix.desktop;
audio/mp4=audacious.desktop;
image/gif=gqview.desktop;
image/jpeg=gqview.desktop;
image/png=gqview.desktop;
inode/directory=pcmanfm.desktop;
text/css=gvim.desktop;
text/html=firefox.desktop;chromium-browser.desktop;
text/plain=leafpad.desktop;gvim.desktop;
text/x-chdr=gvim.desktop;
text/x-csrc=gvim.desktop;
text/x-python=gvim.desktop;
video/mp4=mplayer.desktop;
video/mpeg=mplayer.desktop;
video/quicktime=mplayer.desktop;
video/webm=mplayer.desktop;
video/x-flv=mplayer.desktop;
video/x-matroska=mplayer.desktop;
video/x-ms-wmv=mplayer.desktop;
video/x-msvideo=mplayer.desktop;
video/x-ogm+ogg=mplayer.desktop;
x-scheme-handler/http=firefox.desktop;chromium-browser.desktop;
x-scheme-handler/https=firefox.desktop;chromium-browser.desktop;
x-scheme-handler/ftp=firefox.desktop;chromium-browser.desktop;
x-scheme-handler/feed=thunderbird.desktop;
x-scheme-handler/mailto=thunderbird.desktop;
x-scheme-handler/news=thunderbird.desktop;
x-scheme-handler/nntp=thunderbird.desktop;
x-scheme-handler/snews=thunderbird.desktop;
```

### Add an Application

It's possible to add an application to GNOME Shell that is not in the repositories. For example, I installed the latest version of Eclipse in `/usr/share` and I want to make the application visible in the list of available applications. Just create a file with this content (`~/.local/share/applications/eclipse.desktop`):

```bash
[Desktop Entry]
Categories=Development
Comment=Eclipse
Encoding=UTF-8
Exec=/usr/share/eclipse/eclipse
GenericName=Eclipse
Hidden=false
# Icons: 64x64 png
Icon=/usr/share/eclipse/icon.png
Name=Eclipse
Type=Application
```

All you need to do is reload GNOME Shell (Alt+F2 - r - Enter).

## References

http://gregcor.com/2011/05/07/fix-dual-monitors-in-gnome-3-aka-my-workspaces-are-broken/

[^1]: https://github.com/ssokolow/profile/blob/master/home/.local/share/applications/mimeapps.list
