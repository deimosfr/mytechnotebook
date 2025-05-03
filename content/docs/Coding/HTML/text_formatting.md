---
weight: 999
url: "/Mise_en_forme_du_texte/"
title: "Text Formatting"
description: "Guide to HTML text formatting including paragraphs, line breaks, alignment, importance tags, font styling, and more."
categories:
  - "Linux"
  - "Development"
date: "2007-09-04T22:41:00+02:00"
lastmod: "2007-09-04T22:41:00+02:00"
toc: true
---

## The Tags

### Paragraph Tags

This is a very simple tag: `<p>`  
You need to open it when starting a paragraph, and close it when ending one:

```html
<p>example</p>
```

### Line Break Tags

Note the final slash indicating it's a solitary tag.  
We can play with paragraphs and line breaks like this:

```html
<p>
  <span style="color:#FF0000"> example </span> <br />
  <span style="color:#3366FF">example2</span>
</p>
```

This will simply write:

```
example
example2
```

### Alignment Tags

We can include an alignment attribute in the paragraph tag.  
This is the align attribute. It allows placing text either on the left side of the page (left value), on the right (right value), or centered (center value).  
Example:

```html
<p align="right">example</p>
```

### Importance Tags

We can emphasize a paragraph more.  
For this, we'll use the `<hx>` tags where x represents a natural integer between 1 and 6.

These tag families modify the importance given to elements between the opening and closing tags. The closer to 1, the more important the text. And conversely if we move away from it:

```html
<p> <h1> example </h1> </p>
```

### Italic, Bold and Underlining

To put text in _italic_, use this tag:

```html
<em></em>
```

To make it **bold**, use:

```html
<strong></strong>
```

While for underlining, use:

```html
<u></u>
```

### Subscripts and Superscripts

For superscripts, use this tag:

```html
<sup></sup>
```

While for subscripts, use:

```html
<sub></sub>
```

### BlockQuote

To insert text with a slight margin on the right, you can use this tag. There isn't much to say about it, except that it's possible to nest several inside each other:

```html
<blockquote></blockquote>
```

### Acronym Tags

An acronym is an abbreviation like C.A.F., C.E.O., etc...  
To keep track, we can offer our visitors the translation of these without cluttering the page.  
For this, simply use the `<acronym>` tag.  
It takes the meaning of this acronym as a title argument.  
Example:

```html
the <acronym title="National Railway Company"> SNCF </acronym>
```

This will display "the SNCF". If the internet user places their mouse cursor over it, a small bubble appears displaying the meaning of this acronym.

### Font Tags

There are many fonts around the world.  
The Internet user doesn't have all of them by default.  
If they don't have the font used, well, it will be the default font, defined by the browser, that will be used.
Not great for aesthetics... So be careful not to choose too exotic fonts.

```html
<font face="Name of font to use"> the text </font>
```

You can also modify the size of the font used, thanks to the size argument.
This takes natural integers between 1 and 7 as values.

```html
<font face="Name of font to use" size="5"> the text </font>
```

### Summary of Basic Tags

{{< table "table-hover table-striped" >}}
| **Tag** | **Description** | **Example** |
|---------|----------------|------------|
| `<i>` | to italicize | _test_ |
| `<b>` | to bold | **test** |
| `<u>` | to underline | test |
| `<small>` | to decrease character size | test |
| `<big>` | to increase character size | test |
| `<sub>` | to subscript | test |
| `<sup>` | to superscript | test |
| `<tt>` | for fixed-width font | test |
| `<s>` | to strikethrough | test |
| `<strike>` | same as above | test |
{{< /table >}}
