---
weight: 999
url: "/Les_caractères_spéciaux/"
title: "Special Characters"
description: "Comprehensive guide on special characters for HTML encoding, including ISO and HTML codes for accented letters and symbols."
categories: ["Linux"]
date: "2007-10-02T20:37:00+02:00"
lastmod: "2007-10-02T20:37:00+02:00"
tags: ["Special pages", "View source", "Network", "Printable version", "cd ~", "Servers", "Development", "Windows", "What links here", "Page information"]
toc: true
---

## Introduction

The default encoding format for HTML pages is **UTF-8**, which is the American format. However, this format does not include our Latin accented characters.

It is necessary to conform to this standard by encoding our special characters in a format that this standard can understand. That's why there are two encodings: one ISO in numeric format and the other specific to HTML expressed in natural language.

An ISO code is written as: &amp;#**code**;, while an HTML code is written as: &amp;**name**;.

## General Coding Rules

HTML codes are mnemonic abbreviations (in English) of accented letters.

Syntax: 

&amp; letter + abbreviation;

Example for **É**: 

Rule: &amp; E + acute;

Actual code: &amp; Eacute;

List of the most common abbreviations: 

{{< table "table-hover table-striped" >}}
| Description | HTML Abbreviation |
|-------------|------------------|
| grave accent | `grave` |
| acute accent | `acute` |
| circumflex accent | `circ` |
| cedilla | `cedil` |
| umlaut | `uml` |
| tilde | `tilde` |
{{< /table >}}

## Table of ISO and HTML Codes for Special Characters

