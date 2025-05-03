---
weight: 999
url: "/Chatter_avec_les_personnes_connecté_sur_une_même_machine_via_terminal/"
title: "Chat with Users on the Same Machine via Terminal"
description: "Learn how to communicate with other users logged into the same Linux or Unix system using terminal commands like write and wall."
categories: ["Linux", "Command Line"]
date: "2009-08-08T09:10:00+02:00"
lastmod: "2009-08-08T09:10:00+02:00"
tags: ["chat", "terminal", "write", "wall", "command line"]
toc: true
---

## Introduction

I've been looking several times for a way to chat with other people connected to the same machine as me. And since I'm tired of searching every time how to do it, I'm noting the solutions I found here.

## Solutions

### Chat with a single user

Use the write command:

```bash
write user_you_want_to_talk_to
```

Then type your message and press Enter to validate.

### Chat with all connected users

Use the wall command:

```bash
echo "My message" | wall
```
