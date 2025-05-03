---
weight: 999
url: "/Utilisation_avancée_de_Mediawiki/"
title: "Advanced Usage of MediaWiki"
description: "Advanced usage of MediaWiki features including dynamic tables, table simplification, and sorting capabilities."
categories:
  - Linux
date: "2007-10-13T09:29:00+02:00"
lastmod: "2007-10-13T09:29:00+02:00"
tags:
  - Development
  - Servers
  - Network
  - Resume
  - Solaris
  - "View source"
  - "cd ~"
  - "What links here"
  - "Special pages"
  - search
toc: true
---

## Introduction

We're very far from realizing how powerful MediaWiki truly is. Those who are familiar with Confluence (a wiki solution more targeted at businesses) know that MediaWiki is far ahead and can sometimes be complex when you want to add dynamic features to your wiki.

That's why I'll try to document the advanced uses of MediaWiki here.

If you're just starting with MediaWiki, I recommend first reading [this]({{< ref "docs/Servers/Webapps/Mediawiki/mediawiki_installation_and_configuration.md" >}}).

## Uses

### Tables

#### Dynamic Modifications

You should read the [table writing simplification](#simplification-of-table-writing) before continuing. Edit your template and look at the example below:

```
|-
|align="left"|{{{1}}}
|align="center" {{#switch: {{{2|}}}
 | yes
 | YES
 | Yes = style="background:palegreen"
 | no
 | NO
 | No = style="background:salmon"
 | partial
 | PARTIAL
 | Partial = style="background:skyblue"}}|{{{2}}}
|align="center" {{#switch: {{{3|}}}
 | yes
 | YES
 | Yes = style="background:palegreen"
 | no
 | NO
 | No = style="background:salmon"
 | partial
 | PARTIAL
 | Partial = style="background:skyblue"}}|{{{3}}}
|align="left"|{{{4}}}
```

Let's study the columns:

- 1st: this one should be familiar to you
- 2nd: We use the **#switch** to indicate that we want to change data according to the content of the text in the cell:

I have "yes, YES or Yes" here, if one of them matches, then I apply a different background color.  
If it's "no, NO or No", it's yet another color.  
And if it's "partial, PARTIAL or Partial", then it's still another color.

- 3rd: it's the same as above, for developers, it resembles a type of if statement, so they won't be too lost.
- 4th: a simple column.

#### Simplification of Table Writing

By default, writing tables isn't particularly simple. That's why we'll create a template to provide a simple writing order. We'll call this template "infos" (**{{infos}}**) and populate it like this:

```
|-
|align="left"|{{{1}}}
|align="center"|{{{2}}}
|align="center" style="background:pink"|{{{3}}}
```

Here, I have 3 columns:

- 1st: text left-aligned
- 2nd: centered text
- 3rd: centered text + pink background fill

So far, this is relatively simple. Then for writing your table, you'll proceed like this:

```
{| width="100%" border="1"
!Name
!First Name
!Company
{{infos|Bill|Gates|Microsoft}}
{{infos|Steve|Jobs|Apple}}
{{infos|Linus|Torvald||}}
|}
```

Still with our 3 columns, we tell it to use the infos template, then fill in the fields. When a field is not available, we leave it blank but **be careful not to forget to put a "|" (pipe) even if we have no more information to fill in!**

#### Sorting

MediaWiki has a default sorting function for its tables. Once your table is made, you'll have the possibility to sort alphabetically for example. The best part is that this solution is ultra-simple to implement! See for yourself... at the beginning of your table, just add **class="wikitable sortable"**:

```
{| class="wikitable sortable" width="100%" border="1"
```

Once you've saved or previewed, you'll be able to sort your columns alphabetically.
