---
weight: 999
url: "/Introduction_au_XHTML/"
title: "Introduction to XHTML"
description: "This guide provides an introduction to XHTML, including the basics of tags, comments, minimal structure, and advanced features such as video tags."
categories:
  - Linux
date: "2012-09-30T07:38:00+02:00"
lastmod: "2012-09-30T07:38:00+02:00"
tags:
  - Development
  - GitHub
  - Mac OS X
  - Windows
  - navigation
toc: true
---

## Introduction

XHTML is a markup language used for writing World Wide Web pages. Originally designed as the successor to HTML, XHTML is based on the syntax defined by XML, which is more recent and simpler than the syntax defined by SGML on which HTML is based.

Like many XML-based languages, XHTML begins with the letter X, which represents the word "extensible". Thus, the first document officially describing XHTML is called "XHTML™ 1.0 The Extensible HyperText Markup Language". However, the abbreviation XHTML is a trademark of the World Wide Web Consortium (W3C) and is the only one used in the specifications that followed version 1.0.

## Tags

A tag is a keyword surrounded by angle brackets `<` and `>`.
XHTML requires that all tags be written in lowercase.
Tags always work in pairs: an opening tag and a closing tag.
The two tags will transform everything between them.
And, like every rule, this one has its exceptions... There are indeed some that prefer solitude... They can be recognized by the / just before the closing bracket. An example with the following tag, which allows a line break:

```html
<br />
```

Some tags take what are called arguments or attributes. These are actually additional parameters that allow for a little more variety in color choices, position on the page, etc... They are always placed before the closing angle bracket of the opening tag.

Example:

```html
<body bgcolor="#000000"></body>
```

This is the tag indicating the beginning of the body of the page; this tag sets a black background on the page.

## Comments

Comments are pieces of code that will not be interpreted by the browser.
Let's imagine that several people are working on a site. It is highly likely that each person will program a piece of page in their own corner, and that the parts will be combined later. But if there is a problem and someone asks for help, they will need to understand the code! So, adding some indications can only be useful.
To write comments in XHTML, there is a tag specially designed for the occasion. It is quite particular, because besides being a solitary tag, it is written in its own way:

```html
<!-- Write the comment here -->
```

## Minimal Structure

Let's finally begin:

```html
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.1//EN" "http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd">
<html xmlns="http://www.w3c.org/1999/xhtml" xmllg="fr" lang="fr">
  <head>
    <title>My first Web page</title>
  </head>

  <body>
    This is the body of my page!
  </body>
</html>
```

## Explanations

- First of all, you need to write a rather intimidating line:

```html
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.1//EN" "http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd">
```

This line makes it known that we are using the XHTML language, version 1.1. This is the most recent version to date, I believe (as of June 1st, 2005).

By declaring this DocType, we indicate to the Internet browser that we know the Web standards, and that we are going to use them. Thus, it will interpret all the code very rigorously. In this way, the page will respect the various standards in force, and the display of the page will be the same everywhere. If we don't include it, each browser will try to understand the code in its own way, and variations may be seen depending on the browser used. Hence the importance of this line!

- Then the line:

```html
<html xmlns="http://www.w3c.org/1999/xhtml" xmllg="fr" lang="fr">
  (don't forget the
</html>
that closes it!!)
```

Here, we are retrieving everything we need to program in XHTML from the W3C (World Wide Web Consortium) site. This is the international organization that sets the rules and standards in force on the Web.

The first argument indicates the namespace to use. It's a link to a W3C page that defines the different keywords of the language.

The second argument is the language of the pages. In this case, it will be a French page (value fr). For an English page, you will need to put en, for an Italian page it, etc... This is useful to increase the accessibility of pages: search engines will know what the language of the site is, which will allow for better referencing.

The third argument is not normally mandatory, but it is advised to include it in order to make the pages compatible with older browsers, which are not yet up to date with Internet standards (isn't that right, IE?).

- Tags `<head>` and `</head>`

The head area concerns everything that is not on the page itself. It concerns everything around it: in the browser's title bar, in the status bar, etc.
For info:

```html
<title>My first Web page</title>
```

This is the message in the browser's title bar.

- Tags `<body>` and `</body>`

The body of the page will contain everything that should be displayed on the page, in the main area of the browser.

## Video Tags

Video tags appeared with HTML5. There are quite a few drawbacks regarding video formats for each browser. In short, to make everyone agree, we can embed several video formats within the same video tag. A bit complicated? Let's go with an example:

```html
<video controls="" preload="" width="510">
   <source type='video/ogg; codecs="theora, vorbis"' src="http://www.deimos.fr/blog/wp-content/uploads/2010/05/usb_locker.ogv"></source>
   <source type='video/mp4; codecs="avc1.42E01E, mp4a.40.2"' src="http://www.deimos.fr/blog/wp-content/uploads/2010/05/usb_locker.mp4"></source>
   <source type='video/mp4; codecs="avc1.42E01E, mp4a.40.2"' src="http://www.deimos.fr/blog/wp-content/uploads/2010/05/usb_locker_iphone.mp4"></source>
   Internet Explorer is not up to standard to read this video, use instead <a href="http://getfirefox.com">Firefox</a>
</video>
```

- Line 1: The OGV format is used for videos that can be played with Firefox. Use this preferably because it's free and open source. No need to pay anything in terms of licenses.
- Line 2: The mp4 format is proprietary, but recognized by browsers other than Firefox.
- Line 3: More mp4 but this time with modified dimensions so that the video can be played on iPhone (480x320)
- Line 4: This line indicates to IE users that this awful mess still doesn't read videos in HTML5 format. Coincidentally, it's the last browser not to have these features!

## A Loading Page in HTML

Here's a subtle way to use only HTML for a loading page with just a div and innerHTML. I found this solution [here](https://forums.digitalpoint.com/showthread.php?t=16341). So here's a simple form with the 'Submit' button that will take the desired text once clicked:

```html
<form
  action="target.asp"
  onsubmit="document.getElementById('submitdiv').innerHTML='Please wait...'"
>
  other fields
  <div id="submitdiv">
    <input type="Submit" value="Submit" id="Submit" name="Submit" />
  </div>
</form>
```

## References

1. [https://forums.digitalpoint.com/showthread.php?t=16341](https://forums.digitalpoint.com/showthread.php?t=16341)
