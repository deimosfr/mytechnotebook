---
title: "Debian: Changing Default Tools (e.g., Editor)"
slug: debian-changing-default-tools-eg-editor/
description: "How to change default tools preferences in Debian Linux, particularly the default editor using update-alternatives or manual exports."
categories: ["Linux", "Debian"]
date: "2006-11-28T07:50:00+02:00"
lastmod: "2006-11-28T07:50:00+02:00"
tags: ["Debian", "Configuration", "Editor", "Linux"]
---

A nice command exists in Debian to modify your preferences. To change the default editor, for example, do this:

```bash
update-alternatives --config editor
```

Otherwise, manually, you can do this:

```bash
export EDITOR=vim
visudo
```
