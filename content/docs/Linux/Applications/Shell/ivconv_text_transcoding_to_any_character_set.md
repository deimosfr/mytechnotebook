---
weight: 999
url: "/Ivconv_\\:_Transcodage_de_texte_vers_n'importe_quel_jeu_de_caractères/"
title: "Ivconv: Text Transcoding to Any Character Set"
description: "How to use the iconv command on Debian to transcode text from one character set to another, useful for file conversion between UTF-8 and ISO standards."
categories: ["Debian", "Linux"]
date: "2007-03-08T08:18:00+02:00"
lastmod: "2007-03-08T08:18:00+02:00"
tags: ["Text conversion", "Character encoding", "UTF-8", "ISO", "Command line"]
toc: true
---

Debian provides the iconv command that allows you to transcode text to and from any character set. Example:

```bash
$ iconv -f utf8 -t iso8859-15 fichier_utf8.txt
```

This command will transcode from UTF-8 to ISO-8859-15. It's very useful if you've developed a web page in UTF-8 and realized it too late.
