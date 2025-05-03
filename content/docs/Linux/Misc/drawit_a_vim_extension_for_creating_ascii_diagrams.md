---
weight: 999
url: "/DrawIt_\\:_Une_extension_VIM_pour_faire_des_diagrammes_en_ASCII/"
title: "DrawIt: A VIM Extension for Creating ASCII Diagrams"
description: "Guide on how to install and use DrawIt, a VIM extension that allows creating ASCII diagrams directly in the editor."
categories: ["Linux"]
date: "2010-06-05T20:40:00+02:00"
lastmod: "2010-06-05T20:40:00+02:00"
tags: ["Servers", "Mac OS X", "Windows", "DrawIt", "VIM", "ASCII"]
toc: true
---

![Drawit VIM](/images/ascii-drawing-in-vim-editor-300x257.avif)

## Introduction

[DrawIt](https://vim.sourceforge.net/scripts/script.php?script_id=40) allows you to create ASCII diagrams. It's very practical and avoids the hassle of using unnecessary tools.

If you want to convert your ASCII diagrams to images, use [Ditaa](https://ditaa.sourceforge.net/).

## Installation

Installation is quite simple as there is a small installer for vim:

```bash
cd ~
wget -O DrawIt.vba.tgz "http://www.vim.org/scripts/download_script.php?src_id=8798"
gzip -d DrawIt.vba.tgz
mv DrawIt.vba.tar Drawit.vba
vim Drawit.vba
:so %
:q
```

And that's it.

## Utilization

Using it is also quite straightforward. I'm just copying and pasting the documentation as it's clear enough.

Basic commands:

- Activate draw: \di
- Deactivate: \ds

{{< table "table-striped table-hover" >}}
| Key | Description |
|-----|-------------|
| \<left\> | move and draw left |
| \<right\> | move and draw right, inserting lines/space as needed |
| \<up\> | move and draw up, inserting lines/space as needed |
| \<down\> | move and draw down, inserting lines/space as needed |
| \<s-left\> | move left |
| \<s-right\> | move right, inserting lines/space as needed |
| \<s-up\> | move up, inserting lines/space as needed |
| \<s-down\> | move down, inserting lines/space as needed |
| \<space\> | toggle into and out of erase mode |
| \> | draw -> arrow |
| \< | draw <- arrow |
| ^ | draw ^ arrow |
| v | draw v arrow |
| \<pgdn\> | replace with a \\, move down and right, and insert a \\ |
| \<end\> | replace with a /, move down and left, and insert a / |
| \<pgup\> | replace with a /, move up and right, and insert a / |
| \<home\> | replace with a \\, move up and left, and insert a \\ |
| \\> | draw fat -> arrow |
| \\< | draw fat <- arrow |
| \\^ | draw fat ^ arrow |
| \\v | draw fat v arrow |
| \\a | draw arrow based on corners of visual-block |
| \\b | draw box using visual-block selected region |
| \\e | draw an ellipse inside visual-block |
| \\f | fill a figure with some character |
| \\h | create a canvas for \\a \\b \\e \\l |
| \\l | draw line based on corners of visual block |
| \\s | adds spaces to canvas |
| \<leftmouse\> | select visual block |
| \<s-leftmouse\> | drag and draw with current brush (register) |
| \\ra ... \\rz | replace text with given brush/register |
| \\pa ... | like \\ra ... \\rz, except that blanks are considered to be transparent |
{{< /table >}}
