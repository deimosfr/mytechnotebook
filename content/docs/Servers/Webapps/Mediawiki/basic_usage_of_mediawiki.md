---
weight: 999
url: "/Utilisation_basique_de_Mediawiki/"
title: "Basic Usage of MediaWiki"
description: "A concise guide to the basic features of MediaWiki, including links, lists, formatting, and code blocks"
categories:
  - Linux
  - Development
date: "2007-10-13T09:27:00+02:00"
lastmod: "2007-10-13T09:27:00+02:00"
tags:
  - "cd ~"
  - "View source"
  - "search"
  - "What links here"
  - "Servers"
  - "Special pages"
  - "Network"
  - "Development"
  - "Resume"
toc: true
---

## Introduction

For those who want to get started with MediaWiki, it's a good step. Since the documentation on the site is complex and abundant, here's a small summary of what is necessary to have the basics.

## Links

How to create a link to a workshop page?  
Copy and paste its complete title (from the page header), and insert it in your text in the following format:

```
[[Title]]
```

For example, a link to this help page would be written as:

```
[[Modification of Basic Usage of MediaWiki]]
```

However, if it's a category page, you'll need to add a colon (:) before the title:

```
[[:Title]]
```

For example, a link to the category titled Category:Undetermined solution would be written:

```
[[:Category:Undetermined solution]]
```

To specify the label of an internal link to a workshop page, use the syntax:

```
[[Title|label]]
```

(Note the vertical bar | that separates the URL and the label)

### How to link to an external page

Simply entering the URL in your text is enough to generate a clickable link:

```
http://example.org
```

If you want to specify the label of the link, use the syntax:

```
[http://example.org label]
```

(Note the presence of a single bracket, and the space that separates the URL and the label)

## Lists

### How to create an unordered list (bulleted list)

Place each item of the list on a line beginning with an asterisk \*:

```bash
* item 1
* item 2
* item 3
```

- item 1
- item 2
- item 3

### How to create an ordered list (numbered list)

Place each item of the list on a line beginning with a hash #:

```bash
# item 1
# item 2
# item 3
```

1. item 1
2. item 2
3. item 3

### How to create nested lists

For a second-level list, use two asterisks \* or two hashes #:

```bash
* level 1 item
** level 2 item
```

- level 1 item
  - level 2 item

```bash
# level 1 item
## level 2 item
```

1. level 1 item
   1. level 2 item

(Use 3 asterisks \* or 3 hashes # for a 3rd level list, etc.)

## Bold and italic

### How to italicize text

Surround the expression to be italicized with two single apostrophes '':

In my sentence, _an italicized expression_...

If you want the entire line to be in italics, just put '' at the beginning.

### How to make text bold

Surround the expression to be bolded with three single apostrophes ''':
In my sentence, **a bold expression**...

If you want the entire line to be bold, just put ''' at the beginning.

## Code blocks

### How to create a code element

By simply writing the `<code>` element ;):  
Here is an expression in CSS: `background-color: #fff;`...

### How to create a pre element

By simply writing the `<syntaxhighlight lang=text>` element ;):

```
<syntaxhighlight lang=text>
background-color: #fff;
color: #000
</syntaxhighlight>
...
```

### How to escape wiki syntax

If you want to quote an element of the wiki syntax above without it being interpreted, use the syntax:

```
<nowiki>...</nowiki>
```

## I'm using Wikipedia's wiki syntax, but it doesn't always work. Why?

Unlike regular wikis based on MediaWiki (such as Wikipedia), several syntaxes have been disabled for public editing in the Opquast workshop: particularly the creation of new categories, adding titles and subtitles, and templates. These elements are indeed managed directly by our workflow.
