---
title: "Finding passwords from the Windows SAM database"
slug: finding-passwords-from-the-windows-sam-database/
description: "How to recover lost Windows passwords from the SAM database using Ophcrack."
categories: ["Windows", "Security"]
date: "2007-08-31T15:06:00+02:00"
lastmod: "2007-08-31T15:06:00+02:00"
tags: ["Windows", "Password", "Recovery", "Security", "SAM"]
---

Are you one of those people who have lost their password? And you neither want to reset it through the recovery console, nor format your Windows installation to be able to use your machine.

Well, [Ophcrack](https://ophcrack.sourceforge.net/) is made for you! It will mount your Windows partition in NTFS, then crack the LM and NTLM hashes, and display them before your amazed eyes - all in just a few minutes :-)
