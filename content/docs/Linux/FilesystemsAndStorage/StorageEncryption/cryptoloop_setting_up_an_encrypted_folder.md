---
weight: 999
url: "/Cryptoloop_\\:_Mise_en_place_d'un_dossier_crypté/"
title: "Cryptoloop: Setting Up an Encrypted Folder"
description: "Guide on how to set up Cryptoloop for creating encrypted file containers"
categories: ["Linux"]
date: "2008-03-27T13:34:00+02:00"
lastmod: "2008-03-27T13:34:00+02:00"
tags: ["encryption", "security", "cryptoloop", "cryptofs", "linux"]
toc: true
---

Cryptoloop is the predecessor of CryptoFS. It's a useful system because you create an empty encrypted file (image) of a specific size that can contain multiple files or folders.

You can then insert everything inside and everything is encrypted. You can even create a file the size of your partition to make a pseudo-encrypted volume. The disadvantage is that in case of a system crash, if you left the encrypted file open, you have a 50% chance of corrupting your data placed inside.

Therefore, it should be used with caution. Here is some documentation:

[Documentation on Cryptoloop](/pdf/cryptoloop.pdf)  
[Device mapper and loop](/pdf/device_mapper_et_loop.pdf)
