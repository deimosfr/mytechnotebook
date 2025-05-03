---
weight: 999
url: "/Wine_\\:_Lancer_des_applications_windows_sur_linux/"
title: "Wine: Running Windows Applications on Linux"
description: "A guide on how to use Wine to run Windows applications on Linux systems"
categories: 
  - Linux
date: "2008-06-08T16:44:00+02:00"
lastmod: "2008-06-08T16:44:00+02:00"
tags:
  - Windows
  - Wine
  - Linux
  - Gaming
toc: true
---

## 1 Introduction

Hello everyone! Want to run a Windows application on Linux and finding it complicated? Don't worry! It's actually very easy... just follow this tutorial.

## 2 Installation

As usual:

```bash
apt-get install wine
```

## 3 Configuration

What's next? Do we need to configure lots of things like bottles and all that stuff? Actually no, the latest versions are well-made and the setup is super quick. We'll run winecfg which will create our little environment:

```bash
winecfg
```

Make sure the drives are properly displayed. Then click OK.

## 4 Launching the application

Here, I'll use Spider Solitaire as an example. It's a small game on Windows found at c:\windows\system32\spider.exe. Simply copy it to `/home/$USER/.wine/drive_c/windows/system32`. As you can see, I placed it exactly in the same location as it was on Windows.

Now we're ready to launch it (**with the absolute path**):

```bash
wine /home/$USER/.wine/drive_c/windows/system32/spider.exe
```

Hurray, it works! :-)
