---
weight: 999
url: "/Utiliser_le_pavé_numérique_comme_souris_avec_X/"
title: "Using the Numeric Keypad as a Mouse with X"
description: "This article explains how to use the numeric keypad to control the mouse cursor in X Window System when your mouse isn't working."
categories: ["Linux"]
date: "2008-12-27T23:55:00+02:00"
lastmod: "2008-12-27T23:55:00+02:00"
tags: ["X", "Keyboard", "Mouse", "Linux", "Tips"]
toc: true
---

## Introduction

In the series "things you should know", we have the control of the mouse pointer using the keyboard in X. If your mouse dies, it can be a real hassle - what to do? Simple, activate keyboard control.

## Usage

To activate:

```bash
[CTRL]+[SHIFT]+[NUMLOCK]
```

The pointer can then be moved using the numeric keypad keys.

* Click with [0] and [+]
* Return to normal operation with [CTRL]+[SHIFT]+[NUMLOCK]. Small beeps indicate the operation worked correctly.
* You can also use [5] to accelerate the movement.
