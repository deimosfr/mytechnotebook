---
weight: 999
url: "/Xmodmap_\\:_mapper_tous_les_boutons_de_sa_souris/"
title: "Xmodmap: Map All Your Mouse Buttons"
description: "How to map all the buttons on your mouse using Xmodmap in Linux"
categories: ["Linux", "Ubuntu"]
date: "2009-05-10T09:42:00+02:00"
lastmod: "2009-05-10T09:42:00+02:00"
tags: ["Mouse", "Xmodmap", "Configuration"]
toc: true
---

## Introduction

Since Ubuntu 8.10, I had lost a small functionality - the ability to paste from clipboard using the left button representing a down arrow on a Logitech MX Revolution mouse.

## Solution

Use xmodmap to map all the possibilities of your mouse:

```bash
xmodmap -e "pointer = 1 8 3 4 5 6 7 2 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27 28 29 30 31 32"
```

You can also use BTNX if you want to assign certain buttons to specific applications or simply standard actions.
