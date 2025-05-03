---
weight: 999
url: "/Voir_les_fichiers_et_dossiers_cachés_sous_Mac_OS_X/"
title: "Viewing Hidden Files and Folders in Mac OS X"
description: "How to view hidden files and folders in Mac OS X system."
categories:
  - "MacOS"
date: "2009-11-28T15:59:00+02:00"
lastmod: "2009-11-28T15:59:00+02:00"
tags:
  - "MacOS"
  - "Terminal"
  - "System"
toc: true
---

## Introduction

By default Mac OS X doesn't show any hidden files or folders. If you want to see everything, follow these steps.

## Usage

Show all hidden objects:

```bash
$ defaults write com.apple.Finder AppleShowAllFiles TRUE
```

Swap TRUE with FALSE to turn it off again. Note: Finder must be relaunched afterwards to see the effect. For example like this:

```bash
$ killall Finder && open /System/Library/CoreServices/Finder.app
```
