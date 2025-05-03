---
weight: 999
url: "/CSS_\\:_Les_feuilles_de_style/"
title: "CSS: Style Sheets"
description: "A guide to CSS pseudo-classes and their usage in web development, including link states, text formatting, and page styling."
categories: ["Web", "Development", "CSS"]
date: "2008-10-13T19:48:00+02:00"
lastmod: "2008-10-13T19:48:00+02:00"
tags: ["css", "web design", "pseudo-classes", "html", "styling"]
toc: true
---

## Introduction

When developing a style sheet, it is sometimes useful to apply a style to just one small element (e.g., the first line of a paragraph) or to an element with a certain behavior (e.g., a link being hovered over). Common selectors like "body", "p", or "h1" do not allow for such precise style application. There are in fact other types of selectors, called pseudo-formats or pseudo-classes, that allow for more refined styling.

## Syntax

The name of a pseudo-class is always written in the following way: first the name of the HTML tag concerned, followed by a colon, and finally the mention of the behavior or position. For example:

- a:visited (to apply a particular style to visited links);
- a:hover (to apply a particular style to links being hovered over);
- p:first-letter (to apply a particular style to the first letter of the paragraph).

Note that you cannot modify these names; they are predefined keywords in the CSS format. It is not possible to create your own pseudo-classes.

## Pseudo-classes for Links

The following pseudo-classes apply, as you might guess, only to the `<a>` tag. There are five of them:

- :link, allows you to apply a style to links that have not yet been visited;
- :visited, allows you to apply a style to links that have already been visited;
- :hover, allows you to apply a style to links "hovered over" by the mouse;
- :active, allows you to apply a style to links that are being clicked (the style is applied as long as the finger remains pressed and ceases when the button is released);
- :focus, allows you to apply a style to a link when it is the target of focus (for example, navigating from link to link via the [tab] key means that each link is successively the target of focus).

Note that for correct interpretation of pseudo-formats, to avoid inheritance errors, you should declare them in this order: :link, :visited, :hover, :active.

Be aware that the pseudo-classes :hover, :active and :focus are called "dynamic pseudo-classes", as they allow modifying the style of a tag according to an event, often a user intervention. They can very well be applied to tags other than links: titles, paragraphs, etc.

Note, however, that older browsers do not interpret all these pseudo-classes, particularly :hover, :active and :focus, which are CSS 2 implementations.

## Pseudo-classes for Texts

Text pseudo-classes allow you to apply a style to a well-defined portion of text. Thus, text pseudo-classes are generally used with the paragraph tag `<p>`.
There are two main text pseudo-classes:

- :first-line, which allows you to apply a style to the first line of the paragraph;
- :first-letter, which allows you to apply a style to the first letter of the paragraph (used in the case of drop caps, see the previous tutorial).

You can also choose to insert content (static text or variable content) before or after an element, using the pseudo-classes :before and :after. The content to be inserted is then specified with the content property. The static text to be inserted must be placed in quotes, as follows:

```css
h1:before {
  content: "hello!";
}
```

The above rule has the effect of writing the word "hello!" just before the level 1 heading. It is also possible to insert an image; in this case, content: should be followed by the mention url(image_name.png).

Note that these two pseudo-classes, CSS2 implementations, are not yet interpreted by Internet Explorer.

## Descendant Pseudo-class

A "descendant" pseudo-class allows you to apply a style to the first tag contained in a parent tag pair.
The syntax of this pseudo-class is as follows:

```
parent_tag > first_tag:first-child { property declarations }
```

Let's take the example below:

```css
#header1 > p:first-child {
  color: red;
}
```

And consider the corresponding HTML code portion:

```html
<div id="header1">
  <p>My text here</p>
</div>
```

The text contained between the `<p>...</p>` tags will indeed be written in red, as specified in the style sheet. On the other hand, if you add an image, like this:

```html
<div id="header1">
  <img src="image_name.png" id="logo" alt="product logo" />
  <p>My text here</p>
</div>
```

This won't work anymore because `<p>` is no longer the first tag encountered within the div header1.
Note that Internet Explorer does not yet interpret the :first-child pseudo-class.

## Page Pseudo-classes

It is the @page selector that allows you to define page layout parameters for printing (dimensions, margins, etc.). There are pseudo-classes that allow you to define a specific style for right, left, or the first page of a document. Thus:

- @page:first, allows you to define the style of the first page of a document;
- @page:left, allows you to define the style of left pages;
- @page:right, allows you to define the style of right pages.

## Resources
- http://www.unixgarden.com/index.php/web/les-pseudo-classes-en-css
- [Offering multiple CSS based on browser](/pdf/proposer_plusieurs_css_en_fonction_du_navigateur.pdf)
- [A colorful menu](/pdf/un_menu_haut_en_couleurs.pdf)
