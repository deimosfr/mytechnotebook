---
weight: 999
url: "/preventing_indexed_website/"
title: "Preventing your website from being indexed (disabling robot scans)"
description: "How to prevent search engines from crawling and indexing your website using robots.txt file and meta tags"
categories: ["Development", "Linux", "Servers"]
date: "2009-09-19T21:11:00+02:00"
lastmod: "2009-09-19T21:11:00+02:00"
tags: ["robots.txt", "SEO", "Web Development", "3.1 Deny All", "search"]
toc: true
---

## Introduction

"Robots.txt" is a regular text file that through its name, has special meaning to the majority of "honorable" robots on the web. By defining a few rules in this text file, you can instruct robots to not crawl and index certain files, directories within your site, or at all. For example, you may not want Google to crawl the `/images` directory of your site, as it's both meaningless to you and a waste of your site's bandwidth. "Robots.txt" lets you tell Google just that.

## Index

Another solution to robot.txt consists in editing your index (index.html/index.php...) of your website and entering this content:

```html
<meta name="robots" content="noindex, nofollow" />
or
<meta name="robots" content="none" />
```

## robots.txt

### Deny All

Here's a basic "robots.txt":

```bash
 User-agent: *
 Disallow: /
```

With the above declared, all robots (indicated by "\*") are instructed to not index any of your pages (indicated by "/"). Most likely not what you want, but you get the idea.

### Deny Google Image Bot

You may not want Google's Image bot crawling your site's images and making them searchable online, if just to save bandwidth. The below declaration will do the trick:

```bash
    User-agent: Googlebot-Image
    Disallow: /
```

### Disable all search engine

The following disallows all search engines and robots from crawling select directories and pages:

```bash
    User-agent: *
    Disallow: /cgi-bin/
    Disallow: /privatedir/
    Disallow: /tutorials/blank.htm
```

### Multiple robots restrictions

You can conditionally target multiple robots in "robots.txt." Take a look at the below:

```bash
    User-agent: *
    Disallow: /
    User-agent: Googlebot
    Disallow: /cgi-bin/
    Disallow: /privatedir/
```

This is interesting - here we declare that crawlers in general should not crawl any parts of our site, EXCEPT for Google, which is allowed to crawl the entire site apart from `/cgi-bin/` and `/privatedir/`. So the rules of specificity apply, not inheritance.

### Allow and Disallow restrictions - Method 1

There is a way to use Disallow: to essentially turn it into "Allow all", and that is by not entering a value after the semicolon(:):

```bash
    User-agent: *
    Disallow: /
    User-agent: ia_archiver
    Disallow:
```

Here all crawlers should be prohibited from crawling our site, except for Alexa, which is allowed.

### Allow and Disallow restrictions - Method 2

Finally, some crawlers now support an additional field called "Allow:", most notably, Google. As its name implies, "Allow:" lets you explicitly dictate what files/folders can be crawled. However, this field is currently not part of the "robots.txt" protocol, so use it only if absolutely needed, as it might confuse some less intelligent crawlers.

Per Google's FAQs for webmasters, the below is the preferred way to disallow all crawlers from your site EXCEPT Google:

```bash
    User-agent: *
    Disallow: /
    User-agent: Googlebot
    Allow: /
```

Finally this file (robot.txt) must be uploaded to the root accessible directory of your site, not a subdirectory (eg. www.mysite.com/robot.txt) it is only by following the above rules will search engines interpret the instructions contained in the file.

## Resources
- http://linuxpoison.blogspot.com/2009/02/how-to-create-and-configure-robottxt.html
