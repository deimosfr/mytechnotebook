---
weight: 999
url: "/Xterm_\\:_personnaliser_l'affichage/"
title: "Xterm: Customizing the Display"
description: "How to customize Xterm configuration and display settings for better readability"
categories:
- Linux
date: "2007-05-27T12:25:00+02:00"
lastmod: "2007-05-27T12:25:00+02:00"
tags:
- Linux
- Terminal
- Customization
toc: true
---

## Introduction

When using any shell with a highly customized interface, launching an Xterm terminal can sometimes result in ugly and unreadable colors.

I'd like to share my simple configuration file that I use daily.

## Xdefaults

Simply create a `~/.Xdefaults` file and insert the following content:

```bash
! Deimos Xterm file
! Put the content in ~/.Xdefaults
*xterm*background:       black
*xterm*foreground:       white
*loginShell:            true
 
! Use color for underline attribute
*VT100*colorULMode: on
*VT100*underLine: off
 
! Use color for the bold attribute
*VT100*colorBDMode: on
 
! Love scrollback
*VT100*saveLines: 5000
*VT100*scrollBar: true
```

Now launch a graphical Xterm and the magic will happen :-)