{{< table "table-hover table-striped" >}}
| Character | ISO Code | HTML Abbreviation |
|-----------|----------|------------------|
| **Non-breaking Space** |  |  |
|   | `&#160;` | `&nbsp;` |
| **A** |  |  |
| À | `&#192;` | `&Agrave;` |
| Á | `&#193;` | `&Aacute;` |
| Â | `&#194;` | `&Acirc;` |
| Ã | `&#195;` | `&Atilde;` |
| Ä | `&#196;` | `&Auml;` |
| Å | `&#197;` | `&Aring;` |
| Æ | `&#198;` | `&Aelig;` |
| à | `&#224;` | `&agrave;` |
| á | `&#225;` | `&aacute;` |
| â | `&#226;` | `&acirc;` |
| ã | `&#227;` | `&atilde;` |
| ä | `&#228;` | `&auml;` |
| å | `&#229;` | `&aring;` |
| æ | `&#230;` | `&aelig;` |
| **C** |  |  |
| Ç | `&#199;` | `&Ccedil;` |
| ç | `&#231;` | `&ccedil;` |
| **D** |  |  |
| Ð | `&#208;` | `&ETH;` |
| ð | `&#240;` | `&eth;` |
| **E** |  |  |
| È | `&#200;` | `&Egrave;` |
| É | `&#201;` | `&Eacute;` |
| Ê | `&#202;` | `&Ecirc;` |
| Ë | `&#203;` | `&Euml;` |
| è | `&#232;` | `&egrave;` |
| é | `&#233;` | `&eacute;` |
| ê | `&#234;` | `&ecirc;` |
| ë | `&#235;` | `&euml;` |
| **I** |  |  |
| Ì | `&#204;` | `&Igrave;` |
| Í | `&#205;` | `&Iacute;` |
| Î | `&#206;` | `&Icirc;` |
| Ï | `&#207;` | `&Iuml;` |
| ì | `&#236;` | `&igrave;` |
| í | `&#237;` | `&iacute;` |
| î | `&#238;` | `&icirc;` |
| ï | `&#239;` | `&iuml;` |
| **N** |  |  |
| Ñ | `&#209;` | `&Ntilde;` |
| ñ | `&#241;` | `&ntilde;` |
| **O** |  |  |
| Ò | `&#210;` | `&Ograve;` |
| Ó | `&#211;` | `&Oacute;` |
| Ô | `&#212;` | `&Ocirc;` |
| Õ | `&#213;` | `&Otilde;` |
| Ö | `&#214;` | `&Ouml;` |
| Ø | `&#216;` | `&Oslash;` |
| Œ | `&#140;` | `&OElig;` |
| ò | `&#242;` | `&ograve;` |
| ó | `&#243;` | `&oacute;` |
| ô | `&#244;` | `&ocirc;` |
| õ | `&#245;` | `&otilde;` |
| ö | `&#246;` | `&ouml;` |
| ø | `&#248;` | `&oslash;` |
| œ | `&#156;` | `&oelig;` |
| **S** |  |  |
| Š | `&#138;` | |
| š | `&#154;` | |
| **U** |  |  |
| Ù | `&#217;` | `&Ugrave;` |
| Ú | `&#218;` | `&Uacute;` |
| Û | `&#219;` | `&Ucirc;` |
| Ü | `&#220;` | `&Uuml;` |
| ù | `&#249;` | `&ugrave;` |
| ú | `&#250;` | `&uacute;` |
| û | `&#251;` | `&ucirc;` |
| ü | `&#252;` | `&uuml;` |
| **Y** |  |  |
| Ý | `&#221;` | `&Yacute;` |
| Ÿ | `&#159;` | `&Yuml;` |
| ý | `&#253;` | `&yacute;` |
| ÿ | `&#255;` | `&yuml;` |
| **Z** |  |  |
| Ž | `&#142;` | |
| ž | `&#158;` | |
| **Currency Symbols** |  |  |
| ¢ | `&#162;` | `&cent;` |
| £ | `&#163;` | `&pound;` |
| ¥ | `&#165;` | `&yen;` |
| **Legal Symbols** |  |  |
| ™ | `&#153;` | |
| © | `&#169;` | `&copy;` |
| ® | `&#174;` | `&reg;` |
| **Numerical Symbols** |  |  |
| ‰ | `&#137;` | |
| ª | `&#170;` | `&ordf;` |
| º | `&#186;` | `&ordm;` |
| ¹ | `&#185;` | `&sup1;` |
| ² | `&#178;` | `&sup2;` |
| ³ | `&#179;` | `&sup3;` |
| ¼ | `&#188;` | `&frac14;` |
| ½ | `&#189;` | `&frac12;` |
| ¾ | `&#190;` | `&frac34;` |
| ÷ | `&#247;` | `&divide;` |
| × | `&#215;` | `&times;` |
| > | `&#155;` | `&gt;` |
| < | `&#139;` | `&lt;` |
| ± | `&#177;` | `&plusmn;` |
| **Other Symbols** |  |  |
| & | | `&amp;` |
| ‚ | `&#130;` | |
| ƒ | `&#131;` | |
| „ | `&#132;` | |
| … | `&#133;` | |
| † | `&#134;` | |
| ‡ | `&#135;` | |
| ˆ | `&#136;` | |
| ' | `&#145;` | |
| ' | `&#146;` | |
| " | `&#147;` | |
| " | `&#148;` | |
| • | `&#149;` | |
| – | `&#150;` | |
| — | `&#151;` | |
| ˜ | `&#152;` | |
| ¿ | `&#191;` | `&iquest;` |
| ¡ | `&#161;` | `&iexcl;` |
| ¤ | `&#164;` | `&curren;` |
| ¦ | `&#166;` | `&brvbar;` |
| § | `&#167;` | `&sect;` |
| ¨ | `&#168;` | `&uml;` |
| « | `&#171;` | `&laquo;` |
| » | `&#187;` | `&raquo;` |
| ¬ | `&#172;` | `&not;` |
| ¯ | `&#175;` | |
| ´ | `&#180;` | `&acute;` |
| µ | `&#181;` | `&micro;` |
| ¶ | `&#182;` | `&para;` |
| · | `&#183;` | `&middot;` |
| ¸ | `&#184;` | `&cedil;` |
| Þ | `&#222;` | `&thorn;` |
| ß | `&#223;` | `&szlig;` |
{{< /table >}}
