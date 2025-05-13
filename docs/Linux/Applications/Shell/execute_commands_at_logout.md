---
title: "Execute Commands at Logout"
slug: execute-commands-at-logout/
description: "How to execute commands when a user logs out from a system, such as disconnecting a VPN or performing cleanup tasks."
categories: ["Linux"]
date: "2009-12-11T21:21:00+02:00"
lastmod: "2009-12-11T21:21:00+02:00"
tags: ["Servers", "Linux", "Shell", "Terminal"]
---

## Introduction

It can sometimes be useful to run commands when logging out, such as disconnecting a VPN or similar tasks.

## Usage

Here's an example that you can put in your `.profile` file:

```bash
trap cmd 0
```
